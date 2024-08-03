package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	deleteStandaloneConfigResponse model.StandaloneConfig
)

var DeleteStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   constants.DeleteStandaloneConfigShortDesc,
	Long:    constants.DeleteStandaloneConfigLongDesc,
	Run:     executeDeleteStandaloneConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NameFlag, constants.VersionFlag})
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
}

func prepareDeleteStandaloneConfigRequestConfig() model.SingleConfigReference {
	requestBody := model.SingleConfigReference{
		Organization: organization,
		Namespace:    namespace,
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
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &deleteStandaloneConfigResponse,
	}
}

func init() {
	DeleteStandaloneConfigCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	DeleteStandaloneConfigCmd.MarkFlagRequired(constants.NamespaceFlag)
	DeleteStandaloneConfigCmd.MarkFlagRequired(constants.OrganizationFlag)
	DeleteStandaloneConfigCmd.MarkFlagRequired(constants.NameFlag)
	DeleteStandaloneConfigCmd.MarkFlagRequired(constants.VersionFlag)
}
