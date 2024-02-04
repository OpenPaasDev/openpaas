package platform

import (
	"path/filepath"
	"testing"

	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func Hetzner_Cleanup_NotImplemented(t *testing.T) {
	ctx := context.Background()
	conf, err := conf.Load(filepath.Join("..", "testdata", "config.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, conf)
	hetzner := Hetzner{}

	err = hetzner.Cleanup(ctx, conf)
	require.NoError(t, err)
}

func FindChangesToMake_NoKeysInHetzner(t *testing.T) {
	var keysInHetzner []HetznerSSHKey
	var githubKeys []conf.GithubKey

	githubKeys = []conf.GithubKey{
		{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIH0p7vluS7ldSFDBYx9ZXVQcsJdWIoSTVvqhcakKDQ34 user@example.com", Fingerprint: "ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22"},
		{PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCas+tCHuRci1xIHwkBFvrq/dDm0l3PSBiD+Pm7SOqq23Qg+ZUANcVNSgot0W/NmEHy/9rA78Ps+jHwrwPrliw6TNLYC82LueJHPo1dGioprmKVKn9efYVuFDUubRPr+/CAZXARUOSqMon7xiaxAy51qhIpWrLxcs2HP/G3IW8kVxuwdztyT7D+tz5tfCvH98PlXf6MjWudL8bbTAWU7OEpUVia2pcUlAOXkkOi0ANrx4Ieovmhw7G8/AC0Rn+g3hSf1A45RsODVFq9BezunWbjcNicwV2++/CFpE5fuXT6pRgrBfXgI3P4BVRMxaG4CXLl6uUPrg/8oYoz/uJtxzEwv767YeNICi9RfXjIg0hQLoZfIAZCxYIVZw9A91ZIG8+IP276gG1kHfyMfs2W95dK6Uy/fdF6G3p3PLFUtjSK6dZwQeO4IzltrxujQ26kgMFDYMmD7lDDI3JzLqXy969MpMd60iamXDpgxQ3okZpa7sd5TgtEH+aA8ia58bhSOwE= user@main.com", Fingerprint: "a6:1e:2b:44:9c:84:e2:1c:84:9c:6d:2d:ed:72:ad:16"},
	}

	erase, upload := findChangesToMake(keysInHetzner, githubKeys)
	assert.Equal(t, []string{}, erase)
	assert.Equal(t, githubKeys, upload)
}

func FindChangesToMake_NoKeysInGithub(t *testing.T) {
	var keysInHetzner []HetznerSSHKey
	var githubKeys []conf.GithubKey

	keysInHetzner = []HetznerSSHKey{
		{
			ID:          0,
			Name:        "key1",
			Fingerprint: "ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22",
		},
	}

	erase, upload := findChangesToMake(keysInHetzner, githubKeys)
	assert.Equal(t, []string{"ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22"}, erase)
	assert.Equal(t, githubKeys, upload)
}

func FindChangesToMake_NoCommonKeys(t *testing.T) {
	var keysInHetzner []HetznerSSHKey
	var githubKeys []conf.GithubKey

	keysInHetzner = []HetznerSSHKey{
		{
			ID:          0,
			Name:        "key1",
			Fingerprint: "ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22",
		},
	}

	githubKeys = []conf.GithubKey{
		{PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCas+tCHuRci1xIHwkBFvrq/dDm0l3PSBiD+Pm7SOqq23Qg+ZUANcVNSgot0W/NmEHy/9rA78Ps+jHwrwPrliw6TNLYC82LueJHPo1dGioprmKVKn9efYVuFDUubRPr+/CAZXARUOSqMon7xiaxAy51qhIpWrLxcs2HP/G3IW8kVxuwdztyT7D+tz5tfCvH98PlXf6MjWudL8bbTAWU7OEpUVia2pcUlAOXkkOi0ANrx4Ieovmhw7G8/AC0Rn+g3hSf1A45RsODVFq9BezunWbjcNicwV2++/CFpE5fuXT6pRgrBfXgI3P4BVRMxaG4CXLl6uUPrg/8oYoz/uJtxzEwv767YeNICi9RfXjIg0hQLoZfIAZCxYIVZw9A91ZIG8+IP276gG1kHfyMfs2W95dK6Uy/fdF6G3p3PLFUtjSK6dZwQeO4IzltrxujQ26kgMFDYMmD7lDDI3JzLqXy969MpMd60iamXDpgxQ3okZpa7sd5TgtEH+aA8ia58bhSOwE= user@main.com", Fingerprint: "a6:1e:2b:44:9c:84:e2:1c:84:9c:6d:2d:ed:72:ad:16"},
	}

	erase, upload := findChangesToMake(keysInHetzner, githubKeys)
	assert.Equal(t, []string{"ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22"}, erase)
	assert.Equal(t, githubKeys, upload)
}

func FindChangesToMake_FindCommonKeys(t *testing.T) {
	var keysInHetzner []HetznerSSHKey
	var githubKeys []conf.GithubKey

	keysInHetzner = []HetznerSSHKey{
		{
			ID:          0,
			Name:        "key1",
			Fingerprint: "ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22",
		},
		{
			ID:          1,
			Name:        "key_to_erase",
			Fingerprint: "bb:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22",
		},
	}

	githubKeys = []conf.GithubKey{
		{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIH0p7vluS7ldSFDBYx9ZXVQcsJdWIoSTVvqhcakKDQ34 user@example.com", Fingerprint: "ae:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22"},
		{PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCas+tCHuRci1xIHwkBFvrq/dDm0l3PSBiD+Pm7SOqq23Qg+ZUANcVNSgot0W/NmEHy/9rA78Ps+jHwrwPrliw6TNLYC82LueJHPo1dGioprmKVKn9efYVuFDUubRPr+/CAZXARUOSqMon7xiaxAy51qhIpWrLxcs2HP/G3IW8kVxuwdztyT7D+tz5tfCvH98PlXf6MjWudL8bbTAWU7OEpUVia2pcUlAOXkkOi0ANrx4Ieovmhw7G8/AC0Rn+g3hSf1A45RsODVFq9BezunWbjcNicwV2++/CFpE5fuXT6pRgrBfXgI3P4BVRMxaG4CXLl6uUPrg/8oYoz/uJtxzEwv767YeNICi9RfXjIg0hQLoZfIAZCxYIVZw9A91ZIG8+IP276gG1kHfyMfs2W95dK6Uy/fdF6G3p3PLFUtjSK6dZwQeO4IzltrxujQ26kgMFDYMmD7lDDI3JzLqXy969MpMd60iamXDpgxQ3okZpa7sd5TgtEH+aA8ia58bhSOwE= user@main.com", Fingerprint: "a6:1e:2b:44:9c:84:e2:1c:84:9c:6d:2d:ed:72:ad:16"},
	}

	erase, upload := findChangesToMake(keysInHetzner, githubKeys)
	assert.Equal(t, []string{"bb:dc:ab:c1:b1:b0:21:2b:8a:06:77:ae:9c:9b:4b:22"}, erase)
	assert.Equal(t, []conf.GithubKey{
		{PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCas+tCHuRci1xIHwkBFvrq/dDm0l3PSBiD+Pm7SOqq23Qg+ZUANcVNSgot0W/NmEHy/9rA78Ps+jHwrwPrliw6TNLYC82LueJHPo1dGioprmKVKn9efYVuFDUubRPr+/CAZXARUOSqMon7xiaxAy51qhIpWrLxcs2HP/G3IW8kVxuwdztyT7D+tz5tfCvH98PlXf6MjWudL8bbTAWU7OEpUVia2pcUlAOXkkOi0ANrx4Ieovmhw7G8/AC0Rn+g3hSf1A45RsODVFq9BezunWbjcNicwV2++/CFpE5fuXT6pRgrBfXgI3P4BVRMxaG4CXLl6uUPrg/8oYoz/uJtxzEwv767YeNICi9RfXjIg0hQLoZfIAZCxYIVZw9A91ZIG8+IP276gG1kHfyMfs2W95dK6Uy/fdF6G3p3PLFUtjSK6dZwQeO4IzltrxujQ26kgMFDYMmD7lDDI3JzLqXy969MpMd60iamXDpgxQ3okZpa7sd5TgtEH+aA8ia58bhSOwE= user@main.com", Fingerprint: "a6:1e:2b:44:9c:84:e2:1c:84:9c:6d:2d:ed:72:ad:16"},
	}, upload)
}
