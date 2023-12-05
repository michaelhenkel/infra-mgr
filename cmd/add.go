package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/michaelhenkel/infra-mgr/config"
	"github.com/michaelhenkel/infra-mgr/instance"
	"github.com/michaelhenkel/infra-mgr/lxdmgr"
	"github.com/michaelhenkel/infra-mgr/network"
	networkinterface "github.com/michaelhenkel/infra-mgr/networkInterface"
)

func init() {

}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long:  `add`,
	Run:   add,
}

func add(cmd *cobra.Command, args []string) {
	println("cfgFile: ", cfgFile)
	if cfgFile != "" {
		data, err := os.ReadFile(cfgFile)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		var config config.Config

		if err = yaml.Unmarshal(data, &config); err != nil {
			log.Fatalf("error: %v", err)
		}

		for _, netw := range config.Networks {
			if netw.AllocatedIps == nil {
				netw.AllocatedIps = make(map[uint32]net.IP)
			}
		}

		interfaceMap := make(map[string]*networkinterface.Config)
		for _, instance := range config.Instances {
			for _, intf := range instance.NetworkInterfaces {
				if intf.MacAddress == nil || intf.IpAddress == nil {
					macAddress := networkinterface.GenerateMacAddress().String()
					intf.MacAddress = &macAddress
				}
				if intf.IpAddress == nil {
					for _, network := range config.Networks {
						if network.Name == intf.Network {
							ipAddress := network.AllocateIpAddress()
							intf.IpAddress = &ipAddress
						}
					}
				}
				interfaceMap[intf.Name] = intf
			}
		}

		for _, instance := range config.Instances {
			for _, route := range instance.Routes {
				if intf, ok := interfaceMap[route.Destination.Name]; ok {
					route.Destination.IpAddress = intf.IpAddress
					route.Destination.MacAddress = intf.MacAddress
				} else {
					panic("interface not found")
				}
				for _, nextHop := range route.NextHops {
					if intf, ok := interfaceMap[nextHop.Name]; ok {
						nextHop.Name = intf.Name
						nextHop.IpAddress = intf.IpAddress
						nextHop.MacAddress = intf.MacAddress
					} else {
						panic("interface not found")
					}
				}
			}
		}

		out, err := yaml.Marshal(config)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(out))

		lxdm, err := lxdmgr.New()
		if err != nil {
			panic(err)
		}
		if err := lxdm.CreateInstance("infra1", "host1", "image1"); err != nil {
			panic(err)
		}

	} else {
		sim()
	}
}

func sim() {
	config1 := config.New("infra1")
	host1 := instance.New("host1", 1, "1G", "image1")
	host2 := instance.New("host2", 1, "1G", "image1")
	router1 := instance.New("router1", 1, "1G", "image1")
	router2 := instance.New("router2", 1, "1G", "image1")
	fabric1 := network.New("fabric1", "10.0.0.0/24")
	fabric2 := network.New("fabric2", "10.0.1.0/24")
	access1 := network.New("access1", "10.0.2.0/24")
	access2 := network.New("access2", "10.0.3.0/24")

	router1NetworkInterface1 := networkinterface.New("r1_eth1", fabric1)
	router1NetworkInterface2 := networkinterface.New("r1_eth2", fabric2)
	router1AccessInterface1 := networkinterface.New("r1_eth3", access1)
	router2NetworkInterface1 := networkinterface.New("r2_eth1", fabric1)
	router2NetworkInterface2 := networkinterface.New("r2_eth2", fabric2)
	router2AccessInterface1 := networkinterface.New("r2_eth3", access2)
	host1NetworkInterface1 := networkinterface.New("h1_eth1", access1)
	host2NetworkInterface1 := networkinterface.New("h2_eth1", access2)

	config1.
		AddInstance(host1).
		AddInstance(host2).
		AddInstance(router1).
		AddInstance(router2).
		AddNetwork(fabric1).
		AddNetwork(fabric2).
		AddNetwork(access1).
		AddNetwork(access2)

	router1.AddNetworkInterface(router1NetworkInterface1)
	router1.AddNetworkInterface(router1NetworkInterface2)
	router1.AddNetworkInterface(router1AccessInterface1)
	router2.AddNetworkInterface(router2NetworkInterface1)
	router2.AddNetworkInterface(router2NetworkInterface2)
	router2.AddNetworkInterface(router2AccessInterface1)
	host1.AddNetworkInterface(host1NetworkInterface1)
	host2.AddNetworkInterface(host2NetworkInterface1)

	host1.AddRoute(host2NetworkInterface1, []*networkinterface.Config{router1AccessInterface1})
	host2.AddRoute(host1NetworkInterface1, []*networkinterface.Config{router2AccessInterface1})

	router1.AddRoute(host1NetworkInterface1, []*networkinterface.Config{host1NetworkInterface1})
	router1.AddRoute(host2NetworkInterface1, []*networkinterface.Config{
		router2NetworkInterface1,
		router2NetworkInterface2,
	})

	router2.AddRoute(host2NetworkInterface1, []*networkinterface.Config{host2NetworkInterface1})
	router2.AddRoute(host1NetworkInterface1, []*networkinterface.Config{
		router1NetworkInterface1,
		router1NetworkInterface2,
	})

	out, err := yaml.Marshal(config1)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
