---
name: RunX CLI Architecture
purpose: Describe the component boundaries and data flow of the native RunX CLI.
description: Explains Cobra routing, strict typed YAML, manifest composition, shell-safe execution, lifecycle workers, agent resources, verified upgrades, and Go distribution.
created: 2026-07-12
updated: 2026-07-26
flags:
  - approved
tags:
  - architecture
  - cli
  - go
keywords:
  - runx
  - cobra
  - manifest v2
  - lifecycle workers
owner: runx-architecture
---

# RunX CLI Architecture

## Command Flow

`main.go` supplies immutable build information to `cmd.NewRootCommand`. Each
call constructs a fresh, dependency-injected Cobra tree; there is no package
global router. Cobra owns all parsing, usage, help, and routing. RunX-owned
`run` options must precede the selector because flag parsing stops at the first
positional token. Every later token belongs to the configured child command.

`pkg/manifest` resolves configuration from an explicit `--config`, the
effective working directory's `runx.yaml`, then
`~/.guiho/runx/runx.yaml`. It never searches ancestors. A strict YAML decoder
rejects unknown fields and validates the typed manifest-v2 graph before
`list`, `describe`, `check`, or `run` can use it. Local and bounded GitHub child
catalogs require reciprocal parent declarations; graph depth, document size,
cycles, identifiers, selectors, and UIDs are validated deterministically. UIDs
resolve before canonical selectors and unambiguous group-scoped ID shorthands;
a UID may equal another command's ID, while duplicate ID shorthands remain
ambiguous and fail closed.

Inspection commands render stable text or JSON and never reach the executor.
Only `run` resolves a command leaf and invokes `pkg/executor`. POSIX positional
parameters, PowerShell JSON-backed splatting, and a short-lived cmd wrapper keep
forwarded child values out of generated shell source. Explicit child exit codes
are preserved.

## Lifecycle And Agent Resources

Bare invocation synchronously performs local-only agent bootstrap before
rendering its welcome. It installs both global skill copies and, when inside a
repository, atomically reconciles both existing root instruction files, the one
that exists, or a new `AGENTS.md`. Marker validation, line-ending preservation,
and exact content comparisons make the operation safe and idempotent.

Other ordinary foreground commands schedule hidden update and maintenance
workers. The workers use finite network timeouts, freshness windows, and leases
so concurrent invocations coalesce. They never turn an otherwise successful
foreground command into a network-dependent one.

`pkg/maintenance` compares embedded resources with both supported global skill
locations and atomically reconciles bounded blocks in repository-root
`AGENTS.md` and `CLAUDE.md` according to which files exist. Explicit `runx agent` commands expose the same embedded skill,
instruction, and prompt resources for inspection or repair. Explicit removal
and uninstall paths do not schedule automatic reinstallation.

`pkg/update` maintains strict cached release metadata. `pkg/updater` downloads
the canonical target plus `checksums.txt`, verifies SHA-256 and native format,
then replaces transactionally. Windows stages a post-exit helper because a
running executable cannot replace itself; failure paths preserve or restore the
previous binary.

## Safety Boundaries

- Manifests are trusted executable code, but inspection and dry-run never spawn.
- Help, version, and usage failures do not require a manifest.
- Composition follows only explicit `runx` and `parent` edges.
- Hidden workers never appear in public help and cannot recursively schedule.
- The npm script is a checksum-verifying native bootstrap, not a domain runtime.
- Legacy `source/` TypeScript files are historical references only.

## Distribution

Go produces eight `CGO_ENABLED=0` executables: Linux AMD64, ARM64, ARMv7, and
ARMv6; Darwin AMD64 and ARM64; and Windows AMD64 and ARM64. The release adds the
bundled skill ZIP, instruction Markdown, and checksum manifest for exactly 11
assets. CI and protected publication workflows validate this same contract.
