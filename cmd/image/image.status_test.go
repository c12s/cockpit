package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/c12s/cockpit/model"
)

func TestDescribeImage(t *testing.T) {
	d41 := "sha256:b4e5e68be4ec9bba1526a8bb14096b3ecf52d19a685fa7d71dc35b5e76989531"
	tags := map[string][]string{d41: {"v4.1"}}
	cases := []struct {
		digest string
		wide   bool
		want   string
	}{
		{d41, false, "v4.1 (sha256:b4e5e68be4ec)"},
		{d41, true, "v4.1 (" + d41 + ")"},
		{"sha256:" + strings.Repeat("c", 64), false, "sha256:cccccccccccc"},
		{"", false, "-"},
	}
	for _, c := range cases {
		if got := describeImage(c.digest, tags, c.wide); got != c.want {
			t.Errorf("describeImage(%q, wide=%v) = %q, want %q", c.digest, c.wide, got, c.want)
		}
	}
}

// Columns must line up whatever the value lengths are.
func TestRenderStatusTableAligns(t *testing.T) {
	nodes := []model.Node{
		{ID: "a", OsStatus: &model.OsStatus{Phase: "healthy", BootedVersion: "45.20260830.0", BootedDigest: "sha256:" + strings.Repeat("1", 64)}},
		{ID: "b-much-longer-id", OsStatus: &model.OsStatus{Phase: "rolled_back", BootedVersion: "v4"}},
	}
	out := renderStatusTable(nodes, nil, false)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("want header + 2 rows, got:\n%s", out)
	}
	col := strings.Index(lines[0], "PHASE")
	for _, l := range lines[1:] {
		if len(l) < col || (l[col-1] != ' ') || l[col] == ' ' {
			t.Fatalf("PHASE column not aligned at %d:\n%s", col, out)
		}
	}
}

// Stale and offline nodes show how old their report is; online ones just say so.
func TestStatusLiveness(t *testing.T) {
	old := model.FlexInt64(time.Now().Add(-59 * time.Minute).Unix())
	if got := livenessCell(model.Node{Liveness: "offline", LastSeen: old}); !strings.HasPrefix(got, "offline (59m") || !strings.HasSuffix(got, " ago)") {
		t.Errorf("offline: %q", got)
	}
	if got := livenessCell(model.Node{Liveness: "online", LastSeen: old}); got != "online" {
		t.Errorf("online: %q", got)
	}
	if got := livenessCell(model.Node{}); got != "-" {
		t.Errorf("unknown: %q", got)
	}
	out := renderStatusTable([]model.Node{{ID: "a", Liveness: "online"}}, nil, false)
	if !strings.Contains(strings.Split(out, "\n")[0], "LIVENESS") {
		t.Fatalf("no LIVENESS column:\n%s", out)
	}
}
