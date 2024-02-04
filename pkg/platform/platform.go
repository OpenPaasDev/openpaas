package platform

import (
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"golang.org/x/net/context"
)

type Platform interface {
	// Preparation is intended for changes to the platform that we want to run before any terraform or ansible script. Example: modify allowed ssh keys
	Preparation(context.Context, *conf.Config) error
	// Cleanup is intended for cleanup after we run terraform and ansible scripts.
	Cleanup(context.Context, *conf.Config) error
}
