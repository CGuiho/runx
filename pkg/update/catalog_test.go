package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
	assert.Equal(t, 1, CompareVersions("1.2.0", "1.1.0"))
	assert.Equal(t, -1, CompareVersions("1.0.0", "1.1.0"))
	assert.Equal(t, 0, CompareVersions("1.1.0", "1.1.0"))

	// Prerelease comparison
	assert.Equal(t, 1, CompareVersions("1.0.0", "1.0.0-dev"))
	assert.Equal(t, -1, CompareVersions("1.0.0-alpha", "1.0.0"))
	assert.Equal(t, 1, CompareVersions("v1.5.0", "@guiho/runx@1.4.0"))
}

func TestResolveUpgradePlatform(t *testing.T) {
	pWin, err := ResolveUpgradePlatform("windows", "amd64")
	require.NoError(t, err)
	assert.Equal(t, "windows", pWin.OS)
	assert.Equal(t, "amd64", pWin.Arch)

	pMac, err := ResolveUpgradePlatform("darwin", "arm64")
	require.NoError(t, err)
	assert.Equal(t, "darwin", pMac.OS)
	assert.Equal(t, "arm64", pMac.Arch)

	_, err = ResolveUpgradePlatform("unsupported", "amd64")
	assert.Error(t, err)
}

func TestAssetCandidatesPreferProtocolV1Payload(t *testing.T) {
	tests := []struct {
		platform ReleasePlatform
		want     []string
	}{
		{
			platform: ReleasePlatform{OS: "windows", Arch: "amd64", Target: "runx-windows-amd64"},
			want:     []string{"runx-payload-windows-amd64.exe", "runx-windows-amd64.exe"},
		},
		{
			platform: ReleasePlatform{OS: "linux", Arch: "armv7", Target: "runx-linux-armv7"},
			want:     []string{"runx-payload-linux-armv7", "runx-linux-armv7"},
		},
	}
	for _, testCase := range tests {
		assert.Equal(t, testCase.want, AssetCandidates(testCase.platform))
	}

	release := GitHubRelease{Assets: []GitHubAsset{
		{Name: "runx-windows-amd64.exe", BrowserDownloadURL: "https://example.com/legacy"},
		{Name: "runx-payload-windows-amd64.exe", BrowserDownloadURL: "https://example.com/payload"},
	}}
	asset := FindCompatibleAsset(release, tests[0].platform)
	require.NotNil(t, asset)
	assert.Equal(t, "runx-payload-windows-amd64.exe", asset.Name)
}

func TestFetchReleaseCatalog(t *testing.T) {
	mockReleases := []GitHubRelease{
		{
			TagName:    "v1.0.0",
			Prerelease: false,
			Draft:      false,
			Assets: []GitHubAsset{
				{Name: "runx-payload-windows-amd64.exe", BrowserDownloadURL: "https://example.com/win.exe"},
				{Name: "runx-payload-linux-amd64", BrowserDownloadURL: "https://example.com/linux"},
				{Name: "checksums.txt", BrowserDownloadURL: "https://example.com/checksums"},
			},
		},
		{
			TagName:    "v1.1.0-alpha.1",
			Prerelease: true,
			Draft:      false,
			Assets: []GitHubAsset{
				{Name: "runx-windows-amd64.exe", BrowserDownloadURL: "https://example.com/win-alpha.exe"},
				{Name: "checksums.txt", BrowserDownloadURL: "https://example.com/checksums-alpha"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockReleases)
	}))
	defer server.Close()

	platform := ReleasePlatform{OS: "windows", Arch: "amd64", Target: "runx-windows-amd64"}
	catalog, err := FetchReleaseCatalog(server.URL, platform, "0.9.0", server.Client())
	require.NoError(t, err)

	assert.Equal(t, "0.9.0", catalog.CurrentVersion)
	assert.Equal(t, "1.0.0", catalog.LatestStableVersion)
	assert.Len(t, catalog.Releases, 2)
	assert.True(t, catalog.Releases[0].Prerelease)  // 1.1.0-alpha.1 sorted first
	assert.False(t, catalog.Releases[1].Prerelease) // 1.0.0
	assert.True(t, catalog.Releases[1].LatestStable)
	assert.NotNil(t, catalog.Releases[1].CompatibleAsset)
	assert.Equal(t, "https://example.com/win.exe", catalog.Releases[1].CompatibleAsset.URL)
}
