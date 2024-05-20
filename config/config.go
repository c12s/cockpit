package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Services map[string]string                             `yaml:"services"`
	Gateway  Gateway                                       `yaml:"gateway"`
	Groups   map[string]map[string]map[string]MethodConfig `yaml:"groups"`
}

type Gateway struct {
	Route string `yaml:"route"`
	Port  string `yaml:"port"`
}

type MethodConfig struct {
	MethodRoute string `yaml:"method_route"`
	Type        string `yaml:"type"`
	Service     string `yaml:"service"`
}

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
