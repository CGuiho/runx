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

`cli.ts` defines a Citty command tree. Citty parses options and positionals,
resolves nested commands and aliases, and renders ordinary usage. A bounded
command-map adapter preserves the legacy `runx <selector>` shorthand while
known Citty commands retain precedence. Manifest commands use `manifest.ts` to
discover the nearest `runx.yaml`, parse YAML, validate a strict TypeBox schema,
and resolve a selector. `render.ts` presents text or JSON. `executor.ts` spawns
exactly one configured shell command only for a real run.

## Boundaries

- Inspection paths never call the executor.
- Help, version, and CLI usage failures do not discover a manifest.
- `source/flags.ts` no longer exists; RunX has no second token parser or manual
  execution router behind Citty.
- `agents.ts` copies the embedded or packaged skill into the requested standard
  or Claude-compatible skill location.
- `self-management.ts` queries GitHub Releases and only replaces/removes a
  native executable, never a Bun development process.
- Native builds register `embedded-resources.ts`, which includes the skill.

## Distribution

Bun bundles Citty and compiles standalone binaries without a Node.js runtime.
Direct installers select a GitHub Release asset; CI validates source and the
native release matrix, while protected tags publish release artifacts.
