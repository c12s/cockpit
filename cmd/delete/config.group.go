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
	name                      string
	outputFormat              string
	deleteConfigGroupResponse model.ConfigGroup
)

var DeleteConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   constants.DeleteConfigGroupShortDesc,
	Long:    constants.DeleteConfigGroupLongDesc,
	Run:     executeDeleteConfigGroup,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NameFlag, constants.VersionFlag})
	},
}

func executeDeleteConfigGroup(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteConfigGroupRequest()

	config := sendDeleteConfigGroupRequest(requestBody)

	err := utils.SendHTTPRequest(config)
	if err != nil {
		fmt.Println("Error sending delete config group request:", err)
		os.Exit(1)
	}

	if outputFormat == "" {
		render.RenderResponseAsTabWriter(deleteConfigGroupResponse)
	} else if outputFormat == "yaml" || outputFormat == "json" {
		render.DisplayResponseAsJSONOrYAML(&deleteConfigGroupResponse, outputFormat, "Config group deleted successfully")
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
}

func prepareDeleteConfigGroupRequest() interface{} {
	requestBody := model.ConfigReference{
		Organization: organization,
		Namespace:    namespace,
		Name:         name,
		Version:      version,
	}
	return requestBody
}

func sendDeleteConfigGroupRequest(requestBody interface{}) model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "DeleteConfigGroup")

	return model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &deleteConfigGroupResponse,
	}
}

func init() {
	DeleteConfigGroupCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DeleteConfigGroupCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	DeleteConfigGroupCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	DeleteConfigGroupCmd.Flags().StringVarP(&outputFormat, constants.OutputFlag, constants.OutputShorthandFlag, "", constants.OutputDescription)
	DeleteConfigGroupCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	DeleteConfigGroupCmd.MarkFlagRequired(constants.NamespaceFlag)
	DeleteConfigGroupCmd.MarkFlagRequired(constants.OrganizationFlag)
	DeleteConfigGroupCmd.MarkFlagRequired(constants.NameFlag)
	DeleteConfigGroupCmd.MarkFlagRequired(constants.VersionFlag)
}
