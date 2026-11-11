package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	applySpec   string
	applyDryRun bool
	applyForce  bool
)

var ImageApplyCmd = &cobra.Command{
	Use:     "apply",
	Aliases: []string{"declare"},
	Short:   "Apply a fleet spec: OS images per node, VMs per host",
	Long: `Make the nodes listed in a HostImageSpec match it. Nodes not in the spec
are never touched; removing a node from the spec stops managing it.

  spec.images   named OS images (ref, optional digest: a tag is resolved)
  spec.nodes    node id, optional image, and for VM hosts the complete VM
                set (vms: name, disk, optional digest/vcpus/memoryMb/image)
  spec.rollout  soak, timeout, poll for the per-group canary gate

Order: VM sets on hosts; claim newly enrolled VMs (their ids are derived
from host id + VM name); OS for plain nodes and VMs; OS for VM hosts last.
Per image group the first node is the canary. A rollback or timeout stops
the run; re-running continues. --dry-run prints the plan; --force re-sends
OS updates whose desired image already matches.`,
	Run: executeImageApply,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.SpecFlag})
	},
}

func executeImageApply(cmd *cobra.Command, args []string) {
	spec, err := loadHostImageSpec(applySpec)
	if err != nil {
		applyFail("parsing spec: %v", err)
	}
	if err := validateHostImageSpec(spec); err != nil {
		applyFail("validating spec: %v", err)
	}
	cfg, err := applyRolloutConfig(spec)
	if err != nil {
		applyFail("validating rollout: %v", err)
	}
	cat, err := fetchCatalog()
	if err != nil {
		applyFail("listing images: %v", err)
	}
	images, vms, err := resolveHostImageSpec(spec, cat)
	if err != nil {
		applyFail("resolving spec: %v", err)
	}
	state, err := fetchApplyState(spec)
	if err != nil {
		applyFail("reading control-plane state: %v", err)
	}
	plan, err := buildApplyPlan(spec, images, vms, state, applyForce)
	if err != nil {
		applyFail("%v", err)
	}
	printApplyPlan(spec, plan)
	if plan.empty() || applyDryRun {
		return
	}
	if err := executeApplyPlan(cfg, plan); err != nil {
		rolloutFail("%v", err)
	}
	rolloutLog("apply complete")
}

func applyRolloutConfig(spec model.HostImageSpec) (rolloutConfig, error) {
	cfg := rolloutConfig{org: spec.Metadata.Org, soak: 2 * time.Minute, timeout: 15 * time.Minute, poll: 10 * time.Second, force: applyForce}
	r := spec.Spec.Rollout
	if r == nil {
		return cfg, nil
	}
	for _, f := range []struct {
		name, value string
		dst         *time.Duration
	}{{"soak", r.Soak, &cfg.soak}, {"timeout", r.Timeout, &cfg.timeout}, {"poll", r.Poll, &cfg.poll}} {
		if f.value == "" {
			continue
		}
		d, err := time.ParseDuration(f.value)
		if err != nil || d < 0 {
			return cfg, fmt.Errorf("rollout.%s %q is not a duration", f.name, f.value)
		}
		*f.dst = d
	}
	return cfg, nil
}

func fetchCatalog() (applyCatalog, error) {
	var resp model.ImageListResponse
	if err := gatewayGet("ListImages", nil, &resp); err != nil {
		return applyCatalog{}, err
	}
	cat := applyCatalog{images: map[string]string{}, artifacts: map[string]string{}}
	for _, img := range resp.Images {
		cat.images[img.Ref] = img.Digest
	}
	for _, a := range resp.Artifacts {
		cat.artifacts[a.Name] = a.Checksum
	}
	return cat, nil
}

func fetchApplyState(spec model.HostImageSpec) (applyState, error) {
	state := applyState{poolNodes: map[string]model.Node{}, vms: map[string][]model.VMDesired{}}
	var err error
	if state.orgNodes, err = orgNodesFor(spec.Metadata.Org); err != nil {
		return state, err
	}
	if state.poolNodes, err = poolNodes(); err != nil {
		return state, err
	}
	for _, n := range spec.Spec.Nodes {
		if n.VMs == nil {
			continue
		}
		if _, inOrg := state.orgNodes[n.ID]; !inOrg {
			continue // buildApplyPlan reports it
		}
		var resp model.VMDesiredGetResponse
		if err := gatewayGet("GetDesiredVms", model.VMDesiredGetRequest{HostID: n.ID}, &resp); err != nil {
			return state, fmt.Errorf("VMs of %s: %w", n.ID, err)
		}
		state.vms[n.ID] = resp.VMs
	}
	return state, nil
}

func poolNodes() (map[string]model.Node, error) {
	var resp model.NodesResponse
	if err := gatewayGet("ListNodePool", nil, &resp); err != nil {
		return nil, err
	}
	nodes := make(map[string]model.Node, len(resp.Nodes))
	for _, node := range resp.Nodes {
		nodes[node.ID] = node
	}
	return nodes, nil
}

func executeApplyPlan(cfg rolloutConfig, plan applyPlan) error {
	for _, set := range plan.vmSets {
		rolloutLog("VMs on %s: setting %d VM(s)", nodeLabel(set.hostName, set.hostID), len(set.vms))
		var resp model.VMDesiredSetResponse
		if err := gatewaySend("SetDesiredVms", "POST", model.VMDesiredSetRequest{HostID: set.hostID, VMs: set.vms}, &resp); err != nil {
			return fmt.Errorf("VMs on %s: %w", set.hostID, err)
		}
	}
	for _, claim := range plan.claims {
		if err := claimHostVMs(cfg, claim); err != nil {
			return err
		}
	}
	for i, wave := range plan.waves {
		for _, g := range wave {
			ids := make([]string, len(g.targets))
			for j, t := range g.targets {
				ids[j] = t.id
			}
			rolloutLog("wave %d, image %s (%s): canary %s, then %d more", i+1, g.imageName, shortDigest(g.image.Digest), ids[0], len(ids)-1)
			c := cfg
			c.ref, c.digest = g.image.Ref, g.image.Digest
			if err := rolloutGroup(c, ids[0], ids[1:]); err != nil {
				return fmt.Errorf("wave %d, image %s: %w", i+1, g.imageName, err)
			}
		}
	}
	return nil
}

// claimHostVMs waits until a host's VMs enrolled (they register with their
// derived ids), then claims them by the vm-host label their seed sets.
func claimHostVMs(cfg rolloutConfig, claim claimStep) error {
	host := nodeLabel(claim.hostName, claim.hostID)
	rolloutLog("VMs on %s: waiting for %s to enroll", host, strings.Join(claim.vmNames, ", "))
	deadline := time.Now().Add(cfg.timeout)
	for {
		org, err := orgNodesFor(cfg.org)
		if err != nil {
			return err
		}
		pool, err := poolNodes()
		if err != nil {
			return err
		}
		var waiting []string
		inPool := false
		for i, id := range claim.vmIDs {
			if _, ok := org[id]; ok {
				continue
			}
			if _, ok := pool[id]; ok {
				inPool = true
				continue
			}
			waiting = append(waiting, claim.vmNames[i])
		}
		if len(waiting) == 0 {
			if !inPool {
				return nil
			}
			break
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("VMs on %s did not enroll within %s: %s", host, cfg.timeout, strings.Join(waiting, ", "))
		}
		time.Sleep(cfg.poll)
	}
	queries, err := utils.CreateNodeQuery("vm-host = " + claim.hostID)
	if err != nil {
		return err
	}
	var resp model.ClaimNodesResponse
	if err := gatewaySend("ClaimOwnership", "PATCH", model.ClaimNodesRequest{Org: cfg.org, Query: queries}, &resp); err != nil {
		return fmt.Errorf("claiming VMs of %s: %w", host, err)
	}
	rolloutLog("VMs on %s: claimed %d node(s) into %s", host, len(resp.Nodes), cfg.org)
	return nil
}

// gatewaySend is gatewayGet for writes (POST/PATCH), with the same
// rate-limit retry.
func gatewaySend(method, verb string, body, response interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("reading token: %w", err)
	}
	for attempt := 0; ; attempt++ {
		err = utils.SendHTTPRequest(model.HTTPRequestConfig{
			URL:         clients.BuildURL("core", "v1", method),
			Method:      verb,
			Token:       token,
			Headers:     map[string]string{"Content-Type": "application/json"},
			RequestBody: body,
			Response:    response,
			Timeout:     30 * time.Second,
		})
		if !isRateLimited(err) || attempt >= rateLimitRetries {
			return err
		}
		time.Sleep(rateLimitBackoff)
	}
}

// desiredMatches reports whether the node's declared desired image already
// points at the source. The digest is the immutable unit; a matching ref or
// policy is only required when the source carries one.
func desiredMatches(desired *model.DesiredOsImage, source model.HostImageSource) bool {
	if desired == nil || desired.Digest != source.Digest {
		return false
	}
	if source.Ref != "" && desired.ImageRef != "" && desired.ImageRef != source.Ref {
		return false
	}
	if source.RebootPolicy != "" && desired.RebootPolicy != "" && desired.RebootPolicy != source.RebootPolicy {
		return false
	}
	return true
}

// applyImageUpdate retries while the gateway rate-limits the caller.
func applyImageUpdate(nodeID, ref, digest, policy string) (bool, error) {
	for attempt := 0; ; attempt++ {
		accepted, err := requestImageUpdate(nodeID, ref, digest, policy)
		if isRateLimited(err) && attempt < rateLimitRetries {
			time.Sleep(rateLimitBackoff)
			continue
		}
		return accepted, err
	}
}

// registryFromRef extracts the registry from an image reference when the
// first component looks like a host (a dot, a port, or localhost).
func registryFromRef(ref string) string {
	i := strings.Index(ref, "/")
	if i <= 0 {
		return ""
	}
	host := ref[:i]
	if strings.Contains(host, ".") || strings.Contains(host, ":") || host == "localhost" {
		return host
	}
	return ""
}

// repoFromRef extracts the repository without the registry and tag.
func repoFromRef(ref string) string {
	host := registryFromRef(ref)
	repo := ref
	if host != "" {
		repo = strings.TrimPrefix(ref, host+"/")
	}
	if i := strings.LastIndex(repo, ":"); i > strings.LastIndex(repo, "/") {
		repo = repo[:i]
	}
	return repo
}

// tagFromRef extracts the tag when the last colon is after the last slash.
func tagFromRef(ref string) string {
	i := strings.LastIndex(ref, ":")
	if i > strings.LastIndex(ref, "/") {
		return ref[i+1:]
	}
	return ""
}

func printApplyPlan(spec model.HostImageSpec, plan applyPlan) {
	fmt.Printf("HostImageSpec org=%s: %d node(s) managed\n", spec.Metadata.Org, len(spec.Spec.Nodes))
	if plan.empty() {
		fmt.Println("Up to date: nothing to do.")
		return
	}
	for _, set := range plan.vmSets {
		fmt.Printf("\nVMs on %s:\n", nodeLabel(set.hostName, set.hostID))
		for _, name := range set.added {
			fmt.Printf("  + %s\n", name)
		}
		for _, name := range set.changed {
			fmt.Printf("  ~ %s (changes apply when the VM is re-created)\n", name)
		}
		for _, name := range set.removed {
			fmt.Printf("  - %s (deleted)\n", name)
		}
	}
	for _, claim := range plan.claims {
		fmt.Printf("\nClaim into %s after they enroll (VMs of %s):\n", spec.Metadata.Org, nodeLabel(claim.hostName, claim.hostID))
		for i, id := range claim.vmIDs {
			fmt.Printf("  %s  %s\n", claim.vmNames[i], id)
		}
	}
	titles := [2]string{"OS wave 1 (nodes and VMs)", "OS wave 2 (VM hosts)"}
	for i, wave := range plan.waves {
		if len(wave) == 0 {
			continue
		}
		fmt.Printf("\n%s:\n", titles[i])
		for _, g := range wave {
			fmt.Printf("  image %s = %s@%s\n", g.imageName, g.image.Ref, shortDigest(g.image.Digest))
			for j, t := range g.targets {
				role := ""
				if j == 0 {
					role = "  [canary]"
				}
				fmt.Printf("    %-10s %s  %s -> %s  %s%s\n", orDash(t.name), t.id, describeDesired(t.current), shortDigest(g.image.Digest), describeBooted(t.status), role)
			}
		}
	}
	ids := map[string]bool{}
	for _, t := range plan.targets() {
		ids[t.id] = true
	}
	list := make([]string, 0, len(ids))
	for id := range ids {
		list = append(list, id)
	}
	sort.Strings(list)
	fmt.Printf("\nNeeds os.update policies (then cockpit login) for: %s\n", strings.Join(list, ", "))
	if applyDryRun {
		fmt.Println("Dry run: nothing was changed.")
	}
}

func describeDesired(d *model.DesiredOsImage) string {
	if d == nil || d.Digest == "" {
		return "(none)"
	}
	return shortDigest(d.Digest)
}

func describeBooted(st *model.OsStatus) string {
	if st == nil {
		return "booted: -"
	}
	return fmt.Sprintf("booted: %s %s", orDash(st.BootedVersion), orDash(st.Phase))
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func applyFail(format string, args ...interface{}) {
	fmt.Printf("Error: "+format+"\n", args...)
	os.Exit(1)
}

func init() {
	ImageCmd.AddCommand(ImageApplyCmd)
	f := ImageApplyCmd.Flags()
	f.StringVar(&applySpec, constants.SpecFlag, "", constants.HostImageSpecDescription)
	f.BoolVar(&applyDryRun, "dry-run", false, "Print the plan and exit without changing anything")
	f.BoolVar(&applyForce, "force", false, "Re-send OS updates even when the desired image already matches")
	ImageApplyCmd.MarkFlagRequired(constants.SpecFlag)
}
