package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedHelpTree = `COMMAND TREE

runx
├── list                                  List commands in a RunX configuration.
│   ├── --cwd <path>                          Use this effective working directory.
│   ├── --config <path>                       Use this runx.yaml configuration file.
│   ├── --format <text|json>                  Select output format.
│   ├── --verbose                             Enable diagnostics.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   └── --help-docs                           Emit Markdown documentation for this command.
├── describe                              Describe one catalog command without execution.
│   ├── --cwd <path>                          Use this effective working directory.
│   ├── --config <path>                       Use this runx.yaml configuration file.
│   ├── --format <text|json>                  Select output format.
│   ├── --verbose                             Enable diagnostics.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   └── --help-docs                           Emit Markdown documentation for this command.
├── run                                   Execute one selected catalog command.
│   ├── --cwd <path>                          Use this effective working directory.
│   ├── --config <path>                       Use this runx.yaml configuration file.
│   ├── --format <text|json>                  Select output format.
│   ├── --verbose                             Enable diagnostics.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   ├── --help-docs                           Emit Markdown documentation for this command.
│   ├── --dry-run                             Print the execution plan without spawning.
│   └── --yes                                 Approve a confirmation-gated command.
├── check                                 Validate a RunX configuration without execution.
│   ├── --cwd <path>                          Use this effective working directory.
│   ├── --config <path>                       Use this runx.yaml configuration file.
│   ├── --format <text|json>                  Select output format.
│   ├── --verbose                             Enable diagnostics.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   └── --help-docs                           Emit Markdown documentation for this command.
├── init                                  Create a new YAML RunX configuration.
│   ├── --cwd <path>                          Use this effective working directory.
│   ├── --config <path>                       Use this runx.yaml configuration file.
│   ├── --format <text|json>                  Select output format.
│   ├── --verbose                             Enable diagnostics.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   └── --help-docs                           Emit Markdown documentation for this command.
├── agent                                 Manage RunX agent integration.
│   ├── skill                                 Manage the bundled RunX skill.
│   │   ├── install                               Install the bundled skill into both global tool locations.
│   │   │   ├── --cwd <path>                          Use this effective working directory.
│   │   │   ├── --local                               Use project-local tool directories.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── uninstall                             Remove the bundled skill from both tool locations.
│   │   │   ├── --cwd <path>                          Use this effective working directory.
│   │   │   ├── --local                               Use project-local tool directories.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── update                                Refresh the bundled skill in both tool locations.
│   │   │   ├── --cwd <path>                          Use this effective working directory.
│   │   │   ├── --local                               Use project-local tool directories.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── list                                  List bundled RunX skills.
│   │   │   ├── --filter <keyword>                    Filter skill metadata.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── show                                  Show metadata for one bundled skill.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── --help                                Show command help.
│   │   ├── --help-tree                           Show this command hierarchy.
│   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   └── --help-docs                           Emit Markdown documentation for this command.
│   ├── instruction                           Manage RunX instruction blocks.
│   │   ├── apply                                 Apply the managed instruction block.
│   │   │   ├── --cwd <path>                          Use this effective working directory.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── remove                                Remove the managed instruction block.
│   │   │   ├── --cwd <path>                          Use this effective working directory.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── update                                Replace stale managed instruction content.
│   │   │   ├── --cwd <path>                          Use this effective working directory.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── show                                  Print the raw instruction template.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── --help                                Show command help.
│   │   ├── --help-tree                           Show this command hierarchy.
│   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   └── --help-docs                           Emit Markdown documentation for this command.
│   ├── prompt                                Inspect bundled agent prompts.
│   │   ├── list                                  List bundled RunX prompts.
│   │   │   ├── --names                               Print prompt names only.
│   │   │   ├── --format <text|json>                  Select output format.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── show                                  Print one raw bundled prompt.
│   │   │   ├── --help                                Show command help.
│   │   │   ├── --help-tree                           Show this command hierarchy.
│   │   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   │   └── --help-docs                           Emit Markdown documentation for this command.
│   │   ├── --help                                Show command help.
│   │   ├── --help-tree                           Show this command hierarchy.
│   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   └── --help-docs                           Emit Markdown documentation for this command.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   └── --help-docs                           Emit Markdown documentation for this command.
├── upgrade                               Inspect or upgrade a native RunX executable.
│   ├── check                                 Check whether a newer stable release exists.
│   │   ├── --format <text|json>                  Select output format.
│   │   ├── --help                                Show command help.
│   │   ├── --help-tree                           Show this command hierarchy.
│   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   └── --help-docs                           Emit Markdown documentation for this command.
│   ├── list                                  List RunX releases newest first.
│   │   ├── --page <positive-integer>             Select result page.
│   │   ├── --per-page <positive-integer>         Select page size.
│   │   ├── --pre-releases                        Accepted explicitly; prereleases are always included.
│   │   ├── --arch <x64|arm64>                    Select target architecture.
│   │   ├── --variant <baseline|default|modern>   Select x64 variant.
│   │   ├── --format <text|json>                  Select output format.
│   │   ├── --help                                Show command help.
│   │   ├── --help-tree                           Show this command hierarchy.
│   │   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   │   └── --help-docs                           Emit Markdown documentation for this command.
│   ├── --version <version>                   Select an exact release version.
│   ├── --arch <x64|arm64>                    Select target architecture.
│   ├── --variant <baseline|default|modern>   Select x64 binary variant.
│   ├── --dry-run                             Plan without mutation.
│   ├── --format <text|json>                  Select output format.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   └── --help-docs                           Emit Markdown documentation for this command.
├── uninstall                             Uninstall the native RunX executable.
│   ├── --dry-run                             Print the target without deleting it.
│   ├── --format <text|json>                  Select output format.
│   ├── --help                                Show command help.
│   ├── --help-tree                           Show this command hierarchy.
│   ├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
│   └── --help-docs                           Emit Markdown documentation for this command.
├── --version                             Show the RunX version.
├── --help                                Show command help.
├── --help-tree                           Show this command hierarchy.
├── --help-tree-depth <positive-integer>  Limit help-tree recursion depth.
└── --help-docs                           Emit Markdown documentation for this command.`

func TestHelpTreeCanonicalMatrix(t *testing.T) {
	var buf bytes.Buffer
	RenderHelpTree(&buf, 0)
	actual := strings.TrimSpace(buf.String())
	expected := strings.TrimSpace(expectedHelpTree)

	actualLines := strings.Split(actual, "\n")
	expectedLines := strings.Split(expected, "\n")

	assert.Equal(t, len(expectedLines), len(actualLines), "Line count should match RFC Section 4 matrix")

	for i := 0; i < len(expectedLines) && i < len(actualLines); i++ {
		assert.Equal(t, expectedLines[i], actualLines[i], "Mismatch at line %d", i+1)
	}
}

func TestHelpTreeDepthLimit(t *testing.T) {
	var buf bytes.Buffer
	RenderHelpTree(&buf, 1)
	out := buf.String()

	assert.Contains(t, out, "├── list")
	assert.Contains(t, out, "├── agent")
	assert.NotContains(t, out, "├── skill")
	assert.NotContains(t, out, "├── install")
}

func TestHelpDocsRendering(t *testing.T) {
	var buf bytes.Buffer
	RenderHelpDocs(&buf, listCmd)
	out := buf.String()

	assert.Contains(t, out, "# runx list")
	assert.Contains(t, out, "List commands in a RunX configuration.")
	assert.Contains(t, out, "## Usage")
}
