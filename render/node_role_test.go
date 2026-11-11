package render

import (
	"testing"

	"github.com/c12s/cockpit/model"
)

// A VM is one a c12s host created in its libvirt: its seed sets the vm-host label.
func TestNodeLabel(t *testing.T) {
	vm := model.Node{ID: "bbe0", Labels: []model.Label{{Key: "vm-host", Value: "5478"}}}
	if got := NodeLabel(vm); got != "bbe0 (vm)" {
		t.Errorf("vm: %q", got)
	}
	host := model.Node{ID: "5478", Capabilities: []string{"os_update", "run_vms"}}
	if got := NodeLabel(host); got != "5478" {
		t.Errorf("VM host: %q", got)
	}
}
