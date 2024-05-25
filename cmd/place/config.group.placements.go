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
	config := createPlacementsRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	fmt.Println()
	render.HandleConfigGroupPlacementsResponse(&placementsResponse)
}

func createPlacementsRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "PlaceConfigGroup")

	var requestBody model.PlaceConfigGroupPlacementsRequest
	if strings.HasSuffix(path, ".yaml") {
		requestBody, err = readYAML(path)
	} else if strings.HasSuffix(path, ".json") {
		requestBody, err = readJSON(path)
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
		Response:    &placementsResponse,
	}
}

func readYAML(filePath string) (model.PlaceConfigGroupPlacementsRequest, error) {
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

func readJSON(filePath string) (model.PlaceConfigGroupPlacementsRequest, error) {
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
	PlaceConfigGroupPlacementsCmd.Flags().StringVarP(&path, pathFlag, "p", "", pathDesc)
	PlaceConfigGroupPlacementsCmd.MarkFlagRequired(pathFlag)
}
