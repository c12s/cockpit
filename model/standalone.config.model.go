package model

type StandaloneConfig struct {
	Organization string  `json:"organization" yaml:"organization"`
	Namespace string          `json:"namespace" yaml:"namespace"`
	Name         string  `json:"name" yaml:"name"`
	Version      string  `json:"version" yaml:"version"`
	CreatedAt    string  `json:"createdAt" yaml:"createdAt"`
	ParamSet     []Param `json:"paramSet" yaml:"paramSet"`
}

type StandaloneConfigsResponse struct {
	Configurations []StandaloneConfig `json:"configurations" yaml:"configurations"`
}

type SingleConfigDiffRequest struct {
	Reference SingleConfigReference `json:"reference" yaml:"reference"`
	Diff      SingleConfigReference `json:"diff" yaml:"diff"`
}

type StandaloneConfigDiffResponse struct {
	Diffs []SingleConfigDiff `json:"diffs" yaml:"diffs"`
}

type SingleConfigDiff struct {
	Type string            `json:"type" yaml:"type"`
	Diff map[string]string `json:"diff" yaml:"diff"`
}

type SingleConfigReference struct {
	Namespace    string `json:"namespace" yaml:"namespace"`
	Name         string `json:"name" yaml:"name"`
	Organization string `json:"organization" yaml:"organization"`
	Version      string `json:"version" yaml:"version"`
}

type StandaloneConfigDiffDetail struct {
	Key      string `json:"key" yaml:"key"`
	Value    string `json:"value,omitempty" yaml:"value,omitempty"`
	NewValue string `json:"new_value,omitempty" yaml:"new_value,omitempty"`
	OldValue string `json:"old_value,omitempty" yaml:"old_value,omitempty"`
}

type StandaloneConfigDiff struct {
	Type string                     `json:"type" yaml:"type"`
	Diff StandaloneConfigDiffDetail `json:"diff" yaml:"diff"`
}

type StandaloneConfigStandaloneConfig struct {
	Diffs []StandaloneConfigDiff `json:"diffs" yaml:"diffs"`
}
