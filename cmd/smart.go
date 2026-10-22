package cmd

import "github.com/spf13/cobra"

// newSmartCommand builds the smart entry point: bare smart lists the catalog
// through the canonical list body, and smart with a selector executes it
// through the canonical run body. Only full-name flags are defined, and flag
// parsing stops at the first positional so child arguments after the selector
// are forwarded untouched.
func newSmartCommand(deps Dependencies) *cobra.Command {
	var flags catalogFlags
	var dryRun, yes bool
	command := &cobra.Command{
		Use:   "smart [options] [<uid-or-selector-or-index> [--] [child arguments...]]",
		Short: "List the catalog, or run a command by selector.",
		Long:  "With no selector, list commands in the RunX configuration exactly like runx list. With a selector, execute the selected catalog command exactly like runx run.",
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 0 {
				return runCatalogList(command, deps, flags)
			}
			return runCatalogRun(command, deps, flags, dryRun, yes, args)
		},
	}
	addCatalogFlags(command, &flags)
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Print the execution plan without spawning.")
	command.Flags().BoolVar(&yes, "yes", false, "Approve a confirmation-gated command.")
	command.Flags().SetInterspersed(false)
	return command
}
