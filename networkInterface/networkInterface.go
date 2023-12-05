package networkinterface

import (
	"math/rand"
	"net"

	"github.com/michaelhenkel/infra-mgr/network"
)

type NetworkInterfaceRoute struct {
	Destination Destination
	NextHop     []NextHopInterface
}

type NextHopInterface struct {
	IpAddress  net.IP
	MacAddress string
}

type Destination struct {
	IpAddress string
	Mask      string
}

type Config struct {
	Name       string
	Network    string
	MacAddress *string
	IpAddress  *net.IP
}

func New(name string, network *network.Config) *Config {
	macAddress := GenerateMacAddress().String()
	ipAddress := network.AllocateIpAddress()
	return &Config{
		Name:       name,
		Network:    network.Name,
		MacAddress: &macAddress,
		IpAddress:  &ipAddress,
	}
}

// func randomMacAddress generates a random unicast MAC address.
func GenerateMacAddress() net.HardwareAddr {
	mac := make(net.HardwareAddr, 6)
	mac[0] = 0x02
	mac[1] = 0x00
	mac[2] = 0x00
	mac[3] = byte(rand.Intn(256))
	mac[4] = byte(rand.Intn(256))
	mac[5] = byte(rand.Intn(256))
	return mac
}
