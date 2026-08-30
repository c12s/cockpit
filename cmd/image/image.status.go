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
	imageStatusOrg   string
	imageStatusNodes model.NodesResponse
)

var ImageStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show OS status of allocated nodes",
	Long:    "Display the current OS image, digest, and update phase (healthy / switching / rolled_back) for all nodes of an organization.",
	Run:     executeImageStatus,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.OrganizationFlag})
	},
}

func executeImageStatus(cmd *cobra.Command, args []string) {
	url := clients.BuildURL("core", "v1", "ListOrgOwnedNodes")
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Error reading token:", err)
		os.Exit(1)
	}

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: model.ClaimNodesRequest{Org: imageStatusOrg},
		Response:    &imageStatusNodes,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		fmt.Println("Error fetching nodes:", err)
		os.Exit(1)
	}

	fmt.Printf("%-38s %-12s %-12s %-20s %s\n", "NODE ID", "PHASE", "VERSION", "DIGEST", "LAST ATTEMPT")
	for _, node := range imageStatusNodes.Nodes {
		id := node.ID
		phase := "-"
		version := "-"
		digest := "-"
		attempt := "-"
		if node.OsStatus != nil {
			if node.OsStatus.Phase != "" {
				phase = node.OsStatus.Phase
			}
			if node.OsStatus.BootedVersion != "" {
				version = node.OsStatus.BootedVersion
			}
			if len(node.OsStatus.BootedDigest) > 19 {
				digest = node.OsStatus.BootedDigest[:19] + "..."
			}
			if node.OsStatus.LastAttempt != nil {
				attempt = fmt.Sprintf("%s (%s)", node.OsStatus.LastAttempt.Result, time.Unix(int64(node.OsStatus.LastAttempt.Timestamp), 0).Format("15:04:05"))
			}
		}
		fmt.Printf("%-38s %-12s %-12s %-20s %s\n", id, phase, version, digest, attempt)
	}
}

func init() {
	ImageCmd.AddCommand(ImageStatusCmd)
	ImageStatusCmd.Flags().StringVarP(&imageStatusOrg, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ImageStatusCmd.MarkFlagRequired(constants.OrganizationFlag)
}
