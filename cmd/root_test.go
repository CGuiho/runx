package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CGuiho/runx/pkg/maintenance"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeTest(t *testing.T, cwd string, args ...string) (string, string, error) {
	return executeTestWithInput(t, cwd, "", false, args...)
}

func executeTestWithInput(t *testing.T, cwd, input string, terminal bool, args ...string) (string, string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	home := t.TempDir()
	deps := Dependencies{In: strings.NewReader(input), Out: &stdout, Err: &stderr, Getwd: func() (string, error) { return cwd, nil }, HomeDir: func() (string, error) { return home, nil }, Executable: func() (string, error) { return filepath.Join(cwd, "runx.exe"), nil }, Spawn: func(string, ...string) error { return nil }, IsTerminal: func(io.Reader) bool { return terminal }}
	root := NewRootCommand(deps, BuildInfo{Version: "1.2.3", Commit: "abc", Date: "2026-01-01T00:00:00Z", Target: "runx-windows-amd64"})
	root.SetArgs(args)
	err := root.Execute()
	if err == errDeveloperHelp {
		err = nil
	}
	return stdout.String(), stderr.String(), err
}

func writeManifest(t *testing.T, directory string) string {
	t.Helper()
	path := filepath.Join(directory, "runx.yaml")
	content := `version: "2.0.0"
namespace: "demo"
scripts:
  directory: "scripts"
commands:
  - uid: "hello-command"
    id: "hello"
    summary: "Print hello."
    description: "Print a greeting without side effects."
    command: "echo hello"
    confirm: "never"
  - group: "tools"
    summary: "Tool commands."
    commands:
      - uid: "danger-command"
        id: "danger"
        summary: "Dangerous operation."
        description: "Requires explicit approval."
        command: "echo danger"
        confirm: "always"
`
	require.NoError(t, os.MkdirAll(filepath.Join(directory, "scripts"), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestVersionAndWelcome(t *testing.T) {
	cwd := t.TempDir()
	out, _, err := executeTest(t, cwd, "--version")
	require.NoError(t, err)
	assert.Equal(t, "1.2.3\n", out)
	out, _, err = executeTest(t, cwd)
	require.NoError(t, err)
	assert.Equal(t, "Hello Windows - runx v1.2.3\n", out)
}

func TestBareInvocationBootstrapsAgentIntegrationBeforeWelcome(t *testing.T) {
	repository := t.TempDir()
	home := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(repository, ".git"), 0o755))
	for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
		require.NoError(t, os.WriteFile(filepath.Join(repository, name), []byte("# Preserve\r\n"), 0o644))
	}
	var stdout, stderr bytes.Buffer
	deps := Dependencies{In: strings.NewReader(""), Out: &stdout, Err: &stderr, Getwd: func() (string, error) { return repository, nil }, HomeDir: func() (string, error) { return home, nil }, Executable: func() (string, error) { return filepath.Join(repository, "runx.exe"), nil }, Spawn: func(string, ...string) error { return nil }}
	root := NewRootCommand(deps, BuildInfo{Version: "1.2.3", Target: "runx-windows-amd64"})
	require.NoError(t, root.Execute())
	assert.Equal(t, "Hello Windows - runx v1.2.3\n", stdout.String())
	for _, skillRoot := range []string{".agents", ".claude"} {
		content, err := os.ReadFile(filepath.Join(home, skillRoot, "skills", skillID, "SKILL.md"))
		require.NoError(t, err)
		assert.Contains(t, string(content), "name: guiho-s-runx")
	}
	for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
		content, err := os.ReadFile(filepath.Join(repository, name))
		require.NoError(t, err)
		assert.Contains(t, string(content), "# Preserve\r\n")
		assert.Contains(t, string(content), maintenance.ManagedStart)
	}

	root = NewRootCommand(deps, BuildInfo{Version: "1.2.3", Target: "runx-windows-amd64"})
	require.NoError(t, root.Execute())
	for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
		content, err := os.ReadFile(filepath.Join(repository, name))
		require.NoError(t, err)
		assert.Equal(t, 1, strings.Count(string(content), maintenance.ManagedStart))
	}
}

func TestHelpAndVersionDoNotBootstrapAgentIntegration(t *testing.T) {
	repository := t.TempDir()
	home := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(repository, ".git"), 0o755))
	for _, args := range [][]string{{"--version"}, {"--help"}, {"agent", "--help"}, {"uninstall", "--help"}} {
		var stdout, stderr bytes.Buffer
		deps := Dependencies{In: strings.NewReader(""), Out: &stdout, Err: &stderr, Getwd: func() (string, error) { return repository, nil }, HomeDir: func() (string, error) { return home, nil }, Executable: func() (string, error) { return filepath.Join(repository, "runx.exe"), nil }, Spawn: func(string, ...string) error { return nil }}
		root := NewRootCommand(deps, BuildInfo{Version: "1.2.3", Target: "runx-windows-amd64"})
		root.SetArgs(args)
		err := root.Execute()
		if err == errDeveloperHelp {
			err = nil
		}
		require.NoError(t, err)
	}
	assert.False(t, maintenance.PathExists(filepath.Join(repository, "AGENTS.md")))
	assert.False(t, maintenance.PathExists(filepath.Join(home, ".agents", "skills", skillID)))
}

func TestBareInvocationRejectsMalformedInstructionMarkers(t *testing.T) {
	repository := t.TempDir()
	home := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(repository, ".git"), 0o755))
	agents := filepath.Join(repository, "AGENTS.md")
	original := "# Preserve\n" + maintenance.ManagedStart + "\nunterminated\n"
	require.NoError(t, os.WriteFile(agents, []byte(original), 0o644))
	var stdout, stderr bytes.Buffer
	deps := Dependencies{In: strings.NewReader(""), Out: &stdout, Err: &stderr, Getwd: func() (string, error) { return repository, nil }, HomeDir: func() (string, error) { return home, nil }, Executable: func() (string, error) { return filepath.Join(repository, "runx.exe"), nil }, Spawn: func(string, ...string) error { return nil }}
	root := NewRootCommand(deps, BuildInfo{Version: "1.2.3", Target: "runx-windows-amd64"})
	err := root.Execute()
	require.Error(t, err)
	assert.Equal(t, 5, ExitCode(err))
	assert.Empty(t, stdout.String())
	content, readErr := os.ReadFile(agents)
	require.NoError(t, readErr)
	assert.Equal(t, original, string(content))
}

func TestLiveHelpTreeAndDocs(t *testing.T) {
	cwd := t.TempDir()
	out, _, err := executeTest(t, cwd, "--help-tree")
	require.NoError(t, err)
	assert.Contains(t, out, "COMMAND TREE")
	assert.Contains(t, out, "├── list")
	assert.Contains(t, out, "describe <uid-or-selector-or-index>")
	assert.Contains(t, out, "run [options] <uid-or-selector-or-index> [--] [child arguments...]")
	assert.Contains(t, out, "└── -v, --version")
	assert.NotContains(t, out, "â”œ")
	out, _, err = executeTest(t, cwd, "list", "--help-tree-depth", "1", "--help-tree")
	require.NoError(t, err)
	assert.Contains(t, out, "runx list")
	out, _, err = executeTest(t, cwd, "list", "--help-docs")
	require.NoError(t, err)
	assert.Contains(t, out, "# runx list")
	assert.Contains(t, out, "runx list --help-tree")
	assert.NotContains(t, out, "Auto generated")
	out, _, err = executeTest(t, cwd, "describe", "--help-docs")
	require.NoError(t, err)
	assert.Contains(t, out, "# runx describe")
	out, _, err = executeTest(t, cwd, "run", "--help-tree-depth", "1")
	require.NoError(t, err)
	assert.Contains(t, out, "runx run")
}

func TestOnlyApprovedShortAliases(t *testing.T) {
	root := NewRootCommand(Dependencies{Spawn: func(string, ...string) error { return nil }}, BuildInfo{Version: "1.0.0"})
	seen := map[string]string{}
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		command.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
			if flag.Shorthand != "" {
				seen[command.CommandPath()+":"+flag.Name] = flag.Shorthand
			}
		})
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
	assert.Equal(t, "v", seen["runx:version"])
	for key, value := range seen {
		if key != "runx:version" && !strings.HasSuffix(key, ":help") {
			t.Fatalf("unexpected short alias %s=%s", key, value)
		}
	}
}

func TestCatalogCommandsReadRealManifest(t *testing.T) {
	cwd := t.TempDir()
	path := writeManifest(t, cwd)
	out, _, err := executeTest(t, cwd, "check", "--format", "json")
	require.NoError(t, err)
	var checked map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &checked))
	assert.Equal(t, true, checked["valid"])
	assert.Equal(t, path, checked["configuration"])
	out, _, err = executeTest(t, cwd, "list", "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "hello-command")
	assert.NotContains(t, out, "app:hello")
	out, _, err = executeTest(t, cwd, "describe", "hello-command", "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "Print a greeting")
	out, _, err = executeTest(t, cwd, "describe", "1", "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, `"uid": "hello-command"`)
	out, _, err = executeTest(t, cwd, "run", "--dry-run", "hello-command", "--", "--hostile", "$HOME")
	require.NoError(t, err)
	assert.Contains(t, out, "[0] \"--hostile\"")
	assert.Contains(t, out, "[1] \"$HOME\"")
	out, _, err = executeTest(t, cwd, "run", "--dry-run", "1")
	require.NoError(t, err)
	assert.Contains(t, out, "uid: hello-command")
	_, _, err = executeTest(t, cwd, "run", "--dry-run", "danger-command")
	assert.Error(t, err)
	assert.Equal(t, 2, ExitCode(err))
	_, _, err = executeTest(t, cwd, "run", "--dry-run", "2")
	assert.Error(t, err)
	assert.Equal(t, 2, ExitCode(err))
	out, _, err = executeTest(t, cwd, "run", "--dry-run", "--yes", "danger-command")
	require.NoError(t, err)
	assert.Contains(t, out, "danger-command")
	out, _, err = executeTest(t, cwd, "run", "--dry-run", "--yes", "2")
	require.NoError(t, err)
	assert.Contains(t, out, "danger-command")
	for _, selector := range []string{"0", "01", "3"} {
		_, _, err = executeTest(t, cwd, "run", "--dry-run", selector)
		require.Error(t, err)
		assert.Equal(t, 3, ExitCode(err))
	}
	_, _, err = executeTest(t, cwd, "run", "--dry-run", "--", "-1")
	require.Error(t, err)
	assert.Equal(t, 3, ExitCode(err))
}

func TestRunPromptsForInteractiveConfirmation(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)
	for _, answer := range []string{"y\n", "YES\r\n"} {
		out, prompt, err := executeTestWithInput(t, cwd, answer, true, "run", "--dry-run", "danger-command")
		require.NoError(t, err)
		assert.Contains(t, prompt, "Command danger-command requires confirmation.")
		assert.Contains(t, prompt, "runx run --dry-run --yes danger-command")
		assert.Contains(t, prompt, "Are you sure? [y/N]:")
		assert.Contains(t, out, "uid: danger-command")
	}
}

func TestRunConfirmationDeclinesSafely(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)
	for _, answer := range []string{"\n", "n\n", "no\n", "unexpected\n", ""} {
		out, prompt, err := executeTestWithInput(t, cwd, answer, true, "run", "--dry-run", "danger-command")
		require.Error(t, err)
		assert.Equal(t, 2, ExitCode(err))
		assert.Contains(t, err.Error(), "command danger-command was not authorized")
		assert.Contains(t, err.Error(), "runx run --dry-run --yes danger-command")
		assert.Contains(t, prompt, "Are you sure? [y/N]:")
		assert.Empty(t, out)
	}
}

func TestRunConfirmationNeverPromptsNonInteractiveOrJSONInput(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)

	for _, test := range []struct {
		args     []string
		terminal bool
		retry    string
	}{
		{[]string{"run", "danger-command"}, false, "runx run --yes danger-command"},
		{[]string{"run", "--format", "json", "danger-command"}, true, "runx run --format json --yes danger-command"},
	} {
		out, prompt, err := executeTestWithInput(t, cwd, "yes\n", test.terminal, test.args...)
		require.Error(t, err)
		assert.Equal(t, 2, ExitCode(err))
		assert.Contains(t, err.Error(), "requires confirmation; rerun exactly:")
		assert.Contains(t, err.Error(), test.retry)
		assert.Empty(t, out)
		assert.Empty(t, prompt)
	}
}

func TestRunConfirmationRetryPreservesChildArguments(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)
	_, prompt, err := executeTestWithInput(t, cwd, "n\n", true, "run", "--dry-run", "danger-command", "--", "--force", "two words")
	require.Error(t, err)
	assert.Contains(t, prompt, `runx run --dry-run --yes danger-command -- --force "two words"`)

	out, prompt, err := executeTestWithInput(t, cwd, "", true, "run", "--dry-run", "--yes", "danger-command")
	require.NoError(t, err)
	assert.Empty(t, prompt)
	assert.Contains(t, out, "uid: danger-command")
}

func TestListTextAlignsColumns(t *testing.T) {
	cwd := t.TempDir()
	writeManifest(t, cwd)
	out, _, err := executeTest(t, cwd, "list")
	require.NoError(t, err)
	lines := strings.Split(out, "\n")
	require.GreaterOrEqual(t, len(lines), 6)
	header := lines[3]
	headers := []string{"IDX", "UID", "SELECTOR", "SUMMARY"}
	starts := make([]int, len(headers))
	for index, value := range headers {
		starts[index] = strings.Index(header, value)
		require.NotEqualf(t, -1, starts[index], "header column %s was not found", value)
	}
	expectedRows := [][]string{
		{"1", "hello-command", "hello", "Print hello."},
		{"2", "danger-command", "tools/danger", "Dangerous operation. [confirm]"},
	}
	for rowIndex, row := range lines[4:6] {
		for columnIndex, expected := range expectedRows[rowIndex] {
			end := len(row)
			if columnIndex+1 < len(starts) {
				end = starts[columnIndex+1]
			}
			require.GreaterOrEqual(t, len(row), end)
			assert.Equalf(t, expected, strings.TrimSpace(row[starts[columnIndex]:end]), "row %d column %s is not aligned", rowIndex+1, headers[columnIndex])
		}
	}
}

func TestInspectionNeverExecutesConfiguredCommand(t *testing.T) {
	cwd := t.TempDir()
	marker := filepath.Join(cwd, "spawned")
	path := writeManifest(t, cwd)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	updated := strings.Replace(string(data), `command: "echo hello"`, `command: powershell.exe -NoProfile -Command Set-Content`+" "+marker+` spawned`, 1)
	require.NoError(t, os.WriteFile(path, []byte(updated), 0o644))
	for _, args := range [][]string{{"check"}, {"list"}, {"describe", "hello-command"}, {"run", "--dry-run", "hello-command"}} {
		_, _, err := executeTest(t, cwd, args...)
		require.NoError(t, err)
		_, statErr := os.Stat(marker)
		assert.True(t, os.IsNotExist(statErr))
	}
}

func TestStrictUnknownFieldFailsWithConfigExit(t *testing.T) {
	cwd := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(cwd, "runx.yaml"), []byte("version: \"2.0.0\"\nnamespace: demo\nscripts:\n  directory: scripts\nunknown: true\ncommands: []\n"), 0o644))
	_, _, err := executeTest(t, cwd, "check", "--format", "json")
	require.Error(t, err)
	assert.Equal(t, 3, ExitCode(err))
}

func TestAgentResourcesAndInit(t *testing.T) {
	cwd := t.TempDir()
	out, _, err := executeTest(t, cwd, "agent", "skill", "install", "--local", "--cwd", cwd, "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, ".agents")
	for _, root := range []string{".agents", ".claude"} {
		_, err := os.Stat(filepath.Join(cwd, root, "skills", skillID, "SKILL.md"))
		require.NoError(t, err)
	}
	out, _, err = executeTest(t, cwd, "agent", "prompt", "list", "--names")
	require.NoError(t, err)
	assert.Equal(t, "guiho-i-runx\n", out)
	newProject := t.TempDir()
	out, _, err = executeTest(t, newProject, "init", "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "created")
	_, err = os.Stat(filepath.Join(newProject, "runx.yaml"))
	require.NoError(t, err)
}

func TestEveryPublicScopeSupportsDeveloperContext(t *testing.T) {
	cwd := t.TempDir()
	paths := [][]string{
		{}, {"list"}, {"describe"}, {"run"}, {"check"}, {"init"}, {"agent"},
		{"agent", "skill"}, {"agent", "skill", "install"}, {"agent", "skill", "uninstall"},
		{"agent", "skill", "update"}, {"agent", "skill", "list"}, {"agent", "skill", "show"},
		{"agent", "instruction"}, {"agent", "instruction", "apply"}, {"agent", "instruction", "remove"},
		{"agent", "instruction", "update"}, {"agent", "instruction", "show"}, {"agent", "prompt"},
		{"agent", "prompt", "list"}, {"agent", "prompt", "show"}, {"upgrade"}, {"upgrade", "check"},
		{"upgrade", "list"}, {"uninstall"},
	}
	for _, path := range paths {
		args := append(append([]string{}, path...), "--help-docs")
		out, _, err := executeTest(t, cwd, args...)
		require.NoErrorf(t, err, "help docs failed at %v", path)
		expected := "# runx"
		if len(path) > 0 {
			expected += " " + strings.Join(path, " ")
		}
		assert.Contains(t, out, expected)
	}
}

func TestAgentNamespaceIsIdempotent(t *testing.T) {
	cwd := t.TempDir()
	for _, action := range []string{"install", "update"} {
		_, _, err := executeTest(t, cwd, "agent", "skill", action, "--local", "--cwd", cwd)
		require.NoError(t, err)
	}
	out, _, err := executeTest(t, cwd, "agent", "skill", "list", "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, skillID)
	out, _, err = executeTest(t, cwd, "agent", "skill", "show", skillID, "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "SKILL.md")

	for _, action := range []string{"apply", "update"} {
		_, _, err = executeTest(t, cwd, "agent", "instruction", action, "--cwd", cwd)
		require.NoError(t, err)
	}
	instructionPath := filepath.Join(cwd, "AGENTS.md")
	data, err := os.ReadFile(instructionPath)
	require.NoError(t, err)
	assert.Equal(t, 1, strings.Count(string(data), maintenance.ManagedStart))
	_, _, err = executeTest(t, cwd, "agent", "instruction", "remove", "--cwd", cwd)
	require.NoError(t, err)
	data, err = os.ReadFile(instructionPath)
	require.NoError(t, err)
	assert.NotContains(t, string(data), maintenance.ManagedStart)

	out, _, err = executeTest(t, cwd, "agent", "prompt", "show", "guiho-i-runx")
	require.NoError(t, err)
	assert.Contains(t, out, "RunX Agent Instruction")
	_, _, err = executeTest(t, cwd, "agent", "skill", "uninstall", "--local", "--cwd", cwd)
	require.NoError(t, err)
	for _, root := range []string{".agents", ".claude"} {
		_, statErr := os.Stat(filepath.Join(cwd, root, "skills", skillID))
		assert.True(t, os.IsNotExist(statErr))
	}
}
