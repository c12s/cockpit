package cmd

import (
	"fmt"
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
	placeConfigGroupPlacementsShortDesc = "Place configuration group placements"
	placeConfigGroupPlacementsLongDesc  = `This command places configuration group placements based on the input file.
The input file should be in either YAML or JSON format, containing the details of the configuration group placements.
It reads the file, processes the placements, and applies them accordingly.

Example:
cockpit place config group placements --path 'path to yaml or json file'`

	// Flag Constants
	pathFlag = "path"

	// Flag Shorthand Constants
	pathShorthandFlag = "p"

	// Flag Descriptions
	pathDescription = "Path to the input YAML or JSON file (required)"
)

var (
	path                          string
	groupConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var PlaceConfigGroupPlacementsCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"grp", "gr"},
	Short:   placeConfigGroupPlacementsShortDesc,
	Long:    placeConfigGroupPlacementsLongDesc,
	Run:     executePlaceConfigGroupPlacements,
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
	println()
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
		Timeout:     10 * time.Second,
	})
}

func init() {
	PlaceConfigGroupPlacementsCmd.Flags().StringVarP(&path, pathFlag, pathShorthandFlag, "", pathDescription)
	PlaceConfigGroupPlacementsCmd.MarkFlagRequired(pathFlag)
}
