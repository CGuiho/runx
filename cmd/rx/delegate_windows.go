//go:build windows

package main

import (
	"os"
)

// delegate starts the payload attached to this console with the translated
// arguments. Windows keeps the child's standard handles and exit code; the shell
// remains blocked until the payload completes because StartProcess waits via Wait().
func delegate(payload string, args ...string) (int, error) {
	argv := append([]string{payload}, args...)
	attr := &os.ProcAttr{
		Dir:   "",
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Env:   os.Environ(),
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
