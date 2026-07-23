package update

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCache_ReadWrite(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-cache-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cachePath := filepath.Join(tempDir, "test-cache.json")

	cache := &UpdateCache{
		NewVersionAvailable: true,
		LatestVersion:       "1.2.0",
		UpgradeCommand:      "runx upgrade",
		LastCheck:           time.Now().UTC().Format(time.RFC3339),
	}

	err = WriteCache(cachePath, cache)
	require.NoError(t, err)

	readBack, err := ReadCache(cachePath)
	require.NoError(t, err)
	assert.True(t, readBack.NewVersionAvailable)
	assert.Equal(t, "1.2.0", readBack.LatestVersion)
	assert.Equal(t, "runx upgrade", readBack.UpgradeCommand)
}

func TestUpdateCache_IsCacheFresh(t *testing.T) {
	now := time.Now()
	recent := &UpdateCache{
		LastCheck: now.Add(-1 * time.Hour).UTC().Format(time.RFC3339),
	}
	assert.True(t, IsCacheFresh(recent, 4*time.Hour))

	stale := &UpdateCache{
		LastCheck: now.Add(-5 * time.Hour).UTC().Format(time.RFC3339),
	}
	assert.False(t, IsCacheFresh(stale, 4*time.Hour))

	invalid := &UpdateCache{
		LastCheck: "not-a-date",
	}
	assert.False(t, IsCacheFresh(invalid, 4*time.Hour))
}

func TestUpdateCache_ReadCachedUpdateNotice(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-notice-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cachePath := filepath.Join(tempDir, "cache.json")

	cache := &UpdateCache{
		NewVersionAvailable: true,
		LatestVersion:       "1.5.0",
		UpgradeCommand:      "runx upgrade",
		LastCheck:           time.Now().UTC().Format(time.RFC3339),
	}
	require.NoError(t, WriteCache(cachePath, cache))

	notice, avail := ReadCachedUpdateNotice(cachePath, "1.0.0")
	assert.True(t, avail)
	assert.Contains(t, notice, "v1.5.0")
	assert.Contains(t, notice, "runx upgrade")

	// Same version - no notice
	_, availSame := ReadCachedUpdateNotice(cachePath, "1.5.0")
	assert.False(t, availSame)
}
