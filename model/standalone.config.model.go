package model

type StandaloneConfig struct {
	Organization string  `json:"organization" yaml:"organization"`
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

type SingleConfigDiffResponse struct {
	Diffs []SingleConfigDiff `json:"diffs" yaml:"diffs"`
}

type SingleConfigDiff struct {
	Type string            `json:"type" yaml:"type"`
	Diff map[string]string `json:"diff" yaml:"diff"`
}

type SingleConfigReference struct {
	Name         string `json:"name" yaml:"name"`
	Organization string `json:"organization" yaml:"organization"`
	Version      string `json:"version" yaml:"version"`
}
