package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new YAML RunX configuration.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("runx.yaml"); err == nil {
			return fmt.Errorf("runx.yaml already exists in current directory")
		}

		content := `version: "2"
namespace: "my-app"
description: "My project command catalog"

commands:
  - uid: "app:hello"
    id: "hello"
    description: "Print hello message"
    run: "echo Hello from RunX Go!"
`
		if err := os.WriteFile("runx.yaml", []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create runx.yaml: %w", err)
		}

		fmt.Println("Initialized runx.yaml successfully.")
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&cwdPath, "cwd", "", "Use this effective working directory.")
	initCmd.Flags().StringVar(&cfgFile, "config", "", "Use this runx.yaml configuration file.")
	initCmd.Flags().String("format", "text", "Select output format.")
	initCmd.Flags().Bool("verbose", false, "Enable diagnostics.")

	RootCmd.AddCommand(initCmd)
}
