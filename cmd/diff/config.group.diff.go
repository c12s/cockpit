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
	organization string
	namespace    string
	names        string
	versions     string
	outputFormat string
	diffResponse model.ConfigGroupDiffResponse
)

var DiffConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.ConfigAliases,
	Short:   constants.DiffConfigGroupShortDesc,
	Long:    constants.DiffConfigGroupLongDesc,
	Run:     executeDiffConfigGroup,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NamesFlag, constants.VersionsFlag})
	},
}

func executeDiffConfigGroup(cmd *cobra.Command, args []string) {
	requestBody, err := utils.PrepareConfigDiffRequest(namespace, names, versions, organization)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendDiffRequest(requestBody); err != nil {
		fmt.Println("Error sending config group diff request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(diffResponse)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(diffResponse, outputFormat, "")

		filePath := constants.DiffConfigFilePathYAML
		if outputFormat == "json" {
			filePath = constants.DiffConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&diffResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
}

func sendDiffRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "DiffConfigGroup")

	config := model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &diffResponse,
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}

	return nil
}

func init() {
	DiffConfigGroupCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DiffConfigGroupCmd.Flags().StringVarP(&names, constants.NamesFlag, constants.NamesShorthandFlag, "", constants.ConfigDiffNamesDescription)
	DiffConfigGroupCmd.Flags().StringVarP(&versions, constants.VersionsFlag, constants.VersionsShorthandFlag, "", constants.ConfigDiffVersionsDescription)
	DiffConfigGroupCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	DiffConfigGroupCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	DiffConfigGroupCmd.MarkFlagRequired(constants.NamespaceFlag)
	DiffConfigGroupCmd.MarkFlagRequired(constants.OrganizationFlag)
	DiffConfigGroupCmd.MarkFlagRequired(constants.NamesFlag)
	DiffConfigGroupCmd.MarkFlagRequired(constants.VersionsFlag)
}
