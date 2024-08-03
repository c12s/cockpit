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
	organization        string
	namespace           string
	configGroupResponse model.ConfigGroupsResponse
	outputFormat        string
)

var ListConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   constants.ListStandaloneConfigShortDesc,
	Long:    constants.ListStandaloneConfigLongDesc,
	Run:     executeListConfigGroup,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag})
	},
}

func executeListConfigGroup(cmd *cobra.Command, args []string) {
	config := createListRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		fmt.Println("Error sending config group request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(configGroupResponse.Groups)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&configGroupResponse, outputFormat, "")

		filePath := constants.ListConfigGroupFilePathYAML
		if outputFormat == "json" {
			filePath = constants.ListConfigGroupFilePathJSON
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

func createListRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListConfigGroup")

	requestBody := map[string]string{
		"organization": organization,
		"namespace":    namespace,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &configGroupResponse,
	}
}

func init() {
	ListConfigGroupCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ListConfigGroupCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	ListConfigGroupCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamesShorthandFlag, "", constants.NamespaceDescription)

	ListConfigGroupCmd.MarkFlagRequired(constants.NamespaceFlag)
	ListConfigGroupCmd.MarkFlagRequired(constants.OrganizationFlag)
	ListConfigGroupCmd.MarkFlagRequired(constants.NamespaceFlag)
}
