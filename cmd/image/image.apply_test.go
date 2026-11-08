package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/c12s/cockpit/model"
)

var validDigest = "sha256:" + strings.Repeat("a", 64)

func TestLoadHostImageSpec(t *testing.T) {
	inline := `{"apiVersion":"c12s/os.v1","kind":"HostImageSpec","metadata":{"org":"c12s","name":"edge-base","version":"v1"},"spec":{"image":{"ref":"192.168.122.1:5001/os/edge-base","digest":"` + validDigest + `","rebootPolicy":"immediate"},"target":{"org":"c12s","nodes":["node-1"]}}}`
	spec, err := loadHostImageSpec(inline)
	if err != nil {
		t.Fatalf("loadHostImageSpec(inline): %v", err)
	}
	if spec.Metadata.Org != "c12s" || spec.Metadata.Name != "edge-base" || spec.Metadata.Version != "v1" {
		t.Fatalf("unexpected metadata: %+v", spec.Metadata)
	}
	body := spec.Normalize()
	if body.Image == nil || body.Image.Ref != "192.168.122.1:5001/os/edge-base" || body.Image.Digest != validDigest {
		t.Fatalf("unexpected image: %+v", body.Image)
	}
	if body.Target == nil || len(body.Target.Nodes) != 1 || body.Target.Nodes[0] != "node-1" {
		t.Fatalf("unexpected target: %+v", body.Target)
	}

	path := filepath.Join(t.TempDir(), "host-image-spec.yaml")
	yamlDoc := "apiVersion: c12s/os.v1\nkind: HostImageSpec\nmetadata:\n  org: c12s\nspec:\n  image:\n    ref: reg/os/edge-base\n    digest: " + validDigest + "\n  target:\n    nodes: [n1]\n"
	if err := os.WriteFile(path, []byte(yamlDoc), 0o644); err != nil {
		t.Fatalf("write temp spec: %v", err)
	}
	spec, err = loadHostImageSpec(path)
	if err != nil {
		t.Fatalf("loadHostImageSpec(file): %v", err)
	}
	if spec.Normalize().Image == nil || spec.Normalize().Image.Ref != "reg/os/edge-base" {
		t.Fatalf("file spec did not load: %+v", spec.Normalize())
	}

	if _, err := loadHostImageSpec("not: [valid"); err == nil {
		t.Fatal("expected an error for malformed YAML")
	}

	if _, err := loadHostImageSpec("missing-spec.yaml"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("missing file: expected not found, got %v", err)
	}
}

func TestNormalizeHostImageSpec(t *testing.T) {
	flat := model.HostImageSpec{
		Org:    "c12s",
		Image:  &model.HostImageSource{Ref: "reg/os/edge-base", Digest: validDigest},
		Target: &model.HostImageTarget{Nodes: []string{"n1"}},
	}
	body := flat.Normalize()
	if body.Image == nil || body.Target == nil {
		t.Fatalf("flat form did not merge: %+v", body)
	}

	canonical := model.HostImageSpec{
		Metadata: model.HostImageMetadata{Org: "c12s"},
		Spec: model.HostImageSpecBody{
			Image: &model.HostImageSource{Ref: "canonical", Digest: validDigest},
		},
		Image: &model.HostImageSource{Ref: "flat", Digest: validDigest},
	}
	body = canonical.Normalize()
	if body.Image == nil || body.Image.Ref != "canonical" {
		t.Fatalf("canonical form should win: %+v", body.Image)
	}
}

func TestNormalizeImageSource(t *testing.T) {
	if _, err := normalizeImageSource(nil, nil); err == nil {
		t.Fatal("expected an error when neither image nor build is set")
	}

	image := &model.HostImageSource{Ref: "reg/os/edge-base", Digest: validDigest}
	src, err := normalizeImageSource(image, nil)
	if err != nil {
		t.Fatalf("normalizeImageSource(image): %v", err)
	}
	if src.Ref != image.Ref || src.Digest != validDigest || src.RebootPolicy != "immediate" {
		t.Fatalf("unexpected source: %+v", src)
	}

	build := &model.HostImageBuild{Name: "os/edge-base"}
	src, err = normalizeImageSource(nil, build)
	if err != nil {
		t.Fatalf("normalizeImageSource(build): %v", err)
	}
	if src.Ref != "" || src.Digest != "" {
		t.Fatalf("build-only source should be empty: %+v", src)
	}

	src, err = normalizeImageSource(image, build)
	if err != nil {
		t.Fatalf("normalizeImageSource(image+build): %v", err)
	}
	if src.Ref != image.Ref || src.Digest != validDigest {
		t.Fatalf("image fallback not kept: %+v", src)
	}

	if _, err := normalizeImageSource(&model.HostImageSource{Ref: "reg", Digest: "sha256:abc"}, nil); err == nil {
		t.Fatal("expected an error for a malformed digest")
	}
	if _, err := normalizeImageSource(&model.HostImageSource{Digest: validDigest}, nil); err == nil {
		t.Fatal("expected an error for a missing ref")
	}
	if _, err := normalizeImageSource(&model.HostImageSource{Ref: "reg"}, nil); err == nil {
		t.Fatal("expected an error for a missing digest")
	}
}

func TestNormalizeBuild(t *testing.T) {
	meta := model.HostImageMetadata{Org: "c12s", Name: "edge-base", Version: "v1"}
	image := model.HostImageSource{Ref: "192.168.122.1:5001/os/edge-base:demo-2", Digest: validDigest, RebootPolicy: "immediate"}

	build := &model.HostImageBuild{Context: "tools/bootc", BuildArgs: map[string]string{"OS_VERSION": "v1"}}
	got, err := normalizeBuild(build, meta, image, 30*time.Minute)
	if err != nil {
		t.Fatalf("normalizeBuild: %v", err)
	}
	if got.Name != "os/edge-base" || got.Tag != "demo-2" || got.Registry != "192.168.122.1:5001" {
		t.Fatalf("unexpected build name/tag/registry: %+v", got)
	}
	if got.Containerfile != "Containerfile" || got.Timeout != "30m0s" {
		t.Fatalf("unexpected build defaults: %+v", got)
	}
	if got.BuildArgs["OS_VERSION"] != "v1" {
		t.Fatalf("build args not preserved: %+v", got.BuildArgs)
	}

	got, err = normalizeBuild(&model.HostImageBuild{}, meta, model.HostImageSource{}, time.Minute)
	if err != nil {
		t.Fatalf("normalizeBuild(metadata): %v", err)
	}
	if got.Name != "edge-base" || got.Tag != "v1" || got.Registry != "192.168.122.1:5001" || got.Timeout != "1m0s" {
		t.Fatalf("unexpected metadata defaults: %+v", got)
	}

	if _, err := normalizeBuild(&model.HostImageBuild{Timeout: "nope"}, meta, model.HostImageSource{}, time.Minute); err == nil {
		t.Fatal("expected an error for an invalid build timeout")
	}
	if _, err := normalizeBuild(&model.HostImageBuild{}, model.HostImageMetadata{}, model.HostImageSource{}, time.Minute); err == nil {
		t.Fatal("expected an error when the build name cannot be resolved")
	}
}

func TestNormalizeTarget(t *testing.T) {
	meta := model.HostImageMetadata{Org: "c12s"}
	if _, err := normalizeTarget(nil, meta); err == nil {
		t.Fatal("expected an error for a missing target")
	}

	target := &model.HostImageTarget{Nodes: []string{"n1", " n1 ", "", "n2"}}
	got, err := normalizeTarget(target, meta)
	if err != nil {
		t.Fatalf("normalizeTarget(nodes): %v", err)
	}
	if got.Org != "c12s" || len(got.Nodes) != 2 || got.Nodes[0] != "n1" || got.Nodes[1] != "n2" {
		t.Fatalf("unexpected target: %+v", got)
	}

	target = &model.HostImageTarget{Query: "role = prod", Org: "other"}
	got, err = normalizeTarget(target, meta)
	if err != nil {
		t.Fatalf("normalizeTarget(query): %v", err)
	}
	if got.Org != "other" {
		t.Fatalf("target org should win: %+v", got)
	}

	if _, err := normalizeTarget(&model.HostImageTarget{Nodes: []string{"n1"}, Query: "role = prod"}, meta); err == nil {
		t.Fatal("expected an error when both nodes and query are set")
	}
	if _, err := normalizeTarget(&model.HostImageTarget{}, meta); err == nil {
		t.Fatal("expected an error when neither nodes nor query is set")
	}
	if _, err := normalizeTarget(&model.HostImageTarget{Nodes: []string{" "}}, meta); err == nil {
		t.Fatal("expected an error for an empty nodes list")
	}
	if _, err := normalizeTarget(&model.HostImageTarget{Nodes: []string{"n1"}}, model.HostImageMetadata{}); err == nil {
		t.Fatal("expected an error when the org cannot be resolved")
	}
}

func TestNormalizeRollout(t *testing.T) {
	got, err := normalizeRollout(nil)
	if err != nil || got != nil {
		t.Fatalf("nil rollout should return nil, nil; got %v, %v", got, err)
	}

	if _, err := normalizeRollout(&model.HostImageRollout{}); err == nil {
		t.Fatal("expected an error for a missing canary")
	}

	got, err = normalizeRollout(&model.HostImageRollout{Canary: "n1"})
	if err != nil {
		t.Fatalf("normalizeRollout: %v", err)
	}
	if got.Soak != "2m" || got.Timeout != "15m" || got.Poll != "10s" {
		t.Fatalf("unexpected rollout defaults: %+v", got)
	}

	if _, err := normalizeRollout(&model.HostImageRollout{Canary: "n1", Soak: "nope"}); err == nil {
		t.Fatal("expected an error for an invalid soak")
	}
	if _, err := normalizeRollout(&model.HostImageRollout{Canary: "n1", Timeout: "nope"}); err == nil {
		t.Fatal("expected an error for an invalid timeout")
	}
	if _, err := normalizeRollout(&model.HostImageRollout{Canary: "n1", Poll: "nope"}); err == nil {
		t.Fatal("expected an error for an invalid poll")
	}
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
