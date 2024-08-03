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
	standaloneConfigDiffResponse model.StandaloneConfigDiffResponse
)

var DiffStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.GroupAliases,
	Short:   constants.DiffStandaloneConfigShortDesc,
	Long:    constants.DiffStandaloneConfigLongDesc,
	Run:     executeDiffStandaloneConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NamesFlag, constants.VersionsFlag})
	},
}

func executeDiffStandaloneConfig(cmd *cobra.Command, args []string) {
	requestBody, err := utils.PrepareConfigDiffRequest(namespace, names, versions, organization)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendStandaloneDiffRequest(requestBody); err != nil {
		fmt.Println("Error sending standalone config diff request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(standaloneConfigDiffResponse)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&standaloneConfigDiffResponse, outputFormat, "")

		filePath := constants.DiffStandaloneConfigFilePathYAML
		if outputFormat == "json" {
			filePath = constants.DiffStandaloneConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&standaloneConfigDiffResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
}

func sendStandaloneDiffRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "DiffStandaloneConfig")

	config := model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &standaloneConfigDiffResponse,
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}

	return nil
}

func init() {
	DiffStandaloneConfigCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DiffStandaloneConfigCmd.Flags().StringVarP(&names, constants.NamesFlag, constants.NamesShorthandFlag, "", constants.ConfigDiffNamesDescription)
	DiffStandaloneConfigCmd.Flags().StringVarP(&versions, constants.VersionsFlag, constants.VersionsShorthandFlag, "", constants.ConfigDiffVersionsDescription)
	DiffStandaloneConfigCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	DiffStandaloneConfigCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	DiffStandaloneConfigCmd.MarkFlagRequired(constants.NamespaceFlag)
	DiffStandaloneConfigCmd.MarkFlagRequired(constants.OrganizationFlag)
	DiffStandaloneConfigCmd.MarkFlagRequired(constants.NamesFlag)
	DiffStandaloneConfigCmd.MarkFlagRequired(constants.VersionsFlag)
}
