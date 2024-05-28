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
	getAppConfigShortDesc = "Retrieve and display the configuration"
	getAppConfigLongDesc  = "This command retrieves the specific configuration\n" +
		"displays it in a nicely formatted way, and saves it to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"cockpit get config group --org 'org' --name 'app_config' --version 'v1.0.0'"

	// Flag Constants
	flagName   = "name"
	flagOutput = "output"

	// Flag Shorthand Constants
	shortFlagName   = "n"
	shortFlagOutput = "o"

	// Flag Descriptions
	descName   = "Configuration name (required)"
	descOutput = "Output format (yaml or json)"

	// Path to files
	getConfigFilePathJSON = "./response/config-group/single-config.json"
	getConfigFilePathYAML = "./response/config-group/single-config.yaml"
)

var (
	name              string
	appConfigResponse model.SingleConfigGroupResponse
	outputFormat      string
)
var GetSingleConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: getAppConfigShortDesc,
	Long:  getAppConfigLongDesc,
	Run:   executeGetAppConfig,
}

func executeGetAppConfig(cmd *cobra.Command, args []string) {
	requestBody, err := prepareRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendRequest(requestBody); err != nil {
		fmt.Printf("Error retrieving configuration: %v\n", err)
		os.Exit(1)
	}

	render.RenderResponseToYAMLOrJSON(&appConfigResponse, outputFormat)

	filePath := getConfigFilePathYAML
	if outputFormat == "json" {
		filePath = getConfigFilePathJSON
	}

	if err := utils.SaveConfigResponseToFile(&appConfigResponse, filePath); err != nil {
		fmt.Printf("Failed to save response to file: %v\n", err)
		os.Exit(1)
	}
}

func prepareRequestConfig() (interface{}, error) {
	requestBody := model.ConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}

	return requestBody, nil
}

func sendRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetConfigGroup")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &appConfigResponse,
		Timeout:     10 * time.Second,
	})
}

func init() {
	GetSingleConfigGroupCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	GetSingleConfigGroupCmd.Flags().StringVarP(&name, flagName, shortFlagName, "", descName)
	GetSingleConfigGroupCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)
	GetSingleConfigGroupCmd.Flags().StringVarP(&outputFormat, flagOutput, shortFlagOutput, "yaml", descOutput)

	GetSingleConfigGroupCmd.MarkFlagRequired(flagOrganization)
	GetSingleConfigGroupCmd.MarkFlagRequired(flagName)
	GetSingleConfigGroupCmd.MarkFlagRequired(flagVersion)
}
