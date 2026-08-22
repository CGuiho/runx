package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const allProceedInput = "n\nalways-proceed\nalways-proceed\nalways-proceed\nalways-proceed\n"

func TestInitNoninteractiveWithoutPolicyFailsClosed(t *testing.T) {
	cwd := t.TempDir()
	_, _, err := executeTest(t, cwd, "init", "--format", "json")
	require.Error(t, err)
	assert.Equal(t, 3, ExitCode(err))
}

func TestInitInteractiveWritesGlobalPolicyAndIsIdempotent(t *testing.T) {
	cwd := t.TempDir()
	home := t.TempDir()
	out, _, err := executeTestFull(t, cwd, home, allProceedInput, true, "init", "--format", "json")
	require.NoError(t, err)

	var document struct {
		Status string `json:"status"`
		Policy struct {
			Upgrade string `json:"upgrade"`
			Issues  struct {
				Bugs         string `json:"bugs"`
				Improvements string `json:"improvements"`
				Reviews      string `json:"reviews"`
			} `json:"issues"`
		} `json:"policy"`
		Global  string `json:"global"`
		Project string `json:"project"`
		Report  struct {
			Created []string `json:"created"`
		} `json:"report"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &document))
	assert.Equal(t, "always-proceed", document.Policy.Upgrade)
	assert.Equal(t, "always-proceed", document.Policy.Issues.Bugs)
	assert.Equal(t, "always-proceed", document.Policy.Issues.Improvements)
	assert.Equal(t, "always-proceed", document.Policy.Issues.Reviews)

	raw, err := os.ReadFile(document.Global)
	require.NoError(t, err)
	text := string(raw)
	assert.Contains(t, text, "upgrade: always-proceed")
	assert.Contains(t, text, "runx.global.schema.json")

	// Second run must be idempotent: nothing created, no questions asked.
	out2, _, err := executeTestFull(t, cwd, home, "", false, "init", "--format", "json")
	require.NoError(t, err)
	var document2 map[string]any
	require.NoError(t, json.Unmarshal([]byte(out2), &document2))
	assert.Equal(t, "initialized", document2["status"])
	created2, _ := document2["report"].(map[string]any)["created"].([]any)
	assert.Empty(t, created2, "second init must create nothing")
}

func TestInitHonorsManifestAgentPolicy(t *testing.T) {
	cwd := t.TempDir()
	policyManifest := "version: \"2.0.0\"\nnamespace: demo\nscripts:\n  directory: scripts\nagent:\n  evolution:\n    upgrade: always-proceed\n    issues:\n      bugs: disabled\n      improvements: disabled\n      reviews: disabled\ncommands: []\n"
	require.NoError(t, os.WriteFile(filepath.Join(cwd, "runx.yaml"), []byte(policyManifest), 0o644))
	out, _, err := executeTest(t, cwd, "init", "--format", "json")
	require.NoError(t, err)
	var document struct {
		Policy struct {
			Upgrade string `json:"upgrade"`
			Issues  struct {
				Bugs string `json:"bugs"`
			} `json:"issues"`
		} `json:"policy"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &document))
	assert.Equal(t, "always-proceed", document.Policy.Upgrade)
	assert.Equal(t, "disabled", document.Policy.Issues.Bugs)
}

func TestInitRejectsInvalidPolicyValues(t *testing.T) {
	cwd := t.TempDir()
	bad := "version: \"2.0.0\"\nnamespace: demo\nscripts:\n  directory: scripts\nagent:\n  evolution:\n    upgrade: yes-please\n    issues:\n      bugs: always-ask\n      improvements: always-ask\n      reviews: always-ask\ncommands: []\n"
	require.NoError(t, os.WriteFile(filepath.Join(cwd, "runx.yaml"), []byte(bad), 0o644))
	_, _, err := executeTest(t, cwd, "check", "--format", "json")
	require.Error(t, err)
	assert.True(t, strings.Contains(ExitErrorText(err), "invalid") || ExitCode(err) == 3, "strict decoding must reject invalid policy values")
}

func ExitErrorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
