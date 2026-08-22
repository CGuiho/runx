package installstate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func bytesReader(data []byte) *bytes.Reader { return bytes.NewReader(data) }

// WriteFileAtomic writes data to a temporary sibling and renames it over
// path so readers never observe a partial file.
func WriteFileAtomic(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}
	temporary, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("stage %s: %w", path, err)
	}
	tempName := temporary.Name()
	defer func() { _ = os.Remove(tempName) }()
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("write %s: %w", tempName, err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("sync %s: %w", tempName, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close %s: %w", tempName, err)
	}
	if err := os.Chmod(tempName, mode); err != nil {
		return fmt.Errorf("chmod %s: %w", tempName, err)
	}
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("activate %s: %w", path, err)
	}
	return nil
}

// StagingDir creates a unique CLI-owned staging directory under
// $HOME/.guiho/.temp/runx-<operation>-<unique>/ and returns it together with a
// cleanup function that is safe to call multiple times.
func StagingDir(operation string) (string, func(), error) {
	root, err := TempRoot()
	if err != nil {
		return "", nil, err
	}
	if strings.ContainsAny(operation, `/\`) || operation == "." || operation == ".." || operation == "" {
		return "", nil, fmt.Errorf("invalid staging operation name %q", operation)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", nil, fmt.Errorf("create staging root %s: %w", root, err)
	}
	dir, err := os.MkdirTemp(root, "runx-"+operation+"-")
	if err != nil {
		return "", nil, fmt.Errorf("create staging under %s: %w", root, err)
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return dir, cleanup, nil
}
