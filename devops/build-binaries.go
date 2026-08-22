//go:build ignore

// Command build-binaries compiles the complete Convention 0001 release matrix:
// eight immutable runx-payload-* executables, eight stable launchers, the
// bundled skill archive, the agent instruction, both configuration schemas,
// the artifacts.json ownership manifest, and checksums.txt covering every
// artifact except itself.
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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type target struct {
	name   string
	goos   string
	goarch string
	tuning string
}

var targets = []target{
	{"linux-amd64", "linux", "amd64", "GOAMD64=v1"}, {"linux-arm64", "linux", "arm64", "GOARM64=v8.0"},
	{"linux-armv7", "linux", "arm", "GOARM=7"}, {"linux-armv6", "linux", "arm", "GOARM=6"},
	{"darwin-amd64", "darwin", "amd64", "GOAMD64=v1"}, {"darwin-arm64", "darwin", "arm64", "GOARM64=v8.0"},
	{"windows-amd64", "windows", "amd64", "GOAMD64=v1"}, {"windows-arm64", "windows", "arm64", "GOARM64=v8.0"},
}

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

	manifest := releaseManifest{
		Schema: 1, Protocol: 1, Name: "runx",
		Version: *version, Commit: *commit, BuildDate: stamp.Format(time.RFC3339),
	}
	add := func(id, file, kind, goos, goarch, installedPath string) {
		sum, err := sha256File(filepath.Join(*output, file))
		if err != nil {
			fatalf("checksum %s: %v", file, err)
		}
		manifest.Artifacts = append(manifest.Artifacts, manifestArtifact{
			ID: id, File: file, SHA256: sum, Kind: kind,
			OS: goos, Arch: goarch, Version: *version, InstalledPath: installedPath,
		})
	}

	for _, item := range targets {
		payloadName := fmt.Sprintf("runx-payload-%s%s", item.name, exeSuffix(item.goos))
		launcherName := fmt.Sprintf("runx-launcher-%s%s", item.name, exeSuffix(item.goos))
		buildBinary(filepath.Join(*output, payloadName), item.goos, item.goarch, item.tuning, ".", *version, *commit, stamp)
		fmt.Printf("built %s\n", payloadName)
		buildBinary(filepath.Join(*output, launcherName), item.goos, item.goarch, item.tuning, "./cmd/runx-launcher", *version, *commit, stamp)
		fmt.Printf("built %s\n", launcherName)
		add("payload/"+item.name, payloadName, "payload", item.goos, item.archName(), "./.guiho/runx/versions/"+*version+"/"+payloadFileName(item.goos))
		add("launcher/"+item.name, launcherName, "launcher", item.goos, item.archName(), "./.guiho/bin/runx"+exeSuffix(item.goos))
	}

	skillZip := "guiho-s-runx.zip"
	if err := zipDirectory(filepath.Join("skills", "guiho-s-runx"), filepath.Join(*output, skillZip), stamp); err != nil {
		fatalf("create skill archive: %v", err)
	}
	add("skill/guiho-s-runx", skillZip, "skill-archive", "", "", "./.guiho/runx/resources/skills/guiho-s-runx")

	instruction := "guiho-i-runx.md"
	if err := copyFile(filepath.Join("prompts", "guiho-i-runx.md"), filepath.Join(*output, instruction)); err != nil {
		fatalf("copy instruction: %v", err)
	}
	add("instruction/guiho-i-runx", instruction, "instruction", "", "", "./.guiho/runx/resources/instruction/guiho-i-runx.md")

	for _, prompt := range []string{"runx-install.md", "runx-uninstall.md"} {
		if err := copyFile(filepath.Join("prompts", prompt), filepath.Join(*output, prompt)); err != nil {
			fatalf("copy prompt %s: %v", prompt, err)
		}
		add("prompt/"+strings.TrimSuffix(prompt, ".md"), prompt, "prompt", "", "", "./.guiho/runx/resources/prompts/"+prompt)
	}

	for _, schema := range []string{"runx.schema.json", "runx.global.schema.json"} {
		if err := copyFile(filepath.Join("schemas", schema), filepath.Join(*output, schema)); err != nil {
			fatalf("copy schema %s: %v", schema, err)
		}
		add("schema/"+strings.TrimSuffix(schema, ".json"), schema, "configuration-schema", "", "", "./.guiho/runx/resources/schemas/"+schema)
	}

	sort.Slice(manifest.Artifacts, func(i, j int) bool { return manifest.Artifacts[i].File < manifest.Artifacts[j].File })
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fatalf("encode artifacts.json: %v", err)
	}
	manifestBytes = append(manifestBytes, '\n')
	if err := os.WriteFile(filepath.Join(*output, "artifacts.json"), manifestBytes, 0o644); err != nil {
		fatalf("write artifacts.json: %v", err)
	}

	// Convention 0001: checksums.txt covers every published installation
	// artifact except itself, including the ownership manifest.
	checksumFiles := []string{"artifacts.json"}
	for _, artifact := range manifest.Artifacts {
		checksumFiles = append(checksumFiles, artifact.File)
	}
	sort.Strings(checksumFiles)
	if err := writeChecksums(filepath.Join(*output, "checksums.txt"), *output, checksumFiles); err != nil {
		fatalf("write checksums: %v", err)
	}
	expectedTotal := len(targets)*2 + 6 + 1 /*manifest*/ + 1 /*checksums*/
	fmt.Printf("release matrix complete: %d declared artifacts plus manifest and checksums (%d files)\n",
		len(manifest.Artifacts), expectedTotal)
}

func (t target) archName() string { return t.goarch }

func exeSuffix(goos string) string {
	if goos == "windows" {
		return ".exe"
	}
	return ""
}

func payloadFileName(goos string) string {
	if goos == "windows" {
		return "runx-payload.exe"
	}
	return "runx-payload"
}

func buildBinary(outputPath, goos, goarch, tuning, packagePath, version, commit string, stamp time.Time) {
	targetName := strings.TrimSuffix(filepath.Base(outputPath), ".exe")
	ldflags := strings.Join([]string{
		"-s", "-w",
		"-X", "main.version=" + version,
		"-X", "main.commit=" + commit,
		"-X", "main.buildDate=" + stamp.Format(time.RFC3339),
		"-X", "main.buildTarget=" + targetName,
	}, " ")
	command := exec.Command("go", "build", "-trimpath", "-buildvcs=true", "-ldflags", ldflags, "-o", outputPath, packagePath)
	command.Env = buildEnvironment(goos, goarch, tuning)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		fatalf("build %s: %v", targetName, err)
	}
}

func buildEnvironment(goos, goarch, tuning string) []string {
	blocked := map[string]bool{"GOOS": true, "GOARCH": true, "GOAMD64": true, "GOARM64": true, "GOARM": true, "CGO_ENABLED": true}
	result := []string{}
	for _, value := range os.Environ() {
		key, _, ok := strings.Cut(value, "=")
		if !ok || !blocked[key] {
			result = append(result, value)
		}
	}
	return append(result, "GOOS="+goos, "GOARCH="+goarch, "CGO_ENABLED=0", tuning)
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

func writeChecksums(path, outputDir string, files []string) error {
	output, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(output)
	for _, file := range files {
		sum, err := sha256File(filepath.Join(outputDir, file))
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(writer, "%s  %s\n", sum, file); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return output.Close()
}

func fatalf(format string, values ...any) { fmt.Fprintf(os.Stderr, format+"\n", values...); os.Exit(1) }
