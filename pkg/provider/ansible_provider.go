package provider

import (
	"context"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
)

type Ansible struct{}

type AnsibleConfig struct {
	SudoUser   string            `yaml:"sudo_user"`
	SSHKeyPath string            `yaml:"ssh_key_path"`
	GlobalVars map[string]string `yaml:"global_vars"`
	Playbooks  []Playbook        `yaml:"playbooks"`
}

type Playbook struct {
	Name string            `yaml:"file"`
	Vars map[string]string `yaml:"vars"`
}

func (s *Ansible) Run(ctx context.Context, providerConfig interface{}, inventory *ansible.Inventory) error {
	return nil
}
