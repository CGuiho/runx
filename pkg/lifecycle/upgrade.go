// Package lifecycle implements the Convention 0001 whole-release upgrade for
// protocol-v1 installations: stage under the shared temp root, verify every
// artifact, self-test the staged payload, install an immutable version,
// atomically activate through the stable launcher's pointer, and verify via
// the launcher — never by replacing the running executable.
package lifecycle

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/CGuiho/runx/pkg/maintenance"
	"github.com/CGuiho/runx/pkg/update"
)

// ErrLegacyInstallation reports that the running installation predates
// protocol v1 and must use the legacy direct-replacement updater.
var ErrLegacyInstallation = errors.New("installation predates protocol v1")

// ErrLegacyRelease reports that the selected release does not publish the
// protocol-v1 payload required by a v1 installation.
var ErrLegacyRelease = errors.New("selected release has no protocol-v1 payload asset")

const (
	maxArtifactBytes = 256 << 20
	maxCatalogBytes  = 8 << 20
)

// Options configures a whole-release upgrade.
type Options struct {
	CurrentVersion   string
	RequestedVersion string // exact version, empty means latest stable
	BuildTarget      string
	HTTPClient       *http.Client
	HomeDir          func() (string, error)
	DryRun           bool
}

// Result reports what a successful (or planned) upgrade did.
type Result struct {
	Outcome         string   `json:"outcome"` // "upgraded", "up-to-date", "dry-run"
	PreviousVersion string   `json:"previousVersion"`
	TargetVersion   string   `json:"targetVersion"`
	LauncherPath    string   `json:"launcherPath"`
	PayloadPath     string   `json:"payloadPath"`
	Verified        bool     `json:"verified"`
	AssetsInstalled []string `json:"assetsInstalled,omitempty"`
}

// IsProtocolV1Installation reports whether the home contains a committed
// protocol-v1 pointer.
func IsProtocolV1Installation(homeDir func() (string, error)) bool {
	home, err := homeDir()
	if err != nil || home == "" {
		return false
	}
	pointer, err := installstate.ReadPointerIn(home)
	return err == nil && pointer != nil
}

// UpgradeWholeRelease performs the synchronous protocol-v1 transaction.
func UpgradeWholeRelease(opts Options) (*Result, error) {
	home, err := opts.homeDir()()
	if err != nil {
		return nil, fmt.Errorf("resolve home: %w", err)
	}
	pointer, err := installstate.ReadPointerIn(home)
	if err != nil {
		return nil, fmt.Errorf("read installation pointer: %w", err)
	}
	if pointer == nil {
		return nil, ErrLegacyInstallation
	}

	platform, err := update.ResolveBuildTarget(opts.BuildTarget, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, err
	}
	catalog, err := update.FetchReleaseCatalog("", platform, opts.CurrentVersion, opts.HTTPClient)
	if err != nil {
		return nil, fmt.Errorf("fetch release catalog: %w", err)
	}

	entry, err := selectRelease(catalog, opts.RequestedVersion)
	if err != nil {
		return nil, err
	}
	if entry.Version == pointer.Active && opts.RequestedVersion == "" {
		return &Result{Outcome: "up-to-date", PreviousVersion: pointer.Active, TargetVersion: pointer.Active}, nil
	}
	if entry.CompatibleAsset == nil || !strings.HasPrefix(entry.CompatibleAsset.Name, "runx-payload-") {
		return nil, fmt.Errorf("%w: %s", ErrLegacyRelease, entry.Version)
	}

	result := &Result{PreviousVersion: pointer.Active, TargetVersion: entry.Version}
	if opts.DryRun {
		result.Outcome = "dry-run"
		return result, nil
	}

	tagPath := releaseTagPath(entry.Tag)
	base := "https://github.com/CGuiho/runx/releases/download/" + tagPath
	staging, cleanup, err := installstate.StagingDir("upgrade")
	if err != nil {
		return nil, err
	}
	defer cleanup()

	assets := []string{
		entry.CompatibleAsset.Name,
		launcherAssetName(entry.CompatibleAsset.Name),
		"checksums.txt", "artifacts.json",
		"guiho-s-runx.zip", "guiho-i-runx.md",
		"runx.schema.json", "runx.global.schema.json",
	}
	for _, asset := range assets {
		if err := downloadTo(opts.HTTPClient, base+"/"+asset, filepath.Join(staging, asset)); err != nil {
			return nil, fmt.Errorf("download %s: %w", asset, err)
		}
	}
	if err := verifyChecksums(staging, assets); err != nil {
		return nil, fmt.Errorf("verify staged artifacts: %w", err)
	}

	payloadName := installstate.PayloadName()
	payloadPath := filepath.Join(staging, entry.CompatibleAsset.Name)
	if err := os.Chmod(payloadPath, 0o755); err != nil {
		return nil, err
	}
	stagedVersion, err := runCommandOutput(payloadPath, "--version")
	if err != nil || strings.TrimSpace(stagedVersion) != entry.Version {
		return nil, fmt.Errorf("staged payload version verification failed")
	}
	if err := runCommand(payloadPath, "__self-test"); err != nil {
		return nil, fmt.Errorf("staged payload failed its self-test: %w", err)
	}

	versionDir, err := installstate.VersionDirIn(home, entry.Version)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(versionDir, 0o755); err != nil {
		return nil, err
	}
	installedPayload := filepath.Join(versionDir, payloadName)
	if err := copyFile(payloadPath, installedPayload, 0o755); err != nil {
		return nil, err
	}
	if err := copyFile(filepath.Join(staging, "artifacts.json"), filepath.Join(versionDir, "release-artifacts.json"), 0o644); err != nil {
		return nil, err
	}
	result.PayloadPath = installedPayload

	// Project instruction reconciliation uses the staged skill bytes.
	skillBytes, err := skillFromArchive(filepath.Join(staging, "guiho-s-runx.zip"))
	if err == nil {
		cwd, _ := os.Getwd()
		_, _ = maintenance.MaintainAgentIntegration(cwd, home, skillBytes)
	}

	previousPointer := *pointer
	newPointer := installstate.Pointer{Protocol: installstate.ProtocolVersion, Active: entry.Version, Previous: pointer.Active}
	if err := installstate.WritePointerIn(home, newPointer); err != nil {
		return nil, fmt.Errorf("activate new version: %w", err)
	}

	launcherPath := installstate.LauncherPathIn(home)
	result.LauncherPath = launcherPath
	activatedVersion, verifyErr := runCommandOutput(launcherPath, "--version")
	if verifyErr != nil || strings.TrimSpace(activatedVersion) != entry.Version {
		_ = installstate.WritePointerIn(home, previousPointer)
		_ = os.RemoveAll(versionDir)
		if verifyErr != nil {
			return nil, fmt.Errorf("activated launcher verification failed and was rolled back: %w", verifyErr)
		}
		return nil, fmt.Errorf("activated launcher reports %q, want %q; rolled back", strings.TrimSpace(activatedVersion), entry.Version)
	}
	result.Outcome = "upgraded"
	result.Verified = true
	result.AssetsInstalled = assets
	return result, nil
}

func (o Options) homeDir() func() (string, error) {
	if o.HomeDir != nil {
		return o.HomeDir
	}
	return os.UserHomeDir
}

func selectRelease(catalog *update.ReleaseCatalog, requested string) (*update.ReleaseCatalogEntry, error) {
	for index := range catalog.Releases {
		entry := &catalog.Releases[index]
		if requested != "" && entry.Version == requested {
			return entry, nil
		}
		if requested == "" && entry.LatestStable {
			return entry, nil
		}
	}
	if requested != "" {
		return nil, fmt.Errorf("release %q not found in catalog", requested)
	}
	return nil, fmt.Errorf("no stable release found in catalog")
}

func releaseTagPath(tag string) string {
	replacer := strings.NewReplacer("/", "%2F", "@", "%40")
	return replacer.Replace(tag)
}

func launcherAssetName(payloadAsset string) string {
	name := strings.TrimPrefix(payloadAsset, "runx-payload-")
	if strings.HasSuffix(name, ".exe") {
		return "runx-launcher-" + name
	}
	return "runx-launcher-" + name
}

func downloadTo(client *http.Client, url, destination string) error {
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d for %s", response.StatusCode, url)
	}
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()
	written, err := io.Copy(file, io.LimitReader(response.Body, maxArtifactBytes+1))
	if err != nil {
		return err
	}
	if written > maxArtifactBytes {
		return fmt.Errorf("artifact %s exceeds size limit", destination)
	}
	return nil
}

func verifyChecksums(staging string, assets []string) error {
	data, err := os.ReadFile(filepath.Join(staging, "checksums.txt"))
	if err != nil {
		return err
	}
	expected := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			continue
		}
		expected[strings.TrimPrefix(fields[1], "*")] = strings.ToLower(fields[0])
	}
	for _, asset := range assets {
		want, ok := expected[asset]
		if !ok {
			return fmt.Errorf("checksum entry missing for %s", asset)
		}
		file, err := os.Open(filepath.Join(staging, asset))
		if err != nil {
			return err
		}
		hash := sha256.New()
		_, copyErr := io.Copy(hash, file)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		if got := hex.EncodeToString(hash.Sum(nil)); got != want {
			return fmt.Errorf("checksum mismatch for %s", asset)
		}
	}
	return nil
}

func copyFile(source, destination string, mode os.FileMode) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destination, data, mode)
}

// skillFromArchive extracts SKILL.md bytes from a guiho-s-runx zip without
// extracting to disk. It is best-effort: failure only skips instruction
// reconciliation, never the upgrade itself.
func skillFromArchive(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	stat, err := reader.Stat()
	if err != nil {
		return "", err
	}
	if stat.Size() > maxCatalogBytes {
		return "", fmt.Errorf("skill archive exceeds size limit")
	}
	return zipReadFile(path, "guiho-s-runx/SKILL.md")
}

// runCommandOutput runs an executable and returns its trimmed stdout.
func runCommandOutput(executable string, arguments ...string) (string, error) {
	command := exec.Command(executable, arguments...)
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// runCommand runs an executable discarding output, requiring exit zero.
func runCommand(executable string, arguments ...string) error {
	_, err := runCommandOutput(executable, arguments...)
	return err
}

// zipReadFile reads one member of a zip archive into memory.
func zipReadFile(archivePath, name string) (string, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if file.Name != name {
			continue
		}
		opened, err := file.Open()
		if err != nil {
			return "", err
		}
		defer opened.Close()
		data, err := io.ReadAll(io.LimitReader(opened, maxCatalogBytes))
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", fmt.Errorf("archive %s is missing %s", archivePath, name)
}
