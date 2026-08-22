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
├── reveal <uid-or-selector-or-index>
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
On Windows, `shell: auto` selects the resolved `bash` executable only when the
inherited `MSYSTEM` marker identifies a Git Bash/MSYS caller and the executable
is not the Windows System32/WSL launcher. Otherwise it uses `cmd.exe`; on
non-Windows platforms it uses `sh`. Explicit shell values always bypass this
caller inference, and configured command text and paths are passed unchanged.

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

`check`, `list`, `describe`, `reveal`, help, agent operations, and `run --dry-run`
never spawn a configured command. `--format json` emits one stable JSON document
to stdout; diagnostics use stderr. Human `list` output uses padded columns so
the `IDX`, `UID`, `SELECTOR`, and `SUMMARY` fields align across supported
terminals.

`runx reveal <uid-or-selector-or-index>` accepts exactly one selector and the
catalog location/diagnostic flags `--cwd`, `--config`, and `--verbose`. It
reuses the catalog resolver used by `runx run` and writes the selected manifest
command verbatim followed by one newline to stdout. Reveal does not accept
`--format`, child arguments, `--yes`, or `--dry-run`, and never executes or
confirms the selected command.

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

Protocol-v1 installations upgrade through `pkg/lifecycle`: the whole-release
engine stages every release artifact under `$HOME/.guiho/.temp/runx-upgrade-<id>`,
verifies each against `checksums.txt`, runs the staged payload's hidden
self-test, installs an immutable version directory, atomically swaps the
`current.json` pointer while retaining the previous verified version, verifies
the result through the stable launcher, and rolls back the pointer and version
directory on any failure. The running executable is never replaced. Legacy
direct-binary installations keep the verified staged-replacement updater; every
terminal outcome prints the two-line reinstallation recovery block pinned to
the resolved target version.

## Installers And Uninstallers

`devops/install.sh` detects AMD64, ARM64, ARMv7, ARMv6, and Darwin targets.
`devops/install.ps1` detects Windows AMD64 and ARM64. Both accept full-name
only `--version`/`--channel` (PowerShell `-Version`/`-Channel`) selection with
release-catalog pagination, stage downloads under
`$HOME/.guiho/.temp/runx-install-<id>`, verify checksums and the
`artifacts.json` manifest for every artifact, self-test the payload before
activation, install immutable payloads plus the stable launcher into the
canonical `.guiho` layout, activate via atomic pointer swap with previous-version
rollback, preserve configuration and data on reinstall, update the user PATH
idempotently, and verify the installed launcher version.

`devops/uninstall.sh`, `devops/uninstall.ps1`, and `runx uninstall` share one
uninstallation contract: a REMOVE/PRESERVE plan, `--preserve-config`,
`--preserve-data`, `--dry-run`, and `--yes`; fail-closed noninteractive use;
managed-block-only instruction removal; and preservation of shared
`.guiho/` infrastructure.

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

`devops/build-binaries.go` produces the protocol-v1 matrix: eight immutable
payload executables (`runx-payload-<os>-<arch>` for Linux AMD64/ARM64/ARMv7/ARMv6,
Darwin AMD64/ARM64, Windows AMD64/ARM64), eight stable launchers
(`runx-launcher-<os>-<arch>`), `guiho-s-runx.zip`, `guiho-i-runx.md`,
`runx.schema.json`, `runx.global.schema.json`, an `artifacts.json` ownership
manifest declaring every artifact's digest, kind, and canonical installed path,
and a `checksums.txt` covering every artifact except itself.

All binaries use `CGO_ENABLED=0`; AMD64 uses V1, ARM64 uses V8.0, and 32-bit
ARM uses its named GOARM level. AMD64 V2/V3/V4 and the legacy 11-asset direct
binary names are not supported by the current contract.
