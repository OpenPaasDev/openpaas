package provider

import (
	"context"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
)

type Ansible struct{}

type AnsibleConfig struct {
	Inventory string `yaml:"inventory"`
}

func (s *Ansible) Run(ctx context.Context, globalConf *conf.Config, providerConfig interface{}, inventory *ansible.Inventory) error {
	return nil
}
