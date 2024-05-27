package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	getStandaloneConfigShortDesc = "Retrieve and display a standalone configuration"
	getStandaloneConfigLongDesc  = "This command retrieves a standalone configuration specified by its name, organization, and version\n" +
		"displays it in a nicely formatted way, and saves it to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"get-standalone-config --org 'c12s' --name 'db_config' --version 'v1.0.0'"

	// Path to files
	getStandaloneConfigFilePathJSON = "./response/standalone-config/standalone-config.json"
	getStandaloneConfigFilePathYAML = "./response/standalone-config/standalone-config.yaml"
)

var (
	standaloneConfigResponse model.SingleConfigGroupResponse
)

var GetStandaloneConfigCmd = &cobra.Command{
	Use:   "config",
	Short: getStandaloneConfigShortDesc,
	Long:  getStandaloneConfigLongDesc,
	Run:   executeGetStandaloneConfig,
}

func executeGetStandaloneConfig(cmd *cobra.Command, args []string) {
	config := createStandaloneRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.RenderResponseToYAMLOrJSON(config.Response.(*model.SingleConfigGroupResponse), outputFormat)

	filePath := getStandaloneConfigFilePathYAML
	if outputFormat == "json" {
		filePath = getStandaloneConfigFilePathJSON
	}

	err = utils.SaveConfigResponseToFile(config.Response.(*model.SingleConfigGroupResponse), filePath)
	if err != nil {
		log.Fatalf("Failed to save response to files: %v", err)
	}
}

func createStandaloneRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "GetStandaloneConfig")

	requestBody := model.ConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &standaloneConfigResponse,
	}
}

func init() {
	GetStandaloneConfigCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	GetStandaloneConfigCmd.Flags().StringVarP(&name, flagName, shortFlagName, "", descName)
	GetStandaloneConfigCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)
	GetStandaloneConfigCmd.Flags().StringVarP(&outputFormat, flagOutput, shortFlagOutput, "yaml", descOutput)

	GetStandaloneConfigCmd.MarkFlagRequired(flagOrganization)
	GetStandaloneConfigCmd.MarkFlagRequired(flagName)
	GetStandaloneConfigCmd.MarkFlagRequired(flagVersion)
}
