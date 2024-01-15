package terraform

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Init_Terraform(t *testing.T) {
	ctx := context.Background()
	tf, err := InitTf(ctx, ".", os.Stdin, os.Stderr)
	require.NoError(t, err)
	v, _, err := tf.Version(ctx, false)
	require.NoError(t, err)
	assert.NotNil(t, v)
}
