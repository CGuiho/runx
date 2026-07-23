package update

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type WorkerOptions struct {
	CachePath      string
	CurrentVersion string
	GOOS           string
	GOARCH         string
	APIURL         string
	HTTPClient     *http.Client
	Timeout        time.Duration
	Now            func() time.Time
}

func RunUpdateWorker(opts WorkerOptions) (*UpdateCache, error) {
	if os.Getenv("RUNX_DISABLE_UPDATE_WORKER") == "1" {
		return nil, nil
	}

	goos := opts.GOOS
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := opts.GOARCH
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	currentVer := opts.CurrentVersion
	if currentVer == "" {
		currentVer = "0.8.0-dev"
	}
	cachePath := opts.CachePath
	if cachePath == "" {
		cachePath = GetDefaultCachePath()
	}

	now := opts.Now
	if now == nil {
		now = time.Now
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}

	platform, err := ResolveUpgradePlatform(goos, goarch)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve platform: %w", err)
	}

	catalog, err := FetchReleaseCatalog(opts.APIURL, platform, currentVer, client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release catalog: %w", err)
	}

	latestVersion := catalog.LatestStableVersion
	if latestVersion == "" {
		latestVersion = currentVer
	}

	newVersionAvailable := CompareVersions(latestVersion, currentVer) > 0

	cache := &UpdateCache{
		NewVersionAvailable: newVersionAvailable,
		LatestVersion:       latestVersion,
		LastCheck:           now().UTC().Format(time.RFC3339),
	}
	if newVersionAvailable {
		cache.UpgradeCommand = "runx upgrade"
	}

	if err := WriteCache(cachePath, cache); err != nil {
		return nil, fmt.Errorf("failed to write update cache: %w", err)
	}

	return cache, nil
}

func SpawnUpdateWorker(execPath string, args ...string) error {
	if os.Getenv("RUNX_DISABLE_UPDATE_WORKER") == "1" {
		return nil
	}

	if execPath == "" {
		var err error
		execPath, err = os.Executable()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(execPath, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	configureDetachedProcess(cmd)

	return cmd.Start()
}
