package provider

import (
	"context"

	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
)

func RunAll(ctx context.Context, cnf *conf.Config, inventory *ansible.Inventory) error {
	providers := map[string]Service{
		"ansible": &Ansible{},
		"k3s":     &K3S{},
	}
	for _, provider := range cnf.Providers {
		if service, found := providers[provider.Name]; found {
			err := service.Run(ctx, cnf, inventory)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
