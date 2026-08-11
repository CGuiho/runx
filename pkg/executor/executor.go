package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type ExecutionOptions struct {
	Command     string    `json:"command"`
	Args        []string  `json:"arguments"`
	CWD         string    `json:"cwd"`
	Shell       string    `json:"shell"`
	DryRun      bool      `json:"dryRun"`
	Interactive bool      `json:"interactive"`
	In          io.Reader `json:"-"`
	Out         io.Writer `json:"-"`
	Err         io.Writer `json:"-"`
}

type ExecutionResult struct {
	ExitCode int `json:"exitCode"`
}

type ShellExecution struct {
	Program string
	Args    []string
	Env     []string
	Script  string
}

type shellResolutionDeps struct {
	goos     string
	getenv   func(string) string
	lookPath func(string) (string, error)
}

func BuildShellExecution(shell, command string, args []string) (ShellExecution, error) {
	return buildShellExecution(shell, command, args, shellResolutionDeps{
		goos:     runtime.GOOS,
		getenv:   os.Getenv,
		lookPath: exec.LookPath,
	})
}

func buildShellExecution(shell, command string, args []string, deps shellResolutionDeps) (ShellExecution, error) {
	if deps.goos == "" {
		deps.goos = runtime.GOOS
	}
	if deps.getenv == nil {
		deps.getenv = os.Getenv
	}
	if deps.lookPath == nil {
		deps.lookPath = exec.LookPath
	}

	resolvedProgram := ""
	if shell == "" || shell == "auto" {
		shell, resolvedProgram = resolveAutomaticShell(deps)
	}
	switch shell {
	case "bash", "sh":
		argv := []string{"-lc", command + ` "$@"`, "runx-child"}
		argv = append(argv, args...)
		program := shell
		if resolvedProgram != "" {
			program = resolvedProgram
		}
		return ShellExecution{Program: program, Args: argv}, nil
	case "powershell":
		payload, err := json.Marshal(args)
		if err != nil {
			return ShellExecution{}, err
		}
		program := "pwsh"
		if deps.goos == "windows" {
			program = "powershell.exe"
		}
		script := `$runxForwarded = @(ConvertFrom-Json -InputObject $env:RUNX_FORWARDED_ARGUMENTS_JSON); & { ` + command + ` @runxForwarded }; exit $LASTEXITCODE`
		return ShellExecution{Program: program, Args: []string{"-NoProfile", "-NonInteractive", "-Command", script}, Env: []string{"RUNX_FORWARDED_ARGUMENTS_JSON=" + string(payload)}}, nil
	case "cmd":
		env := make([]string, 0, len(args))
		references := make([]string, 0, len(args))
		for index, value := range args {
			key := "RUNX_FORWARDED_ARGUMENT_" + strconv.Itoa(index)
			env = append(env, key+"="+value)
			references = append(references, `"%`+key+`%"`)
		}
		script := "@echo off\r\nsetlocal DisableDelayedExpansion\r\n" + strings.TrimSpace(command+" "+strings.Join(references, " ")) + "\r\nexit /b %ERRORLEVEL%\r\n"
		return ShellExecution{Program: "cmd.exe", Args: []string{"/d", "/v:off", "/s", "/c"}, Env: env, Script: script}, nil
	default:
		return ShellExecution{}, fmt.Errorf("unsupported shell %q", shell)
	}
}

func resolveAutomaticShell(deps shellResolutionDeps) (shell, program string) {
	if deps.goos != "windows" {
		return "sh", ""
	}
	if !isMSYSCaller(deps.getenv("MSYSTEM")) {
		return "cmd", ""
	}

	bash, err := deps.lookPath("bash")
	if err != nil || strings.TrimSpace(bash) == "" || isWindowsSystemBashLauncher(bash) {
		return "cmd", ""
	}
	return "bash", bash
}

func isMSYSCaller(marker string) bool {
	marker = strings.ToUpper(strings.TrimSpace(marker))
	if marker == "" {
		return false
	}
	return marker == "MSYS" || marker == "MSYS2" || strings.HasPrefix(marker, "MINGW") || strings.HasPrefix(marker, "UCRT") || strings.HasPrefix(marker, "CLANG")
}

func isWindowsSystemBashLauncher(path string) bool {
	path = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(path, "/", "\\")))
	path = strings.TrimRight(path, "\\")
	return strings.HasSuffix(path, `\windows\system32\bash.exe`) ||
		strings.HasSuffix(path, `\windows\sysnative\bash.exe`) ||
		strings.HasSuffix(path, `\windows\system32\wsl.exe`) ||
		strings.HasSuffix(path, `\windows\sysnative\wsl.exe`)
}

func ExecuteCommand(opts ExecutionOptions) (*ExecutionResult, error) {
	if opts.DryRun {
		return &ExecutionResult{ExitCode: 0}, nil
	}
	execution, err := BuildShellExecution(opts.Shell, opts.Command, opts.Args)
	if err != nil {
		return nil, err
	}
	if opts.In == nil {
		opts.In = os.Stdin
	}
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if opts.Err == nil {
		opts.Err = os.Stderr
	}

	if execution.Script != "" {
		directory := os.TempDir()
		if opts.CWD != "" {
			directory = opts.CWD
		}
		file, err := os.CreateTemp(directory, "runx-command-*.cmd")
		if err != nil {
			return nil, fmt.Errorf("create command transport script: %w", err)
		}
		path := file.Name()
		if _, err := file.WriteString(execution.Script); err != nil {
			file.Close()
			os.Remove(path)
			return nil, err
		}
		if err := file.Close(); err != nil {
			os.Remove(path)
			return nil, err
		}
		defer os.Remove(path)
		execution.Args = append(execution.Args, filepath.Clean(path))
	}

	command := exec.Command(execution.Program, execution.Args...)
	command.Dir = opts.CWD
	command.Env = append(os.Environ(), execution.Env...)
	command.Stdin, command.Stdout, command.Stderr = opts.In, opts.Out, opts.Err
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return &ExecutionResult{ExitCode: exitError.ExitCode()}, nil
		}
		return nil, fmt.Errorf("execute configured command: %w", err)
	}
	return &ExecutionResult{ExitCode: 0}, nil
}
