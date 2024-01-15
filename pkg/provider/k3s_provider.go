package provider

import (
	"context"
	"fmt"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
)

type K3S struct {
}

type K3sSettings struct {
	NetworkInterface string `yaml:"internal_network_interface_name"`
	ServerGroup      string `yaml:"server_group"`
	AgentGroup       string `yaml:"agent_group"`
}

func (s *K3S) Run(context.Context, *conf.Config, *ansible.Inventory) error {
	fmt.Println("Run K3S installer")
	return nil
}
