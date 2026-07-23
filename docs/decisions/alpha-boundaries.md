---
name: RunX Alpha Boundaries
purpose: Preserve the accepted product and technical decisions that constrain the first RunX implementation.
description: Records the manifest-first, local-only, safety, distribution, and public licensing decisions for the alpha.
created: 2026-07-12
flags:
  - accepted
tags:
  - decisions
keywords:
  - runx
  - yaml
  - safety
owner: runx-decisions
---

# RunX Alpha Boundaries

> **Current manifest note:** The alpha prohibition on manifest composition is
> superseded by [RunX Manifest V2 Composition](manifest-v2-composition.md).
> Other execution, safety, and distribution boundaries remain historical
> constraints unless a later decision explicitly supersedes them.

- Use one versioned YAML manifest named `runx.yaml`; do not merge manifests.
- RunX is implemented in Bun/TypeScript but consumes commands from any project
  language through a configured local shell.
- Default `runx` behavior is a home page; `runx list` is explicit discovery.
- Keep `runx r` as an alias for explicit `runx run`.
- Treat manifests as trusted executable code and keep secret storage outside
  the manifest.
- Use explicit confirmation fields rather than guessing safety from groups.
- Start as MIT-licensed open source with CI only and no npm publishing job.
- Keep task graphs, variables, remote execution, caching, and automatic imports
  out of the alpha.
