package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildShellExecutionWindowsAutoUsesResolvedGitBash(t *testing.T) {
	lookupCalls := 0
	deps := shellResolutionDeps{
		goos: "windows",
		getenv: func(name string) string {
			if name == "MSYSTEM" {
				return "MINGW64"
			}
			return ""
		},
		lookPath: func(name string) (string, error) {
			lookupCalls++
			assert.Equal(t, "bash", name)
			return `C:\Program Files\Git\usr\bin\bash.exe`, nil
		},
	}

	execution, err := buildShellExecution("auto", "node path-probe", []string{"/c/GUIHO"}, deps)

	require.NoError(t, err)
	assert.Equal(t, `C:\Program Files\Git\usr\bin\bash.exe`, execution.Program)
	assert.Equal(t, []string{"-lc", `node path-probe "$@"`, "runx-child", "/c/GUIHO"}, execution.Args)
	assert.Empty(t, execution.Env)
	assert.Empty(t, execution.Script)
	assert.Equal(t, 1, lookupCalls)
}

func TestBuildShellExecutionWindowsAutoFallsBackToCmdWithoutMSYSCaller(t *testing.T) {
	lookupCalled := false
	deps := shellResolutionDeps{
		goos:   "windows",
		getenv: func(string) string { return "" },
		lookPath: func(string) (string, error) {
			lookupCalled = true
			return "", nil
		},
	}

	execution, err := buildShellExecution("auto", "echo ready", []string{"/c/GUIHO"}, deps)

	require.NoError(t, err)
	assert.Equal(t, "cmd.exe", execution.Program)
	assert.Equal(t, []string{"/d", "/v:off", "/s", "/c"}, execution.Args)
	assert.Contains(t, execution.Script, `"%RUNX_FORWARDED_ARGUMENT_0%"`)
	assert.Contains(t, execution.Script, "echo ready")
	assert.False(t, lookupCalled)
}

func TestBuildShellExecutionWindowsAutoFallsBackWhenBashCannotBeResolved(t *testing.T) {
	tests := []struct {
		name     string
		lookPath func(string) (string, error)
	}{
		{
			name: "missing",
			lookPath: func(string) (string, error) {
				return "", errors.New("bash not found")
			},
		},
		{
			name: "empty result",
			lookPath: func(string) (string, error) {
				return "", nil
			},
		},
		{
			name: "system32 launcher",
			lookPath: func(string) (string, error) {
				return `C:\Windows\System32\bash.exe`, nil
			},
		},
		{
			name: "sysnative launcher",
			lookPath: func(string) (string, error) {
				return `C:/Windows/Sysnative/bash.exe`, nil
			},
		},
		{
			name: "system32 WSL launcher",
			lookPath: func(string) (string, error) {
				return `C:\Windows\System32\wsl.exe`, nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deps := shellResolutionDeps{
				goos:     "windows",
				getenv:   func(string) string { return "UCRT64" },
				lookPath: test.lookPath,
			}

			execution, err := buildShellExecution("auto", "echo ready", nil, deps)

			require.NoError(t, err)
			assert.Equal(t, "cmd.exe", execution.Program)
			assert.Equal(t, []string{"/d", "/v:off", "/s", "/c"}, execution.Args)
		})
	}
}

func TestBuildShellExecutionExplicitShellsBypassWindowsAutoInference(t *testing.T) {
	shells := []struct {
		name       string
		expected   string
		wantScript bool
	}{
		{name: "cmd", expected: "cmd.exe", wantScript: true},
		{name: "powershell", expected: "powershell.exe"},
		{name: "bash", expected: "bash"},
		{name: "sh", expected: "sh"},
	}

	for _, test := range shells {
		t.Run(test.name, func(t *testing.T) {
			lookupCalled := false
			deps := shellResolutionDeps{
				goos:   "windows",
				getenv: func(string) string { return "MINGW64" },
				lookPath: func(string) (string, error) {
					lookupCalled = true
					return `C:\Program Files\Git\usr\bin\bash.exe`, nil
				},
			}

			execution, err := buildShellExecution(test.name, "echo ready", []string{"/c/GUIHO"}, deps)

			require.NoError(t, err)
			assert.Equal(t, test.expected, execution.Program)
			assert.Equal(t, test.wantScript, execution.Script != "")
			assert.False(t, lookupCalled)
		})
	}
}

func TestBuildShellExecutionNonWindowsAutoRemainsSh(t *testing.T) {
	lookupCalled := false
	deps := shellResolutionDeps{
		goos:   "linux",
		getenv: func(string) string { return "MINGW64" },
		lookPath: func(string) (string, error) {
			lookupCalled = true
			return "/usr/bin/bash", nil
		},
	}

	execution, err := buildShellExecution("auto", "echo ready", nil, deps)

	require.NoError(t, err)
	assert.Equal(t, "sh", execution.Program)
	assert.Equal(t, []string{"-lc", `echo ready "$@"`, "runx-child"}, execution.Args)
	assert.False(t, lookupCalled)
}

func TestIsMSYSCallerRecognizesGitBashMarkers(t *testing.T) {
	tests := map[string]bool{
		"MSYS":       true,
		"MSYS2":      true,
		"MINGW32":    true,
		"MINGW64":    true,
		"UCRT64":     true,
		"CLANG64":    true,
		"CLANGARM64": true,
		"":           false,
		"Windows_NT": false,
		"WSL_DISTRO": false,
	}

	for marker, expected := range tests {
		assert.Equal(t, expected, isMSYSCaller(marker), marker)
	}
}
