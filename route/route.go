package route

import (
	"net"

	"github.com/michaelhenkel/infra-mgr/infrastructure"
	"github.com/michaelhenkel/infra-mgr/network"
	networkinterface "github.com/michaelhenkel/infra-mgr/networkInterface"
	"github.com/michaelhenkel/infra-mgr/object"
)

type Route struct {
	Config      *Config
	Destination net.IPNet
	NextHops    []NextHop
}

type NextHop struct {
	IpAddress        net.IP
	Instance         string
	NetworkInterface string
}

type Config struct {
	Name        string
	Source      Source
	Destination interface{}
	NextHops    []NextHop
}

type Source struct {
	Instance         string
	NetworkInterface string
}

func New(name string, source Source, destination interface{}, nexthops []NextHop) *Config {
	return &Config{
		Name:        name,
		Source:      source,
		Destination: destination,
		NextHops:    nexthops,
	}
}

func (c *Config) Build(config *Config, infrastructure *infrastructure.Infrastructure) *Route {
	nextHops := []NextHop{}
	for _, nexthop := range config.NextHops {
		inst, ok := infrastructure.Instances[nexthop.Instance]
		if !ok {
			panic("instance not found")
		}
		networkInterface, ok := inst.NetworkInterfaces[nexthop.NetworkInterface]
		if !ok {
			panic("network interface not found")
		}
		nextHops = append(nextHops, NextHop{
			IpAddress:        networkInterface.IpAddress,
			NetworkInterface: nexthop.NetworkInterface,
		})
	}
	switch v := config.Destination.(type) {
	case *network.Config:
		dstNetwork, ok := infrastructure.Networks[v.Name]
		if !ok {
			panic("network not found")
		}
		return &Route{
			Config:      config,
			Destination: *dstNetwork.Subnet,
			NextHops:    nextHops,
		}
	case *networkinterface.Config:
		for _, inst := range infrastructure.Instances {
			intf, ok := inst.NetworkInterfaces[v.Name]
			if !ok {
				continue
			} else {
				return &Route{
					Config:      config,
					Destination: net.IPNet{IP: intf.IpAddress, Mask: net.CIDRMask(32, 32)},
					NextHops:    nextHops,
				}
			}
		}
	}
	return &Route{
		Config: config,
	}
}

func (n *Route) Add(o object.Object) {

}

func (n *Route) ObjectType() object.ObjectType {
	return object.Route
}

func (n *Route) GetName() string {
	return n.Config.Name
}
