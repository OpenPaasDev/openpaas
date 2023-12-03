package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	ssh "github.com/helloyi/go-sshclient"

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
	inventory, err := ansible.GenerateInventory(cnf)
	if err != nil {
		panic(err)
	}

	// TODO this is running stuff on the servers

	serverIps := []string{}
	for k := range inventory.All.Children["servers"].Hosts {
		serverIps = append(serverIps, k)
	}
	clients := []string{}
	for k := range inventory.All.Children["clients"].Hosts {
		clients = append(clients, k)
	}

	var wg sync.WaitGroup
	for _, ip := range serverIps {
		wg.Add(1)
		fmt.Println(fmt.Sprintf("%s:22", ip))
		go func(ip string) {
			client, err := ssh.DialWithKey(fmt.Sprintf("%s:22", ip), cnf.CloudProviderConfig.User, cnf.CloudProviderConfig.SSHKey)
			defer client.Close() //nolint
			if err != nil {
				panic(err)
			}
			script := client.Cmd("sudo apt-get update").Cmd("sudo apt-get upgrade -y")
			script.SetStdio(os.Stdout, os.Stderr)
			err = script.Run()
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(ip)
	}
	wg.Wait()

}
