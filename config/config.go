package config

import (
	"github.com/michaelhenkel/infra-mgr/infrastructure"
	"github.com/michaelhenkel/infra-mgr/instance"
	"github.com/michaelhenkel/infra-mgr/network"
	networkinterface "github.com/michaelhenkel/infra-mgr/networkInterface"
	"github.com/michaelhenkel/infra-mgr/route"
)

type Config struct {
	Name       string
	Instances  []*instance.Config
	Networks   []*network.Config
	Interfaces []*networkinterface.Config
	Routes     []*route.Config
}

func New(name string) *Config {
	return &Config{
		Name:       name,
		Instances:  []*instance.Config{},
		Networks:   []*network.Config{},
		Interfaces: []*networkinterface.Config{},
		Routes:     []*route.Config{},
	}
}

func (c *Config) Build() *infrastructure.Infrastructure {
	infrastructure := infrastructure.New(c.Name)

	for _, instanceConfig := range c.Instances {
		instance := instanceConfig.Build()
		infrastructure.Add(instance)
	}

	for _, networkConfig := range c.Networks {
		network := networkConfig.New()
		infrastructure.Add(network)
	}

	for _, interfaceConfig := range c.Interfaces {
		networkInterface := interfaceConfig.Build()
		if inst, ok := infrastructure.Instances[interfaceConfig.Instance]; !ok {
			panic("instance not found")
		} else {
			inst.Add(networkInterface)
		}

		if network, ok := infrastructure.Networks[interfaceConfig.Network]; !ok {
			panic("network not found")
		} else {
			network.Add(networkInterface)
		}
	}

	for _, routeConfig := range c.Routes {
		route := routeConfig.Build(routeConfig, infrastructure)
		if inst, ok := infrastructure.Instances[routeConfig.Source.Instance]; !ok {
			panic("instance not found")
		} else {
			if networkInterface, ok := inst.NetworkInterfaces[routeConfig.Source.NetworkInterface]; !ok {
				panic("network interface not found")
			} else {
				networkInterface.Add(route)
			}
		}
	}
	return infrastructure
}
