---
name: Implement RunX Alpha Command Catalog
purpose: Define the expected outcome, constraints, and completion signals for task 1 in TODO.md.
description: Describes the first RunX CLI release and the safety, documentation, and validation constraints that govern it.
created: 2026-07-12
flags:
  - implementation
tags:
  - todo
  - cli
keywords:
  - runx
  - alpha
owner: runx-todo
---

# Implement RunX Alpha Command Catalog

## Todo Index

- Task: `1. Implement RunX Alpha Command Catalog`
- Status: completed
- Index: [TODO.md](../../TODO.md)

## Outcome

Deliver a documented, locally executable RunX alpha that meets the approved
command-catalog acceptance criteria without publishing a package.

## Watch-outs

- Preserve read-only behavior for list, describe, check, and dry run paths.
- Never place secrets in a manifest or run confirmation-gated commands without
  `--yes`.
- Keep every new module and document represented by XDocs metadata.
- Use Mirror for the version transition and keep commits one file at a time.
