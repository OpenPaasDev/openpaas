package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	ssh "github.com/helloyi/go-sshclient"
	"github.com/spf13/cobra"

	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
	"github.com/OpenPaasDev/core/pkg/terraform"
)

func main() {

	err := os.Setenv("ANSIBLE_HOST_KEY_CHECKING", "False")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "openpaas",
		Short: "sets up the openpaas",
		Long:  `openpaas`,
		Run: func(cmd *cobra.Command, args []string) {
			e := cmd.Help()
			if e != nil {
				panic(e)
			}
		},
	}

	rootCmd.AddCommand()

	ctx := context.Background()
	cnf, inventory, err := initStack(ctx, "config.yaml")
	if err != nil {
		panic(err)
	}
	updateNodes(cnf, inventory)
}

func initStack(ctx context.Context, file string) (*conf.Config, *ansible.Inventory, error) {
	cnf, err := conf.Load(file)
	if err != nil {
		return nil, nil, err
	}

	err = terraform.GenerateTerraform(cnf)
	if err != nil {
		return nil, nil, err
	}

	tf, err := terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), os.Stdout, os.Stderr)
	if err != nil {
		return nil, nil, err
	}

	err = tf.Apply(ctx, conf.LoadTFExecVars())
	if err != nil {
		return nil, nil, err
	}
	os.Remove(filepath.Join(cnf.BaseDir, "inventory-output.json")) //nolint
	f, err := os.OpenFile(filepath.Join(cnf.BaseDir, "inventory-output.json"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		e := f.Close()
		if e != nil {
			panic(e)
		}
	}()
	tf, err = terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), f, os.Stderr)
	if err != nil {
		return nil, nil, err
	}
	_, err = tf.Output(ctx)
	if err != nil {
		return nil, nil, err
	}
	inventory, err := ansible.GenerateInventory(cnf)
	if err != nil {
		return nil, nil, err
	}
	return cnf, inventory, nil
}

func updateNodes(cnf *conf.Config, inventory *ansible.Inventory) {
	serverIps := []string{}
	for k := range inventory.All.Children {
		for _, v := range inventory.All.Children[k].Hosts {
			serverIps = append(serverIps, v.PublicIP)
		}
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
