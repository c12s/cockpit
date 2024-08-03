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

var DeleteNamespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: constants.DeleteNamespaceShortDesc,
	Long:  constants.DeleteNamespaceLongDesc,
	Run:   executeDeleteNamespace,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NameFlag, constants.OrganizationFlag})
	},
}

func executeDeleteNamespace(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteNamespaceRequest()

	if err := sendDeleteNamespaceRequestConfig(requestBody); err != nil {
		fmt.Println("Error sending delete namespace request:", err)
		os.Exit(1)
	}

	fmt.Println("Namespace deleted successfully!")
}

func prepareDeleteNamespaceRequest() interface{} {
	return struct {
		OrgId string `json:"orgId"`
		Name  string `json:"name"`
	}{
		OrgId: organization,
		Name:  name,
	}
}

func sendDeleteNamespaceRequestConfig(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "RemoveNamespace")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
	})
}

func init() {
	DeleteNamespaceCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DeleteNamespaceCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)

	DeleteNamespaceCmd.MarkFlagRequired(constants.OrganizationFlag)
	DeleteNamespaceCmd.MarkFlagRequired(constants.NameFlag)
}
