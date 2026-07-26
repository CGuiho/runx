//go:build ignore

package main

import (
	"archive/zip"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var expected = []string{"checksums.txt", "guiho-i-runx.md", "guiho-s-runx.zip", "runx-darwin-amd64", "runx-darwin-arm64", "runx-linux-amd64", "runx-linux-arm64", "runx-linux-armv6", "runx-linux-armv7", "runx-windows-amd64.exe", "runx-windows-arm64.exe"}

func main() {
	directory := flag.String("directory", "bundle", "asset directory")
	flag.Parse()
	entries, err := os.ReadDir(*directory)
	if err != nil {
		fatalf("read assets: %v", err)
	}
	observed := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			observed = append(observed, entry.Name())
		}
	}
	sort.Strings(observed)
	sort.Strings(expected)
	if strings.Join(observed, "\n") != strings.Join(expected, "\n") {
		fatalf("expected exactly 11 assets\nexpected: %s\nobserved: %s", strings.Join(expected, ", "), strings.Join(observed, ", "))
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
	fmt.Println("verified exactly 11 release assets and every SHA-256 checksum")
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
	if len(seen) != 10 {
		return fmt.Errorf("expected 10 checksum entries, got %d", len(seen))
	}
	return nil
}
func fatalf(format string, values ...any) { fmt.Fprintf(os.Stderr, format+"\n", values...); os.Exit(1) }
