package provider

import (
	"context"
	"os"
	"strings"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/util"
	"gopkg.in/yaml.v3"
)

type Ansible struct {
	makeClient func(inventoryFile string, sudoUser string) ansible.Client
}

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
	ansibleClient := s.makeClient(inventory.Path, conf.SudoUser)
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
	yamlData, err := yaml.Marshal(providerConfig)
	if err != nil {
		return nil, err
	}
	var cnf AnsibleConfig
	err = yaml.Unmarshal(yamlData, &cnf)
	if err != nil {
		return nil, err
	}
	return &cnf, nil
}

func generateVarsFile(vars map[string]string, globalVars map[string]string) (string, error) {
	outputMap := make(map[string]string)
	for k, v := range globalVars {
		if v == strings.ToUpper(v) && os.Getenv(v) != "" {
			outputMap[k] = os.Getenv(v)
		} else {
			outputMap[k] = v
		}
	}
	for k, v := range vars {
		if v == strings.ToUpper(v) && os.Getenv(v) != "" {
			outputMap[k] = os.Getenv(v)
		} else {
			outputMap[k] = v
		}
	}
	yamlData, err := yaml.Marshal(outputMap)
	if err != nil {
		return "", err
	}
	fileName := util.RandString(15)
	err = os.WriteFile(fileName, yamlData, 0644) //nolint
	if err != nil {
		return "", err
	}
	return fileName, nil
}
