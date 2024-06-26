package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	placeStandaloneConfigPlacementsShortDesc = "Place standalone configuration placements"
	placeStandaloneConfigPlacementsLongDesc  = `This command places standalone configuration placements based on the input file.
The input file should be in either YAML or JSON format, containing the details of the standalone configuration placements.
It reads the file, processes the placements, and applies them accordingly.

Example:
- cockpit place standalone config placements --path 'path to yaml or json file'`
)

var (
	standaloneConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var PlaceStandaloneConfigPlacementsCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   placeStandaloneConfigPlacementsShortDesc,
	Long:    placeStandaloneConfigPlacementsLongDesc,
	Run:     executePlaceStandaloneConfigPlacements,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{pathFlag})
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
	println()
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
		Timeout:     10 * time.Second,
	})
}

func init() {
	PlaceStandaloneConfigPlacementsCmd.Flags().StringVarP(&path, pathFlag, pathShorthandFlag, "", pathDescription)
	PlaceStandaloneConfigPlacementsCmd.MarkFlagRequired(pathFlag)
}
