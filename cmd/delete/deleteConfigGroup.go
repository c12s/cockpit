package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	deleteConfigGroupShortDesc = "Delete a configuration group version"
	deleteConfigGroupLongDesc  = "This command deletes a specified configuration group version\n" +
		"and displays the deleted configuration group details in JSON format.\n\n" +
		"Example:\n" +
		"delete-config-group --org 'org' --name 'app_config' --version 'v1.0.0'"

	// Flag Constants
	flagName = "name"

	// Flag Shorthand Constants
	shortFlagName = "n"

	// Flag Descriptions
	descName = "Configuration group name (required)"
)

var (
	name           string
	deleteResponse model.ConfigGroup
)

var DeleteConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: deleteConfigGroupShortDesc,
	Long:  deleteConfigGroupLongDesc,
	Run:   executeDeleteConfigGroup,
}

func executeDeleteConfigGroup(cmd *cobra.Command, args []string) {
	config := createDeleteConfigGroupRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	displayResponseAsJSON(&deleteResponse)
}

func createDeleteConfigGroupRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "DeleteConfigGroup")

	requestBody := model.ConfigGroupReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}

	return model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &deleteResponse,
	}
}

func displayResponseAsJSON(response *model.ConfigGroup) {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("Error converting response to JSON: %v\n", err)
		return
	}
	fmt.Println("Deleted Config Group (JSON):")
	fmt.Println(string(jsonData))
	fmt.Println("Configuration Group deleted successfully!")
}

func init() {
	DeleteConfigGroupCmd.Flags().StringVarP(&organization, flagOrg, shortFlagOrg, "", descOrg)
	DeleteConfigGroupCmd.Flags().StringVarP(&name, flagName, shortFlagName, "", descName)
	DeleteConfigGroupCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)

	DeleteConfigGroupCmd.MarkFlagRequired(flagOrg)
	DeleteConfigGroupCmd.MarkFlagRequired(flagName)
	DeleteConfigGroupCmd.MarkFlagRequired(flagVersion)
}
