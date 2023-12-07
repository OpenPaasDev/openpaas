package ansible

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnsible(t *testing.T) {
	currentUser, err := user.Current()
	require.NoError(t, err)

	ansibleClient := NewClient(filepath.Join("testdata", "inventory"), filepath.Join("testdata", "secrets"), currentUser.Username, filepath.Join("testdata", "secrets"))
	err = ansibleClient.Run(filepath.Join("testdata", "ansible.yml"))
	require.NoError(t, err)
}
