package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

// LoadDocument strictly decodes YAML bytes into the requested document type,
// rejecting unknown fields, then applies semantic validation.
func LoadDocument(data []byte, target any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode configuration: %w", err)
	}
	return Validate(target)
}

// LoadProject reads and validates the project configuration file.
func LoadProject(path string) (ProjectConfig, error) {
	var document ProjectConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("read project configuration %s: %w", path, err)
	}
	if err := LoadDocument(data, &document); err != nil {
		return ProjectConfig{}, err
	}
	return document, nil
}

// LoadGlobal reads and validates the global configuration file.
func LoadGlobal(path string) (GlobalConfig, error) {
	var document GlobalConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return GlobalConfig{}, fmt.Errorf("read global configuration %s: %w", path, err)
	}
	if err := LoadDocument(data, &document); err != nil {
		return GlobalConfig{}, err
	}
	return document, nil
}

// WriteProject renders a project configuration with a version-pinned schema
// reference and writes it atomically enough for a config file (create + rename
// is deferred to callers that need durability guarantees).
func WriteProject(path, version string, document ProjectConfig) error {
	document.Schema = SchemaURL(version, ProjectSchemaFile)
	if err := Validate(document); err != nil {
		return err
	}
	return writeYAML(path, document)
}

// WriteGlobal renders a global configuration with a version-pinned schema
// reference.
func WriteGlobal(path, version string, document GlobalConfig) error {
	document.Schema = SchemaURL(version, GlobalSchemaFile)
	if err := Validate(document); err != nil {
		return err
	}
	return writeYAML(path, document)
}

func writeYAML(path string, document any) error {
	buffer := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buffer)
	encoder.SetIndent(2)
	if err := encoder.Encode(document); err != nil {
		return fmt.Errorf("encode configuration: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, buffer.Bytes(), 0o644)
}
