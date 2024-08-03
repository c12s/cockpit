package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"

	"github.com/spf13/cobra"
)

var (
	response any
)

var GetNamespaceCmd = &cobra.Command{
	Use: "namespace",
	Run: executeGetNamespace,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.OrganizationFlag, constants.NameFlag})
	},
}

func executeGetNamespace(cmd *cobra.Command, args []string) {
	requestBody := prepareNamespaceRequestConfig()

	if err := sendNamespaceRequest(requestBody); err != nil {
		fmt.Println("Error sending get namespace request", err)
		os.Exit(1)
	}

	if err := utils.SaveYAMLOrJSONResponseToFile(response, constants.ResponseDirPathJSON+"namespace.yaml"); err != nil {
		fmt.Printf("Failed to save response to YAML file: %v\n", err)
		os.Exit(1)
	}
}

func prepareNamespaceRequestConfig() interface{} {
	return struct {
		OrgId string `json:"orgId"`
		Name  string `json:"name"`
	}{
		OrgId: organization,
		Name:  name,
	}
}

func sendNamespaceRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetNamespace")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &response,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetNamespaceCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	GetNamespaceCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)

	GetNamespaceCmd.MarkFlagRequired(constants.OrganizationFlag)
	GetNamespaceCmd.MarkFlagRequired(constants.NameFlag)
}
