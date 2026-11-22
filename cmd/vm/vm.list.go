package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	vmListHost string
	vmListResp model.VMDesiredGetResponse
)

var VmsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the desired VM set of a host node",
	Run:     executeVmsList,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.HostFlag})
	},
}

func executeVmsList(cmd *cobra.Command, args []string) {
	url := clients.BuildURL("core", "v1", "GetDesiredVms")
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Error reading token:", err)
		os.Exit(1)
	}

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: model.VMDesiredGetRequest{HostID: vmListHost},
		Response:    &vmListResp,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		fmt.Println("Error fetching desired VMs:", err)
		os.Exit(1)
	}

	// NODE ID is derived (host id + VM name), the id the VM registers with:
	// grant os.update for it before the VM exists.
	fmt.Printf("%-12s %-36s  %-28s %-20s %-6s %s\n", "NAME", "NODE ID", "QCOW2", "DIGEST", "VCPUS", "MEMORY (MB)")
	for _, vm := range vmListResp.VMs {
		fmt.Printf("%-12s %-36s  %-28s %-20s %-6d %d\n", vm.Name, vmID(vmListHost, vm.Name), vm.Qcow2URL, shortDigest(vm.Digest), vm.VCPUs, vm.MemoryMB)
	}
}

func vmID(hostID, name string) string {
	id, err := utils.VMNodeID(hostID, name)
	if err != nil {
		return "-"
	}
	return id
}

func shortDigest(digest string) string {
	if len(digest) > 19 {
		return digest[:19] + "..."
	}
	return digest
}

func init() {
	VmsCmd.AddCommand(VmsListCmd)
	VmsListCmd.Flags().StringVar(&vmListHost, constants.HostFlag, "", constants.HostDescription)
	VmsListCmd.MarkFlagRequired(constants.HostFlag)
}
