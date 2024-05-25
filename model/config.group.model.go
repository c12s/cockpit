package model

type PlaceConfigGroupPlacementsRequest struct {
	Config    ConfigGroupReference `json:"config" yaml:"config"`
	Namespace string               `json:"namespace" yaml:"namespace"`
	Strategy  struct {
		Name  string `json:"name" yaml:"name"`
		Query []struct {
			LabelKey string `json:"labelKey" yaml:"labelKey"`
			ShouldBe string `json:"shouldBe" yaml:"shouldBe"`
			Value    string `json:"value" yaml:"value"`
		} `json:"query" yaml:"query"`
	} `json:"strategy" yaml:"strategy"`
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

type ConfigGroupsResponse struct {
	Groups []ConfigGroup `json:"groups" yaml:"groups"`
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

type ConfigGroupPlacementsResponse struct {
	Tasks []Task `json:"tasks" yaml:"tasks"`
}

type Task struct {
	ID         string `json:"id" yaml:"id"`
	Node       string `json:"node" yaml:"node"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	Status     string `json:"status" yaml:"status"`
	AcceptedAt string `json:"acceptedAt" yaml:"acceptedAt"`
	ResolvedAt string `json:"resolvedAt" yaml:"resolvedAt"`
}

type Param struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}
