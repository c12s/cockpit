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

type DeleteConfigGroupResponse struct {
	Organization string     `json:"organization" yaml:"organization"`
	Name         string     `json:"name" yaml:"name"`
	Version      string     `json:"version" yaml:"version"`
	CreatedAt    string     `json:"createdAt" yaml:"createdAt"`
	ParamSets    []ParamSet `json:"paramSets" yaml:"paramSets"`
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

type ConfigGroupReference struct {
	Name         string `json:"name" yaml:"name"`
	Organization string `json:"organization" yaml:"organization"`
	Version      string `json:"version" yaml:"version"`
}

type ConfigGroupDiffRequest struct {
	Reference ConfigGroupReference `json:"reference" yaml:"reference"`
	Diff      ConfigGroupReference `json:"diff" yaml:"diff"`
}

type ConfigGroupDiffResponse struct {
	Diffs map[string]interface{} `json:"diffs" yaml:"diffs"`
}
