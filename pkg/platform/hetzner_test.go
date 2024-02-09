package platform

import (
	"errors"
	"testing"

	"github.com/OpenPaasDev/openpaas/pkg/conf"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestRunPreparationLogic_UpdatesConfigAndKeys(t *testing.T) {
	ctx := context.Background()
	cnf, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	githubKeys := []string{"ssh-rsa AAAAB3Nza1234567890AAAAB3Nza user@example.com", "ssh-rsa BBBBA3Nza1234567890BBBBA3Nza user@main.com"}
	hetznerKeys := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 1, Name: GHKeyPrefix + "old_key", Fingerprint: "abc"},
		{ID: 2, Name: GHKeyPrefix + "old_key2", Fingerprint: "abcd"},
	}
	hetznerKeys2 := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 3, Name: GHKeyPrefix + "new_key", Fingerprint: "abc"},
		{ID: 4, Name: GHKeyPrefix + "new_key2", Fingerprint: "abcd"},
	}
	var erasedId []uint32
	var uploadedKeys []string
	var uploadedNames []string

	called := false
	getHetznerKeys := func() ([]HetznerSSHKey, error) {
		if !called {
			called = true
			return hetznerKeys, nil
		} else {
			return hetznerKeys2, nil
		}
	}
	getGithubKeys := func(context.Context, string) ([]string, error) {
		return githubKeys, nil
	}
	eraseHetznerKey := func(k HetznerSSHKey) error {
		erasedId = append(erasedId, k.ID)
		return nil
	}
	uploadHetznerKey := func(key string, name string) error {
		uploadedKeys = append(uploadedKeys, key)
		uploadedNames = append(uploadedNames, name)
		return nil
	}
	getPublicIp := func(context.Context) (string, error) {
		return "1.1.1.1", nil
	}

	err = runPreparationLogic(ctx, cnf, getHetznerKeys, getGithubKeys, eraseHetznerKey, uploadHetznerKey, getPublicIp)
	require.NoError(t, err)
	assert.Equal(t, []uint32{1, 2}, erasedId)
	assert.Equal(t, githubKeys, uploadedKeys)
	assert.Equal(t, []string{"gh-key-xample.com", "gh-key-r@main.com"}, uploadedNames)
	assert.Equal(t, []string{"3", "4"}, cnf.CloudProviderConfig.ProviderSettings["ssh_keys"])
	assert.Equal(t, []string{"85.4.84.201/32", "1.1.1.1/32"}, cnf.CloudProviderConfig.AllowedIPs)
}

func TestRunPreparationLogic_PropagatesReadKeysErrors(t *testing.T) {
	ctx := context.Background()
	cnf, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	githubKeys := []string{"ssh-rsa AAAAB3Nza1234567890AAAAB3Nza user@example.com", "ssh-rsa BBBBA3Nza1234567890BBBBA3Nza user@main.com"}
	var erasedId []uint32
	var uploadedKeys []string
	var uploadedNames []string

	getHetznerKeys := func() ([]HetznerSSHKey, error) {
		return nil, errors.New("Fail")
	}
	getGithubKeys := func(context.Context, string) ([]string, error) {
		return githubKeys, nil
	}
	eraseHetznerKey := func(k HetznerSSHKey) error {
		erasedId = append(erasedId, k.ID)
		return nil
	}
	uploadHetznerKey := func(key string, name string) error {
		uploadedKeys = append(uploadedKeys, key)
		uploadedNames = append(uploadedNames, name)
		return nil
	}
	getPublicIp := func(context.Context) (string, error) {
		return "1.1.1.1", nil
	}

	err = runPreparationLogic(ctx, cnf, getHetznerKeys, getGithubKeys, eraseHetznerKey, uploadHetznerKey, getPublicIp)
	require.Error(t, err)
}

func TestRunPreparationLogic_PropagatesReadGithubErrors(t *testing.T) {
	ctx := context.Background()
	cnf, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	expected := errors.New("Fail")

	hetznerKeys := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 1, Name: GHKeyPrefix + "old_key", Fingerprint: "abc"},
		{ID: 2, Name: GHKeyPrefix + "old_key2", Fingerprint: "abcd"},
	}
	hetznerKeys2 := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 3, Name: GHKeyPrefix + "new_key", Fingerprint: "abc"},
		{ID: 4, Name: GHKeyPrefix + "new_key2", Fingerprint: "abcd"},
	}
	var erasedId []uint32
	var uploadedKeys []string
	var uploadedNames []string

	called := false
	getHetznerKeys := func() ([]HetznerSSHKey, error) {
		if !called {
			called = true
			return hetznerKeys, nil
		} else {
			return hetznerKeys2, nil
		}
	}
	getGithubKeys := func(context.Context, string) ([]string, error) {
		return nil, expected
	}
	eraseHetznerKey := func(k HetznerSSHKey) error {
		erasedId = append(erasedId, k.ID)
		return nil
	}
	uploadHetznerKey := func(key string, name string) error {
		uploadedKeys = append(uploadedKeys, key)
		uploadedNames = append(uploadedNames, name)
		return nil
	}
	getPublicIp := func(context.Context) (string, error) {
		return "1.1.1.1", nil
	}
	err = runPreparationLogic(ctx, cnf, getHetznerKeys, getGithubKeys, eraseHetznerKey, uploadHetznerKey, getPublicIp)
	require.ErrorIs(t, expected, err)
}

func TestRunPreparationLogic_PropagatesEraseKeyErrors(t *testing.T) {
	ctx := context.Background()
	cnf, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	expected := errors.New("Fail")
	githubKeys := []string{"ssh-rsa AAAAB3Nza1234567890AAAAB3Nza user@example.com", "ssh-rsa BBBBA3Nza1234567890BBBBA3Nza user@main.com"}
	hetznerKeys := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 1, Name: GHKeyPrefix + "old_key", Fingerprint: "abc"},
		{ID: 2, Name: GHKeyPrefix + "old_key2", Fingerprint: "abcd"},
	}
	hetznerKeys2 := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 3, Name: GHKeyPrefix + "new_key", Fingerprint: "abc"},
		{ID: 4, Name: GHKeyPrefix + "new_key2", Fingerprint: "abcd"},
	}
	var uploadedKeys []string
	var uploadedNames []string

	called := false
	getHetznerKeys := func() ([]HetznerSSHKey, error) {
		if !called {
			called = true
			return hetznerKeys, nil
		} else {
			return hetznerKeys2, nil
		}
	}
	getGithubKeys := func(context.Context, string) ([]string, error) {
		return githubKeys, nil
	}
	eraseHetznerKey := func(k HetznerSSHKey) error {
		return expected
	}
	uploadHetznerKey := func(key string, name string) error {
		uploadedKeys = append(uploadedKeys, key)
		uploadedNames = append(uploadedNames, name)
		return nil
	}
	getPublicIp := func(context.Context) (string, error) {
		return "1.1.1.1", nil
	}
	err = runPreparationLogic(ctx, cnf, getHetznerKeys, getGithubKeys, eraseHetznerKey, uploadHetznerKey, getPublicIp)
	require.ErrorIs(t, expected, err)
}

func TestRunPreparationLogic_PropagatesUploadKeyErrors(t *testing.T) {
	ctx := context.Background()
	cnf, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	expected := errors.New("Fail")
	githubKeys := []string{"ssh-rsa AAAAB3Nza1234567890AAAAB3Nza user@example.com", "ssh-rsa BBBBA3Nza1234567890BBBBA3Nza user@main.com"}
	hetznerKeys := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 1, Name: GHKeyPrefix + "old_key", Fingerprint: "abc"},
		{ID: 2, Name: GHKeyPrefix + "old_key2", Fingerprint: "abcd"},
	}
	hetznerKeys2 := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 3, Name: GHKeyPrefix + "new_key", Fingerprint: "abc"},
		{ID: 4, Name: GHKeyPrefix + "new_key2", Fingerprint: "abcd"},
	}
	var erasedId []uint32

	called := false
	getHetznerKeys := func() ([]HetznerSSHKey, error) {
		if !called {
			called = true
			return hetznerKeys, nil
		} else {
			return hetznerKeys2, nil
		}
	}
	getGithubKeys := func(context.Context, string) ([]string, error) {
		return githubKeys, nil
	}
	eraseHetznerKey := func(k HetznerSSHKey) error {
		erasedId = append(erasedId, k.ID)
		return nil
	}
	uploadHetznerKey := func(key string, name string) error {
		return expected
	}
	getPublicIp := func(context.Context) (string, error) {
		return "1.1.1.1", nil
	}

	err = runPreparationLogic(ctx, cnf, getHetznerKeys, getGithubKeys, eraseHetznerKey, uploadHetznerKey, getPublicIp)
	require.ErrorIs(t, expected, err)
}

func TestRunPreparationLogic_PropagatesGetIPErrors(t *testing.T) {
	ctx := context.Background()
	cnf, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	expected := errors.New("Fail")
	githubKeys := []string{"ssh-rsa AAAAB3Nza1234567890AAAAB3Nza user@example.com", "ssh-rsa BBBBA3Nza1234567890BBBBA3Nza user@main.com"}
	hetznerKeys := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 1, Name: GHKeyPrefix + "old_key", Fingerprint: "abc"},
		{ID: 2, Name: GHKeyPrefix + "old_key2", Fingerprint: "abcd"},
	}
	hetznerKeys2 := []HetznerSSHKey{
		{ID: 0, Name: "old_key", Fingerprint: "abc"},
		{ID: 3, Name: GHKeyPrefix + "new_key", Fingerprint: "abc"},
		{ID: 4, Name: GHKeyPrefix + "new_key2", Fingerprint: "abcd"},
	}
	var erasedId []uint32
	var uploadedKeys []string
	var uploadedNames []string

	called := false
	getHetznerKeys := func() ([]HetznerSSHKey, error) {
		if !called {
			called = true
			return hetznerKeys, nil
		} else {
			return hetznerKeys2, nil
		}
	}
	getGithubKeys := func(context.Context, string) ([]string, error) {
		return githubKeys, nil
	}
	eraseHetznerKey := func(k HetznerSSHKey) error {
		erasedId = append(erasedId, k.ID)
		return nil
	}
	uploadHetznerKey := func(key string, name string) error {
		uploadedKeys = append(uploadedKeys, key)
		uploadedNames = append(uploadedNames, name)
		return nil
	}
	getPublicIp := func(context.Context) (string, error) {
		return "", expected
	}

	err = runPreparationLogic(ctx, cnf, getHetznerKeys, getGithubKeys, eraseHetznerKey, uploadHetznerKey, getPublicIp)
	require.ErrorIs(t, expected, err)
}

func TestFetchGitHubKeys(t *testing.T) {
	ctx := context.Background()
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	payload := `ssh-rsa AAAAB3Nza... user@example.com

ssh-rsa BBBBA3Nza... user@main.com
`

	httpmock.RegisterResponder("GET", "https://github.com/wfaler.keys",
		httpmock.NewStringResponder(200, payload))

	keys, err := fetchGitHubKeys(ctx, "wfaler")
	require.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.Contains(t, keys[0], "ssh-rsa AAAAB3Nza")
	assert.Contains(t, keys[1], "ssh-rsa BBBBA3Nza")
}

func TestFetchGitHubKeys_NotFound(t *testing.T) {
	ctx := context.Background()
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	httpmock.RegisterResponder("GET", "https://github.com/wfaler.keys",
		httpmock.NewStringResponder(404, "Not Found"))

	_, err := fetchGitHubKeys(ctx, "wfaler")
	require.Error(t, err)
}

func TestFetchGitHubKeys_Error(t *testing.T) {
	ctx := context.Background()
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	httpmock.RegisterResponder("GET", "https://github.com/wfaler.keys",
		httpmock.NewErrorResponder(errors.New("internal server error")))

	_, err := fetchGitHubKeys(ctx, "wfaler")
	require.Error(t, err)
}
