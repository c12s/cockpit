package helper

import (
	"encoding/json"
	"errors"
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

func searchMutate(file *model.Constellations) []request.Region {
	regions := []request.Region{}
	s := constructStrategy(file.Strategy)
	if s == nil {
		fmt.Println("Error parsing strategy") //TODO: shuld return error or panic
	}

	p := constructPayload(file.Payload, file.Kind)
	if p == nil {
		fmt.Println("Error parsing payload") //TODO: shuld return error or panic
	}

	sl := constructSelector(file.Selector)
	if sl == nil {
		fmt.Println("Error parsing selector") //TODO: shuld return error or panic
	}

	clusters := []request.Cluster{}
	c := request.Cluster{
		ID:       "*",
		Payload:  p,
		Strategy: *s,
		Selector: *sl,
	}
	clusters = append(clusters, c)

	r := request.Region{
		ID:       "*",
		Clusters: clusters,
	}
	regions = append(regions, r)

	return regions
}

func detailMutate(file *model.Constellations) []request.Region {
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
	return regions
}

func convertFile(file *model.Constellations) (*request.MutateRequest, error) {
	err, ctx := GetContext()
	if err != nil {
		return nil, err
	}

	if ctx.Context.User == "" {
		return nil, errors.New("Please login to continue")
	}

	regions := []request.Region{}
	if len(file.Region) > 0 {
		regions = append(regions, detailMutate(file)...)
	} else {
		regions = append(regions, searchMutate(file)...)
	}

	var namespace = "default"
	if file.MTData.Namespace != "" {
		namespace = file.MTData.Namespace
	} else if ctx.Context.Namespace != "" {
		namespace = ctx.Context.Namespace
	}

	return &request.MutateRequest{
		Version: file.Version,
		Request: ctx.Context.User,
		Regions: regions,
		Kind:    file.Kind,
		MTData: request.Metadata{
			TaskName:     file.MTData.TaskName,
			Timestamp:    timestamp(),
			Namespace:    namespace,
			ForceNSQueue: file.MTData.ForceNSQueue,
			Queue:        file.MTData.Queue,
		},
	}, nil
}

func convertNFile(file *model.NConstellations) (*request.NMutateRequest, error) {
	labels := map[string]string{}
	for key, value := range file.Payload[LABELS] {
		labels[key] = value
	}

	metadata := request.Metadata{
		TaskName:     file.MTData.TaskName,
		Timestamp:    timestamp(),
		Namespace:    file.MTData.Namespace,
		ForceNSQueue: file.MTData.ForceNSQueue,
		Queue:        file.MTData.Queue,
	}

	err, ctx := GetContext()
	if err != nil {
		return nil, err
	}

	if ctx.Context.User == "" {
		return nil, errors.New("Please login to continue")
	}

	return &request.NMutateRequest{
		Version: file.Version,
		Request: ctx.Context.User,
		Name:    file.Payload[NAMESPACE][NAME],
		Labels:  labels,
		Kind:    NAMESPACES,
		MTData:  metadata,
	}, nil
}

func FileToJSON(file interface{}) (string, error) {
	switch v := file.(type) {
	case *model.Constellations:
		data, err1 := convertFile(v)
		if err1 != nil {
			return "", err1
		}

		dat, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return string(dat), nil
	case *model.NConstellations:
		data, err1 := convertNFile(v)
		if err1 != nil {
			return "", err1
		}

		dat, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return string(dat), nil
	}

	return "NOT VALID", nil
}
