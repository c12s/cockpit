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
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	putConfigGroupShortDesc = "Send a configuration group to the server"
	putConfigGroupLongDesc  = "This command sends a configuration group read from a file (JSON or YAML)\n" +
		"to the server and displays the server's response in JSON format.\n\n" +
		"Example:\n" +
		"put-config-group --path 'path to yaml or JSON file'"
)

var (
	filePath    string
	putResponse model.ConfigGroup
)

var PutConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: putConfigGroupShortDesc,
	Long:  putConfigGroupLongDesc,
	Run:   executePutConfigGroup,
}

func executePutConfigGroup(cmd *cobra.Command, args []string) {
	configData, err := readConfigFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	config := createPutRequestConfig(configData)

	err = utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	displayResponseAsJSON(&putResponse)
}

func readConfigFile(path string) (map[string]interface{}, error) {
	var configData map[string]interface{}

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if strings.HasSuffix(path, ".yaml") {
		err = yaml.Unmarshal(fileContent, &configData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
		}
	} else if strings.HasSuffix(path, ".json") {
		err = json.Unmarshal(fileContent, &configData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file format")
	}

	return configData, nil
}

func createPutRequestConfig(configData map[string]interface{}) model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "PutConfigGroup")

	return model.HTTPRequestConfig{
		Method:      "POST",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: configData,
		Response:    &putResponse,
	}
}

func displayResponseAsJSON(response *model.ConfigGroup) {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("Error converting response to JSON: %v\n", err)
		return
	}
	fmt.Println("Config Group Response (JSON):")
	fmt.Println(string(jsonData))
}

func init() {
	PutConfigGroupCmd.Flags().StringVarP(&filePath, "path", "p", "", "Path to the configuration file (required)")
	PutConfigGroupCmd.MarkFlagRequired("path")
}
