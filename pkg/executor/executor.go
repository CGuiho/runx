package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ExecutionOptions specifies options for executing a catalog command.
type ExecutionOptions struct {
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	CWD         string   `json:"cwd"`
	DryRun      bool     `json:"dry_run"`
	Interactive bool     `json:"interactive"`
}

// ExecutionResult captures execution status.
type ExecutionResult struct {
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output,omitempty"`
}

// ExecuteCommand runs a command string with forwarded arguments based on host OS.
func ExecuteCommand(opts ExecutionOptions) (*ExecutionResult, error) {
	if opts.DryRun {
		payload, _ := json.MarshalIndent(opts, "", "  ")
		fmt.Println(string(payload))
		return &ExecutionResult{ExitCode: 0}, nil
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = buildPowerShellCommand(opts.Command, opts.Args)
	} else {
		cmd = buildPOSIXCommand(opts.Command, opts.Args)
	}

	if opts.CWD != "" {
		cmd.Dir = opts.CWD
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return &ExecutionResult{ExitCode: exitErr.ExitCode()}, nil
		}
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	return &ExecutionResult{ExitCode: 0}, nil
}

func buildPOSIXCommand(commandStr string, args []string) *exec.Cmd {
	fullCmd := commandStr
	if len(args) > 0 {
		escapedArgs := make([]string, len(args))
		for i, a := range args {
			escapedArgs[i] = fmt.Sprintf("%q", a)
		}
		fullCmd = fmt.Sprintf("%s %s", commandStr, strings.Join(escapedArgs, " "))
	}
	return exec.Command("sh", "-c", fullCmd)
}

func buildPowerShellCommand(commandStr string, args []string) *exec.Cmd {
	fullCmd := commandStr
	if len(args) > 0 {
		escapedArgs := make([]string, len(args))
		for i, a := range args {
			escapedArgs[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(a, "'", "''"))
		}
		fullCmd = fmt.Sprintf("%s %s", commandStr, strings.Join(escapedArgs, " "))
	}
	return exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", fullCmd)
}
