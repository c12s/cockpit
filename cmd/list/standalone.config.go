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
	listStandaloneConfigShortDesc = "List all standalone configurations"
	listStandaloneConfigLongDesc  = "This command retrieves a list of all standalone configurations for a given organization\n" +
		"displays them in a nicely formatted way, and saves them to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"cockpit list standalone config --org 'org'"

	listStandaloneConfigFilePathJSON = "./response/standalone-config/list-standalone-config.json"
	listStandaloneConfigFilePathYAML = "./response/standalone-config/list-standalone-config.yaml"
)

var (
	listResponse model.StandaloneConfigsResponse
)

var ListStandaloneConfigCmd = &cobra.Command{
	Use:   "config",
	Short: listStandaloneConfigShortDesc,
	Long:  listStandaloneConfigLongDesc,
	Run:   executeListStandaloneConfig,
}

func executeListStandaloneConfig(cmd *cobra.Command, args []string) {
	config := createListStandaloneRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.RenderResponseToYAMLOrJSON(config.Response.(*model.StandaloneConfigsResponse), outputFormat)

	filePath := listStandaloneConfigFilePathYAML
	if outputFormat == "json" {
		filePath = listStandaloneConfigFilePathJSON
	}

	err = utils.SaveConfigResponseToFile(config.Response.(*model.StandaloneConfigsResponse), filePath)
	if err != nil {
		log.Fatalf("Failed to save response to file: %v", err)
	}
}

func createListStandaloneRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListStandaloneConfig")

	requestBody := map[string]string{
		"organization": organization,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &listResponse,
	}
}

func init() {
	ListStandaloneConfigCmd.Flags().StringVarP(&organization, orgFlag, orgFlagShortHand, "", orgDesc)
	ListStandaloneConfigCmd.Flags().StringVarP(&outputFormat, outputFlag, outputFlagShortHand, "yaml", outputDesc)

	ListStandaloneConfigCmd.MarkFlagRequired("org")
}
