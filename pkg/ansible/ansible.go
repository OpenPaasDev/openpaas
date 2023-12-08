package ansible

import (
	"fmt"
	"os"

	"github.com/OpenPaasDev/core/pkg/runtime"
)

type Client interface {
	Run(file string) error
}

type ansibleClient struct {
	inventory   string
	secretsFile string
	user        string
	configPath  string
}

func NewClient(inventory, user, configPath, secretsFile string) Client {
	return &ansibleClient{
		inventory:   inventory,
		secretsFile: secretsFile,
		user:        user,
		configPath:  configPath,
	}
}

func (client *ansibleClient) Run(file string) error {
	if client.secretsFile != "" && client.configPath != "" {
		return runtime.Exec(&runtime.EmptyEnv{}, fmt.Sprintf("ansible-playbook %s -i %s -u %s -e @%s -e @%s", file, client.inventory, client.user, client.secretsFile, client.configPath), os.Stdout)
	}
	if client.secretsFile == "" && client.configPath == "" {
		return runtime.Exec(&runtime.EmptyEnv{}, fmt.Sprintf("ansible-playbook %s -i %s -u %s", file, client.inventory, client.user), os.Stdout)
	}
	if client.secretsFile == "" {
		return runtime.Exec(&runtime.EmptyEnv{}, fmt.Sprintf("ansible-playbook %s -i %s -u %s -e @%s", file, client.inventory, client.user, client.configPath), os.Stdout)
	}
	return fmt.Errorf("insufficient client configuration provided")
}
