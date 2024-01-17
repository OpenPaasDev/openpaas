/*
Core provisions infrastructure based on the provided config.
Underneath, it uses Terraform and Ansible to provision and configure the servers.

Without an explicit file, it processes ./config.yaml by default.

NOTE: the executable sets the Environment variable ANSIBLE_HOST_KEY_CHECKING to False during execution.
This disables host key checking on SSH connections.

Usage:

	core [flags] [file]

The flags are:

	-h,--help
	    Shows this output.

Environment Variables:

	HETZNER_TOKEN
		The Hetzner access token to use to apply the Terraform configuration
*/
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	ssh "github.com/helloyi/go-sshclient"
	"github.com/spf13/cobra"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"github.com/OpenPaasDev/openpaas/pkg/provider"
	"github.com/OpenPaasDev/openpaas/pkg/state"
	"github.com/OpenPaasDev/openpaas/pkg/terraform"
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

	rootCmd.AddCommand(bootstrap(), syncCmd())
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func syncCmd() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync your platform",
		Long:  `Sync the platform`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cnf, inv, err := initStack(ctx, configFile)
			if err != nil {
				panic(err)
			}
			d := state.Init(cnf.BaseDir)

			err = d.Sync(cnf, inv)
			if err != nil {
				panic(err)
			}
			err = provider.RunAll(ctx, cnf, inv)
			if err != nil {
				panic(err)
			}
		},
	}

	addFlags(cmd, &configFile)

	return cmd
}

func bootstrap() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap your platform",
		Long:  `bootstrap the platform`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cnf, inv, err := initStack(ctx, configFile)
			if err != nil {
				panic(err)
			}
			updateNodes(cnf, inv)
			d := state.Init(cnf.BaseDir)

			err = d.Sync(cnf, inv)
			if err != nil {
				panic(err)
			}

			err = provider.RunAll(ctx, cnf, inv)
			if err != nil {
				panic(err)
			}
		},
	}

	addFlags(cmd, &configFile)

	return cmd
}

func addFlags(cmd *cobra.Command, file *string) {
	cmd.Flags().StringVarP(file, "config.file", "f", "", "configuration file")

	err := cmd.MarkFlagRequired("config.file")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

	//TODO we initialise here and then again in line 175, is this needed? why?
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
