package cmd

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

type treeItem struct {
	id, name, description string
	command               *cobra.Command
	flag                  bool
}

// RenderCommandTree renders the hierarchy from command downward. With
// includeGlobalFlags, inherited global flags repeat under every descendant;
// without it (the convention default), global flags appear once on the root
// and are hidden from descendants.
func RenderCommandTree(command *cobra.Command, maxDepth int, includeGlobalFlags bool) string {
	var output strings.Builder
	output.WriteString("COMMAND TREE\n\n")
	fmt.Fprintf(&output, "%s  %s\n", command.CommandPath(), command.Short)
	renderTreeChildren(&output, command, "", 1, maxDepth, command.Parent() == nil || includeGlobalFlags)
	return output.String()
}

func renderTreeChildren(output *strings.Builder, command *cobra.Command, prefix string, depth, maxDepth int, showInherited bool) {
	if maxDepth > 0 && depth > maxDepth {
		return
	}
	items := commandTreeItems(command, showInherited)
	for index, item := range items {
		last := index == len(items)-1
		branch, next := "├── ", prefix+"│   "
		if last {
			branch, next = "└── ", prefix+"    "
		}
		fmt.Fprintf(output, "%s%s%s", prefix, branch, item.name)
		if item.description != "" {
			fmt.Fprintf(output, "  %s", item.description)
		}
		output.WriteByte('\n')
		if item.command != nil {
			renderTreeChildren(output, item.command, next, depth+1, maxDepth, showInherited)
		}
	}
}

func commandTreeItems(command *cobra.Command, showInherited bool) []treeItem {
	items := []treeItem{}
	command.InitDefaultHelpFlag()
	if help := command.Flags().Lookup("help"); help != nil {
		help.Usage = "Show command help."
	}
	for _, child := range command.Commands() {
		if !child.Hidden {
			items = append(items, treeItem{id: child.Name(), name: child.Use, description: child.Short, command: child})
		}
	}
	order := commandOrder(command.Name())
	sort.SliceStable(items, func(i, j int) bool { return order[items[i].id] < order[items[j].id] })
	known := map[string]bool{}
	add := func(flags *pflag.FlagSet) {
		flags.VisitAll(func(flag *pflag.Flag) {
			if flag.Hidden || known[flag.Name] {
				return
			}
			known[flag.Name] = true
			name := "--" + flag.Name
			if flag.Shorthand != "" {
				name = "-" + flag.Shorthand + ", " + name
			}
			if flag.NoOptDefVal == "" {
				name += " <" + flagValueLabel(flag.Name) + ">"
			}
			items = append(items, treeItem{id: flag.Name, name: name, description: flag.Usage, flag: true})
		})
	}
	add(command.LocalNonPersistentFlags())
	if showInherited {
		add(command.InheritedFlags())
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].flag != items[j].flag {
			return !items[i].flag
		}
		if items[i].flag {
			return items[i].name < items[j].name
		}
		return false
	})
	return items
}

func configureDeveloperContext(root *cobra.Command, helpRequested func(*cobra.Command) bool) {
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		if !command.Hidden && strings.TrimSpace(command.Example) == "" {
			path := command.CommandPath()
			command.Example = strings.Join([]string{path + " --help", path + " --help-tree", path + " --help-docs"}, "\n")
		}
		validateArgs := command.Args
		command.Args = func(current *cobra.Command, args []string) error {
			if helpRequested(current) {
				return nil
			}
			if validateArgs == nil {
				return nil
			}
			return validateArgs(current, args)
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
}

func commandOrder(parent string) map[string]int {
	orders := map[string][]string{
		"runx":  {"list", "describe", "run", "check", "init", "agent", "upgrade", "uninstall"},
		"agent": {"skill", "instruction", "prompt"}, "skill": {"install", "uninstall", "update", "list", "show"},
		"instruction": {"apply", "remove", "update", "show"}, "prompt": {"list", "show"}, "upgrade": {"check", "list"},
	}
	result := map[string]int{}
	for index, name := range orders[parent] {
		result[name] = index
	}
	return result
}

func flagValueLabel(name string) string {
	switch name {
	case "cwd", "config":
		return "path"
	case "format":
		return "text|json"
	case "help-tree-depth", "page", "size":
		return "positive-integer"
	case "version":
		return "version"
	case "filter":
		return "keyword"
	default:
		return "value"
	}
}

func RenderHelpDocs(command *cobra.Command) (string, error) {
	var output bytes.Buffer
	command.DisableAutoGenTag = true
	if err := doc.GenMarkdown(command, &output); err != nil {
		return "", fmt.Errorf("generate Markdown help: %w", err)
	}
	return output.String(), nil
}
