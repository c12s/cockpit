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
	standaloneConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var PlaceStandaloneConfigPlacementsCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   constants.PlaceStandaloneConfigPlacementsShortDesc,
	Long:    constants.PlaceStandaloneConfigPlacementsLongDesc,
	Run:     executePlaceStandaloneConfigPlacements,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
}

func executePlaceStandaloneConfigPlacements(cmd *cobra.Command, args []string) {
	requestBody, err := prepareStandaloneConfigPlacementsRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendStandaloneConfigPlacementsRequest(requestBody); err != nil {
		fmt.Println("Error sending standalone configuration request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(standaloneConfigPlacementsResponse.Tasks)
}

func prepareStandaloneConfigPlacementsRequestConfig() (interface{}, error) {
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

func sendStandaloneConfigPlacementsRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "PlaceStandaloneConfig")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		Token:       token,
		RequestBody: requestBody,
		Response:    &standaloneConfigPlacementsResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	PlaceStandaloneConfigPlacementsCmd.Flags().StringVarP(&path, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)
	PlaceStandaloneConfigPlacementsCmd.MarkFlagRequired(constants.FilePathFlag)
}
