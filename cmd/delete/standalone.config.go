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
	deleteStandaloneConfigShortDesc = "Delete a standalone configuration version"
	deleteStandaloneConfigLongDesc  = `This command deletes a specified standalone configuration version and displays the deleted configuration details in JSON format.
The user can specify the organization, standalone configuration name, and version to delete the configuration. The output can be formatted as either JSON or YAML based on user preference.

Example:
- cockpit delete standalone config --org 'c12s' --name 'db_config' --version 'v1.0.1'`
)

var (
	deleteStandaloneConfigResponse model.StandaloneConfig
)

var DeleteStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   deleteStandaloneConfigShortDesc,
	Long:    deleteStandaloneConfigLongDesc,
	Run:     executeDeleteStandaloneConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{organizationFlag, nameFlag, versionFlag})
	},
}

func executeDeleteStandaloneConfig(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteStandaloneConfigRequestConfig()

	config := sendDeleteStandaloneConfigRequestConfig(requestBody)

	err := utils.SendHTTPRequest(config)
	if err != nil {
		fmt.Println("Error sending delete standalone config request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(deleteStandaloneConfigResponse)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&deleteStandaloneConfigResponse, outputFormat, "Config group deleted successfully")
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
	fmt.Println()
}

func prepareDeleteStandaloneConfigRequestConfig() model.SingleConfigReference {
	requestBody := model.SingleConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}
	return requestBody
}

func sendDeleteStandaloneConfigRequestConfig(requestBody interface{}) model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "DeleteStandaloneConfig")

	return model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &deleteStandaloneConfigResponse,
	}
}

func init() {
	DeleteStandaloneConfigCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&name, nameFlag, nameShorthandFlag, "", nameDescription)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	DeleteStandaloneConfigCmd.MarkFlagRequired(organizationFlag)
	DeleteStandaloneConfigCmd.MarkFlagRequired(nameFlag)
	DeleteStandaloneConfigCmd.MarkFlagRequired(versionFlag)
}
