package provider

import (
	"context"

	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
)

type Ansible struct{}

func (s *Ansible) Run(context.Context, *conf.Config, *ansible.Inventory) error {
	return nil
}
