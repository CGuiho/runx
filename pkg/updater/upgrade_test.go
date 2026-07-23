package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CGuiho/runx/pkg/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateNativeBinary(t *testing.T) {
	winData := []byte{0x4D, 0x5A, 0x90, 0x00}
	assert.NoError(t, ValidateNativeBinary(winData, "windows"))

	linData := []byte{0x7F, 'E', 'L', 'F', 0x02}
	assert.NoError(t, ValidateNativeBinary(linData, "linux"))

	macData := []byte{0xCF, 0xFA, 0xED, 0xFE}
	assert.NoError(t, ValidateNativeBinary(macData, "darwin"))

	invalidData := []byte("plain text file")
	assert.Error(t, ValidateNativeBinary(invalidData, "windows"))
	assert.Error(t, ValidateNativeBinary(invalidData, "linux"))
	assert.Error(t, ValidateNativeBinary(invalidData, "darwin"))
}

func TestUpgradeSelf_UpToDate(t *testing.T) {
	mockReleases := []update.GitHubRelease{
		{
			TagName:    "v1.0.0",
			Prerelease: false,
			Draft:      false,
			Assets: []update.GitHubAsset{
				{Name: "runx-windows-x64.exe", BrowserDownloadURL: "https://example.com/win.exe"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mockReleases)
	}))
	defer server.Close()

	opts := UpgradeOptions{
		CurrentVersion: "1.0.0",
		GOOS:           "windows",
		GOARCH:         "amd64",
		APIURL:         server.URL,
	}

	env, err := UpgradeSelf(opts)
	require.NoError(t, err)
	require.NotNil(t, env)

	assert.Equal(t, "up-to-date", env.Outcome)
	assert.Equal(t, "1.0.0", env.Result.InstalledVersion)
}

func TestUpgradeSelf_DryRun(t *testing.T) {
	mockReleases := []update.GitHubRelease{
		{
			TagName:    "v2.0.0",
			Prerelease: false,
			Draft:      false,
			Assets: []update.GitHubAsset{
				{Name: "runx-windows-x64.exe", BrowserDownloadURL: "https://example.com/win.exe"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mockReleases)
	}))
	defer server.Close()

	opts := UpgradeOptions{
		DryRun:         true,
		CurrentVersion: "1.0.0",
		GOOS:           "windows",
		GOARCH:         "amd64",
		APIURL:         server.URL,
	}

	env, err := UpgradeSelf(opts)
	require.NoError(t, err)
	require.NotNil(t, env)

	assert.Equal(t, "dry-run", env.Outcome)
	assert.Equal(t, "2.0.0", env.Plan.TargetVersion)
}

func TestUpgradeSelf_FullUpgrade(t *testing.T) {
	mockReleases := []update.GitHubRelease{
		{
			TagName:    "v2.0.0",
			Prerelease: false,
			Draft:      false,
			Assets: []update.GitHubAsset{
				{Name: "runx-windows-x64.exe", BrowserDownloadURL: "https://example.com/win2.exe"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mockReleases)
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "runx-upgrade-full-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	execPath := filepath.Join(tempDir, "runx.exe")
	require.NoError(t, os.WriteFile(execPath, []byte{0x4D, 0x5A, 0x01}, 0755))

	opts := UpgradeOptions{
		CurrentVersion: "1.0.0",
		ExecutablePath: execPath,
		GOOS:           "windows",
		GOARCH:         "amd64",
		APIURL:         server.URL,
		DownloadFunc: func(url string) ([]byte, error) {
			return []byte{0x4D, 0x5A, 0x02}, nil // Valid MZ binary
		},
		VerifyFunc: func(path, expectedVersion string) error {
			return nil
		},
		FileOps: MockFileOps{},
	}

	env, err := UpgradeSelf(opts)
	require.NoError(t, err)
	require.NotNil(t, env)

	assert.Equal(t, "upgraded", env.Outcome)
	assert.Equal(t, "2.0.0", env.Result.InstalledVersion)
}
