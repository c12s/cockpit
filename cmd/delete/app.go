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

var DeleteAppCmd = &cobra.Command{
	Use:   "app",
	Short: constants.DeleteAppShortDesc,
	Long:  constants.DeleteAppLongDesc,
	Run:   executeDeleteApp,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.NameFlag, constants.OrganizationFlag})
	},
}

func executeDeleteApp(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteAppRequest()

	if err := sendDeleteAppRequestConfig(requestBody); err != nil {
		fmt.Println("Error sending delete app request:", err)
		os.Exit(1)
	}

	fmt.Println("App deleted successfully!")
}

func prepareDeleteAppRequest() interface{} {
	return struct {
		OrgId     string `json:"orgId"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}{
		OrgId:     organization,
		Name:      name,
		Namespace: namespace,
	}
}

func sendDeleteAppRequestConfig(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "RemoveApp")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
	})
}

func init() {
	DeleteAppCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DeleteAppCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	DeleteAppCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	DeleteAppCmd.MarkFlagRequired(constants.OrganizationFlag)
	DeleteAppCmd.MarkFlagRequired(constants.NameFlag)
	DeleteAppCmd.MarkFlagRequired(constants.NamespaceFlag)
}
