package model

import (
	"gopkg.in/yaml.v2"
)

type CContext struct {
	Context *Content `yaml:"context"`
}

type Content struct {
	Version   string `yaml:"version"`
	Address   string `yaml:"Address"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace"`
}

func Marshall(c *CContext) (error, string) {
	d, err := yaml.Marshal(c)
	if err != nil {
		return err, ""
	}

	return nil, string(d)
}
