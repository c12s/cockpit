package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	latestMetricsShortDesc = "Retrieve and display the latest metrics"
	latestMetricsLongDesc  = `This command fetches the latest metrics for a specific node and displays them in a nicely formatted way.
The user can specify the node ID to retrieve the metrics. 

Example:
- cockpit get nodes metrics --node-id 'nodeID'`

	// Flag Constants
	nodeIDFlag = "node-id"
	allFlag    = "all-services"
	sortFlag   = "sort"

	// Flag Shorthand Constants
	nodeIDShorthandFlag = "n"
	sortShorthandFlag   = "s"
	allShorthandFlag    = "a"

	// Flag Descriptions
	nodeIDDescription = "Node ID (required)"
	allDescription    = "Display all metrics (optional)"
	sortDescription   = "Sort metrics by 'cpu' or 'memory'"
	metricsBaseURL    = "http://localhost:8086/api/metrics-api/latest-node-data/"
)

var (
	nodeID string
	all    bool
	sortBy string
)

var LatestMetricsCmd = &cobra.Command{
	Use:     "metrics",
	Aliases: aliases.MetricsAliases,
	Short:   latestMetricsShortDesc,
	Long:    latestMetricsLongDesc,
	Run:     executeLatestMetrics,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{nodeIDFlag})
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
	LatestMetricsCmd.Flags().StringVarP(&nodeID, nodeIDFlag, nodeIDShorthandFlag, "", nodeIDDescription)
	LatestMetricsCmd.Flags().BoolVarP(&all, allFlag, allShorthandFlag, false, allDescription)
	LatestMetricsCmd.Flags().StringVarP(&sortBy, sortFlag, sortShorthandFlag, "cpu", sortDescription)

	LatestMetricsCmd.MarkFlagRequired(nodeIDFlag)
}
