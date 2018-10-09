package model

type MutateFile struct {
	Content Constellations `yaml:"constellations"`
}

// Model for parsing yml file
type Constellations struct {
	Version   string                       `yaml:"version"`
	Kind      string                       `yaml:"kind"`
	Namespace string                       `yaml:"namespace"`
	Payload   map[string][]string          `yaml:"payload"`
	Strategy  map[string]string            `yaml:"strategy"`
	Selector  map[string]map[string]string `yaml:"selector"`
	Region    map[string]Region            `yaml:"region"`
}

type Region struct {
	Strategy map[string]string            `yaml:"strategy"`
	Payload  map[string][]string          `yaml:"payload"`
	Selector map[string]map[string]string `yaml:"selector"`
	Cluster  map[string]Cluster           `yaml:"cluster"`
}

type Cluster struct {
	Strategy map[string]string            `yaml:"strategy"`
	Payload  map[string][]string          `yaml:"payload"`
	Selector map[string]map[string]string `yaml:"selector"`
}
