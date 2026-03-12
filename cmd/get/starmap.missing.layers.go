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
	layers                 []string
	chartMissingLayersResp model.GetMissingLayersResp
)

var GetStarmapMissingLayers = &cobra.Command{
	Use:   "mslayers",
	Short: constants.GetStarmapChartMissingLayersShortDesc,
	Long:  constants.GetMissingLayersLongDesc,
	Run:   executeGetChartMissingLayers,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlagsAny(cmd, []string{constants.IdFlag, constants.MaintainerFlag, constants.NamespaceFlag, constants.LayersFlag})
	},
}

func executeGetChartMissingLayers(cmd *cobra.Command, args []string) {
	req := prepareGetChartMissingLayersReq()

	if err := sendGetMissingLayersRequest(req); err != nil {
		fmt.Println("Error sending get missing layers request", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(chartMissingLayersResp, "yaml", "")
}

func prepareGetChartMissingLayersReq() interface{} {
	requestBody := model.GetStarmapChartMissingLayersReq{
		ChartId:       chartId,
		Maintainer:    maintainer,
		Namespace:     namespace,
		Layers:        layers,
		SchemaVersion: schemaVersion,
	}

	return requestBody
}

func sendGetMissingLayersRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetMissingLayers")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &chartMissingLayersResp,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetStarmapMissingLayers.Flags().StringVarP(&chartId, constants.IdFlag, "", "", constants.ChartIdDescription)
	GetStarmapMissingLayers.Flags().StringVarP(&maintainer, constants.MaintainerFlag, constants.MaintainerShorthandFlag, "", constants.MaintainerDescription)
	GetStarmapMissingLayers.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	GetStarmapMissingLayers.Flags().StringSliceVarP(&layers, constants.LayersFlag, "", nil, constants.ChartMissingLayers)
	GetStarmapMissingLayers.Flags().StringVarP(&schemaVersion, constants.SchemaVersionFlag, constants.VersionShorthandFlag, "", constants.ChartVersionDescription)

	GetStarmapMissingLayers.MarkFlagRequired(constants.MaintainerFlag)
	GetStarmapMissingLayers.MarkFlagRequired(constants.NamespaceFlag)
	GetStarmapMissingLayers.MarkFlagRequired(constants.IdFlag)
}
