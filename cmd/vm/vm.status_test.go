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
