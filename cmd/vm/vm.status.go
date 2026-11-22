package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	vmStatusHost string
	vmStatusOrg  string
	vmStatusResp model.VMStatusGetResponse
)

var VmsStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show observed VM status on a host node",
	Long: "Display the state of the VMs the host node currently manages (running, shut off, ...). " +
		"With --org, BOOTED OS shows each VM's own OS report (version and phase), found by the VM's node id.",
	Run: executeVmsStatus,
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

	// a VM reports its OS on its own node record (id = UUIDv5(host, name));
	// the host only knows whether the domain runs
	var orgNodes, poolNodes map[string]model.Node
	if vmStatusOrg != "" {
		orgNodes = nodesByID(token, "ListOrgOwnedNodes", model.ClaimNodesRequest{Org: vmStatusOrg})
	}
	poolNodes = nodesByID(token, "ListNodePool", nil)

	fmt.Printf("%-12s %-36s  %-12s %s\n", "NAME", "NODE ID", "STATE", "BOOTED OS")
	for _, status := range vmStatusResp.Statuses {
		state := status.State
		if state == "" {
			state = "-"
		}
		bootedOs := bootedOS(vmID(vmStatusHost, status.Name), orgNodes, poolNodes)
		fmt.Printf("%-12s %-36s  %-12s %s\n", status.Name, vmID(vmStatusHost, status.Name), state, bootedOs)
	}
}

// bootedOS describes a VM's OS from its node record. orgNodes is nil when no
// --org was given.
func bootedOS(id string, orgNodes, poolNodes map[string]model.Node) string {
	if node, ok := orgNodes[id]; ok {
		st := node.OsStatus
		if st == nil || (st.BootedVersion == "" && st.Phase == "") {
			return "- (no OS report yet)"
		}
		version, phase := st.BootedVersion, st.Phase
		if version == "" {
			version = "?"
		}
		return strings.TrimSpace(version + " " + phase)
	}
	if _, ok := poolNodes[id]; ok {
		return "- (in pool, not claimed)"
	}
	if orgNodes == nil {
		return "- (add --org to look it up)"
	}
	return "- (not enrolled yet)"
}

// nodesByID lists nodes through the gateway; a failed lookup only leaves the
// BOOTED OS column empty.
func nodesByID(token, method string, body interface{}) map[string]model.Node {
	var resp model.NodesResponse
	err := utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         clients.BuildURL("core", "v1", method),
		Method:      "GET",
		Token:       token,
		RequestBody: body,
		Response:    &resp,
		Timeout:     30 * time.Second,
	})
	nodes := map[string]model.Node{}
	if err != nil {
		return nodes
	}
	for _, node := range resp.Nodes {
		nodes[node.ID] = node
	}
	return nodes
}

func init() {
	VmsCmd.AddCommand(VmsStatusCmd)
	VmsStatusCmd.Flags().StringVar(&vmStatusHost, constants.HostFlag, "", constants.HostDescription)
	VmsStatusCmd.Flags().StringVarP(&vmStatusOrg, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", "Organization the VMs are claimed into (fills BOOTED OS)")
	VmsStatusCmd.MarkFlagRequired(constants.HostFlag)
}
