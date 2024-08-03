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

var PutAppResourcesCmd = &cobra.Command{
	Use: "app resources",
	Run: executePutAppResources,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
}

func executePutAppResources(cmd *cobra.Command, args []string) {
	requestBody, err := prepareAppRequestBody()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendPutAppResourcesRequest(requestBody); err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}

	fmt.Println("App resources updated successfully")
}

func prepareAppRequestBody() (map[string]any, error) {
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

func sendPutAppResourcesRequest(requestBody map[string]any) error {
	config, err := prepareAppRequestConfig(requestBody)
	if err != nil {
		return fmt.Errorf("error creating request config: %v", err)
	}

	err = utils.SendHTTPRequest(config)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	return nil
}

func prepareAppRequestConfig(requestBody map[string]any) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "SetAppResources")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: requestBody,
		Token:       token,
		Timeout:     30 * time.Second,
	}, nil
}

func init() {
	PutAppResourcesCmd.Flags().StringVarP(&filePath, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)
	PutAppResourcesCmd.MarkFlagRequired(constants.FilePathFlag)
}
