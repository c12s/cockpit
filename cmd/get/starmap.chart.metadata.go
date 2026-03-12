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
	maintainer              string
	schemaVersion           string
	chartByMetadataResponse model.GetChartResp
)

var GetStarmapChartByMetadata = &cobra.Command{
	Use:     "metadata",
	Aliases: aliases.StarmapGetChartByMetadataAliases,
	Short:   constants.GetStarmapChartByMetadataShortDesc,
	Long:    constants.GetStarmapMetadataLongDesc,
	Run:     executeGetChartByMetadata,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.MaintainerFlag, constants.NameFlag, constants.NamespaceFlag})
	},
}

func executeGetChartByMetadata(cmd *cobra.Command, args []string) {
	req := prepareGetChartMetadataReq()

	if err := sendGetByMetadataRequest(req); err != nil {
		fmt.Println("Error sending get chart by metadata request", err)
		os.Exit(1)
	}

	render.DisplayResponseAsJSONOrYAML(chartByMetadataResponse, "yaml", "")
}

func prepareGetChartMetadataReq() interface{} {
	requestBody := model.GetStarmapChartByMetadataReq{
		Maintainer:    maintainer,
		Name:          name,
		Namespace:     namespace,
		SchemaVersion: schemaVersion,
	}

	return requestBody
}

func sendGetByMetadataRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "GetChartMetadata")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: requestBody,
		Response:    &chartByMetadataResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	GetStarmapChartByMetadata.Flags().StringVarP(&maintainer, constants.MaintainerFlag, constants.MaintainerShorthandFlag, "", constants.MaintainerDescription)
	GetStarmapChartByMetadata.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.ChartNameDescription)
	GetStarmapChartByMetadata.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	GetStarmapChartByMetadata.Flags().StringVarP(&version, constants.SchemaVersionFlag, constants.VersionShorthandFlag, "", constants.ChartVersionDescription)

	GetStarmapChartByMetadata.MarkFlagRequired(constants.MaintainerFlag)
	GetStarmapChartByMetadata.MarkFlagRequired(constants.NameFlag)
	GetStarmapChartByMetadata.MarkFlagRequired(constants.NamespaceFlag)
}
