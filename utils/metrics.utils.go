package utils

import (
	"strings"

	"github.com/c12s/cockpit/model"
)

func GetServiceCpuUsage(metric model.MetricData) float64 {
	if strings.Contains(metric.Metric["__name__"], "cpu") {
		return StringToFloat(metric.Value[1].(string))
	}
	return 0.0
}

func GetServiceMemoryUsage(metric model.MetricData) float64 {
	if strings.Contains(metric.Metric["__name__"], "ram_usage") {
		return StringToFloat(metric.Value[1].(string))
	}
	return 0.0
}

func GetServiceDiskUsage(metric model.MetricData) float64 {
	if strings.Contains(metric.Metric["__name__"], "disk_usage") {
		return StringToFloat(metric.Value[1].(string))
	}
	return 0.0
}

func GetServiceNetworkReceive(metric model.MetricData) float64 {
	if strings.Contains(metric.Metric["__name__"], "network_receive") {
		return StringToFloat(metric.Value[1].(string))
	}
	return 0.0
}

func GetServiceNetworkTransmit(metric model.MetricData) float64 {
	if strings.Contains(metric.Metric["__name__"], "network_transmit") {
		return StringToFloat(metric.Value[1].(string))
	}
	return 0.0
}

func GetServiceBandwidth(metric model.MetricData) float64 {
	return GetServiceNetworkReceive(metric) + GetServiceNetworkTransmit(metric)
}

func GetNodeBandwidth(metricsMap map[string]map[string]float64) float64 {
	return metricsMap["custom_node_network_receive_mb"]["network_receive"] + metricsMap["custom_node_network_transmit_mb"]["network_transmit"]
}
