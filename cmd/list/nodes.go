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
	nodesShortDescription = "Retrieve a list of all available nodes"
	nodesLongDescription  = `Retrieve a comprehensive list of all available nodes in the system.
These nodes can be allocated to your organization based on your requirements.
You can use a query to filter the nodes using operators like >, =, !=, and < with the label values.

Examples:
- cockpit list nodes --query 'labelKey >||=||!=||< value'
- cockpit list nodes --query 'memory-totalGB > 2'`
)

var (
	query         string
	org           string
	nodesResponse model.NodesResponse
)

var NodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: aliases.NodesAliases,
	Short:   nodesShortDescription,
	Long:    nodesLongDescription,
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

	render.RenderNodes(nodesResponse.Nodes)
	println()
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
		Timeout:     10 * time.Second,
	})
}

func init() {
	NodesCmd.Flags().StringVarP(&query, queryFlag, queryShorthandFlag, "", queryFlag)
}
