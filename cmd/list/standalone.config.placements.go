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
	listStandaloneConfigPlacementsLongDesc  = "This command retrieves the standalone configuration placements from a specified endpoint\n" +
		"and displays them in a nicely formatted way.\n\n" +
		"Example:\n" +
		"list-standalone-config-placements --org 'c12s' --name 'app_config' --version 'v1.0.0'"
)

var (
	standaloneConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var ListStandaloneConfigPlacementsCmd = &cobra.Command{
	Use:   "placements",
	Short: listStandaloneConfigPlacementsShortDesc,
	Long:  listStandaloneConfigPlacementsLongDesc,
	Run:   executeListStandaloneConfigPlacements,
}

func executeListStandaloneConfigPlacements(cmd *cobra.Command, args []string) {
	requestBody, err := prepareStandalonePlacementsRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendStandalonePlacementsRequest(requestBody); err != nil {
		fmt.Printf("Error retrieving standalone configuration placements: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	render.HandleConfigPlacementsResponse(&standaloneConfigPlacementsResponse)
}

func prepareStandalonePlacementsRequestConfig() (interface{}, error) {
	requestBody := model.ConfigReference{
		Name:         name,
		Organization: organization,
		Version:      version,
	}

	return requestBody, nil
}

func sendStandalonePlacementsRequest(requestBody interface{}) error {
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
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&organization, orgFlag, orgFlagShortHand, "", orgDesc)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&name, nameFlag, nameFlagShortHand, "", nameDesc)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&version, versionFlag, versionFlagShortHand, "", versionDesc)

	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(orgFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(nameFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(versionFlag)
}
