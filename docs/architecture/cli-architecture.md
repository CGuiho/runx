---
name: RunX CLI Architecture
purpose: Describe the component boundaries and data flow of the RunX alpha CLI.
description: Explains Citty routing, the run argument boundary, manifest discovery, TypeBox validation, shell-safe execution, startup rendering, agent installation, and native distribution.
created: 2026-07-12
flags:
  - approved
tags:
  - architecture
  - cli
keywords:
  - runx
  - citty
  - typebox
  - bun
owner: runx-architecture
---

# RunX CLI Architecture

## Flow

`cli.ts` defines one Citty command tree. `execution-arguments.ts` identifies the
selector boundary only for `runx run`, then Citty parses the RunX-owned prefix,
resolves nested commands, and renders usage. The immutable post-selector array
is passed separately to the executor. Manifest commands use
`configuration.ts` to resolve `runx.yaml` by explicit path, effective cwd, then
global fallback; TypeBox validates the complete shape. `render.ts` presents
text or JSON. `executor.ts` spawns exactly one configured shell command only for
a real `runx run`; POSIX positional parameters, PowerShell JSON-backed
splatting, and a short-lived cmd wrapper keep child values out of shell source.

`welcome.ts` purely renders bare invocation. `update-cache.ts` validates cached
SemVer before supplying an optional post-body warning and completes only local
worker scheduling in the foreground; the detached worker owns bounded remote
work.

After an ordinary command routes through Citty, `agent-maintenance.ts` receives
the decoded effective cwd and detaches a hidden worker. The worker compares
embedded skill content with both global copies and atomically reconciles one
nearest `AGENTS.md` block. Explicit agent-resource and uninstall commands are
excluded so intentional removal is not reversed.

## Boundaries

- Inspection paths never call the executor.
- RunX options precede the selector; post-selector flags belong to the child.
- Dry-run text and JSON expose forwarded arguments without spawning.
- Help, version, and CLI usage failures do not discover a manifest.
- `source/flags.ts` no longer exists; RunX has no second token parser or manual
  execution router behind Citty.
- `agents.ts` owns explicit dual-tool actions plus idempotent, legacy-aware
  automatic skill and instruction reconciliation.
- `storage.ts` performs same-directory temporary writes and atomic replacement
  for automatic resource maintenance.
- Hidden update and agent-maintenance routes never appear in public help and
  cannot recursively schedule themselves.
- `self-management.ts` queries GitHub Releases and only replaces/removes a
  native executable, never a Bun development process.
- Native builds register `embedded-resources.ts`, which includes the skill.

## Distribution

Bun bundles Citty and compiles standalone binaries without a Node.js runtime.
Direct installers select a GitHub Release asset; CI validates source and the
native release matrix, while protected tags publish release artifacts.
