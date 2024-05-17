package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	shortLabelDescription = "Add a label to a node."
	longLabelDescription  = "This command allows you to add a new label to a specified node, enhancing node metadata. \n" +
		"Provide a key-value pair to define the label. If the label already exists, its value will be updated to the new specified value."

	tokenFilePath   = "token.txt"
	contentTypeJSON = "application/json"
	authHeader      = "Authorization"
	bearer          = "Bearer "

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
	descKey    = "Label key"
	descValue  = "Label value"
	descNodeID = "Node ID"
	descOrg    = "Organization"
)

func getToken() (string, error) {
	token, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}
	return string(token), nil
}

var LabelsCmd = &cobra.Command{
	Use:   "label",
	Short: shortLabelDescription,
	Long:  longLabelDescription,
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		valueStr, _ := cmd.Flags().GetString("value")
		nodeId, _ := cmd.Flags().GetString("nodeId")
		org, _ := cmd.Flags().GetString("org")

		token, err := getToken()
		if err != nil {
			fmt.Println("Error getting token:", err)
			return
		}

		value, url, err := determineValueTypeAndURL(valueStr)
		if err != nil {
			fmt.Println("Error determining value type:", err)
			return
		}

		labelInput := createLabelInput(key, value, nodeId, org)

		if err := sendLabelRequest(labelInput, url, token); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Label added or updated successfully.")
	},
}

func determineValueTypeAndURL(valueStr string) (interface{}, string, error) {
	labelEndpointBase := clients.Clients.Gateway + "/apis/core/v1/labels"

	if floatValue, err := strconv.ParseFloat(valueStr, 64); err == nil {
		labelEndpointFloat64 := labelEndpointBase + "/float64"
		return floatValue, labelEndpointFloat64, nil
	}
	if boolValue, err := strconv.ParseBool(valueStr); err == nil {
		labelEndpointBool := labelEndpointBase + "/bool"
		return boolValue, labelEndpointBool, nil
	}
	labelEndpointString := labelEndpointBase + "/string"
	return valueStr, labelEndpointString, nil
}

func createLabelInput(key string, value interface{}, nodeId string, org string) model.LabelInput {
	return model.LabelInput{
		Label: struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value,omitempty"`
		}{
			Key:   key,
			Value: value,
		},
		NodeID: nodeId,
		Org:    org,
	}
}

func sendLabelRequest(input model.LabelInput, url string, token string) error {
	payload, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error creating JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set(authHeader, bearer+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		return fmt.Errorf("error response from server: %s", body)
	}

	return nil
}

func init() {
	LabelsCmd.Flags().StringP(flagKey, flagShorthandKey, "", descKey)
	LabelsCmd.Flags().StringP(flagValue, flagShorthandValue, "", descValue)
	LabelsCmd.Flags().StringP(flagNodeID, flagShorthandNodeID, "", descNodeID)
	LabelsCmd.Flags().StringP(flagOrg, flagShorthandOrg, "", descOrg)

	LabelsCmd.MarkFlagRequired(flagKey)
	LabelsCmd.MarkFlagRequired(flagValue)
	LabelsCmd.MarkFlagRequired(flagNodeID)
	LabelsCmd.MarkFlagRequired(flagOrg)
}
