package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveDefaultsToAlwaysAsk(t *testing.T) {
	resolved, err := Resolve(ProjectConfig{}, GlobalConfig{})
	require.NoError(t, err)
	assert.Equal(t, DefaultPolicy, resolved.Upgrade)
	assert.Equal(t, DefaultPolicy, resolved.Issues.Bugs)
	assert.Equal(t, DefaultPolicy, resolved.Issues.Improvements)
	assert.Equal(t, DefaultPolicy, resolved.Issues.Reviews)
}

func TestResolveProjectOverridesGlobalPerField(t *testing.T) {
	global := GlobalConfig{Agent: Agent{Evolution: Evolution{
		Upgrade: PolicyAlwaysProceed,
		Issues:  IssuesPolicy{Bugs: PolicyAlwaysProceed, Improvements: PolicyDisabled, Reviews: PolicyAlwaysAsk},
	}}}
	project := ProjectConfig{Agent: Agent{Evolution: Evolution{
		Upgrade: PolicyDisabled,
		Issues:  IssuesPolicy{Bugs: PolicyAlwaysAsk},
	}}}
	resolved, err := Resolve(project, global)
	require.NoError(t, err)
	assert.Equal(t, PolicyDisabled, resolved.Upgrade)             // project wins
	assert.Equal(t, PolicyAlwaysAsk, resolved.Issues.Bugs)        // project wins
	assert.Equal(t, PolicyDisabled, resolved.Issues.Improvements) // global inherited
	assert.Equal(t, PolicyAlwaysAsk, resolved.Issues.Reviews)     // global inherited
}

func TestResolveRejectsUnknownPolicyValues(t *testing.T) {
	for _, invalid := range []string{"true", "yes", "always", "ALWAYS-ASK", "maybe"} {
		_, err := Resolve(ProjectConfig{Agent: Agent{Evolution: Evolution{Upgrade: invalid}}}, GlobalConfig{})
		require.Error(t, err, "value %q must fail", invalid)
	}
}

func TestLoadDocumentRejectsUnknownFields(t *testing.T) {
	data := []byte("agent:\n  evolution:\n    upgrade: always-proceed\n    unknownField: x\n")
	var document ProjectConfig
	err := LoadDocument(data, &document)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknownField")
}

func TestLoadDocumentRoundTrip(t *testing.T) {
	document := ProjectConfig{Agent: Agent{Evolution: Evolution{
		Upgrade: PolicyAlwaysProceed,
		Issues:  IssuesPolicy{Bugs: PolicyDisabled, Improvements: PolicyAlwaysAsk, Reviews: PolicyAlwaysProceed},
	}}}
	path := filepath.Join(t.TempDir(), "runx.yaml")
	require.NoError(t, WriteProject(path, "1.2.3", document))
	loaded, err := LoadProject(path)
	require.NoError(t, err)
	assert.Equal(t, document.Agent, loaded.Agent)
	assert.Contains(t, loaded.Schema, "/releases/download/v1.2.3/"+ProjectSchemaFile)
}

func TestWriteGlobalPinsGlobalSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "runx.global.yaml")
	require.NoError(t, WriteGlobal(path, "0.13.0", GlobalConfig{}))
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(raw), "https://github.com/CGuiho/runx/releases/download/v0.13.0/runx.global.schema.json")
}
