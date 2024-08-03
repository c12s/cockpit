package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"

	"github.com/spf13/cobra"
)

var AllocatedNodesCmd = &cobra.Command{
	Use:     "allocated",
	Aliases: aliases.AllocatedAliases,
	Short:   constants.AllocatedNodesShortDesc,
	Long:    constants.AllocatedNodesLongDesc,
	Run:     executeAllocatedNodes,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.OrganizationFlag})
	},
}

func executeAllocatedNodes(cmd *cobra.Command, args []string) {
	requestBody, url, err := prepareAllocatedRequest(query)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendAllocatedNodeRequest(requestBody, url); err != nil {
		fmt.Println("Error sending list allocated nodes request:", err)
		os.Exit(1)
	}

	if details {
		render.RenderNodes(nodesResponse.Nodes)
	} else {
		render.RenderNodesTabWriter(nodesResponse.Nodes)
	}
}

func prepareAllocatedRequest(query string) (interface{}, string, error) {
	var request model.ClaimNodesRequest
	var allocatedNodesURL string

	if query == "" {
		allocatedNodesURL = clients.BuildURL("core", "v1", "ListOrgOwnedNodes")
		request.Org = org
	} else {
		allocatedNodesURL = clients.BuildURL("core", "v1", "QueryOrgOwnedNodes")
		nodeQueries, err := utils.CreateNodeQuery(query)
		if err != nil {
			return nil, "", err
		}
		request = model.ClaimNodesRequest{
			Org:   org,
			Query: nodeQueries,
		}
	}
	return request, allocatedNodesURL, nil
}

func sendAllocatedNodeRequest(requestBody interface{}, url string) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("Error reading token: %v", err)
	}
	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		Headers:     map[string]string{"Content-Type": "application/json"},
		RequestBody: requestBody,
		Response:    &nodesResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	AllocatedNodesCmd.Flags().StringVarP(&org, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	AllocatedNodesCmd.Flags().StringVarP(&query, constants.QueryFlag, constants.QueryFlagShorthandFlag, "", constants.NodeQueryDescription)
	AllocatedNodesCmd.Flags().BoolVarP(&details, "details", "d", false, "Display detailed node information")

	AllocatedNodesCmd.MarkFlagRequired(constants.OrganizationFlag)
}
