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
	deleteStandaloneConfigShortDesc = "Delete a standalone configuration version"
	deleteStandaloneConfigLongDesc  = "This command deletes a specified standalone configuration version\n" +
		"and displays the deleted configuration details in JSON format.\n\n" +
		"Example:\n" +
		"delete-standalone-config --org 'c12s' --name 'db_config' --version 'v1.0.1'"
)

var (
	deleteStandaloneConfigResponse model.StandaloneConfig
)

var DeleteStandaloneConfigCmd = &cobra.Command{
	Use:   "config",
	Short: deleteStandaloneConfigShortDesc,
	Long:  deleteStandaloneConfigLongDesc,
	Run:   executeDeleteStandaloneConfig,
}

func executeDeleteStandaloneConfig(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteStandaloneConfigRequestConfig()

	config := sendDeleteStandaloneConfigRequestConfig(requestBody)

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.DisplayResponse(&deleteStandaloneConfigResponse, output, "Standalone deleted successfully")
}

func prepareDeleteStandaloneConfigRequestConfig() model.SingleConfigReference {
	requestBody := model.SingleConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}
	return requestBody
}

func sendDeleteStandaloneConfigRequestConfig(requestBody interface{}) model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "DeleteStandaloneConfig")

	return model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &deleteStandaloneConfigResponse,
	}
}

func init() {
	DeleteStandaloneConfigCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&name, flagName, shortFlagName, "", descName)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)
	DeleteStandaloneConfigCmd.Flags().StringVarP(&output, flagOutput, shortFlagOutput, "yaml", descOutput)

	DeleteStandaloneConfigCmd.MarkFlagRequired(flagOrganization)
	DeleteStandaloneConfigCmd.MarkFlagRequired(flagName)
	DeleteStandaloneConfigCmd.MarkFlagRequired(flagVersion)
}
