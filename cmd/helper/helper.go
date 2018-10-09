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

// func inside(key string) bool {
// 	for _, item := range allowed_payloads {
// 		if key == item {
// 			return true
// 		}
// 	}

// 	return false
// }

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

			p := extractPayload(file.Payload, region.Payload, cluster.Payload)
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

	return &request.MutateRequest{
		Request:   "user_name_email_or_something_else",
		Timestamp: timestamp(),
		Regions:   regions,
		Namespace: file.Namespace,
	}
}

func FileToJSON(file *model.Constellations) (string, error) {
	data := convertFile(file)
	dat, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
