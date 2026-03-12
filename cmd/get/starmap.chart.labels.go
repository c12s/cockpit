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
	labels                 map[string]string = make(map[string]string)
	chartsByLabelsResponse model.GetChartsLabelsResp
)

var GetStarmapChartsByLabels = &cobra.Command{
	Use:   "labels",
	Short: constants.GetStarmapChartsByLabelsShortDesc,
	Long:  constants.GetStarmapLabelsLongDesc,
	Run:   executeGetChartsByLabels,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlagsAny(cmd, []string{constants.MaintainerFlag, constants.NamespaceFlag, constants.LabelsFlag})
	},
}

func executeGetChartsByLabels(cmd *cobra.Command, args []string) {
	req := prepareGetChartsLabelsReq()

	if err := sendGetByLabelsRequest(req); err != nil {
		fmt.Println("Error sending get chart by id request", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(chartsByLabelsResponse, "yaml", "")
}

func prepareGetChartsLabelsReq() interface{} {
	requestBody := model.GetStarmapChartsByLabelsReq{
		Labels:     labels,
		Maintainer: maintainer,
		Namespace:  namespace,
	}

	return requestBody
}

func sendGetByLabelsRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetChartsLabels")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &chartsByLabelsResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetStarmapChartsByLabels.Flags().StringVarP(&maintainer, constants.MaintainerFlag, constants.MaintainerShorthandFlag, "", constants.MaintainerDescription)
	GetStarmapChartsByLabels.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	GetStarmapChartsByLabels.Flags().StringToStringVarP(&labels, constants.LabelsFlag, constants.LabelsShorthandFlag, nil, constants.ChartLabelsDescription)

	GetStarmapChartsByLabels.MarkFlagRequired(constants.MaintainerFlag)
	GetStarmapChartsByLabels.MarkFlagRequired(constants.LabelsFlag)
	GetStarmapChartsByLabels.MarkFlagRequired(constants.NamespaceFlag)
}
