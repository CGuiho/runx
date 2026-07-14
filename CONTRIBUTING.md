---
name: Contributing to RunX
purpose: Explain the repository rules for contributing safely to RunX.
description: Covers Bun validation, documentation, agent-skill, and release workflow expectations.
created: 2026-07-12
flags: []
tags:
  - documentation
  - contributing
keywords:
  - runx
  - contributing
owner: runx
---

# Contributing to RunX

RunX is a Bun/TypeScript CLI. Keep changes focused and preserve the distinction
between catalog inspection and command execution.

1. Install dependencies with `bun install`.
2. Run `bun run typecheck` and `bun test` before opening a pull request.
3. Update `DOCS.md`, the bundled `guiho-s-runx` skill, and XDocs descriptors
   whenever CLI behavior, manifest fields, installers, or agent workflows change.
4. Preserve the approved tag-driven release workflow: GitHub Release binaries
   and npm trusted publishing run only from protected version tags in the
   `production` environment.
5. Use Mirror for version planning and version application; do not edit a
   release version manually.
