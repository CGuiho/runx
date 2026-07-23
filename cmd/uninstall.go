package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	uninstallDryRun bool
	uninstallFormat string
)

type UninstallResult struct {
	Target      string `json:"target"`
	DryRun      bool   `json:"dry_run"`
	Uninstalled bool   `json:"uninstalled"`
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the native RunX executable.",
	Long:  "Removes the installed native RunX binary from the system PATH.",
	RunE: func(cmd *cobra.Command, args []string) error {
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to locate executable: %w", err)
		}

		if uninstallDryRun {
			if uninstallFormat == "json" || jsonFmt {
				res := UninstallResult{
					Target:      execPath,
					DryRun:      true,
					Uninstalled: false,
				}
				data, _ := json.MarshalIndent(res, "", "  ")
				fmt.Println(string(data))
				return nil
			}
			fmt.Printf("Dry-run: Would remove native RunX executable at %s\n", execPath)
			return nil
		}

		if err := os.Remove(execPath); err != nil {
			return fmt.Errorf("failed to remove executable at %s: %w", execPath, err)
		}

		if uninstallFormat == "json" || jsonFmt {
			res := UninstallResult{
				Target:      execPath,
				DryRun:      false,
				Uninstalled: true,
			}
			data, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Successfully uninstalled native RunX executable at %s\n", execPath)
		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallDryRun, "dry-run", false, "Print the target without deleting it.")
	uninstallCmd.Flags().StringVar(&uninstallFormat, "format", "text", "Select output format.")

	RootCmd.AddCommand(uninstallCmd)
}
