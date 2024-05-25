package cmd

import (
	"fmt"
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
	allocatedNodesLongDescription  = "You can search for nodes organization has allocated. \n" +
		"You can add query to search nodes by labels.\n\n" +
		"Example:\n" +
		"nodes-allocated --org 'labelKey >||=||< value'" +
		"nodes-allocated --query 'memory-totalGB > 2'"

	// Flag Constants
	orgFlag   = "org"
	queryFlag = "query"

	// Flag Shorthand Constants
	orgFlagShortHand   = "r"
	queryFlagShortHand = "q"

	// Flag Descriptions
	orgDesc   = "Organization name (required)"
	queryDesc = "Query label for finding specific nodes"
)

var AllocatedNodesCmd = &cobra.Command{
	Use:   "allocated",
	Short: allocatedNodesShortDescription,
	Long:  allocatedNodesLongDescription,
	Run:   executeAllocatedNodes,
}

func executeAllocatedNodes(cmd *cobra.Command, args []string) {
	requestBody, url, err := prepareAllocatedRequest(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := sendAllocatedNodeRequest(requestBody, url); err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}

	render.RenderNodes(nodesResponse.Nodes)
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
	AllocatedNodesCmd.Flags().StringVarP(&org, orgFlag, orgFlagShortHand, "", orgDesc)
	AllocatedNodesCmd.Flags().StringVarP(&query, queryFlag, queryFlagShortHand, "", queryDesc)
	AllocatedNodesCmd.MarkFlagRequired(orgFlag)
}
