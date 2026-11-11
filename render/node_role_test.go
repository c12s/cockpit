package render

import (
	"testing"

	"github.com/c12s/cockpit/model"
)

// A nested VM is one a c12s host created: its seed sets the vm-host label.
func TestNodeLabel(t *testing.T) {
	vm := model.Node{ID: "bbe0", Labels: []model.Label{{Key: "vm-host", Value: "5478"}}}
	if got := NodeLabel(vm); got != "bbe0 (nested VM)" {
		t.Errorf("nested VM: %q", got)
	}
	host := model.Node{ID: "5478", Capabilities: []string{"os_update", "run_vms"}}
	if got := NodeLabel(host); got != "5478" {
		t.Errorf("VM host: %q", got)
	}
}
