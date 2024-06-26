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
	getStandaloneConfigShortDesc = "Retrieve and display a standalone configuration"
	getStandaloneConfigLongDesc  = `This command retrieves a standalone configuration specified by its name, organization, and version, displays it in a nicely formatted way, and saves it to both YAML and JSON files.
The user can specify the organization, standalone configuration name, and version to retrieve the configuration details. The response can be formatted as either YAML or JSON based on user preference.

Example:
- cockpit get standalone config --org 'org' --name 'db_config' --version 'v1.0.0'`

	// Path to files
	getStandaloneConfigFilePathJSON = "./response/standalone-config/standalone-config.json"
	getStandaloneConfigFilePathYAML = "./response/standalone-config/standalone-config.yaml"
)

var (
	standaloneConfigResponse model.StandaloneConfig
)

var GetStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   getStandaloneConfigShortDesc,
	Long:    getStandaloneConfigLongDesc,
	Run:     executeGetStandaloneConfig,
}

func executeGetStandaloneConfig(cmd *cobra.Command, args []string) {
	requestBody := prepareStandaloneRequestConfig()

	if err := sendStandaloneRequest(requestBody); err != nil {
		fmt.Println("Error sending standalone config request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(standaloneConfigResponse)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&standaloneConfigResponse, outputFormat, "")

		filePath := getStandaloneConfigFilePathYAML
		if outputFormat == "json" {
			filePath = getStandaloneConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&standaloneConfigResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
	fmt.Println()
}

func prepareStandaloneRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}

	return requestBody
}

func sendStandaloneRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetStandaloneConfig")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &standaloneConfigResponse,
		Timeout:     10 * time.Second,
	})
}
func init() {
	GetStandaloneConfigCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	GetStandaloneConfigCmd.Flags().StringVarP(&name, nameFlag, nameShorthandFlag, "", nameDescription)
	GetStandaloneConfigCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)
	GetStandaloneConfigCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	GetStandaloneConfigCmd.MarkFlagRequired(organizationFlag)
	GetStandaloneConfigCmd.MarkFlagRequired(nameFlag)
	GetStandaloneConfigCmd.MarkFlagRequired(versionFlag)
}
