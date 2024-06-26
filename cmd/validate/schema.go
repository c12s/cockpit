package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	validateSchemaVersionShortDesc = "Validate a schema version"
	validateSchemaVersionLongDesc  = `This command validates a schema version with the given configuration.
The user specifies the organization, schema name, version, and path to the YAML or JSON configuration file.
It reads the configuration file and validates the schema version against it.

Example:
- cockpit validate schema --org 'org' --schema-name 'schema' --version 'v1.0.0' --path '/path/to/config.yaml'`

	// Flag Constants
	organizationFlag = "org"
	schemaNameFlag   = "schema-name"
	versionFlag      = "version"
	configPathFlag   = "path"

	// Flag Shorthand Constants
	organizationShorthandFlag = "r"
	schemaNameShorthandFlag   = "s"
	versionShorthandFlag      = "v"
	configPathShorthandFlag   = "p"

	// Flag Descriptions
	organizationDescription = "Organization name (required)"
	schemaNameDescription   = "Schema name (required)"
	versionDescription      = "Schema version (required)"
	configPathDescription   = "Path to the YAML configuration file (required)"
)

var (
	organization string
	schemaName   string
	version      string
	configPath   string
)

var ValidateSchemaVersionCmd = &cobra.Command{
	Use:     "schema",
	Aliases: aliases.SchemaAliases,
	Short:   validateSchemaVersionShortDesc,
	Long:    validateSchemaVersionLongDesc,
	Run:     executeValidateSchemaVersion,
}

func executeValidateSchemaVersion(cmd *cobra.Command, args []string) {
	requestBody, err := prepareValidateSchemaRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendValidateSchemaRequest(requestBody); err != nil {
		fmt.Println("Error sending validate schema request:", err)
		os.Exit(1)
	}

	fmt.Println("Schema validated successfully!")
	println()
}

func prepareValidateSchemaRequestConfig() (interface{}, error) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	schemaDetails := model.SchemaDetails{
		Organization: organization,
		SchemaName:   schemaName,
		Version:      version,
	}

	requestBody := struct {
		SchemaDetails model.SchemaDetails `json:"schema_details"`
		Configuration string              `json:"configuration"`
	}{
		SchemaDetails: schemaDetails,
		Configuration: string(configData),
	}

	return requestBody, nil
}

func sendValidateSchemaRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "ValidateConfiguration")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
	})
}

func init() {
	ValidateSchemaVersionCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	ValidateSchemaVersionCmd.Flags().StringVarP(&schemaName, schemaNameFlag, schemaNameShorthandFlag, "", schemaNameDescription)
	ValidateSchemaVersionCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)
	ValidateSchemaVersionCmd.Flags().StringVarP(&configPath, configPathFlag, configPathShorthandFlag, "", configPathDescription)

	ValidateSchemaVersionCmd.MarkFlagRequired(organizationFlag)
	ValidateSchemaVersionCmd.MarkFlagRequired(schemaNameFlag)
	ValidateSchemaVersionCmd.MarkFlagRequired(versionFlag)
	ValidateSchemaVersionCmd.MarkFlagRequired(configPathFlag)
}
