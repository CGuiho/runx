// Package launcher implements stable-launcher dispatch: resolve the committed
// active payload from the current pointer and delegate to it exactly.
package launcher

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CGuiho/runx/pkg/installstate"
)

// PayloadPath resolves the executable path of the active payload named by the
// pointer. It fails closed when the pointer is missing or the payload does not
// exist.
func PayloadPath(pointer *installstate.Pointer) (string, error) {
	if pointer == nil {
		return "", fmt.Errorf("no active RunX pointer: run the installer or upgrade")
	}
	version, err := installstate.SanitizeVersion(pointer.Active)
	if err != nil {
		return "", err
	}
	dir, err := installstate.VersionDir(version)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, installstate.PayloadName())
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("active payload %s is missing: %w", path, err)
	}
	return path, nil
}

// FallbackPath resolves the retained rollback payload when the pointer records
// a previous version. Returns "" when no fallback exists.
func FallbackPath(pointer *installstate.Pointer) (string, error) {
	if pointer == nil || pointer.Previous == "" {
		return "", nil
	}
	version, err := installstate.SanitizeVersion(pointer.Previous)
	if err != nil {
		return "", err
	}
	dir, err := installstate.VersionDir(version)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, installstate.PayloadName())
	if _, err := os.Stat(path); err != nil {
		return "", nil
	}
	return path, nil
}
