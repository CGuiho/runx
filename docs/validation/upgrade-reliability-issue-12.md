---
name: RunX Upgrade Reliability Issue 12 Validation
purpose: Record current-main evidence for GitHub issue 12.
description: Captures focused and full transaction, catalog, routing, native, live GitHub, XDocs, and Git integrity checks.
created: 2026-07-19
flags:
  - validated
tags:
  - validation
  - cli
  - reliability
keywords:
  - RunX
  - issue 12
  - 50 tests
  - 14 assets
owner: runx-validation
---

# RunX Upgrade Reliability Issue 12 Validation

## Summary

All local completion gates for GitHub issue #12 passed after correcting the
complete-list default and Citty parent/subcommand routing. No validation
blocker remains.

## Commands Run

| Command or check | Result |
| --- | --- |
| `bun run typecheck` | Passed |
| focused catalog, transaction, reporting, and CLI tests | Passed: 31 tests, 222 assertions |
| `bun test --timeout 30000` | Passed: 50 tests, 342 assertions |
| `bun run build` | Passed |
| `bun run binary` | Passed |
| `bun run binaries` | Passed: twelve native targets |
| `bun run verify-assets` | Passed: exactly fourteen unique assets |
| live source `upgrade list --format json` | Passed: one JSON document, ten current published releases, current/latest/asset metadata |
| live source `upgrade check --format json` | Passed: one JSON document, no parent upgrade envelope |
| strict XDocs metadata and doctor | Passed for source and docs: zero errors and zero warnings |
| `git diff --check` | Passed |

## Transaction Evidence

- The self-management suite renames a running Windows executable, places the
  new binary at the canonical path, executes that path with `--version`, and
  reports success only after the exact target is observed.
- Controlled second-rename, target-version mismatch, and rollback-failure tests
  prove nonzero structured outcomes and deterministic restoration behavior.
- Plan and phase callbacks prove plan output precedes download work and
  replacement starts only after download and native validation.

## Catalog Evidence

- GitHub `Link` pagination, SemVer/prerelease ordering, channel labels,
  malformed payloads, candidate ordering, and missing compatible assets are
  covered by focused tests.
- A real Citty-route regression follows two mocked release pages, includes an
  alpha release without an opt-in flag, parses exactly one JSON document, and
  proves the parent upgrade action did not execute.
- Live GitHub-backed list and check routes each returned one parseable document
  with exit code zero.

## Skipped Checks

- No installed user executable was replaced.
- No package publication, GitHub Release creation, Mirror version, or tag was
  performed during this issue unit.

## References

- [Implementation review](../reviews/implementation/upgrade-reliability-issue-12-review.md)
- [Implementation plan](../plans/upgrade-reliability-implementation.md)
