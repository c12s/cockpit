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
)

var DeleteSchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: aliases.SchemaAliases,
	Short:   constants.DeleteSchemaShortDesc,
	Long:    constants.DeleteSchemaLongDesc,
	Run:     executeDeleteSchema,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.SchemaNameFlag, constants.OrganizationFlag, constants.VersionFlag})
	},
}

func executeDeleteSchema(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteSchemaRequest()

	if err := sendDeleteRequestConfig(requestBody); err != nil {
		fmt.Println("Error sending delete schema request:", err)
		os.Exit(1)
	}

	fmt.Println("Schema deleted successfully!")
}

func prepareDeleteSchemaRequest() interface{} {
	schemaDetails := model.SchemaDetails{
		Organization: organization,
		SchemaName:   schemaName,
		Version:      version,
		Namespace:    namespace,
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
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
	})
}

func init() {
	DeleteSchemaCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DeleteSchemaCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	DeleteSchemaCmd.Flags().StringVarP(&schemaName, constants.SchemaNameFlag, constants.SchemaNameShorthandFlag, "", constants.SchemaNameDescription)
	DeleteSchemaCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)

	DeleteSchemaCmd.MarkFlagRequired(constants.OrganizationFlag)
	DeleteSchemaCmd.MarkFlagRequired(constants.NamespaceFlag)
	DeleteSchemaCmd.MarkFlagRequired(constants.SchemaNameFlag)
	DeleteSchemaCmd.MarkFlagRequired(constants.VersionFlag)
}
