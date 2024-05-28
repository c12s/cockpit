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
	diffStandaloneConfigShortDesc = "Compare two standalone configurations"
	diffStandaloneConfigLongDesc  = "This command compares two standalone configurations specified by their names and versions\n" +
		"displays the differences in a nicely formatted way, and saves them to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"cockpit diff standalone config --org 'org' --names 'db_config|db_config' --versions 'v1.0.0|v1.0.1'"

	// Path to files
	diffStandaloneFilePathJSON = "./response/standalone-config/standalone-config-diff.json"
	diffStandaloneFilePathYAML = "./response/standalone-config/standalone-config-diff.yaml"
)

var (
	singleConfigDiffResponse model.SingleConfigDiffResponse
)

var DiffStandaloneConfigCmd = &cobra.Command{
	Use:   "config",
	Short: diffStandaloneConfigShortDesc,
	Long:  diffStandaloneConfigLongDesc,
	Run:   executeDiffStandaloneConfig,
}

func executeDiffStandaloneConfig(cmd *cobra.Command, args []string) {
	requestBody, err := prepareStandaloneDiffRequest()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendStandaloneDiffRequest(requestBody); err != nil {
		fmt.Printf("Error comparing standalone configurations: %v\n", err)
		os.Exit(1)
	}

	render.RenderResponseToYAMLOrJSON(&singleConfigDiffResponse, outputFormat)

	filePath := diffStandaloneFilePathYAML
	if outputFormat == "json" {
		filePath = diffStandaloneFilePathJSON
	}

	if err := utils.SaveConfigResponseToFile(&singleConfigDiffResponse, filePath); err != nil {
		fmt.Printf("Failed to save response to file: %v\n", err)
		os.Exit(1)
	}
}

func prepareStandaloneDiffRequest() (interface{}, error) {
	namesList := strings.Split(names, "|")
	versionsList := strings.Split(versions, "|")

	if len(namesList) != 2 || len(versionsList) != 2 {
		return nil, fmt.Errorf("invalid names or versions format. Please use 'name1|name2' and 'version1|version2'")
	}

	requestBody := model.SingleConfigDiffRequest{
		Reference: model.SingleConfigReference{
			Name:         namesList[0],
			Organization: organization,
			Version:      versionsList[0],
		},
		Diff: model.SingleConfigReference{
			Name:         namesList[1],
			Organization: organization,
			Version:      versionsList[1],
		},
	}

	return requestBody, nil
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
		Response:    &singleConfigDiffResponse,
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}

	return nil
}

func init() {
	DiffStandaloneConfigCmd.Flags().StringVarP(&organization, flagOrg, shortFlagOrg, "", descOrg)
	DiffStandaloneConfigCmd.Flags().StringVarP(&names, flagNames, shortFlagNames, "", descNames)
	DiffStandaloneConfigCmd.Flags().StringVarP(&versions, flagVersions, shortFlagVersions, "", descVersions)
	DiffStandaloneConfigCmd.Flags().StringVarP(&outputFormat, flagOutput, shortFlagOutput, "yaml", descOutput)

	DiffStandaloneConfigCmd.MarkFlagRequired(flagOrg)
	DiffStandaloneConfigCmd.MarkFlagRequired(flagNames)
	DiffStandaloneConfigCmd.MarkFlagRequired(flagVersions)
}
