package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"

	"github.com/spf13/cobra"
)

var (
	organization   string
	schemaName     string
	version        string
	namespace      string
	schemaResponse model.SchemaResponse
)

var GetSchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: aliases.SchemaAliases,
	Short:   constants.GetSchemaShortDesc,
	Long:    constants.GetSchemaLongDesc,
	Run:     executeGetSchema,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.SchemaNameFlag, constants.VersionFlag})
	},
}

func executeGetSchema(cmd *cobra.Command, args []string) {
	requestBody := prepareSchemaRequestConfig()

	if err := sendSchemaRequest(requestBody); err != nil {
		fmt.Println("Error sending get schema request", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(schemaResponse.SchemaData, "yaml", "")
	if err := utils.SaveSchemaResponseToYAML(&schemaResponse, constants.GetSchemaFilePathYAML); err != nil {
		fmt.Printf("Failed to save response to YAML file: %v\n", err)
		os.Exit(1)
	}
}

func prepareSchemaRequestConfig() interface{} {
	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: model.SchemaDetails{
			Organization: organization,
			SchemaName:   schemaName,
			Version:      version,
			Namespace:    namespace,
		},
	}

	return requestBody
}

func sendSchemaRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetConfigSchema")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &schemaResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetSchemaCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	GetSchemaCmd.Flags().StringVarP(&schemaName, constants.SchemaNameFlag, constants.SchemaNameShorthandFlag, "", constants.SchemaNameDescription)
	GetSchemaCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	GetSchemaCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	GetSchemaCmd.MarkFlagRequired(constants.OrganizationFlag)
	GetSchemaCmd.MarkFlagRequired(constants.SchemaNameFlag)
	GetSchemaCmd.MarkFlagRequired(constants.VersionFlag)
	GetSchemaCmd.MarkFlagRequired(constants.NamespaceFlag)
}
