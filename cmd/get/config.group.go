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
	getAppConfigShortDesc = "Retrieve and display the configuration"
	getAppConfigLongDesc  = `This command retrieves a specific configuration by its organization, name, and version, displays it in a nicely formatted way, and saves it to both YAML and JSON files.
The user can specify the organization, configuration name, and version to retrieve the configuration details. The response can be formatted as either YAML or JSON based on user preference.

Example:
- cockpit get config group --org 'org' --name 'app_config' --version 'v1.0.0'`

	// Flag Constants
	nameFlag   = "name"
	outputFlag = "output"

	// Flag Shorthand Constants
	nameShorthandFlag   = "n"
	outputShorthandFlag = "o"

	// Flag Descriptions
	nameDescription   = "Configuration name (required)"
	outputDescription = "Output format (yaml or json)"

	// Path to files
	getConfigFilePathJSON = "./response/config-group/single-config.json"
	getConfigFilePathYAML = "./response/config-group/single-config.yaml"
)

var (
	name                string
	configGroupResponse model.ConfigGroup
	outputFormat        string
)
var GetSingleConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   getAppConfigShortDesc,
	Long:    getAppConfigLongDesc,
	Run:     executeGetAppConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{organizationFlag, nameFlag, versionFlag})
	},
}

func executeGetAppConfig(cmd *cobra.Command, args []string) {
	requestBody := prepareRequestConfig()

	if err := sendConfigGroupRequest(requestBody); err != nil {
		fmt.Println("Error sending config group request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(configGroupResponse)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&configGroupResponse, outputFormat, "")

		filePath := getConfigFilePathYAML
		if outputFormat == "json" {
			filePath = getConfigFilePathJSON
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

func prepareRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}

	return requestBody
}

func sendConfigGroupRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetConfigGroup")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &configGroupResponse,
		Timeout:     10 * time.Second,
	})
}

func init() {
	GetSingleConfigGroupCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	GetSingleConfigGroupCmd.Flags().StringVarP(&name, nameFlag, nameShorthandFlag, "", nameDescription)
	GetSingleConfigGroupCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)
	GetSingleConfigGroupCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	GetSingleConfigGroupCmd.MarkFlagRequired(organizationFlag)
	GetSingleConfigGroupCmd.MarkFlagRequired(nameFlag)
	GetSingleConfigGroupCmd.MarkFlagRequired(versionFlag)
}
