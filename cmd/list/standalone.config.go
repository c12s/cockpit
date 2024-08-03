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
	listStandaloneConfigResponse model.StandaloneConfigsResponse
)

var ListStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   constants.ListStandaloneConfigShortDesc,
	Long:    constants.ListStandaloneConfigLongDesc,
	Run:     executeListStandaloneConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag})
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

		filePath := constants.ListStandaloneConfigFilePathYAML
		if outputFormat == "json" {
			filePath = constants.ListStandaloneConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&listStandaloneConfigResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
}

func createListStandaloneRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListStandaloneConfig")

	requestBody := map[string]string{
		"namespace":    namespace,
		"organization": organization,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &listStandaloneConfigResponse,
	}
}

func init() {
	ListStandaloneConfigCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ListStandaloneConfigCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	ListStandaloneConfigCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamesShorthandFlag, "", constants.NamespaceDescription)

	ListStandaloneConfigCmd.MarkFlagRequired(constants.OrganizationFlag)
	ListStandaloneConfigCmd.MarkFlagRequired(constants.NamespaceFlag)
}
