// Command rx is the short alias launcher installed at $HOME/.guiho/bin/rx(.exe).
// It translates ergonomics (bare rx -> runx list, rx <selector> -> runx run <selector>,
// version and help passthrough) and delegates to the committed active RunX payload.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/CGuiho/runx/pkg/launcher"
)

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

func translateArgs(rawArgs []string) []string {
	if len(rawArgs) == 0 {
		return []string{"list"}
	}
	if isRootPassthrough(rawArgs) {
		return rawArgs
	}

	hasPositional := false
	for i := 0; i < len(rawArgs); {
		arg := rawArgs[i]
		if arg == "--" {
			if i+1 < len(rawArgs) {
				hasPositional = true
			}
			break
		}
		if strings.HasPrefix(arg, "-") {
			if (arg == "--cwd" || arg == "--config" || arg == "--format" || arg == "--help-tree-depth") && i+1 < len(rawArgs) {
				i += 2
			} else {
				i++
			}
		} else {
			hasPositional = true
			break
		}
	}

	if !hasPositional {
		return append([]string{"list"}, rawArgs...)
	}
	return append([]string{"run"}, rawArgs...)
}

func main() {
	args := translateArgs(os.Args[1:])
	pointer, err := installstate.ReadPointer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "rx launcher: %v\n", err)
		os.Exit(5)
	}
	payload, err := launcher.PayloadPath(pointer)
	if err != nil {
		fallback, fallbackErr := launcher.FallbackPath(pointer)
		if fallbackErr != nil || fallback == "" {
			fmt.Fprintf(os.Stderr, "rx launcher: %v\n", err)
			os.Exit(5)
		}
		payload = fallback
	}
	code, err := delegate(payload, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rx launcher: %v\n", err)
		os.Exit(5)
	}
	os.Exit(code)
}
