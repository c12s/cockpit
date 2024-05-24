package model

type Param struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type ParamSet struct {
	Name     string  `json:"name" yaml:"name"`
	ParamSet []Param `json:"paramSet" yaml:"paramSet"`
}

type ConfigGroup struct {
	Organization string     `json:"organization" yaml:"organization"`
	Name         string     `json:"name" yaml:"name"`
	Version      string     `json:"version" yaml:"version"`
	CreatedAt    string     `json:"createdAt" yaml:"createdAt"`
	ParamSets    []ParamSet `json:"paramSets" yaml:"paramSets"`
}

type SingleConfigGroupRequest struct {
	Name         string `json:"name" yaml:"name"`
	Organization string `json:"organization" yaml:"organization"`
	Version      string `json:"version" yaml:"version"`
}

type ConfigGroupsResponse struct {
	Groups []ConfigGroup `json:"groups" yaml:"groups"`
}

type SingleConfigGroupResponse struct {
	Organization string `json:"organization" yaml:"organization"`
	Name         string `json:"name" yaml:"name"`
	Version      string `json:"version" yaml:"version"`
	CreatedAt    string `json:"createdAt" yaml:"createdAt"`
	ParamSets    []struct {
		Name     string `json:"name" yaml:"name"`
		ParamSet []struct {
			Key   string `json:"key" yaml:"key"`
			Value string `json:"value" yaml:"value"`
		} `json:"paramSet" yaml:"paramSet"`
	} `json:"paramSets" yaml:"paramSets"`
}
