package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	diffStandaloneConfigShortDesc = "Compare two standalone configurations"
	diffStandaloneConfigLongDesc  = `This command compares two standalone configurations specified by their names and versions, displays the differences in a nicely formatted way, and saves them to both YAML and JSON files.
The user can specify the organization, names, and versions of the two standalone configurations to be compared. The differences between the configurations will be highlighted and saved in the specified format.

Example:
- cockpit diff standalone config --org 'org' --names 'db_config|db_config' --versions 'v1.0.0|v1.0.1'
- cockpit diff standalone config --org 'org' --names 'db_config' --versions 'v1.0.0|v1.0.1'
- cockpit diff standalone config --org 'org' --names 'db_config|db_config' --versions 'v1.0.0`

	// Path to files
	diffStandaloneFilePathJSON = "./response/standalone-config/standalone-config-diff.json"
	diffStandaloneFilePathYAML = "./response/standalone-config/standalone-config-diff.yaml"
)

var (
	standaloneConfigDiffResponse model.StandaloneConfigDiffResponse
)

var DiffStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"grp", "gr"},
	Short:   diffStandaloneConfigShortDesc,
	Long:    diffStandaloneConfigLongDesc,
	Run:     executeDiffStandaloneConfig,
}

func executeDiffStandaloneConfig(cmd *cobra.Command, args []string) {
	requestBody, err := utils.PrepareConfigDiffRequest(names, versions, organization)
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

		filePath := diffStandaloneFilePathYAML
		if outputFormat == "json" {
			filePath = diffStandaloneFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&standaloneConfigDiffResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}
	fmt.Println()
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
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &standaloneConfigDiffResponse,
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}

	return nil
}

func init() {
	DiffStandaloneConfigCmd.Flags().StringVarP(&organization, organizationFlag, orgShorthandFlag, "", orgDescription)
	DiffStandaloneConfigCmd.Flags().StringVarP(&names, namesFlag, namesShorthandFlag, "", namesDescription)
	DiffStandaloneConfigCmd.Flags().StringVarP(&versions, versionsFlag, versionShorthandFlag, "", versionsDescription)
	DiffStandaloneConfigCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	DiffStandaloneConfigCmd.MarkFlagRequired(organizationFlag)
	DiffStandaloneConfigCmd.MarkFlagRequired(namesFlag)
	DiffStandaloneConfigCmd.MarkFlagRequired(versionsFlag)
}
