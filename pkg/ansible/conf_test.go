package ansible

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenPaasDev/core/pkg/conf"
	"github.com/OpenPaasDev/core/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestGenerateInventory(t *testing.T) {
	config, err := conf.Load(filepath.Join("..", "testdata", "config.yaml"))
	assert.NoError(t, err)

	folder := util.RandString(8)
	config.BaseDir = folder
	err = os.MkdirAll(folder, 0700)
	assert.NoError(t, err)
	defer func() {
		e := os.RemoveAll(filepath.Join(folder))
		assert.NoError(t, e)
	}()
	src := filepath.Join("testdata", "inventory-output.json")
	dest := filepath.Join(folder, "inventory-output.json")

	bytesRead, err := os.ReadFile(filepath.Clean(src))

	assert.NoError(t, err)

	err = os.WriteFile(dest, bytesRead, 0600)

	assert.NoError(t, err)

	inventory, err := GenerateInventory(config)
	assert.NoError(t, err)

	assert.FileExists(t, filepath.Join(folder, "inventory"))

	assert.NotEmpty(t, inventory.GetAllPrivateHosts())

	_, err = LoadInventory(filepath.Join(folder, "inventory"))
	assert.NoError(t, err)

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
	assert.Equal(t, host.Mounts[0].MountPath, "/mnt/HC_Volume_29747974")
	assert.Equal(t, host.Mounts[0].Path, "/opt/nomad_server_data")
	assert.Equal(t, host.Mounts[0].Owner, "www-data")

	assert.Equal(t, host.HostName, "venue-servers-1")
	assert.Equal(t, host.PrivateIP, "10.0.1.2")
	assert.Equal(t, host.PublicIP, "138.201.186.150")

}
