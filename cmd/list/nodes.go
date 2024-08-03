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

var (
	query         string
	org           string
	details       bool
	nodesResponse model.NodesResponse
)

var NodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: aliases.NodesAliases,
	Short:   constants.ListNodesShortDesc,
	Long:    constants.ListNodesLongDesc,
	Run:     executeRetrieveNodes,
}

func executeRetrieveNodes(cmd *cobra.Command, args []string) {
	requestBody, url, err := prepareRequest(query)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendNodeRequest(requestBody, url); err != nil {
		fmt.Println("Error sending list nodes request:", err)
		os.Exit(1)
	}

	if details {
		render.RenderNodes(nodesResponse.Nodes)
	} else {
		render.RenderNodesTabWriter(nodesResponse.Nodes)
	}
}

func prepareRequest(query string) (interface{}, string, error) {
	if query == "" {
		return nil, clients.BuildURL("core", "v1", "ListNodePool"), nil
	}

	nodeQueries, err := utils.CreateNodeQuery(query)
	if err != nil {
		return nil, "", err
	}
	requestBody := map[string][]model.NodeQuery{"query": nodeQueries}
	url := clients.BuildURL("core", "v1", "QueryNodePool")
	return requestBody, url, nil
}

func sendNodeRequest(requestBody interface{}, url string) error {
	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Headers:     map[string]string{"Content-Type": "application/json"},
		RequestBody: requestBody,
		Response:    &nodesResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	NodesCmd.Flags().StringVarP(&query, constants.QueryFlag, constants.QueryFlagShorthandFlag, "", constants.NodeQueryDescription)
	NodesCmd.Flags().BoolVarP(&details, "details", "d", false, "Display detailed node information")
}
