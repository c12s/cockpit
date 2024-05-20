package cmd

import (
	"encoding/json"
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
		"nodes-allocated --org \"org\" --query '[{\"labelKey\": \"labelKey\", \"shouldBe\": \">||<||==\", \"value\": \"2\"}]'"

	// Flag Constants
	orgFlag   = "org"
	queryFlag = "query"

	// Flag Shorthand Constants
	orgFlagShortHand   = "o"
	queryFlagShortHand = "q"

	// Flag Descriptions
	orgDesc   = "Organization name (required)"
	queryDesc = "Query JSON for node allocation"
)

var AllocatedNodesCmd = &cobra.Command{
	Use:   "allocated",
	Short: allocatedNodesShortDescription,
	Long:  allocatedNodesLongDescription,
	Run:   executeAllocatedNodes,
}

func executeAllocatedNodes(cmd *cobra.Command, args []string) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	var allocatedNodesURL string
	var request model.ClaimNodesRequest
	if query == "" {
		allocatedNodesURL = clients.BuildURL("core", "v1", "ListOrgOwnedNodes")
		request.Org = org
	} else {
		allocatedNodesURL = clients.BuildURL("core", "v1", "QueryOrgOwnedNodes")
		var nodeQueries []model.NodeQuery
		if err := json.Unmarshal([]byte(query), &nodeQueries); err != nil {
			fmt.Printf("Error parsing query JSON: %v\n", err)
			os.Exit(1)
		}
		request = model.ClaimNodesRequest{
			Org:   org,
			Query: nodeQueries,
		}
	}

	var nodesResponse model.NodesResponse
	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         allocatedNodesURL,
		Method:      "GET",
		RequestBody: request,
		Response:    &nodesResponse,
		Token:       token,
		Timeout:     10 * time.Second,
	})

	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}

	if len(nodesResponse.Nodes) == 0 {
		fmt.Println("No nodes were found.")
	} else {
		fmt.Println("")
		render.RenderNodes(nodesResponse.Nodes)
	}
	println()
}

func init() {
	AllocatedNodesCmd.Flags().StringVarP(&org, orgFlag, orgFlagShortHand, "", orgDesc)
	AllocatedNodesCmd.Flags().StringVarP(&query, queryFlag, queryFlagShortHand, "", queryDesc)
	AllocatedNodesCmd.MarkFlagRequired(orgFlag)
}
