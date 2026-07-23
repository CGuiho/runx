/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type BinaryTarget struct {
	GOOS      string
	GOARCH    string
	GOAMD64   string
	AssetName string
}

var BinaryTargets = []BinaryTarget{
	{GOOS: "linux", GOARCH: "arm64", GOAMD64: "", AssetName: "runx-linux-arm64"},
	{GOOS: "linux", GOARCH: "amd64", GOAMD64: "v1", AssetName: "runx-linux-x64"},
	{GOOS: "linux", GOARCH: "amd64", GOAMD64: "v1", AssetName: "runx-linux-x64-baseline"},
	{GOOS: "linux", GOARCH: "amd64", GOAMD64: "v3", AssetName: "runx-linux-x64-modern"},
	{GOOS: "darwin", GOARCH: "arm64", GOAMD64: "", AssetName: "runx-darwin-arm64"},
	{GOOS: "darwin", GOARCH: "amd64", GOAMD64: "v1", AssetName: "runx-darwin-x64"},
	{GOOS: "darwin", GOARCH: "amd64", GOAMD64: "v1", AssetName: "runx-darwin-x64-baseline"},
	{GOOS: "darwin", GOARCH: "amd64", GOAMD64: "v3", AssetName: "runx-darwin-x64-modern"},
	{GOOS: "windows", GOARCH: "arm64", GOAMD64: "", AssetName: "runx-windows-arm64.exe"},
	{GOOS: "windows", GOARCH: "amd64", GOAMD64: "v1", AssetName: "runx-windows-x64.exe"},
	{GOOS: "windows", GOARCH: "amd64", GOAMD64: "v1", AssetName: "runx-windows-x64-baseline.exe"},
	{GOOS: "windows", GOARCH: "amd64", GOAMD64: "v3", AssetName: "runx-windows-x64-modern.exe"},
}

var AgentAssetNames = []string{
	"guiho-s-runx.md",
	"guiho-i-runx.md",
}

func getExpectedReleaseAssetNames() []string {
	var names []string
	for _, target := range BinaryTargets {
		names = append(names, target.AssetName)
	}
	names = append(names, AgentAssetNames...)
	return names
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func getRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return dir
	}
	parent := filepath.Dir(dir)
	if _, err := os.Stat(filepath.Join(parent, "go.mod")); err == nil {
		return parent
	}
	return dir
}

func main() {
	expected := getExpectedReleaseAssetNames()
	if len(expected) != 14 {
		log.Fatalf("RunX release matrix must contain exactly fourteen unique assets, got %d", len(expected))
	}
	assetSet := make(map[string]bool)
	for _, name := range expected {
		assetSet[name] = true
	}
	if len(assetSet) != 14 {
		log.Fatalf("RunX release matrix contains duplicate asset names")
	}

	root := getRootDir()
	binDir := filepath.Join(root, "bin")

	if err := os.RemoveAll(binDir); err != nil {
		log.Fatalf("failed to clean bin directory: %v", err)
	}
	if err := os.MkdirAll(binDir, 0755); err != nil {
		log.Fatalf("failed to create bin directory: %v", err)
	}

	for _, target := range BinaryTargets {
		output := filepath.Join(binDir, target.AssetName)

		cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", output, ".")
		cmd.Dir = root

		env := os.Environ()
		env = append(env, "CGO_ENABLED=0")
		env = append(env, fmt.Sprintf("GOOS=%s", target.GOOS))
		env = append(env, fmt.Sprintf("GOARCH=%s", target.GOARCH))
		if target.GOAMD64 != "" {
			env = append(env, fmt.Sprintf("GOAMD64=%s", target.GOAMD64))
		}
		cmd.Env = env

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to build %s:\n%s\n%v", target.AssetName, string(out), err)
		}

		info, err := os.Stat(output)
		if err != nil || info.Size() == 0 {
			log.Fatalf("Built binary is empty or missing: %s", target.AssetName)
		}

		fmt.Printf("built: %s\n", target.AssetName)
	}

	skillSrc := filepath.Join(root, "skills", "guiho-s-runx", "SKILL.md")
	skillDst := filepath.Join(binDir, "guiho-s-runx.md")
	if err := copyFile(skillSrc, skillDst); err != nil {
		log.Fatalf("failed to copy skill asset: %v", err)
	}

	promptSrc := filepath.Join(root, "prompts", "guiho-i-runx.md")
	promptDst := filepath.Join(binDir, "guiho-i-runx.md")
	if err := copyFile(promptSrc, promptDst); err != nil {
		log.Fatalf("failed to copy prompt asset: %v", err)
	}

	entries, err := os.ReadDir(binDir)
	if err != nil {
		log.Fatalf("failed to read bin directory: %v", err)
	}

	var observed []string
	for _, entry := range entries {
		if !entry.IsDir() {
			observed = append(observed, entry.Name())
		}
	}
	sort.Strings(observed)

	expectedSorted := append([]string(nil), expected...)
	sort.Strings(expectedSorted)

	if strings.Join(observed, ",") != strings.Join(expectedSorted, ",") {
		log.Fatalf("Release asset mismatch.\nExpected: %s\nObserved: %s", strings.Join(expectedSorted, ", "), strings.Join(observed, ", "))
	}

	fmt.Println("verified exactly 14 release assets")
}
