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
