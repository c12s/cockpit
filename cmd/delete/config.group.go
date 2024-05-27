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
	deleteConfigGroupShortDesc = "Delete a configuration group version"
	deleteConfigGroupLongDesc  = "This command deletes a specified configuration group version\n" +
		"and displays the deleted configuration group details in JSON or YAML format.\n\n" +
		"Example:\n" +
		"delete-config-group --org 'org' --name 'app_config' --version 'v1.0.0'"

	// Flag Constants
	flagName   = "name"
	flagOutput = "output"

	// Flag Shorthand Constants
	shortFlagName   = "n"
	shortFlagOutput = "o"

	// Flag Descriptions
	descName   = "Configuration group name (required)"
	descOutput = "Output format (json or yaml)"
)

var (
	name           string
	output         string
	deleteResponse model.ConfigGroup
)

var DeleteConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: deleteConfigGroupShortDesc,
	Long:  deleteConfigGroupLongDesc,
	Run:   executeDeleteConfigGroup,
}

func executeDeleteConfigGroup(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteConfigGroupRequest()

	config := sendDeleteConfigGroupRequest(requestBody)

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.DisplayResponse(&deleteResponse, output, "Config group deleted successfully")
}

func prepareDeleteConfigGroupRequest() interface{} {
	requestBody := model.ConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}
	return requestBody
}

func sendDeleteConfigGroupRequest(requestBody interface{}) model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "DeleteConfigGroup")

	return model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &deleteResponse,
	}
}

func init() {
	DeleteConfigGroupCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	DeleteConfigGroupCmd.Flags().StringVarP(&name, flagName, shortFlagName, "", descName)
	DeleteConfigGroupCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)
	DeleteConfigGroupCmd.Flags().StringVarP(&output, flagOutput, shortFlagOutput, "yaml", descOutput)

	DeleteConfigGroupCmd.MarkFlagRequired(flagOrganization)
	DeleteConfigGroupCmd.MarkFlagRequired(flagName)
	DeleteConfigGroupCmd.MarkFlagRequired(flagVersion)
}
