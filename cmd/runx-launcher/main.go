// Command runx-launcher is the stable launcher installed at
// $HOME/.guiho/bin/runx(.exe). It resolves the committed active payload and
// delegates with identical arguments, environment, working directory, and
// standard streams, then exits with the payload's exit code. On start failure
// it falls back to the retained previous verified payload.
package main

import (
	"fmt"
	"os"

	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/CGuiho/runx/pkg/launcher"
)

func main() {
	pointer, err := installstate.ReadPointer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "runx launcher: %v\n", err)
		os.Exit(5)
	}
	payload, err := launcher.PayloadPath(pointer)
	if err != nil {
		fallback, fallbackErr := launcher.FallbackPath(pointer)
		if fallbackErr != nil || fallback == "" {
			fmt.Fprintf(os.Stderr, "runx launcher: %v\n", err)
			os.Exit(5)
		}
		payload = fallback
	}
	code, err := delegate(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runx launcher: %v\n", err)
		os.Exit(5)
	}
	os.Exit(code)
}
