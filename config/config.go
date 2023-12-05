package config

import (
	"github.com/michaelhenkel/infra-mgr/instance"
	"github.com/michaelhenkel/infra-mgr/network"
)

type Config struct {
	Name      string
	Instances []*instance.Config
	Networks  []*network.Config
}

func New(name string) *Config {
	return &Config{
		Name:      name,
		Instances: []*instance.Config{},
		Networks:  []*network.Config{},
	}
}

func (c *Config) AddInstance(i *instance.Config) *Config {
	c.Instances = append(c.Instances, i)
	return c
}

func (c *Config) AddNetwork(n *network.Config) *Config {
	c.Networks = append(c.Networks, n)
	return c
}
