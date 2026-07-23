---
name: RunX Manifest V2 Composition Validation
purpose: Record implementation and release evidence for GitHub issue 26.
description: Tracks schema, composition, CLI, complete tests, builds, assets, XDocs, CI, release, installers, and public acceptance.
created: 2026-07-23
flags:
  - validated
  - released
tags:
  - validation
  - cli
keywords:
  - manifest v2
  - issue 26
owner: runx-validation
---

# RunX Manifest V2 Composition Validation

## Scope

Validated namespace identity, recursive colocated groups, local and GitHub
child mounts, aliases, reciprocal parents, selectors, provenance, cwd behavior,
breaking migration, CLI surfaces, and distribution. Issue 22 is excluded.

## Evidence

| Check | Result |
| --- | --- |
| `bun run typecheck` | Passed. |
| `bun test` | Passed: 90 tests, 571 assertions. |
| Focused configuration/init/CLI tests | Passed: 32 tests, 221 assertions, including all eleven independent review blocker categories. |
| `bun run build` | Passed. |
| `bun run binaries` | Passed: twelve native targets. |
| `bun run verify-assets` | Passed: exactly fourteen assets, including both `.md` files. |
| `xdocs scan --strict` | Passed. |
| `xdocs doctor` | Passed: zero errors and zero warnings. |
| Independent implementation review | Cleared exact pushed HEAD `376404c`; all eleven direct adversarial probes passed. |
| Main correction CI | Passed: [30047576118](https://github.com/CGuiho/runx/actions/runs/30047576118), including the full Windows suite with explicit process-test budgets. |
| Release CI | Passed: [30047684096](https://github.com/CGuiho/runx/actions/runs/30047684096). |
| Publish | Passed: [30047684524](https://github.com/CGuiho/runx/actions/runs/30047684524); exact fourteen assets, exact-version installer, npm publication, and version-scoped notes. |
| GitHub issue | [Issue 26](https://github.com/CGuiho/runx/issues/26) closed as completed with evidence. |

## Release Outcome

RunX 0.7.0 delivered the accepted manifest-v2 runtime but exposed a parallel
CI race by comparing an unreleased package version with the mutable latest
release. RunX 0.7.1 corrected workflow ownership: ordinary CI performs a generic
latest-release smoke, while publish verifies the immutable tag and exact version
only after the fourteen assets exist. RunX 0.7.2 then gave the two Windows
process tests explicit fifteen-second budgets so runner startup variance cannot
produce false five-second failures. Final Linux and Windows CI, npm provenance,
the exact fourteen assets, tag-pinned public installation, version-scoped notes,
and issue closure are complete.
