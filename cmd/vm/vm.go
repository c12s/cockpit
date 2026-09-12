package cmd

import (
	"github.com/spf13/cobra"
)

// VmsCmd groups the VM-management commands (Phase 4). They only apply to
// nodes reporting the run_vms capability; the CP stores the desired VM set on
// the host's node record and the host's Star reconciler converges to it.
var VmsCmd = &cobra.Command{
	Use:     "vms",
	Aliases: []string{"vm"},
	Short:   "Manage libvirt VMs on host nodes",
	Long:    "Declare, list and inspect bootc VMs managed by host nodes with the run_vms capability.",
}
