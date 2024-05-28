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
	deleteNodeLabelsShortDesc = "Delete a label from a node"
	deleteNodeLabelsLongDesc  = "Delete a specific label from a node using its key.\n" +
		"Provide the node ID, organization, and label key to remove the label.\n\n" +
		"Example:\n" +
		"cockpit delete label --nodeId 'nodeID' --org 'org' --key 'labelKey' "

	// Flag Constants
	flagNodeID = "nodeId"
	flagKey    = "key"

	// Flag Shorthand Constants
	shortFlagNodeID = "n"
	shortFlagKey    = "k"

	// Flag Descriptions
	descNodeID = "Node ID (required)"
	descKey    = "Label key (required)"
)

var (
	nodeId       string
	org          string
	key          string
	nodeResponse model.NodeResponse
)

var DeleteNodeLabelsCmd = &cobra.Command{
	Use:   "label",
	Short: deleteNodeLabelsShortDesc,
	Long:  deleteNodeLabelsLongDesc,
	Run:   executeDeleteLabel,
}

func executeDeleteLabel(cmd *cobra.Command, args []string) {
	err := sendDeleteLabelRequest()

	if err != nil {
		fmt.Printf("Error deleting label: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	render.RenderNode(nodeResponse.Node)
	fmt.Println("Label deleted successfully.")
	println()
}

func prepareLabelRequest() model.DeleteLabelInput {
	return model.DeleteLabelInput{
		LabelKey: key,
		NodeID:   nodeId,
		Org:      org,
	}
}

func sendDeleteLabelRequest() error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	deleteLabelURL := clients.BuildURL("core", "v1", "DeleteLabel")

	input := prepareLabelRequest()

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         deleteLabelURL,
		Method:      "DELETE",
		RequestBody: input,
		Token:       token,
		Response:    &nodeResponse,
		Timeout:     10 * time.Second,
	})
}

func init() {
	DeleteNodeLabelsCmd.Flags().StringVarP(&nodeId, flagNodeID, shortFlagNodeID, "", descNodeID)
	DeleteNodeLabelsCmd.Flags().StringVarP(&org, flagOrganization, shortFlagOrganization, "", descOrganization)
	DeleteNodeLabelsCmd.Flags().StringVarP(&key, flagKey, shortFlagKey, "", descKey)

	DeleteNodeLabelsCmd.MarkFlagRequired(flagNodeID)
	DeleteNodeLabelsCmd.MarkFlagRequired(flagOrganization)
	DeleteNodeLabelsCmd.MarkFlagRequired(flagKey)
}
