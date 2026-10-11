package cmd

import (
	"testing"

	"github.com/c12s/cockpit/model"
)

func TestEvaluateGate(t *testing.T) {
	const want, other = "sha256:new", "sha256:old"
	node := func(phase, booted string, last *model.OsAttempt) model.Node {
		return model.Node{OsStatus: &model.OsStatus{Phase: phase, BootedDigest: booted, LastAttempt: last}}
	}
	cases := []struct {
		name string
		node model.Node
		want gateVerdict
	}{
		{"no status yet", model.Node{}, gatePending},
		{"still on old image", node("healthy", other, nil), gatePending},
		{"booted, greenboot running", node("switching", want, nil), gatePending},
		{"booted and healthy", node("healthy", want, nil), gateHealthy},
		{"rolled back this digest", node("rolled_back", other, &model.OsAttempt{Digest: want, Result: "rolled_back"}), gateRolledBack},
		{"older rollback of another digest", node("rolled_back", other, &model.OsAttempt{Digest: other, Result: "rolled_back"}), gatePending},
		{"pull failed, retrying", node("degraded", other, &model.OsAttempt{Digest: want, Result: "failed"}), gatePending},
	}
	for _, c := range cases {
		if got := evaluateGate(c.node, want); got != c.want {
			t.Errorf("%s: got %s, want %s", c.name, got, c.want)
		}
	}
}
