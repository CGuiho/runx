package manifest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/CGuiho/runx/pkg/config"
	"go.yaml.in/yaml/v3"
)

var (
	identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	semverPattern     = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$`)
)

func ParseManifestFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest %q: %w", path, err)
	}
	return ParseManifestBytes(data)
}

func ParseManifestBytes(data []byte) (*Manifest, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode manifest YAML: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("decode manifest YAML: multiple documents are not supported")
		}
		return nil, fmt.Errorf("decode manifest YAML: %w", err)
	}
	if err := ValidateManifestSchema(&manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func ValidateManifestSchema(manifest *Manifest) error {
	if !semverPattern.MatchString(manifest.Version) || !strings.HasPrefix(manifest.Version, "2.") {
		return fmt.Errorf("unsupported manifest version %q: semantic version 2.x is required", manifest.Version)
	}
	if manifest.Agent != nil {
		if err := config.Validate(*manifest.Agent); err != nil {
			return fmt.Errorf("manifest agent policy: %w", err)
		}
	}
	if !identifierPattern.MatchString(manifest.Namespace) {
		return fmt.Errorf("manifest namespace %q must match %s", manifest.Namespace, identifierPattern)
	}
	if strings.TrimSpace(manifest.Scripts.Directory) == "" {
		return fmt.Errorf("manifest scripts.directory is required")
	}
	if filepath.IsAbs(manifest.Scripts.Directory) || !isContainedRelative(manifest.Scripts.Directory) || filepath.Clean(manifest.Scripts.Directory) == "." {
		return fmt.Errorf("manifest scripts.directory must be a relative subdirectory")
	}
	return validateEntries(manifest.Commands, manifest.Namespace, true)
}

func validateEntries(entries []Command, namespace string, topLevel bool) error {
	seen := map[string]bool{}
	for _, entry := range entries {
		leaf := entry.ID != "" || entry.UID != "" || entry.Command != ""
		group := entry.Group != ""
		if leaf == group {
			return fmt.Errorf("each commands entry must be exactly one command leaf or group")
		}
		if leaf {
			if !identifierPattern.MatchString(entry.UID) || !identifierPattern.MatchString(entry.ID) {
				return fmt.Errorf("command uid %q and id %q must use identifier-safe names", entry.UID, entry.ID)
			}
			if strings.TrimSpace(entry.Summary) == "" || strings.TrimSpace(entry.Description) == "" || strings.TrimSpace(entry.Command) == "" {
				return fmt.Errorf("command %q requires summary, description, and command", entry.UID)
			}
			if len(entry.Commands) != 0 || entry.RunX != "" {
				return fmt.Errorf("command %q cannot contain commands or runx", entry.UID)
			}
			if entry.CWD != "" && (filepath.IsAbs(entry.CWD) || !isContainedRelative(entry.CWD)) {
				return fmt.Errorf("command %q cwd must remain inside its catalog directory", entry.UID)
			}
			if entry.Shell != "" && entry.Shell != "auto" && entry.Shell != "bash" && entry.Shell != "sh" && entry.Shell != "powershell" && entry.Shell != "cmd" {
				return fmt.Errorf("command %q has unsupported shell %q", entry.UID, entry.Shell)
			}
			if entry.Confirm != "" && entry.Confirm != "never" && entry.Confirm != "always" {
				return fmt.Errorf("command %q confirm must be never or always", entry.UID)
			}
			if seen[entry.ID] {
				return fmt.Errorf("duplicate sibling command or group name %q", entry.ID)
			}
			seen[entry.ID] = true
			if topLevel && entry.ID == namespace {
				return fmt.Errorf("namespace %q conflicts with a top-level command or group", namespace)
			}
			continue
		}

		if !identifierPattern.MatchString(entry.Group) || strings.TrimSpace(entry.Summary) == "" {
			return fmt.Errorf("group %q requires an identifier-safe group and summary", entry.Group)
		}
		if seen[entry.Group] {
			return fmt.Errorf("duplicate sibling command or group name %q", entry.Group)
		}
		seen[entry.Group] = true
		if topLevel && entry.Group == namespace {
			return fmt.Errorf("namespace %q conflicts with a top-level command or group", namespace)
		}
		hasCommands := entry.Commands != nil
		hasRunX := entry.RunX != ""
		if hasCommands == hasRunX {
			return fmt.Errorf("group %q must define exactly one of commands or runx", entry.Group)
		}
		if hasCommands {
			if err := validateEntries(entry.Commands, namespace, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func isContainedRelative(value string) bool {
	clean := filepath.Clean(value)
	return clean != ".." && !strings.HasPrefix(clean, ".."+string(filepath.Separator))
}
