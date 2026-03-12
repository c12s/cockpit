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
	allChartsResp model.GetChartsLabelsResp
)

var GetStarmapAllCharts = &cobra.Command{
	Use:   "charts",
	Short: constants.GetStarmapAllChartsShortDesc,
	Long:  constants.GetStarmapAllChartsLongDesc,
	Run:   executeGetAllCharts,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{})
	},
}

func executeGetAllCharts(cmd *cobra.Command, args []string) {
	req := ""

	if err := sendGetAllChartsRequest(req); err != nil {
		fmt.Println("Error sending get all charts", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(allChartsResp, "yaml", "")
}

func sendGetAllChartsRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetCharts")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &allChartsResp,
		Timeout:     30 * time.Second,
	})
}

func init() {}
