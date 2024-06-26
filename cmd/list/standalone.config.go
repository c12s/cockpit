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
	listStandaloneConfigShortDesc = "List all standalone configurations"
	listStandaloneConfigLongDesc  = `This command retrieves a list of all standalone configurations for a given organization,
displays them in a nicely formatted way, and saves them to both YAML and JSON files.
You can choose the output format by specifying either 'yaml' or 'json'.

Examples:
- cockpit list standalone config --org 'org' --output 'json'
- cockpit list standalone config --org 'org' --output 'yaml'`

	listStandaloneConfigFilePathJSON = "./response/standalone-config/list-standalone-config.json"
	listStandaloneConfigFilePathYAML = "./response/standalone-config/list-standalone-config.yaml"
)

var (
	listStandaloneConfigResponse model.StandaloneConfigsResponse
)

var ListStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   listStandaloneConfigShortDesc,
	Long:    listStandaloneConfigLongDesc,
	Run:     executeListStandaloneConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{organizationFlag})
	},
}

func executeListStandaloneConfig(cmd *cobra.Command, args []string) {
	config := createListStandaloneRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		fmt.Println("Error sending standalone config request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(listStandaloneConfigResponse.Configurations)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&listStandaloneConfigResponse, outputFormat, "")

		filePath := listStandaloneConfigFilePathYAML
		if outputFormat == "json" {
			filePath = listStandaloneConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&listStandaloneConfigResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
	fmt.Println()
}

func createListStandaloneRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListStandaloneConfig")

	requestBody := map[string]string{
		"organization": organization,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &listStandaloneConfigResponse,
	}
}

func init() {
	ListStandaloneConfigCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	ListStandaloneConfigCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	ListStandaloneConfigCmd.MarkFlagRequired(organizationFlag)
}
