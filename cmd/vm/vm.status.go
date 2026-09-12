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
	vmStatusHost string
	vmStatusResp model.VMStatusGetResponse
)

var VmsStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show observed VM status on a host node",
	Long:    "Display the state of the VMs the host node currently manages (running, shut off, ...).",
	Run:     executeVmsStatus,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.HostFlag})
	},
}

func executeVmsStatus(cmd *cobra.Command, args []string) {
	url := clients.BuildURL("core", "v1", "GetVmStatuses")
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Error reading token:", err)
		os.Exit(1)
	}

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: model.VMStatusGetRequest{HostID: vmStatusHost},
		Response:    &vmStatusResp,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		fmt.Println("Error fetching VM statuses:", err)
		os.Exit(1)
	}

	fmt.Printf("%-20s %-12s %s\n", "NAME", "STATE", "BOOTED OS")
	for _, status := range vmStatusResp.Statuses {
		state := status.State
		if state == "" {
			state = "-"
		}
		bootedOs := status.BootedOS
		if bootedOs == "" {
			bootedOs = "-"
		}
		fmt.Printf("%-20s %-12s %s\n", status.Name, state, bootedOs)
	}
}

func init() {
	VmsCmd.AddCommand(VmsStatusCmd)
	VmsStatusCmd.Flags().StringVar(&vmStatusHost, constants.HostFlag, "", constants.HostDescription)
	VmsStatusCmd.MarkFlagRequired(constants.HostFlag)
}
