package conf

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, conf)
	fmt.Println(conf)

	assert.Equal(t, "config", conf.BaseDir)
	assert.Equal(t, "hetzner", conf.DC)
	assert.Equal(t, "root", conf.CloudProviderConfig.User)
	assert.Equal(t, []string{"wfaler", "pvillega"}, conf.CloudProviderConfig.GithubIds)
	assert.Equal(t, "/home/wfaler/.ssh/id_rsa", conf.CloudProviderConfig.SSHKey)
	assert.Equal(t, "hetzner", conf.CloudProviderConfig.Provider)

	assert.Len(t, conf.ServerGroups, 2)
	assert.Equal(t, "cpx31", conf.ServerGroups["clients"].InstanceType)
	assert.Equal(t, 2, conf.ServerGroups["clients"].Num)
	assert.Equal(t, 20, conf.ServerGroups["servers"].Volumes[0].Size)
	assert.Equal(t, "data_vol", conf.ServerGroups["servers"].Volumes[0].Name)
	assert.Equal(t, "/opt/nomad_server_data", conf.ServerGroups["servers"].Volumes[0].Path)
	assert.Equal(t, []string{"consul"}, conf.ServerGroups["servers"].Aliases)

}

func TestLoadProviders(t *testing.T) {
	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, conf)
	fmt.Println(conf)

	assert.Len(t, conf.Providers, 2)
}

func TestLoadProviderConfig(t *testing.T) {
	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, conf)

	provider, err := LoadTFVarsConfig(*conf)
	require.NoError(t, err)
	assert.NotNil(t, provider)
	hetzner := provider.ProviderConfig.(HetznerSettings)

	expected := HetznerSettings{
		SSHKeys:          []string{"123456"},
		Location:         "nbg1",
		LoadBalancerType: "lb11",
		ResourceNames: HetznerResourceNames{
			BaseServerName: "nomad-srv",
			FirewallName:   "dev_firewall",
			NetworkName:    "dev_network",
		},
	}

	assert.Equal(t, expected, hetzner)
	assert.Equal(t, []string{"85.4.84.201/32"}, conf.CloudProviderConfig.AllowedIPs)
}

func TestLoadTFExecVars(t *testing.T) {
	theVar := LoadTFExecVars()

	assert.NotNil(t, theVar)
}
