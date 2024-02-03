package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
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
			runner := provider.DefaultRunner()
			err = runner.RunAll(ctx, cnf, inv)

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
			runner := provider.DefaultRunner()
			err = runner.RunAll(ctx, cnf, inv)
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
