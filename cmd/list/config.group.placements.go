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
	name                          string
	version                       string
	groupConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var ListConfigGroupPlacementsCmd = &cobra.Command{
	Use:     "placements",
	Aliases: aliases.PlacementAliases,
	Short:   constants.ListConfigGroupPlacementsShortDesc,
	Long:    constants.ListConfigGroupPlacementsLongDesc,
	Run:     executeListConfigGroupPlacements,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NameFlag, constants.VersionFlag})
	},
}

func executeListConfigGroupPlacements(cmd *cobra.Command, args []string) {
	requestBody := preparePlacementsRequestConfig()

	if err := sendPlacementsRequest(requestBody); err != nil {
		fmt.Println("Error sending config group placements request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(groupConfigPlacementsResponse.Tasks)
}

func preparePlacementsRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Name:         name,
		Namespace:    namespace,
		Organization: organization,
		Version:      version,
	}

	return requestBody
}

func sendPlacementsRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "ListPlacementTaskByConfigGroup")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &groupConfigPlacementsResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	ListConfigGroupPlacementsCmd.MarkFlagRequired(constants.OrganizationFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(constants.NameFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(constants.VersionFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(constants.NamespaceFlag)
}
