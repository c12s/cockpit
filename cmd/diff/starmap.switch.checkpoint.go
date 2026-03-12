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
	chartId                  string
	maintainer               string
	oldVersion               string
	newVersion               string
	layers                   []string
	switchCheckpointResponse model.SwitchCheckpointResp
)

var SwitchCheckpointCmd = &cobra.Command{
	Use:     "swchp",
	Aliases: aliases.StarmapSwitchCheckpointAliases,
	Short:   constants.StarmapSwitchCheckpointShortDesc,
	Long:    constants.StarmapSwitchCheckpointLongDesc,
	Run:     executeSwitchCheckpoint,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlagsAny(cmd, []string{constants.IdFlag, constants.MaintainerFlag, constants.NamespaceFlag, constants.LayersFlag, constants.OldVersionFlag, constants.NewVersionFlag})
	},
}

func executeSwitchCheckpoint(cmd *cobra.Command, args []string) {
	req := prepareSwitchCheckpointRequest()

	if err := sendSwitchCheckpointRequest(req); err != nil {
		fmt.Println("Error sending switch checkpoint starmap request:", err)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("Starmap switch checkpoint successfully checked!")
	render.DisplayResponseAsJSONOrYAML(switchCheckpointResponse, "yaml", "")
}

func prepareSwitchCheckpointRequest() interface{} {
	requestBody := model.SwitchCheckpointReq{
		ChartId:    chartId,
		Maintainer: maintainer,
		Namespace:  namespace,
		OldVersion: oldVersion,
		NewVersion: newVersion,
		Layers:     layers,
	}
	return requestBody
}

func sendSwitchCheckpointRequest(requestBody interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "SwitchCheckpoint")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: requestBody,
		Token:       token,
		Timeout:     30 * time.Second,
		Response:    &switchCheckpointResponse,
	})
}

func init() {
	SwitchCheckpointCmd.Flags().StringVarP(&chartId, constants.IdFlag, "", "", constants.ChartIdDescription)
	SwitchCheckpointCmd.Flags().StringVarP(&namespace, constants.NamespaceFlag, constants.NamespaceShorthandFlag, "", constants.NamespaceDescription)
	SwitchCheckpointCmd.Flags().StringVarP(&maintainer, constants.MaintainerFlag, constants.MaintainerShorthandFlag, "", constants.MaintainerDescription)
	SwitchCheckpointCmd.Flags().StringVarP(&oldVersion, constants.OldVersionFlag, "", "", constants.OldChartVersion)
	SwitchCheckpointCmd.Flags().StringVarP(&newVersion, constants.NewVersionFlag, "", "", constants.NewChartVersion)
	SwitchCheckpointCmd.Flags().StringSliceVarP(&layers, constants.LayersFlag, "", nil, constants.ChartLayers)

	SwitchCheckpointCmd.MarkFlagRequired(constants.IdFlag)
	SwitchCheckpointCmd.MarkFlagRequired(constants.NamespaceFlag)
	SwitchCheckpointCmd.MarkFlagRequired(constants.MaintainerFlag)
	SwitchCheckpointCmd.MarkFlagRequired(constants.OldVersionFlag)
	SwitchCheckpointCmd.MarkFlagRequired(constants.NewVersionFlag)
	SwitchCheckpointCmd.MarkFlagRequired(constants.LayersFlag)
}
