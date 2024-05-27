package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	putStandaloneConfigShortDesc = "Send a standalone configuration to the server"
	putStandaloneConfigLongDesc  = "This command sends a standalone configuration read from a file (JSON or YAML)\n" +
		"to the server and displays the server's response in the same format as the input file.\n\n" +
		"Example:\n" +
		"put-standalone-config --path 'path to yaml or JSON file'"
)

var (
	standaloneConfigPutResponse model.StandaloneConfig
)
var PutStandaloneConfigCmd = &cobra.Command{
	Use:   "config",
	Short: putStandaloneConfigShortDesc,
	Long:  putStandaloneConfigLongDesc,
	Run:   executePutStandaloneConfig,
}

func executePutStandaloneConfig(cmd *cobra.Command, args []string) {
	configData, err := prepareStandaloneConfigData(filePath)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	if err := sendStandaloneConfigData(configData, &standaloneConfigPutResponse); err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	displayStandaloneConfigResponse(&standaloneConfigPutResponse, inputFormat)
}

func prepareStandaloneConfigData(path string) (map[string]interface{}, error) {
	var configData map[string]interface{}

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if strings.HasSuffix(path, ".yaml") {
		inputFormat = "yaml"
		err = yaml.Unmarshal(fileContent, &configData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
		}
	} else if strings.HasSuffix(path, ".json") {
		inputFormat = "json"
		err = json.Unmarshal(fileContent, &configData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file format")
	}

	return configData, nil
}

func sendStandaloneConfigData(requestBody interface{}, response interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "PutStandaloneConfig")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "POST",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    response,
	})
}

func displayStandaloneConfigResponse(response *model.StandaloneConfig, format string) {
	if format == "json" {
		utils.DisplayResponseAsJSON(response, "Standalone Config Response (JSON):")
	} else {
		utils.DisplayResponseAsYAML(response, "Standalone Config Response (YAML):")
	}
}

func init() {
	PutStandaloneConfigCmd.Flags().StringVarP(&filePath, "path", "p", "", "Path to the configuration file (required)")
	PutStandaloneConfigCmd.MarkFlagRequired("path")
}
