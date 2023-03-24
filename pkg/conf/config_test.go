package conf

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	fmt.Println(conf)

	assert.Equal(t, "config", conf.BaseDir)
	assert.Equal(t, "hetzner", conf.DC)
	assert.Equal(t, "ens10", conf.CloudProviderConfig.NetworkInterface)
	assert.Equal(t, "root", conf.CloudProviderConfig.User)
	assert.Equal(t, "hetzner", conf.CloudProviderConfig.Provider)

	assert.Equal(t, 2, len(conf.ServerGroups))
	assert.Equal(t, "cpx31", conf.ServerGroups["clients"].InstanceType)
	assert.Equal(t, 2, conf.ServerGroups["clients"].Num)
	assert.Equal(t, 20, conf.ServerGroups["servers"].Volumes[0].Size)
	assert.Equal(t, "data_vol", conf.ServerGroups["servers"].Volumes[0].Name)
	assert.Equal(t, "/opt/nomad_server_data", conf.ServerGroups["servers"].Volumes[0].Path)

}

func TestLoadProviderConfig(t *testing.T) {
	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	provider, err := LoadTFVarsConfig(*conf)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	hetzner := provider.ProviderConfig.(HetznerSettings)

	expected := HetznerSettings{
		SSHKeys:  []string{"wille.faler@gmail.com"},
		Location: "nbg1",
		LoadBalancer: LoadBalancerSettings{
			Enabled:      true,
			InstanceType: "lb11",
			ServerGroup:  "clients",
		},
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
