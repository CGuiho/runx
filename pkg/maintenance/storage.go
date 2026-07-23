package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

func ReadTextIfExists(path string) (string, error) {
	if !PathExists(path) {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteTextFile(path string, content string) error {
	if err := EnsureDirectory(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func WriteTextFileAtomic(path string, content string) error {
	if err := EnsureDirectory(filepath.Dir(path)); err != nil {
		return err
	}
	tmpPath := fmt.Sprintf("%s.tmp.%d.%d", path, os.Getpid(), time.Now().UnixNano())
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

func RemovePath(path string) error {
	return os.RemoveAll(path)
}

func HomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

func GlobalRunXDirectory() string {
	return filepath.Join(HomeDirectory(), ".guiho", "runx")
}
