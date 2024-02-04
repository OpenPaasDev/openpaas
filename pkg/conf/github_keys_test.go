package conf

import (
	"errors"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func UpdateConfigWithGithubKeys_NoKeysInConfig(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	payload := `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIH0p7vluS7ldSFDBYx9ZXVQcsJdWIoSTVvqhcakKDQ34 user@example.com
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCas+tCHuRci1xIHwkBFvrq/dDm0l3PSBiD+Pm7SOqq23Qg+ZUANcVNSgot0W/NmEHy/9rA78Ps+jHwrwPrliw6TNLYC82LueJHPo1dGioprmKVKn9efYVuFDUubRPr+/CAZXARUOSqMon7xiaxAy51qhIpWrLxcs2HP/G3IW8kVxuwdztyT7D+tz5tfCvH98PlXf6MjWudL8bbTAWU7OEpUVia2pcUlAOXkkOi0ANrx4Ieovmhw7G8/AC0Rn+g3hSf1A45RsODVFq9BezunWbjcNicwV2++/CFpE5fuXT6pRgrBfXgI3P4BVRMxaG4CXLl6uUPrg/8oYoz/uJtxzEwv767YeNICi9RfXjIg0hQLoZfIAZCxYIVZw9A91ZIG8+IP276gG1kHfyMfs2W95dK6Uy/fdF6G3p3PLFUtjSK6dZwQeO4IzltrxujQ26kgMFDYMmD7lDDI3JzLqXy969MpMd60iamXDpgxQ3okZpa7sd5TgtEH+aA8ia58bhSOwE= user@main.com`

	httpmock.RegisterResponder("GET", "https://github.com/exampleUser.keys",
		httpmock.NewStringResponder(200, payload))

	// initial config without ssh_keys entry defined
	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, conf)

	// remove ssh keys defined
	conf.CloudProviderConfig.ProviderSettings["ssh_keys"] = []string{}
	expected := []string{"ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22", "a6:1e:2b:44:9c:84:e2:1c:84:9c:6d:2d:ed:72:ad:16"}

	conf, err = UpdateConfigWithGithubKeys(conf)
	assert.NoError(t, err)
	assert.Contains(t, conf.CloudProviderConfig.ProviderSettings["ssh_keys"], expected)
}

func UpdateConfigWithGithubKeys_KeysInConfig(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	payload := `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIH0p7vluS7ldSFDBYx9ZXVQcsJdWIoSTVvqhcakKDQ34 user@example.com
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCas+tCHuRci1xIHwkBFvrq/dDm0l3PSBiD+Pm7SOqq23Qg+ZUANcVNSgot0W/NmEHy/9rA78Ps+jHwrwPrliw6TNLYC82LueJHPo1dGioprmKVKn9efYVuFDUubRPr+/CAZXARUOSqMon7xiaxAy51qhIpWrLxcs2HP/G3IW8kVxuwdztyT7D+tz5tfCvH98PlXf6MjWudL8bbTAWU7OEpUVia2pcUlAOXkkOi0ANrx4Ieovmhw7G8/AC0Rn+g3hSf1A45RsODVFq9BezunWbjcNicwV2++/CFpE5fuXT6pRgrBfXgI3P4BVRMxaG4CXLl6uUPrg/8oYoz/uJtxzEwv767YeNICi9RfXjIg0hQLoZfIAZCxYIVZw9A91ZIG8+IP276gG1kHfyMfs2W95dK6Uy/fdF6G3p3PLFUtjSK6dZwQeO4IzltrxujQ26kgMFDYMmD7lDDI3JzLqXy969MpMd60iamXDpgxQ3okZpa7sd5TgtEH+aA8ia58bhSOwE= user@main.com`

	httpmock.RegisterResponder("GET", "https://github.com/exampleUser.keys",
		httpmock.NewStringResponder(200, payload))

	// initial config without ssh_keys entry defined
	conf, err := Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, conf)

	// we expect the github keys to be added to the ssh keys defined in the config
	configKeys := conf.CloudProviderConfig.ProviderSettings["ssh_keys"].([]string)
	expected := append(configKeys, []string{"ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22", "a6:1e:2b:44:9c:84:e2:1c:84:9c:6d:2d:ed:72:ad:16"}...)

	conf, err = UpdateConfigWithGithubKeys(conf)
	assert.NoError(t, err)
	assert.Contains(t, conf.CloudProviderConfig.ProviderSettings["ssh_keys"], expected)
}

func TestFetchGitHubKeys(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	payload := `ssh-rsa AAAAB3Nza... user@example.com
ssh-rsa BBBBA3Nza... user@main.com`

	httpmock.RegisterResponder("GET", "https://github.com/exampleUser.keys",
		httpmock.NewStringResponder(200, payload))

	keys, err := fetchGitHubKeys("exampleUser")
	assert.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.Contains(t, keys[0], "ssh-rsa AAAAB3Nza")
	assert.Contains(t, keys[1], "ssh-rsa BBBBA3Nza")
}

func TestFetchGitHubKeys_NotFound(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	httpmock.RegisterResponder("GET", "https://github.com/exampleUser.keys",
		httpmock.NewStringResponder(404, "Not Found"))

	_, err := fetchGitHubKeys("exampleUser")
	assert.Error(t, err)
}

func TestFetchGitHubKeys_Error(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	httpmock.RegisterResponder("GET", "https://github.com/exampleUser.keys",
		httpmock.NewErrorResponder(errors.New("internal server error")))

	_, err := fetchGitHubKeys("exampleUser")
	assert.Error(t, err)
}

func TestGenerateFingerprint(t *testing.T) {
	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCas+tCHuRci1xIHwkBFvrq/dDm0l3PSBiD+Pm7SOqq23Qg+ZUANcVNSgot0W/NmEHy/9rA78Ps+jHwrwPrliw6TNLYC82LueJHPo1dGioprmKVKn9efYVuFDUubRPr+/CAZXARUOSqMon7xiaxAy51qhIpWrLxcs2HP/G3IW8kVxuwdztyT7D+tz5tfCvH98PlXf6MjWudL8bbTAWU7OEpUVia2pcUlAOXkkOi0ANrx4Ieovmhw7G8/AC0Rn+g3hSf1A45RsODVFq9BezunWbjcNicwV2++/CFpE5fuXT6pRgrBfXgI3P4BVRMxaG4CXLl6uUPrg/8oYoz/uJtxzEwv767YeNICi9RfXjIg0hQLoZfIAZCxYIVZw9A91ZIG8+IP276gG1kHfyMfs2W95dK6Uy/fdF6G3p3PLFUtjSK6dZwQeO4IzltrxujQ26kgMFDYMmD7lDDI3JzLqXy969MpMd60iamXDpgxQ3okZpa7sd5TgtEH+aA8ia58bhSOwE="
	expectedFingerprint := "a6:1e:2b:44:9c:84:e2:1c:84:9c:6d:2d:ed:72:ad:16"
	fingerprint, err := generateFingerprint(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedFingerprint, fingerprint)

	// check ed keys
	key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIH0p7vluS7ldSFDBYx9ZXVQcsJdWIoSTVvqhcakKDQ34"
	expectedFingerprint = "ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22"
	fingerprint, err = generateFingerprint(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedFingerprint, fingerprint)
}
