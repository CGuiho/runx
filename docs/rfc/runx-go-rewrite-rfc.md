---
name: RunX Go Rewrite RFC (Issue #22)
purpose: Master technical specification, architecture, implementation plan, and verification roadmap for porting RunX from Bun/TypeScript to Go using Cobra and Viper with a strict YAML-only policy and TypeScript coexistence rule.
description: Exhaustive RFC document for GitHub Issue 22, detailing the Cobra command tree, Viper YAML-only configuration management, Go standard library tools, Manifest V2 compliance, argument transport, detached update worker, TypeScript coexistence directive, complete implementation plan, cross-compilation release matrix, environment variables, exit codes, and verification criteria.
created: 2026-07-24
flags:
  - accepted
tags:
  - rfc
  - go
  - cobra
  - viper
  - yaml-only
  - architecture
  - cli
  - implementation-plan
  - coexistence
keywords:
  - issue 22
  - go rewrite
  - cobra
  - viper
  - yaml only
  - golang
  - runx
  - rfc 0034
  - help tree
  - implementation plan
  - typescript coexistence
owner: runx-architecture
---

# GUIHO RFC: RunX Go Rewrite Specification & Master Implementation Plan (Issue #22)

## 1. Executive Summary & Goals

This Request for Comments (RFC) serves as the **master technical specification** and **executable implementation plan** for **GitHub Issue [#22](https://github.com/CGuiho/runx/issues/22)**: rewriting the **RunX CLI** (`@guiho/runx`) from Bun/TypeScript to **Go (Golang)** using **Cobra** (`github.com/spf13/cobra`) as the primary CLI framework and **Viper** (`github.com/spf13/viper`) for configuration resolution under a strict **YAML-Only Policy**.

### Strategic Objectives

1. **Sub-10ms Cold Startup**: Eliminate JavaScript runtime bootstrap overhead to achieve instant (<10ms) command resolution across all invocations.
2. **Compact Binary Size**: Reduce single-executable binary sizes from ~50 MiB (Bun compiled runtime) down to ~8–12 MiB per target platform.
3. **Strict YAML-Only Directive**: Completely drop all `.toml` configuration file formats. RunX exclusively parses and supports YAML (`.yaml` / `.yml`) for all manifest files (`runx.yaml`) and configuration surfaces.
4. **TypeScript Preservation & Coexistence Policy**: **Do not delete any existing TypeScript code** (`source/`, `scripts/`, `library/`, `package.json`, `tsconfig.json`). The Go codebase (`main.go`, `cmd/`, `pkg/`, `embed/`, `go.mod`, `go.sum`) coexists directly alongside the TypeScript code in the same repository. The TypeScript codebase remains untouched until a future decommission date explicitly authorized by the developer.
5. **Cobra Command Framework**: Leverage Cobra for command routing (`runx`, `list`, `describe`, `run`, `check`, `init`, `agent`, `upgrade`, `uninstall`), flag handling, subcommands, usage/help rendering, and shell completion generation.
6. **Viper & YAML Schema Configuration**: Use Viper and `gopkg.in/yaml.v3` for strict `runx.yaml` decoding (`KnownFields(true)`), environment variable bindings (`RUNX_*`), and configuration precedence (Flags > Env > YAML Config File > Defaults).
7. **Command Tree Matrix Parity**: The migration is complete when running `runx --help-tree` against the compiled Go executable prints the exact 100% compliant command tree matrix specified in Section 4.

---

## 2. Mandatory Technology Stack & Dependencies

To guarantee high performance, memory safety, zero external C-dependencies (`CGO_ENABLED=0`), and full RFC 0034 CLI parity, the Go rewrite of RunX uses the following curated set of tools and libraries:

### Core Frameworks & Configuration
1. **Cobra (`github.com/spf13/cobra`)**
   - **Role**: Command tree definition (`runx`, `list`, `describe`, `run`, `check`, `init`, `agent`, `upgrade`, `uninstall`), subcommands, POSIX flag routing, usage/help formatting, `--help-tree` recursive hierarchy printer, and shell completion generation.
2. **Viper (`github.com/spf13/viper`)**
   - **Role**: Hierarchical configuration management with deterministic precedence: Flags > Environment Variables (`RUNX_*`) > `runx.yaml` config file > Key/value defaults. Configured strictly for YAML parsing (`SetConfigType("yaml")`).

### Serialization & Validation
3. **`gopkg.in/yaml.v3`**
   - **Role**: Strict YAML parser for `runx.yaml` (Manifest V2). Configured with `Decoder.KnownFields(true)` to reject unknown or malformed manifest fields. TOML parsing is completely excluded.
4. **`encoding/json`** (Go Standard Library)
   - **Role**: High-speed JSON serialization for `--json` outputs, dry-run argument reporting, and PowerShell argument transport splatting.

### Standard Library Power Tools
5. **`embed` (`go:embed`)**
   - **Role**: Embeds CLI-owned agent skills (`guiho-s-runx.SKILL.md`), agent instructions, and documentation assets directly inside the compiled native binary.
6. **`os/exec` & `syscall`**
   - **Role**: Safe child command execution and background worker detachment (`SysProcAttr{Setsid: true}` on POSIX / `CREATE_NEW_PROCESS_GROUP` on Windows).
7. **`crypto/sha256` & `crypto/subtle`**
   - **Role**: Cryptographic hash calculation and constant-time checksum verification during `runx upgrade` binary verification.
8. **`net/http` with `time.Duration` Timeouts**
   - **Role**: Lightweight HTTP client with custom timeouts and user-agent headers for release discovery (`api.github.com`) and asset downloads.
9. **`path/filepath` & `os`**
   - **Role**: Atomic file operations (write-to-temp + rename) for skill reconciliation, cache updates, and binary upgrades without partial write corruption.

### Testing & Build DevOps
10. **Testify (`github.com/stretchr/testify`)**
    - **Role**: Assertion helpers (`assert`, `require`, `suite`) for testing parity against the existing 90 Bun unit test cases.
11. **`CGO_ENABLED=0` Pure-Go Compilation**
    - **Role**: Guarantees completely static, self-contained binaries with zero C-library dependencies.

---

## 3. Go Module & Coexistent Repository Architecture

The Go codebase lives in the repository root alongside the existing TypeScript code:

```text
runx/
├── main.go                       # Primary Go CLI entry point
├── go.mod                        # Go module (github.com/CGuiho/runx)
├── go.sum                        # Go dependency checksums
├── cmd/                          # Cobra Command Tree & Flag Routing
│   ├── root.go                   # Root command 'runx' & Viper YAML configuration
│   ├── list.go                   # 'runx list' command
│   ├── describe.go               # 'runx describe' command
│   ├── run.go                    # 'runx run' command with post-selector argument forwarding
│   ├── check.go                  # 'runx check' command
│   ├── init.go                   # 'runx init' command (interactive YAML catalog generator)
│   ├── agent.go                  # 'runx agent' root command
│   ├── agent_skill.go            # 'runx agent skill' (install, uninstall, update, list, show)
│   ├── agent_instruction.go      # 'runx agent instruction' (apply, remove, update, show)
│   ├── agent_prompt.go           # 'runx agent prompt' (list, show)
│   ├── upgrade.go                # 'runx upgrade' root command & options
│   ├── upgrade_check.go          # 'runx upgrade check' subcommand
│   ├── upgrade_list.go           # 'runx upgrade list' subcommand
│   ├── uninstall.go              # 'runx uninstall' command
│   └── helptree.go               # '--help-tree', '--help-tree-depth', '--help-docs' handlers
├── pkg/
│   ├── manifest/                 # Manifest V2 YAML parser, composition engine & Viper mapping
│   ├── executor/                 # Cross-platform command spawner & argument transport
│   ├── update/                   # Detached update worker & cache management
│   ├── maintenance/              # Automatic agent maintenance worker
│   ├── welcome/                  # Platform-aware welcome window renderer
│   └── updater/                  # Self-upgrade, binary verification & rollback engine
├── embed/
│   └── skills/                   # Embedded agent skills (`go:embed`)
├── devops/
│   ├── build-binaries.go         # Go multi-target build script
│   └── verify-release-assets.go  # 14-asset verification script
│
├── source/                       # [PRESERVED] Existing TypeScript source files (DO NOT DELETE)
├── package.json                  # [PRESERVED] Existing package configuration (DO NOT DELETE)
├── bun.lock                      # [PRESERVED] Existing Bun lockfile (DO NOT DELETE)
└── tsconfig.json                 # [PRESERVED] Existing TypeScript configuration (DO NOT DELETE)
```

---

## 4. Canonical Command Tree Contract (`runx --help-tree`)

The Go migration will be considered **fully complete** when executing `runx --help-tree` against the compiled Go binary produces the exact hierarchy shown below:

```text
COMMAND TREE

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
└── --help-docs                           Emit Markdown documentation for this command.
```

---

## 5. Detailed Subcommand & Module Specifications

### 5.1. `runx list`
- **Purpose**: Lists all available command leaves and group hierarchies resolved from `runx.yaml`.
- **Flags**: `--cwd`, `--config`, `--format` (`text`|`json`), `--verbose`.
- **JSON Format**: Emits structured array of command objects containing `uid`, `full_id`, `description`, `run`, `cwd`, and `confirm`.

### 5.2. `runx describe <uid>`
- **Purpose**: Inspects metadata and execution configuration for a single catalog command without running it.
- **Flags**: `--cwd`, `--config`, `--format` (`text`|`json`), `--verbose`.

### 5.3. `runx run <selector> [--] [child arguments...]`
- **Purpose**: Executes a resolved catalog command.
- **Argument Forwarding Grammar**:
  - Tokens before `<selector>` belong to RunX (`--dry-run`, `--yes`, `--cwd`, `--config`, `--format`).
  - All tokens after `<selector>` are forwarded losslessly to the child process.
  - Option delimiter `--` is stripped if placed directly after `<selector>`.
- **Transport Drivers**:
  - **POSIX** (`sh -c`): Encapsulates parameters safely into positional array strings.
  - **Windows PowerShell**: Constructs JSON-backed splat arrays without parameter string interpolation.
  - **Windows cmd**: Short-lived environment context wrapper with delayed expansion disabled.
- **Dry-Run Mode (`--dry-run`)**: Prints resolved command execution payload (JSON or text) without spawning child processes.

### 5.4. `runx check`
- **Purpose**: Parses and validates `runx.yaml` against Manifest V2 schema rules without executing any command.
- **Validation Rules**:
  - Requires `version: "2"` and non-empty `namespace`.
  - Rejects unknown top-level or command fields via `KnownFields(true)`.
  - Validates group colocation depth (max 32).
  - Validates reciprocal child catalog mounts (`runx` and `parent` fields).

### 5.5. `runx init`
- **Purpose**: Generates a valid starter `runx.yaml` configuration in the target directory.
- **Rules**: Fails safely if `runx.yaml` already exists unless explicitly forced.

### 5.6. `runx agent`
- **Subcommands**:
  - `agent skill`: `install`, `uninstall`, `update`, `list`, `show`.
  - `agent instruction`: `apply`, `remove`, `update`, `show`.
  - `agent prompt`: `list`, `show`.
- **Embedded Source**: Uses `go:embed` to read bundled skill files (`embed/skills/guiho-s-runx.SKILL.md`) directly from binary memory.

### 5.7. `runx upgrade`
- **Subcommands**: `check`, `list`.
- **Flags**: `--version`, `--arch`, `--variant`, `--dry-run`, `--format`.
- **Atomic Process Replacement**:
  - Downloads target binary from GitHub Releases.
  - Verifies native executable magic header and SHA256 checksum.
  - Performs two-phase file rename on Windows (`.exe` -> `.old.exe` -> new `.exe`).
  - Executes fail-safe rollback if verification fails.

### 5.8. `runx uninstall`
- **Purpose**: Removes the installed native RunX binary from the system PATH.

---

## 6. Environment Variables Matrix

Viper binds environment variables with the `RUNX_` prefix:

| Environment Variable | Equivalent Flag / Config | Usage / Meaning |
| --- | --- | --- |
| `RUNX_CONFIG` | `--config <path>` | Path to explicit `runx.yaml` configuration file. |
| `RUNX_CWD` | `--cwd <path>` | Target effective working directory for execution. |
| `RUNX_FORMAT` | `--format <text\|json>` | Output format selection (`text` or `json`). |
| `RUNX_CACHE_DIR` | N/A | Override directory for update worker TTL cache. |
| `RUNX_NO_UPDATE_CHECK` | N/A | Disable background update worker scheduling when set to `1` or `true`. |
| `RUNX_NO_AGENT_MAINTENANCE` | N/A | Disable automatic agent skill reconciliation when set to `1` or `true`. |

---

## 7. Exit Code Policy

RunX enforces deterministic exit code behavior:

| Exit Code | Classification | Cause / Scenario |
| --- | --- | --- |
| `0` | **Success** | Command completed cleanly, `--help`, `--version`, `--dry-run`, or check passed. |
| `1` | **User / Usage / Config Error** | Missing required flags, invalid manifest YAML syntax, unresolvable command selector, unknown flag. |
| `2` | **Execution Failure** | Child process failed (preserves child process non-zero exit status where available) or transport spawn error. |
| `3` | **Network / Upgrade Error** | Transport failure during `runx upgrade` download or catalog fetch error. |

---

## 8. Master Phase-by-Phase Implementation Plan

### Phase 1: Project Scaffolding & Cobra Command Tree
- [x] Initialize `go.mod` module (`github.com/CGuiho/runx`).
- [x] Install Go dependencies (`github.com/spf13/cobra`, `github.com/spf13/viper`, `gopkg.in/yaml.v3`, `github.com/stretchr/testify`).
- [x] Create `embed/embed.go` with `go:embed` asset binding for `guiho-s-runx.SKILL.md`.
- [x] Scaffold root Cobra command tree (`cmd/root.go`, `run.go`, `init.go`, `check.go`, `list.go`, `upgrade.go`, `agent.go`).
- [x] Compile and verify basic binary execution (`go build -o bin/runx-go.exe .`).

### Phase 2: Manifest V2 Parser & Strict YAML Decoder (`pkg/manifest`)
- [x] Create `pkg/manifest/types.go` struct definitions.
- [x] Implement strict YAML parser in `pkg/manifest/parser.go` with `KnownFields(true)`.
- [x] Implement group composition & indexing in `pkg/manifest/composition.go`.
- [x] Write Testify unit test suite in `pkg/manifest/manifest_test.go`.
- [x] Run `go test ./pkg/manifest/...` (Passing cleanly).

### Phase 3: Argument Transport & Executor (`pkg/executor`)
- [x] Implement POSIX parameter transport (`buildPOSIXCommand`).
- [x] Implement PowerShell JSON-backed splatting (`buildPowerShellCommand`).
- [x] Connect `pkg/manifest` and `pkg/executor` inside `cmd/run.go`.
- [x] Implement `--dry-run` execution plan reporting.
- [x] Run `go test ./pkg/executor/...` (Passing cleanly).

### Phase 4: Full Subcommand Tree & Help Tree Generator (`cmd/`)
- [ ] Implement `cmd/describe.go` for single command inspection.
- [ ] Implement `cmd/uninstall.go` with `--dry-run` support.
- [ ] Implement `cmd/helptree.go` to render the exact `runx --help-tree` matrix matching Section 4.
- [ ] Implement full subcommand flag sets (`--cwd`, `--config`, `--format`, `--verbose`, `--help-docs`, `--help-tree-depth`).
- [ ] Implement `cmd/agent_skill.go`, `cmd/agent_instruction.go`, and `cmd/agent_prompt.go`.

### Phase 5: Detached Workers & Upgrade Engine (`pkg/update`, `pkg/maintenance`, `pkg/updater`)
- [ ] Implement `pkg/update/worker.go` for background update checks via process detachment (`SysProcAttr` on POSIX / `CREATE_NEW_PROCESS_GROUP` on Windows).
- [ ] Implement `pkg/maintenance/maintenance.go` for non-blocking automatic agent skill and `AGENTS.md` block reconciliation.
- [ ] Implement `pkg/updater/upgrade.go` for atomic binary replacement and rollback.

### Phase 6: Test Parity, Multi-Target Build DevOps & Verification
- [ ] Port remaining Bun test cases to Go tests (`go test ./...`) achieving 90/90 test parity.
- [ ] Update `devops/build-binaries.go` for Go multi-target cross-compilation (`CGO_ENABLED=0`).
- [ ] Verify 14-asset release matrix and direct installers (`devops/install.sh`, `devops/install.ps1`).

---

## 9. Release Matrix & Asset Delivery

Go static cross-compilation (`CGO_ENABLED=0`) covers all 12 platform targets:

| Target OS | Target Arch | Binary Name | Release Asset Name |
| --- | --- | --- | --- |
| Windows | `amd64` | `runx-windows-amd64.exe` | `runx-windows-amd64.exe` |
| Windows | `arm64` | `runx-windows-arm64.exe` | `runx-windows-arm64.exe` |
| Linux | `amd64` | `runx-linux-amd64` | `runx-linux-amd64` |
| Linux | `arm64` | `runx-linux-arm64` | `runx-linux-arm64` |
| Linux | `riscv64` | `runx-linux-riscv64` | `runx-linux-riscv64` |
| macOS | `amd64` | `runx-darwin-amd64` | `runx-darwin-amd64` |
| macOS | `arm64` | `runx-darwin-arm64` | `runx-darwin-arm64` |

Plus 2 agent markdown release assets (`guiho-s-runx.SKILL.md`, `guiho-s-runx.INSTRUCTIONS.md`), for a total of **14 release assets**.

---

## 10. Verification Plan

### Automated Verification
1. **Go Unit & Integration Tests**:
   ```bash
   go test -v ./...
   ```
2. **Help Tree Parity Check**:
   ```bash
   go run . --help-tree
   ```
   Must output the exact hierarchy defined in Section 4.

3. **Existing Bun Check Baseline**:
   ```bash
   bun run check
   ```

### Manual Verification
1. Verify `runx run <selector> --dry-run` against complex multi-group `runx.yaml` catalogs.
2. Test `runx upgrade` replacement and rollback on Windows and Linux.
3. Test direct installer bootstrap (`devops/install.sh | bash` and `devops/install.ps1`).
