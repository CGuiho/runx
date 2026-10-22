// Command rxs is the short launcher installed at $HOME/.guiho/bin/rxs(.exe).
// It delegates to the committed active RunX payload: bare rxs runs the payload
// smart subcommand (catalog listing), root entry flags pass through unchanged,
// and every other argument list forwards as smart <args> so the payload owns
// the list-or-run translation.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/CGuiho/runx/pkg/launcher"
)

// isRootPassthrough reports whether the argument list belongs to the payload
// root: the version and help entry points that smart must not own.
func isRootPassthrough(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "-v", "--version", "-h", "--help", "--help-tree", "--help-docs", "--help-tree-depth", "--help-tree-global-flags":
		return true
	}
	if strings.HasPrefix(args[0], "--help-tree-depth=") {
		return true
	}
	if len(args) == 1 && args[0] == "--color" {
		return true
	}
	return false
}

// translateArgs maps rxs ergonomics onto the payload: root entry flags reach
// the payload root unchanged, and everything else forwards through the smart
// subcommand so the payload parses its own options before the selector.
func translateArgs(rawArgs []string) []string {
	if len(rawArgs) == 0 {
		return []string{"smart"}
	}
	if isRootPassthrough(rawArgs) {
		return rawArgs
	}
	return append([]string{"smart"}, rawArgs...)
}

func main() {
	args := translateArgs(os.Args[1:])
	pointer, err := installstate.ReadPointer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "rxs launcher: %v\n", err)
		os.Exit(5)
	}
	payload, err := launcher.PayloadPath(pointer)
	if err != nil {
		fallback, fallbackErr := launcher.FallbackPath(pointer)
		if fallbackErr != nil || fallback == "" {
			fmt.Fprintf(os.Stderr, "rxs launcher: %v\n", err)
			os.Exit(5)
		}
		payload = fallback
	}
	code, err := delegate(payload, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rxs launcher: %v\n", err)
		os.Exit(5)
	}
	os.Exit(code)
}
