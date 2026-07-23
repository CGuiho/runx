---
name: Implement RunX Manifest V2 Composition
purpose: Define completion signals for TODO task 13 and GitHub issue 26.
description: Tracks recursive colocated commands, namespaces, parent-child catalogs, migration, tests, documentation, and public release acceptance.
created: 2026-07-23
flags:
  - in-progress
tags:
  - todo
  - cli
keywords:
  - manifest v2
  - issue 26
owner: runx-todo
---

# Implement RunX Manifest V2 Composition

## Todo Index

- Task: `13. Implement RunX Manifest V2 Composition`
- Status: in progress
- Index: [TODO.md](../../TODO.md)

## Outcome

RunX manifests colocate commands and nested groups, compose explicitly linked
local or GitHub child catalogs under renameable namespace aliases, and reject
ambiguous, cyclic, non-reciprocal, legacy, or unsafe graphs before execution.

## Acceptance Signals

- Manifest v2 TypeBox decoding and semantic graph tests pass.
- Init and every inspection/execution surface use the composed v2 model.
- Full suite, native builds, exact assets, XDocs, and public installers pass.
- A public release demonstrates nested local child execution and registry
  provenance; issue 26 receives evidence and closes.

## External Trackers

- GitHub: [CGuiho/runx#26](https://github.com/CGuiho/runx/issues/26) - status: open
