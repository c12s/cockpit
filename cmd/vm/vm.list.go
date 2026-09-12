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

	fmt.Printf("%-20s %-48s %-20s %-6s %s\n", "NAME", "QCOW2", "DIGEST", "VCPUS", "MEMORY (MB)")
	for _, vm := range vmListResp.VMs {
		fmt.Printf("%-20s %-48s %-20s %-6d %d\n", vm.Name, vm.Qcow2URL, shortDigest(vm.Digest), vm.VCPUs, vm.MemoryMB)
	}
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
