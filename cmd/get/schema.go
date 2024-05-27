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
	getSchemaShortDesc = "Retrieve and display the schema"
	getSchemaLongDesc  = "This command retrieves the schema from a specified endpoint\n" +
		"displays it in a nicely formatted way, and saves it to a YAML file.\n\n" +
		"Example:\n" +
		"get-schema --org 'org' --schema_name 'schema_name' --version 'v1.0.0'"

	// Flag Constants
	flagOrganization = "org"
	flagSchemaName   = "schema_name"
	flagVersion      = "version"

	// Flag Shorthand Constants
	shortFlagOrganization = "r"
	shortFlagSchemaName   = "s"
	shortFlagVersion      = "v"

	// Flag Descriptions
	descOrganization = "Organization name (required)"
	descSchemaName   = "Schema name (required)"
	descVersion      = "Schema version (required)"

	//Path to file
	saveSchemaToFile = "response/schema/schema.yaml"
)

var (
	organization   string
	schemaName     string
	version        string
	schemaResponse model.SchemaResponse
)

var GetSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: getSchemaShortDesc,
	Long:  getSchemaLongDesc,
	Run:   executeGetSchema,
}

func executeGetSchema(cmd *cobra.Command, args []string) {
	requestBody, err := prepareSchemaRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendSchemaRequest(requestBody); err != nil {
		fmt.Printf("Error retrieving schema: %v\n", err)
		os.Exit(1)
	}

	render.HandleSchemaResponse(&schemaResponse)

	if err := utils.SaveSchemaResponseToYAML(&schemaResponse, saveSchemaToFile); err != nil {
		fmt.Printf("Failed to save response to YAML file: %v\n", err)
		os.Exit(1)
	}
}

func prepareSchemaRequestConfig() (interface{}, error) {
	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: model.SchemaDetails{
			Organization: organization,
			SchemaName:   schemaName,
			Version:      version,
		},
	}

	return requestBody, nil
}

func sendSchemaRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetConfigSchema")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &schemaResponse,
		Timeout:     10 * time.Second,
	})
}

func init() {
	GetSchemaCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	GetSchemaCmd.Flags().StringVarP(&schemaName, flagSchemaName, shortFlagSchemaName, "", descSchemaName)
	GetSchemaCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)

	GetSchemaCmd.MarkFlagRequired(flagOrganization)
	GetSchemaCmd.MarkFlagRequired(flagSchemaName)
	GetSchemaCmd.MarkFlagRequired(flagVersion)
}
