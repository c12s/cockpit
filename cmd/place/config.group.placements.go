package cmd

import (
	"fmt"
	"os"
	"strings"
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
	path                          string
	groupConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var PlaceConfigGroupPlacementsCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   constants.PlaceConfigGroupPlacementsShortDesc,
	Long:    constants.PlaceConfigGroupPlacementsLongDesc,
	Run:     executePlaceConfigGroupPlacements,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
}

func executePlaceConfigGroupPlacements(cmd *cobra.Command, args []string) {
	requestBody, err := preparePlacementsRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendPlacementsRequest(requestBody); err != nil {
		fmt.Println("Error sending config group placements request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(groupConfigPlacementsResponse.Tasks)
}

func preparePlacementsRequestConfig() (interface{}, error) {
	var requestBody model.PlaceConfigGroupPlacementsRequest
	var err error
	if strings.HasSuffix(path, ".yaml") {
		err = utils.ReadYAML(path, &requestBody)
	} else if strings.HasSuffix(path, ".json") {
		err = utils.ReadJSON(path, &requestBody)
	} else {
		return nil, fmt.Errorf("invalid file format. Please provide a YAML or JSON file")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}

	return requestBody, nil
}

func sendPlacementsRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "PlaceConfigGroup")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		Token:       token,
		RequestBody: requestBody,
		Response:    &groupConfigPlacementsResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	PlaceConfigGroupPlacementsCmd.Flags().StringVarP(&path, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)
	PlaceConfigGroupPlacementsCmd.MarkFlagRequired(constants.FilePathFlag)
}
