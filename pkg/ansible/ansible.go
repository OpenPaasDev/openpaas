package ansible

import (
	"fmt"
	"github.com/OpenPaasDev/openpaas/pkg/runtime"
	"os"
)

type Client interface {
	Run(playbookFile, varFile string) error
}

type ansibleClient struct {
	inventory string
	user      string
}

func NewClient(inventory, user string) Client {
	return &ansibleClient{
		inventory: inventory,
		user:      user,
	}
}

func (client *ansibleClient) Run(playbookFile string, varFile string) error {
	if varFile != "" {
		return runtime.Exec(&runtime.EmptyEnv{}, fmt.Sprintf("ansible-playbook %s -i %s -u %s -e @%s ", playbookFile, client.inventory, client.user, varFile), os.Stdout)
	}
	return runtime.Exec(&runtime.EmptyEnv{}, fmt.Sprintf("ansible-playbook %s -i %s -u %s", playbookFile, client.inventory, client.user), os.Stdout)
}
