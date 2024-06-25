package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func RenderNodeMetrics(metrics model.MetricResponse, sortBy string) {
	nodeMetrics := metrics.FilterNodeMetrics()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	// Print headers with better alignment
	fmt.Fprintf(w, "Service\tMetric\tTotal\tUsed\tAvailable\n")

	metricsMap := make(map[string]map[string]float64)

	for _, data := range nodeMetrics {
		metricName := data.Metric["__name__"]
		for _, value := range data.Values {
			val := utils.StringToFloat(value[1].(string))
			if metricsMap[metricName] == nil {
				metricsMap[metricName] = make(map[string]float64)
			}
			if strings.Contains(metricName, "total") {
				metricsMap[metricName]["total"] = val
			} else if strings.Contains(metricName, "usage") {
				metricsMap[metricName]["used"] = val
			}
		}
	}

	printNodeMetrics(w, "Node", "cpu", metricsMap["custom_node_cpu_usage_percentage_total"], "%")
	printNodeMetrics(w, "Node", "disk", metricsMap["custom_node_disk_total_gb_total"], "GB", metricsMap["custom_node_disk_usage_gb_total"])
	printNodeMetrics(w, "Node", "memory", metricsMap["custom_node_ram_total_mb_total"], "MB", metricsMap["custom_node_ram_usage_mb_total"])
}

func RenderServiceMetrics(metrics model.MetricResponse, sortBy string) {
	serviceMetrics := metrics.FilterServiceMetrics()

	// Sort service metrics based on the sortBy flag
	sort.Slice(serviceMetrics, func(i, j int) bool {
		if sortBy == "memory" {
			return utils.GetServiceMemoryUsage(serviceMetrics[i]) > utils.GetServiceMemoryUsage(serviceMetrics[j])
		}
		return utils.GetServiceCpuUsage(serviceMetrics[i]) > utils.GetServiceCpuUsage(serviceMetrics[j])
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	// Print headers for service metrics
	fmt.Fprintf(w, "Service\tCPU\tTotal Memory\tUsed Memory\n")

	serviceMap := make(map[string]map[string]float64)

	for _, data := range serviceMetrics {
		service := data.Metric["service_name"]
		if service == "" {
			service = "Node"
		}
		metricName := data.Metric["__name__"]
		for _, value := range data.Values {
			val := utils.StringToFloat(value[1].(string))
			if serviceMap[service] == nil {
				serviceMap[service] = make(map[string]float64)
			}
			if strings.Contains(metricName, "cpu") {
				serviceMap[service]["cpu"] = val
			} else if strings.Contains(metricName, "usage") {
				serviceMap[service]["used_memory"] = val
			}
		}
	}

	services := make([]string, 0, len(serviceMap))
	for service := range serviceMap {
		services = append(services, service)
	}

	sort.Slice(services, func(i, j int) bool {
		if sortBy == "memory" {
			return serviceMap[services[i]]["used_memory"] > serviceMap[services[j]]["used_memory"]
		}
		return serviceMap[services[i]]["cpu"] > serviceMap[services[j]]["cpu"]
	})

	for _, service := range services {
		cpu := fmt.Sprintf("%.2f %%", serviceMap[service]["cpu"])
		usedMemory := fmt.Sprintf("%.2f MB", serviceMap[service]["used_memory"])
		fmt.Fprintf(w, "%s\t%s\tN/A\t%s\n", service, cpu, usedMemory)
	}
}

func printNodeMetrics(w *tabwriter.Writer, service, metric string, totalValues map[string]float64, unit string, usedValues ...map[string]float64) {
	total, used := totalValues["total"], 0.0
	if len(usedValues) > 0 {
		used = usedValues[0]["used"]
	}

	available := total - used

	if metric == "cpu" {
		fmt.Fprintf(w, "%s\t%s\t-\t%.2f %s\t-\n", service, metric, total, unit)
	} else {
		if total == 0 {
			fmt.Fprintf(w, "%s\t%s\t0.0 %s\t0.0 %s\t0.0 %s\n", service, metric, unit, unit, unit)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%.2f %s\t%.2f %s\t%.2f %s\n", service, metric, total, unit, used, unit, available, unit)
		}
	}
}
