package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/CGuiho/runx/cmd"
)

var (
	version     = "dev"
	commit      = "unknown"
	buildDate   = "unknown"
	buildTarget = "development"
)

func main() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}
	err := cmd.Execute(cmd.BuildInfo{Version: version, Commit: commit, Date: buildDate, Target: buildTarget})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cmd.ExitCode(err))
	}
}
