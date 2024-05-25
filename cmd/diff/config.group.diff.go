package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"io/ioutil"
	"log"
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
		"diff-config-group --org 'org' --names 'name1|name2' --versions 'version1|version2'"

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
	diffConfigFilePathJSON = "./config_group_files/config-group-diff.json"
	diffConfigFilePathYAML = "./config_group_files/config-group-diff.yaml"
)

var (
	organization string
	names        string
	versions     string
	diffResponse model.ConfigGroupDiffResponse
	outputFormat string
)

var DiffConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: diffConfigGroupShortDesc,
	Long:  diffConfigGroupLongDesc,
	Run:   executeDiffConfigGroup,
}

func executeDiffConfigGroup(cmd *cobra.Command, args []string) {
	config := createDiffRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.HandleConfigGroupDiffResponse(config.Response.(*model.ConfigGroupDiffResponse), outputFormat)

	err = saveDiffConfigGroupResponseToFiles(config.Response.(*model.ConfigGroupDiffResponse))
	if err != nil {
		log.Fatalf("Failed to save response to files: %v", err)
	}
}

func createDiffRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "DiffConfigGroup")

	namesList := strings.Split(names, "|")
	versionsList := strings.Split(versions, "|")

	if len(namesList) != 2 || len(versionsList) != 2 {
		log.Fatalf("Invalid names or versions format. Please use 'name1|name2' and 'version1|version2'")
	}

	requestBody := model.ConfigGroupDiffRequest{
		Reference: model.ConfigGroupReference{
			Name:         namesList[0],
			Organization: organization,
			Version:      versionsList[0],
		},
		Diff: model.ConfigGroupReference{
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
		Response:    &diffResponse,
	}
}

func saveDiffConfigGroupResponseToFiles(response *model.ConfigGroupDiffResponse) error {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to convert to JSON: %v", err)
		}
		err = ioutil.WriteFile(diffConfigFilePathJSON, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write JSON file: %v", err)
		}
		fmt.Printf("Config group diff saved to %s\n", diffConfigFilePathJSON)
	} else {
		yamlData, err := utils.MarshalConfigGroupDiffResponseToYAML(response)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}
		err = ioutil.WriteFile(diffConfigFilePathYAML, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write YAML file: %v", err)
		}
		fmt.Printf("Config group diff saved to %s\n", diffConfigFilePathYAML)
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
