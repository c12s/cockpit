package utils

import "testing"

// Same vector as star's TestVMNodeID and Python's uuid.uuid5.
func TestVMNodeID(t *testing.T) {
	got, err := VMNodeID("5c1e7d2a-1b3c-4d5e-8f90-123456789abc", "vm-01")
	if err != nil || got != "f210cdd5-ce6a-561e-a06a-29929fc319bd" {
		t.Fatalf("VMNodeID = %q, %v", got, err)
	}
	if _, err := VMNodeID("nope", "vm-01"); err == nil {
		t.Fatal("malformed host id must fail")
	}
}
