package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"log"
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
		"diff-standalone-config --org 'c12s' --names 'db_config|db_config' --versions 'v1.0.0|v1.0.1'"

	// Path to files
	diffStandaloneFilePathJSON = "./standalone_config_files/standalone-config-diff.json"
	diffStandaloneFilePathYAML = "./standalone_config_files/standalone-config-diff.yaml"
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
	config := createStandaloneDiffRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.HandleSingleConfigDiffResponse(config.Response.(*model.SingleConfigDiffResponse), outputFormat)

	filePath := diffStandaloneFilePathYAML
	if outputFormat == "json" {
		filePath = diffStandaloneFilePathJSON
	}

	err = utils.SaveResponseToFile(config.Response.(*model.SingleConfigDiffResponse), filePath)
	if err != nil {
		log.Fatalf("Failed to save response to file: %v", err)
	}
}

func createStandaloneDiffRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "DiffStandaloneConfig")

	namesList := strings.Split(names, "|")
	versionsList := strings.Split(versions, "|")

	if len(namesList) != 2 || len(versionsList) != 2 {
		log.Fatalf("Invalid names or versions format. Please use 'name1|name2' and 'version1|version2'")
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

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &singleConfigDiffResponse,
	}
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
