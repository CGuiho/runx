package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeCmd(t *testing.T) {
	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "runx.yaml")

	content := `version: "2"
namespace: "test"
description: "Test manifest"

commands:
  - uid: "test:hello"
    id: "hello"
    description: "Print hello"
    run: "echo Hello"
`
	err := os.WriteFile(yamlPath, []byte(content), 0644)
	require.NoError(t, err)

	t.Run("valid command describe text format", func(t *testing.T) {
		cfgFile = yamlPath
		err := describeCmd.RunE(describeCmd, []string{"test:hello"})
		assert.NoError(t, err)
	})

	t.Run("unknown command describe", func(t *testing.T) {
		cfgFile = yamlPath
		err := describeCmd.RunE(describeCmd, []string{"unknown:cmd"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found in manifest")
	})
}
