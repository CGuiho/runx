// Package installstate owns the canonical RunX installation layout under the
// shared GUIHO home: the stable launcher path, immutable versioned payloads,
// the atomic current-version pointer, and the installed-artifacts ledger.
package installstate

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// ProtocolVersion is the frozen launcher/payload protocol understood by this
// build. Protocol 1 is the first Convention 0001 contract.
const ProtocolVersion = 1

func homeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	if home == "" {
		return "", fmt.Errorf("resolve user home: empty path")
	}
	return home, nil
}

// GUISafeHome returns $HOME/.guiho.
func GUISafeHome() (string, error) {
	home, err := homeDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".guiho"), nil
}

// CLIDir returns $HOME/.guiho/runx.
func CLIDir() (string, error) {
	root, err := GUISafeHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "runx"), nil
}

// BinDir returns the shared launcher directory $HOME/.guiho/bin.
func BinDir() (string, error) {
	root, err := GUISafeHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "bin"), nil
}

// VersionsDir returns $HOME/.guiho/runx/versions.
func VersionsDir() (string, error) {
	dir, err := CLIDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "versions"), nil
}

// VersionDir returns the immutable payload directory for one version.
func VersionDir(version string) (string, error) {
	if !isValidVersionName(version) {
		return "", fmt.Errorf("invalid version %q", version)
	}
	dir, err := VersionsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, version), nil
}

// CurrentPointerPath returns $HOME/.guiho/runx/current.json.
func CurrentPointerPath() (string, error) {
	dir, err := CLIDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "current.json"), nil
}

// InstalledLedgerPath returns $HOME/.guiho/runx/installed-artifacts.json.
func InstalledLedgerPath() (string, error) {
	dir, err := CLIDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "installed-artifacts.json"), nil
}

// TempRoot returns the shared staging root $HOME/.guiho/.temp.
func TempRoot() (string, error) {
	root, err := GUISafeHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, ".temp"), nil
}

// LauncherName returns the platform launcher file name.
func LauncherName() string {
	if runtime.GOOS == "windows" {
		return "runx.exe"
	}
	return "runx"
}

// PayloadName returns the platform payload file name inside a version
// directory.
func PayloadName() string {
	if runtime.GOOS == "windows" {
		return "runx-payload.exe"
	}
	return "runx-payload"
}

// LauncherPath returns $HOME/.guiho/bin/runx(.exe).
func LauncherPath() (string, error) {
	bin, err := BinDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(bin, LauncherName()), nil
}

// In-variants resolve canonical locations for an explicitly provided home
// directory so lifecycle code can inject homes instead of relying on process
// environment state.

func dirIn(home string, parts ...string) string {
	return filepath.Join(append([]string{home, ".guiho"}, parts...)...)
}

// GUISafeHomeIn returns <home>/.guiho.
func GUISafeHomeIn(home string) string { return filepath.Join(home, ".guiho") }

// CLIDirIn returns <home>/.guiho/runx.
func CLIDirIn(home string) string { return dirIn(home, "runx") }

// BinDirIn returns <home>/.guiho/bin.
func BinDirIn(home string) string { return dirIn(home, "bin") }

// VersionsDirIn returns <home>/.guiho/runx/versions.
func VersionsDirIn(home string) string { return dirIn(home, "runx", "versions") }

// VersionDirIn returns the immutable version directory under an explicit home.
func VersionDirIn(home, version string) (string, error) {
	if !isValidVersionName(version) {
		return "", fmt.Errorf("invalid version %q", version)
	}
	return filepath.Join(VersionsDirIn(home), version), nil
}

// CurrentPointerPathIn returns <home>/.guiho/runx/current.json.
func CurrentPointerPathIn(home string) string { return filepath.Join(CLIDirIn(home), "current.json") }

// InstalledLedgerPathIn returns <home>/.guiho/runx/installed-artifacts.json.
func InstalledLedgerPathIn(home string) string {
	return filepath.Join(CLIDirIn(home), "installed-artifacts.json")
}

// LauncherPathIn returns <home>/.guiho/bin/runx(.exe).
func LauncherPathIn(home string) string { return filepath.Join(BinDirIn(home), LauncherName()) }
