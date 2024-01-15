package provider

import (
	"context"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
)

type Ansible struct{}

func (s *Ansible) Run(context.Context, *conf.Config, *ansible.Inventory) error {
	return nil
}
