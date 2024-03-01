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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/OpenPaasDev/openpaas/pkg/util"
	"github.com/spf13/cobra"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"github.com/OpenPaasDev/openpaas/pkg/platform"
	"github.com/OpenPaasDev/openpaas/pkg/provider"
	"github.com/OpenPaasDev/openpaas/pkg/terraform"
)

// constants for versions of apps to install as dependencies
// we don't have a version for Ansible as it is installed via brew or pip

const HcloudVersion = "1.42.0"
const TerraformVersion = "1.7.4"

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
	rootCmd.AddCommand(initCmd())
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

	cmd.Flags().StringVarP(&configFile, "config.file", "f", "./config.yaml", "OpenPaaS configuration file to use")
	addTerraformFlag(cmd, &terraformVersion)

	return cmd
}

func addTerraformFlag(cmd *cobra.Command, terraformVersion *string) {
	cmd.Flags().StringVarP(terraformVersion, "terraform.version", "t", TerraformVersion, "Terraform version to use")
}

func initCmd() *cobra.Command {
	var terraformVersion string
	var hcloudVersion string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Install required dependencies for OpenPaaS",
		Long:  `Install required dependencies for OpenPaaS. This includes tooling like Terraform and Ansible, plus any other required dependencies.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := installDependencies(ctx, terraformVersion, hcloudVersion)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		},
	}

	addTerraformFlag(cmd, &terraformVersion)
	cmd.Flags().StringVar(&hcloudVersion, "hcloud.version", HcloudVersion, "Hcloud version to install")

	return cmd
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

func initTerraform(ctx context.Context, cnf *conf.Config, terraformExecPath string) (*ansible.Inventory, error) {
	err := terraform.GenerateTerraform(cnf)
	if err != nil {
		return nil, err
	}

	//TODO we initialise here and then again in line 175, is this needed? why?
	tf, err := terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), terraformExecPath, os.Stdout, os.Stderr)
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
	tf, err = terraform.InitTf(ctx, filepath.Join(cnf.BaseDir, "terraform"), terraformExecPath, f, os.Stderr)
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

func installDependencies(ctx context.Context, terraformVersion string, hcloudVersion string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("Error initialising dependencies: Windows OS not supported")
		return errors.ErrUnsupported
	}

	_, err := terraform.GetTerraformExecutablePath(ctx, terraformVersion)
	if err != nil {
		return err
	}

	err = installAnsible()
	if err != nil {
		return err
	}

	err = installHcloud(ctx, hcloudVersion)
	if err != nil {
		return err
	}

	return nil
}

func installAnsible() error {
	if !util.IsAnsibleInstalled() {
		// Brew can also be installed in linux it seems :)
		if util.IsBrewInstalled() {
			fmt.Println("'ansible': Installing latest version available in brew")
			err := util.RunCmd("brew", "install", "ansible")
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("'ansible': Installing via pip")
			if util.IsPipInstalled() {
				err := util.RunCmd("pip", "install", "--user", "ansible")
				if err != nil {
					return err
				}
			} else {
				return errors.New("'ansible': Error initialising dependencies: pip not installed. Please install pip or install ansible separately")
			}
		}
		// verify ansible runs
		if !util.IsAnsibleInstalled() {
			return errors.New("Error initialising dependencies: ansible not installed")
		}
	}
	return nil
}

func installHcloud(ctx context.Context, hcloudVersion string) error {
	if !util.IsHCloudInstalled() {
		// Brew can also be installed in linux it seems :)
		if util.IsBrewInstalled() {
			fmt.Println("'hcloud': Installing latest version available in brew")
			err := util.RunCmd("brew", "install", "hcloud")
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("'hcloud': Installing version %s from artefacts", hcloudVersion)
			extension := "tar.gz"
			destPath := "/usr/local/bin"
			tempPath := "/tmp/hcloud.tar.gz"
			system := runtime.GOOS
			arch := runtime.GOARCH
			urlTemplate := "https://github.com/hetznercloud/cli/releases/download/v%s/hcloud-%s-%s.%s"
			downloadUrl := fmt.Sprintf(urlTemplate, hcloudVersion, system, arch, extension)

			if err := util.DownloadFile(ctx, downloadUrl, tempPath); err != nil {
				fmt.Println("Error downloading file:", err)
				return err
			}

			if err := util.ExtractTarGz(tempPath, destPath); err != nil {
				fmt.Println("Error extracting file:", err)
				return err
			}
		}
		// verify hcloud runs
		if !util.IsHCloudInstalled() {
			return errors.New("Error initialising dependencies: hcloud not installed")
		}
	}
	return nil
}
