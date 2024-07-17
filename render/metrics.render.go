package render

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
)

func RenderNodeMetrics(metrics model.MetricResponse, sortBy string, infraType string) {
	nodeMetrics := metrics.FilterNodeMetrics()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintf(w, "Service\tMetric\tTotal\tUsed\tAvailable\tNetwork Receive\tNetwork Transmit\tBandwidth\n")

	metricsMap := make(map[string]map[string]float64)

	for _, data := range nodeMetrics {
		metricName := data.Metric["__name__"]
		val := utils.StringToFloat(data.Value[1].(string))
		if metricsMap[metricName] == nil {
			metricsMap[metricName] = make(map[string]float64)
		}
		switch {
		case strings.Contains(metricName, "total"):
			metricsMap[metricName]["total"] = val
		case strings.Contains(metricName, "usage"), strings.Contains(metricName, "available"):
			metricsMap[metricName]["used"] = val
		case strings.Contains(metricName, "network_receive"):
			metricsMap[metricName]["network_receive"] = val
		case strings.Contains(metricName, "network_transmit"):
			metricsMap[metricName]["network_transmit"] = val
		}
	}

	printNodeMetrics(w, infraType, "cpu", nil, "%", metricsMap["custom_node_cpu_usage_percentage"])
	printNodeMetrics(w, infraType, "disk", metricsMap["custom_node_disk_total_gb"], "GB", metricsMap["custom_node_disk_usage_gb"])
	printMemoryMetrics(w, infraType, "memory", metricsMap["custom_node_ram_total_mb"], metricsMap["custom_node_ram_available_mb"], "MB")
	printNetworkMetrics(w, infraType, "network", metricsMap["custom_node_network_receive_mb"], metricsMap["custom_node_network_transmit_mb"], "MB")
}

func RenderServiceMetrics(metrics model.MetricResponse, sortBy string, cluster bool) {
	serviceMetrics := metrics.FilterServiceMetrics()

	sort.Slice(serviceMetrics, func(i, j int) bool {
		switch sortBy {
		case "memory":
			return utils.GetServiceMemoryUsage(serviceMetrics[i]) > utils.GetServiceMemoryUsage(serviceMetrics[j])
		case "disk":
			return utils.GetServiceDiskUsage(serviceMetrics[i]) > utils.GetServiceDiskUsage(serviceMetrics[j])
		case "network receive":
			return utils.GetServiceNetworkReceive(serviceMetrics[i]) > utils.GetServiceNetworkReceive(serviceMetrics[j])
		case "network transmit":
			return utils.GetServiceNetworkTransmit(serviceMetrics[i]) > utils.GetServiceNetworkTransmit(serviceMetrics[j])
		case "bandwidth":
			return utils.GetServiceBandwidth(serviceMetrics[i]) > utils.GetServiceBandwidth(serviceMetrics[j])
		default:
			return utils.GetServiceCpuUsage(serviceMetrics[i]) > utils.GetServiceCpuUsage(serviceMetrics[j])
		}
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintf(w, "Service\tCPU\tTotal Memory\tUsed Memory\tDisk Usage\tNetwork Receive\tNetwork Transmit\tBandwidth\n")

	serviceMap := make(map[string]map[string]float64)

	for _, data := range serviceMetrics {
		service := data.Metric["service_name"]
		if strings.Contains(data.Metric["__name__"], "cpu_usage") {
			service = data.Metric["name"]
		}
		if cluster {
			if strings.Contains(data.Metric["__name__"], "cpu_usage") {
				service = data.Metric["nodeID"] + ": " + data.Metric["name"]
			} else {
				service = data.Metric["nodeID"] + ": " + data.Metric["service_name"]
			}
		}
		if service == "" || strings.HasSuffix(service, ": ") {
			continue
		}
		metricName := data.Metric["__name__"]
		val := utils.StringToFloat(data.Value[1].(string))
		if serviceMap[service] == nil {
			serviceMap[service] = make(map[string]float64)
		}
		switch {
		case strings.Contains(metricName, "cpu"):
			serviceMap[service]["cpu"] = val
		case strings.Contains(metricName, "ram_usage"):
			serviceMap[service]["used_memory"] = val
		case strings.Contains(metricName, "disk_usage"):
			serviceMap[service]["disk_usage"] = val
		case strings.Contains(metricName, "network_receive"):
			serviceMap[service]["network_receive"] = val
		case strings.Contains(metricName, "network_transmit"):
			serviceMap[service]["network_transmit"] = val
		}
	}

	services := make([]string, 0, len(serviceMap))
	for service := range serviceMap {
		services = append(services, service)
	}

	sort.Slice(services, func(i, j int) bool {
		switch sortBy {
		case "memory":
			return serviceMap[services[i]]["used_memory"] > serviceMap[services[j]]["used_memory"]
		case "disk":
			return serviceMap[services[i]]["disk_usage"] > serviceMap[services[j]]["disk_usage"]
		case "network_receive":
			return serviceMap[services[i]]["network_receive"] > serviceMap[services[j]]["network_receive"]
		case "network_transmit":
			return serviceMap[services[i]]["network_transmit"] > serviceMap[services[j]]["network_transmit"]
		case "bandwidth":
			return (serviceMap[services[i]]["network_receive"] + serviceMap[services[i]]["network_transmit"]) > (serviceMap[services[j]]["network_receive"] + serviceMap[services[j]]["network_transmit"])
		default:
			return serviceMap[services[i]]["cpu"] > serviceMap[services[j]]["cpu"]
		}
	})

	for _, service := range services {
		cpu := fmt.Sprintf("%.2f %%", serviceMap[service]["cpu"])
		usedMemory := fmt.Sprintf("%.2f MB", serviceMap[service]["used_memory"])
		diskUsage := fmt.Sprintf("%.2f MB", serviceMap[service]["disk_usage"])
		networkReceive := fmt.Sprintf("%.4f MB", serviceMap[service]["network_receive"])
		networkTransmit := fmt.Sprintf("%.4f MB", serviceMap[service]["network_transmit"])
		bandwidth := fmt.Sprintf("%.4f MB", serviceMap[service]["network_receive"]+serviceMap[service]["network_transmit"])
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", service, cpu, "N/A", usedMemory, diskUsage, networkReceive, networkTransmit, bandwidth)
	}
}

func printNodeMetrics(w *tabwriter.Writer, service, metric string, totalValues map[string]float64, unit string, usedValues ...map[string]float64) {
	var total, used, networkReceive, networkTransmit float64
	if totalValues != nil {
		total = totalValues["total"]
	}
	if len(usedValues) > 0 {
		used = usedValues[0]["used"]
	}
	if len(usedValues) > 1 {
		networkReceive = usedValues[1]["network_receive"]
		networkTransmit = usedValues[1]["network_transmit"]
	}

	available := total - used
	bandwidth := networkReceive + networkTransmit

	if metric == "cpu" {
		fmt.Fprintf(w, "%s\t%s\t-\t%.2f %s\t-\t-\t-\t-\n", service, metric, used, unit)
	} else if metric == "disk" {
		fmt.Fprintf(w, "%s\t%s\t%.4f %s\t%.2f %s\t%.2f %s\t-\t-\t-\n", service, metric, total, unit, used, unit, available, unit)
	} else {
		if total == 0 {
			fmt.Fprintf(w, "%s\t%s\t0.0000 %s\t0.0000 %s\t0.0000 %s\t-\t-\t-\n", service, metric, unit, unit, unit)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%.2f %s\t%.2f %s\t%.2f %s\t%.4f %s\t%.4f %s\t%.4f %s\n", service, metric, total, unit, used, unit, available, unit, networkReceive, unit, networkTransmit, unit, bandwidth, unit)
		}
	}
}

func printMemoryMetrics(w *tabwriter.Writer, service, metric string, totalValues map[string]float64, availableValues map[string]float64, unit string) {
	total := totalValues["total"]
	available := availableValues["used"]
	used := total - available

	fmt.Fprintf(w, "%s\t%s\t%.2f %s\t%.2f %s\t%.2f %s\t-\t-\t-\n", service, metric, total, unit, used, unit, available, unit)
}

func printNetworkMetrics(w *tabwriter.Writer, service, metric string, receiveValues map[string]float64, transmitValues map[string]float64, unit string) {
	receive := receiveValues["network_receive"]
	transmit := transmitValues["network_transmit"]
	bandwidth := receive + transmit

	fmt.Fprintf(w, "%s\t%s\t-\t-\t-\t%.4f %s\t%.4f %s\t%.4f %s\n", service, metric, receive, unit, transmit, unit, bandwidth, unit)
}
