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
	timelineResponse model.GetChartsLabelsResp
)

var GetStarmapChartTimeline = &cobra.Command{
	Use:   "timeline",
	Short: constants.GetStarmapChartTimelineShortDesc,
	Long:  constants.GetStarmapChartTimelineLongDesc,
	Run:   executeChartTimeline,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.IdFlag, constants.MaintainerFlag, constants.NamespaceFlag})
	},
}

func executeChartTimeline(cmd *cobra.Command, args []string) {
	req := prepareChartTimelineReq()

	if err := sendTimelineRequest(req); err != nil {
		fmt.Println("Error sending get chart timeline request", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(timelineResponse, "yaml", "")
}

func prepareChartTimelineReq() interface{} {
	requestBody := model.GetStarmapChartATimelineReq{
		ChartId:    chartId,
		Maintainer: maintainer,
		Namespace:  namespace,
	}

	return requestBody
}

func sendTimelineRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "Timeline")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &timelineResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetStarmapChartTimeline.Flags().StringVarP(&maintainer, constants.MaintainerFlag, constants.MaintainerShorthandFlag, "", constants.MaintainerDescription)
	GetStarmapChartTimeline.Flags().StringVarP(&chartId, constants.IdFlag, "", "", constants.ChartIdDescription)
	GetStarmapChartTimeline.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)

	GetStarmapChartTimeline.MarkFlagRequired(constants.MaintainerFlag)
	GetStarmapChartTimeline.MarkFlagRequired(constants.NameFlag)
	GetStarmapChartTimeline.MarkFlagRequired(constants.NamespaceFlag)
}
