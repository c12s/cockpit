package model

import (
	"github.com/c12s/kuiper/pkg/api"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
)

func ConfigsFromDomain(configs map[string]string) []*api.Config {
	resp := make([]*api.Config, 0)
	for key, value := range configs {
		resp = append(resp, &api.Config{
			Key:   key,
			Value: value,
		})
	}
	return resp
}

func ConfigsToDomain(configs []*api.Config) map[string]string {
	resp := make(map[string]string)
	for _, config := range configs {
		resp[config.Key] = config.Value
	}
	return resp
}

func QueriesFromDomain(queries []Query) []*magnetarapi.Query {
	resp := make([]*magnetarapi.Query, len(queries))
	for i, query := range queries {
		resp[i] = &magnetarapi.Query{
			LabelKey: query.Key,
			ShouldBe: query.ShouldBe,
			Value:    query.Value,
		}
	}
	return resp
}

func PermKindFromDomain(kind string) oortapi.Permission_PermissionKind {
	switch kind {
	case "allow":
		return oortapi.Permission_ALLOW
	case "deny":
		return oortapi.Permission_DENY
	default:
		return oortapi.Permission_DENY
	}
}
