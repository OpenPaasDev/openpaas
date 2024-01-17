package provider

import (
	"context"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
)

func RunAll(ctx context.Context, cnf *conf.Config, inventory *ansible.Inventory) error {
	providers := map[string]Service{
		"ansible": &Ansible{},
	}
	for k, providerConfig := range cnf.Providers {
		if _, ok := providers[k]; ok {
			err := providers[k].Run(ctx, cnf, providerConfig, inventory)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
