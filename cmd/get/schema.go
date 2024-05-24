package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"io/ioutil"
	"log"
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
	shortFlagOrganization = "o"
	shortFlagSchemaName   = "s"
	shortFlagVersion      = "v"

	// Flag Descriptions
	descOrganization = "Organization name (required)"
	descSchemaName   = "Schema name (required)"
	descVersion      = "Schema version (required)"
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
	config := createSchemaRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.HandleSchemaResponse(config.Response.(*model.SchemaResponse))

	err = saveSchemaResponseToYAML(config.Response.(*model.SchemaResponse))
	if err != nil {
		log.Fatalf("Failed to save response to YAML file: %v", err)
	}
}

func createSchemaRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "GetConfigSchema")

	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: model.SchemaDetails{
			Organization: organization,
			SchemaName:   schemaName,
			Version:      version,
		},
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &schemaResponse,
	}
}

func saveSchemaResponseToYAML(response *model.SchemaResponse) error {
	if response.SchemaData.Schema != "" {
		yamlData, err := utils.MarshalSchemaResponse(response)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}

		fileName := fmt.Sprintf("./schemas/schema.yaml")
		err = ioutil.WriteFile(fileName, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write YAML file: %v", err)
		}

		fmt.Printf("Schema saved to %s\n", fileName)
	}
	return nil
}

func init() {
	GetSchemaCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	GetSchemaCmd.Flags().StringVarP(&schemaName, flagSchemaName, shortFlagSchemaName, "", descSchemaName)
	GetSchemaCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)

	GetSchemaCmd.MarkFlagRequired(flagOrganization)
	GetSchemaCmd.MarkFlagRequired(flagSchemaName)
	GetSchemaCmd.MarkFlagRequired(flagVersion)
}
