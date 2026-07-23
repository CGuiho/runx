package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CGuiho/runx/embed"
	"github.com/spf13/cobra"
)

const (
	skillID          = "guiho-s-runx"
	skillDescription = "Inspect, validate, describe, and safely execute RunX command catalogs."
)

type SkillInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type SkillMetadata struct {
	Version string `json:"version"`
	Package string `json:"package"`
}

type SkillDetails struct {
	ID          string        `json:"id"`
	Path        string        `json:"path"`
	Description string        `json:"description"`
	Metadata    SkillMetadata `json:"metadata"`
}

var (
	skillLocalFlag  bool
	skillCwdFlag    string
	skillFormatFlag string
	skillFilterFlag string
)

var agentSkillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage the bundled RunX skill",
}

var agentSkillInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the bundled skill into both global tool locations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSkillMutation(cmd, "installed")
	},
}

var agentSkillUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the bundled skill from both tool locations",
	RunE: func(cmd *cobra.Command, args []string) error {
		scope := "global"
		if skillLocalFlag {
			scope = "local"
		}
		dirs, err := getSkillDirectories(scope, skillCwdFlag)
		if err != nil {
			return fmt.Errorf("Agent skill mutation failed: %w", err)
		}
		var removed []string
		for _, dir := range dirs {
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("Agent skill mutation failed: %w", err)
			}
			removed = append(removed, dir)
		}
		return writeFormatted(cmd, skillFormatFlag, map[string][]string{"removed": removed})
	},
}

var agentSkillUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Refresh the bundled skill in both tool locations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSkillMutation(cmd, "updated")
	},
}

var agentSkillListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bundled RunX skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		skills := []SkillInfo{
			{ID: skillID, Description: skillDescription},
		}
		filter := strings.TrimSpace(strings.ToLower(skillFilterFlag))
		if filter != "" {
			var filtered []SkillInfo
			for _, s := range skills {
				text := strings.ToLower(s.ID + " " + s.Description)
				if strings.Contains(text, filter) {
					filtered = append(filtered, s)
				}
			}
			skills = filtered
		}
		return writeFormatted(cmd, skillFormatFlag, skills)
	},
}

var agentSkillShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show metadata for one bundled skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if id != skillID {
			return fmt.Errorf("Unknown RunX skill: %s", id)
		}
		details := SkillDetails{
			ID:          skillID,
			Path:        "skills/" + skillID + "/SKILL.md",
			Description: skillDescription,
			Metadata: SkillMetadata{
				Version: "0.1.0",
				Package: "@guiho/runx",
			},
		}
		return writeFormatted(cmd, skillFormatFlag, details)
	},
}

func init() {
	agentSkillInstallCmd.Flags().BoolVar(&skillLocalFlag, "local", false, "Use project-local tool directories.")
	agentSkillInstallCmd.Flags().StringVar(&skillCwdFlag, "cwd", "", "Use this effective working directory.")
	agentSkillInstallCmd.Flags().StringVar(&skillFormatFlag, "format", "text", "Select output format.")

	agentSkillUninstallCmd.Flags().BoolVar(&skillLocalFlag, "local", false, "Use project-local tool directories.")
	agentSkillUninstallCmd.Flags().StringVar(&skillCwdFlag, "cwd", "", "Use this effective working directory.")
	agentSkillUninstallCmd.Flags().StringVar(&skillFormatFlag, "format", "text", "Select output format.")

	agentSkillUpdateCmd.Flags().BoolVar(&skillLocalFlag, "local", false, "Use project-local tool directories.")
	agentSkillUpdateCmd.Flags().StringVar(&skillCwdFlag, "cwd", "", "Use this effective working directory.")
	agentSkillUpdateCmd.Flags().StringVar(&skillFormatFlag, "format", "text", "Select output format.")

	agentSkillListCmd.Flags().StringVar(&skillFilterFlag, "filter", "", "Filter skill metadata.")
	agentSkillListCmd.Flags().StringVar(&skillFormatFlag, "format", "text", "Select output format.")

	agentSkillShowCmd.Flags().StringVar(&skillFormatFlag, "format", "text", "Select output format.")

	agentSkillCmd.AddCommand(agentSkillInstallCmd)
	agentSkillCmd.AddCommand(agentSkillUninstallCmd)
	agentSkillCmd.AddCommand(agentSkillUpdateCmd)
	agentSkillCmd.AddCommand(agentSkillListCmd)
	agentSkillCmd.AddCommand(agentSkillShowCmd)

	agentCmd.AddCommand(agentSkillCmd)
}

func runSkillMutation(cmd *cobra.Command, label string) error {
	scope := "global"
	if skillLocalFlag {
		scope = "local"
	}
	skillContent, err := getBundledSkill()
	if err != nil {
		return fmt.Errorf("Agent skill mutation failed: %w", err)
	}

	dirs, err := getSkillDirectories(scope, skillCwdFlag)
	if err != nil {
		return fmt.Errorf("Agent skill mutation failed: %w", err)
	}

	var paths []string
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("Agent skill mutation failed: %w", err)
		}
		targetPath := filepath.Join(dir, "SKILL.md")
		if err := os.WriteFile(targetPath, []byte(skillContent), 0644); err != nil {
			return fmt.Errorf("Agent skill mutation failed: %w", err)
		}
		paths = append(paths, targetPath)
	}

	return writeFormatted(cmd, skillFormatFlag, map[string][]string{label: paths})
}

func getSkillDirectories(scope string, cwd string) ([]string, error) {
	var root string
	if scope == "global" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		root = home
	} else {
		root = cwd
		if root == "" {
			var err error
			root, err = os.Getwd()
			if err != nil {
				return nil, err
			}
		}
		var err error
		root, err = filepath.Abs(root)
		if err != nil {
			return nil, err
		}
	}

	return []string{
		filepath.Join(root, ".agents", "skills", skillID),
		filepath.Join(root, ".claude", "skills", skillID),
	}, nil
}

func getBundledSkill() (string, error) {
	data, err := embed.FS.ReadFile("skills/guiho-s-runx.SKILL.md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
