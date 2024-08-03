package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	organization string
	schemaName   string
	version      string
	namespace    string
	filePath     string
)

var CreateSchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: aliases.SchemaAliases,
	Short:   constants.CreateSchemaShortDesc,
	Long:    constants.CreateSchemaLongDesc,
	Run:     executeCreateSchema,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.SchemaNameFlag, constants.VersionFlag, constants.FilePathFlag})
	},
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
}

func createSchemaRequest(schema string) map[string]interface{} {
	schemaDetails := model.SchemaDetails{
		Organization: organization,
		SchemaName:   schemaName,
		Version:      version,
		Namespace:    namespace,
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
		Timeout:     30 * time.Second,
	}, nil
}

func init() {
	CreateSchemaCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	CreateSchemaCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	CreateSchemaCmd.Flags().StringVarP(&schemaName, constants.SchemaNameFlag, constants.SchemaNameShorthandFlag, "", constants.SchemaNameDescription)
	CreateSchemaCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	CreateSchemaCmd.Flags().StringVarP(&filePath, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)

	CreateSchemaCmd.MarkFlagRequired(constants.SchemaNameFlag)
	CreateSchemaCmd.MarkFlagRequired(constants.NamespaceFlag)
	CreateSchemaCmd.MarkFlagRequired(constants.VersionFlag)
	CreateSchemaCmd.MarkFlagRequired(constants.FilePathFlag)
}
