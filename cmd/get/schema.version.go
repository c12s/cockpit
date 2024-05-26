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
	getSchemaVersionShortDesc = "Retrieve and display schema versions"
	getSchemaVersionLongDesc  = "This command retrieves schema versions for specific schema, \n" +
		"displays them in a nicely formatted way, and saves them to a YAML file. \n\n" +
		"Example:\n" +
		"get-schema-version --org 'org' --schema_name 'schema_name'"

	//Path to file
	saveSchemaVersionToFile = "response/schema/schema-version.yaml"
)

var (
	schemaVersionResponse model.SchemaVersionResponse
)

var GetSchemaVersionCmd = &cobra.Command{
	Use:   "version",
	Short: getSchemaVersionShortDesc,
	Long:  getSchemaVersionLongDesc,
	Run:   executeGetSchemaVersion,
}

func executeGetSchemaVersion(cmd *cobra.Command, args []string) {
	config := createSchemaVersionRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.HandleSchemaVersionResponse(config.Response.(*model.SchemaVersionResponse))

	err = utils.SaveVersionResponseToYAML(config.Response.(*model.SchemaVersionResponse), saveSchemaVersionToFile)
	if err != nil {
		log.Fatalf("Failed to save response to YAML file: %v", err)
	}
}

func createSchemaVersionRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "GetConfigSchemaVersions")

	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: model.SchemaDetails{
			Organization: organization,
			SchemaName:   schemaName,
		},
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &schemaVersionResponse,
	}
}

func init() {
	GetSchemaVersionCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	GetSchemaVersionCmd.Flags().StringVarP(&schemaName, flagSchemaName, shortFlagSchemaName, "", descSchemaName)

	GetSchemaVersionCmd.MarkFlagRequired(flagOrganization)
	GetSchemaVersionCmd.MarkFlagRequired(flagSchemaName)
}
