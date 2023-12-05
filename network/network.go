package network

import (
	"net"
)

type Config struct {
	Name         string
	Subnet       string
	Gateway      net.IP
	AllocatedIps map[uint32]net.IP
}

func New(name string, subnet string) *Config {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		panic(err)
	}
	ipnetUint32 := IpToUint32(ipnet.IP.To4())
	gatewayUint32 := ipnetUint32 + 1
	gateway := Uint32ToIp(gatewayUint32)
	return &Config{
		Name:         name,
		Subnet:       subnet,
		Gateway:      gateway,
		AllocatedIps: make(map[uint32]net.IP),
	}
}

func (c *Config) AllocateIpAddress() net.IP {

	ipUint32 := IpToUint32(c.Gateway.To4())
	_, ipnet, err := net.ParseCIDR(c.Subnet)
	if err != nil {
		panic(err)
	}

	for k := range c.AllocatedIps {
		if k == ipUint32 {
			ipUint32++
			if !ipnet.Contains(Uint32ToIp(ipUint32)) {
				panic("No more IP addresses available")
			}
		} else {
			break
		}
	}

	c.AllocatedIps[ipUint32] = Uint32ToIp(ipUint32)
	return Uint32ToIp(ipUint32)
}

func IpToUint32(ip net.IP) uint32 {
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

func Uint32ToIp(ipUint32 uint32) net.IP {
	return net.IPv4(byte(ipUint32>>24), byte(ipUint32>>16), byte(ipUint32>>8), byte(ipUint32))
}
