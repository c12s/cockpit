package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

const metricsBaseURL = "http://localhost:8086/api/metrics-api/latest-node-data/"
const clusterMetricsBaseURL = "http://localhost:8086/api/metrics-api/latest-cluster-data/"

var (
	nodeID    string
	clusterID string
	all       bool
	sortBy    string
)

var LatestMetricsCmd = &cobra.Command{
	Use:     "metrics",
	Aliases: aliases.MetricsAliases,
	Short:   constants.LatestMetricsShortDesc,
	Long:    constants.LatestMetricsLongDesc,
	Run:     executeLatestMetrics,
	// PreRunE: func(cmd *cobra.Command, args []string) error {
	// 	return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag})
	// },
}

func executeLatestMetrics(cmd *cobra.Command, args []string) {
	if nodeID == "" && clusterID == "" {
		fmt.Println("Either node ID or cluster ID are required")
		os.Exit(1)
	}

	url := metricsBaseURL + nodeID
	infraType := "Node"
	cluster := false
	if clusterID != "" {
		url = clusterMetricsBaseURL + clusterID
		infraType = "Cluster"
		cluster = true
	}

	metricsResponse, err := fetchMetrics(url)
	if err != nil {
		fmt.Println("Error fetching metrics:", err)
		os.Exit(1)
	}
	render.RenderNodeMetrics(metricsResponse, sortBy, infraType)

	if all {
		fmt.Println()
		render.RenderServiceMetrics(metricsResponse, sortBy, cluster)
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
	LatestMetricsCmd.Flags().StringVarP(&clusterID, constants.ClusterIdFlag, constants.ClusterIdShorthandFlag, "", constants.ClusterIdDescription)
	LatestMetricsCmd.Flags().BoolVarP(&all, constants.AllServicesFlag, constants.AllServicesShorthandFlag, false, constants.AllServicesDescription)
	LatestMetricsCmd.Flags().StringVarP(&sortBy, constants.SortByFlag, constants.SortShorthandFlag, "cpu", constants.SortMetricsDescription)
}
