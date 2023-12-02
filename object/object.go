package object

type Object interface {
	GetName() string
	Add(Object)
	ObjectType() ObjectType
}

type ObjectType string

const (
	Instance         ObjectType = "instance"
	Network          ObjectType = "network"
	NetworkInterface ObjectType = "networkinterface"
	Route            ObjectType = "route"
)
