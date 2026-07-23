/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var ExpectedReleaseAssetNames = []string{
	"runx-linux-arm64",
	"runx-linux-x64",
	"runx-linux-x64-baseline",
	"runx-linux-x64-modern",
	"runx-darwin-arm64",
	"runx-darwin-x64",
	"runx-darwin-x64-baseline",
	"runx-darwin-x64-modern",
	"runx-windows-arm64.exe",
	"runx-windows-x64.exe",
	"runx-windows-x64-baseline.exe",
	"runx-windows-x64-modern.exe",
	"guiho-s-runx.md",
	"guiho-i-runx.md",
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

func verifyMarkdownAsset(binDir, assetName, expectedName string) {
	path := filepath.Join(binDir, assetName)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read markdown asset %s: %v", assetName, err)
	}
	if len(data) == 0 {
		log.Fatalf("Release Markdown asset is empty: %s", assetName)
	}
	if len(data) >= 2 && data[0] == 0x4d && data[1] == 0x5a {
		log.Fatalf("Release Markdown asset has executable header: %s", assetName)
	}
	if bytes.IndexByte(data, 0) != -1 {
		log.Fatalf("Release Markdown asset contains binary NUL bytes: %s", assetName)
	}
	if !utf8.Valid(data) {
		log.Fatalf("Release Markdown asset is not valid UTF-8: %s", assetName)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		log.Fatalf("Release Markdown asset does not start with YAML frontmatter: %s", assetName)
	}
	pattern := fmt.Sprintf(`(?m)^name:\s*%s\s*$`, regexp.QuoteMeta(expectedName))
	matched, err := regexp.MatchString(pattern, text)
	if err != nil || !matched {
		log.Fatalf("Release Markdown asset has invalid frontmatter identity: %s", assetName)
	}
}

func main() {
	root := getRootDir()
	binDir := filepath.Join(root, "bin")

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

	expected := append([]string(nil), ExpectedReleaseAssetNames...)
	sort.Strings(expected)

	if strings.Join(observed, ",") != strings.Join(expected, ",") {
		log.Fatalf("Expected exactly %d release assets.\nExpected: %s\nObserved: %s",
			len(expected), strings.Join(expected, ", "), strings.Join(observed, ", "))
	}

	markdownAssets := []struct {
		assetName    string
		expectedName string
	}{
		{"guiho-s-runx.md", "guiho-s-runx"},
		{"guiho-i-runx.md", "guiho-i-runx"},
	}

	for _, item := range markdownAssets {
		verifyMarkdownAsset(binDir, item.assetName, item.expectedName)
	}

	outputStruct := struct {
		Count  int      `json:"count"`
		Assets []string `json:"assets"`
	}{
		Count:  len(observed),
		Assets: observed,
	}

	jsonBytes, err := json.MarshalIndent(outputStruct, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal verification output: %v", err)
	}

	fmt.Println(string(jsonBytes))
}
