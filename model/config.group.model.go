package model

type AppConfigRequest struct {
	Name         string `json:"name" yaml:"name"`
	Organization string `json:"organization" yaml:"organization"`
	Version      string `json:"version" yaml:"version"`
}

type AppConfigResponse struct {
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
