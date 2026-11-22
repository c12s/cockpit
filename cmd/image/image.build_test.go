package cmd

import "testing"

func TestParseBuildArgs(t *testing.T) {
	got, err := parseBuildArgs([]string{"OS_VERSION=v3.18", "BASE=localhost:5001/os/edge-base:v3.17"})
	if err != nil {
		t.Fatalf("parseBuildArgs: %v", err)
	}
	if got["OS_VERSION"] != "v3.18" || got["BASE"] != "localhost:5001/os/edge-base:v3.17" {
		t.Fatalf("unexpected map: %v", got)
	}

	if _, err := parseBuildArgs([]string{"NO_EQUALS"}); err == nil {
		t.Fatal("expected an error for a value without '='")
	}
	if _, err := parseBuildArgs([]string{"=v"}); err == nil {
		t.Fatal("expected an error for an empty key")
	}
	if got, err := parseBuildArgs(nil); err != nil || got != nil {
		t.Fatalf("nil input should return nil, nil; got %v, %v", got, err)
	}
}
