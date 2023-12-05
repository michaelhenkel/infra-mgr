package lxdmgr

import (
	incus "github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/api"
)

type LxdMgr struct {
	Conn incus.InstanceServer
}

func New() (*LxdMgr, error) {
	if conn, err := incus.ConnectIncusUnix("/var/snap/lxd/common/lxd/unix.socket", &incus.ConnectionArgs{}); err != nil {
		return nil, err
	} else {
		return &LxdMgr{
			Conn: conn,
		}, nil
	}
}

func (l *LxdMgr) CreateInstance(name string, image string, profile string) error {

	d, err := incus.ConnectSimpleStreams("https://images.linuxcontainers.org", nil)
	if err != nil {
		return err
	}

	// Resolve the alias
	alias, _, err := d.GetImageAlias("ubuntu/23.10/default")
	if err != nil {
		return err
	}

	img, _, err := d.GetImage(alias.Target)
	if err != nil {
		return err
	}

	req := api.InstancesPost{
		Name: name,
		Source: api.InstanceSource{
			Type:        "image",
			Fingerprint: img.Fingerprint,
			Project:     "default",
			Alias:       "ubuntu/23.10",
		},
		Type: api.InstanceType("virtual-machine"),
	}

	op, err := l.Conn.CreateInstanceFromImage(d, *img, req)
	if err != nil {
		return err
	}
	err = op.Wait()
	if err != nil {
		return err
	}

	reqState := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op2, err := l.Conn.UpdateInstanceState(name, reqState, "")
	if err != nil {
		return err
	}

	// Wait for the operation to complete
	err = op2.Wait()
	if err != nil {
		return err
	}
	return nil

}
