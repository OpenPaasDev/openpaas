package provider

import (
	"context"
	"os"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
)

type Ansible struct{}

type AnsibleConfig struct {
	SudoUser   string            `yaml:"sudo_user"`
	GlobalVars map[string]string `yaml:"global_vars"`
	Playbooks  []Playbook        `yaml:"playbooks"`
}

type Playbook struct {
	Name string            `yaml:"file"`
	Vars map[string]string `yaml:"vars"`
}

func (s *Ansible) Run(ctx context.Context, providerConfig interface{}, inventory *ansible.Inventory) error {
	conf, err := asAnsibleConfig(providerConfig)
	if err != nil {
		return err
	}
	ansibleClient := ansible.NewClient(inventory.Path, conf.SudoUser)
	for _, playbook := range conf.Playbooks {
		varsFile, err := generateVarsFile(playbook.Vars, conf.GlobalVars)
		if err != nil {
			return err
		}
		if varsFile != "" {
			defer os.Remove(varsFile) //nolint
		}
		err = ansibleClient.Run(playbook.Name, varsFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func asAnsibleConfig(providerConfig interface{}) (*AnsibleConfig, error) {
	return nil, nil
}

func generateVarsFile(vars map[string]string, globalVars map[string]string) (string, error) {
	return "", nil
}
