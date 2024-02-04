package provider

import (
	"fmt"
	"path/filepath"
	"testing"

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
	t.Skip("Not implemented")
}

func TestAnsibleProvider(t *testing.T) {
	t.Skip("Not implemented")
}

func TestRunProviders(t *testing.T) {
	t.Skip("Not implemented")
}
