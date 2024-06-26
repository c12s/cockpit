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
	deleteNodeLabelsShortDesc = "Delete a label from a node"
	deleteNodeLabelsLongDesc  = `Delete a specific label from a node using its key.
This command allows the user to remove a label from a node by specifying the node ID, organization, and the label key. The response includes the updated node details.

Example:
- cockpit delete label --node-id 'nodeID' --org 'org' --key 'labelKey'`

	// Flag Constants
	nodeIdFlag = "node-id"
	keyFlag    = "key"

	// Flag Shorthand Constants
	flagShorthandFlag = "n"
	keyShorthandFlag  = "k"

	// Flag Descriptions
	nodeIdDescription = "Node ID (required)"
	keyDescription    = "Label key (required)"
)

var (
	nodeId       string
	org          string
	key          string
	nodeResponse model.NodeResponse
)

var DeleteNodeLabelsCmd = &cobra.Command{
	Use:     "label",
	Aliases: aliases.LabelAliases,
	Short:   deleteNodeLabelsShortDesc,
	Long:    deleteNodeLabelsLongDesc,
	Run:     executeDeleteLabel,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{nodeIdFlag, organizationFlag, keyFlag})
	},
}

func executeDeleteLabel(cmd *cobra.Command, args []string) {
	err := sendDeleteLabelRequest()

	if err != nil {
		fmt.Println("Error sending delete node label request:", err)
		os.Exit(1)
	}

	render.RenderNode(nodeResponse.Node)
	println()
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
	DeleteNodeLabelsCmd.Flags().StringVarP(&nodeId, nodeIdFlag, flagShorthandFlag, "", nodeIdDescription)
	DeleteNodeLabelsCmd.Flags().StringVarP(&org, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	DeleteNodeLabelsCmd.Flags().StringVarP(&key, keyFlag, keyShorthandFlag, "", keyDescription)

	DeleteNodeLabelsCmd.MarkFlagRequired(nodeIdFlag)
	DeleteNodeLabelsCmd.MarkFlagRequired(organizationFlag)
	DeleteNodeLabelsCmd.MarkFlagRequired(keyFlag)
}
