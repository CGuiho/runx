package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type GitHubRelease struct {
	TagName     string        `json:"tag_name"`
	Draft       bool          `json:"draft"`
	Prerelease  bool          `json:"prerelease"`
	PublishedAt *string       `json:"published_at"`
	CreatedAt   *string       `json:"created_at"`
	Assets      []GitHubAsset `json:"assets"`
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ReleaseCatalog struct {
	SchemaVersion       int                   `json:"schemaVersion"`
	Command             string                `json:"command"`
	CurrentVersion      string                `json:"currentVersion"`
	LatestStableVersion string                `json:"latestStableVersion"`
	Releases            []ReleaseCatalogEntry `json:"releases"`
}

type ReleaseCatalogEntry struct {
	Tag             string        `json:"tag"`
	Version         string        `json:"version"`
	Channel         string        `json:"channel"`
	Prerelease      bool          `json:"prerelease"`
	PublishedAt     *string       `json:"publishedAt"`
	Current         bool          `json:"current"`
	LatestStable    bool          `json:"latestStable"`
	CompatibleAsset *ReleaseAsset `json:"compatibleAsset"`
}

type ReleaseAsset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ReleasePlatform struct {
	OS      string
	Arch    string
	Variant string
}

func ResolveUpgradePlatform(goos, goarch string) (ReleasePlatform, error) {
	os := ""
	switch goos {
	case "windows", "win32":
		os = "windows"
	case "darwin", "macOS", "macos":
		os = "darwin"
	case "linux":
		os = "linux"
	default:
		return ReleasePlatform{}, fmt.Errorf("self-upgrade is not supported on %s", goos)
	}

	arch := ""
	switch goarch {
	case "x64", "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	default:
		return ReleasePlatform{}, fmt.Errorf("self-upgrade is not supported on architecture %s", goarch)
	}

	return ReleasePlatform{
		OS:      os,
		Arch:    arch,
		Variant: "baseline",
	}, nil
}

func NormalizeReleaseVersion(tag string) string {
	v := strings.TrimPrefix(tag, "@guiho/runx@")
	v = strings.TrimPrefix(v, "v")
	return v
}

func AssetCandidates(platform ReleasePlatform) []string {
	suffix := ""
	if platform.OS == "windows" {
		suffix = ".exe"
	}
	if platform.Arch == "arm64" {
		return []string{fmt.Sprintf("runx-%s-arm64%s", platform.OS, suffix)}
	}
	prefix := fmt.Sprintf("runx-%s-x64", platform.OS)
	switch platform.Variant {
	case "baseline":
		return []string{prefix + "-baseline" + suffix, prefix + suffix, prefix + "-modern" + suffix}
	case "modern":
		return []string{prefix + "-modern" + suffix, prefix + suffix, prefix + "-baseline" + suffix}
	default:
		return []string{prefix + suffix, prefix + "-baseline" + suffix, prefix + "-modern" + suffix}
	}
}

func FindCompatibleAsset(release GitHubRelease, platform ReleasePlatform) *ReleaseAsset {
	candidates := AssetCandidates(platform)
	for _, candidate := range candidates {
		for _, asset := range release.Assets {
			if asset.Name == candidate {
				return &ReleaseAsset{
					Name: asset.Name,
					URL:  asset.BrowserDownloadURL,
				}
			}
		}
	}
	return nil
}

func ReleaseChannel(version string, prerelease bool) string {
	if !prerelease && !strings.Contains(version, "-") {
		return "stable"
	}
	parts := strings.SplitN(version, "-", 2)
	if len(parts) > 1 {
		preParts := strings.Split(parts[1], ".")
		if len(preParts) > 0 && preParts[0] != "" {
			return preParts[0]
		}
	}
	return "prerelease"
}

func CompareVersions(v1, v2 string) int {
	norm1 := NormalizeReleaseVersion(v1)
	norm2 := NormalizeReleaseVersion(v2)

	maj1, min1, pat1, pre1 := parseSemver(norm1)
	maj2, min2, pat2, pre2 := parseSemver(norm2)

	if maj1 != maj2 {
		if maj1 < maj2 {
			return -1
		}
		return 1
	}
	if min1 != min2 {
		if min1 < min2 {
			return -1
		}
		return 1
	}
	if pat1 != pat2 {
		if pat1 < pat2 {
			return -1
		}
		return 1
	}

	if pre1 == "" && pre2 != "" {
		return 1 // release > prerelease
	}
	if pre1 != "" && pre2 == "" {
		return -1 // prerelease < release
	}
	if pre1 != pre2 {
		if pre1 < pre2 {
			return -1
		}
		return 1
	}
	return 0
}

func parseSemver(v string) (major, minor, patch int, prerelease string) {
	parts := strings.SplitN(v, "-", 2)
	if len(parts) > 1 {
		prerelease = parts[1]
	}
	nums := strings.Split(parts[0], ".")
	if len(nums) > 0 {
		major, _ = strconv.Atoi(nums[0])
	}
	if len(nums) > 1 {
		minor, _ = strconv.Atoi(nums[1])
	}
	if len(nums) > 2 {
		patch, _ = strconv.Atoi(nums[2])
	}
	return
}

func FetchReleaseCatalog(apiURL string, platform ReleasePlatform, currentVersion string, client *http.Client) (*ReleaseCatalog, error) {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	if apiURL == "" {
		apiURL = "https://api.github.com/repos/CGuiho/runx/releases"
	}

	var allReleases []GitHubRelease
	page := 1
	hasNext := true

	for hasNext {
		url := fmt.Sprintf("%s?per_page=100&page=%d", apiURL, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve RunX releases page %d: %w", page, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("could not read releases page %d: %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("could not retrieve RunX releases page %d: HTTP %d", page, resp.StatusCode)
		}

		var pageReleases []GitHubRelease
		if err := json.Unmarshal(body, &pageReleases); err != nil {
			return nil, fmt.Errorf("could not retrieve RunX releases page %d: malformed release", page)
		}

		for _, rel := range pageReleases {
			if !rel.Draft {
				allReleases = append(allReleases, rel)
			}
		}

		linkHeader := resp.Header.Get("Link")
		hasNext = hasNextLink(linkHeader)
		page++
	}

	normCurrent := NormalizeReleaseVersion(currentVersion)
	var catalogEntries []ReleaseCatalogEntry

	for _, rel := range allReleases {
		v := NormalizeReleaseVersion(rel.TagName)
		isPre := rel.Prerelease || strings.Contains(v, "-")
		pubAt := rel.PublishedAt
		if pubAt == nil {
			pubAt = rel.CreatedAt
		}

		entry := ReleaseCatalogEntry{
			Tag:             rel.TagName,
			Version:         v,
			Channel:         ReleaseChannel(v, isPre),
			Prerelease:      isPre,
			PublishedAt:     pubAt,
			Current:         v == normCurrent,
			LatestStable:    false,
			CompatibleAsset: FindCompatibleAsset(rel, platform),
		}
		catalogEntries = append(catalogEntries, entry)
	}

	sort.Slice(catalogEntries, func(i, j int) bool {
		return CompareVersions(catalogEntries[i].Version, catalogEntries[j].Version) > 0
	})

	latestStableVersion := ""
	for _, entry := range catalogEntries {
		if !entry.Prerelease {
			latestStableVersion = entry.Version
			break
		}
	}

	for i := range catalogEntries {
		if latestStableVersion != "" && catalogEntries[i].Version == latestStableVersion {
			catalogEntries[i].LatestStable = true
		}
	}

	return &ReleaseCatalog{
		SchemaVersion:       1,
		Command:             "runx upgrade list",
		CurrentVersion:      currentVersion,
		LatestStableVersion: latestStableVersion,
		Releases:            catalogEntries,
	}, nil
}

func hasNextLink(link string) bool {
	if link == "" {
		return false
	}
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
	return re.MatchString(link)
}
