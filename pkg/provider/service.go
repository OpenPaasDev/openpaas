package provider

import (
	"context"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
)

type Service interface {
	Run(context.Context, *conf.Config, *ansible.Inventory) error
}
