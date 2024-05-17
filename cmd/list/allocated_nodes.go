package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const (
	allocatedNodesShortDescription = "List organization nodes"
	allocatedNodesLongDescription  = "You can search for nodes organization has allocated. \n" +
		"You can add query to search nodes by labels."
	orgFlag        = "org"
	queryFlag      = "query"
	orgFlagShort   = "o"
	queryFlagShort = "q"
	orgFlagValue   = ""
	queryFlagValue = ""
)

var (
	org string
)

var AllocatedNodesCmd = &cobra.Command{
	Use:   "allocated",
	Short: allocatedNodesShortDescription,
	Long:  allocatedNodesLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			fmt.Printf("Error reading token: %v\n", err)
			os.Exit(1)
		}

		var url string
		var request model.ClaimNodesRequest

		if query == "" {
			url = clients.Clients.Gateway + "/apis/core/v1/nodes/allocated"
			request = model.ClaimNodesRequest{
				Org: org,
			}
		} else {
			url = clients.Clients.Gateway + "/apis/core/v1/nodes/allocated/query_match"
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

		if err := sendRequestForAllocatedNodes(url, request, string(token)); err != nil {
			fmt.Printf("Error processing nodes: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	AllocatedNodesCmd.Flags().StringVarP(&org, orgFlag, orgFlagShort, orgFlagValue, "Organization name")
	AllocatedNodesCmd.Flags().StringVarP(&query, queryFlag, queryFlagShort, queryFlagValue, "Query JSON for node allocation")
	AllocatedNodesCmd.MarkFlagRequired(orgFlag)
}

func sendRequestForAllocatedNodes(url string, requestBody model.ClaimNodesRequest, token string) error {
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestJSON))
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
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		return fmt.Errorf("request failed with status %d %s: %s", resp.StatusCode, resp.Status, string(bodyBytes))
	}

	var nodesResponse model.NodesResponse
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
