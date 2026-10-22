package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertEntryParity runs the same tail arguments under two entry points and
// requires byte-identical stdout, stderr, and exit behavior.
func assertEntryParity(t *testing.T, cwd, home, first, second string, input string, terminal bool, tail ...string) {
	t.Helper()
	firstOut, firstStderr, firstErr := executeTestFull(t, cwd, home, input, terminal, append([]string{first}, tail...)...)
	secondOut, secondStderr, secondErr := executeTestFull(t, cwd, home, input, terminal, append([]string{second}, tail...)...)
	assert.Equal(t, secondOut, firstOut, "stdout differs for %s %v", first, tail)
	assert.Equal(t, secondStderr, firstStderr, "stderr differs for %s %v", first, tail)
	if secondErr == nil {
		assert.NoError(t, firstErr, "error differs for %s %v", first, tail)
	} else {
		require.Error(t, firstErr, "error differs for %s %v", first, tail)
		assert.Equal(t, secondErr.Error(), firstErr.Error(), "error differs for %s %v", first, tail)
		assert.Equal(t, ExitCode(secondErr), ExitCode(firstErr), "exit code differs for %s %v", first, tail)
	}
}

func TestSmartBareMatchesListOutput(t *testing.T) {
	cwd := t.TempDir()
	path := writeManifest(t, cwd)
	home := t.TempDir()
	other := t.TempDir()
	for _, tail := range [][]string{
		{},
		{"--format", "json"},
		{"--verbose"},
		{"--format", "json", "--verbose"},
		{"--cwd", cwd},
		{"--config", path},
		{"--cwd", cwd, "--config", path, "--format", "json"},
		{"--cwd", other},
	} {
		assertEntryParity(t, cwd, home, "smart", "list", "", false, tail...)
	}
}

func TestSmartSelectorMatchesRunBehavior(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)
	home := t.TempDir()
	for _, tail := range [][]string{
		{"--dry-run", "hello-command"},
		{"--dry-run", "hello-command", "--", "--hostile", "$HOME"},
		{"--dry-run", "--format", "json", "hello-command", "--", "dry-run"},
		{"--dry-run", "--format", "json", "hello-command"},
		{"--dry-run", "--yes", "danger-command"},
		{"--dry-run", "danger-command"},
		{"hello-command"},
		{"hello-command", "--", "extra argument"},
		{"--yes", "danger-command"},
	} {
		assertEntryParity(t, cwd, home, "smart", "run", "", false, tail...)
	}
	assertEntryParity(t, cwd, home, "smart", "run", "y\n", true, "--dry-run", "danger-command")
	assertEntryParity(t, cwd, home, "smart", "run", "n\n", true, "--dry-run", "danger-command", "--", "--force", "two words")
}

func TestSmartPreservesChildExitCodes(t *testing.T) {
	cwd := t.TempDir()
	writeManifestWithCommand(t, cwd, "exit 7")
	home := t.TempDir()
	assertEntryParity(t, cwd, home, "smart", "run", "", false, "hello-command")
	_, _, err := executeTestFull(t, cwd, home, "", false, "smart", "hello-command")
	require.Error(t, err)
	assert.Equal(t, 7, ExitCode(err))
}

func TestSmartFlagsStopAtFirstPositional(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)
	// --format after the selector is a child argument, not a smart flag.
	out, _, err := executeTest(t, cwd, "smart", "--dry-run", "hello-command", "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "uid: hello-command")
	assert.Contains(t, out, "[0] \"--format\"")
	assert.Contains(t, out, "[1] \"json\"")
	// A delimiter is consumed exactly like run's child forwarding.
	out, _, err = executeTest(t, cwd, "smart", "--dry-run", "hello-command", "--", "dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "[0] \"dry-run\"")
	// Flags before the selector still parse natively.
	out, _, err = executeTest(t, cwd, "smart", "--format", "json", "--dry-run", "hello-command")
	require.NoError(t, err)
	assert.Contains(t, out, `"uid": "hello-command"`)
}

func TestSmartBareWithFlagsStillLists(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)
	out, _, err := executeTest(t, cwd, "smart", "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "hello-command")
	assert.Contains(t, out, `"commands"`)
	assert.NotContains(t, out, `"arguments"`)
	out, _, err = executeTest(t, cwd, "smart", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "RunX commands")
}

func TestSmartHelpDescribesSmartCommand(t *testing.T) {
	cwd := t.TempDir()
	for _, args := range [][]string{{"smart", "--help"}, {"smart", "-h"}} {
		out, _, err := executeTest(t, cwd, args...)
		require.NoError(t, err)
		assert.Contains(t, out, "runx smart")
		assert.Contains(t, out, "With no selector, list commands in the RunX configuration exactly like runx list.")
	}
	out, _, err := executeTest(t, cwd, "smart", "--help-tree", "--help-tree-depth", "1")
	require.NoError(t, err)
	assert.Contains(t, out, "COMMAND TREE")
	assert.Contains(t, out, "runx smart")
	assert.NotContains(t, out, "runx run")
	out, _, err = executeTest(t, cwd, "smart", "--help-docs")
	require.NoError(t, err)
	assert.Contains(t, out, "# runx smart")
	assert.Contains(t, out, "runx smart --help-tree")
	for _, args := range [][]string{{"smart", "--version"}, {"smart", "-v"}} {
		out, _, err := executeTest(t, cwd, args...)
		require.Error(t, err)
		assert.Equal(t, 2, ExitCode(err))
		assert.Empty(t, out)
	}
}

func TestSmartInspectionBoundary(t *testing.T) {
	for _, test := range []struct {
		name     string
		args     []string
		code     int
		workers  []string
		executes bool
	}{
		{"bare-list", []string{"smart", "--cwd", "__project__", "--config", "runx.yaml"}, 0, nil, false},
		{"list-json", []string{"smart", "--cwd", "__project__", "--config", "runx.yaml", "--format", "json"}, 0, nil, false},
		{"dry-run", []string{"smart", "--cwd", "__project__", "--config", "runx.yaml", "--dry-run", "hello-command"}, 0, nil, false},
		{"dry-run-gated", []string{"smart", "--cwd", "__project__", "--config", "runx.yaml", "--dry-run", "danger-command"}, 2, nil, false},
		{"real-run", []string{"smart", "--cwd", "__project__", "--config", "runx.yaml", "hello-command"}, 0, []string{"__update-worker", "__maintenance-worker"}, true},
		{"real-run-explicit", []string{"smart", "--cwd", "__project__", "--config", "runx.yaml", "--dry-run=false", "hello-command"}, 0, []string{"__update-worker", "__maintenance-worker"}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newBoundaryFixture(t, true)
			before := boundarySnapshot(t, fixture.root)
			root := NewRootCommand(fixture.deps, BuildInfo{Version: "1.2.3"})
			args := append([]string{}, test.args...)
			for index := range args {
				if args[index] == "__project__" {
					args[index] = fixture.project
				}
			}
			root.SetArgs(args)
			err := root.Execute()
			if err == errDeveloperHelp {
				err = nil
			}
			require.Equal(t, test.code, ExitCode(err), "%v: %v", args, err)
			assert.Equal(t, test.workers, fixture.workers, "unexpected lifecycle workers for %v", args)
			if !test.executes {
				assert.Zero(t, fixture.executableCalls, "inspection resolved a worker executable")
				assert.Equal(t, before, boundarySnapshot(t, fixture.root), "inspection changed the project or global-home fixture")
			} else {
				data, err := readChildSentinel(t, fixture.project)
				require.NoError(t, err)
				assert.Contains(t, data, "child")
			}
		})
	}
}

func readChildSentinel(t *testing.T, project string) (string, error) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(project, "child-sentinel"))
	return string(content), err
}

func writeManifestWithCommand(t *testing.T, directory, command string) string {
	t.Helper()
	path := writeManifest(t, directory)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	updated := strings.Replace(string(data), `command: "echo hello"`, `command: "`+command+`"`, 1)
	require.NotEqual(t, string(data), updated)
	require.NoError(t, os.WriteFile(path, []byte(updated), 0o644))
	return path
}
