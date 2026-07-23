package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunUpdateWorker(t *testing.T) {
	mockReleases := []GitHubRelease{
		{
			TagName:    "v2.0.0",
			Prerelease: false,
			Draft:      false,
			Assets: []GitHubAsset{
				{Name: "runx-windows-x64.exe", BrowserDownloadURL: "https://example.com/win2.exe"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockReleases)
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "runx-worker-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cachePath := filepath.Join(tempDir, "cache.json")

	opts := WorkerOptions{
		CachePath:      cachePath,
		CurrentVersion: "1.0.0",
		GOOS:           "windows",
		GOARCH:         "amd64",
		APIURL:         server.URL,
		HTTPClient:     server.Client(),
		Timeout:        5 * time.Second,
	}

	cache, err := RunUpdateWorker(opts)
	require.NoError(t, err)
	require.NotNil(t, cache)

	assert.True(t, cache.NewVersionAvailable)
	assert.Equal(t, "2.0.0", cache.LatestVersion)
	assert.Equal(t, "runx upgrade", cache.UpgradeCommand)

	// Verify written file
	savedCache, err := ReadCache(cachePath)
	require.NoError(t, err)
	assert.Equal(t, cache.LatestVersion, savedCache.LatestVersion)
}
