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
	schemaVersionResponse model.SchemaVersionResponse
)

var GetSchemaVersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: aliases.VersionAliases,
	Short:   constants.GetSchemaVersionShortDesc,
	Long:    constants.GetSchemaVersionLongDesc,
	Run:     executeGetSchemaVersion,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.SchemaNameFlag})
	},
}

func executeGetSchemaVersion(cmd *cobra.Command, args []string) {
	requestBody, err := prepareSchemaVersionRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendSchemaVersionRequest(requestBody); err != nil {
		fmt.Printf("Error retrieving schema versions: %v\n", err)
		os.Exit(1)
	}

	render.RenderSchemaVersionsTabWriter(schemaVersionResponse.SchemaVersions)
	if err := utils.SaveVersionResponseToYAML(&schemaVersionResponse, constants.GetSchemaVersionFilePathYAML); err != nil {
		fmt.Printf("Failed to save response to YAML file: %v\n", err)
		os.Exit(1)
	}
}

func prepareSchemaVersionRequestConfig() (interface{}, error) {
	requestBody := model.SchemaDetailsRequest{
		SchemaDetails: model.SchemaDetails{
			Organization: organization,
			Namespace:    namespace,
			SchemaName:   schemaName,
		},
	}

	return requestBody, nil
}

func sendSchemaVersionRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetConfigSchemaVersions")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &schemaVersionResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetSchemaVersionCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	GetSchemaVersionCmd.Flags().StringVarP(&schemaName, constants.SchemaNameFlag, constants.SchemaNameShorthandFlag, "", constants.SchemaNameDescription)
	GetSchemaVersionCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	GetSchemaVersionCmd.MarkFlagRequired(constants.NamespaceFlag)
	GetSchemaVersionCmd.MarkFlagRequired(constants.OrganizationFlag)
	GetSchemaVersionCmd.MarkFlagRequired(constants.SchemaNameFlag)
}
