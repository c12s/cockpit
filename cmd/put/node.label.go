package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
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
	Aliases: aliases.LabelAliases,
	Short:   constants.ShortLabelDesc,
	Long:    constants.LongLabelDesc,
	Run:     executeLabelCommand,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.KeyFlag, constants.ValueFlag, constants.NodeIdFlag, constants.OrganizationFlag})
	},
}

func executeLabelCommand(cmd *cobra.Command, args []string) {
	originalValue := value

	formattedValue, url, err := determineValueTypeAndURL(value)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	labelInput := createLabelInput(key, formattedValue, nodeId, org)
	err = sendLabelRequest(labelInput, url)
	if err != nil {
		fmt.Println("Error sending add node label request:", err)
		os.Exit(1)
	}

	fmt.Printf("Label %s with value %s: added or updated successfully.\n", key, originalValue)
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
		Timeout:     30 * time.Second,
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
	LabelsCmd.Flags().StringVarP(&key, constants.KeyFlag, constants.KeyShorthandFlag, "", constants.LabelKeyDescription)
	LabelsCmd.Flags().StringVarP(&value, constants.ValueFlag, constants.ValueShorthandFlag, "", constants.LabelValueDescription)
	LabelsCmd.Flags().StringVarP(&nodeId, constants.NodeIdFlag, constants.NodeIdShorthandFlag, "", constants.NodeIdDescription)
	LabelsCmd.Flags().StringVarP(&org, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)

	LabelsCmd.MarkFlagRequired(constants.KeyFlag)
	LabelsCmd.MarkFlagRequired(constants.ValueFlag)
	LabelsCmd.MarkFlagRequired(constants.NodeIdFlag)
	LabelsCmd.MarkFlagRequired(constants.OrganizationFlag)
}
