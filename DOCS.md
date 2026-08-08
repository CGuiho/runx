---
name: RunX CLI Reference
purpose: Define the complete current behavior of the native RunX CLI.
description: Cobra routing, strict manifest v2, configuration precedence, deterministic output, execution, agents, lifecycle workers, upgrades, installers, and release assets.
created: 2026-07-12
owner: runx
flags: []
tags:
  - cli
  - go
  - documentation
keywords:
  - runx
  - cobra
  - manifest v2
  - release
---

# RunX CLI Reference

## Runtime Authority

`main.go`, `cmd/`, and `pkg/` are the RunX production implementation. Cobra is
the only command and flag router. The Node npm entrypoint downloads and delegates
to a native Go release; it contains no domain logic. Retained Bun/TypeScript
files are legacy reference and are not used by CI or release publication.

The build injects `version`, `commit`, `buildDate`, and `buildTarget`. The target
is retained during self-upgrade so Linux ARMv6 and ARMv7 installations select
the same compatible architecture.

## Command Tree

```text
runx
├── list
├── describe <uid-or-selector-or-index>
├── run [options] <uid-or-selector-or-index> [--] [child arguments...]
├── check
├── init
├── agent
│   ├── skill install|uninstall|update|list|show
│   ├── instruction apply|remove|update|show
│   └── prompt list|show
├── upgrade
│   ├── check
│   └── list
└── uninstall
```

Every scope exposes `-h`/`--help`, `--help-tree`,
`--help-tree-depth <positive-integer>`, and `--help-docs`. Root alone exposes
`-v`/`--version`. The Unicode tree and Markdown help traverse the live Cobra
tree, so separate static command metadata cannot drift.

## Startup

With no arguments RunX prints:

```text
Hello Windows - runx v<version>
```

A successful bare invocation first performs a synchronous, local-only agent
bootstrap. It installs the embedded skill under both
`~/.agents/skills/guiho-s-runx` and `~/.claude/skills/guiho-s-runx`, locates the
current repository root, and reconciles the bounded RunX instruction block in
both existing `AGENTS.md` and `CLAUDE.md`, the one file that exists, or a newly
created `AGENTS.md`. It preserves unmanaged bytes and the existing line-ending
style, rejects malformed or duplicate markers, uses atomic replacement, and is
idempotent. A non-repository directory still receives global skills but no
instruction file. Help, version, agent-management, and uninstall paths skip the
bootstrap.

A newer stable version may add a validated notice read from
`~/.guiho/runx/cache.json`. Foreground startup performs no remote request.
Ordinary invocations start hidden detached workers for a bounded release check
and automatic agent maintenance. A cache lease coalesces simultaneous update
workers and stale leases recover automatically. Internal worker commands are
hidden from help and cannot recursively start workers.

## Configuration Resolution

Commands that use a catalog resolve exactly one YAML file:

1. `--config <path>`, relative to effective cwd when needed;
2. `<effective-cwd>/runx.yaml`;
3. `~/.guiho/runx/runx.yaml`.

RunX never searches ancestors and never merges configuration files. Successful
human and JSON inspection identifies the absolute configuration path.

## Manifest v2

The root fields are:

- `version`: semantic version whose major is `2`;
- `namespace`: identifier matching `^[a-z][a-z0-9-]*$`;
- `scripts.directory`: a relative subdirectory contained by the catalog;
- optional `parent`: a relative path or full HTTPS GitHub blob/raw URL;
- `commands`: recursive command leaves and groups.

A command leaf requires `uid`, `id`, `summary`, `description`, and `command`.
Optional fields are `cwd`, `shell`, `tags`, and `confirm`. Shell is `auto`,
`bash`, `sh`, `powershell`, or `cmd`; confirmation is `never` or `always`.

A group requires `group`, `summary`, and exactly one of nested `commands` or a
`runx` child reference. Child catalogs must declare the reciprocal parent.
Relative foreign references cannot escape their GitHub owner/repository/ref.
Foreign reads have a 10-second client deadline, a one-MiB limit, cycle
protection, and a maximum graph depth of 32.

`runx init` writes `.scripts` as the default `scripts.directory`; an explicitly
configured directory remains unchanged. Canonical selectors use slash-separated
group paths. UIDs are globally unique, canonical selectors are unique, and
unqualified ID shorthands are available only when their ID has one owner. A UID
may equal another command's group-scoped ID: exact UID lookup wins before a
canonical selector or unique ID shorthand. Duplicate unqualified IDs remain
ambiguous and fail instead of selecting an arbitrary command. Exact textual
identities resolve before a canonical positive-decimal `IDX` fallback. Numeric
indexes refer to the current resolved `runx list` order and are intended for
interactive use; stable UIDs remain the automation contract. Command and
scripts paths are containment-validated before `check`, `list`, or `describe`
succeeds.

YAML uses `go.yaml.in/yaml/v3` with `KnownFields(true)`. Multiple documents and
unknown fields are rejected before semantic validation.

## Inspection And Execution

`check`, `list`, `describe`, help, agent operations, and `run --dry-run` never
spawn a configured command. `--format json` emits one stable JSON document to
stdout; diagnostics use stderr. Human `list` output uses padded columns so the
`IDX`, `UID`, `SELECTOR`, and `SUMMARY` fields align across supported terminals.

Only `runx run` executes a selected command. RunX-owned options precede the
selector. Flag parsing stops at the selector and every later token is forwarded
without RunX reinterpretation. POSIX shells use positional parameters,
PowerShell uses JSON-backed splatting, and cmd uses environment-backed argument
transport. Child values are never interpolated into generated shell source.

For `confirm: always`, an interactive text invocation prints the exact
`runx run --yes <selector>` retry and asks `Are you sure? [y/N]`. Only `y` or
`yes`, case-insensitively, authorizes execution; Enter, EOF, and every other
answer decline without spawning. Noninteractive and JSON invocations never
prompt and fail closed with the same exact retry. The retry preserves RunX
options, the selected identity, the child delimiter, and child arguments.
Supplying `--yes` before the selector bypasses the prompt. Child process exit
codes are preserved.

## Agent Resources

The Go binary embeds the RunX skill and instruction prompt. Skill install/update
targets both `.agents/skills/guiho-s-runx` and
`.claude/skills/guiho-s-runx`; `--local` selects the effective project instead
of the user home. Writes use temporary files and Windows-safe atomic replacement.

Instruction actions preserve all unmanaged text and reconcile the bounded RunX
block in `AGENTS.md`, `CLAUDE.md`, both existing files, or a newly created
`AGENTS.md`. Malformed markers fail without modifying the instruction file.
Background automatic maintenance is silent and failure-isolated.

## Updates And Self-Upgrade

`upgrade check` and `upgrade list` read GitHub Releases through a finite-timeout
HTTP client. Listing exhausts API pagination before applying the visible
`--page`/`--size` slice. JSON retains complete release and compatible-asset
metadata.

`upgrade` resolves the embedded build target, downloads the exact binary and
`checksums.txt`, verifies SHA-256 and native executable magic, and only then
changes the installed executable. Unix replacement uses same-filesystem staged
renames with verification and rollback. A running Windows executable starts a
hidden helper that waits for the parent process to exit, replaces and verifies
the binary, restores the backup on failure, and records a recovery error log.
Every terminal failure reports an exact-version installer recovery command.

## Installers And npm Bootstrap

`devops/install.sh` detects AMD64, ARM64, ARMv7, ARMv6, and Darwin targets.
`devops/install.ps1` detects Windows AMD64 and ARM64. Both display target
metadata, use download progress, verify binary and skill archive checksums,
replace transactionally, install both global skill copies, reconcile existing
project instructions, and verify the installed version.

RunX 0.8 uses the retired Bun release contract and cannot discover current
native releases. When its updater reports that 0.8 is already current, migrate
with the unpinned direct installer rather than its pinned 0.8 recovery command:

```bash
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash
hash -r
runx --version
```

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
runx --version
```

Restart an existing shell if it still resolves the old executable. Git Bash can
load the installer-managed path immediately with `source ~/.bashrc`.

`scripts/runx-bin.mjs` is a Node-compatible npm bootstrap. It downloads the
package version's native artifact and checksum manifest into
`~/.guiho/runx/npm/<version>/`, verifies SHA-256, and delegates stdio,
environment, arguments, signals, and exit status.

## Exit Codes

| Code | Meaning |
| ---: | --- |
| `0` | Success. |
| `1` | Unexpected or operational failure. |
| `2` | Usage, flag, or approval failure. |
| `3` | Configuration resolution, decoding, or semantic validation failure. |
| `4` | Remote release or network failure. |
| `5` | Installation, upgrade, or filesystem mutation failure. |
| child | `run` preserves the configured child process exit code. |

## Release Contract

`devops/build-binaries.go` produces exactly:

```text
runx-linux-amd64
runx-linux-arm64
runx-linux-armv7
runx-linux-armv6
runx-darwin-amd64
runx-darwin-arm64
runx-windows-amd64.exe
runx-windows-arm64.exe
guiho-s-runx.zip
guiho-i-runx.md
checksums.txt
```

All binaries use `CGO_ENABLED=0`; AMD64 uses V1, ARM64 uses V8.0, and 32-bit
ARM uses its named GOARM level. AMD64 V2/V3/V4 and the former 14-asset matrix
are not supported by the current contract.
