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
	standaloneConfigResponse model.StandaloneConfig
)

var GetStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   constants.GetStandaloneConfigShortDesc,
	Long:    constants.GetStandaloneConfigLongDesc,
	Run:     executeGetStandaloneConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NameFlag, constants.VersionFlag})
	},
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

		filePath := constants.GetStandaloneConfigFilePathYAML
		if outputFormat == "json" {
			filePath = constants.GetStandaloneConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&standaloneConfigResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
}

func prepareStandaloneRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Organization: organization,
		Namespace:    namespace,
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
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetStandaloneConfigCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	GetStandaloneConfigCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	GetStandaloneConfigCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	GetStandaloneConfigCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	GetStandaloneConfigCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	GetStandaloneConfigCmd.MarkFlagRequired(constants.NamespaceFlag)
	GetStandaloneConfigCmd.MarkFlagRequired(constants.OrganizationFlag)
	GetStandaloneConfigCmd.MarkFlagRequired(constants.NameFlag)
	GetStandaloneConfigCmd.MarkFlagRequired(constants.VersionFlag)
}
