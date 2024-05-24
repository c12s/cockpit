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
	getSchemaVersionShortDesc = "Retrieve and display schema versions"
	getSchemaVersionLongDesc  = "This command retrieves schema versions for specific schema, \n" +
		"displays them in a nicely formatted way, and saves them to a YAML file. \n\n" +
		"Example:\n" +
		"get-schema-version --org 'org' --schema_name 'schema_name'"
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

	err = saveVersionResponseToYAML(config.Response.(*model.SchemaVersionResponse))
	if err != nil {
		log.Fatalf("Failed to save response to YAML file: %v", err)
	}
}

func saveVersionResponseToYAML(response *model.SchemaVersionResponse) error {
	yamlData, err := utils.MarshalSchemaVersionResponse(response)
	if err != nil {
		return fmt.Errorf("failed to convert to YAML: %v", err)
	}

	fileName := fmt.Sprintf("./schema_files/schema_versions.yaml")
	err = ioutil.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write YAML file: %v", err)
	}

	fmt.Printf("Schema saved to %s\n", fileName)
	return nil
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
