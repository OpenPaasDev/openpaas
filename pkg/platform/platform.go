package platform

import (
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"golang.org/x/net/context"
)

type Platform interface {
	// intended for changes to the platform that we want to run before any terraform or ansible script. Example: modify allowed ssh keys
	Preparation(context.Context, *conf.Config) error
	// intended for cleanup after we run terraform and ansible scripts.
	Cleanup(context.Context, *conf.Config) error
}
