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
	diffConfigGroupLongDesc  = "This command compares two configuration groups specified by their names and versions\n" +
		"displays the differences in a nicely formatted way, and saves them to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"diff config group --org 'org' --names 'name1|name2' --versions 'version1|version2'"

	// Flag Constants
	flagOrg      = "org"
	flagNames    = "names"
	flagVersions = "versions"
	flagOutput   = "output"

	// Flag Shorthand Constants
	shortFlagOrg      = "r"
	shortFlagNames    = "n"
	shortFlagVersions = "v"
	shortFlagOutput   = "o"

	// Flag Descriptions
	descOrg      = "Organization (required)"
	descNames    = "Configuration group names separated by '|' (required)"
	descVersions = "Configuration group versions separated by '|' (required)"
	descOutput   = "Output format (yaml or json)"

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
	Use:   "group",
	Short: diffConfigGroupShortDesc,
	Long:  diffConfigGroupLongDesc,
	Run:   executeDiffConfigGroup,
}

func executeDiffConfigGroup(cmd *cobra.Command, args []string) {
	requestBody, err := prepareDiffRequest()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendDiffRequest(requestBody); err != nil {
		fmt.Printf("Error comparing configuration groups: %v\n", err)
		os.Exit(1)
	}

	render.RenderResponseToYAMLOrJSON(&diffResponse, outputFormat)

	filePath := diffConfigFilePathYAML
	if outputFormat == "json" {
		filePath = diffConfigFilePathJSON
	}

	if err := utils.SaveConfigResponseToFile(&diffResponse, filePath); err != nil {
		fmt.Printf("Failed to save response to file: %v\n", err)
		os.Exit(1)
	}
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
	DiffConfigGroupCmd.Flags().StringVarP(&organization, flagOrg, shortFlagOrg, "", descOrg)
	DiffConfigGroupCmd.Flags().StringVarP(&names, flagNames, shortFlagNames, "", descNames)
	DiffConfigGroupCmd.Flags().StringVarP(&versions, flagVersions, shortFlagVersions, "", descVersions)
	DiffConfigGroupCmd.Flags().StringVarP(&outputFormat, flagOutput, shortFlagOutput, "yaml", descOutput)

	DiffConfigGroupCmd.MarkFlagRequired(flagOrg)
	DiffConfigGroupCmd.MarkFlagRequired(flagNames)
	DiffConfigGroupCmd.MarkFlagRequired(flagVersions)
}
