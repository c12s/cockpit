package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	claimNodesShortDescription = "Claim nodes for an organization based on specific criteria"
	claimNodesLongDescription  = "Claims nodes for an organization based on a defined query that specifies criteria like labels.\n\n" +
		"Example:\n" +
		"claim-nodes --org 'myOrg' --query '[{\"labelKey\": \"key\", \"shouldBe\": \"<||=||>\", \"value\": \"value\"}]'"

	// Flag Constants
	orgFlag   = "org"
	queryFlag = "query"

	// Flag Shorthand Constants
	orgFlagShortHand   = "o"
	queryFlagShortHand = "q"

	// Flag Descriptions
	orgDesc   = "Organization name (required)"
	queryDesc = "Query in JSON format specifying node selection criteria (required)"
)

var (
	org      string
	query    string
	response model.ClaimNodesResponse
)

var ClaimNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: claimNodesShortDescription,
	Long:  claimNodesLongDescription,
	Run:   executeClaimNodes,
}

func executeClaimNodes(cmd *cobra.Command, args []string) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	request := model.ClaimNodesRequest{
		Org:   org,
		Query: []model.NodeQuery{},
	}
	if err := json.Unmarshal([]byte(query), &request.Query); err != nil {
		fmt.Printf("Error parsing query JSON: %v\n", err)
		return
	}

	if err := claimNodes(request, token); err != nil {
		fmt.Printf("Error claiming nodes: %v\n", err)
		os.Exit(1)
	}

	if len(response.Nodes) == 0 {
		fmt.Println("No nodes were found.")
	} else {
		render.RenderNodes(response.Nodes)
	}
	fmt.Println()
}

func claimNodes(request model.ClaimNodesRequest, token string) error {
	claimNodesURL := clients.BuildURL("core", "v1", "ClaimOwnership")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         claimNodesURL,
		Method:      "PATCH",
		Token:       token,
		RequestBody: request,
		Response:    &response,
		Timeout:     10 * time.Second,
	})
}
func init() {
	ClaimNodesCmd.Flags().StringVarP(&org, orgFlag, orgFlagShortHand, "", orgDesc)
	ClaimNodesCmd.Flags().StringVarP(&query, queryFlag, queryFlagShortHand, "", queryDesc)
	ClaimNodesCmd.MarkFlagRequired(orgFlag)
	ClaimNodesCmd.MarkFlagRequired(queryFlag)
}
