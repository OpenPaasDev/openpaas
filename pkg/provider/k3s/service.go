package k3s

import (
	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
)

type Service struct {
}

type K3sSettings struct {
	NetworkInterface string `yaml:"internal_network_interface_name"`
	ServerGroup      string `yaml:"server_group"`
	AgentGroup       string `yaml:"agent_group"`
}

func (s *Service) Build(cnf conf.Config, inventory ansible.Inventory) error {
	return nil
}
