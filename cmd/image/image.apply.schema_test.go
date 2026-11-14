package cmd

import (
	"strings"
	"testing"
)

func TestQuasarVerdict(t *testing.T) {
	applySchemaNamespace, applySchemaVersion = "default", "v1.0.0"
	if err := quasarVerdict("c12s", validateConfigurationResponse{IsValid: true}); err != nil {
		t.Fatalf("valid: %v", err)
	}
	err := quasarVerdict("c12s", validateConfigurationResponse{Message: "spec.nodes.0.id: Does not match pattern; spec.rollout.soak: Invalid type"})
	if err == nil || !strings.Contains(err.Error(), "\n  - spec.rollout.soak") {
		t.Fatalf("invalid: %v", err)
	}
	err = quasarVerdict("c12s", validateConfigurationResponse{Status: 3, Message: "No schema with key 'c12s/default/host-image-spec/v1.0.0' found!"})
	if err == nil || !strings.Contains(err.Error(), "cockpit create schema --org c12s --namespace default --schema-name host-image-spec --version v1.0.0") {
		t.Fatalf("missing schema: %v", err)
	}
}

func TestSpecOrgIsLenient(t *testing.T) {
	doc := []byte("metadata:\n  org: c12s\nspec:\n  nodes:\n    - {id: x, canery: true}\n")
	if got := specOrg(doc); got != "c12s" {
		t.Fatalf("got %q", got)
	}
	if got := specOrg([]byte("not: [yaml")); got != "" {
		t.Fatalf("broken yaml: got %q", got)
	}
}
