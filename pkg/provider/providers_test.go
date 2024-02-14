package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnsibleProviderCanReadConfig(t *testing.T) {
	cnf, err := conf.Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	inf := cnf.Providers["ansible"]
	require.NotNil(t, inf)
	fmt.Println(inf)
	ansibleConf, err := asAnsibleConfig(inf)
	require.NoError(t, err)
	assert.NotNil(t, ansibleConf)
	assert.Equal(t, "root", ansibleConf.SudoUser)
	assert.Equal(t, "bar", ansibleConf.GlobalVars["foo"])
	assert.Equal(t, "playbooks/k3s.yml", ansibleConf.Playbooks[0].Name)
	assert.Equal(t, "baz", ansibleConf.Playbooks[0].Vars["foo"])
	assert.Equal(t, "playbooks/postgres.yml", ansibleConf.Playbooks[1].Name)
}

func TestGenerateAnsibleVarsFile(t *testing.T) {
	globalVars := map[string]string{"foo": "bar", "baz": "qux"}
	playbookVars := map[string]string{"foo": "baz", "qux": "QUUX_VAR"}
	err := os.Setenv("QUUX_VAR", "quux")
	require.NoError(t, err)
	varsFile, err := generateVarsFile(playbookVars, globalVars)
	require.NoError(t, err)
	assert.FileExists(t, varsFile)        //nolint
	content, err := os.ReadFile(varsFile) //nolint
	require.NoError(t, err)
	assert.Contains(t, string(content), "foo: baz")
	assert.Contains(t, string(content), "qux: quux")
	assert.Contains(t, string(content), "baz: qux")
	err = os.Unsetenv("QUUX_VAR")
	require.NoError(t, err)
	err = os.Remove(varsFile)
	require.NoError(t, err)
}

func TestAnsibleProvider(t *testing.T) {
	mockAnsible := &MockAnsibleClient{
		RunFunc: func(playbook string, varsFile string) error {
			if playbook != "playbooks/k3s.yml" && playbook != "playbooks/postgres.yml" && !strings.Contains(playbook, "default-playbook-") {
				return fmt.Errorf("unexpected playbook: %s", playbook)
			}
			assert.FileExists(t, varsFile) //nolint
			return nil
		},
	}
	ansibleProvider := &Ansible{
		makeClient: func(inventoryFile string, sudoUser string) ansible.Client {
			return mockAnsible
		},
	}
	inv := &ansible.Inventory{Path: "inventory"}
	cnf, err := conf.Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	inf := cnf.Providers["ansible"]
	err = ansibleProvider.Run(context.Background(), inf, inv)
	require.NoError(t, err)
	assert.Len(t, mockAnsible.RunCalls(), 2+len(defaultPlaybooks))
}

func TestRunProviders(t *testing.T) {
	mockAnsible := &MockAnsibleClient{
		RunFunc: func(playbook string, varsFile string) error {
			if playbook != "playbooks/k3s.yml" && playbook != "playbooks/postgres.yml" && !strings.Contains(playbook, "default-playbook-") {
				return fmt.Errorf("unexpected playbook: %s", playbook)
			}
			assert.FileExists(t, varsFile) //nolint
			return nil
		},
	}
	runner := &defaultRunner{
		providers: map[string]Service{
			"ansible": &Ansible{
				makeClient: func(inventoryFile string, sudoUser string) ansible.Client {
					return mockAnsible
				},
			},
		},
	}
	inv := &ansible.Inventory{Path: "inventory"}
	cnf, err := conf.Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	err = runner.RunAll(context.Background(), cnf, inv)
	require.NoError(t, err)
	assert.Len(t, mockAnsible.RunCalls(), 2+len(defaultPlaybooks))
}

func TestPrepareDefaultPlaybooks(t *testing.T) {
	testPlaybooks := []string{"playbookContent1", "playbookContent2", "playbookContent3"}
	playbooks, err := prepareDefaultPlaybooks(testPlaybooks)
	require.NoError(t, err)
	assert.Len(t, playbooks, 3)
	for _, pb := range playbooks {
		assert.FileExists(t, pb.Name) //nolint
		assert.Contains(t, pb.Name, "default-playbook-")
		assert.Contains(t, pb.Name, ".yml")
		assert.Equal(t, map[string]string{}, pb.Vars)
	}
}
