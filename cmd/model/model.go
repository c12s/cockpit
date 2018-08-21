package model

type MutateFile struct {
	Content Constellations `yaml:"constellations"`
}

type Constellations struct {
	Version  string              `yaml:"version"`
	Kind     string              `yaml:"kind"`
	Payload  map[string][]string `yaml:"payload"`
	Strategy map[string]string   `yaml:"strategy"`
	Selector Selection           `yaml:"selector"`
	Region   map[string]Region   `yaml:"region"`
}

type Region struct {
	Cluster map[string]Cluster `yaml:"cluster"`
}

type Cluster struct {
	Strategy map[string]string   `yaml:"strategy"`
	Payload  map[string][]string `yaml:"payload"`
	Selector Selection           `yaml:"selector"`
}

type Selection struct {
	Labels  map[string]string `yaml:"labels"`
	Compare map[string]string `yaml:"compare"`
}
