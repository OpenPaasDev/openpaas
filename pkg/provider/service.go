package provider

import (
	"context"

	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
)

type Service interface {
	Run(context.Context, *conf.Config, *ansible.Inventory) error
}
