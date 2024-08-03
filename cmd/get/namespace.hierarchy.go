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

var GetNamespaceHierarchyCmd = &cobra.Command{
	Use: "hierarchy",
	Run: executeGetNamespaceHierarchy,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.OrganizationFlag})
	},
}

func executeGetNamespaceHierarchy(cmd *cobra.Command, args []string) {
	requestBody := prepareNamespaceHierarchyRequestConfig()

	if err := sendNamespaceHierarchyRequest(requestBody); err != nil {
		fmt.Println("Error sending get namespace hierarchy request", err)
		os.Exit(1)
	}

	if err := utils.SaveYAMLOrJSONResponseToFile(response, constants.ResponseDirPathJSON+"namespace_hierarchy.yaml"); err != nil {
		fmt.Printf("Failed to save response to JSON file: %v\n", err)
		os.Exit(1)
	}
}

func prepareNamespaceHierarchyRequestConfig() interface{} {
	return struct {
		OrgId string `json:"orgId"`
	}{
		OrgId: organization,
	}
}

func sendNamespaceHierarchyRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetNamespaceHierarchy")

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
	GetNamespaceHierarchyCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)

	GetNamespaceHierarchyCmd.MarkFlagRequired(constants.OrganizationFlag)
}
