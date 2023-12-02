package networkinterface

import (
	"net"

	"github.com/michaelhenkel/infra-mgr/object"
	"github.com/michaelhenkel/infra-mgr/route"
)

type NetworkInterface struct {
	Config    *Config
	IpAddress net.IP
	Routes    map[*net.IPNet]string
}

type Config struct {
	Name     string
	Network  string
	Instance string
}

func New(name string, network object.Object, instance object.Object) *Config {
	return &Config{
		Name:     name,
		Network:  network.GetName(),
		Instance: instance.GetName(),
	}
}

func (config *Config) Build() *NetworkInterface {
	return &NetworkInterface{
		Config: config,
		Routes: make(map[*net.IPNet]string),
	}
}

func (n *NetworkInterface) Add(o object.Object) {
	switch o.ObjectType() {
	case object.Route:
		n.Routesroute := o.(*route.Route)

	}
}

func (n *NetworkInterface) ObjectType() object.ObjectType {
	return object.NetworkInterface
}

func (n *NetworkInterface) GetName() string {
	return n.Config.Name
}

func (c *Config) GetDestination() *net.IPNet {
	return nil
}
