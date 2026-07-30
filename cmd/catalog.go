package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/CGuiho/runx/pkg/executor"
	"github.com/CGuiho/runx/pkg/manifest"
	"github.com/spf13/cobra"
)

type catalogFlags struct {
	cwd, config, format string
	verbose             bool
}

func addCatalogFlags(command *cobra.Command, flags *catalogFlags) {
	command.Flags().StringVar(&flags.cwd, "cwd", "", "Use this effective working directory.")
	command.Flags().StringVar(&flags.config, "config", "", "Use this runx.yaml configuration file.")
	command.Flags().StringVar(&flags.format, "format", "text", "Select output format: text or json.")
	command.Flags().BoolVar(&flags.verbose, "verbose", false, "Enable diagnostics.")
}

func validateFormat(value string) error {
	if value != "text" && value != "json" {
		return withExitCode(2, fmt.Errorf("--format must be text or json"))
	}
	return nil
}

func loadCatalog(command *cobra.Command, deps Dependencies, flags catalogFlags) (*manifest.Catalog, error) {
	if err := validateFormat(flags.format); err != nil {
		return nil, err
	}
	home, _ := deps.HomeDir()
	if flags.cwd == "" {
		flags.cwd, _ = deps.Getwd()
	}
	ctx, cancel := context.WithTimeout(command.Context(), 15*time.Second)
	defer cancel()
	catalog, err := manifest.Load(ctx, manifest.LoadOptions{CWD: flags.cwd, ConfigPath: flags.config, HomeDir: home, HTTPClient: deps.HTTPClient})
	if err != nil {
		return nil, withExitCode(3, err)
	}
	if flags.verbose && flags.format == "text" {
		fmt.Fprintf(command.ErrOrStderr(), "configuration file loaded: %s\n", catalog.Path)
	}
	return catalog, nil
}

func writeJSON(command *cobra.Command, value any) error {
	encoder := json.NewEncoder(command.OutOrStdout())
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func newCheckCommand(deps Dependencies) *cobra.Command {
	var flags catalogFlags
	command := &cobra.Command{Use: "check", Short: "Validate a RunX configuration without execution.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		catalog, err := loadCatalog(command, deps, flags)
		if err != nil {
			return err
		}
		result := struct {
			Valid         bool   `json:"valid"`
			Configuration string `json:"configuration"`
			Namespace     string `json:"namespace"`
			Commands      int    `json:"commands"`
			Groups        int    `json:"groups"`
			Children      int    `json:"children"`
		}{true, catalog.Path, catalog.Namespace, len(catalog.Commands), len(catalog.Groups), len(catalog.Children)}
		if flags.format == "json" {
			return writeJSON(command, result)
		}
		fmt.Fprintf(command.OutOrStdout(), "configuration file loaded: %s\nValidation passed: manifest v2 is valid (%d commands, %d groups, %d child catalogs).\n", catalog.Path, len(catalog.Commands), len(catalog.Groups), len(catalog.Children))
		return nil
	}}
	addCatalogFlags(command, &flags)
	return command
}

func newListCommand(deps Dependencies) *cobra.Command {
	var flags catalogFlags
	command := &cobra.Command{Use: "list", Short: "List commands in a RunX configuration.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		catalog, err := loadCatalog(command, deps, flags)
		if err != nil {
			return err
		}
		commands := manifest.SortedCommands(catalog.Commands)
		if flags.format == "json" {
			return writeJSON(command, struct {
				Configuration string                     `json:"configuration"`
				Namespace     string                     `json:"namespace"`
				Commands      []manifest.ResolvedCommand `json:"commands"`
			}{catalog.Path, catalog.Namespace, commands})
		}
		fmt.Fprintf(command.OutOrStdout(), "configuration file loaded: %s\nRunX commands\n\n", catalog.Path)
		table := tabwriter.NewWriter(command.OutOrStdout(), 0, 4, 2, ' ', 0)
		fmt.Fprintln(table, "IDX\tUID\tSELECTOR\tSUMMARY")
		for _, item := range commands {
			suffix := ""
			if item.Confirm == "always" {
				suffix = " [confirm]"
			}
			fmt.Fprintf(table, "%d\t%s\t%s\t%s%s\n", item.Index, item.UID, item.Selector, item.Summary, suffix)
		}
		return table.Flush()
	}}
	addCatalogFlags(command, &flags)
	return command
}

func newDescribeCommand(deps Dependencies) *cobra.Command {
	var flags catalogFlags
	command := &cobra.Command{Use: "describe <uid-or-selector-or-index>", Short: "Describe one catalog command without execution.", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		catalog, err := loadCatalog(command, deps, flags)
		if err != nil {
			return err
		}
		selected, ok := catalog.Resolve(args[0])
		if !ok {
			return withExitCode(3, fmt.Errorf("command %q was not found", args[0]))
		}
		if flags.format == "json" {
			return writeJSON(command, selected)
		}
		fmt.Fprintf(command.OutOrStdout(), "configuration file loaded: %s\n%s (%s)\n\n%s\n\nindex: %d\ncwd: %s\nshell: %s\nconfirmation: %s\ntags: %s\ncommand: %s\n", catalog.Path, selected.UID, selected.Selector, selected.Description, selected.Index, selected.CWD, selected.Shell, selected.Confirm, joinedTags(selected.Tags), selected.Command)
		return nil
	}}
	addCatalogFlags(command, &flags)
	return command
}

func newRunCommand(deps Dependencies) *cobra.Command {
	var flags catalogFlags
	var dryRun, yes bool
	command := &cobra.Command{Use: "run [options] <uid-or-selector-or-index> [--] [child arguments...]", Short: "Execute one selected catalog command.", Args: cobra.MinimumNArgs(1), RunE: func(command *cobra.Command, args []string) error {
		catalog, err := loadCatalog(command, deps, flags)
		if err != nil {
			return err
		}
		selected, ok := catalog.Resolve(args[0])
		if !ok {
			return withExitCode(3, fmt.Errorf("command %q was not found", args[0]))
		}
		childArgs := append([]string{}, args[1:]...)
		if len(childArgs) > 0 && childArgs[0] == "--" {
			childArgs = childArgs[1:]
		}
		if selected.Confirm == "always" && !yes {
			return withExitCode(2, fmt.Errorf("command %s requires --yes before the selector", selected.UID))
		}
		plan := struct {
			UID       string   `json:"uid"`
			Selector  string   `json:"selector"`
			Manifest  string   `json:"manifest"`
			CWD       string   `json:"cwd"`
			Shell     string   `json:"shell"`
			Command   string   `json:"command"`
			Arguments []string `json:"arguments"`
		}{selected.UID, selected.Selector, selected.CatalogPath, selected.CWD, selected.Shell, selected.Command, childArgs}
		if dryRun {
			if flags.format == "json" {
				return writeJSON(command, plan)
			}
			fmt.Fprintf(command.OutOrStdout(), "uid: %s\nselector: %s\nmanifest: %s\ncwd: %s\nshell: %s\ncommand: %s\narguments:\n", plan.UID, plan.Selector, plan.Manifest, plan.CWD, plan.Shell, plan.Command)
			if len(childArgs) == 0 {
				fmt.Fprintln(command.OutOrStdout(), "  (none)")
			} else {
				for index, value := range childArgs {
					fmt.Fprintf(command.OutOrStdout(), "  [%d] %q\n", index, value)
				}
			}
			return nil
		}
		result, err := executor.ExecuteCommand(executor.ExecutionOptions{Command: selected.Command, Args: childArgs, CWD: selected.CWD, Shell: selected.Shell, In: command.InOrStdin(), Out: command.OutOrStdout(), Err: command.ErrOrStderr()})
		if err != nil {
			return withExitCode(1, err)
		}
		if result.ExitCode != 0 {
			return &ChildExitError{Code: result.ExitCode}
		}
		return nil
	}}
	addCatalogFlags(command, &flags)
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Print the execution plan without spawning.")
	command.Flags().BoolVar(&yes, "yes", false, "Approve a confirmation-gated command.")
	command.Flags().SetInterspersed(false)
	return command
}

func newInitCommand(deps Dependencies) *cobra.Command {
	var flags catalogFlags
	command := &cobra.Command{Use: "init", Short: "Create a new YAML RunX configuration.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(flags.format); err != nil {
			return err
		}
		cwd, err := effectiveCWD(deps, flags.cwd)
		if err != nil {
			return withExitCode(5, err)
		}
		path := flags.config
		if path == "" {
			path = "runx.yaml"
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		path, _ = filepath.Abs(path)
		if _, err := os.Stat(path); err == nil {
			return withExitCode(5, fmt.Errorf("configuration already exists: %s", path))
		} else if !os.IsNotExist(err) {
			return withExitCode(5, err)
		}
		namespace := normalizeNamespace(filepath.Base(filepath.Dir(path)))
		content := fmt.Sprintf("version: \"2.0.0\"\n\nnamespace: %q\n\nscripts:\n  directory: \"scripts\"\n\ncommands: []\n", namespace)
		if _, err := manifest.ParseManifestBytes([]byte(content)); err != nil {
			return withExitCode(3, err)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return withExitCode(5, err)
		}
		temporary, err := os.CreateTemp(filepath.Dir(path), ".runx-*.yaml")
		if err != nil {
			return withExitCode(5, err)
		}
		temporaryPath := temporary.Name()
		defer os.Remove(temporaryPath)
		if _, err := temporary.WriteString(content); err != nil {
			temporary.Close()
			return withExitCode(5, err)
		}
		if err := temporary.Close(); err != nil {
			return withExitCode(5, err)
		}
		if err := os.Rename(temporaryPath, path); err != nil {
			return withExitCode(5, err)
		}
		result := struct {
			Status    string `json:"status"`
			Path      string `json:"path"`
			Namespace string `json:"namespace"`
		}{"created", path, namespace}
		if flags.format == "json" {
			return writeJSON(command, result)
		}
		fmt.Fprintf(command.OutOrStdout(), "Created RunX configuration: %s\n", path)
		return nil
	}}
	addCatalogFlags(command, &flags)
	return command
}

func joinedTags(tags []string) string {
	if len(tags) == 0 {
		return "none"
	}
	sorted := append([]string{}, tags...)
	sort.Strings(sorted)
	return strings.Join(sorted, ", ")
}
func normalizeNamespace(value string) string {
	value = strings.ToLower(value)
	value = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if value == "" {
		return "runx"
	}
	if value[0] < 'a' || value[0] > 'z' {
		return "n-" + value
	}
	return value
}
