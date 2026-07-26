package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type workerLease struct {
	PID       int       `json:"pid"`
	CreatedAt time.Time `json:"createdAt"`
}

type WorkerOptions struct {
	CachePath      string
	CurrentVersion string
	GOOS           string
	GOARCH         string
	BuildTarget    string
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
		currentVer = "dev"
	}
	cachePath := opts.CachePath
	if cachePath == "" {
		cachePath = GetDefaultCachePath()
	}
	if cache, err := ReadCache(cachePath); err == nil && IsCacheFresh(cache, CacheTTL) {
		return cache, nil
	}
	release, acquired, err := acquireWorkerLease(cachePath+".lease", 30*time.Second)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, nil
	}
	defer release()

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

	platform, err := ResolveBuildTarget(opts.BuildTarget, goos, goarch)
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

func acquireWorkerLease(path string, staleAfter time.Duration) (func(), bool, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, false, err
	}
	create := func() (*os.File, error) { return os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) }
	file, err := create()
	if err != nil {
		if !os.IsExist(err) {
			return nil, false, err
		}
		data, readErr := os.ReadFile(path)
		var lease workerLease
		if readErr == nil && json.Unmarshal(data, &lease) == nil && time.Since(lease.CreatedAt) >= 0 && time.Since(lease.CreatedAt) < staleAfter {
			return func() {}, false, nil
		}
		if removeErr := os.Remove(path); removeErr != nil && !os.IsNotExist(removeErr) {
			return nil, false, removeErr
		}
		file, err = create()
		if err != nil {
			if os.IsExist(err) {
				return func() {}, false, nil
			}
			return nil, false, err
		}
	}
	leaseData, _ := json.Marshal(workerLease{PID: os.Getpid(), CreatedAt: time.Now().UTC()})
	if _, err := file.Write(leaseData); err != nil {
		file.Close()
		os.Remove(path)
		return nil, false, err
	}
	if err := file.Close(); err != nil {
		os.Remove(path)
		return nil, false, err
	}
	return func() { _ = os.Remove(path) }, true, nil
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
