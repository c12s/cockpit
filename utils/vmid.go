package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// VMNodeID is a VM's node id: UUIDv5 with the host's node id as namespace
// and the VM name as name. star derives the same id when it creates the VM
// (star internal/vm VMNodeID), so a VM's id is known before it exists.
func VMNodeID(hostID, vmName string) (string, error) {
	ns, err := uuid.Parse(hostID)
	if err != nil {
		return "", fmt.Errorf("host id %q is not a UUID: %w", hostID, err)
	}
	return uuid.NewSHA1(ns, []byte(vmName)).String(), nil
}
