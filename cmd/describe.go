package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/CGuiho/runx/pkg/manifest"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe <uid>",
	Short: "Describe one catalog command without execution.",
	Long:  "Inspect metadata and execution configuration for a single catalog command without running it.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

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

		formatVal, _ := cmd.Flags().GetString("format")
		if formatVal == "json" || jsonFmt {
			data, err := json.MarshalIndent(targetCmd, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to serialize command to JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("UID:          %s\n", targetCmd.UID)
		fmt.Printf("Full ID:      %s\n", targetCmd.FullID)
		fmt.Printf("Description:  %s\n", targetCmd.Description)
		fmt.Printf("Run:          %s\n", targetCmd.Run)
		if targetCmd.CWD != "" {
			fmt.Printf("CWD:          %s\n", targetCmd.CWD)
		}
		if targetCmd.Confirm != "" {
			fmt.Printf("Confirm:      %s\n", targetCmd.Confirm)
		}
		fmt.Printf("Catalog Path: %s\n", targetCmd.CatalogPath)

		return nil
	},
}

func init() {
	describeCmd.Flags().StringVar(&cwdPath, "cwd", "", "Use this effective working directory.")
	describeCmd.Flags().StringVar(&cfgFile, "config", "", "Use this runx.yaml configuration file.")
	describeCmd.Flags().String("format", "text", "Select output format.")
	describeCmd.Flags().Bool("verbose", false, "Enable diagnostics.")

	RootCmd.AddCommand(describeCmd)
}
