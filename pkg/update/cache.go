package update

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UpdateCache struct {
	NewVersionAvailable bool   `json:"newVersionAvailable"`
	LatestVersion       string `json:"latestVersion"`
	UpgradeCommand      string `json:"upgradeCommand,omitempty"`
	LastCheck           string `json:"lastCheck"`
}

const CacheTTL = 4 * time.Hour

func GetDefaultCachePath() string {
	if envPath := os.Getenv("RUNX_CACHE_PATH"); envPath != "" {
		return envPath
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".runx", "update-cache.json")
}

func ReadCache(path string) (*UpdateCache, error) {
	if path == "" {
		path = GetDefaultCachePath()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cache UpdateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("invalid update cache JSON: %w", err)
	}
	return &cache, nil
}

func WriteCache(path string, cache *UpdateCache) error {
	if path == "" {
		path = GetDefaultCachePath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := fmt.Sprintf("%s.tmp.%d", path, time.Now().UnixNano())
	if err := os.WriteFile(tmpPath, append(data, '\n'), 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

func IsCacheFresh(cache *UpdateCache, ttl time.Duration) bool {
	if cache == nil || cache.LastCheck == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, cache.LastCheck)
	if err != nil {
		return false
	}
	return time.Since(t) >= 0 && time.Since(t) < ttl
}

func ReadCachedUpdateNotice(path string, currentVersion string) (string, bool) {
	cache, err := ReadCache(path)
	if err != nil || cache == nil {
		return "", false
	}

	if cache.NewVersionAvailable && cache.LatestVersion != "" {
		if CompareVersions(cache.LatestVersion, currentVersion) > 0 {
			cmd := cache.UpgradeCommand
			if cmd == "" {
				cmd = "runx upgrade"
			}
			notice := fmt.Sprintf("  ⚠ New version available: v%s\n    Run `%s` to update.", strings.TrimPrefix(cache.LatestVersion, "v"), cmd)
			return notice, true
		}
	}
	return "", false
}
