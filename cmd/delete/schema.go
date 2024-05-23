package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	deleteSchemaShortDesc = "Delete a schema"
	deleteSchemaLongDesc  = "This command deletes a schema from the specified organization.\n\n" +
		"Example:\n" +
		"delete-schema --org 'c12s' --schema_name 'schema' --version 'v1.0.1'"

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
	organization string
	schemaName   string
	version      string
)

var DeleteSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: deleteSchemaShortDesc,
	Long:  deleteSchemaLongDesc,
	Run:   executeDeleteSchema,
}

func executeDeleteSchema(cmd *cobra.Command, args []string) {
	config, err := createDeleteRequestConfig()
	if err != nil {
		fmt.Printf("Error creating request config: %v\n", err)
		os.Exit(1)
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		fmt.Printf("Error deleting schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Schema deleted successfully!")
}

func createDeleteRequestConfig() (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "DeleteConfigSchema")

	schemaDetails := model.SchemaDetails{
		Organization: organization,
		SchemaName:   schemaName,
		Version:      version,
	}

	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: schemaDetails,
	}

	return model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
	}, nil
}

func init() {
	DeleteSchemaCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	DeleteSchemaCmd.Flags().StringVarP(&schemaName, flagSchemaName, shortFlagSchemaName, "", descSchemaName)
	DeleteSchemaCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)

	DeleteSchemaCmd.MarkFlagRequired(flagOrganization)
	DeleteSchemaCmd.MarkFlagRequired(flagSchemaName)
	DeleteSchemaCmd.MarkFlagRequired(flagVersion)
}
