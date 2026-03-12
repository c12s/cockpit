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
	path             string
	responsePutChart model.PutChartResp
)

var PutStarmapCmd = &cobra.Command{
	Use:   "put",
	Short: constants.CreateStarmapShortDesc,
	Long:  constants.CreateStarmapLongDesc,
	Run:   executeCreateStarmap,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
}

func executeCreateStarmap(cmd *cobra.Command, args []string) {
	starmap, err := utils.ReadStarmapFile(path)
	if err != nil {
		fmt.Println("Error reading starchart file:", err)
		os.Exit(1)
	}

	config, err := prepareCreateStarmapRequest(starmap)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		fmt.Println("Error sending put starmap request:", err)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("Starmap chart created successfully!")
	render.DisplayResponseAsJSONOrYAML(responsePutChart, "yaml", "")
}

func prepareCreateStarmapRequest(requestBody any) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "PutChart")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: requestBody,
		Token:       token,
		Timeout:     30 * time.Second,
		Response:    &responsePutChart,
	}, nil
}

func init() {
	PutStarmapCmd.Flags().StringVarP(&path, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)

	PutStarmapCmd.MarkFlagRequired(constants.FilePathFlag)
}
