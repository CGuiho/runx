package updater

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockFileOps struct {
	RenameFunc         func(from, to string) error
	RemoveFunc         func(path string) error
	MakeExecutableFunc func(path string) error
}

func (m MockFileOps) Rename(from, to string) error {
	if m.RenameFunc != nil {
		return m.RenameFunc(from, to)
	}
	return os.Rename(from, to)
}

func (m MockFileOps) Remove(path string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(path)
	}
	return os.RemoveAll(path)
}

func (m MockFileOps) MakeExecutable(path string) error {
	if m.MakeExecutableFunc != nil {
		return m.MakeExecutableFunc(path)
	}
	return nil
}

func TestPerformReplacementAndRollback_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-rollback-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	execPath := filepath.Join(tempDir, "runx")
	require.NoError(t, os.WriteFile(execPath, []byte("old-binary-content"), 0755))

	newBinaryData := []byte("new-binary-content")

	verifyCalled := false
	verifyFunc := func(path, expectedVersion string) error {
		verifyCalled = true
		assert.Equal(t, "1.5.0", expectedVersion)
		return nil
	}

	res, errCode, err := PerformReplacementAndRollback(
		execPath,
		newBinaryData,
		"1.5.0",
		"linux",
		MockFileOps{},
		verifyFunc,
	)

	require.NoError(t, err)
	assert.Empty(t, errCode)
	require.NotNil(t, res)
	assert.Equal(t, "1.5.0", res.InstalledVersion)
	assert.True(t, verifyCalled)

	content, err := os.ReadFile(execPath)
	require.NoError(t, err)
	assert.Equal(t, "new-binary-content", string(content))
}

func TestPerformReplacementAndRollback_VerificationFailureRollback(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-rollback-fail-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	execPath := filepath.Join(tempDir, "runx")
	require.NoError(t, os.WriteFile(execPath, []byte("original-binary"), 0755))

	newBinaryData := []byte("broken-binary")

	verifyFunc := func(path, expectedVersion string) error {
		return errors.New("verification check failed: binary crashed")
	}

	res, errCode, err := PerformReplacementAndRollback(
		execPath,
		newBinaryData,
		"2.0.0",
		"linux",
		MockFileOps{},
		verifyFunc,
	)

	assert.Error(t, err)
	assert.Equal(t, "verification_failed", errCode)
	assert.Nil(t, res)

	// Verify original file was restored via rollback
	content, err := os.ReadFile(execPath)
	require.NoError(t, err)
	assert.Equal(t, "original-binary", string(content))
}
