package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

const (
	deleteNodeLabelsShortDesc = "Delete a label from a node"
	deleteNodeLabelsLongDesc  = "Delete a specific label from a node using its key.\n" +
		"Provide the node ID, organization, and label key to remove the label.\n\n" +
		"Example:\n" +
		"delete-node-labels --nodeId <nodeID> --org <organization> --key <labelKey>"
	tokenFile   = "token.txt"
	contentType = "application/json"
	authHeader  = "Authorization"

	// Flag names
	flagNodeID = "nodeId"
	flagOrg    = "org"
	flagKey    = "key"

	// Short flag names
	shortFlagNodeID = "n"
	shortFlagOrg    = "o"
	shortFlagKey    = "k"

	// Flag descriptions
	descNodeID = "Node ID"
	descOrg    = "Organization"
	descKey    = "Label key"
)

var DeleteNodeLabelsCmd = &cobra.Command{
	Use:   "label",
	Short: deleteNodeLabelsShortDesc,
	Long:  deleteNodeLabelsLongDesc,
	Run:   executeDeleteLabel,
}

func executeDeleteLabel(cmd *cobra.Command, args []string) {
	nodeId, _ := cmd.Flags().GetString(flagNodeID)
	org, _ := cmd.Flags().GetString(flagOrg)
	key, _ := cmd.Flags().GetString(flagKey)

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		return
	}

	if err := deleteLabel(nodeId, org, key, string(token)); err != nil {
		fmt.Printf("Error deleting label: %v\n", err)
		return
	}

	fmt.Println("Label deleted successfully.")
}

func deleteLabel(nodeId, org, key, token string) error {
	apiEndpointLabels := clients.Clients.Gateway + "/apis/core/v1/labels"

	input := model.DeleteLabelInput{
		LabelKey: key,
		NodeID:   nodeId,
		Org:      org,
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error creating JSON payload: %w", err)
	}

	req, err := http.NewRequest("DELETE", apiEndpointLabels, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set(authHeader, "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		return fmt.Errorf("error response from server: %s", body)
	}

	return nil
}

func init() {
	DeleteNodeLabelsCmd.Flags().StringP(flagNodeID, shortFlagNodeID, "", descNodeID)
	DeleteNodeLabelsCmd.Flags().StringP(flagOrg, shortFlagOrg, "", descOrg)
	DeleteNodeLabelsCmd.Flags().StringP(flagKey, shortFlagKey, "", descKey)

	DeleteNodeLabelsCmd.MarkFlagRequired(flagNodeID)
	DeleteNodeLabelsCmd.MarkFlagRequired(flagOrg)
	DeleteNodeLabelsCmd.MarkFlagRequired(flagKey)
}
