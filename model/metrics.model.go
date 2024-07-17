package model

import (
	"strings"
)

type MetricResponse struct {
	Status int          `json:"status"`
	Data   []MetricData `json:"data"`
}

type MetricData struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

func (m *MetricResponse) FilterNodeMetrics() []MetricData {
	var nodeMetrics []MetricData
	for _, d := range m.Data {
		if name, ok := d.Metric["__name__"]; ok && strings.HasPrefix(name, "custom_node_") {
			nodeMetrics = append(nodeMetrics, d)
		}
	}
	return nodeMetrics
}

func (m *MetricResponse) FilterServiceMetrics() []MetricData {
	var serviceMetrics []MetricData
	for _, d := range m.Data {
		if name, ok := d.Metric["__name__"]; ok && strings.HasPrefix(name, "custom_service_") {
			serviceMetrics = append(serviceMetrics, d)
		}
	}
	return serviceMetrics
}
