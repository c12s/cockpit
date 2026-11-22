package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	namespace    string
	schemaName   string
	version      string
	configPath   string
)

var ValidateSchemaVersionCmd = &cobra.Command{
	Use:     "schema",
	Aliases: aliases.SchemaAliases,
	Short:   constants.ValidateSchemaVersionShortDesc,
	Long:    constants.ValidateSchemaVersionLongDesc,
	Run:     executeValidateSchemaVersion,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.SchemaNameFlag, constants.VersionFlag, constants.FilePathFlag})
	},
}

func executeValidateSchemaVersion(cmd *cobra.Command, args []string) {
	requestBody, err := prepareValidateSchemaRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	resp, err := sendValidateSchemaRequest(requestBody)
	if err != nil {
		fmt.Println("Error sending validate schema request:", err)
		os.Exit(1)
	}
	if !resp.IsValid {
		fmt.Println("Configuration is not valid:")
		for _, m := range strings.Split(resp.Message, "; ") {
			fmt.Println("  -", m)
		}
		os.Exit(1)
	}
	fmt.Println("Configuration is valid against the schema.")
}

// validateSchemaResponse is Quasar's reply (proto JSON: false fields omitted).
type validateSchemaResponse struct {
	Message string `json:"message"`
	IsValid bool   `json:"isValid"`
}

func prepareValidateSchemaRequestConfig() (interface{}, error) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	schemaDetails := model.SchemaDetails{
		Organization: organization,
		Namespace:    namespace,
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

func sendValidateSchemaRequest(requestBody interface{}) (validateSchemaResponse, error) {
	var resp validateSchemaResponse
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return resp, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "ValidateConfiguration")

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &resp,
	})
	return resp, err
}

func init() {
	ValidateSchemaVersionCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ValidateSchemaVersionCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	ValidateSchemaVersionCmd.Flags().StringVarP(&schemaName, constants.SchemaNameFlag, constants.SchemaNameShorthandFlag, "", constants.SchemaNameDescription)
	ValidateSchemaVersionCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	ValidateSchemaVersionCmd.Flags().StringVarP(&configPath, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)

	ValidateSchemaVersionCmd.MarkFlagRequired(constants.OrganizationFlag)
	ValidateSchemaVersionCmd.MarkFlagRequired(constants.NamespaceFlag)
	ValidateSchemaVersionCmd.MarkFlagRequired(constants.SchemaNameFlag)
	ValidateSchemaVersionCmd.MarkFlagRequired(constants.VersionFlag)
	ValidateSchemaVersionCmd.MarkFlagRequired(constants.FilePathFlag)
}
