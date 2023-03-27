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

	err = tf.Apply(ctx, conf.LoadTFExecVars())
	if err != nil {
		panic(err)
	}
	os.Remove(filepath.Join(cnf.BaseDir, "inventory-output.json")) //nolint
	f, err := os.OpenFile(filepath.Join(cnf.BaseDir, "inventory-output.json"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		e := f.Close()
		if e != nil {
			panic(e)
		}
	}()
	tf, err = terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), f, os.Stderr)
	if err != nil {
		panic(err)
	}
	_, err = tf.Output(ctx)
	if err != nil {
		panic(err)
	}
	_, err = ansible.GenerateInventory(cnf)
	if err != nil {
		panic(err)
	}
}
