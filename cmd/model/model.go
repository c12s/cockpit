package model

type MutateFile struct {
	Content Constellations `yaml:"constellations"`
}

type Constellations struct {
	Version  string     `yaml:"version"`
	Kind     string     `yaml:"kind"`
	Regions  []string   `yaml:"regions"`
	Clusters []string   `yaml:"clusters"`
	Labels   []string   `yaml:"labels"`
	Data     MutateData `yaml:"data"`
}

type MutateData struct {
	File []string `yaml:"file"`
	Env  []string `yaml:"env"`
}
