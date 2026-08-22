package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/CGuiho/runx/pkg/config"
	"github.com/CGuiho/runx/pkg/maintenance"
	"github.com/CGuiho/runx/pkg/manifest"
	"github.com/spf13/cobra"
)

type initReport struct {
	Created   []string `json:"created"`
	Upgraded  []string `json:"upgraded"`
	Verified  []string `json:"verified"`
	Unchanged []string `json:"unchanged"`
}

func (r *initReport) mark(group, path string) {
	switch group {
	case "created":
		r.Created = append(r.Created, path)
	case "upgraded":
		r.Upgraded = append(r.Upgraded, path)
	case "verified":
		r.Verified = append(r.Verified, path)
	default:
		r.Unchanged = append(r.Unchanged, path)
	}
}

// newInitCommand reconciles RunX for the current project following the
// Convention 0001 initialization sequence: global skills, managed instruction
// block, global configuration, project catalog manifest, and the complete
// agent-evolution policy. It is idempotent and never discards user values.
func newInitCommand(deps Dependencies, info BuildInfo) *cobra.Command {
	var flags catalogFlags
	command := &cobra.Command{Use: "init", Short: "Initialize or reconcile RunX for the current project.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(flags.format); err != nil {
			return err
		}
		cwd, err := effectiveCWD(deps, flags.cwd)
		if err != nil {
			return withExitCode(5, err)
		}
		home, err := deps.HomeDir()
		if err != nil {
			return withExitCode(5, err)
		}
		report := &initReport{Created: []string{}, Upgraded: []string{}, Verified: []string{}, Unchanged: []string{}}

		// 1-2. Verify bundled agent skills in both global destinations.
		skill, err := bundledSkill()
		if err == nil {
			if _, err := maintenance.MaintainAgentIntegration(cwd, home, skill); err != nil {
				return withExitCode(5, fmt.Errorf("reconcile agent resources: %w", err))
			}
			for _, dir := range maintenance.SkillDirectories(home) {
				report.mark("verified", filepath.Join(dir, "SKILL.md"))
			}
		}

		// 3. Global configuration baseline.
		globalPath := filepath.Join(home, ".guiho", "runx", "runx.global.yaml")
		globalDoc, globalErr := config.LoadGlobal(globalPath)
		switch {
		case globalErr == nil:
			report.mark("verified", globalPath)
		case errors.Is(globalErr, fs.ErrNotExist):
			globalDoc = config.GlobalConfig{}
			if err := writeGlobalConfig(globalPath, info.Version, globalDoc); err != nil {
				return withExitCode(5, err)
			}
			report.mark("created", globalPath)
		default:
			return withExitCode(3, globalErr)
		}

		// 4. Project catalog manifest.
		manifestPath := flags.config
		if manifestPath == "" {
			manifestPath = filepath.Join(cwd, "runx.yaml")
		} else if !filepath.IsAbs(manifestPath) {
			manifestPath = filepath.Join(cwd, manifestPath)
		}
		manifestPath, _ = filepath.Abs(manifestPath)
		projectManifest, manifestStatus, err := ensureProjectManifest(manifestPath, cwd)
		if err != nil {
			return withExitCode(3, err)
		}
		report.mark(manifestStatus, manifestPath)

		// 5-6. Complete the agent-evolution policy.
		policyChanged, err := reconcileEvolutionPolicy(deps, command, report, globalPath, &globalDoc, projectManifest)
		if err != nil {
			return withExitCode(3, err)
		}
		if policyChanged {
			if err := writeGlobalConfig(globalPath, info.Version, globalDoc); err != nil {
				return withExitCode(5, err)
			}
			report.mark("upgraded", globalPath)
		}
		resolved, err := config.Resolve(projectConfigFrom(projectManifest), globalDoc)
		if err != nil {
			return withExitCode(3, err)
		}

		if flags.format == "json" {
			return writeJSON(command, map[string]any{
				"status":  "initialized",
				"policy":  resolved,
				"report":  report,
				"project": manifestPath,
				"global":  globalPath,
			})
		}
		out := command.OutOrStdout()
		printGroup(out, "Created", report.Created)
		printGroup(out, "Upgraded", report.Upgraded)
		printGroup(out, "Verified", report.Verified)
		printGroup(out, "Unchanged", report.Unchanged)
		fmt.Fprintf(out, "agent.evolution.upgrade=%s issues.bugs=%s issues.improvements=%s issues.reviews=%s\n",
			resolved.Upgrade, resolved.Issues.Bugs, resolved.Issues.Improvements, resolved.Issues.Reviews)
		fmt.Fprintln(out, "[OK] RunX initialized.")
		return nil
	}}
	addCatalogFlags(command, &flags)
	return command
}

func printGroup(out io.Writer, name string, paths []string) {
	fmt.Fprintf(out, "%s:\n", name)
	if len(paths) == 0 {
		fmt.Fprintln(out, "  (nothing)")
		return
	}
	for _, path := range paths {
		fmt.Fprintf(out, "  %s\n", path)
	}
}

// ensureProjectManifest loads the existing catalog manifest or creates a
// minimal valid one. It returns the parsed manifest plus "created" or
// "verified".
func ensureProjectManifest(path, cwd string) (*manifest.Manifest, string, error) {
	if data, err := os.ReadFile(path); err == nil {
		parsed, err := manifest.ParseManifestBytes(data)
		if err != nil {
			return nil, "", fmt.Errorf("existing manifest %s: %w", path, err)
		}
		return parsed, "verified", nil
	} else if !os.IsNotExist(err) {
		return nil, "", err
	}
	namespace := normalizeNamespace(filepath.Base(filepath.Dir(path)))
	content := fmt.Sprintf("version: \"2.0.0\"\n\nnamespace: %q\n\nscripts:\n  directory: \".scripts\"\n\ncommands: []\n", namespace)
	if _, err := manifest.ParseManifestBytes([]byte(content)); err != nil {
		return nil, "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, "", err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".runx-*.yaml")
	if err != nil {
		return nil, "", err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if _, err := temporary.WriteString(content); err != nil {
		temporary.Close()
		return nil, "", err
	}
	if err := temporary.Close(); err != nil {
		return nil, "", err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return nil, "", err
	}
	parsed, err := manifest.ParseManifestFile(path)
	if err != nil {
		return nil, "", err
	}
	return parsed, "created", nil
}

func projectPolicy(m *manifest.Manifest) *config.Agent {
	if m == nil {
		return nil
	}
	return m.Agent
}

// reconcileEvolutionPolicy completes missing policy fields. Interactive
// terminals are asked; noninteractive invocations without a complete policy
// fail closed instead of guessing.
func reconcileEvolutionPolicy(deps Dependencies, command *cobra.Command, report *initReport, globalPath string, globalDoc *config.GlobalConfig, projectManifest *manifest.Manifest) (bool, error) {
	missing := policyMissingRaw(projectPolicy(projectManifest), globalDoc)
	if len(missing) == 0 {
		report.mark("verified", "agent.evolution policy")
		return false, nil
	}
	interactive := deps.IsTerminal != nil && deps.IsTerminal(deps.In)
	if !interactive {
		return false, fmt.Errorf("agent.evolution is not fully configured (%s); run runx init in an interactive terminal or set the values in %s", strings.Join(missing, ", "), globalPath)
	}
	out := command.ErrOrStderr()
	reader := bufio.NewReader(deps.In)
	fmt.Fprintln(out, "RunX is governed by an agent-evolution policy controlling upgrades and GitHub issue creation.")
	fmt.Fprintln(out, "Each field accepts exactly one of:")
	fmt.Fprintln(out, "  disabled        do not perform the governed action on your behalf")
	fmt.Fprintln(out, "  always-ask      explain, ask, and wait for approval each time")
	fmt.Fprintln(out, "  always-proceed  perform the action, then report the result")
	fmt.Fprintln(out, "Recommendation: always-proceed so the agent keeps RunX upgraded and contributes bugs, improvements, and reviews.")
	fmt.Fprintln(out, "Answers are stored in the GLOBAL configuration: "+globalPath)

	ask := func(prompt string, current string) (string, error) {
		suffix := ""
		if current != "" {
			suffix = " [" + current + "]"
		}
		fmt.Fprintf(out, "%s%s: ", prompt, suffix)
		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			return "", fmt.Errorf("read answer: %w", err)
		}
		answer := strings.ToLower(strings.TrimSpace(line))
		if answer == "" && current != "" {
			return current, nil
		}
		return normalizePolicyValue(answer)
	}

	applyAll := ""
	fmt.Fprint(out, "Apply always-proceed to all agent-evolution actions? [y/N]: ")
	line, _ := reader.ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		applyAll = config.PolicyAlwaysProceed
	}

	fields := []struct {
		label   string
		current string
		set     func(string)
	}{
		{"agent.evolution.upgrade", rawValue(projectPolicy(projectManifest), globalDoc, "upgrade"), func(v string) { globalDoc.Agent.Evolution.Upgrade = v }},
		{"agent.evolution.issues.bugs", rawValue(projectPolicy(projectManifest), globalDoc, "bugs"), func(v string) { globalDoc.Agent.Evolution.Issues.Bugs = v }},
		{"agent.evolution.issues.improvements", rawValue(projectPolicy(projectManifest), globalDoc, "improvements"), func(v string) { globalDoc.Agent.Evolution.Issues.Improvements = v }},
		{"agent.evolution.issues.reviews", rawValue(projectPolicy(projectManifest), globalDoc, "reviews"), func(v string) { globalDoc.Agent.Evolution.Issues.Reviews = v }},
	}
	for _, field := range fields {
		value := applyAll
		if value == "" {
			var askErr error
			value, askErr = ask(field.label, field.current)
			if askErr != nil {
				return false, askErr
			}
		}
		field.set(value)
	}
	// Validate before writing.
	if _, err := config.Resolve(projectConfigFrom(projectManifest), *globalDoc); err != nil {
		return false, err
	}
	return true, nil
}

func rawValue(project *config.Agent, global *config.GlobalConfig, field string) string {
	if project != nil {
		switch field {
		case "upgrade":
			if project.Evolution.Upgrade != "" {
				return project.Evolution.Upgrade
			}
		case "bugs":
			if project.Evolution.Issues.Bugs != "" {
				return project.Evolution.Issues.Bugs
			}
		case "improvements":
			if project.Evolution.Issues.Improvements != "" {
				return project.Evolution.Issues.Improvements
			}
		case "reviews":
			if project.Evolution.Issues.Reviews != "" {
				return project.Evolution.Issues.Reviews
			}
		}
	}
	switch field {
	case "upgrade":
		return global.Agent.Evolution.Upgrade
	case "bugs":
		return global.Agent.Evolution.Issues.Bugs
	case "improvements":
		return global.Agent.Evolution.Issues.Improvements
	case "reviews":
		return global.Agent.Evolution.Issues.Reviews
	}
	return ""
}

func normalizePolicyValue(answer string) (string, error) {
	switch answer {
	case "":
		return config.DefaultPolicy, nil
	case config.PolicyDisabled, config.PolicyAlwaysAsk, config.PolicyAlwaysProceed:
		return answer, nil
	default:
		return "", fmt.Errorf("invalid value %q: must be one of disabled, always-ask, always-proceed", answer)
	}
}

func policyMissingRaw(project *config.Agent, global *config.GlobalConfig) []string {
	missing := []string{}
	for _, field := range []string{"upgrade", "bugs", "improvements", "reviews"} {
		if rawValue(project, global, field) == "" {
			missing = append(missing, field)
		}
	}
	return missing
}

func writeGlobalConfig(path, version string, document config.GlobalConfig) error {
	if version == "" || version == "dev" {
		version = "0.0.0-dev"
	}
	return config.WriteGlobal(path, version, document)
}

func projectConfigFrom(m *manifest.Manifest) config.ProjectConfig {
	if m == nil || m.Agent == nil {
		return config.ProjectConfig{}
	}
	return config.ProjectConfig{Agent: *m.Agent}
}
