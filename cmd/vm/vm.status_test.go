package cmd

import (
	"testing"

	"github.com/c12s/cockpit/model"
)

func TestBootedOS(t *testing.T) {
	org := map[string]model.Node{
		"a": {ID: "a", OsStatus: &model.OsStatus{BootedVersion: "v4.1", Phase: "healthy"}},
		"b": {ID: "b"},
	}
	pool := map[string]model.Node{"p": {ID: "p"}}
	cases := map[string]string{
		"a": "v4.1 healthy",
		"b": "- (no OS report yet)",
		"p": "- (in pool, not claimed)",
		"x": "- (not enrolled yet)",
	}
	for id, want := range cases {
		if got := bootedOS(id, org, pool); got != want {
			t.Errorf("bootedOS(%s) = %q, want %q", id, got, want)
		}
	}
	if got := bootedOS("a", nil, nil); got != "- (add --org to look it up)" {
		t.Errorf("without --org: %q", got)
	}
}

func TestCheckVMHost(t *testing.T) {
	vmHost := model.Node{ID: "n1", Capabilities: []string{"os_update", "run_vms"}}
	plain := model.Node{ID: "h1", Capabilities: []string{"os_update"}}
	if err := checkVMHost(vmHost, true, 2); err != nil {
		t.Errorf("VM host must pass: %v", err)
	}
	if err := checkVMHost(plain, true, 2); err == nil {
		t.Error("a host without run_vms must be refused")
	}
	if err := checkVMHost(plain, true, 0); err != nil {
		t.Errorf("clearing the VM list is always allowed: %v", err)
	}
	if err := checkVMHost(model.Node{}, false, 2); err == nil {
		t.Error("an unknown host must be refused")
	}
}
