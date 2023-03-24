package main

import (
	"github.com/OpenPaasDev/core/pkg/conf"
	"github.com/OpenPaasDev/core/pkg/terraform"
)

func main() {
	cnf, err := conf.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	err = terraform.GenerateTerraform(cnf)
	if err != nil {
		panic(err)
	}
}
