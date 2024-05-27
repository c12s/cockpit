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
	putConfigGroupShortDesc = "Send a configuration group to the server"
	putConfigGroupLongDesc  = "This command sends a configuration group read from a file (JSON or YAML)\n" +
		"to the server and displays the server's response in the same format as the input file.\n\n" +
		"Example:\n" +
		"put-config-group --path 'path to yaml or JSON file'"
)

var (
	filePath    string
	putResponse model.ConfigGroup
	inputFormat string
)

var PutConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: putConfigGroupShortDesc,
	Long:  putConfigGroupLongDesc,
	Run:   executePutConfigGroup,
}

func executePutConfigGroup(cmd *cobra.Command, args []string) {
	configData, err := prepareConfigGroupData(filePath)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	if err := sendConfigGroupData(configData, &putResponse); err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	displayConfigGroupResponse(&putResponse, inputFormat)
}

func prepareConfigGroupData(path string) (map[string]interface{}, error) {
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

func sendConfigGroupData(requestBody interface{}, response interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "PutConfigGroup")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "POST",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    response,
	})
}

func displayConfigGroupResponse(response *model.ConfigGroup, format string) {
	if format == "json" {
		utils.DisplayResponseAsJSON(response, "Config Group Response (JSON):")
	} else {
		utils.DisplayResponseAsYAML(response, "Config Group Response (YAML):")
	}
}

func init() {
	PutConfigGroupCmd.Flags().StringVarP(&filePath, "path", "p", "", "Path to the configuration file (required)")
	PutConfigGroupCmd.MarkFlagRequired("path")
}
