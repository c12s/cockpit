package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	claimNodesShortDescription = "Claim nodes for an organization based on specific criteria"
	claimNodesLongDescription  = "Claims nodes for an organization based on a defined query that specifies criteria like labels.\n\n" +
		"Example:\n" +
		"claim-nodes --org 'myOrg' --query '[{\"labelKey\": \"key\", \"shouldBe\": \"<||=||>\", \"value\": \"value\"}]'"
	orgFlag     = "org"
	queryFlag   = "query"
	orgFlagSh   = "o"
	queryFlagSh = "q"
)

var (
	org   string
	query string
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
		fmt.Printf("%v", err)
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
	}
}

func init() {
	ClaimNodesCmd.Flags().StringVarP(&org, orgFlag, orgFlagSh, "", "Organization name (required)")
	ClaimNodesCmd.Flags().StringVarP(&query, queryFlag, queryFlagSh, "", "Query in JSON format specifying node selection criteria (required)")
	ClaimNodesCmd.MarkFlagRequired(orgFlag)
	ClaimNodesCmd.MarkFlagRequired(queryFlag)
}

func claimNodes(request model.ClaimNodesRequest, token string) error {
	claimNodesURL := clients.Clients.Gateway + "/apis/core/v1/nodes"
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("PATCH", claimNodesURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %s", resp.Status)
	}

	var nodesResponse model.ClaimNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&nodesResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if len(nodesResponse.Nodes) == 0 {
		fmt.Println("No nodes were found.")
	} else {
		model.RenderNodes(nodesResponse.Nodes)
	}

	return nil
}
