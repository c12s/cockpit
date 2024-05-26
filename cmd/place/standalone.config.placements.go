package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	placeStandaloneConfigPlacementsShortDesc = "Place standalone configuration placements"
	placeStandaloneConfigPlacementsLongDesc  = "This command places standalone configuration placements based on the input file.\n\n" +
		"Example:\n" +
		"place-standalone-config-placements --path 'path to yaml or json file'"
)

var (
	standaloneConfigPlacementsResponse model.ConfigGroupPlacementsResponse
)

var PlaceStandaloneConfigPlacementsCmd = &cobra.Command{
	Use:   "config",
	Short: placeStandaloneConfigPlacementsShortDesc,
	Long:  placeStandaloneConfigPlacementsLongDesc,
	Run:   executePlaceStandaloneConfigPlacements,
}

func executePlaceStandaloneConfigPlacements(cmd *cobra.Command, args []string) {
	config := createStandaloneConfigPlacementsRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	fmt.Println()
	render.HandleConfigGroupPlacementsResponse(&placementsResponse)
}

func createStandaloneConfigPlacementsRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "PlaceStandaloneConfig")

	var requestBody model.PlaceConfigGroupPlacementsRequest
	if strings.HasSuffix(path, ".yaml") {
		requestBody, err = readStandaloneConfigFromYAML(path)
	} else if strings.HasSuffix(path, ".json") {
		requestBody, err = readStandaloneConfigFromJSON(path)
	} else {
		log.Fatalf("Invalid file format. Please provide a YAML or JSON file.")
	}

	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	return model.HTTPRequestConfig{
		Method:      "POST",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &standaloneConfigPlacementsResponse,
	}
}

func readStandaloneConfigFromYAML(filePath string) (model.PlaceConfigGroupPlacementsRequest, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return model.PlaceConfigGroupPlacementsRequest{}, err
	}

	var request model.PlaceConfigGroupPlacementsRequest
	err = yaml.Unmarshal(data, &request)
	if err != nil {
		return model.PlaceConfigGroupPlacementsRequest{}, err
	}

	return request, nil
}

func readStandaloneConfigFromJSON(filePath string) (model.PlaceConfigGroupPlacementsRequest, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return model.PlaceConfigGroupPlacementsRequest{}, err
	}

	var request model.PlaceConfigGroupPlacementsRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		return model.PlaceConfigGroupPlacementsRequest{}, err
	}

	return request, nil
}

func init() {
	PlaceStandaloneConfigPlacementsCmd.Flags().StringVarP(&path, pathFlag, "p", "", pathDesc)
	PlaceStandaloneConfigPlacementsCmd.MarkFlagRequired(pathFlag)
}
