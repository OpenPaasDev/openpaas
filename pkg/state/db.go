package state

import (
	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
)

type Db struct {
	folder string
}

func Init(folder string) *Db {
	return &Db{folder: folder}
}

func (d *Db) Sync(config *conf.Config, inventory *ansible.Inventory) error {
	return nil
}
