package provider

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
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

//go:embed defaults/hardening-ubuntu-2204.yml
var hardeningUbuntu string

// add them in the order they need to be applied
var defaultPlaybooks []string = []string{hardeningUbuntu}

func (s *Ansible) Run(ctx context.Context, providerConfig interface{}, inventory *ansible.Inventory) error {
	conf, err := asAnsibleConfig(providerConfig)
	if err != nil {
		return err
	}
	ansibleClient := s.makeClient(inventory.Path, conf.SudoUser)

	corePlaybooks, err := prepareDefaultPlaybooks(defaultPlaybooks)
	if err != nil {
		return err
	}

	// runs all playbooks provided by the user after we have run the default ones
	conf.Playbooks = append(corePlaybooks, conf.Playbooks...)
	for _, playbook := range conf.Playbooks {
		err = s.runPlaybook(playbook, conf, ansibleClient)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Ansible) runPlaybook(playbook Playbook, conf *AnsibleConfig, ansibleClient ansible.Client) error {
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
	return nil
}

// prepareDefaultPlaybooks runs the default playbooks in the `defaultPlaybooks` slice, in order
func prepareDefaultPlaybooks(playbooks []string) ([]Playbook, error) {
	var ymlFiles []Playbook
	// first create temporal files for the default playbooks
	for _, playbookContent := range playbooks {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "default-playbook-*.yml")
		defer tempFile.Close() //nolint
		if err != nil {
			fmt.Printf("Error creating temp file %s: %s\n", tempFile.Name(), err)
			return nil, err
		}

		// Write the content to the temporary file
		_, err = tempFile.WriteString(playbookContent)
		if err != nil {
			fmt.Printf("Error writing to temp file %s: %s\n", tempFile.Name(), err)
			return nil, err
		}

		// Get the absolute path of the temporary file
		absPath, err := filepath.Abs(tempFile.Name())
		if err != nil {
			fmt.Printf("Error getting absolute path to temp file %s: %s\n", tempFile.Name(), err)
			return nil, err
		}

		// Store the absolute path in the slice
		ymlFiles = append(ymlFiles, Playbook{Name: absPath, Vars: map[string]string{}})
	}

	return ymlFiles, nil
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
