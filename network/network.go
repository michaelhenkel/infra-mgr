package network

import (
	"net"

	networkinterface "github.com/michaelhenkel/infra-mgr/networkInterface"
	"github.com/michaelhenkel/infra-mgr/object"
)

type Network struct {
	Config            *Config
	Subnet            *net.IPNet
	Gateway           net.IP
	NetworkInterfaces map[string]*networkinterface.NetworkInterface
	AllocatedIps      map[uint32]net.IP
}

type Config struct {
	Name   string
	Subnet string
}

func New(name string, subnet string) *Config {
	return &Config{
		Name:   name,
		Subnet: subnet,
	}
}

func (config *Config) Add(object.Object) {
}

func (config *Config) GetName() string {
	return config.Name
}

func (config *Config) ObjectType() object.ObjectType {
	return object.Network
}

func (config *Config) New() *Network {

	_, ipnet, err := net.ParseCIDR(config.Subnet)
	if err != nil {
		panic(err)
	}
	ipnetUint32 := IpToUint32(ipnet.IP.To4())
	gatewayUint32 := ipnetUint32 + 1
	gateway := Uint32ToIp(gatewayUint32)
	return &Network{
		Subnet:            ipnet,
		Gateway:           gateway,
		Config:            config,
		NetworkInterfaces: make(map[string]*networkinterface.NetworkInterface),
		AllocatedIps:      map[uint32]net.IP{gatewayUint32: gateway},
	}
}

func (n *Network) Add(o object.Object) {
	switch o.ObjectType() {
	case object.NetworkInterface:
		intf := o.(*networkinterface.NetworkInterface)
		ipUint32 := IpToUint32(n.Gateway.To4())

		for k := range n.AllocatedIps {
			if k == ipUint32 {
				ipUint32++
				if !n.Subnet.Contains(Uint32ToIp(ipUint32)) {
					panic("No more IP addresses available")
				}
			} else {
				break
			}
		}

		n.AllocatedIps[ipUint32] = Uint32ToIp(ipUint32)
		intf.IpAddress = Uint32ToIp(ipUint32)
		n.NetworkInterfaces[o.GetName()] = intf
	}
}

func (n *Network) ObjectType() object.ObjectType {
	return object.Network
}

func (n *Network) GetName() string {
	return n.Config.Name
}

func (c *Config) GetDestination() *net.IPNet {
	return nil
}

func IpToUint32(ip net.IP) uint32 {
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

func Uint32ToIp(ipUint32 uint32) net.IP {
	return net.IPv4(byte(ipUint32>>24), byte(ipUint32>>16), byte(ipUint32>>8), byte(ipUint32))
}
