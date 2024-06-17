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
	getSchemaLongDesc  = `This command retrieves the schema from a specified organization and specific version, displays it in a nicely formatted way, and saves it to a YAML file.
The user can specify the organization, schema name, and version to retrieve the schema details. The response will be displayed in a tabular format and saved as a YAML file.

Example:
- cockpit get schema --org 'org' --schema_name 'schema_name' --version 'v1.0.0'`

	// Flag Constants
	organizationFlag = "org"
	schemaNameFlag   = "schema_name"
	versionFlag      = "version"

	// Flag Shorthand Constants
	organizationShorthandFlag = "r"
	schemaNameShorthandFlag   = "s"
	versionShorthandFlag      = "v"

	// Flag Descriptions
	organizationDescription = "Organization name (required)"
	schemaNameDescription   = "Schema name (required)"
	versionDescription      = "Schema version (required)"

	//Path to file
	getSchemaFilePath = "response/schema/schema.yaml"
)

var (
	organization   string
	schemaName     string
	version        string
	schemaResponse model.SchemaResponse
)

var GetSchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: []string{"schem", "schemaa", "sch", "sche"},
	Short:   getSchemaShortDesc,
	Long:    getSchemaLongDesc,
	Run:     executeGetSchema,
}

func executeGetSchema(cmd *cobra.Command, args []string) {
	requestBody := prepareSchemaRequestConfig()

	if err := sendSchemaRequest(requestBody); err != nil {
		fmt.Println("Error sending get schema request", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(schemaResponse.SchemaData)
	println()
	if err := utils.SaveSchemaResponseToYAML(&schemaResponse, getSchemaFilePath); err != nil {
		fmt.Printf("Failed to save response to YAML file: %v\n", err)
		os.Exit(1)
	}
	println()
}

func prepareSchemaRequestConfig() interface{} {
	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: model.SchemaDetails{
			Organization: organization,
			SchemaName:   schemaName,
			Version:      version,
		},
	}

	return requestBody
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
	GetSchemaCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	GetSchemaCmd.Flags().StringVarP(&schemaName, schemaNameFlag, schemaNameShorthandFlag, "", schemaNameDescription)
	GetSchemaCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)

	GetSchemaCmd.MarkFlagRequired(organizationFlag)
	GetSchemaCmd.MarkFlagRequired(schemaNameFlag)
	GetSchemaCmd.MarkFlagRequired(versionFlag)
}
