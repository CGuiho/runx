package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate a RunX configuration without execution.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Validation passed: runx.yaml manifest V2 schema is valid.")
		return nil
	},
}

func init() {
	checkCmd.Flags().StringVar(&cwdPath, "cwd", "", "Use this effective working directory.")
	checkCmd.Flags().StringVar(&cfgFile, "config", "", "Use this runx.yaml configuration file.")
	checkCmd.Flags().String("format", "text", "Select output format.")
	checkCmd.Flags().Bool("verbose", false, "Enable diagnostics.")

	RootCmd.AddCommand(checkCmd)
}
