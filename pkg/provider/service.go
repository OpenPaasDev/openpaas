package provider

import (
	"context"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
)

type Service interface {
	// Run runs the provider, takes context, global config, provider config and inventory.
	// provider config is the interface argument because it is specific to the provider
	Run(context.Context, interface{}, *ansible.Inventory) error
}
