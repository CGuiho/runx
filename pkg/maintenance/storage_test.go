package maintenance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_BasicOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-storage-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "sub", "test.txt")

	assert.False(t, PathExists(filePath))

	err = WriteTextFile(filePath, "hello world")
	require.NoError(t, err)
	assert.True(t, PathExists(filePath))

	content, err := ReadTextIfExists(filePath)
	require.NoError(t, err)
	assert.Equal(t, "hello world", content)

	// Atomic write test
	err = WriteTextFileAtomic(filePath, "updated content")
	require.NoError(t, err)

	updatedContent, err := ReadTextIfExists(filePath)
	require.NoError(t, err)
	assert.Equal(t, "updated content", updatedContent)

	// Remove test
	err = RemovePath(filePath)
	require.NoError(t, err)
	assert.False(t, PathExists(filePath))
}
