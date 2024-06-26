package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	deleteSchemaShortDesc = "Delete a schema"
	deleteSchemaLongDesc  = `This command deletes a schema version from the specified organization.
The user must provide the organization name, schema name, and version to delete the schema. This ensures that the specified schema version is removed from the system.

Example:
- cockpit delete schema --org 'c12s' --schema-name 'schema' --version 'v1.0.1'`

	// Flag Constants
	organizationFlag = "org"
	schemaNameFlag   = "schema-name"
	versionFlag      = "version"

	// Flag Shorthand Constants
	organizationShorthandFlag = "r"
	schemaNameShorthandFlag   = "s"
	versionShorthandFlag      = "v"

	// Flag Descriptions
	organizationDescription = "Organization name (required)"
	schemaNameDescription   = "Schema name (required)"
	versionDescription      = "Schema version (required)"
)

var (
	organization string
	schemaName   string
	version      string
)

var DeleteSchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: aliases.SchemaAliases,
	Short:   deleteSchemaShortDesc,
	Long:    deleteSchemaLongDesc,
	Run:     executeDeleteSchema,
}

func executeDeleteSchema(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteSchemaRequest()

	if err := sendDeleteRequestConfig(requestBody); err != nil {
		fmt.Println("Error sending delete schema request:", err)
		os.Exit(1)
	}

	fmt.Println("Schema deleted successfully!")
	println()
}

func prepareDeleteSchemaRequest() interface{} {
	schemaDetails := model.SchemaDetails{
		Organization: organization,
		SchemaName:   schemaName,
		Version:      version,
	}

	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: schemaDetails,
	}

	return requestBody
}

func sendDeleteRequestConfig(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "DeleteConfigSchema")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
	})
}

func init() {
	DeleteSchemaCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	DeleteSchemaCmd.Flags().StringVarP(&schemaName, schemaNameFlag, schemaNameShorthandFlag, "", schemaNameDescription)
	DeleteSchemaCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)

	DeleteSchemaCmd.MarkFlagRequired(organizationFlag)
	DeleteSchemaCmd.MarkFlagRequired(schemaNameFlag)
	DeleteSchemaCmd.MarkFlagRequired(versionFlag)
}
