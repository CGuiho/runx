package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// newConventionUninstallCommand implements the shared Convention 0001
// uninstallation contract: a REMOVE/PRESERVE plan, full-name preservation
// flags, mandatory confirmation in interactive terminals, and fail-closed
// behavior for noninteractive invocations without --yes.
func newConventionUninstallCommand(deps Dependencies) *cobra.Command {
	var preserveConfig, preserveData, dryRun, assumeYes, jsonOutput bool
	command := &cobra.Command{Use: "uninstall", Short: "Uninstall RunX and every artifact it owns.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		cwd, err := deps.Getwd()
		if err != nil {
			return withExitCode(1, err)
		}
		plan, err := buildUninstallPlan(cwd, preserveConfig, preserveData)
		if err != nil {
			return withExitCode(5, err)
		}

		if jsonOutput {
			if err := writeJSON(command, plan); err != nil {
				return err
			}
		} else {
			renderUninstallPlan(command, plan)
		}

		if dryRun {
			if !jsonOutput {
				fmt.Fprintln(command.OutOrStdout(), "Dry run: nothing was changed.")
			}
			return nil
		}

		if !assumeYes {
			if !isatty.IsTerminal(os.Stdin.Fd()) || !isatty.IsTerminal(os.Stdout.Fd()) {
				return withExitCode(2, fmt.Errorf("noninteractive uninstall requires --yes"))
			}
			fmt.Fprint(command.OutOrStdout(), "Proceed with uninstallation? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.ToLower(strings.TrimSpace(answer))
			if answer != "y" && answer != "yes" {
				fmt.Fprintln(command.OutOrStdout(), "Aborted.")
				return nil
			}
		}

		executed, err := executeUninstallPlan(plan)
		if err != nil {
			return withExitCode(5, err)
		}
		if !jsonOutput {
			for _, path := range executed {
				fmt.Fprintf(command.OutOrStdout(), "removed: %s\n", path)
			}
			fmt.Fprintln(command.OutOrStdout(), "[OK] RunX uninstalled.")
		}
		return nil
	}}
	command.Flags().BoolVar(&preserveConfig, "preserve-config", false, "Preserve runx.global.yaml and the current project's runx.yaml.")
	command.Flags().BoolVar(&preserveData, "preserve-data", false, "Preserve persistent data paths and CLI-owned databases.")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Print the removal and preservation plan without changing anything.")
	command.Flags().BoolVar(&assumeYes, "yes", false, "Confirm the displayed plan without an interactive confirmation prompt.")
	command.Flags().BoolVar(&jsonOutput, "format-json", false, "Emit one stable JSON document instead of text.")
	return command
}

type uninstallPlan struct {
	Remove   []string `json:"remove"`
	Preserve []string `json:"preserve"`
}

func buildUninstallPlan(cwd string, preserveConfig, preserveData bool) (*uninstallPlan, error) {
	cliDir, err := installstate.CLIDir()
	if err != nil {
		return nil, err
	}
	launcherPath, err := installstate.LauncherPath()
	if err != nil {
		return nil, err
	}
	guihoRoot, err := installstate.GUISafeHome()
	if err != nil {
		return nil, err
	}
	tempRoot, err := installstate.TempRoot()
	if err != nil {
		return nil, err
	}
	binDir, err := installstate.BinDir()
	if err != nil {
		return nil, err
	}

	plan := &uninstallPlan{Remove: []string{}, Preserve: []string{
		guihoRoot + " (shared GUIHO home)",
		binDir + " (shared launcher directory)",
		tempRoot + " (shared staging root)",
		"user PATH entry for " + binDir,
	}}

	appendIfExists := func(list *[]string, path string) {
		if _, err := os.Stat(path); err == nil {
			*list = append(*list, path)
		}
	}

	appendIfExists(&plan.Remove, launcherPath)
	if _, err := os.Stat(cliDir); err == nil {
		if preserveConfig || preserveData {
			plan.Preserve = append(plan.Preserve,
				filepath.Join(cliDir, "runx.global.yaml"),
				filepath.Join(cliDir, "data"),
				cliDir+" directory itself (preserved children remain)")
			appendIfExists(&plan.Remove, filepath.Join(cliDir, "versions"))
			appendIfExists(&plan.Remove, filepath.Join(cliDir, "resources"))
			appendIfExists(&plan.Remove, filepath.Join(cliDir, "current.json"))
			appendIfExists(&plan.Remove, filepath.Join(cliDir, "installed-artifacts.json"))
		} else {
			plan.Remove = append(plan.Remove, cliDir)
		}
	}
	for _, skillsRoot := range []string{".agents", ".claude"} {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		appendIfExists(&plan.Remove, filepath.Join(home, skillsRoot, "skills", "guiho-s-runx"))
	}
	for _, marker := range instructionMarkerNames() {
		path := filepath.Join(cwd, marker)
		if fileExists(path) {
			plan.Remove = append(plan.Remove, fmt.Sprintf("RunX managed block in %s", path))
		}
	}
	projectManifest := filepath.Join(cwd, "runx.yaml")
	if fileExists(projectManifest) {
		if preserveConfig {
			plan.Preserve = append(plan.Preserve, projectManifest)
		} else {
			plan.Remove = append(plan.Remove, projectManifest)
		}
	}
	return plan, nil
}

func instructionMarkerNames() []string {
	if runtime.GOOS == "windows" {
		return []string{"AGENTS.md", "CLAUDE.md"}
	}
	return []string{"AGENTS.md", "CLAUDE.md"}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func renderUninstallPlan(command *cobra.Command, plan *uninstallPlan) {
	out := command.OutOrStdout()
	fmt.Fprintln(out, "Uninstallation plan")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "REMOVE:")
	if len(plan.Remove) == 0 {
		fmt.Fprintln(out, "  (nothing)")
	} else {
		for _, item := range plan.Remove {
			fmt.Fprintf(out, "  %s\n", item)
		}
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "PRESERVE:")
	for _, item := range plan.Preserve {
		fmt.Fprintf(out, "  %s\n", item)
	}
	fmt.Fprintln(out)
}

// removeOwnedPath fails closed unless the target is inside the CLI-owned home,
// is the launcher, or is a managed skill projection.
func removeOwnedPath(path string) error {
	cliDir, err := installstate.CLIDir()
	if err != nil {
		return err
	}
	launcherPath, err := installstate.LauncherPath()
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cleanPath := filepath.Clean(path)
	cleanCLI := filepath.Clean(cliDir)
	if cleanPath == cleanLauncherPath(launcherPath) ||
		cleanPath == filepath.Clean(filepath.Join(home, ".agents", "skills", "guiho-s-runx")) ||
		cleanPath == filepath.Clean(filepath.Join(home, ".claude", "skills", "guiho-s-runx")) ||
		isWithin(cleanPath, cleanCLI) {
		return os.RemoveAll(path)
	}
	return fmt.Errorf("refusing to remove unowned path: %s", path)
}

func cleanLauncherPath(launcher string) string { return filepath.Clean(launcher) }

func isWithin(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel != ".." && !filepath.IsAbs(rel)
}

func executeUninstallPlan(plan *uninstallPlan) ([]string, error) {
	executed := []string{}
	for _, target := range plan.Remove {
		switch {
		case strings.HasPrefix(target, "RunX managed block in "):
			file := strings.TrimPrefix(target, "RunX managed block in ")
			if err := removeManagedInstructionBlock(file); err != nil {
				return executed, err
			}
			executed = append(executed, target)
		case filepath.Base(target) == "guiho-s-runx":
			if err := removeOwnedPath(target); err != nil && !os.IsNotExist(err) {
				return executed, err
			}
			executed = append(executed, target)
		default:
			if err := removeOwnedPath(target); err != nil && !os.IsNotExist(err) {
				return executed, err
			}
			executed = append(executed, target)
		}
	}
	return executed, nil
}

// removeManagedInstructionBlock removes only the bounded RunX block and never
// deletes AGENTS.md itself.
func removeManagedInstructionBlock(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	kept := make([]string, 0, len(lines))
	inside := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "<!-- BEGIN RUNX"):
			inside = true
			continue
		case strings.HasPrefix(trimmed, "<!-- END RUNX"):
			inside = false
			continue
		}
		if !inside {
			kept = append(kept, line)
		}
	}
	output := strings.Join(kept, "\n")
	output = strings.TrimRight(output, "\n") + "\n"
	return os.WriteFile(path, []byte(output), 0o644)
}
