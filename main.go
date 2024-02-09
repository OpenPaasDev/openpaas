/*
OpenPaaS provisions infrastructure based on the provided config,
while leveraging existing IaC tooling.

The intention is to simplify the configuration for a majority of application, which
can be server with efficient, traditional infrastructure. This covers a majority of
companies.

Using tools like Terraform and Ansible underneath means that you can inspect the
output and understand what is going to happen, without having to learn a new IaC tool.

Currently, it defaults to Hetzner, but it can easily be expanded to target other platforms.

NOTE: the executable sets the Environment variable ANSIBLE_HOST_KEY_CHECKING to False during execution.
This disables host key checking on SSH connections.

Usage:

	openpaas [flags]
	openpaas [command]

Available Commands:

	bootstrap   Bootstrap your platform
	completion  Generate the autocompletion script for the specified shell
	help        Help about any command
	sync        Sync your platform

Flags:

	-h, --help   help for openpaas

Use "openpaas [command] --help" for more information about a command.

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

	"github.com/spf13/cobra"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"github.com/OpenPaasDev/openpaas/pkg/platform"
	"github.com/OpenPaasDev/openpaas/pkg/provider"
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
		Short: "Sets up the OpenPaas platform",
		Long: `OpenPaaS provisions infrastructure based on the provided config, 
while leveraging existing IaC tooling.

The intention is to simplify the configuration for a majority of application, which 
can be server with efficient, traditional infrastructure. This covers a majority of 
companies. 

Using tools like Terraform and Ansible underneath means that you can inspect the 
output and understand what is going to happen, without having to learn a new IaC tool.

Currently, it defaults to Hetzner, but it can easily be expanded to target other platforms.`,
		Run: func(cmd *cobra.Command, args []string) {
			e := cmd.Help()
			if e != nil {
				panic(e)
			}
		},
	}

	rootCmd.AddCommand(syncCmd())
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func syncCmd() *cobra.Command {
	var configFile string
	var terraformVersion string
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync your platform",
		Long:  `Syncs any changes in your platform configuration with the deployment.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cnf, inv, err := initStack(ctx, configFile, terraformVersion)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			runner := provider.DefaultRunner()
			err = runner.RunAll(ctx, cnf, inv)

			if err != nil {
				panic(err)
			}
		},
	}

	addFlags(cmd, &configFile, &terraformVersion)

	return cmd
}

func addFlags(cmd *cobra.Command, file *string, terraformVersion *string) {
	cmd.Flags().StringVarP(file, "config.file", "f", "./config.yaml", "OpenPaaS configuration file to use")
	cmd.Flags().StringVarP(terraformVersion, "terraform.version", "t", "1.7.2", "Terraform version to use")
}

func initStack(ctx context.Context, file string, terraformVersion string) (*conf.Config, *ansible.Inventory, error) {
	cnf, err := conf.Load(file)
	if err != nil {
		return nil, nil, err
	}

	err = platform.RunPreparation(ctx, cnf)
	if err != nil {
		return nil, nil, err
	}

	inventory, err := initTerraform(ctx, cnf, terraformVersion)
	if err != nil {
		return nil, nil, err
	}

	err = platform.RunCleanup(ctx, cnf)
	if err != nil {
		return nil, nil, err
	}

	return cnf, inventory, nil
}

func initTerraform(ctx context.Context, cnf *conf.Config, terraformVersion string) (*ansible.Inventory, error) {
	err := terraform.GenerateTerraform(cnf)
	if err != nil {
		return nil, err
	}

	//TODO we initialise here and then again in line 175, is this needed? why?
	tf, err := terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), terraformVersion, os.Stdout, os.Stderr)
	if err != nil {
		return nil, err
	}

	err = tf.Apply(ctx, conf.LoadTFExecVars())
	if err != nil {
		return nil, err
	}
	os.Remove(filepath.Join(cnf.BaseDir, "inventory-output.json")) //nolint
	f, err := os.OpenFile(filepath.Join(cnf.BaseDir, "inventory-output.json"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := f.Close()
		if e != nil {
			panic(e)
		}
	}()
	// Terraform initialised again, this time with output to file f, to capture outputs as json
	tf, err = terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), terraformVersion, f, os.Stderr)
	if err != nil {
		return nil, err
	}

	_, err = tf.Output(ctx)
	if err != nil {
		return nil, err
	}
	inventory, err := ansible.GenerateInventory(cnf)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}
