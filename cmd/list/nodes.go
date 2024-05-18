package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	nodesShortDescription = "Retrieve a list of all available nodes"
	nodesLongDescription  = "Retrieve a comprehensive list of all available nodes in the system.\n" +
		"These nodes can be allocated to your organization based on your requirements.\n\n" +
		"Example:\n" +
		"nodes --query '[{\"labelKey\": \"labelKey\", \"shouldBe\": \"> || < || =\", \"value\": \"2\"}]'"
)

var (
	query string
)

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: nodesShortDescription,
	Long:  nodesLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := utils.ReadTokenFromFile()
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}

		url := clients.Clients.Gateway + "/apis/core/v1/nodes/available"
		if query != "" {
			url += "/query_match"
			var nodeQueries []model.NodeQuery
			if err := json.Unmarshal([]byte(query), &nodeQueries); err != nil {
				fmt.Printf("Error parsing query JSON: %v\n", err)
				os.Exit(1)
			}
			request := map[string][]model.NodeQuery{"query": nodeQueries}
			err = sendRequest(url, token, request)
		} else {
			err = sendRequest(url, token, nil)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	NodesCmd.Flags().StringVarP(&query, "query", "q", "", "Query JSON for node allocation")
}

func sendRequest(url, token string, body interface{}) error {
	var req *http.Request
	var err error

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}
		req, err = http.NewRequest("GET", url, bytes.NewBuffer(bodyJSON))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
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
