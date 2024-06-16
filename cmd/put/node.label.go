package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

const (
	shortLabelDescription = "Add a label to a node."
	longLabelDescription  = `This command allows you to add a new label to a specified node, enhancing node metadata.
Provide a key-value pair to define the label. If the label already exists, its value will be updated to the new specified value.
The command supports different types of values: strings, boolean, and floating-point numbers.
The input format determines the appropriate type and URL for the request.

Examples:
- cockpit put label --key 'env' --value 'production' --nodeId 'nodeId' --org 'org'
- cockpit put label --key 'active' --value 'true' --nodeId 'nodeId' --org 'org'
- cockpit put label --key 'cpu' --value '2.5' --nodeId 'nodeId' --org 'org'`

	// Flag Constants
	keyFlag          = "key"
	valueFlag        = "value"
	nodeIdFlag       = "nodeId"
	organizationFlag = "org"

	// Flag Shorthand Constants
	keyShorthandFlag          = "k"
	valueShorthandFlag        = "v"
	nodeIdShorthandFlag       = "n"
	organizationShorthandFlag = "r"

	// Flag Descriptions
	keyDescription          = "Label key (required)"
	valueDescription        = "Label value (required)"
	nodeIdDescription       = "Node ID (required)"
	organizationDescription = "Organization (required)"
)

var (
	nodeId       string
	org          string
	key          string
	value        string
	nodeResponse model.NodeResponse
)

var LabelsCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"lbl", "lab"},
	Short:   shortLabelDescription,
	Long:    longLabelDescription,
	Run:     executeLabelCommand,
}

func executeLabelCommand(cmd *cobra.Command, args []string) {
	value, url, err := determineValueTypeAndURL(value)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	labelInput := createLabelInput(key, value, nodeId, org)
	err = sendLabelRequest(labelInput, url)
	if err != nil {
		fmt.Println("Error sending add node label request:", err)
		os.Exit(1)
	}

	render.RenderNode(nodeResponse.Node)
	fmt.Println("Label added or updated successfully.")
	println()
}

func sendLabelRequest(input model.LabelInput, url string) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: input,
		Token:       token,
		Response:    &nodeResponse,
		Timeout:     10 * time.Second,
	})
}

func determineValueTypeAndURL(valueStr string) (interface{}, string, error) {
	if floatValue, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return floatValue, clients.BuildURL("core", "v1", "PutFloat64Label"), nil
	}
	if boolValue, err := strconv.ParseBool(valueStr); err == nil {
		return boolValue, clients.BuildURL("core", "v1", "PutBoolLabel"), nil
	}
	return valueStr, clients.BuildURL("core", "v1", "PutStringLabel"), nil
}

func createLabelInput(key string, value interface{}, nodeId string, org string) model.LabelInput {
	return model.LabelInput{
		Label: model.Label{
			Key:   key,
			Value: value,
		},
		NodeID: nodeId,
		Org:    org,
	}
}

func init() {
	LabelsCmd.Flags().StringVarP(&key, keyFlag, keyShorthandFlag, "", keyDescription)
	LabelsCmd.Flags().StringVarP(&value, valueFlag, valueShorthandFlag, "", valueDescription)
	LabelsCmd.Flags().StringVarP(&nodeId, nodeIdFlag, nodeIdShorthandFlag, "", nodeIdDescription)
	LabelsCmd.Flags().StringVarP(&org, organizationFlag, organizationShorthandFlag, "", organizationDescription)

	LabelsCmd.MarkFlagRequired(keyFlag)
	LabelsCmd.MarkFlagRequired(valueFlag)
	LabelsCmd.MarkFlagRequired(nodeIdFlag)
	LabelsCmd.MarkFlagRequired(organizationFlag)
}
