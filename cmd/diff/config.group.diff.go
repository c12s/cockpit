package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	diffConfigGroupShortDesc = "Compare two configuration groups"
	diffConfigGroupLongDesc  = `This command compares two configuration groups specified by their names and versions, displays the differences in a nicely formatted way, and saves them to both YAML and JSON files.
The user can specify the organization, names, and versions of the two configuration groups to be compared. The differences between the groups will be highlighted and saved in the specified format.

Example:
- cockpit diff config group --org 'org' --names 'name1|name2' --versions 'version1|version2'
- cockpit diff config group --org 'org' --names 'name' --versions 'version1|version2'
- cockpit diff config group --org 'org' --names 'name1|name2' --versions 'version'`

	// Flag Constants
	organizationFlag = "org"
	namesFlag        = "names"
	versionsFlag     = "versions"
	outputFlag       = "output"

	// Flag Shorthand Constants
	orgShorthandFlag     = "r"
	namesShorthandFlag   = "n"
	versionShorthandFlag = "v"
	outputShorthandFlag  = "o"

	// Flag Descriptions
	orgDescription      = "Organization (required)"
	namesDescription    = "Configuration group names separated by '|' (required)"
	versionsDescription = "Configuration group versions separated by '|' (required)"
	outputDescription   = "Output format (yaml or json)"

	// Path to files
	diffConfigFilePathJSON = "./response/config-group/config-group-diff.json"
	diffConfigFilePathYAML = "./response/config-group/config-group-diff.yaml"
)

var (
	organization string
	names        string
	versions     string
	outputFormat string
	diffResponse model.ConfigGroupDiffResponse
)

var DiffConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"conf", "cnfg", "cfg"},
	Short:   diffConfigGroupShortDesc,
	Long:    diffConfigGroupLongDesc,
	Run:     executeDiffConfigGroup,
}

func executeDiffConfigGroup(cmd *cobra.Command, args []string) {
	requestBody, err := utils.PrepareConfigDiffRequest(names, versions, organization)
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

		filePath := diffConfigFilePathYAML
		if outputFormat == "json" {
			filePath = diffConfigFilePathJSON
		}

		if err := utils.SaveYAMLOrJSONResponseToFile(&diffResponse, filePath); err != nil {
			fmt.Println("Failed to save response to file:", err)
			println()
			os.Exit(1)
		}
	} else {
		println("Invalid output format. Expected 'yaml' or 'json'.")
	}

	println()
}

func prepareDiffRequest() (interface{}, error) {
	namesList := strings.Split(names, "|")
	versionsList := strings.Split(versions, "|")

	if len(namesList) != 2 || len(versionsList) != 2 {
		return nil, fmt.Errorf("invalid names or versions format. Please use 'name1|name2' and 'version1|version2'")
	}

	requestBody := model.ConfigGroupDiffRequest{
		Reference: model.ConfigReference{
			Name:         namesList[0],
			Organization: organization,
			Version:      versionsList[0],
		},
		Diff: model.ConfigReference{
			Name:         namesList[1],
			Organization: organization,
			Version:      versionsList[1],
		},
	}

	return requestBody, nil
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
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &diffResponse,
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}

	return nil
}

func init() {
	DiffConfigGroupCmd.Flags().StringVarP(&organization, organizationFlag, orgShorthandFlag, "", orgDescription)
	DiffConfigGroupCmd.Flags().StringVarP(&names, namesFlag, namesShorthandFlag, "", namesDescription)
	DiffConfigGroupCmd.Flags().StringVarP(&versions, versionsFlag, versionShorthandFlag, "", versionsDescription)
	DiffConfigGroupCmd.Flags().StringVarP(&outputFormat, outputFlag, outputShorthandFlag, "", outputDescription)

	DiffConfigGroupCmd.MarkFlagRequired(organizationFlag)
	DiffConfigGroupCmd.MarkFlagRequired(namesFlag)
	DiffConfigGroupCmd.MarkFlagRequired(versionsFlag)
}
