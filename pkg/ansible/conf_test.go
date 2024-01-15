package ansible

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"github.com/OpenPaasDev/openpaas/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateInventory(t *testing.T) {
	config, err := conf.Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)

	folder := util.RandString(8)
	config.BaseDir = folder
	err = os.MkdirAll(folder, 0700)
	require.NoError(t, err)
	defer func() {
		e := os.RemoveAll(filepath.Join(folder))
		require.NoError(t, e)
	}()
	src := filepath.Join("testdata", "inventory-output.json")
	dest := filepath.Join(folder, "inventory-output.json")

	bytesRead, err := os.ReadFile(filepath.Clean(src))

	require.NoError(t, err)

	err = os.WriteFile(dest, bytesRead, 0600)

	require.NoError(t, err)

	inventory, err := GenerateInventory(config)
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(folder, "inventory"))

	assert.NotEmpty(t, inventory.GetAllPrivateHosts())

	_, err = LoadInventory(filepath.Join(folder, "inventory"))
	require.NoError(t, err)

	var host AnsibleHost
	for _, v := range inventory.All.Children["servers"].Hosts {
		if v.ID == "30421332" {
			host = v
			break
		}
	}

	assert.Len(t, inventory.All.Children["clients"].Hosts, 2)
	assert.Len(t, inventory.All.Children["servers"].Hosts, 3)
	assert.Len(t, inventory.All.Children["consul"].Hosts, 3)
	assert.Len(t, host.Mounts, 1)
	assert.Equal(t, "/mnt/HC_Volume_29747974", host.Mounts[0].MountPath)
	assert.Equal(t, "/opt/nomad_server_data", host.Mounts[0].Path)
	assert.Equal(t, "www-data", host.Mounts[0].Owner)

	assert.Equal(t, "venue-servers-1", host.HostName)
	assert.Equal(t, "10.0.1.2", host.PrivateIP)
	assert.Equal(t, "138.201.186.150", host.PublicIP)

}
