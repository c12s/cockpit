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
	searchResp  model.LayersResp
	description string
	tags        map[string]string
)

var SearchStarmap = &cobra.Command{
	Use:   "search",
	Short: constants.SearchStarmapShortDesc,
	Long:  constants.SearchStarmapLongDesc,
	Run:   executeSearchStarmap,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NameFlag, constants.DescriptionFlag, constants.TagFlag})
	},
}

func executeSearchStarmap(cmd *cobra.Command, args []string) {
	req := prepareSearchStarmapReq()

	if err := sendSearchRequest(req); err != nil {
		fmt.Println("Error sending search starmap request", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(searchResp, "yaml", "")
}

func prepareSearchStarmapReq() interface{} {
	requestBody := model.GetStarmapChartATimelineReq{
		ChartId:    chartId,
		Maintainer: maintainer,
		Namespace:  namespace,
	}

	return requestBody
}

func sendSearchRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "Search")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &searchResp,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetStarmapChartTimeline.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	GetStarmapChartTimeline.Flags().StringVarP(&description, constants.DescriptionFlag, "", "", constants.ChartDescription)
	GetStarmapChartTimeline.Flags().StringToStringVarP(&tags, constants.TagFlag, "", nil, constants.StarmapLayerTags)
}
