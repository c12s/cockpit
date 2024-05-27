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
	requestBody, err := prepareStandaloneRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendStandaloneRequest(requestBody); err != nil {
		fmt.Printf("Error retrieving standalone configuration: %v\n", err)
		os.Exit(1)
	}

	render.RenderResponseToYAMLOrJSON(&standaloneConfigResponse, outputFormat)

	filePath := getStandaloneConfigFilePathYAML
	if outputFormat == "json" {
		filePath = getStandaloneConfigFilePathJSON
	}

	if err := utils.SaveConfigResponseToFile(&standaloneConfigResponse, filePath); err != nil {
		fmt.Printf("Failed to save response to file: %v\n", err)
		os.Exit(1)
	}
}

func prepareStandaloneRequestConfig() (interface{}, error) {
	requestBody := model.ConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}

	return requestBody, nil
}

func sendStandaloneRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetStandaloneConfig")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &standaloneConfigResponse,
		Timeout:     10 * time.Second,
	})
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
