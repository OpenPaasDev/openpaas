package provider

import (
	"context"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
)

type Runner interface {
	RunAll(ctx context.Context, cnf *conf.Config, inventory *ansible.Inventory) error
}
type defaultRunner struct {
	providers map[string]Service
}

func DefaultRunner() Runner {
	return &defaultRunner{
		providers: map[string]Service{
			"ansible": &Ansible{makeClient: ansible.NewClient},
		},
	}
}

func (runner *defaultRunner) RunAll(ctx context.Context, cnf *conf.Config, inventory *ansible.Inventory) error {
	for k, providerConfig := range cnf.Providers {
		if _, ok := runner.providers[k]; ok {
			err := runner.providers[k].Run(ctx, providerConfig, inventory)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
