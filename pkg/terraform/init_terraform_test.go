package terraform

import (
	"context"
	"os"
	"testing"

	"github.com/OpenPaasDev/openpaas/pkg/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Init_Terraform(t *testing.T) {
	ctx := context.Background()
	tf, err := InitTf(ctx, ".", "1.7.3", os.Stdin, os.Stderr)
	require.NoError(t, err)
	v, _, err := tf.Version(ctx, false)
	require.NoError(t, err)
	assert.NotNil(t, v)
}

func Test_Get_Terraform_Executable_Path(t *testing.T) {
	ctx := context.Background()
	execPath, err := GetTerraformExecutablePath(ctx, "1.7.3")
	require.NoError(t, err)
	err = util.RunCmd(execPath, "--version")
	require.NoError(t, err)
}
