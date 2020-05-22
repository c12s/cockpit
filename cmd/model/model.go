package model

import (
	"gopkg.in/yaml.v2"
)

type MutateFile struct {
	Content Constellations `yaml:"constellations"`
}

// Model for parsing yml file
type Constellations struct {
	Version  string                   `yaml:"version"`
	Kind     string                   `yaml:"kind"`
	MTData   Metadata                 `yaml:"metadata"`
	Payload  map[string]yaml.MapSlice `yaml:"payload"`
	Strategy Strategy                 `yaml:"strategy"`
	Selector map[string]yaml.MapSlice `yaml:"selector"`
	Region   map[string]Region        `yaml:"region"`
}

type Metadata struct {
	TaskName     string `yaml:"taskName"`
	Namespace    string `yaml:"namespace"`
	Queue        string `yaml:"queue"`
	ForceNSQueue bool   `yaml:"forceNamespaceQueue"`
}

type Region struct {
	Strategy Strategy                 `yaml:"strategy"`
	Payload  map[string]yaml.MapSlice `yaml:"payload"` //map[string]string
	Selector map[string]yaml.MapSlice `yaml:"selector"`
	Cluster  map[string]Cluster       `yaml:"cluster"`
}

type Cluster struct {
	Strategy Strategy                 `yaml:"strategy"`
	Payload  map[string]yaml.MapSlice `yaml:"payload"`
	Selector map[string]yaml.MapSlice `yaml:"selector"`
}

type NConstellations struct {
	Version string                       `yaml:"version"`
	Kind    string                       `yaml:"kind"`
	MTData  Metadata                     `yaml:"metadata"`
	Payload map[string]map[string]string `yaml:"payload"`
}

type NMutateFile struct {
	Content NConstellations `yaml:"constellations"`
}

type Strategy struct {
	Type     string            `yaml:"type"`
	Update   string            `yaml:"update"`
	Interval string            `yaml:"interval"`
	Retry    map[string]string `yaml:"retry"`
}

type RolesFile struct {
	Content Roles `yaml:"constellations"`
}

type Roles struct {
	Version string   `yaml:"version"`
	Kind    string   `yaml:"kind"`
	MTData  Metadata `yaml:"metadata"`
	Payload Rules    `yaml:"rules"`
}

type Rules struct {
	User      string   `yaml:"user"`
	Resources []string `yaml:"resources"`
	Verbs     []string `yaml:"verbs"`
}
