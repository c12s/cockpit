package cmd

import (
	"strings"
	"testing"

	"github.com/c12s/cockpit/model"
)

func TestRollbackTarget(t *testing.T) {
	healthyOn := func(rollbackImage, rollbackDigest string, last *model.OsAttempt) model.Node {
		return model.Node{ID: "n1", OsStatus: &model.OsStatus{
			Phase: "healthy", BootedDigest: "sha256:new",
			RollbackImage: rollbackImage, RollbackDigest: rollbackDigest, LastAttempt: last,
		}}
	}
	ref, digest, err := rollbackTarget(healthyOn("192.168.122.1:5001/os/edge-base:v4.0@sha256:old", "sha256:old", nil))
	if err != nil || ref != "192.168.122.1:5001/os/edge-base:v4.0" || digest != "sha256:old" {
		t.Fatalf("plain rollback: %q %q %v (the @digest must be stripped from the ref)", ref, digest, err)
	}

	cases := []struct {
		name string
		node model.Node
		want string
	}{
		{"no OS report", model.Node{ID: "n1"}, "no OS report"},
		{"no previous deployment", healthyOn("", "", nil), "no previous deployment"},
		// after a greenboot rollback the other deployment is the rejected
		// image: going "back" to it would boot the broken image again
		{"previous is the rejected image", healthyOn("r/os:broken", "sha256:bad",
			&model.OsAttempt{Digest: "sha256:bad", Result: "rolled_back"}), "rejected"},
		{"update in flight", model.Node{ID: "n1", OsStatus: &model.OsStatus{Phase: "switching", RollbackImage: "r/os:v4.0", RollbackDigest: "sha256:old"}}, "in progress"},
	}
	for _, c := range cases {
		if _, _, err := rollbackTarget(c.node); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want error containing %q, got %v", c.name, c.want, err)
		}
	}
}
