package manifest_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CGuiho/runx/pkg/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validManifest = `version: "2.0.0"
namespace: "test-app"
scripts:
  directory: "scripts"
commands:
  - uid: "hello-command"
    id: "hello"
    summary: "Say hello."
    description: "Say hello from a test."
    command: "echo Hello"
  - group: "tools"
    summary: "Tool commands."
    commands:
      - uid: "test-command"
        id: "test"
        summary: "Run tests."
        description: "Run all tests."
        command: "go test ./..."
`

func TestStrictManifestAndIndex(t *testing.T) {
	value, err := manifest.ParseManifestBytes([]byte(validManifest))
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", value.Version)
	index, err := manifest.IndexManifest(value, "runx.yaml")
	require.NoError(t, err)
	assert.Equal(t, "echo Hello", index["hello-command"].Command)
	assert.Equal(t, "go test ./...", index["tools/test"].Command)
}
func TestUnknownAndSemanticFieldsFail(t *testing.T) {
	_, err := manifest.ParseManifestBytes([]byte(validManifest + "unknown: true\n"))
	assert.Error(t, err)
	_, err = manifest.ParseManifestBytes([]byte("version: \"1.0.0\"\nnamespace: app\nscripts:\n  directory: scripts\ncommands: []\n"))
	assert.Error(t, err)
	_, err = manifest.ParseManifestBytes([]byte("version: \"2.0.0\"\nnamespace: app\nscripts:\n  directory: ..\ncommands: []\n"))
	assert.Error(t, err)
}
func TestConfigurationPrecedenceAndNoParentSearch(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "a", "b")
	require.NoError(t, os.MkdirAll(cwd, 0o755))
	home := filepath.Join(root, "home")
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".guiho", "runx"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "runx.yaml"), []byte(validManifest), 0o644))
	global := filepath.Join(home, ".guiho", "runx", "runx.yaml")
	require.NoError(t, os.WriteFile(global, []byte(validManifest), 0o644))
	path, err := manifest.ResolveConfigPath(cwd, "", home)
	require.NoError(t, err)
	assert.Equal(t, global, path)
	local := filepath.Join(cwd, "runx.yaml")
	require.NoError(t, os.WriteFile(local, []byte(validManifest), 0o644))
	path, err = manifest.ResolveConfigPath(cwd, "", home)
	require.NoError(t, err)
	assert.Equal(t, local, path)
	explicit := filepath.Join(root, "explicit.yaml")
	require.NoError(t, os.WriteFile(explicit, []byte(validManifest), 0o644))
	path, err = manifest.ResolveConfigPath(cwd, explicit, home)
	require.NoError(t, err)
	assert.Equal(t, explicit, path)
}
func TestLocalCompositionAndReciprocity(t *testing.T) {
	directory := t.TempDir()
	parent := filepath.Join(directory, "runx.yaml")
	child := filepath.Join(directory, "child.yaml")
	parentText := `version: "2.0.0"
namespace: parent
scripts:
  directory: scripts
commands:
  - group: worker
    summary: Worker commands.
    runx: child.yaml
`
	childText := `version: "2.0.0"
namespace: child
scripts:
  directory: scripts
parent: runx.yaml
commands:
  - uid: child-build
    id: build
    summary: Build child.
    description: Build the child project.
    command: go build ./...
`
	require.NoError(t, os.WriteFile(parent, []byte(parentText), 0o644))
	require.NoError(t, os.WriteFile(child, []byte(childText), 0o644))
	catalog, err := manifest.Load(context.Background(), manifest.LoadOptions{CWD: directory})
	require.NoError(t, err)
	resolved, ok := catalog.Resolve("worker/build")
	require.True(t, ok)
	assert.Equal(t, "child-build", resolved.UID)
	resolved, ok = catalog.Resolve("1")
	require.True(t, ok)
	assert.Equal(t, "child-build", resolved.UID)
	for _, selector := range []string{"0", "-1", "01", "2"} {
		_, ok = catalog.Resolve(selector)
		assert.False(t, ok)
	}
	assert.Len(t, catalog.Children, 1)
}

func TestUIDCanEqualIDAcrossCommandsAndResolvesByUID(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "runx.yaml")
	content := `version: "2.0.0"
namespace: app
scripts:
  directory: .scripts
commands:
  - group: cli
    summary: CLI commands.
    commands:
      - uid: cli-test
        id: test
        summary: Test the CLI.
        description: Run the CLI test suite.
        command: go test ./cmd
  - group: go
    summary: Go commands.
    commands:
      - uid: test
        id: test
        summary: Test Go.
        description: Run the Go test suite.
        command: go test ./...
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	parsed, err := manifest.ParseManifestBytes([]byte(content))
	require.NoError(t, err)
	indexed, err := manifest.IndexManifest(parsed, path)
	require.NoError(t, err)
	assert.Equal(t, "test", indexed["test"].UID, "the exact UID must win over the ambiguous ID shorthand")

	catalog, err := manifest.Load(context.Background(), manifest.LoadOptions{CWD: directory})
	require.NoError(t, err)
	selected, ok := catalog.Resolve("test")
	require.True(t, ok)
	assert.Equal(t, "test", selected.UID)
	assert.Equal(t, "go/test", selected.Selector)
	selected, ok = catalog.Resolve("cli/test")
	require.True(t, ok)
	assert.Equal(t, "cli-test", selected.UID)
	selected, ok = catalog.Resolve("go/test")
	require.True(t, ok)
	assert.Equal(t, "test", selected.UID)
	selected, ok = catalog.Resolve("1")
	require.True(t, ok)
	assert.Equal(t, "cli-test", selected.UID)
	selected, ok = catalog.Resolve("2")
	require.True(t, ok)
	assert.Equal(t, "test", selected.UID)
}

func TestAmbiguousIDShorthandDoesNotResolveWithoutUID(t *testing.T) {
	directory := t.TempDir()
	content := `version: "2.0.0"
namespace: app
scripts:
  directory: scripts
commands:
  - group: cli
    summary: CLI commands.
    commands:
      - uid: cli-test
        id: test
        summary: Test the CLI.
        description: Run the CLI test suite.
        command: go test ./cmd
  - group: go
    summary: Go commands.
    commands:
      - uid: go-test
        id: test
        summary: Test Go.
        description: Run the Go test suite.
        command: go test ./...
`
	require.NoError(t, os.WriteFile(filepath.Join(directory, "runx.yaml"), []byte(content), 0o644))
	catalog, err := manifest.Load(context.Background(), manifest.LoadOptions{CWD: directory})
	require.NoError(t, err)
	_, ok := catalog.Resolve("test")
	assert.False(t, ok, "an ambiguous unqualified ID must not select an arbitrary command")
}

func TestCompositionRejectsTrueIdentityCollisions(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "duplicate UID",
			content: `version: "2.0.0"
namespace: app
scripts:
  directory: scripts
commands:
  - group: cli
    summary: CLI commands.
    commands:
      - uid: same
        id: test
        summary: Test the CLI.
        description: Run the CLI test suite.
        command: go test ./cmd
  - group: go
    summary: Go commands.
    commands:
      - uid: same
        id: verify
        summary: Verify Go.
        description: Run the Go verification suite.
        command: go test ./...
`,
		},
		{
			name: "UID collides with canonical selector",
			content: `version: "2.0.0"
namespace: app
scripts:
  directory: scripts
commands:
  - uid: top
    id: one
    summary: Top command.
    description: Run the top command.
    command: echo top
  - group: nested
    summary: Nested commands.
    commands:
      - uid: one
        id: two
        summary: Nested command.
        description: Run the nested command.
        command: echo nested
`,
		},
		{
			name: "ID collides with canonical selector",
			content: `version: "2.0.0"
namespace: app
scripts:
  directory: scripts
commands:
  - uid: top
    id: one
    summary: Top command.
    description: Run the top command.
    command: echo top
  - group: nested
    summary: Nested commands.
    commands:
      - uid: nested-command
        id: one
        summary: Nested command.
        description: Run the nested command.
        command: echo nested
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(directory, "runx.yaml"), []byte(test.content), 0o644))
			_, err := manifest.Load(context.Background(), manifest.LoadOptions{CWD: directory})
			require.Error(t, err)
		})
	}
}
