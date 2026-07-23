package manifest

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseManifestFile reads and parses a runx.yaml manifest file with strict field checking.
func ParseManifestFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file %s: %w", path, err)
	}

	return ParseManifestBytes(data)
}

// ParseManifestBytes decodes raw YAML bytes into a Manifest struct, enforcing strict known fields.
func ParseManifestBytes(data []byte) (*Manifest, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var m Manifest
	if err := decoder.Decode(&m); err != nil {
		return nil, fmt.Errorf("manifest YAML syntax or schema error: %w", err)
	}

	if err := ValidateManifestSchema(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// ValidateManifestSchema enforces semantic Manifest V2 validation rules.
func ValidateManifestSchema(m *Manifest) error {
	if m.Version != "2" {
		return fmt.Errorf("unsupported manifest version '%s': only version '2' is supported", m.Version)
	}
	if m.Namespace == "" {
		return fmt.Errorf("manifest requires a non-empty 'namespace' field")
	}
	if len(m.Commands) == 0 {
		return fmt.Errorf("manifest commands array cannot be empty")
	}

	return nil
}
