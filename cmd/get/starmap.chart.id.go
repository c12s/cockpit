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
	chartId           string
	chartByIdResponse model.GetChartResp
)

var GetStarmapChartById = &cobra.Command{
	Use:   "id",
	Short: constants.GetStarmapChartByIdShortDesc,
	Long:  constants.GetStarmapIdLongDesc,
	Run:   executeGetChartById,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.IdFlag, constants.MaintainerFlag, constants.NamespaceFlag})
	},
}

func executeGetChartById(cmd *cobra.Command, args []string) {
	req := prepareGetChartIdReq()

	if err := sendGetByIdRequest(req); err != nil {
		fmt.Println("Error sending get chart by id request", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(chartByIdResponse, "yaml", "")
}

func prepareGetChartIdReq() interface{} {
	requestBody := model.GetStarmapChartByIdReq{
		ChartId:       chartId,
		Maintainer:    maintainer,
		Namespace:     namespace,
		SchemaVersion: schemaVersion,
	}

	return requestBody
}

func sendGetByIdRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetChartId")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &chartByIdResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetStarmapChartById.Flags().StringVarP(&maintainer, constants.MaintainerFlag, constants.MaintainerShorthandFlag, "", constants.MaintainerDescription)
	GetStarmapChartById.Flags().StringVarP(&chartId, constants.IdFlag, "", "", constants.ChartIdDescription)
	GetStarmapChartById.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	GetStarmapChartById.Flags().StringVarP(&schemaVersion, constants.SchemaVersionFlag, constants.VersionShorthandFlag, "", constants.ChartVersionDescription)

	GetStarmapChartById.MarkFlagRequired(constants.MaintainerFlag)
	GetStarmapChartById.MarkFlagRequired(constants.NameFlag)
	GetStarmapChartById.MarkFlagRequired(constants.NamespaceFlag)
}
