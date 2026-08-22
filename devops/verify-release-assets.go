//go:build ignore

// Command verify-release-assets proves a built release satisfies the
// Convention 0001 contract: the exact declared artifact set, an artifacts.json
// manifest whose digests match the files, checksums.txt covering every other
// artifact, and a complete skill archive.
package main

import (
	"archive/zip"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type manifestArtifact struct {
	ID            string   `json:"id"`
	File          string   `json:"file"`
	SHA256        string   `json:"sha256"`
	Kind          string   `json:"kind"`
	OS            string   `json:"os,omitempty"`
	Arch          string   `json:"arch,omitempty"`
	Version       string   `json:"version"`
	InstalledPath string   `json:"installedPath"`
	Projections   []string `json:"projections,omitempty"`
}

type releaseManifest struct {
	Schema    int                `json:"schema"`
	Protocol  int                `json:"protocol"`
	Name      string             `json:"name"`
	Version   string             `json:"version"`
	Commit    string             `json:"commit"`
	BuildDate string             `json:"buildDate"`
	Artifacts []manifestArtifact `json:"artifacts"`
}

func main() {
	directory := flag.String("directory", "bundle", "asset directory")
	flag.Parse()
	raw, err := os.ReadFile(filepath.Join(*directory, "artifacts.json"))
	if err != nil {
		fatalf("read artifacts.json: %v", err)
	}
	var manifest releaseManifest
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		fatalf("decode artifacts.json: %v", err)
	}
	if manifest.Schema != 1 || manifest.Protocol != 1 || manifest.Name != "runx" || manifest.Version == "" {
		fatalf("artifacts.json header is incomplete: %+v", manifest)
	}
	if len(manifest.Artifacts) != 20 {
		fatalf("expected exactly 20 declared artifacts, got %d", len(manifest.Artifacts))
	}

	entries, err := os.ReadDir(*directory)
	if err != nil {
		fatalf("read assets: %v", err)
	}
	observed := map[string]bool{}
	for _, entry := range entries {
		if !entry.IsDir() {
			observed[entry.Name()] = true
		}
	}
	expected := map[string]bool{"artifacts.json": true, "checksums.txt": true}
	for _, artifact := range manifest.Artifacts {
		expected[artifact.File] = true
	}
	for name := range observed {
		if !expected[name] {
			fatalf("unexpected release asset %s", name)
		}
	}
	for name := range expected {
		if !observed[name] {
			fatalf("missing declared asset %s", name)
		}
	}

	for _, artifact := range manifest.Artifacts {
		sum, err := sha256File(filepath.Join(*directory, artifact.File))
		if err != nil {
			fatalf("hash %s: %v", artifact.File, err)
		}
		if !strings.EqualFold(sum, artifact.SHA256) {
			fatalf("manifest digest mismatch for %s", artifact.File)
		}
	}

	if err := verifyChecksums(*directory); err != nil {
		fatalf("verify checksums: %v", err)
	}
	archive, err := zip.OpenReader(filepath.Join(*directory, "guiho-s-runx.zip"))
	if err != nil {
		fatalf("open skill archive: %v", err)
	}
	defer archive.Close()
	found := false
	for _, file := range archive.File {
		if file.Name == "guiho-s-runx/SKILL.md" {
			found = true
		}
	}
	if !found {
		fatalf("skill archive is missing guiho-s-runx/SKILL.md")
	}
	fmt.Printf("verified protocol-v1 release v%s: 20 declared artifacts, manifest digests, checksums, and skill archive\n", manifest.Version)
}

func verifyChecksums(directory string) error {
	file, err := os.Open(filepath.Join(directory, "checksums.txt"))
	if err != nil {
		return err
	}
	defer file.Close()
	seen := map[string]bool{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			return fmt.Errorf("invalid checksum line %q", scanner.Text())
		}
		name := strings.TrimPrefix(fields[1], "*")
		if name == "checksums.txt" {
			return fmt.Errorf("checksums.txt must not hash itself")
		}
		input, err := os.Open(filepath.Join(directory, name))
		if err != nil {
			return err
		}
		hash := sha256.New()
		_, copyErr := io.Copy(hash, input)
		closeErr := input.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		actual := hex.EncodeToString(hash.Sum(nil))
		if !strings.EqualFold(actual, fields[0]) {
			return fmt.Errorf("checksum mismatch for %s", name)
		}
		seen[name] = true
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(seen) != 21 {
		return fmt.Errorf("expected 21 checksum entries, got %d", len(seen))
	}
	return nil
}

func sha256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func fatalf(format string, values ...any) { fmt.Fprintf(os.Stderr, format+"\n", values...); os.Exit(1) }
