package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	allocatedNodesShortDescription = "List organization nodes"
	allocatedNodesLongDescription  = `This command allows you to list all nodes allocated to a specified organization.
You can also use a query to search for nodes based on their labels.
The query format allows you to filter nodes using operators like >, =, !=, and < with the label values.

Examples:
- cockpit list nodes allocated --org 'org' --query 'labelKey >||=||!=||< value'
- cockpit list nodes allocated --org 'org' --query 'memory-totalGB > 2'`

	// Flag Constants
	organizationFlag = "org"
	queryFlag        = "query"

	// Flag Shorthand Constants
	organizationShorthandFlag = "r"
	queryShorthandFlag        = "q"

	// Flag Descriptions
	organizationDescription = "Organization name (required)"
	queryDescription        = "Query label for finding specific nodes"
)

var AllocatedNodesCmd = &cobra.Command{
	Use:     "allocated",
	Aliases: aliases.AllocatedAliases,
	Short:   allocatedNodesShortDescription,
	Long:    allocatedNodesLongDescription,
	Run:     executeAllocatedNodes,
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

	render.RenderNodes(nodesResponse.Nodes)
	println()
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
		Timeout:     10 * time.Second,
	})
}

func init() {
	AllocatedNodesCmd.Flags().StringVarP(&org, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	AllocatedNodesCmd.Flags().StringVarP(&query, queryFlag, queryShorthandFlag, "", queryDescription)
	AllocatedNodesCmd.MarkFlagRequired(organizationFlag)
}
