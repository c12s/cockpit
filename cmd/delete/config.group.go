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
	deleteConfigGroupShortDesc = "Delete a configuration group version"
	deleteConfigGroupLongDesc  = `This command deletes a specified configuration group version and displays the deleted configuration group details in JSON or YAML format.
The user can specify the organization, the configuration group name, and the version to be deleted. The output can be formatted as either JSON or YAML based on user preference.

Example:
- cockpit delete config group --org 'org' --name 'app_config' --version 'v1.0.0'`

	// Flag Constants
	nameFlag   = "name"
	outputFlag = "output"

	// Flag Shorthand Constants
	nameShorthandFlag   = "n"
	outputShorthandFlag = "o"

	// Flag Descriptions
	nameDescription   = "Configuration group name (required)"
	outputDescription = "Output format (json or yaml)"
)

var (
	name                      string
	outputFormat              string
	deleteConfigGroupResponse model.ConfigGroup
)

var DeleteConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   deleteConfigGroupShortDesc,
	Long:    deleteConfigGroupLongDesc,
	Run:     executeDeleteConfigGroup,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{organizationFlag, nameFlag, versionFlag})
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
	fmt.Println()
}

func prepareDeleteConfigGroupRequest() interface{} {
	requestBody := model.ConfigReference{
		Organization: organization,
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
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &deleteConfigGroupResponse,
	}
}

func init() {
	DeleteConfigGroupCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	DeleteConfigGroupCmd.Flags().StringVarP(&name, nameFlag, nameShorthandFlag, "", nameDescription)
	DeleteConfigGroupCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)
	DeleteConfigGroupCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	DeleteConfigGroupCmd.MarkFlagRequired(organizationFlag)
	DeleteConfigGroupCmd.MarkFlagRequired(nameFlag)
	DeleteConfigGroupCmd.MarkFlagRequired(versionFlag)
}
