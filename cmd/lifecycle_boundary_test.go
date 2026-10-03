package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CGuiho/runx/pkg/maintenance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type boundaryFixture struct {
	root, project, home string
	deps                Dependencies
	workers             []string
	executableCalls     int
}

func newBoundaryFixture(t *testing.T, existingResources bool) *boundaryFixture {
	t.Helper()
	fixture := &boundaryFixture{root: t.TempDir()}
	fixture.project = filepath.Join(fixture.root, "project")
	fixture.home = filepath.Join(fixture.root, "home")
	for _, directory := range []string{filepath.Join(fixture.project, ".git"), filepath.Join(fixture.project, "nested"), fixture.home} {
		require.NoError(t, os.MkdirAll(directory, 0o755))
	}
	// These are isolated fixture values, not the caller's resource homes.
	t.Setenv("HOME", fixture.home)
	t.Setenv("USERPROFILE", fixture.home)
	t.Setenv("RUNX_DISABLE_UPDATE_WORKER", "0")
	t.Setenv("RUNX_DISABLE_AGENT_MAINTENANCE_WORKER", "0")
	writeManifest(t, fixture.project)
	manifestPath := filepath.Join(fixture.project, "runx.yaml")
	data, err := os.ReadFile(manifestPath)
	require.NoError(t, err)
	data = bytes.Replace(data, []byte("echo hello"), []byte("echo child > child-sentinel"), 1)
	require.NoError(t, os.WriteFile(manifestPath, data, 0o644))
	files := map[string]string{
		"project/README.md":                 "ordinary project prose\n",
		"project/xdocs.yaml":                "schema: 1\nai:\n  mode: auto\n",
		"home/.guiho/runx/runx.global.yaml": "agent: {}\n",
		"home/.guiho/runx/cache.json":       "stale-cache-sentinel\n",
		"home/.agents/agents/custom.md":     "unrelated agent sentinel\n",
	}
	if existingResources {
		for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
			files["project/"+name] = "# Keep\r\n" + maintenance.ManagedStart + "\r\nstale block\r\n" + maintenance.ManagedEnd + "\r\n"
		}
		for _, directory := range []string{".agents", ".claude"} {
			files["home/"+directory+"/skills/guiho-s-runx/SKILL.md"] = "stale skill sentinel\n"
			files["home/"+directory+"/skills/guiho-s-runx/custom.txt"] = "preserve skill extra\n"
		}
	}
	for name, content := range files {
		path := filepath.Join(fixture.root, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	fixture.deps = Dependencies{
		In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{},
		Getwd:   func() (string, error) { return filepath.Join(fixture.project, "nested"), nil },
		HomeDir: func() (string, error) { return fixture.home, nil },
		Executable: func() (string, error) {
			fixture.executableCalls++
			return filepath.Join(fixture.root, "runx"), nil
		},
		Spawn: func(_ string, args ...string) error {
			fixture.workers = append(fixture.workers, args[0])
			// Execute side effects synchronously at the production scheduling
			// seam: a detached no-op mock could hide the exact regression.
			require.NoError(t, os.WriteFile(filepath.Join(fixture.home, args[0]+"-sentinel"), []byte("worker ran\n"), 0o644))
			if args[0] == "__maintenance-worker" {
				skill, err := bundledSkill()
				require.NoError(t, err)
				_, err = maintenance.MaintainAgentIntegration(fixture.project, fixture.home, skill)
				require.NoError(t, err)
			} else if args[0] == "__update-worker" {
				require.NoError(t, os.WriteFile(filepath.Join(fixture.home, ".guiho", "runx", "cache.json"), []byte("worker replaced cache\n"), 0o644))
			}
			return nil
		},
	}
	return fixture
}

// Capture bytes, file/directory set, permissions and modification times, so
// identical-content rewrites, temporary leftovers and resource creation count.
func boundarySnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		value := fmt.Sprintf("%s:%d:", info.Mode(), info.ModTime().UnixNano())
		if !entry.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += string(data)
		}
		result[relative] = value
		return nil
	})
	require.NoError(t, err)
	return result
}

func TestInspectionLifecycleFilesystemBoundary(t *testing.T) {
	for _, existing := range []bool{false, true} {
		for _, test := range []struct {
			name string
			args []string
			code int
		}{
			{"check-text", []string{"check"}, 0},
			{"check-json", []string{"check", "--format", "json"}, 0},
			{"list-text", []string{"list"}, 0},
			{"list-json", []string{"list", "--format", "json"}, 0},
			{"describe-text", []string{"describe", "hello-command"}, 0},
			{"describe-json", []string{"describe", "--format", "json", "hello-command"}, 0},
			{"reveal", []string{"reveal", "hello-command"}, 0},
			{"dry-run", []string{"run", "--dry-run", "hello-command"}, 0},
			{"dry-run-json", []string{"run", "--dry-run=true", "--format", "json", "hello-command", "--", "--dry-run=false"}, 0},
			{"dry-run-gated", []string{"run", "--dry-run", "danger-command"}, 2},
			{"missing-selector", []string{"describe", "missing"}, 3},
			{"help", []string{"--help"}, 0},
			{"version", []string{"--version"}, 0},
			{"help-docs", []string{"run", "--help-docs"}, 0},
		} {
			t.Run(fmt.Sprintf("existing=%t/%s", existing, test.name), func(t *testing.T) {
				fixture := newBoundaryFixture(t, existing)
				before := boundarySnapshot(t, fixture.root)
				root := NewRootCommand(fixture.deps, BuildInfo{Version: "1.2.3"})
				args := append([]string{}, test.args...)
				if test.args[0] != "--help" && test.args[0] != "--version" && test.name != "help-docs" {
					// Keep process cwd nested and resolve the explicitly selected
					// project catalog without accidentally bootstrapping either.
					args = append([]string{args[0], "--cwd", fixture.project, "--config", "runx.yaml"}, args[1:]...)
				}
				root.SetArgs(args)
				err := root.Execute()
				if err == errDeveloperHelp {
					err = nil
				}
				require.Equal(t, test.code, ExitCode(err), "%v: %v", args, err)
				assert.Empty(t, fixture.workers, "inspection scheduled a lifecycle worker")
				assert.Zero(t, fixture.executableCalls, "inspection resolved a worker executable")
				assert.Equal(t, before, boundarySnapshot(t, fixture.root), "inspection changed the project or global-home fixture")
			})
		}
	}
}

func TestRealRunRetainsLifecycleAndChildExecution(t *testing.T) {
	for _, runArgs := range [][]string{{"hello-command"}, {"--dry-run=false", "hello-command"}, {"hello-command", "--dry-run"}} {
		t.Run(strings.Join(runArgs, " "), func(t *testing.T) {
			fixture := newBoundaryFixture(t, true)
			before := boundarySnapshot(t, fixture.root)
			root := NewRootCommand(fixture.deps, BuildInfo{Version: "1.2.3"})
			root.SetArgs(append([]string{"run", "--cwd", fixture.project}, runArgs...))
			require.NoError(t, root.Execute())
			assert.Equal(t, []string{"__update-worker", "__maintenance-worker"}, fixture.workers)
			data, err := os.ReadFile(filepath.Join(fixture.project, "child-sentinel"))
			require.NoError(t, err)
			assert.Contains(t, string(data), "child")
			assert.NotEqual(t, before, boundarySnapshot(t, fixture.root), "sentinels failed to detect real execution and worker/resource writes")
		})
	}
}
