package platform

import (
	"errors"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

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
