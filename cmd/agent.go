package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CGuiho/runx/embed"
	"github.com/CGuiho/runx/pkg/maintenance"
	"github.com/spf13/cobra"
)

const skillID = "guiho-s-runx"

func newAgentCommand(deps Dependencies) *cobra.Command {
	command := &cobra.Command{Use: "agent", Short: "Manage RunX agent integration.", Args: cobra.NoArgs, RunE: helpCommand}
	command.AddCommand(newAgentSkillCommand(deps), newAgentInstructionCommand(deps), newAgentPromptCommand())
	return command
}

func helpCommand(command *cobra.Command, _ []string) error { return command.Help() }

func newAgentSkillCommand(deps Dependencies) *cobra.Command {
	root := &cobra.Command{Use: "skill", Short: "Manage the bundled RunX skill.", Args: cobra.NoArgs, RunE: helpCommand}
	root.AddCommand(newSkillMutationCommand(deps, "install"), newSkillMutationCommand(deps, "update"), newSkillUninstallCommand(deps))
	var listFormat, filter string
	list := &cobra.Command{Use: "list", Short: "List bundled RunX skills.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(listFormat); err != nil {
			return err
		}
		value := struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		}{skillID, "Inspect, validate, describe, and safely execute RunX command catalogs."}
		if filter != "" && !strings.Contains(strings.ToLower(value.ID+" "+value.Description), strings.ToLower(filter)) {
			if listFormat == "json" {
				return writeJSON(command, []any{})
			}
			return nil
		}
		if listFormat == "json" {
			return writeJSON(command, []any{value})
		}
		fmt.Fprintf(command.OutOrStdout(), "%s  %s\n", value.ID, value.Description)
		return nil
	}}
	list.Flags().StringVar(&filter, "filter", "", "Filter skill metadata.")
	list.Flags().StringVar(&listFormat, "format", "text", "Select output format: text or json.")
	var showFormat string
	show := &cobra.Command{Use: "show <id>", Short: "Show metadata for one bundled skill.", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		if err := validateFormat(showFormat); err != nil {
			return err
		}
		if args[0] != skillID {
			return withExitCode(2, fmt.Errorf("unknown RunX skill %q", args[0]))
		}
		value := struct {
			ID          string `json:"id"`
			Path        string `json:"path"`
			Description string `json:"description"`
		}{skillID, "skills/guiho-s-runx/SKILL.md", "Inspect, validate, describe, and safely execute RunX command catalogs."}
		if showFormat == "json" {
			return writeJSON(command, value)
		}
		fmt.Fprintf(command.OutOrStdout(), "%s\npath: %s\n%s\n", value.ID, value.Path, value.Description)
		return nil
	}}
	show.Flags().StringVar(&showFormat, "format", "text", "Select output format: text or json.")
	root.AddCommand(list, show)
	return root
}

func skillDirectories(deps Dependencies, local bool, cwd string) ([]string, error) {
	var root string
	var err error
	if local {
		root, err = effectiveCWD(deps, cwd)
	} else {
		root, err = deps.HomeDir()
	}
	if err != nil {
		return nil, err
	}
	return []string{filepath.Join(root, ".agents", "skills", skillID), filepath.Join(root, ".claude", "skills", skillID)}, nil
}

func newSkillMutationCommand(deps Dependencies, action string) *cobra.Command {
	var local bool
	var cwd, format string
	command := &cobra.Command{Use: action, Short: map[string]string{"install": "Install the bundled skill into both tool locations.", "update": "Refresh the bundled skill in both tool locations."}[action], Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(format); err != nil {
			return err
		}
		content, err := bundledSkill()
		if err != nil {
			return withExitCode(5, err)
		}
		directories, err := skillDirectories(deps, local, cwd)
		if err != nil {
			return withExitCode(5, err)
		}
		paths := []string{}
		for _, directory := range directories {
			target := filepath.Join(directory, "SKILL.md")
			if err := maintenance.WriteTextFileAtomic(target, content); err != nil {
				return withExitCode(5, err)
			}
			paths = append(paths, target)
		}
		result := map[string]any{"action": action, "paths": paths}
		if format == "json" {
			return writeJSON(command, result)
		}
		for _, path := range paths {
			fmt.Fprintf(command.OutOrStdout(), "%s: %s\n", action, path)
		}
		return nil
	}}
	command.Flags().BoolVar(&local, "local", false, "Use project-local tool directories.")
	command.Flags().StringVar(&cwd, "cwd", "", "Use this effective working directory.")
	command.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
	return command
}

func newSkillUninstallCommand(deps Dependencies) *cobra.Command {
	var local bool
	var cwd, format string
	command := &cobra.Command{Use: "uninstall", Short: "Remove the bundled skill from both tool locations.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(format); err != nil {
			return err
		}
		directories, err := skillDirectories(deps, local, cwd)
		if err != nil {
			return withExitCode(5, err)
		}
		removed := []string{}
		for _, directory := range directories {
			if err := os.RemoveAll(directory); err != nil {
				return withExitCode(5, err)
			}
			removed = append(removed, directory)
		}
		if format == "json" {
			return writeJSON(command, map[string]any{"removed": removed})
		}
		for _, path := range removed {
			fmt.Fprintf(command.OutOrStdout(), "removed: %s\n", path)
		}
		return nil
	}}
	command.Flags().BoolVar(&local, "local", false, "Use project-local tool directories.")
	command.Flags().StringVar(&cwd, "cwd", "", "Use this effective working directory.")
	command.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
	return command
}

func instructionTargets(deps Dependencies, cwd string) ([]string, error) {
	root, err := effectiveCWD(deps, cwd)
	if err != nil {
		return nil, err
	}
	agents, claude := filepath.Join(root, "AGENTS.md"), filepath.Join(root, "CLAUDE.md")
	aExists, cExists := maintenance.PathExists(agents), maintenance.PathExists(claude)
	if aExists && cExists {
		return []string{agents, claude}, nil
	}
	if cExists {
		return []string{claude}, nil
	}
	return []string{agents}, nil
}

func newAgentInstructionCommand(deps Dependencies) *cobra.Command {
	root := &cobra.Command{Use: "instruction", Short: "Manage RunX instruction blocks.", Args: cobra.NoArgs, RunE: helpCommand}
	for _, action := range []string{"apply", "update", "remove"} {
		action := action
		var cwd, format string
		command := &cobra.Command{Use: action, Short: action + " the managed RunX instruction block.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
			if err := validateFormat(format); err != nil {
				return err
			}
			targets, err := instructionTargets(deps, cwd)
			if err != nil {
				return withExitCode(5, err)
			}
			changed := []string{}
			for _, path := range targets {
				existing, readErr := maintenance.ReadTextIfExists(path)
				if readErr != nil {
					return withExitCode(5, readErr)
				}
				next, updateErr := maintenance.ReplaceManagedBlockStrict(existing, maintenance.DefaultInstructionBlock())
				if action == "remove" {
					next, updateErr = maintenance.RemoveManagedBlockStrict(existing)
				}
				if updateErr != nil {
					return withExitCode(5, fmt.Errorf("update managed RunX block in %s: %w", path, updateErr))
				}
				if next != existing {
					if err := maintenance.WriteTextFileAtomic(path, next); err != nil {
						return withExitCode(5, err)
					}
					changed = append(changed, path)
				}
			}
			if format == "json" {
				return writeJSON(command, map[string]any{"action": action, "paths": changed})
			}
			for _, path := range changed {
				fmt.Fprintf(command.OutOrStdout(), "%s: %s\n", action, path)
			}
			return nil
		}}
		command.Flags().StringVar(&cwd, "cwd", "", "Use this effective working directory.")
		command.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
		root.AddCommand(command)
	}
	root.AddCommand(&cobra.Command{Use: "show", Short: "Print the raw instruction template.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		fmt.Fprint(command.OutOrStdout(), maintenance.DefaultInstructionBlock())
		return nil
	}})
	return root
}

func newAgentPromptCommand() *cobra.Command {
	root := &cobra.Command{Use: "prompt", Short: "Inspect bundled agent prompts.", Args: cobra.NoArgs, RunE: helpCommand}
	var names bool
	var format string
	list := &cobra.Command{Use: "list", Short: "List bundled RunX prompts.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(format); err != nil {
			return err
		}
		entries, err := embed.FS.ReadDir("prompts")
		if err != nil {
			return err
		}
		ids := []string{}
		details := []map[string]string{}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || strings.HasSuffix(entry.Name(), ".xdocs.md") {
				continue
			}
			data, err := embed.FS.ReadFile("prompts/" + entry.Name())
			if err != nil {
				continue
			}
			id := strings.TrimSuffix(entry.Name(), ".md")
			if fid := promptIDFromFrontmatter(string(data)); fid != "" {
				id = fid
			}
			ids = append(ids, id)
			details = append(details, map[string]string{"id": id, "path": "prompts/" + entry.Name()})
		}
		if names {
			if format == "json" {
				return writeJSON(command, ids)
			}
			for _, id := range ids {
				fmt.Fprintln(command.OutOrStdout(), id)
			}
			return nil
		}
		if format == "json" {
			return writeJSON(command, details)
		}
		for _, d := range details {
			fmt.Fprintf(command.OutOrStdout(), "%s  %s\n", d["id"], d["path"])
		}
		return nil
	}}
	list.Flags().BoolVar(&names, "names", false, "Print prompt names only.")
	list.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
	show := &cobra.Command{Use: "show <id>", Short: "Print one raw bundled prompt.", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		path := "prompts/" + args[0] + ".md"
		if data, err := embed.FS.ReadFile(path); err == nil {
			fmt.Fprint(command.OutOrStdout(), string(data))
			return nil
		}
		entries, err := embed.FS.ReadDir("prompts")
		if err != nil {
			return withExitCode(2, fmt.Errorf("unknown RunX prompt %q", args[0]))
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || strings.HasSuffix(entry.Name(), ".xdocs.md") {
				continue
			}
			data, err := embed.FS.ReadFile("prompts/" + entry.Name())
			if err != nil {
				continue
			}
			if promptIDFromFrontmatter(string(data)) == args[0] {
				fmt.Fprint(command.OutOrStdout(), string(data))
				return nil
			}
		}
		return withExitCode(2, fmt.Errorf("unknown RunX prompt %q", args[0]))
	}}
	root.AddCommand(list, show)
	return root
}

func promptIDFromFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---") {
		return ""
	}
	rest := content[3:]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return ""
	}
	frontmatter := rest[:end]
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "name:"))
			value = strings.Trim(value, "\"'")
			return value
		}
	}
	return ""
}
