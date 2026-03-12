package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	path                string
	responseUpdateChart model.PutChartResp
)

var UpdateStarmapCmd = &cobra.Command{
	Use:   "update",
	Short: constants.UpdateStarmapShortDesc,
	Long:  constants.UpdateStarmapLongDesc,
	Run:   executeCreateStarmap,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
}

func executeCreateStarmap(cmd *cobra.Command, args []string) {
	starmap, err := utils.ReadStarmapFile(path)
	if err != nil {
		fmt.Println("Error reading starchart update file:", err)
		os.Exit(1)
	}

	config, err := prepareUpdateStarmapChartRequest(starmap)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		fmt.Println("Error sending update starmap request:", err)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("Starmap chart updated successfully!")
	render.DisplayResponseAsJSONOrYAML(responseUpdateChart, "yaml", "")
}

func prepareUpdateStarmapChartRequest(requestBody any) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "UpdateChart")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "PATCH",
		RequestBody: requestBody,
		Token:       token,
		Timeout:     30 * time.Second,
		Response:    &responseUpdateChart,
	}, nil
}

func init() {
	UpdateStarmapCmd.Flags().StringVarP(&path, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)

	UpdateStarmapCmd.MarkFlagRequired(constants.FilePathFlag)
}
