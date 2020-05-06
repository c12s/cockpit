package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type CContext struct {
	Context *Content `yaml:"context"`
}

type Content struct {
	Version   string `yaml:"version"`
	Address   string `yaml:"Address"`
	User      string `yaml:"user"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
}

func (c *CContext) Address() string {
	return c.Context.Address
}

func (c *CContext) Version() string {
	return c.Context.Version
}

func (c *CContext) User() string {
	return c.Context.User
}

func (c *CContext) Token() string {
	return c.Context.Token
}

func Marshall(c *CContext) (error, string) {
	d, err := yaml.Marshal(c)
	if err != nil {
		return err, ""
	}

	return nil, string(d)
}

func Context(address string) (error, *CContext) {
	data, err := ioutil.ReadFile(address)
	if err != nil {
		return err, nil
	}

	var c CContext
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err, nil
	}

	return nil, &c
}
