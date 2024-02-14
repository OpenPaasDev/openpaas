package conf

import (
	"fmt"
	"os"
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
	assert.Equal(t, "hetzner", conf.CloudProviderConfig.Provider)
	assert.Equal(t, []string{"wfaler"}, conf.CloudProviderConfig.GithubIds)
	assert.Equal(t, "s3", conf.TfState.Backend)

	// we don't have env vars set, so access and secret keys show the string provided
	expected := map[string]string{"endpoint": "endpoint_to_s3_compatible_storage",
		"bucket":     "bucket_name",
		"region":     "auto",
		"access_key": "env_var_access_key",
		"secret_key": "env_var_secret_key",
	}
	assert.Equal(t, expected, conf.TfState.Config)

	assert.Len(t, conf.ServerGroups, 2)
	assert.Equal(t, "cpx31", conf.ServerGroups["clients"].InstanceType)
	assert.Equal(t, 2, conf.ServerGroups["clients"].Num)
	assert.Equal(t, 20, conf.ServerGroups["servers"].Volumes[0].Size)
	assert.Equal(t, "data_vol", conf.ServerGroups["servers"].Volumes[0].Name)
	assert.Equal(t, "/opt/nomad_server_data", conf.ServerGroups["servers"].Volumes[0].Path)
	assert.Equal(t, []string{"consul"}, conf.ServerGroups["servers"].Aliases)

}

func TestLoadConfigWithEnvVars(t *testing.T) {
	err := os.Setenv("env_var_access_key", "access_key")
	require.NoError(t, err)
	err = os.Setenv("env_var_secret_key", "secret_key")
	require.NoError(t, err)

	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, conf)
	fmt.Println(conf)

	assert.Equal(t, "config", conf.BaseDir)
	assert.Equal(t, "hetzner", conf.DC)
	assert.Equal(t, "hetzner", conf.CloudProviderConfig.Provider)
	assert.Equal(t, []string{"wfaler"}, conf.CloudProviderConfig.GithubIds)
	assert.Equal(t, "s3", conf.TfState.Backend)

	// we have env vars set, so access and secret keys are replaced accordingly
	expected := map[string]string{"endpoint": "endpoint_to_s3_compatible_storage",
		"bucket":     "bucket_name",
		"region":     "auto",
		"access_key": "access_key",
		"secret_key": "secret_key",
	}
	assert.Equal(t, expected, conf.TfState.Config)

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

	assert.Len(t, conf.Providers, 1)
}

func TestLoadTFExecVars(t *testing.T) {
	theVar := LoadTFExecVars()

	assert.NotNil(t, theVar)
}
