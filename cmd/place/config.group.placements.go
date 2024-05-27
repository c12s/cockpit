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
	placeConfigGroupPlacementsLongDesc  = "This command places configuration group placements based on the input file.\n\n" +
		"Example:\n" +
		"place-config-group-placements --path 'path to yaml of json file'"

	// Flag Constants
	pathFlag = "path"

	// Flag Descriptions
	pathDesc = "Path to the input YAML or JSON file (required)"
)

var (
	path               string
	placementsResponse model.ConfigGroupPlacementsResponse
)

var PlaceConfigGroupPlacementsCmd = &cobra.Command{
	Use:   "group",
	Short: placeConfigGroupPlacementsShortDesc,
	Long:  placeConfigGroupPlacementsLongDesc,
	Run:   executePlaceConfigGroupPlacements,
}

func executePlaceConfigGroupPlacements(cmd *cobra.Command, args []string) {
	requestBody, err := preparePlacementsRequestConfig()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendPlacementsRequest(requestBody); err != nil {
		fmt.Printf("Error placing configuration group placements: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	render.HandleConfigPlacementsResponse(&placementsResponse)
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
		Response:    &placementsResponse,
		Timeout:     10 * time.Second,
	})
}

func init() {
	PlaceConfigGroupPlacementsCmd.Flags().StringVarP(&path, pathFlag, "p", "", pathDesc)
	PlaceConfigGroupPlacementsCmd.MarkFlagRequired(pathFlag)
}
