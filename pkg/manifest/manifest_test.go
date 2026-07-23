package manifest_test

import (
	"testing"

	"github.com/CGuiho/runx/pkg/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseManifestValidYAML(t *testing.T) {
	yamlData := []byte(`
version: "2"
namespace: "test-app"
description: "Test manifest"

commands:
  - uid: "test:hello"
    id: "hello"
    description: "Say hello"
    run: "echo Hello"
`)

	m, err := manifest.ParseManifestBytes(yamlData)
	require.NoError(t, err)
	assert.Equal(t, "2", m.Version)
	assert.Equal(t, "test-app", m.Namespace)
	assert.Len(t, m.Commands, 1)

	index, err := manifest.IndexManifest(m, "runx.yaml")
	require.NoError(t, err)
	assert.Contains(t, index, "test-app:hello")
	assert.Equal(t, "echo Hello", index["test-app:hello"].Run)
}

func TestParseManifestInvalidUnknownField(t *testing.T) {
	yamlData := []byte(`
version: "2"
namespace: "test-app"
unknown_field: "invalid"

commands:
  - id: "hello"
    run: "echo Hello"
`)

	_, err := manifest.ParseManifestBytes(yamlData)
	assert.Error(t, err, "Strict YAML parser must reject unknown_field")
}

func TestParseManifestUnsupportedVersion(t *testing.T) {
	yamlData := []byte(`
version: "1"
namespace: "test-app"

commands:
  - id: "hello"
    run: "echo Hello"
`)

	_, err := manifest.ParseManifestBytes(yamlData)
	assert.Error(t, err, "Version 1 should be rejected")
}
