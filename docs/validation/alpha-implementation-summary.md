---
name: RunX Alpha Implementation Summary
purpose: Summarize the first implemented RunX alpha and preserve its validation evidence.
description: Records delivered CLI behavior, agent and distribution support, checks run, and release boundaries.
created: 2026-07-12
flags:
  - validated
tags:
  - validation
  - implementation
keywords:
  - runx
  - alpha
  - summary
owner: runx-validation
---

# RunX Alpha Implementation Summary

## Implementation Summary

RunX is implemented as the open-source `@guiho/runx` Bun/TypeScript CLI. The
alpha discovers one nearest `runx.yaml`, validates it strictly with TypeBox,
lists documented commands through `runx list`, and runs an explicitly selected
command through `runx run` or `runx r`. A bare selector remains available as a
human shorthand; no-argument `runx` now shows the requested home page and usage.

The implementation includes UID/group-ID/index selector handling, JSON output,
describe, check, dry run, confirmation gates, help/help-tree/help-docs, bundled
`guiho-s-runx` skill installation, native upgrade and uninstall paths, Windows
and macOS/Linux direct installers, MIT licensing, and GitHub Actions CI. CI
validates only; npm and release publishing are intentionally absent.

## Validation Evidence

- `bun run typecheck` passed.
- `bun test` passed with four tests covering manifest discovery, selector
  behavior, flags, the `r` alias, dry runs, and real local command execution.
- `bun run build` passed and emitted the ignored TypeScript library output.
- `bun run binary` passed and emitted the ignored Windows native executable.
- The Bun source CLI displayed `--help`, `--help-tree`, and `--help-docs`.
- `xdocs tree`, `xdocs scan`, and `xdocs doctor --format json` passed. The
  doctor reported only metadata warnings for the special-purpose skill file;
  there were no descriptor, document, or tree errors.
- `mirror config check` passed. Mirror applied the local `0.1.0` minor release
  and created `@guiho/runx@0.1.0`; the configuration did not push or publish.

## Release Boundary

The initial external capability warranted and received a Mirror minor transition
from `0.0.0` to `0.1.0`. It created the configured local version commit and Git
tag without pushing either because `mirror.config.toml` sets `push = false`.
