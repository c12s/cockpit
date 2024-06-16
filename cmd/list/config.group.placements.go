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
	listConfigGroupPlacementsLongDesc  = `This command retrieves all configuration group placements from a specified organization,
displays them in a nicely formatted way, and allows you to see the placements in detail.

Examples:
- cockpit list config group placements --org 'org' --name 'app_config' --version 'v1.0.0'
- cockpit list config group placements --org 'org' --name 'db_config' --version 'v2.0.0'`

	// Flag Constants
	nameFlag    = "name"
	versionFlag = "version"

	// Flag Shorthand Constants
	nameShorthandFlag    = "n"
	versionShorthandFlag = "v"

	// Flag Descriptions
	nameDescription    = "Configuration group name (required)"
	versionDescription = "Configuration group version (required)"
)

var (
	name                          string
	version                       string
	groupConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var ListConfigGroupPlacementsCmd = &cobra.Command{
	Use:     "placements",
	Aliases: []string{"placement", "placementss"},
	Short:   listConfigGroupPlacementsShortDesc,
	Long:    listConfigGroupPlacementsLongDesc,
	Run:     executeListConfigGroupPlacements,
}

func executeListConfigGroupPlacements(cmd *cobra.Command, args []string) {
	requestBody := preparePlacementsRequestConfig()

	if err := sendPlacementsRequest(requestBody); err != nil {
		fmt.Println("Error sending config group placements request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(groupConfigPlacementsResponse.Tasks)
	println()
}

func preparePlacementsRequestConfig() interface{} {
	requestBody := model.ConfigReference{
		Name:         name,
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
		Timeout:     10 * time.Second,
	})
}

func init() {
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&organization, organizationFlag, organizationShorthandFlag, "", organizationDescription)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&name, nameFlag, nameShorthandFlag, "", nameDescription)
	ListConfigGroupPlacementsCmd.Flags().StringVarP(&version, versionFlag, versionShorthandFlag, "", versionDescription)

	ListConfigGroupPlacementsCmd.MarkFlagRequired(organizationFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(nameFlag)
	ListConfigGroupPlacementsCmd.MarkFlagRequired(versionFlag)
}
