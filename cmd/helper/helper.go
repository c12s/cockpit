package helper

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/c12s/cockpit/cmd/model/request"
	"path/filepath"
	"strings"
	"time"
)

func timestamp() int64 {
	return time.Now().UnixNano()
}

func constructFileKey(path string) string {
	name := filepath.Base(path)
	return strings.Join([]string{"file", name}, "://")
}

func convertFile(file *model.Constellations) *request.MutateRequest {
	regions := []request.Region{}
	for regionID, region := range file.Region {
		clusters := []request.Cluster{}
		for clusterID, cluster := range region.Cluster {
			s := extractStrategy(file.Strategy, region.Strategy, cluster.Strategy)
			if s == nil {
				fmt.Println("Error parsing strategy") //TODO: shuld return error or panic
			}

			p := extractPayload(file.Payload, region.Payload, cluster.Payload, file.Kind)
			if p == nil {
				fmt.Println("Error parsing payload") //TODO: shuld return error or panic
			}

			sl := extractSelector(file.Selector, region.Selector, cluster.Selector)
			if sl == nil {
				fmt.Println("Error parsing selector") //TODO: shuld return error or panic
			}

			c := request.Cluster{
				ID:       clusterID,
				Payload:  p,
				Strategy: *s,
				Selector: *sl,
			}
			clusters = append(clusters, c)
		}

		r := request.Region{
			ID:       regionID,
			Clusters: clusters,
		}
		regions = append(regions, r)
	}

	var namespace = "default"
	if file.MTData.Namespace != "" {
		namespace = file.MTData.Namespace
	}

	metadata := request.Metadata{
		Version:      file.Version,
		TaskName:     file.MTData.TaskName,
		Timestamp:    timestamp(),
		Namespace:    namespace,
		ForceNSQueue: file.MTData.ForceNSQueue,
	}

	return &request.MutateRequest{
		Request: "user_name_email_or_something_else",
		Regions: regions,
		Kind:    file.Kind,
		MTData:  metadata,
	}
}

func convertNFile(file *model.NConstellations) *request.NMutateRequest {
	labels := map[string]string{}
	for key, value := range file.Payload[LABELS] {
		labels[key] = value
	}

	metadata := request.Metadata{
		Version:      file.Version,
		TaskName:     file.MTData.TaskName,
		Timestamp:    timestamp(),
		Namespace:    file.MTData.Namespace,
		ForceNSQueue: file.MTData.ForceNSQueue,
	}

	return &request.NMutateRequest{
		Request: "user_name_email_or_something_else",
		Labels:  labels,
		Kind:    NAMESPACES,
		MTData:  metadata,
	}
}

func FileToJSON(file interface{}) (string, error) {
	switch v := file.(type) {
	case *model.Constellations:
		data := convertFile(v)
		dat, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return string(dat), nil
	case *model.NConstellations:
		data := convertNFile(v)
		dat, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return string(dat), nil
	}

	return "NOT VALID", nil
}
