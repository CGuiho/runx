//go:build !windows

package main

import (
	"os"
	"syscall"
)

// delegate starts payload with the launcher's arguments, environment, working
// directory, and standard streams, waits for it, and returns its exit code.
func delegate(payload string) (int, error) {
	argv := append([]string{payload}, os.Args[1:]...)
	attr := &os.ProcAttr{
		Dir:   "",
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Env:   os.Environ(),
		Sys:   &syscall.SysProcAttr{},
	}
	process, err := os.StartProcess(payload, argv, attr)
	if err != nil {
		return 0, err
	}
	state, err := process.Wait()
	if err != nil {
		return 0, err
	}
	return state.ExitCode(), nil
}
