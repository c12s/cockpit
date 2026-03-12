package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	id            string
	kind          string
	maintainer    string
	schemaVersion string
)

var DeleteStarmapChartCmd = &cobra.Command{
	Use:   "delete",
	Short: constants.DeleteStarmapChartShortDesc,
	Long:  constants.DeleteStarmapChartLongDesc,
	Run:   executeDeleteStarmapChart,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.IdFlag, constants.MaintainerFlag, constants.NameFlag, constants.NamespaceFlag, constants.SchemaVersionFlag, constants.KindFlag})
	},
}

func executeDeleteStarmapChart(cmd *cobra.Command, args []string) {
	requestBody := prepareDeleteStarmapChartRequest()

	if err := sendDeleteStarmapChartRequest(requestBody); err != nil {
		fmt.Println("Error sending delete chart request:", err)
		os.Exit(1)
	}

	fmt.Println("Chart deleted successfully!")
}

func prepareDeleteStarmapChartRequest() interface{} {
	requestBody := model.DeleteChartReq{
		Id:            id,
		Maintainer:    maintainer,
		Namespace:     namespace,
		Name:          name,
		Kind:          kind,
		SchemaVersion: schemaVersion,
	}

	return requestBody
}

func sendDeleteStarmapChartRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "DeleteChart")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		Method:      "DELETE",
		URL:         url,
		Token:       token,
		Timeout:     30 * time.Second,
		RequestBody: requestBody,
	})
}

func init() {
	DeleteStarmapChartCmd.Flags().StringVarP(&id, constants.IdFlag, "", "", constants.ChartIdDescription)
	DeleteStarmapChartCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	DeleteStarmapChartCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.ChartNameDescription)
	DeleteStarmapChartCmd.Flags().StringVarP(&kind, constants.KindFlag, constants.KindsShorthandFlag, "", constants.ChartKindDescription)
	DeleteStarmapChartCmd.Flags().StringVarP(&maintainer, constants.MaintainerFlag, constants.MaintainerShorthandFlag, "", constants.MaintainerDescription)
	DeleteStarmapChartCmd.Flags().StringVarP(&schemaVersion, constants.SchemaVersionFlag, constants.VersionShorthandFlag, "", constants.ChartVersionDescription)

	DeleteStarmapChartCmd.MarkFlagRequired(constants.IdFlag)
	DeleteStarmapChartCmd.MarkFlagRequired(constants.NamespaceFlag)
	DeleteStarmapChartCmd.MarkFlagRequired(constants.NameFlag)
	DeleteStarmapChartCmd.MarkFlagRequired(constants.KindFlag)
	DeleteStarmapChartCmd.MarkFlagRequired(constants.MaintainerFlag)
	DeleteStarmapChartCmd.MarkFlagRequired(constants.SchemaVersionFlag)
}
