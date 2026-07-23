package cmd

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	helpTree       bool
	helpTreeDepth  int
	helpDocs       bool
	exitOnHelpTree = true
)

type FlagNode struct {
	Flag        string
	Description string
}

type TreeNode struct {
	Name        string
	Description string
	Flags       []FlagNode
	Children    []*TreeNode
}

// BuildCanonicalTree returns the exact tree node structure matching Section 4 of the RFC.
func BuildCanonicalTree() *TreeNode {
	return &TreeNode{
		Name:        "runx",
		Description: "",
		Flags: []FlagNode{
			{"--version", "Show the RunX version."},
			{"--help", "Show command help."},
			{"--help-tree", "Show this command hierarchy."},
			{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
			{"--help-docs", "Emit Markdown documentation for this command."},
		},
		Children: []*TreeNode{
			{
				Name:        "list",
				Description: "List commands in a RunX configuration.",
				Flags: []FlagNode{
					{"--cwd <path>", "Use this effective working directory."},
					{"--config <path>", "Use this runx.yaml configuration file."},
					{"--format <text|json>", "Select output format."},
					{"--verbose", "Enable diagnostics."},
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
				},
			},
			{
				Name:        "describe",
				Description: "Describe one catalog command without execution.",
				Flags: []FlagNode{
					{"--cwd <path>", "Use this effective working directory."},
					{"--config <path>", "Use this runx.yaml configuration file."},
					{"--format <text|json>", "Select output format."},
					{"--verbose", "Enable diagnostics."},
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
				},
			},
			{
				Name:        "run",
				Description: "Execute one selected catalog command.",
				Flags: []FlagNode{
					{"--cwd <path>", "Use this effective working directory."},
					{"--config <path>", "Use this runx.yaml configuration file."},
					{"--format <text|json>", "Select output format."},
					{"--verbose", "Enable diagnostics."},
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
					{"--dry-run", "Print the execution plan without spawning."},
					{"--yes", "Approve a confirmation-gated command."},
				},
			},
			{
				Name:        "check",
				Description: "Validate a RunX configuration without execution.",
				Flags: []FlagNode{
					{"--cwd <path>", "Use this effective working directory."},
					{"--config <path>", "Use this runx.yaml configuration file."},
					{"--format <text|json>", "Select output format."},
					{"--verbose", "Enable diagnostics."},
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
				},
			},
			{
				Name:        "init",
				Description: "Create a new YAML RunX configuration.",
				Flags: []FlagNode{
					{"--cwd <path>", "Use this effective working directory."},
					{"--config <path>", "Use this runx.yaml configuration file."},
					{"--format <text|json>", "Select output format."},
					{"--verbose", "Enable diagnostics."},
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
				},
			},
			{
				Name:        "agent",
				Description: "Manage RunX agent integration.",
				Flags: []FlagNode{
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
				},
				Children: []*TreeNode{
					{
						Name:        "skill",
						Description: "Manage the bundled RunX skill.",
						Flags: []FlagNode{
							{"--help", "Show command help."},
							{"--help-tree", "Show this command hierarchy."},
							{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
							{"--help-docs", "Emit Markdown documentation for this command."},
						},
						Children: []*TreeNode{
							{
								Name:        "install",
								Description: "Install the bundled skill into both global tool locations.",
								Flags: []FlagNode{
									{"--cwd <path>", "Use this effective working directory."},
									{"--local", "Use project-local tool directories."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "uninstall",
								Description: "Remove the bundled skill from both tool locations.",
								Flags: []FlagNode{
									{"--cwd <path>", "Use this effective working directory."},
									{"--local", "Use project-local tool directories."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "update",
								Description: "Refresh the bundled skill in both tool locations.",
								Flags: []FlagNode{
									{"--cwd <path>", "Use this effective working directory."},
									{"--local", "Use project-local tool directories."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "list",
								Description: "List bundled RunX skills.",
								Flags: []FlagNode{
									{"--filter <keyword>", "Filter skill metadata."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "show",
								Description: "Show metadata for one bundled skill.",
								Flags: []FlagNode{
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
						},
					},
					{
						Name:        "instruction",
						Description: "Manage RunX instruction blocks.",
						Flags: []FlagNode{
							{"--help", "Show command help."},
							{"--help-tree", "Show this command hierarchy."},
							{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
							{"--help-docs", "Emit Markdown documentation for this command."},
						},
						Children: []*TreeNode{
							{
								Name:        "apply",
								Description: "Apply the managed instruction block.",
								Flags: []FlagNode{
									{"--cwd <path>", "Use this effective working directory."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "remove",
								Description: "Remove the managed instruction block.",
								Flags: []FlagNode{
									{"--cwd <path>", "Use this effective working directory."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "update",
								Description: "Replace stale managed instruction content.",
								Flags: []FlagNode{
									{"--cwd <path>", "Use this effective working directory."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "show",
								Description: "Print the raw instruction template.",
								Flags: []FlagNode{
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
						},
					},
					{
						Name:        "prompt",
						Description: "Inspect bundled agent prompts.",
						Flags: []FlagNode{
							{"--help", "Show command help."},
							{"--help-tree", "Show this command hierarchy."},
							{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
							{"--help-docs", "Emit Markdown documentation for this command."},
						},
						Children: []*TreeNode{
							{
								Name:        "list",
								Description: "List bundled RunX prompts.",
								Flags: []FlagNode{
									{"--names", "Print prompt names only."},
									{"--format <text|json>", "Select output format."},
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
							{
								Name:        "show",
								Description: "Print one raw bundled prompt.",
								Flags: []FlagNode{
									{"--help", "Show command help."},
									{"--help-tree", "Show this command hierarchy."},
									{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
									{"--help-docs", "Emit Markdown documentation for this command."},
								},
							},
						},
					},
				},
			},
			{
				Name:        "upgrade",
				Description: "Inspect or upgrade a native RunX executable.",
				Flags: []FlagNode{
					{"--version <version>", "Select an exact release version."},
					{"--arch <x64|arm64>", "Select target architecture."},
					{"--variant <baseline|default|modern>", "Select x64 binary variant."},
					{"--dry-run", "Plan without mutation."},
					{"--format <text|json>", "Select output format."},
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
				},
				Children: []*TreeNode{
					{
						Name:        "check",
						Description: "Check whether a newer stable release exists.",
						Flags: []FlagNode{
							{"--format <text|json>", "Select output format."},
							{"--help", "Show command help."},
							{"--help-tree", "Show this command hierarchy."},
							{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
							{"--help-docs", "Emit Markdown documentation for this command."},
						},
					},
					{
						Name:        "list",
						Description: "List RunX releases newest first.",
						Flags: []FlagNode{
							{"--page <positive-integer>", "Select result page."},
							{"--per-page <positive-integer>", "Select page size."},
							{"--pre-releases", "Accepted explicitly; prereleases are always included."},
							{"--arch <x64|arm64>", "Select target architecture."},
							{"--variant <baseline|default|modern>", "Select x64 variant."},
							{"--format <text|json>", "Select output format."},
							{"--help", "Show command help."},
							{"--help-tree", "Show this command hierarchy."},
							{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
							{"--help-docs", "Emit Markdown documentation for this command."},
						},
					},
				},
			},
			{
				Name:        "uninstall",
				Description: "Uninstall the native RunX executable.",
				Flags: []FlagNode{
					{"--dry-run", "Print the target without deleting it."},
					{"--format <text|json>", "Select output format."},
					{"--help", "Show command help."},
					{"--help-tree", "Show this command hierarchy."},
					{"--help-tree-depth <positive-integer>", "Limit help-tree recursion depth."},
					{"--help-docs", "Emit Markdown documentation for this command."},
				},
			},
		},
	}
}

// RenderHelpTree writes the formatted command tree to w.
func RenderHelpTree(w io.Writer, maxDepth int) {
	root := BuildCanonicalTree()
	fmt.Fprintln(w, "COMMAND TREE")
	fmt.Fprintln(w)
	fmt.Fprintln(w, root.Name)
	renderChildren(w, root, "", 1, maxDepth)
}

func renderChildren(w io.Writer, node *TreeNode, prefix string, currentDepth, maxDepth int) {
	totalChildren := len(node.Children) + len(node.Flags)

	for i := 0; i < totalChildren; i++ {
		isLast := (i == totalChildren-1)
		itemPrefix := "├── "
		if isLast {
			itemPrefix = "└── "
		}

		if i < len(node.Children) {
			child := node.Children[i]
			left := prefix + itemPrefix + child.Name
			fmt.Fprintln(w, formatHelpLine(left, child.Description, currentDepth))

			if maxDepth == 0 || currentDepth < maxDepth {
				nextPrefix := prefix + "│   "
				if isLast {
					nextPrefix = prefix + "    "
				}
				renderChildren(w, child, nextPrefix, currentDepth+1, maxDepth)
			}
		} else {
			flagIdx := i - len(node.Children)
			flag := node.Flags[flagIdx]
			left := prefix + itemPrefix + flag.Flag
			fmt.Fprintln(w, formatHelpLine(left, flag.Description, currentDepth))
		}
	}
}

func formatHelpLine(left, desc string, depth int) string {
	targetCol := 42 + 4*(depth-1)
	width := utf8.RuneCountInString(left)
	if width < targetCol {
		padding := strings.Repeat(" ", targetCol-width)
		return left + padding + desc
	}
	return left + "  " + desc
}

// RenderHelpDocs writes Markdown documentation for a Cobra command to w.
func RenderHelpDocs(w io.Writer, cmd *cobra.Command) {
	fmt.Fprintf(w, "# %s\n\n", cmd.CommandPath())
	if cmd.Short != "" {
		fmt.Fprintf(w, "%s\n\n", cmd.Short)
	}
	fmt.Fprintf(w, "## Usage\n\n```bash\n%s\n```\n\n", cmd.UseLine())
	fmt.Fprintln(w, "## Flags")
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(w, "- `--%s`: %s\n", f.Name, f.Usage)
	})
}
