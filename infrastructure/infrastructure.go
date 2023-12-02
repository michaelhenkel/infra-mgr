package infrastructure

import (
	"github.com/michaelhenkel/infra-mgr/instance"
	"github.com/michaelhenkel/infra-mgr/network"
	"github.com/michaelhenkel/infra-mgr/object"
)

type Infrastructure struct {
	Name      string
	Instances map[string]*instance.Instance
	Networks  map[string]*network.Network
}

func New(name string) *Infrastructure {
	return &Infrastructure{
		Instances: make(map[string]*instance.Instance),
		Networks:  make(map[string]*network.Network),
	}
}

func (i *Infrastructure) Add(o object.Object) {
	switch o.ObjectType() {
	case object.Instance:
		i.Instances[o.GetName()] = o.(*instance.Instance)
	case object.Network:
		i.Networks[o.GetName()] = o.(*network.Network)
	}
}
