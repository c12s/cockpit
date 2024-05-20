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
	longLabelDescription  = "This command allows you to add a new label to a specified node, enhancing node metadata. \n" +
		"Provide a key-value pair to define the label. If the label already exists, its value will be updated to the new specified value.\n\n" +
		"Example:\n" +
		"put-label --key \"newLabel\" --value \"value||true||25.00\" --nodeId \"nodeId\" --org \"orgId\""

	// Flag Constants
	flagKey    = "key"
	flagValue  = "value"
	flagNodeID = "nodeId"
	flagOrg    = "org"

	// Flag Shorthand Constants
	flagShorthandKey    = "k"
	flagShorthandValue  = "v"
	flagShorthandNodeID = "n"
	flagShorthandOrg    = "o"

	// Flag Descriptions
	descKey    = "Label key (required)"
	descValue  = "Label value (required)"
	descNodeID = "Node ID (required)"
	descOrg    = "Organization (required)"
)

var (
	nodeId       string
	org          string
	key          string
	value        string
	nodeResponse model.NodeResponse
)

var LabelsCmd = &cobra.Command{
	Use:   "label",
	Short: shortLabelDescription,
	Long:  longLabelDescription,
	Run:   executeLabelCommand,
}

func executeLabelCommand(cmd *cobra.Command, args []string) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	value, url, err := determineValueTypeAndURL(value)
	if err != nil {
		fmt.Printf("Error determining value type: %v\n", err)
		os.Exit(1)
	}

	labelInput := createLabelInput(key, value, nodeId, org)
	err = sendLabelRequest(labelInput, url, token)
	if err != nil {
		fmt.Printf("Error processing label: %v\n", err)
		os.Exit(1)
	}

	render.RenderNode(nodeResponse.Node)
	fmt.Println("Label added or updated successfully.")
	println()
}

func sendLabelRequest(input model.LabelInput, url, token string) error {
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
	LabelsCmd.Flags().StringVarP(&key, flagKey, flagShorthandKey, "", descKey)
	LabelsCmd.Flags().StringVarP(&value, flagValue, flagShorthandValue, "", descValue)
	LabelsCmd.Flags().StringVarP(&nodeId, flagNodeID, flagShorthandNodeID, "", descNodeID)
	LabelsCmd.Flags().StringVarP(&org, flagOrg, flagShorthandOrg, "", descOrg)

	LabelsCmd.MarkFlagRequired(flagKey)
	LabelsCmd.MarkFlagRequired(flagValue)
	LabelsCmd.MarkFlagRequired(flagNodeID)
	LabelsCmd.MarkFlagRequired(flagOrg)
}
