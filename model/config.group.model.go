package model

type Query struct {
	LabelKey string `json:"labelKey" yaml:"labelKey"`
	ShouldBe string `json:"shouldBe" yaml:"shouldBe"`
	Value    string `json:"value" yaml:"value"`
}

type PlaceConfigGroupPlacementsRequest struct {
	Config   ConfigReference `json:"config" yaml:"config"`
	Strategy struct {
		Name       string  `json:"name" yaml:"name"`
		Query      []Query `json:"query" yaml:"query"`
		Percentage int     `json:"percentage" yaml:"percentage"`
	} `json:"strategy" yaml:"strategy"`
}

type ParamSet struct {
	Name     string  `json:"name" yaml:"name"`
	ParamSet []Param `json:"paramSet" yaml:"paramSet"`
}

type ConfigGroup struct {
	Organization string     `json:"organization" yaml:"organization"`
	Namespace    string     `json:"namespace" yaml:"namespace"`
	Name         string     `json:"name" yaml:"name"`
	Version      string     `json:"version" yaml:"version"`
	CreatedAt    string     `json:"createdAt" yaml:"createdAt"`
	ParamSets    []ParamSet `json:"paramSets" yaml:"paramSets"`
}

type ConfigGroupsResponse struct {
	Groups []ConfigGroup `json:"groups" yaml:"groups"`
}

type ConfigReference struct {
	Name         string `json:"name" yaml:"name"`
	Namespace    string `json:"namespace" yaml:"namespace"`
	Organization string `json:"organization" yaml:"organization"`
	Version      string `json:"version" yaml:"version"`
}

type ConfigGroupDiffRequest struct {
	Reference ConfigReference `json:"reference" yaml:"reference"`
	Diff      ConfigReference `json:"diff" yaml:"diff"`
}

type ConfigGroupPlacementsResponse struct {
	Tasks []Task `json:"tasks" yaml:"tasks"`
}

type Task struct {
	ID         string `json:"id" yaml:"id"`
	Node       string `json:"node" yaml:"node"`
	Status     string `json:"status" yaml:"status"`
	AcceptedAt string `json:"acceptedAt" yaml:"acceptedAt"`
	ResolvedAt string `json:"resolvedAt" yaml:"resolvedAt"`
}

type Param struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type ConfigGroupDiffDetail struct {
	Key      string `json:"key" yaml:"key"`
	Value    string `json:"value,omitempty" yaml:"value,omitempty"`
	NewValue string `json:"new_value,omitempty" yaml:"new_value,omitempty"`
	OldValue string `json:"old_value,omitempty" yaml:"old_value,omitempty"`
}

type ConfigGroupDiff struct {
	Type string                `json:"type" yaml:"type"`
	Diff ConfigGroupDiffDetail `json:"diff" yaml:"diff"`
}

type ConfigGroupDiffResponse struct {
	Diffs map[string]struct {
		Diffs []ConfigGroupDiff `json:"diffs" yaml:"diffs"`
	} `json:"diffs" yaml:"diffs"`
}
