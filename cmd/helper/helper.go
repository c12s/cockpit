package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/c12s/cockpit/cmd/model/request"
	"gopkg.in/yaml.v2"
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
	s := constructStrategy(&file.Strategy)
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
		if len(region.Cluster) > 0 {
			for clusterID, cluster := range region.Cluster {
				s := extractStrategy(&file.Strategy, &region.Strategy, &cluster.Strategy)
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

				fmt.Println(2)
				c := request.Cluster{
					ID:       clusterID,
					Payload:  p,
					Strategy: *s,
					Selector: *sl,
				}
				clusters = append(clusters, c)
			}
		} else {
			emptyPS := map[string]yaml.MapSlice{}
			// emptyS := map[string]string{}
			s := extractStrategy(&file.Strategy, &region.Strategy, nil)
			if s == nil {
				fmt.Println("Error parsing strategy") //TODO: shuld return error or panic
			}

			p := extractPayload(file.Payload, region.Payload, emptyPS, file.Kind)
			if p == nil {
				fmt.Println("Error parsing payload") //TODO: shuld return error or panic
			}

			sl := extractSelector(file.Selector, region.Selector, emptyPS)
			if sl == nil {
				fmt.Println("Error parsing selector") //TODO: shuld return error or panic
			}

			c := request.Cluster{
				ID:       "*",
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
	err, ctx := GetContext()
	if err != nil {
		return nil, err
	}

	if ctx.Context.User == "" {
		return nil, errors.New("Please login to continue")
	}

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

	return &request.NMutateRequest{
		Version: file.Version,
		Request: ctx.Context.User,
		Name:    file.Payload[NAMESPACE][NAME],
		Labels:  labels,
		Kind:    NAMESPACES,
		MTData:  metadata,
	}, nil
}

func convertRoleFile(file *model.Roles) (*request.RMutateRequest, error) {
	err, ctx := GetContext()
	if err != nil {
		return nil, err
	}

	if ctx.Context.User == "" {
		return nil, errors.New("Please login to continue")
	}

	return &request.RMutateRequest{
		Version: file.Version,
		Request: ctx.Context.User, // user who sent request
		Kind:    ROLES,
		MTData: request.Metadata{
			TaskName:     file.MTData.TaskName,
			Timestamp:    timestamp(),
			Namespace:    file.MTData.Namespace,
			ForceNSQueue: file.MTData.ForceNSQueue,
			Queue:        file.MTData.Queue,
		},
		Payload: request.Rules{
			User:      file.Payload.User, // user to change rules for
			Resources: file.Payload.Resources,
			Verbs:     file.Payload.Verbs,
		},
	}, nil
}

func convertUFile(file *model.NConstellations) (*request.UMutateRequest, error) {
	err, ctx := GetContext()
	if err != nil {
		return nil, err
	}

	if ctx.Context.User == "" {
		return nil, errors.New("Please login to continue")
	}

	metadata := request.Metadata{
		TaskName:     file.MTData.TaskName,
		Timestamp:    timestamp(),
		Namespace:    file.MTData.Namespace,
		ForceNSQueue: file.MTData.ForceNSQueue,
		Queue:        file.MTData.Queue,
	}

	return &request.UMutateRequest{
		Version: file.Version,
		Request: ctx.Context.User, // user who sent request
		Info:    file.Payload["info"],
		Labels:  file.Payload[LABELS],
		Kind:    USERS,
		MTData:  metadata,
	}, nil
}

func convertTopologyFile(file *model.Topology) (*request.TMutateRequest, error) {
	err, ctx := GetContext()
	if err != nil {
		return nil, err
	}

	if ctx.Context.User == "" {
		return nil, errors.New("Please login to continue")
	}

	metadata := request.Metadata{
		TaskName:     file.MTData.TaskName,
		Timestamp:    timestamp(),
		Namespace:    file.MTData.Namespace,
		ForceNSQueue: file.MTData.ForceNSQueue,
		Queue:        file.MTData.Queue,
	}

	regions := []request.TRegion{}
	for regionid, clusters := range file.Payload.Topology {
		tclusters := []request.TCluster{}
		for clusterid, nodes := range clusters {
			tnodes := []request.TNode{}
			retention := ""
			for nodeid, value := range nodes {
				if nodeid != "retention" {
					labels := map[string]string{}
					for _, v := range value["selector"].(map[interface{}]interface{}) {
						for vk, vv := range v.(map[interface{}]interface{}) {
							labels[vk.(string)] = vv.(string)
						}
					}

					tnodes = append(tnodes, request.TNode{
						ID:     nodeid,
						Name:   value["name"].(string),
						Labels: labels,
					})
				} else {
					retention = value["period"].(string)
				}
			}
			tclusters = append(tclusters, request.TCluster{
				ID:        clusterid,
				Retention: retention,
				Nodes:     tnodes,
			})
		}
		regions = append(regions, request.TRegion{
			ID:       regionid,
			Clusters: tclusters,
		})
	}
	topology := request.Topology{
		Name:    file.Payload.Name,
		Labels:  file.Payload.Selector["labels"],
		Regions: regions,
	}

	return &request.TMutateRequest{
		Version: file.Version,
		Request: ctx.Context.User, // user who sent request
		Kind:    TOPOLOGY,
		MTData:  metadata,
		Payload: topology,
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
		if v.Kind == "Namespaces" {
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

		data, err1 := convertUFile(v)
		if err1 != nil {
			return "", err1
		}

		dat, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(dat), nil
	case *model.Roles:
		data, err1 := convertRoleFile(v)
		if err1 != nil {
			return "", err1
		}

		dat, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return string(dat), nil
	case *model.Topology:
		data, err1 := convertTopologyFile(v)
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
