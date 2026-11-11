package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// This file is the pure part of `image apply`: loading and validating the
// spec, resolving tags and disk names, and diffing the spec against the
// control plane's state into a plan. Nothing here talks to the gateway.

const runVmsCapability = "run_vms"

// vmNodeID is a VM's node id: UUIDv5(host id, VM name). star derives the same
// id when it creates the VM (star internal/vm VMNodeID), so VM ids are known
// before the VMs exist.
func vmNodeID(hostID, vmName string) (string, error) {
	ns, err := uuid.Parse(hostID)
	if err != nil {
		return "", fmt.Errorf("host id %q is not a UUID: %w", hostID, err)
	}
	return uuid.NewSHA1(ns, []byte(vmName)).String(), nil
}

// loadHostImageSpec reads a path or an inline JSON document (JSON is valid
// YAML). Unknown fields are errors, so a misspelled key — or the removed
// single-image form — fails instead of being ignored.
func loadHostImageSpec(raw string) (model.HostImageSpec, error) {
	var spec model.HostImageSpec
	data := []byte(raw)
	if !strings.Contains(raw, "\n") && !strings.HasPrefix(strings.TrimSpace(raw), "{") {
		b, err := os.ReadFile(raw)
		if err != nil {
			return spec, err
		}
		data = b
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&spec); err != nil {
		if strings.Contains(err.Error(), "not found in type") {
			return spec, fmt.Errorf("%w (spec takes images, nodes and rollout; the single-image form with image/build/target was removed)", err)
		}
		return spec, err
	}
	return spec, nil
}

func validateHostImageSpec(spec model.HostImageSpec) error {
	if spec.Kind != "" && spec.Kind != "HostImageSpec" {
		return fmt.Errorf("kind %q: want HostImageSpec", spec.Kind)
	}
	if spec.Metadata.Org == "" {
		return errors.New("metadata.org is required")
	}
	for name, img := range spec.Spec.Images {
		if img.Ref == "" {
			return fmt.Errorf("image %s: ref is required", name)
		}
		if img.Digest != "" {
			if err := utils.ValidateDigest(img.Digest); err != nil {
				return fmt.Errorf("image %s: %w", name, err)
			}
		}
	}
	if len(spec.Spec.Nodes) == 0 {
		return errors.New("spec.nodes is empty")
	}
	imageKnown := func(name string) bool { _, ok := spec.Spec.Images[name]; return ok }
	seen := map[string]bool{}
	for i, n := range spec.Spec.Nodes {
		label := nodeLabel(n.Name, n.ID)
		if _, err := uuid.Parse(n.ID); err != nil {
			return fmt.Errorf("nodes[%d]: id %q is not a node UUID", i, n.ID)
		}
		if seen[n.ID] {
			return fmt.Errorf("nodes[%d]: duplicate node id %s", i, n.ID)
		}
		seen[n.ID] = true
		if n.Image != "" && !imageKnown(n.Image) {
			return fmt.Errorf("node %s: image %q is not in spec.images", label, n.Image)
		}
		if n.VMs == nil {
			continue
		}
		names := map[string]bool{}
		for _, vm := range *n.VMs {
			if vm.Name == "" {
				return fmt.Errorf("node %s: a VM has no name", label)
			}
			if names[vm.Name] {
				return fmt.Errorf("node %s: duplicate VM %s", label, vm.Name)
			}
			names[vm.Name] = true
			if vm.Disk == "" {
				return fmt.Errorf("VM %s on %s: disk is required", vm.Name, label)
			}
			if vm.Digest != "" {
				if err := utils.ValidateDigest(vm.Digest); err != nil {
					return fmt.Errorf("VM %s on %s: %w", vm.Name, label, err)
				}
			}
			if vm.Image != "" && !imageKnown(vm.Image) {
				return fmt.Errorf("VM %s on %s: image %q is not in spec.images", vm.Name, label, vm.Image)
			}
			id, _ := vmNodeID(n.ID, vm.Name)
			if seen[id] {
				return fmt.Errorf("VM %s on %s: its id %s is also listed as a node", vm.Name, label, id)
			}
			seen[id] = true
		}
	}
	return nil
}

// applyCatalog is what the registry and the disk registry hold: "repo:tag" ->
// manifest digest, and qcow2 name -> sha256 (from `image list`).
type applyCatalog struct {
	images    map[string]string
	artifacts map[string]string
}

// resolveHostImageSpec pins every image to a digest and every VM disk to a
// checksum. A tag without a digest is looked up in the registry; a disk name
// without a digest in the disk registry. A full disk URL needs its digest in
// the file.
func resolveHostImageSpec(spec model.HostImageSpec, cat applyCatalog) (map[string]model.HostImageSource, map[string][]model.VMDesired, error) {
	images := make(map[string]model.HostImageSource, len(spec.Spec.Images))
	for name, img := range spec.Spec.Images {
		if img.Digest == "" {
			key := repoFromRef(img.Ref) + ":" + tagFromRef(img.Ref)
			if tagFromRef(img.Ref) == "" {
				return nil, nil, fmt.Errorf("image %s: %s has neither a tag nor a digest", name, img.Ref)
			}
			digest, ok := cat.images[key]
			if !ok {
				return nil, nil, fmt.Errorf("image %s: tag %s not found in the registry", name, img.Ref)
			}
			img.Digest = digest
		}
		images[name] = img
	}
	vms := map[string][]model.VMDesired{}
	for _, n := range spec.Spec.Nodes {
		if n.VMs == nil {
			continue
		}
		set := make([]model.VMDesired, 0, len(*n.VMs))
		for _, vm := range *n.VMs {
			digest := vm.Digest
			if digest == "" {
				if strings.Contains(vm.Disk, "://") {
					return nil, nil, fmt.Errorf("VM %s: disk %s is a URL; give its digest", vm.Name, vm.Disk)
				}
				sum, ok := cat.artifacts[vm.Disk]
				if !ok {
					return nil, nil, fmt.Errorf("VM %s: disk %s not found in the disk registry", vm.Name, vm.Disk)
				}
				digest = sum
			}
			set = append(set, model.VMDesired{Name: vm.Name, Qcow2URL: vm.Disk, Digest: digest, VCPUs: vm.VCPUs, MemoryMB: vm.MemoryMB})
		}
		vms[n.ID] = set
	}
	return images, vms, nil
}

// applyState is the control plane's view: org and pool nodes, and the stored
// VM set of every host the spec declares VMs for.
type applyState struct {
	orgNodes  map[string]model.Node
	poolNodes map[string]model.Node
	vms       map[string][]model.VMDesired
}

type vmSetChange struct {
	hostID, hostName        string
	vms                     []model.VMDesired
	added, removed, changed []string
}

// claimStep claims the VMs of one host into the org once they enrolled.
type claimStep struct {
	hostID, hostName string
	vmIDs, vmNames   []string
}

type osTarget struct {
	id, name string
	current  *model.DesiredOsImage
	status   *model.OsStatus
}

type osGroup struct {
	imageName string
	image     model.HostImageSource
	targets   []osTarget
}

// applyPlan is what apply will do, in order: VM sets, claims, OS wave 1
// (plain nodes and VMs), OS wave 2 (VM hosts).
type applyPlan struct {
	vmSets []vmSetChange
	claims []claimStep
	waves  [2][]osGroup
}

func (p applyPlan) empty() bool {
	return len(p.vmSets) == 0 && len(p.claims) == 0 && len(p.waves[0]) == 0 && len(p.waves[1]) == 0
}

// targets lists every node the plan sends an OS update to.
func (p applyPlan) targets() []osTarget {
	var out []osTarget
	for _, wave := range p.waves {
		for _, g := range wave {
			out = append(out, g.targets...)
		}
	}
	return out
}

func buildApplyPlan(spec model.HostImageSpec, images map[string]model.HostImageSource, vms map[string][]model.VMDesired, state applyState, force bool) (applyPlan, error) {
	var plan applyPlan
	var waves [2]groupBuilder

	for _, n := range spec.Spec.Nodes {
		label := nodeLabel(n.Name, n.ID)
		node, ok := state.orgNodes[n.ID]
		if !ok {
			return plan, fmt.Errorf("node %s is not in org %s (claim it first)", label, spec.Metadata.Org)
		}
		isHost := n.VMs != nil
		if n.Image != "" {
			if !hasCapability(node, osUpdateCapability) {
				return plan, fmt.Errorf("node %s has no %s capability", label, osUpdateCapability)
			}
			wave := 0
			if isHost {
				wave = 1
			}
			img := images[n.Image]
			if force || !desiredMatches(node.DesiredOs, img) {
				waves[wave].add(n.Image, img, osTarget{id: n.ID, name: n.Name, current: node.DesiredOs, status: node.OsStatus})
			}
		}
		if !isHost {
			continue
		}
		if !hasCapability(node, runVmsCapability) {
			return plan, fmt.Errorf("node %s declares VMs but has no %s capability", label, runVmsCapability)
		}
		if change, differs := diffVMSet(state.vms[n.ID], vms[n.ID]); differs {
			change.hostID, change.hostName = n.ID, n.Name
			plan.vmSets = append(plan.vmSets, change)
		}
		claim := claimStep{hostID: n.ID, hostName: n.Name}
		for _, vm := range *n.VMs {
			if vm.Image == "" {
				continue
			}
			id, err := vmNodeID(n.ID, vm.Name)
			if err != nil {
				return plan, err
			}
			vmNode, inOrg := state.orgNodes[id]
			if !inOrg {
				claim.vmIDs = append(claim.vmIDs, id)
				claim.vmNames = append(claim.vmNames, vm.Name)
			}
			img := images[vm.Image]
			if force || !inOrg || !desiredMatches(vmNode.DesiredOs, img) {
				waves[0].add(vm.Image, img, osTarget{id: id, name: vm.Name, current: vmNode.DesiredOs, status: vmNode.OsStatus})
			}
		}
		if len(claim.vmIDs) > 0 {
			plan.claims = append(plan.claims, claim)
		}
	}
	plan.waves[0], plan.waves[1] = waves[0].groups, waves[1].groups
	return plan, nil
}

// groupBuilder collects targets per image, in order of first appearance.
type groupBuilder struct {
	groups []osGroup
}

func (b *groupBuilder) add(name string, img model.HostImageSource, target osTarget) {
	for i := range b.groups {
		if b.groups[i].imageName == name {
			b.groups[i].targets = append(b.groups[i].targets, target)
			return
		}
	}
	b.groups = append(b.groups, osGroup{imageName: name, image: img, targets: []osTarget{target}})
}

// diffVMSet compares a host's stored VM set with the declared one by name.
func diffVMSet(current, declared []model.VMDesired) (vmSetChange, bool) {
	change := vmSetChange{vms: declared}
	byName := make(map[string]model.VMDesired, len(current))
	for _, vm := range current {
		byName[vm.Name] = vm
	}
	for _, vm := range declared {
		old, ok := byName[vm.Name]
		switch {
		case !ok:
			change.added = append(change.added, vm.Name)
		case old != vm:
			change.changed = append(change.changed, vm.Name)
		}
		delete(byName, vm.Name)
	}
	for name := range byName {
		change.removed = append(change.removed, name)
	}
	sort.Strings(change.removed)
	differs := len(change.added)+len(change.removed)+len(change.changed) > 0
	return change, differs
}

func nodeLabel(name, id string) string {
	if name == "" {
		return id
	}
	return name + " (" + id + ")"
}
