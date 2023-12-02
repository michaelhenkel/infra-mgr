package instance

import (
	networkinterface "github.com/michaelhenkel/infra-mgr/networkInterface"
	"github.com/michaelhenkel/infra-mgr/object"
)

type Instance struct {
	Config            *Config
	NetworkInterfaces map[string]*networkinterface.NetworkInterface
}

type Config struct {
	Name   string
	Vcpu   int
	Memory string
	Image  string
}

func New(name string, vcpu int, memory string, image string) *Config {
	return &Config{
		Name:   name,
		Vcpu:   vcpu,
		Memory: memory,
		Image:  image,
	}
}

func (config *Config) Add(object.Object) {

}

func (config *Config) GetName() string {
	return config.Name
}

func (config *Config) ObjectType() object.ObjectType {
	return object.Instance
}

func (config *Config) Build() *Instance {
	return &Instance{
		Config:            config,
		NetworkInterfaces: make(map[string]*networkinterface.NetworkInterface),
	}
}

func (i *Instance) Add(o object.Object) {
	switch o.ObjectType() {
	case object.NetworkInterface:
		i.NetworkInterfaces[o.GetName()] = o.(*networkinterface.NetworkInterface)
	}
}

func (i *Instance) ObjectType() object.ObjectType {
	return object.Instance
}

func (i *Instance) GetName() string {
	return i.Config.Name
}
