package provider

import (
	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
)

type Service interface {
	Build(cnf conf.Config, inventory ansible.Inventory) error
}
