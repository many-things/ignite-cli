package mdrun_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/mdrun"
)

func TestRun(t *testing.T) {
	err := mdrun.Run("testdata")

	require.NoError(t, err)
}
