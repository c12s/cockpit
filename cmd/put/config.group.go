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
	filePath               string
	configGroupPutResponse model.ConfigGroup
)

var PutConfigGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: aliases.GroupAliases,
	Short:   constants.PutConfigGroupShortDesc,
	Long:    constants.PutConfigGroupLongDesc,
	Run:     executePutConfigGroup,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
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
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
		Response:    &configGroupPutResponse,
	})
}

func init() {
	PutConfigGroupCmd.Flags().StringVarP(&filePath, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)
	PutConfigGroupCmd.MarkFlagRequired(constants.FilePathFlag)
}
