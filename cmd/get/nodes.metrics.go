package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const metricsBaseURL = "http://localhost:8086/api/metrics-api/latest-node-data/"

var (
	nodeID string
	all    bool
	sortBy string
)

var LatestMetricsCmd = &cobra.Command{
	Use:     "metrics",
	Aliases: aliases.MetricsAliases,
	Short:   constants.LatestMetricsShortDesc,
	Long:    constants.LatestMetricsLongDesc,
	Run:     executeLatestMetrics,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag})
	},
}

func executeLatestMetrics(cmd *cobra.Command, args []string) {
	if nodeID == "" {
		fmt.Println("Node ID is required")
		os.Exit(1)
	}

	url := metricsBaseURL + nodeID
	metricsResponse, err := fetchMetrics(url)
	if err != nil {
		fmt.Println("Error fetching metrics:", err)
		os.Exit(1)
	}

	render.RenderNodeMetrics(metricsResponse, sortBy)

	if all {
		fmt.Println()
		render.RenderServiceMetrics(metricsResponse, sortBy)
	}
}

func fetchMetrics(url string) (model.MetricResponse, error) {
	var metricsResponse model.MetricResponse
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return metricsResponse, fmt.Errorf("error reading token: %v", err)
	}

	config := model.HTTPRequestConfig{
		URL:      url,
		Method:   "GET",
		Token:    token,
		Response: &metricsResponse,
		Timeout:  10 * time.Second,
	}

	if err := utils.SendHTTPRequest(config); err != nil {
		return metricsResponse, err
	}

	return metricsResponse, nil
}

func init() {
	LatestMetricsCmd.Flags().StringVarP(&nodeID, constants.NodeIdFlag, constants.NodeIdShorthandFlag, "", constants.NodeIdDescription)
	LatestMetricsCmd.Flags().BoolVarP(&all, constants.AllServicesFlag, constants.AllServicesShorthandFlag, false, constants.AllServicesDescription)
	LatestMetricsCmd.Flags().StringVarP(&sortBy, constants.SortByFlag, constants.SortShorthandFlag, "cpu", constants.SortMetricsDescription)

	LatestMetricsCmd.MarkFlagRequired(constants.NodeIdFlag)
}
