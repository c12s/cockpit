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
	listConfigGroupShortDesc = "Retrieve and display the configuration groups"
	listConfigGroupLongDesc  = "This command retrieves the configuration groups from a specified endpoint\n" +
		"displays them in a nicely formatted way, and saves them to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"list-config-group --organization 'org'"

	// Flag Constants
	outputFlag = "output"

	// Flag Shorthand Constants
	outputFlagShortHand = "o"

	// Flag Descriptions
	outputDesc = "Output format (yaml or json)"

	// Path to files
	listConfigFilePathJSON = "./response/config-group/list-config.json"
	listConfigFilePathYAML = "./response/config-group/list-config.yaml"
)

var (
	organization        string
	configGroupResponse model.ConfigGroupsResponse
	outputFormat        string
)

var ListConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: listConfigGroupShortDesc,
	Long:  listConfigGroupLongDesc,
	Run:   executeListConfigGroup,
}

func executeListConfigGroup(cmd *cobra.Command, args []string) {
	config := createListRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.RenderResponseToYAMLOrJSON(config.Response.(*model.ConfigGroupsResponse), outputFormat)

	filePath := listConfigFilePathYAML
	if outputFormat == "json" {
		filePath = listConfigFilePathJSON
	}

	err = utils.SaveConfigResponseToFile(config.Response.(*model.ConfigGroupsResponse), filePath)
	if err != nil {
		log.Fatalf("Failed to save response to files: %v", err)
	}
}

func createListRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ListConfigGroup")

	requestBody := map[string]string{
		"organization": organization,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &configGroupResponse,
	}
}

func init() {
	ListConfigGroupCmd.Flags().StringVarP(&organization, orgFlag, orgFlagShortHand, "", orgDesc)
	ListConfigGroupCmd.Flags().StringVarP(&outputFormat, outputFlag, outputFlagShortHand, "yaml", outputDesc)

	ListConfigGroupCmd.MarkFlagRequired(orgFlag)
}
