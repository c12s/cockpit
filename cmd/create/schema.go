package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"time"
)

var (
	organization string
	schemaName   string
	version      string
	filePath     string
)

const (
	createSchemaShortDescription = "Create a schema for an organization"
	createSchemaLongDescription  = "Creates a schema for an organization by providing schema details and the path to a YAML file containing the schema definition.\n\n" +
		"Example:\n" +
		"create-schema --organization 'c12s' --schema_name 'schema' --version 'v1.0.0' --path '/path/to/schema.yaml'"

	organizationFlag = "org"
	schemaNameFlag   = "schema_name"
	versionFlag      = "version"
	filePathFlag     = "path"

	organizationFlagShorthand = "r"
	schemaNameFlagShorthand   = "s"
	versionFlagShorthand      = "v"
	filePathFlagShorthand     = "p"

	organizationDescription = "Organization name (required)"
	schemaNameDescription   = "Schema name (required)"
	versionDescription      = "Schema version (required)"
	filePathDescription     = "Path to the YAML file containing the schema definition (required)"
)

var CreateSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: createSchemaShortDescription,
	Long:  createSchemaLongDescription,
	Run:   executeCreateSchema,
}

func executeCreateSchema(cmd *cobra.Command, args []string) {
	schema, err := readSchemaFile(filePath)
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	requestBody := createRequestBody(schema)
	config, err := createRequestConfig(requestBody)
	if err != nil {
		fmt.Printf("Error creating request config: %v\n", err)
		os.Exit(1)
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		fmt.Printf("Error creating schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Schema created successfully!")
}

func createRequestBody(schema string) map[string]interface{} {
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

func createRequestConfig(requestBody map[string]interface{}) (model.HTTPRequestConfig, error) {
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

func readSchemaFile(filePath string) (string, error) {
	schema, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema file: %v", err)
	}
	return string(schema), nil
}

func init() {
	CreateSchemaCmd.Flags().StringVarP(&organization, organizationFlag, organizationFlagShorthand, "", organizationDescription)
	CreateSchemaCmd.Flags().StringVarP(&schemaName, schemaNameFlag, schemaNameFlagShorthand, "", schemaNameDescription)
	CreateSchemaCmd.Flags().StringVarP(&version, versionFlag, versionFlagShorthand, "", versionDescription)
	CreateSchemaCmd.Flags().StringVarP(&filePath, filePathFlag, filePathFlagShorthand, "", filePathDescription)

	CreateSchemaCmd.MarkFlagRequired(organizationFlag)
	CreateSchemaCmd.MarkFlagRequired(schemaNameFlag)
	CreateSchemaCmd.MarkFlagRequired(versionFlag)
	CreateSchemaCmd.MarkFlagRequired(filePathFlag)
}
