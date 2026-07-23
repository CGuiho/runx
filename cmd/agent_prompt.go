package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/CGuiho/runx/embed"
	"github.com/spf13/cobra"
)

const (
	promptID          = "guiho-i-runx"
	promptDescription = "Guide an agent through safe RunX catalog inspection and execution."
)

type PromptInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

var (
	promptNamesFlag  bool
	promptFormatFlag string
)

var agentPromptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Inspect bundled agent prompts",
}

var agentPromptListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bundled RunX prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		isJSON := jsonFmt || promptFormatFlag == "json"
		if promptNamesFlag {
			names := []string{promptID}
			if isJSON {
				bytes, err := json.MarshalIndent(names, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), promptID)
			}
			return nil
		}

		prompts := []PromptInfo{
			{ID: promptID, Description: promptDescription},
		}
		return writeFormatted(cmd, promptFormatFlag, prompts)
	},
}

var agentPromptShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Print one raw bundled prompt",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if id != promptID {
			return fmt.Errorf("Unknown RunX prompt: %s", id)
		}
		promptContent, err := getBundledPrompt()
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), promptContent)
		return nil
	},
}

func init() {
	agentPromptListCmd.Flags().BoolVar(&promptNamesFlag, "names", false, "Print prompt names only.")
	agentPromptListCmd.Flags().StringVar(&promptFormatFlag, "format", "text", "Select output format.")

	agentPromptCmd.AddCommand(agentPromptListCmd)
	agentPromptCmd.AddCommand(agentPromptShowCmd)

	agentCmd.AddCommand(agentPromptCmd)
}

func getBundledPrompt() (string, error) {
	data, err := embed.FS.ReadFile("prompts/guiho-i-runx.md")
	if err != nil {
		return "", fmt.Errorf("Bundled guiho-i-runx prompt is unavailable: %w", err)
	}
	return string(data), nil
}
