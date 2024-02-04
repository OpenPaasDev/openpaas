package platform

import (
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"golang.org/x/net/context"
)

var platforms = map[string]Platform{
	"hetzner": &Hetzner{},
}

func RunPreparation(ctx context.Context, cnf *conf.Config) error {
	if platform, found := platforms[cnf.DC]; found {
		err := platform.Preparation(ctx, cnf)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunCleanup(ctx context.Context, cnf *conf.Config) error {
	if platform, found := platforms[cnf.DC]; found {
		err := platform.Cleanup(ctx, cnf)
		if err != nil {
			return err
		}
	}
	return nil
}
