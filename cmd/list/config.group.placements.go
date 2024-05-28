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
	listConfigGroupPlacementsShortDesc = "Retrieve and display the configuration group placements"
	listConfigGroupPlacementsLongDesc  = "This command retrieves the configuration group placements from a specified organization \n" +
		"displays them in a nicely formatted way.\n\n" +
		"Example:\n" +
		"cockpit list config group placements --org 'org' --name 'app_config' --version 'v1.0.0'"

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
	requestBody, err := preparePlacementsRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendPlacementsRequest(requestBody); err != nil {
		fmt.Printf("Error retrieving configuration group placements: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	render.HandleConfigPlacementsResponse(&placementsResponse)
}

func preparePlacementsRequestConfig() (interface{}, error) {
	requestBody := model.ConfigReference{
		Name:         name,
		Organization: organization,
		Version:      version,
	}

	return requestBody, nil
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
		Response:    &placementsResponse,
		Timeout:     10 * time.Second,
	})
}

func init() {
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&organization, orgFlag, orgFlagShortHand, "", orgDesc)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&name, nameFlag, nameFlagShortHand, "", nameDesc)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&version, versionFlag, versionFlagShortHand, "", versionDesc)

	ListConfigGroupPlacementsCmd.MarkFlagRequired(orgFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(nameFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(versionFlag)
}
