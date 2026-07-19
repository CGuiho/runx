---
name: RunX CLI Architecture
purpose: Describe the component boundaries and data flow of the RunX alpha CLI.
description: Explains Citty routing, manifest discovery, TypeBox validation, selector resolution, execution, agent installation, and native distribution.
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

`cli.ts` defines one Citty command tree. Citty parses options and positionals,
resolves nested commands, and renders usage. Manifest commands use
`configuration.ts` to resolve `runx.yaml` by explicit path, effective cwd, then
global fallback; TypeBox validates the complete shape. `render.ts` presents
text or JSON. `executor.ts` spawns exactly one configured shell command only for
a real `runx run`.

After an ordinary command routes through Citty, `agent-maintenance.ts` receives
the decoded effective cwd and detaches a hidden worker. The worker compares
embedded skill content with both global copies and atomically reconciles one
nearest `AGENTS.md` block. Explicit agent-resource and uninstall commands are
excluded so intentional removal is not reversed.

## Boundaries

- Inspection paths never call the executor.
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
