package instance

import (
	networkinterface "github.com/michaelhenkel/infra-mgr/networkInterface"
)

type Config struct {
	Name              string
	Vcpu              int
	Memory            string
	Image             string
	Routes            []NetworkInterfaceRoute
	NetworkInterfaces []*networkinterface.Config
}

func New(name string, vcpu int, memory string, image string) *Config {
	return &Config{
		Name:              name,
		Vcpu:              vcpu,
		Memory:            memory,
		Image:             image,
		Routes:            []NetworkInterfaceRoute{},
		NetworkInterfaces: []*networkinterface.Config{},
	}
}

type NetworkInterfaceRoute struct {
	Destination *networkinterface.Config
	NextHops    []*networkinterface.Config
}

func (c *Config) AddRoute(networkInterface *networkinterface.Config, nextNetworkInterfaces []*networkinterface.Config) {
	networkInterfaceRoute := NetworkInterfaceRoute{
		Destination: networkInterface,
		NextHops:    nextNetworkInterfaces,
	}
	c.Routes = append(c.Routes, networkInterfaceRoute)
}

func (c *Config) AddNetworkInterface(networkInterface *networkinterface.Config) {
	c.NetworkInterfaces = append(c.NetworkInterfaces, networkInterface)
}

func (c *Config) Build() {

}
