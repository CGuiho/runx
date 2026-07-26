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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type target struct{ name, goos, goarch, tuning string }

var targets = []target{
	{"runx-linux-amd64", "linux", "amd64", "GOAMD64=v1"}, {"runx-linux-arm64", "linux", "arm64", "GOARM64=v8.0"},
	{"runx-linux-armv7", "linux", "arm", "GOARM=7"}, {"runx-linux-armv6", "linux", "arm", "GOARM=6"},
	{"runx-darwin-amd64", "darwin", "amd64", "GOAMD64=v1"}, {"runx-darwin-arm64", "darwin", "arm64", "GOARM64=v8.0"},
	{"runx-windows-amd64.exe", "windows", "amd64", "GOAMD64=v1"}, {"runx-windows-arm64.exe", "windows", "arm64", "GOARM64=v8.0"},
}

func main() {
	version := flag.String("version", "", "semantic version embedded in every binary")
	commit := flag.String("commit", "unknown", "source commit")
	buildDate := flag.String("build-date", time.Now().UTC().Format(time.RFC3339), "stable RFC3339 timestamp")
	output := flag.String("output", "bundle", "output directory")
	flag.Parse()
	if *version == "" {
		fatalf("--version is required")
	}
	stamp, err := time.Parse(time.RFC3339, *buildDate)
	if err != nil {
		fatalf("parse --build-date: %v", err)
	}
	if err := os.RemoveAll(*output); err != nil {
		fatalf("clean output: %v", err)
	}
	if err := os.MkdirAll(*output, 0o755); err != nil {
		fatalf("create output: %v", err)
	}
	assets := []string{}
	for _, item := range targets {
		path := filepath.Join(*output, item.name)
		targetName := strings.TrimSuffix(item.name, ".exe")
		ldflags := strings.Join([]string{"-s", "-w", "-X", "main.version=" + *version, "-X", "main.commit=" + *commit, "-X", "main.buildDate=" + *buildDate, "-X", "main.buildTarget=" + targetName}, " ")
		command := exec.Command("go", "build", "-trimpath", "-buildvcs=true", "-ldflags", ldflags, "-o", path, ".")
		command.Env = buildEnvironment(item)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		fmt.Printf("building %s\n", item.name)
		if err := command.Run(); err != nil {
			fatalf("build %s: %v", item.name, err)
		}
		assets = append(assets, path)
	}
	skill := filepath.Join(*output, "guiho-s-runx.zip")
	if err := zipDirectory(filepath.Join("skills", "guiho-s-runx"), skill, stamp); err != nil {
		fatalf("create skill archive: %v", err)
	}
	assets = append(assets, skill)
	instruction := filepath.Join(*output, "guiho-i-runx.md")
	if err := copyFile(filepath.Join("prompts", "guiho-i-runx.md"), instruction); err != nil {
		fatalf("copy instruction: %v", err)
	}
	assets = append(assets, instruction)
	sort.Strings(assets)
	if err := writeChecksums(filepath.Join(*output, "checksums.txt"), assets); err != nil {
		fatalf("write checksums: %v", err)
	}
	if len(assets)+1 != 11 {
		fatalf("release must contain 11 artifacts, got %d", len(assets)+1)
	}
	fmt.Println("release matrix complete: 8 binaries and 3 supporting artifacts")
}

func buildEnvironment(item target) []string {
	blocked := map[string]bool{"GOOS": true, "GOARCH": true, "GOAMD64": true, "GOARM64": true, "GOARM": true, "CGO_ENABLED": true}
	result := []string{}
	for _, value := range os.Environ() {
		key, _, ok := strings.Cut(value, "=")
		if !ok || !blocked[key] {
			result = append(result, value)
		}
	}
	return append(result, "GOOS="+item.goos, "GOARCH="+item.goarch, "CGO_ENABLED=0", item.tuning)
}
func zipDirectory(source, destination string, stamp time.Time) error {
	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	archive := zip.NewWriter(output)
	files := []string{}
	if err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Strings(files)
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relative)
		header.Method = zip.Deflate
		header.Modified = stamp
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		input, err := os.Open(path)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(writer, input)
		closeErr := input.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
	if err := archive.Close(); err != nil {
		return err
	}
	return output.Close()
}
func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		return err
	}
	return output.Close()
}
func writeChecksums(path string, assets []string) error {
	output, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(output)
	for _, asset := range assets {
		input, err := os.Open(asset)
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
		if _, err := fmt.Fprintf(writer, "%s  %s\n", hex.EncodeToString(hash.Sum(nil)), filepath.Base(asset)); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return output.Close()
}
func fatalf(format string, values ...any) { fmt.Fprintf(os.Stderr, format+"\n", values...); os.Exit(1) }
