package model

type ConfigGroup struct {
	Name    string            `yaml:"Name"`
	OrgId   string            `yaml:"OrgId"`
	Version int32             `yaml:"Version"`
	Configs map[string]string `yaml:"Configs,omitempty"`
}

type Query struct {
	Key      string `yaml:"Key"`
	ShouldBe string `yaml:"ShouldBe"`
	Value    string `yaml:"Value"`
}
