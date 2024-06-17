package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	putConfigGroupShortDesc = "Send a configuration group to the server"
	putConfigGroupLongDesc  = `This command sends a configuration group read from a file (JSON or YAML) to the server.
It processes the file and uploads the configuration group, displaying the server's response in the same format as the input file.

Example:
- cockpit put config group --path 'path to yaml or JSON file'`
)

var (
	filePath               string
	configGroupPutResponse model.ConfigGroup
)

var PutConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"grp", "gr"},
	Short:   putConfigGroupShortDesc,
	Long:    putConfigGroupLongDesc,
	Run:     executePutConfigGroup,
}

func executePutConfigGroup(cmd *cobra.Command, args []string) {
	configData, err := utils.PrepareRequestBodyFromYAMLOrJSON(filePath)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendConfigGroupData(configData); err != nil {
		fmt.Println("Error sending config group request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(configGroupPutResponse)
	println()
}

func sendConfigGroupData(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "PutConfigGroup")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "POST",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &configGroupPutResponse,
	})
}

func init() {
	PutConfigGroupCmd.Flags().StringVarP(&filePath, "path", "p", "", "Path to the configuration file (required)")
	PutConfigGroupCmd.MarkFlagRequired("path")
}
