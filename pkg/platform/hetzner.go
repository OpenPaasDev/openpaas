package platform

import (
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"golang.org/x/net/context"
)

type Hetzner struct{}

func (s *Hetzner) Preparation(context.Context, *conf.Config) error { return nil }

func (s *Hetzner) Cleanup(context.Context, *conf.Config) error { return nil }
