package executor_test

import (
	"testing"

	"github.com/CGuiho/runx/pkg/executor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDryRunExecution(t *testing.T) {
	opts := executor.ExecutionOptions{
		Command: "echo Hello",
		Args:    []string{"--foo", "bar"},
		DryRun:  true,
	}

	res, err := executor.ExecuteCommand(opts)
	require.NoError(t, err)
	assert.Equal(t, 0, res.ExitCode)
}
