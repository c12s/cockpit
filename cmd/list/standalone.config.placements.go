package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"log"
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
	config := createStandalonePlacementsRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	fmt.Println()
	render.HandleConfigPlacementsResponse(&standaloneConfigPlacementsResponse)
}

func createStandalonePlacementsRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListPlacementTaskByStandaloneConfig")

	requestBody := model.ConfigReference{
		Name:         name,
		Organization: organization,
		Version:      version,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &standaloneConfigPlacementsResponse,
	}
}

func init() {
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&organization, orgFlag, orgFlagShortHand, "", orgDesc)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&name, nameFlag, nameFlagShortHand, "", nameDesc)
	ListStandaloneConfigPlacementsCmd.Flags().StringVarP(&version, versionFlag, versionFlagShortHand, "", versionDesc)

	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(orgFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(nameFlag)
	ListStandaloneConfigPlacementsCmd.MarkFlagRequired(versionFlag)
}
