package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
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
	addCatalogLocationFlags(command, flags)
	command.Flags().StringVar(&flags.format, "format", "text", "Select output format: text or json.")
}

func addCatalogLocationFlags(command *cobra.Command, flags *catalogFlags) {
	command.Flags().StringVar(&flags.cwd, "cwd", "", "Use this effective working directory.")
	command.Flags().StringVar(&flags.config, "config", "", "Use this runx.yaml configuration file.")
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

func newRevealCommand(deps Dependencies) *cobra.Command {
	var flags catalogFlags
	flags.format = "text"
	command := &cobra.Command{Use: "reveal <uid-or-selector-or-index>", Short: "Print one catalog command without execution.", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		catalog, err := loadCatalog(command, deps, flags)
		if err != nil {
			return err
		}
		selected, ok := catalog.Resolve(args[0])
		if !ok {
			return withExitCode(3, fmt.Errorf("command %q was not found", args[0]))
		}
		_, err = fmt.Fprint(command.OutOrStdout(), selected.Command, "\n")
		return err
	}}
	addCatalogLocationFlags(command, &flags)
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
			retry := confirmationRetry(command, flags, dryRun, args[0], childArgs)
			if flags.format == "text" && deps.IsTerminal(command.InOrStdin()) {
				approved, promptErr := promptForConfirmation(command, selected.UID, retry)
				if promptErr != nil {
					return withExitCode(2, fmt.Errorf("could not read confirmation: %w", promptErr))
				}
				if !approved {
					return withExitCode(2, fmt.Errorf("command %s was not authorized; to run without prompting: %s", selected.UID, retry))
				}
			} else {
				return withExitCode(2, fmt.Errorf("command %s requires confirmation; rerun exactly: %s", selected.UID, retry))
			}
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

func promptForConfirmation(command *cobra.Command, uid, retry string) (bool, error) {
	fmt.Fprintf(command.ErrOrStderr(), "Command %s requires confirmation.\nTo skip this prompt, run:\n  %s\nAre you sure? [y/N]: ", uid, retry)
	answer, err := readConfirmationAnswer(command.InOrStdin())
	if err != nil {
		return false, err
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

func readConfirmationAnswer(reader io.Reader) (string, error) {
	var answer strings.Builder
	buffer := []byte{0}
	for {
		count, err := reader.Read(buffer)
		if count > 0 {
			switch buffer[0] {
			case '\n':
				return answer.String(), nil
			case '\r':
			default:
				if answer.Len() < 64 {
					answer.WriteByte(buffer[0])
				}
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return answer.String(), nil
			}
			return "", err
		}
		if count == 0 {
			return "", io.ErrNoProgress
		}
	}
}

func confirmationRetry(command *cobra.Command, flags catalogFlags, dryRun bool, selector string, childArgs []string) string {
	parts := []string{"runx", "run"}
	for _, item := range []struct {
		name  string
		value string
	}{
		{"cwd", flags.cwd},
		{"config", flags.config},
		{"format", flags.format},
	} {
		if command.Flags().Changed(item.name) {
			parts = append(parts, "--"+item.name, quoteRetryArgument(item.value))
		}
	}
	if command.Flags().Changed("verbose") {
		if flags.verbose {
			parts = append(parts, "--verbose")
		} else {
			parts = append(parts, "--verbose=false")
		}
	}
	if command.Flags().Changed("dry-run") {
		if dryRun {
			parts = append(parts, "--dry-run")
		} else {
			parts = append(parts, "--dry-run=false")
		}
	}
	parts = append(parts, "--yes", quoteRetryArgument(selector))
	if len(childArgs) > 0 {
		parts = append(parts, "--")
		for _, argument := range childArgs {
			parts = append(parts, quoteRetryArgument(argument))
		}
	}
	return strings.Join(parts, " ")
}

func quoteRetryArgument(value string) string {
	if value != "" && strings.IndexFunc(value, func(character rune) bool {
		return !((character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') ||
			strings.ContainsRune("_-./\\:@%+=,", character))
	}) == -1 {
		return value
	}
	return strconv.Quote(value)
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
