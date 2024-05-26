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
	listConfigGroupPlacementsShortDesc = "Retrieve and display the configuration group placements"
	listConfigGroupPlacementsLongDesc  = "This command retrieves the configuration group placements from a specified endpoint\n" +
		"displays them in a nicely formatted way.\n\n" +
		"Example:\n" +
		"list-config-group-placements --org 'org' --name 'app_config' --version 'v1.0.0'"

	// Flag Constants
	nameFlag    = "name"
	versionFlag = "version"

	// Flag Shorthand Constants
	nameFlagShortHand    = "n"
	versionFlagShortHand = "v"

	// Flag Descriptions
	nameDesc    = "Configuration group name (required)"
	versionDesc = "Configuration group version (required)"
)

var (
	name               string
	version            string
	placementsResponse model.ConfigGroupPlacementsResponse
)

var ListConfigGroupPlacementsCmd = &cobra.Command{
	Use:   "placements",
	Short: listConfigGroupPlacementsShortDesc,
	Long:  listConfigGroupPlacementsLongDesc,
	Run:   executeListConfigGroupPlacements,
}

func executeListConfigGroupPlacements(cmd *cobra.Command, args []string) {
	config := createPlacementsRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	fmt.Println()
	render.HandleConfigPlacementsResponse(&placementsResponse)
}

func createPlacementsRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListPlacementTaskByConfigGroup")

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
		Response:    &placementsResponse,
	}
}

func init() {
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&organization, orgFlag, orgFlagShortHand, "", orgDesc)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&name, nameFlag, nameFlagShortHand, "", nameDesc)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&version, versionFlag, versionFlagShortHand, "", versionDesc)

	ListConfigGroupPlacementsCmd.MarkFlagRequired(orgFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(nameFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(versionFlag)
}
