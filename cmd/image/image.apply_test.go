package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/c12s/cockpit/model"
)

var validDigest = "sha256:" + strings.Repeat("a", 64)

const (
	hostID  = "5c1e7d2a-1b3c-4d5e-8f90-123456789abc" // VM host
	plainID = "7f3c0000-0000-4000-8000-000000000001" // plain node
	vm01ID  = "f210cdd5-ce6a-561e-a06a-29929fc319bd" // UUIDv5(hostID, "vm-01")
)

var (
	digestA = "sha256:" + strings.Repeat("a", 64)
	digestB = "sha256:" + strings.Repeat("b", 64)
	diskSum = "sha256:" + strings.Repeat("d", 64)
)

const fleetYAML = `
apiVersion: c12s/os.v1
kind: HostImageSpec
metadata:
  org: c12s
spec:
  images:
    stable:
      ref: 192.168.122.1:5001/os/edge-base:v4.1
  nodes:
    - id: ` + plainID + `
      name: host-1
      image: stable
    - id: ` + hostID + `
      name: node-1
      image: stable
      vms:
        - name: vm-01
          disk: edge-base-v4.0.qcow2
          vcpus: 1
          memoryMb: 1024
          image: stable
  rollout:
    soak: 1m
`

// Same vector as star's TestVMNodeID and Python's uuid.uuid5.
func TestVMNodeID(t *testing.T) {
	got, err := vmNodeID(hostID, "vm-01")
	if err != nil || got != vm01ID {
		t.Fatalf("vmNodeID = %q, %v; want %s", got, err, vm01ID)
	}
	if _, err := vmNodeID("nope", "vm-01"); err == nil {
		t.Fatal("malformed host id must fail")
	}
}

func TestLoadHostImageSpec(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fleet.yaml")
	if err := os.WriteFile(path, []byte(fleetYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	spec, err := loadHostImageSpec(path)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Metadata.Org != "c12s" || len(spec.Spec.Nodes) != 2 || spec.Spec.Images["stable"].Ref == "" {
		t.Fatalf("unexpected spec %+v", spec)
	}
	if spec.Spec.Nodes[0].VMs != nil {
		t.Fatal("a node without vms: must leave VMs untouched (nil)")
	}
	if vms := spec.Spec.Nodes[1].VMs; vms == nil || len(*vms) != 1 {
		t.Fatal("vms list not parsed")
	}

	empty, err := loadHostImageSpec(`{"kind":"HostImageSpec","metadata":{"org":"c12s"},"spec":{"nodes":[{"id":"` + hostID + `","vms":[]}]}}`)
	if err != nil {
		t.Fatal(err)
	}
	if vms := empty.Spec.Nodes[0].VMs; vms == nil || len(*vms) != 0 {
		t.Fatal("vms: [] must mean an empty set, not absent")
	}

	_, err = loadHostImageSpec(`{"kind":"HostImageSpec","metadata":{"org":"c12s"},"spec":{"image":{"ref":"x","digest":"` + digestA + `"}}}`)
	if err == nil || !strings.Contains(err.Error(), "images") {
		t.Fatalf("old single-image form must be rejected with a hint, got %v", err)
	}
}

func TestValidateHostImageSpec(t *testing.T) {
	base := func() model.HostImageSpec {
		spec, err := loadHostImageSpec(fleetYAML)
		if err != nil {
			t.Fatal(err)
		}
		return spec
	}
	cases := []struct {
		name   string
		mutate func(*model.HostImageSpec)
		want   string
	}{
		{"valid", func(*model.HostImageSpec) {}, ""},
		{"no org", func(s *model.HostImageSpec) { s.Metadata.Org = "" }, "metadata.org"},
		{"duplicate node", func(s *model.HostImageSpec) { s.Spec.Nodes[1].ID = plainID }, "duplicate"},
		{"unknown image", func(s *model.HostImageSpec) { s.Spec.Nodes[0].Image = "nope" }, "nope"},
		{"bad digest", func(s *model.HostImageSpec) {
			s.Spec.Images["stable"] = model.HostImageSource{Ref: "r:t", Digest: "sha256:abc"}
		}, "digest"},
		{"duplicate vm", func(s *model.HostImageSpec) {
			vms := append(*s.Spec.Nodes[1].VMs, (*s.Spec.Nodes[1].VMs)[0])
			s.Spec.Nodes[1].VMs = &vms
		}, "vm-01"},
		{"vm without disk", func(s *model.HostImageSpec) { (*s.Spec.Nodes[1].VMs)[0].Disk = "" }, "disk"},
		{"vm unknown image", func(s *model.HostImageSpec) { (*s.Spec.Nodes[1].VMs)[0].Image = "nope" }, "nope"},
		{"malformed id", func(s *model.HostImageSpec) { s.Spec.Nodes[1].ID = "node-1" }, "node-1"},
	}
	for _, c := range cases {
		spec := base()
		c.mutate(&spec)
		err := validateHostImageSpec(spec)
		if c.want == "" && err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
		}
		if c.want != "" && (err == nil || !strings.Contains(err.Error(), c.want)) {
			t.Errorf("%s: want error containing %q, got %v", c.name, c.want, err)
		}
	}
}

func testCatalog() applyCatalog {
	return applyCatalog{
		images:    map[string]string{"os/edge-base:v4.1": digestB},
		artifacts: map[string]string{"edge-base-v4.0.qcow2": diskSum},
	}
}

func TestResolveHostImageSpec(t *testing.T) {
	spec, _ := loadHostImageSpec(fleetYAML)
	images, vms, err := resolveHostImageSpec(spec, testCatalog())
	if err != nil {
		t.Fatal(err)
	}
	if images["stable"].Digest != digestB {
		t.Fatalf("tag not resolved: %+v", images["stable"])
	}
	if got := vms[hostID][0]; got.Digest != diskSum || got.Qcow2URL != "edge-base-v4.0.qcow2" {
		t.Fatalf("disk not resolved: %+v", got)
	}

	missing := testCatalog()
	delete(missing.images, "os/edge-base:v4.1")
	if _, _, err := resolveHostImageSpec(spec, missing); err == nil || !strings.Contains(err.Error(), "v4.1") {
		t.Fatalf("unknown tag must fail, got %v", err)
	}
	noDisk := testCatalog()
	delete(noDisk.artifacts, "edge-base-v4.0.qcow2")
	if _, _, err := resolveHostImageSpec(spec, noDisk); err == nil || !strings.Contains(err.Error(), "edge-base-v4.0.qcow2") {
		t.Fatalf("unknown disk must fail, got %v", err)
	}
}

func osNode(id, digest string) model.Node {
	n := model.Node{ID: id, Capabilities: []string{osUpdateCapability, "run_vms"}}
	if digest != "" {
		n.DesiredOs = &model.DesiredOsImage{ImageRef: "192.168.122.1:5001/os/edge-base:v4.1", Digest: digest}
	}
	return n
}

func planFor(t *testing.T, state applyState) applyPlan {
	t.Helper()
	spec, _ := loadHostImageSpec(fleetYAML)
	images, vms, err := resolveHostImageSpec(spec, testCatalog())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := buildApplyPlan(spec, images, vms, state, false)
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func TestBuildApplyPlanFreshFleet(t *testing.T) {
	plan := planFor(t, applyState{
		orgNodes: map[string]model.Node{plainID: osNode(plainID, digestA), hostID: osNode(hostID, digestA)},
		vms:      map[string][]model.VMDesired{},
	})
	if len(plan.vmSets) != 1 || plan.vmSets[0].hostID != hostID || strings.Join(plan.vmSets[0].added, ",") != "vm-01" {
		t.Fatalf("vm set: %+v", plan.vmSets)
	}
	if len(plan.claims) != 1 || plan.claims[0].hostID != hostID || plan.claims[0].vmIDs[0] != vm01ID {
		t.Fatalf("claims: %+v", plan.claims)
	}
	if ids := targetIDs(plan.waves[0]); ids != plainID+","+vm01ID {
		t.Fatalf("wave 1 = %s", ids)
	}
	if ids := targetIDs(plan.waves[1]); ids != hostID {
		t.Fatalf("wave 2 (hosts) = %s", ids)
	}
}

func TestBuildApplyPlanConverged(t *testing.T) {
	vm := osNode(vm01ID, digestB)
	plan := planFor(t, applyState{
		orgNodes: map[string]model.Node{plainID: osNode(plainID, digestB), hostID: osNode(hostID, digestB), vm01ID: vm},
		vms: map[string][]model.VMDesired{hostID: {{
			Name: "vm-01", Qcow2URL: "edge-base-v4.0.qcow2", Digest: diskSum, VCPUs: 1, MemoryMB: 1024,
		}}},
	})
	if !plan.empty() {
		t.Fatalf("converged fleet must plan nothing: %+v", plan)
	}
}

func TestBuildApplyPlanVMChangesAndPool(t *testing.T) {
	plan := planFor(t, applyState{
		orgNodes:  map[string]model.Node{plainID: osNode(plainID, digestB), hostID: osNode(hostID, digestB)},
		poolNodes: map[string]model.Node{vm01ID: {ID: vm01ID}},
		vms: map[string][]model.VMDesired{hostID: {
			{Name: "vm-01", Qcow2URL: "edge-base-v4.0.qcow2", Digest: diskSum, VCPUs: 2, MemoryMB: 1024},
			{Name: "vm-03", Qcow2URL: "edge-base-v4.0.qcow2", Digest: diskSum},
		}},
	})
	set := plan.vmSets[0]
	if strings.Join(set.changed, ",") != "vm-01" || strings.Join(set.removed, ",") != "vm-03" || len(set.added) != 0 {
		t.Fatalf("vm diff: %+v", set)
	}
	if len(plan.claims) != 1 {
		t.Fatalf("pool VM must be claimed: %+v", plan.claims)
	}
	if ids := targetIDs(plan.waves[0]); ids != vm01ID {
		t.Fatalf("only the VM needs an OS update: %s", ids)
	}
	if len(plan.waves[1]) != 0 {
		t.Fatalf("host already on the image: %+v", plan.waves[1])
	}
}

func TestBuildApplyPlanRejectsNodeOutsideOrg(t *testing.T) {
	spec, _ := loadHostImageSpec(fleetYAML)
	images, vms, _ := resolveHostImageSpec(spec, testCatalog())
	_, err := buildApplyPlan(spec, images, vms, applyState{orgNodes: map[string]model.Node{hostID: osNode(hostID, "")}}, false)
	if err == nil || !strings.Contains(err.Error(), plainID) {
		t.Fatalf("node outside the org must fail, got %v", err)
	}
}

func TestBuildApplyPlanForce(t *testing.T) {
	spec, _ := loadHostImageSpec(fleetYAML)
	images, vms, _ := resolveHostImageSpec(spec, testCatalog())
	state := applyState{
		orgNodes: map[string]model.Node{plainID: osNode(plainID, digestB), hostID: osNode(hostID, digestB), vm01ID: osNode(vm01ID, digestB)},
		vms: map[string][]model.VMDesired{hostID: {{
			Name: "vm-01", Qcow2URL: "edge-base-v4.0.qcow2", Digest: diskSum, VCPUs: 1, MemoryMB: 1024,
		}}},
	}
	plan, err := buildApplyPlan(spec, images, vms, state, true)
	if err != nil {
		t.Fatal(err)
	}
	if ids := targetIDs(plan.waves[0]); ids != plainID+","+vm01ID {
		t.Fatalf("--force must re-send matching nodes: %s", ids)
	}
}

func targetIDs(groups []osGroup) string {
	var ids []string
	for _, g := range groups {
		for _, target := range g.targets {
			ids = append(ids, target.id)
		}
	}
	return strings.Join(ids, ",")
}

func TestImageRefParts(t *testing.T) {
	cases := []struct {
		ref      string
		registry string
		repo     string
		tag      string
	}{
		{"192.168.122.1:5001/os/edge-base:demo-2", "192.168.122.1:5001", "os/edge-base", "demo-2"},
		{"192.168.122.1:5001/os/edge-base", "192.168.122.1:5001", "os/edge-base", ""},
		{"os/edge-base:demo-2", "", "os/edge-base", "demo-2"},
		{"edge-base:demo-2", "", "edge-base", "demo-2"},
		{"localhost:5001/os/edge-base:demo-2", "localhost:5001", "os/edge-base", "demo-2"},
		{"edge-base", "", "edge-base", ""},
	}
	for _, c := range cases {
		if got := registryFromRef(c.ref); got != c.registry {
			t.Errorf("registryFromRef(%q) = %q, want %q", c.ref, got, c.registry)
		}
		if got := repoFromRef(c.ref); got != c.repo {
			t.Errorf("repoFromRef(%q) = %q, want %q", c.ref, got, c.repo)
		}
		if got := tagFromRef(c.ref); got != c.tag {
			t.Errorf("tagFromRef(%q) = %q, want %q", c.ref, got, c.tag)
		}
	}
}

func TestDesiredMatches(t *testing.T) {
	source := model.HostImageSource{Ref: "reg/os/edge-base:demo-2", Digest: validDigest, RebootPolicy: "immediate"}
	if desiredMatches(nil, source) {
		t.Fatal("nil desired should not match")
	}
	if !desiredMatches(&model.DesiredOsImage{ImageRef: source.Ref, Digest: source.Digest, RebootPolicy: source.RebootPolicy}, source) {
		t.Fatal("matching desired should match")
	}
	if desiredMatches(&model.DesiredOsImage{ImageRef: source.Ref, Digest: source.Digest, RebootPolicy: "manual"}, source) {
		t.Fatal("a different policy should not match")
	}
	if desiredMatches(&model.DesiredOsImage{ImageRef: "other", Digest: source.Digest, RebootPolicy: source.RebootPolicy}, source) {
		t.Fatal("a different ref should not match")
	}
	if desiredMatches(&model.DesiredOsImage{ImageRef: source.Ref, Digest: "sha256:other", RebootPolicy: source.RebootPolicy}, source) {
		t.Fatal("a different digest should not match")
	}
	emptyPolicy := model.HostImageSource{Ref: source.Ref, Digest: source.Digest}
	if !desiredMatches(&model.DesiredOsImage{ImageRef: source.Ref, Digest: source.Digest, RebootPolicy: "manual"}, emptyPolicy) {
		t.Fatal("an empty source policy should not constrain the desired policy")
	}
}
