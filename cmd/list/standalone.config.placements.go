package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	listStandaloneConfigPlacementsShortDesc = "Retrieve and display the standalone configuration placements"
	listStandaloneConfigPlacementsLongDesc  = `This command retrieves all standalone configuration placements from a specified organization,
displays them in a nicely formatted way, and allows you to see the placements in detail.

Examples:
- cockpit list standalone config placements --org 'org' --name 'app_config' --version 'v1.0.0'
- cockpit list standalone config placements --org 'org' --name 'db_config' --version 'v2.0.0'`
)

var (
	standaloneConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var ListStandaloneConfigPlacementsCmd = &cobra.Command{
	Use:     "placements",
	Aliases: []string{"placement", "placementss"},
	Short:   listStandaloneConfigPlacementsShortDesc,
	Long:    listStandaloneConfigPlacementsLongDesc,
	Run:     executeListStandaloneConfigPlacements,
}

func executeListStandaloneConfigPlacements(cmd *cobra.Command, args []string) {
	requestBody := prepareStandalonePlacementsRequestConfig()

	if err := sendStandaloneConfigPlacementsRequest(requestBody); err != nil {
		fmt.Println("Error sending standalone config request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(standaloneConfigPlacementsResponse.Tasks)
	println()
}

func prepareStandalonePlacementsRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Name:         name,
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
		Timeout:     10 * time.Second,
	})
}

func init() {
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&name, nameFlag, nameShorthandFlag, "", nameDescription)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)

	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(organizationFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(nameFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(versionFlag)
}
