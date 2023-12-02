package main

import (
	"fmt"

	"github.com/michaelhenkel/infra-mgr/config"
	"github.com/michaelhenkel/infra-mgr/instance"
	"github.com/michaelhenkel/infra-mgr/network"
	networkinterface "github.com/michaelhenkel/infra-mgr/networkInterface"
	"github.com/michaelhenkel/infra-mgr/route"
	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("Hello World")

	config1 := config.New("infra1")

	host1 := instance.New("host1", 1, "1G", "image1")
	config1.Instances = append(config1.Instances, host1)

	host2 := instance.New("host2", 1, "1G", "image1")
	config1.Instances = append(config1.Instances, host2)

	router1 := instance.New("router1", 1, "1G", "image1")
	config1.Instances = append(config1.Instances, router1)

	router2 := instance.New("router2", 1, "1G", "image1")
	config1.Instances = append(config1.Instances, router2)

	r1r2Network1 := network.New("r1r2_1", "10.0.0.0/24")
	config1.Networks = append(config1.Networks, r1r2Network1)

	r1r2Network2 := network.New("r1r2_2", "10.0.1.0/24")
	config1.Networks = append(config1.Networks, r1r2Network2)

	accessNetwork1 := network.New("access1", "10.0.2.0/24")
	config1.Networks = append(config1.Networks, accessNetwork1)

	accessNetwork2 := network.New("access2", "10.0.3.0/24")
	config1.Networks = append(config1.Networks, accessNetwork2)

	router1NetworkInterface1 := networkinterface.New("r1_eth1", r1r2Network1, router1)
	config1.Interfaces = append(config1.Interfaces, router1NetworkInterface1)

	router1NetworkInterface2 := networkinterface.New("r1_eth2", r1r2Network2, router1)
	config1.Interfaces = append(config1.Interfaces, router1NetworkInterface2)

	router1AccessInterface1 := networkinterface.New("r1_eth3", accessNetwork1, router1)
	config1.Interfaces = append(config1.Interfaces, router1AccessInterface1)

	router2NetworkInterface1 := networkinterface.New("r2_eth1", r1r2Network1, router2)
	config1.Interfaces = append(config1.Interfaces, router2NetworkInterface1)

	router2NetworkInterface2 := networkinterface.New("r2_eth2", r1r2Network2, router2)
	config1.Interfaces = append(config1.Interfaces, router2NetworkInterface2)

	router2AccessInterface1 := networkinterface.New("r2_eth3", accessNetwork2, router2)
	config1.Interfaces = append(config1.Interfaces, router2AccessInterface1)

	host1NetworkInterface1 := networkinterface.New("h1_eth1", accessNetwork1, host1)
	config1.Interfaces = append(config1.Interfaces, host1NetworkInterface1)

	host2NetworkInterface1 := networkinterface.New("h2_eth1", accessNetwork2, host2)
	config1.Interfaces = append(config1.Interfaces, host2NetworkInterface1)

	host1accessNetwork2Route := route.New("host1accessNetwork2Route",
		route.Source{
			Instance:         "host1",
			NetworkInterface: "h1_eth1",
		}, accessNetwork2,
		[]route.NextHop{
			{
				Instance:         "router1",
				NetworkInterface: "r1_eth3",
			},
		})
	config1.Routes = append(config1.Routes, host1accessNetwork2Route)

	host2accessNetwork1Route := route.New("host2accessNetwork1Route",
		route.Source{
			Instance:         "host2",
			NetworkInterface: "h2_eth1",
		}, accessNetwork1,
		[]route.NextHop{
			{
				Instance:         "router2",
				NetworkInterface: "r2_eth3",
			},
		})
	config1.Routes = append(config1.Routes, host2accessNetwork1Route)

	out, err := yaml.Marshal(config1)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

	infrastructure := config1.Build()

	out2, err := yaml.Marshal(infrastructure)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out2))

}
