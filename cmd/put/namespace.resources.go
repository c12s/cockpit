package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var PutNamespaceResourcesCmd = &cobra.Command{
	Use: "namespace resources",
	Run: executePutNamespaceResources,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
}

func executePutNamespaceResources(cmd *cobra.Command, args []string) {
	requestBody, err := prepareNamespaceRequestBody()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendPutNamespaceResourcesRequest(requestBody); err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}

	fmt.Println("App resources updated successfully")
}

func prepareNamespaceRequestBody() (map[string]any, error) {
	var requestBody map[string]any
	var err error
	if strings.HasSuffix(filePath, ".yaml") {
		err = utils.ReadYAML(filePath, &requestBody)
	} else if strings.HasSuffix(filePath, ".json") {
		err = utils.ReadJSON(filePath, &requestBody)
	} else {
		return requestBody, fmt.Errorf("invalid file format. Please provide a YAML or JSON file")
	}

	if err != nil {
		return requestBody, fmt.Errorf("failed to read input file: %v", err)
	}

	return requestBody, nil
}

func sendPutNamespaceResourcesRequest(requestBody map[string]any) error {
	config, err := prepareNamespaceRequestConfig(requestBody)
	if err != nil {
		return fmt.Errorf("error creating request config: %v", err)
	}

	err = utils.SendHTTPRequest(config)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	return nil
}

func prepareNamespaceRequestConfig(requestBody map[string]any) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "SetNamespaceResources")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: requestBody,
		Token:       token,
		Timeout:     30 * time.Second,
	}, nil
}

func init() {
	PutNamespaceResourcesCmd.Flags().StringVarP(&filePath, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)
	PutNamespaceResourcesCmd.MarkFlagRequired(constants.FilePathFlag)
}
