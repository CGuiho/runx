package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const (
	managedStart         = "<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->"
	managedEnd           = "<!-- END RUNX -->"
	mojibakeManagedStart = "<!-- BEGIN RUNX \u00e2\u20ac\u201d DO NOT EDIT THIS SECTION -->"
	legacyManagedStart   = "<!-- BEGIN RUNX AGENT INSTRUCTIONS -->"
	legacyManagedEnd     = "<!-- END RUNX AGENT INSTRUCTIONS -->"
)

var (
	instructionCwdFlag    string
	instructionFormatFlag string
)

var agentInstructionCmd = &cobra.Command{
	Use:   "instruction",
	Short: "Manage RunX instruction blocks",
}

var agentInstructionApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the managed instruction block",
	RunE: func(cmd *cobra.Command, args []string) error {
		return reconcileInstructions(cmd, "updated")
	},
}

var agentInstructionRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the managed instruction block",
	RunE: func(cmd *cobra.Command, args []string) error {
		targets, err := getInstructionTargets(instructionCwdFlag)
		if err != nil {
			return fmt.Errorf("Agent instruction mutation failed: %w", err)
		}
		var changed []string
		for _, targetPath := range targets {
			if !pathExists(targetPath) {
				continue
			}
			data, err := os.ReadFile(targetPath)
			if err != nil {
				return fmt.Errorf("Agent instruction mutation failed: %w", err)
			}
			existing := string(data)
			next := removeManagedBlock(existing)
			if next != existing {
				if err := os.WriteFile(targetPath, []byte(next), 0644); err != nil {
					return fmt.Errorf("Agent instruction mutation failed: %w", err)
				}
				changed = append(changed, targetPath)
			}
		}
		return writeFormatted(cmd, instructionFormatFlag, map[string][]string{"removed": changed})
	},
}

var agentInstructionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Replace stale managed instruction content",
	RunE: func(cmd *cobra.Command, args []string) error {
		return reconcileInstructions(cmd, "updated")
	},
}

var agentInstructionShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the raw instruction template",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprint(cmd.OutOrStdout(), getInstructionBlock())
		return nil
	},
}

func init() {
	agentInstructionApplyCmd.Flags().StringVar(&instructionCwdFlag, "cwd", "", "Use this effective working directory.")
	agentInstructionApplyCmd.Flags().StringVar(&instructionFormatFlag, "format", "text", "Select output format.")

	agentInstructionRemoveCmd.Flags().StringVar(&instructionCwdFlag, "cwd", "", "Use this effective working directory.")
	agentInstructionRemoveCmd.Flags().StringVar(&instructionFormatFlag, "format", "text", "Select output format.")

	agentInstructionUpdateCmd.Flags().StringVar(&instructionCwdFlag, "cwd", "", "Use this effective working directory.")
	agentInstructionUpdateCmd.Flags().StringVar(&instructionFormatFlag, "format", "text", "Select output format.")

	agentInstructionCmd.AddCommand(agentInstructionApplyCmd)
	agentInstructionCmd.AddCommand(agentInstructionRemoveCmd)
	agentInstructionCmd.AddCommand(agentInstructionUpdateCmd)
	agentInstructionCmd.AddCommand(agentInstructionShowCmd)

	agentCmd.AddCommand(agentInstructionCmd)
}

func reconcileInstructions(cmd *cobra.Command, label string) error {
	targets, err := getInstructionTargets(instructionCwdFlag)
	if err != nil {
		return fmt.Errorf("Agent instruction mutation failed: %w", err)
	}
	var changed []string
	block := getInstructionBlock()

	for _, targetPath := range targets {
		var existing string
		if pathExists(targetPath) {
			data, err := os.ReadFile(targetPath)
			if err != nil {
				return fmt.Errorf("Agent instruction mutation failed: %w", err)
			}
			existing = string(data)
		}
		next := replaceManagedBlock(existing, block)
		if next != existing {
			if err := os.WriteFile(targetPath, []byte(next), 0644); err != nil {
				return fmt.Errorf("Agent instruction mutation failed: %w", err)
			}
			changed = append(changed, targetPath)
		}
	}

	return writeFormatted(cmd, instructionFormatFlag, map[string][]string{label: changed})
}

func getInstructionBlock() string {
	return managedStart + "\n" +
		"## RunX Command Catalog\n\n" +
		"Load the `guiho-s-runx` skill whenever discovering commands, creating or\n" +
		"updating catalog entries, validating `runx.yaml`, inspecting command details,\n" +
		"or executing RunX commands.\n" +
		"Start with `runx check --format json` and `runx list --format json`, select\n" +
		"stable UIDs, use `runx describe <uid>`, and run\n" +
		"`runx run --dry-run <uid>` before unfamiliar or side-effecting work.\n" +
		"RunX options precede the selector; post-selector tokens belong to the child.\n" +
		managedEnd + "\n"
}

func removeKnownManagedBlocks(existing string, includeLeadingWhitespace bool) string {
	pairs := [][2]string{
		{managedStart, managedEnd},
		{mojibakeManagedStart, managedEnd},
		{legacyManagedStart, legacyManagedEnd},
	}
	output := existing
	for _, p := range pairs {
		start, end := p[0], p[1]
		prefix := ""
		if includeLeadingWhitespace {
			prefix = `\s*`
		}
		pattern := prefix + regexp.QuoteMeta(start) + `[\s\S]*?` + regexp.QuoteMeta(end) + `\s*`
		re := regexp.MustCompile(pattern)
		repl := ""
		if includeLeadingWhitespace {
			repl = "\n"
		}
		output = re.ReplaceAllString(output, repl)
	}
	return output
}

func replaceManagedBlock(existing, block string) string {
	stripped := strings.TrimRight(removeKnownManagedBlocks(existing, false), "\r\n\t ")
	if stripped == "" {
		return block
	}
	return stripped + "\n\n" + block
}

func removeManagedBlock(existing string) string {
	next := strings.TrimRight(removeKnownManagedBlocks(existing, true), "\r\n\t ")
	if next == "" {
		return ""
	}
	return next + "\n"
}

func getInstructionTargets(cwd string) ([]string, error) {
	root := cwd
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	agents := filepath.Join(root, "AGENTS.md")
	claude := filepath.Join(root, "CLAUDE.md")

	agentsExists := pathExists(agents)
	claudeExists := pathExists(claude)

	if agentsExists && claudeExists {
		return []string{agents, claude}, nil
	}
	if claudeExists {
		return []string{claude}, nil
	}
	return []string{agents}, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
