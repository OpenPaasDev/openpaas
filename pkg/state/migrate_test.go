package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Migrate(t *testing.T) {
	err := os.MkdirAll("testdata", 0650)
	defer func() {
		e := os.Remove(filepath.Join("testdata", "state.db"))
		require.NoError(t, e)
	}()
	require.NoError(t, err)
	err = Migrate("testdata")
	require.NoError(t, err)

}
