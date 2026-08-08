---
name: RunX Go Rewrite RFC (Issue #22)
purpose: Define the authoritative completed migration from Bun/TypeScript/Citty to Go/Cobra.
description: Final architecture, manifest, startup, agent, upgrade, distribution, release, and verification contract for the native RunX CLI.
created: 2026-07-24
updated: 2026-07-26
flags:
  - accepted
  - implemented
tags:
  - rfc
  - go
  - cobra
  - cli
keywords:
  - issue 22
  - go rewrite
  - manifest v2
  - release matrix
owner: runx-rfc
---

# GUIHO RFC: RunX Go Rewrite

## Status And Authority

This RFC is accepted and implemented. It supersedes the initial issue-22 draft
where that draft required Viper, retained Bun as production authority, or
described the retired 14-asset/V3 release matrix.

The authoritative implementation is `main.go`, `cmd/`, `pkg/`, and `embed/`.
Cobra is the only command router. Configuration uses typed Go structs and
strict `go.yaml.in/yaml/v3` decoding; Viper is intentionally absent because
implicit environment and fallback behavior conflicts with RunX's exact
single-file precedence.

The former TypeScript implementation remains in `source/` as historical
reference. It is not built, tested, packaged as domain logic, or used by CI and
release publication. The Node npm script is only a checksum-verifying native
binary downloader and process delegate.

## Product Contract

RunX provides a language-agnostic `runx.yaml` command catalog. Inspection is
side-effect free; only `runx run` may execute a configured command. The CLI is
pre-1.0 and the migration intentionally removes incompatible legacy aliases,
implicit configuration behavior, V1 manifests, and old asset names.

### Command Tree

```text
runx
├── list
├── describe <uid-or-selector>
├── run [options] <selector> [--] [child arguments...]
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

Every scope has `-h`/`--help`, `--help-tree`,
`--help-tree-depth <positive-integer>`, and `--help-docs`. Root alone has
`-v`/`--version`. Tree and Markdown output traverse the live Cobra tree.

### Startup

A bare invocation synchronously and idempotently installs the embedded skill in
both global tool locations and reconciles the bounded RunX block at the current
repository root before it prints `Hello Windows - runx v<version>`. Both
existing instruction files are updated; otherwise the one existing file is
used, or `AGENTS.md` is created. Unmanaged content and line endings are
preserved, malformed markers fail safely, and no catalog command or network
request runs. Help, version, agent-management, uninstall, and non-repository
paths skip repository bootstrap. Foreground startup otherwise reads the typed
cache, optionally prints a newer stable notice, and starts recursion-safe hidden
workers. Update checks use a cache-scoped lease, finite timeout, stale recovery,
and atomic cache replacement.

### Configuration And Manifest

Configuration precedence is exact:

1. explicit `--config`;
2. effective-cwd `runx.yaml`;
3. `~/.guiho/runx/runx.yaml`.

No ancestor search, file merging, or implicit environment override is allowed.

Manifest v2 requires semantic version 2.x, an identifier-safe namespace, a
contained scripts subdirectory, and recursive commands. Leaves require stable
UID and ID, summary, description, and command. Groups require exactly one of
nested commands or a child `runx` reference. Strict decoding rejects unknown
fields and multiple YAML documents.

Local and GitHub child catalogs require reciprocal parents. Foreign catalogs
are HTTPS GitHub blob/raw resources, limited to one MiB and ten seconds, and
cannot escape an owner/repository/ref through relative references. Graph depth
is limited to 32 and cycles fail closed. UIDs remain globally unique and
canonical selectors remain unique across the composed graph. IDs are scoped by
their containing group; an unqualified ID shorthand is available only when it
has one owner. Exact UIDs resolve before canonical selectors or ID shorthands,
so a UID may equal another command's ID without ambiguity. Command cwd and
scripts paths must remain contained by their catalog base.

### Output And Execution

JSON mode writes exactly one deterministic document to stdout. Diagnostics use
stderr. `check`, `list`, `describe`, help, agent commands, and dry runs never
spawn catalog commands.

RunX-owned `run` flags precede the selector. Flag parsing stops at the selector;
all later tokens are forwarded losslessly. POSIX shells use positional
parameters, PowerShell uses JSON-backed splatting, and cmd uses environment-
backed transport. `confirm: always` requires `--yes`. Child exit codes are
preserved.

### Agent Resources

The production binary embeds the RunX skill and prompt. Explicit skill actions
target both `.agents` and `.claude`, globally by default and locally with
`--local`. Instruction actions use one bounded managed block while preserving
all unrelated text. Writes use temporary files and atomic replacement.

### Upgrade And Installation

Release reads use typed GitHub responses and finite HTTP clients. The embedded
build target selects the exact upgrade asset, preserving ARMv6 versus ARMv7.
Every upgrade downloads and validates `checksums.txt` before replacement.
Unix replacement uses same-filesystem staged renames and rollback. Windows
starts a hidden helper that waits for the current process to exit, then replaces,
verifies, cleans up, or restores the previous binary.

Both installers detect the exact platform target, display metadata and download
progress, verify binary and skill ZIP checksums, replace transactionally,
install both global skill copies, reconcile existing instructions, and verify
the installed version.

## Release Matrix Decision

The current shared GUIHO Go CLI engineering contract resolves the earlier
asset drift. RunX produces eight pure-Go binaries:

| Asset | Build controls |
| --- | --- |
| `runx-linux-amd64` | `GOOS=linux GOARCH=amd64 GOAMD64=v1 CGO_ENABLED=0` |
| `runx-linux-arm64` | `GOOS=linux GOARCH=arm64 GOARM64=v8.0 CGO_ENABLED=0` |
| `runx-linux-armv7` | `GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0` |
| `runx-linux-armv6` | `GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0` |
| `runx-darwin-amd64` | `GOOS=darwin GOARCH=amd64 GOAMD64=v1 CGO_ENABLED=0` |
| `runx-darwin-arm64` | `GOOS=darwin GOARCH=arm64 GOARM64=v8.0 CGO_ENABLED=0` |
| `runx-windows-amd64.exe` | `GOOS=windows GOARCH=amd64 GOAMD64=v1 CGO_ENABLED=0` |
| `runx-windows-arm64.exe` | `GOOS=windows GOARCH=arm64 GOARM64=v8.0 CGO_ENABLED=0` |

`guiho-s-runx.zip`, `guiho-i-runx.md`, and `checksums.txt` make exactly 11
release artifacts. The checksum manifest is deterministic and does not hash
itself. AMD64 V2/V3/V4 variants are excluded without benchmarks and a separate
contract decision.

## Implementation Checklist

### Phase 1: Production Architecture

- [x] Pin the Go module and dependencies.
- [x] Make `main.go` a thin build-info and exit-code entrypoint.
- [x] Construct fresh Cobra trees through `NewRootCommand` with injected I/O,
  HTTP, clock, filesystem location, executable, and worker launch dependencies.
- [x] Remove Viper and global mutable command routing.

### Phase 2: Manifest And Execution

- [x] Strictly decode and semantically validate manifest v2.
- [x] Implement exact configuration precedence without parent search.
- [x] Compose reciprocal local and bounded GitHub child catalogs.
- [x] Validate identity collisions, graph bounds, scripts, cwd, shell, and
  confirmation before execution.
- [x] Preserve post-selector child arguments and child exit codes.

### Phase 3: CLI And Agent Contract

- [x] Implement real `check`, `list`, `describe`, `run`, and `init` behavior.
- [x] Generate help tree and Markdown from the live Cobra tree.
- [x] Support only `-h` and root `-v` short aliases.
- [x] Embed skill and prompt resources and implement idempotent agent actions.
- [x] Bootstrap global skills and repository instructions on successful bare invocation.

### Phase 4: Lifecycle And Upgrade

- [x] Wire cached notices and detached update/maintenance workers.
- [x] Coalesce update workers and atomically replace typed cache state.
- [x] Resolve upgrades from embedded build targets.
- [x] Verify release checksums and native formats before replacement.
- [x] Implement Unix rollback and staged Windows post-exit replacement.

### Phase 5: Distribution And Automation

- [x] Build and verify the standard 11 artifacts.
- [x] Update Bash and PowerShell installers for all standard targets.
- [x] Make the npm entrypoint a checksum-verifying native bootstrap only.
- [x] Make CI and protected-tag publication validate the Go implementation and
  exact Go release matrix without Bun/TypeScript authority.

### Phase 6: Documentation And Verification

- [x] Reconcile AGENTS, README, DOCS, TODO, XDocs, and validation records.
- [x] Run formatting, module-tidy, complete tests, vet, native smokes, XDocs,
  and all eight cross-builds.
- [x] Defer Mirror version application, tag, push, publication, and external
  release approval to a separately authorized release task.

## Exit Policy

RunX uses `0` success, `1` unexpected operation failure, `2` usage or approval
failure, `3` configuration failure, `4` release/network failure, and `5`
installation/filesystem failure. `run` preserves explicit child exit codes.
