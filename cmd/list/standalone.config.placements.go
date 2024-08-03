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
	standaloneConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var ListStandaloneConfigPlacementsCmd = &cobra.Command{
	Use:     "placements",
	Aliases: aliases.PlacementAliases,
	Short:   constants.ListStandaloneConfigPlacementsShortDesc,
	Long:    constants.ListStandaloneConfigPlacementsLongDesc,
	Run:     executeListStandaloneConfigPlacements,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NamespaceFlag, constants.OrganizationFlag, constants.NameFlag, constants.VersionFlag})
	},
}

func executeListStandaloneConfigPlacements(cmd *cobra.Command, args []string) {
	requestBody := prepareStandalonePlacementsRequestConfig()

	if err := sendStandaloneConfigPlacementsRequest(requestBody); err != nil {
		fmt.Println("Error sending standalone config request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(standaloneConfigPlacementsResponse.Tasks)
}

func prepareStandalonePlacementsRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Name:         name,
		Namespace:    namespace,
		Organization: organization,
		Version:      version,
	}

	return requestBody
}

func sendStandaloneConfigPlacementsRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "ListPlacementTaskByStandaloneConfig")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &standaloneConfigPlacementsResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&organization, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&version, constants.VersionFlag, constants.VersionShorthandFlag, "", constants.VersionDescription)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(constants.OrganizationFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(constants.NameFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(constants.VersionFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(constants.NamespaceFlag)
}
