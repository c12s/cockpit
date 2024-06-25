package utils

import (
	"github.com/c12s/cockpit/model"
	"strings"
)

func GetServiceCpuUsage(metric model.MetricData) float64 {
	for _, value := range metric.Values {
		if strings.Contains(metric.Metric["__name__"], "cpu") {
			return StringToFloat(value[1].(string))
		}
	}
	return 0.0
}

func GetServiceMemoryUsage(metric model.MetricData) float64 {
	for _, value := range metric.Values {
		if strings.Contains(metric.Metric["__name__"], "ram_usage") {
			return StringToFloat(value[1].(string))
		}
	}
	return 0.0
}
