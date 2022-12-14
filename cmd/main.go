package main

import (
	"context"
	"fmt"
	"os"

	"github.com/OpenPaaSDev/openpaas/internal"
	"github.com/OpenPaaSDev/openpaas/internal/conf"
	"github.com/OpenPaaSDev/openpaas/internal/runtime"
	"github.com/spf13/cobra"
)

func main() {
	err := os.Setenv("ANSIBLE_HOST_KEY_CHECKING", "False")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !runtime.HasDependencies() {
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "openpaas",
		Short: "sets up the openpaas",
		Long:  `openpaas - setups up Consul, Vault & Nomad with ACL & Service Mesh enabled`,
		Run: func(cmd *cobra.Command, args []string) {
			err = cmd.Help()
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.AddCommand(sync(), envRC())

	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func sync() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "bootstraps and starts a cluster or syncs the cluster to its desired state",
		Long:  `bootstraps and starts a cluster or syncs the cluster to its desired state`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := conf.Load(configFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = internal.Bootstrap(context.Background(), config, configFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addFlags(cmd, &configFile)

	return cmd
}

func envRC() *cobra.Command {
	var configFile string
	var targetDir string
	cmd := &cobra.Command{
		Use:   "genenv",
		Short: "Generate env file to source for your environment",
		Long:  `Generate env file to source for your environment`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := conf.Load(configFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = internal.GenerateEnvFile(config, targetDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addFlags(cmd, &configFile)
	cmd.Flags().StringVarP(&targetDir, "target.dir", "t", "", "target directory of .envrc file")

	err := cmd.MarkFlagRequired("config.file")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
