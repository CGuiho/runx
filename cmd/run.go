package cmd

import (
	"fmt"

	"github.com/CGuiho/runx/pkg/executor"
	"github.com/CGuiho/runx/pkg/manifest"
	"github.com/spf13/cobra"
)

var (
	dryRun  bool
	yesFlag bool
	cwdPath string
)

var runCmd = &cobra.Command{
	Use:   "run <selector> [--] [child arguments...]",
	Short: "Execute one selected catalog command.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		childArgs := args[1:]

		configPath := cfgFile
		if configPath == "" {
			configPath = "runx.yaml"
		}

		m, err := manifest.ParseManifestFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to load manifest: %w", err)
		}

		index, err := manifest.IndexManifest(m, configPath)
		if err != nil {
			return fmt.Errorf("failed to index manifest: %w", err)
		}

		targetCmd, ok := index[selector]
		if !ok {
			return fmt.Errorf("command selector or UID '%s' not found in manifest", selector)
		}

		opts := executor.ExecutionOptions{
			Command: targetCmd.Run,
			Args:    childArgs,
			CWD:     cwdPath,
			DryRun:  dryRun,
		}

		res, err := executor.ExecuteCommand(opts)
		if err != nil {
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf("command exited with code %d", res.ExitCode)
		}

		return nil
	},
}

func init() {
	runCmd.Flags().StringVar(&cwdPath, "cwd", "", "Use this effective working directory.")
	runCmd.Flags().StringVar(&cfgFile, "config", "", "Use this runx.yaml configuration file.")
	runCmd.Flags().String("format", "text", "Select output format.")
	runCmd.Flags().Bool("verbose", false, "Enable diagnostics.")
	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print the execution plan without spawning.")
	runCmd.Flags().BoolVar(&yesFlag, "yes", false, "Approve a confirmation-gated command.")

	RootCmd.AddCommand(runCmd)
}
