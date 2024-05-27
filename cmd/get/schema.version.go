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
	requestBody, err := prepareSchemaVersionRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendSchemaVersionRequest(requestBody); err != nil {
		fmt.Printf("Error retrieving schema versions: %v\n", err)
		os.Exit(1)
	}

	render.HandleSchemaVersionResponse(&schemaVersionResponse)

	if err := utils.SaveVersionResponseToYAML(&schemaVersionResponse, saveSchemaVersionToFile); err != nil {
		fmt.Printf("Failed to save response to YAML file: %v\n", err)
		os.Exit(1)
	}
}

func prepareSchemaVersionRequestConfig() (interface{}, error) {
	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: model.SchemaDetails{
			Organization: organization,
			SchemaName:   schemaName,
		},
	}

	return requestBody, nil
}

func sendSchemaVersionRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetConfigSchemaVersions")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &schemaVersionResponse,
		Timeout:     10 * time.Second,
	})
}

func init() {
	GetSchemaVersionCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	GetSchemaVersionCmd.Flags().StringVarP(&schemaName, flagSchemaName, shortFlagSchemaName, "", descSchemaName)

	GetSchemaVersionCmd.MarkFlagRequired(flagOrganization)
	GetSchemaVersionCmd.MarkFlagRequired(flagSchemaName)
}
