package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	listConfigGroupShortDesc = "Retrieve and display the configuration groups"
	listConfigGroupLongDesc  = `This command retrieves all configuration groups from a specified organization,
displays them in a nicely formatted way, and saves them to both YAML and JSON files.
You can choose the output format by specifying either 'yaml' or 'json'.

Examples:
- cockpit list config group --organization 'org' --output 'json'
- cockpit list config group --organization 'org' --output 'yaml'`

	// Flag Constants
	outputFlag = "output"

	// Flag Shorthand Constants
	outputShorthandFlag = "o"

	// Flag Descriptions
	outputDescription = "Output format (yaml or json)"

	// Path to files
	listConfigFilePathJSON = "./response/config-group/list-config.json"
	listConfigFilePathYAML = "./response/config-group/list-config.yaml"
)

var (
	organization        string
	configGroupResponse model.ConfigGroupsResponse
	outputFormat        string
)

var ListConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   listConfigGroupShortDesc,
	Long:    listConfigGroupLongDesc,
	Run:     executeListConfigGroup,
}

func executeListConfigGroup(cmd *cobra.Command, args []string) {
	config := createListRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		fmt.Println("Error sending config group request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(configGroupResponse.Groups)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&configGroupResponse, outputFormat, "")

		filePath := listConfigFilePathYAML
		if outputFormat == "json" {
			filePath = listConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&configGroupResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
	fmt.Println()
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

func init() {
	ListConfigGroupCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	ListConfigGroupCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	ListConfigGroupCmd.MarkFlagRequired(organizationFlag)
}
