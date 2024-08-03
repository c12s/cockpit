package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"

	"github.com/spf13/cobra"
)

var (
	standaloneConfigPutResponse model.StandaloneConfig
)
var PutStandaloneConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: aliases.ConfigAliases,
	Short:   constants.PutStandaloneConfigShortDesc,
	Long:    constants.PutStandaloneConfigLongDesc,
	Run:     executePutStandaloneConfig,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
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
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &standaloneConfigPutResponse,
	})
}

func init() {
	PutStandaloneConfigCmd.Flags().StringVarP(&filePath, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)
	PutStandaloneConfigCmd.MarkFlagRequired(constants.FilePathFlag)
}
