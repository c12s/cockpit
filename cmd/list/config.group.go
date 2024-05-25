package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	listConfigGroupShortDesc = "Retrieve and display the configuration groups"
	listConfigGroupLongDesc  = "This command retrieves the configuration groups from a specified endpoint\n" +
		"displays them in a nicely formatted way, and saves them to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"list-config-group --organization 'org'"

	// Flag Constants
	outputFlag = "output"

	// Flag Shorthand Constants
	outputFlagShortHand = "o"

	// Flag Descriptions
	outputDesc = "Output format (yaml or json)"

	// Path to files
	listConfigFilePathJSON = "./config_group_files/list-config.json"
	listConfigFilePathYAML = "./config_group_files/list-config.yaml"
)

var (
	organization        string
	configGroupResponse model.ConfigGroupsResponse
	outputFormat        string
)

var ListConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: listConfigGroupShortDesc,
	Long:  listConfigGroupLongDesc,
	Run:   executeListConfigGroup,
}

func executeListConfigGroup(cmd *cobra.Command, args []string) {
	config := createListRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.HandleConfigGroupResponse(config.Response.(*model.ConfigGroupsResponse), outputFormat)

	err = saveConfigGroupResponseToFiles(config.Response.(*model.ConfigGroupsResponse))
	if err != nil {
		log.Fatalf("Failed to save response to files: %v", err)
	}
}

func createListRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListConfigGroup")

	requestBody := map[string]string{
		"organization": organization,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &configGroupResponse,
	}
}

func saveConfigGroupResponseToFiles(response *model.ConfigGroupsResponse) error {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to convert to JSON: %v", err)
		}
		err = ioutil.WriteFile(listConfigFilePathJSON, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write JSON file: %v", err)
		}
		fmt.Printf("Config group saved to %s\n", listConfigFilePathJSON)
	} else {
		yamlData, err := utils.MarshalConfigGroupResponseToYAML(response)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}
		err = ioutil.WriteFile(listConfigFilePathYAML, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write YAML file: %v", err)
		}
		fmt.Printf("Config group saved to %s\n", listConfigFilePathYAML)
	}

	return nil
}

func init() {
	ListConfigGroupCmd.Flags().StringVarP(&organization, orgFlag, orgFlagShortHand, "", orgDesc)
	ListConfigGroupCmd.Flags().StringVarP(&outputFormat, outputFlag, outputFlagShortHand, "yaml", outputDesc)

	ListConfigGroupCmd.MarkFlagRequired(orgFlag)
}
