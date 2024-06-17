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
	putStandaloneConfigShortDesc = "Send a standalone configuration to the server"
	putStandaloneConfigLongDesc  = `This command sends a standalone configuration read from a file (JSON or YAML) to the server.
It processes the file and uploads the standalone configuration, displaying the server's response in the same format as the input file.

Example:
cockpit put standalone config --path 'path to yaml or JSON file'`
)

var (
	standaloneConfigPutResponse model.StandaloneConfig
)
var PutStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf", "cnfg", "cfg"},
	Short:   putStandaloneConfigShortDesc,
	Long:    putStandaloneConfigLongDesc,
	Run:     executePutStandaloneConfig,
}

func executePutStandaloneConfig(cmd *cobra.Command, args []string) {
	configData, err := utils.PrepareRequestBodyFromYAMLOrJSON(filePath)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendStandaloneConfigData(configData); err != nil {
		fmt.Println("Error sending standalone config request:", err)
		os.Exit(1)
	}

	render.RenderResponseAsTabWriter(standaloneConfigPutResponse)
	println()
}

func sendStandaloneConfigData(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "PutStandaloneConfig")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "POST",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &standaloneConfigPutResponse,
	})
}

func init() {
	PutStandaloneConfigCmd.Flags().StringVarP(&filePath, "path", "p", "", "Path to the configuration file (required)")
	PutStandaloneConfigCmd.MarkFlagRequired("path")
}