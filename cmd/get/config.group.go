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
	name                string
	configGroupResponse model.ConfigGroup
	outputFormat        string
)

var GetSingleConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   constants.GetAppConfigShortDesc,
	Long:    constants.GetAppConfigLongDesc,
	Run:     executeGetAppConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NameFlag, constants.VersionFlag})
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

		filePath := constants.GetConfigGroupFilePathYAML
		if outputFormat == "json" {
			filePath = constants.GetConfigGroupFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&configGroupResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
}

func prepareRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Organization: organization,
		Namespace:    namespace,
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
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetSingleConfigGroupCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	GetSingleConfigGroupCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	GetSingleConfigGroupCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	GetSingleConfigGroupCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	GetSingleConfigGroupCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	GetSingleConfigGroupCmd.MarkFlagRequired(constants.NamespaceFlag)
	GetSingleConfigGroupCmd.MarkFlagRequired(constants.OrganizationFlag)
	GetSingleConfigGroupCmd.MarkFlagRequired(constants.NameFlag)
	GetSingleConfigGroupCmd.MarkFlagRequired(constants.VersionFlag)
}
