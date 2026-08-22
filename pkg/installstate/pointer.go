package installstate

import (
	"fmt"
	"regexp"
	"strings"
)

var versionPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`)

func isValidVersionName(version string) bool {
	return versionPattern.MatchString(version)
}

// Pointer is the strictly decoded, atomically replaced active-version pointer.
type Pointer struct {
	Protocol int    `json:"protocol"`
	Active   string `json:"active"`
	Previous string `json:"previous,omitempty"`
}

// Validate enforces the pointer contract: protocol 1 only, valid SemVer
// versions, and a distinct retained fallback.
func (p Pointer) Validate() error {
	if p.Protocol != ProtocolVersion {
		return fmt.Errorf("unsupported launcher protocol %d: want %d", p.Protocol, ProtocolVersion)
	}
	if !isValidVersionName(p.Active) {
		return fmt.Errorf("pointer active version %q is not a valid SemVer", p.Active)
	}
	if p.Previous != "" {
		if !isValidVersionName(p.Previous) {
			return fmt.Errorf("pointer previous version %q is not a valid SemVer", p.Previous)
		}
		if p.Previous == p.Active {
			return fmt.Errorf("pointer previous version must differ from active")
		}
	}
	return nil
}

// VersionDirName returns the on-disk directory name for an activated version.
func VersionDirName(version string) string { return version }

// SanitizeVersion rejects path separators and traversal before any use as a
// filesystem component.
func SanitizeVersion(version string) (string, error) {
	if strings.ContainsAny(version, `/\`) || strings.Contains(version, "..") || !isValidVersionName(version) {
		return "", fmt.Errorf("unsafe or invalid version %q", version)
	}
	return version, nil
}
