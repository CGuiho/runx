---
name: RunX CLI Architecture
purpose: Describe the component boundaries and data flow of the RunX alpha CLI.
description: Explains manifest discovery, TypeBox validation, selector resolution, execution, agent installation, and native distribution.
created: 2026-07-12
flags:
  - approved
tags:
  - architecture
  - cli
keywords:
  - runx
  - typebox
  - bun
owner: runx-architecture
---

# RunX CLI Architecture

## Flow

`cli.ts` parses a command and global flags. Manifest commands use `manifest.ts`
to discover the nearest `runx.yaml`, parse YAML, validate a strict TypeBox
schema, and resolve a selector. `render.ts` presents text or JSON. `executor.ts`
spawns exactly one configured shell command only for a real run.

## Boundaries

- Inspection paths never call the executor.
- `agents.ts` copies the embedded or packaged skill into the requested standard
  or Claude-compatible skill location.
- `self-management.ts` queries GitHub Releases and only replaces/removes a
  native executable, never a Bun development process.
- Native builds register `embedded-resources.ts`, which includes the skill.

## Distribution

Bun compiles standalone binaries. Direct installers select a GitHub Release
asset; CI validates source and a native build but does not publish artifacts.
