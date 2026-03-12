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
	responseExtendChart model.PutChartResp
)

var ExtendStarmapCmd = &cobra.Command{
	Use:   "extend",
	Short: constants.ExtendStarmapShortDesc,
	Long:  constants.CreateStarmapLongDesc,
	Run:   executeStarmapExtendChart,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.FilePathFlag})
	},
}

func executeStarmapExtendChart(cmd *cobra.Command, args []string) {
	starmap, err := utils.ReadStarmapFile(path)
	if err != nil {
		fmt.Println("Error reading starchart extend file:", err)
		os.Exit(1)
	}

	config, err := prepareExtendStarmapRequest(starmap)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		fmt.Println("Error sending extend starmap request:", err)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("Starmap chart extended successfully!")
	render.DisplayResponseAsJSONOrYAML(responseExtendChart, "yaml", "")
}

func prepareExtendStarmapRequest(requestBody any) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "Extend")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: requestBody,
		Token:       token,
		Timeout:     30 * time.Second,
		Response:    &responseExtendChart,
	}, nil
}

func init() {
	ExtendStarmapCmd.Flags().StringVarP(&path, constants.FilePathFlag, constants.FilePathShorthandFlag, "", constants.FilePathDescription)

	ExtendStarmapCmd.MarkFlagRequired(constants.FilePathFlag)
}
