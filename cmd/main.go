package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
	"github.com/OpenPaasDev/core/pkg/terraform"
)

func main() {
	ctx := context.Background()
	cnf, err := conf.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	err = terraform.GenerateTerraform(cnf)
	if err != nil {
		panic(err)
	}

	tf, err := terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	err = tf.Apply(ctx)
	if err != nil {
		panic(err)
	}
	_, err = ansible.GenerateInventory(cnf)
	if err != nil {
		panic(err)
	}
}
