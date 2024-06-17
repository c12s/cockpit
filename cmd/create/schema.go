package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	createSchemaShortDescription = "Create a schema for an organization"
	createSchemaLongDescription  = `Creates a schema for an organization by providing schema details and the path to a YAML or JSON file containing the schema definition.
Schemas define the structure of configuration data that can be used across various services and applications within the organization. This command uploads and saves the schema to the server.

Example:
cockpit create schema --org 'org' --schema_name 'schema' --version 'v1.0.0' --path 'path to yaml or json file'`

	// Flag Constants
	organizationFlag = "org"
	schemaNameFlag   = "schema_name"
	versionFlag      = "version"
	filePathFlag     = "path"

	// Flag Shorthand Constants
	organizationShorthandFlag = "r"
	schemaNameShorthandFlag   = "s"
	versionShorthandFlag      = "v"
	filePathShorthandFlag     = "p"

	// Flag Descriptions
	organizationDescription = "Organization name (required)"
	schemaNameDescription   = "Schema name (required)"
	versionDescription      = "Schema version (required)"
	filePathDescription     = "Path to the YAML file containing the schema definition (required)"
)

var (
	organization string
	schemaName   string
	version      string
	filePath     string
)

var CreateSchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: []string{"schem", "schemaa", "sche"},
	Short:   createSchemaShortDescription,
	Long:    createSchemaLongDescription,
	Run:     executeCreateSchema,
}

func executeCreateSchema(cmd *cobra.Command, args []string) {
	schema, err := utils.ReadSchemaFile(filePath)
	if err != nil {
		fmt.Println("Error reading schema file:", err)
		os.Exit(1)
	}

	requestBody := createSchemaRequest(schema)
	config, err := prepareSchemaRequest(requestBody)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		fmt.Println("Error sending create schema request:", err)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("Schema created successfully!")
	fmt.Println()
}

func createSchemaRequest(schema string) map[string]interface{} {
	schemaDetails := model.SchemaDetails{
		Organization: organization,
		SchemaName:   schemaName,
		Version:      version,
	}

	return map[string]interface{}{
		"schema_details": schemaDetails,
		"schema":         schema,
	}
}

func prepareSchemaRequest(requestBody map[string]interface{}) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "SaveConfigSchema")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: requestBody,
		Token:       token,
		Timeout:     10 * time.Second,
	}, nil
}

func init() {
	CreateSchemaCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	CreateSchemaCmd.Flags().StringVarP(&schemaName, schemaNameFlag, schemaNameShorthandFlag, "", schemaNameDescription)
	CreateSchemaCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)
	CreateSchemaCmd.Flags().StringVarP(&filePath, filePathFlag, filePathShorthandFlag, "", filePathDescription)

	CreateSchemaCmd.MarkFlagRequired(organizationFlag)
	CreateSchemaCmd.MarkFlagRequired(schemaNameFlag)
	CreateSchemaCmd.MarkFlagRequired(versionFlag)
	CreateSchemaCmd.MarkFlagRequired(filePathFlag)
}
