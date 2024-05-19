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
	nodesShortDescription = "Retrieve a list of all available nodes"
	nodesLongDescription  = "Retrieve a comprehensive list of all available nodes in the system.\n" +
		"These nodes can be allocated to your organization based on your requirements.\n\n" +
		"Example:\n" +
		"nodes --query '[{\"labelKey\": \"labelKey\", \"shouldBe\": \"> || < || =\", \"value\": \"2\"}]'"
)

var (
	query         string
	org           string
	nodesResponse model.NodesResponse
)

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: nodesShortDescription,
	Long:  nodesLongDescription,
	Run:   executeRetrieveNodes,
}

func executeRetrieveNodes(cmd *cobra.Command, args []string) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.AvailableNodesEndpoint
	var nodeQueries []model.NodeQuery
	var requestBody interface{}
	if query != "" {
		url = clients.AvailableNodesQueryEndpoint
		if err := json.Unmarshal([]byte(query), &nodeQueries); err != nil {
			fmt.Printf("Error parsing query JSON: %v\n", err)
			println()
			os.Exit(1)
		}
		requestBody = map[string][]model.NodeQuery{"query": nodeQueries}
	}

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Headers:     map[string]string{"Content-Type": "application/json"},
		RequestBody: requestBody,
		Response:    &nodesResponse,
		Token:       token,
		Timeout:     10 * time.Second,
	})
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		println()
		os.Exit(1)
	}

	if len(nodesResponse.Nodes) == 0 {
		fmt.Println("No nodes were found.")
	} else {
		render.RenderNodes(nodesResponse.Nodes)
	}
	println()
}

func init() {
	NodesCmd.Flags().StringVarP(&query, queryFlag, queryFlagShortHand, "", queryFlag)
}
