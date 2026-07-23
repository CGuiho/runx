package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands in a RunX configuration.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available commands:")
		fmt.Println("  app:hello - Print hello message")
		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&cwdPath, "cwd", "", "Use this effective working directory.")
	listCmd.Flags().StringVar(&cfgFile, "config", "", "Use this runx.yaml configuration file.")
	listCmd.Flags().String("format", "text", "Select output format.")
	listCmd.Flags().Bool("verbose", false, "Enable diagnostics.")

	RootCmd.AddCommand(listCmd)
}
