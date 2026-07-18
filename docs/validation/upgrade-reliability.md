---
name: RunX Upgrade Reliability Validation
purpose: Record release-ready verification evidence for GitHub issues 12 and 13.
description: Results for release pagination, progress reporting, transactional replacement, rollback, recovery commands, direct installers, native builds, live release smokes, and XDocs integrity.
created: 2026-07-15
flags:
  - superseded
  - needs-revalidation
tags:
  - validation
  - cli
  - reliability
keywords:
  - runx upgrade
  - upgrade list
  - windows replacement
  - recovery install
  - issue 12
  - issue 13
owner: runx-validation
---

# RunX Upgrade Reliability Validation

## Current Status

The first validation pass below is superseded by independent review findings.
The implementation has since been corrected for exact JSON envelope fidelity,
post-backup rename/rollback state, downgrade prevention, installer failure
classification, Ubuntu portability, and executable controlled installer tests.
Candidate verification now has a bounded ten-second timeout in the CLI and both
direct installers. Those corrections still require the complete release gate
to be rerun before this document or the TODO task can return to
`validated`/`completed`.

No package publication, version apply, or tag operation was performed.

## Superseded Automated Gate

| Command | Result |
| --- | --- |
| `bun run typecheck` | Passed. |
| `bun test` | Passed: 37 tests, 0 failures, 245 assertions across 8 files. |
| `bun run build` | Passed; TypeScript library output compiled. |
| `bun run binary` | Passed; local native executable compiled. |
| `bun run binaries` | Passed; all 12 Linux, macOS, and Windows release assets compiled and were non-empty. |
| `bun test devops/installers.spec.ts` | Superseded: the earlier suite relied partly on source-pattern checks. The replacement suite executes controlled native installer and recovery scenarios and awaits revalidation. |

Generated `library/` and `bin/` outputs remain ignored and are not part of the
committed change.

## Superseded Behavior Evidence

- Release-catalog tests follow `Link: rel="next"`, reject malformed/non-2xx
  responses, apply SemVer prerelease precedence, retain invalid tags last, and
  enforce baseline/default/modern asset candidates.
- Reporting tests prove the complete plan precedes `Downloading`, which
  precedes the download body and then `Validating`, `Replacing`, and
  `Verifying`.
- The Windows test runs a copied executable, keeps it mapped in a live process,
  renames it, installs the new canonical executable immediately, verifies the
  exact version, and allows only old-image cleanup to be deferred.
- Failure tests reject non-native downloads before mutation and restore the
  original executable after target verification failure.
- Bun-launcher failure returns a recovery command pinned to the resolved target
  instead of losing recovery at the generic CLI error boundary.
- Installer tests enforce unique temporary paths, exact target versions,
  canonical verification, rollback language/behavior, and separate stop
  commands.

## Superseded Compiled Binary Smoke

The compiled `runx-windows-x64-baseline.exe` reported version `0.2.7`.
`upgrade --dry-run` queried live GitHub release metadata and printed the full
plan before any mutation, followed by the pinned `0.2.7` PowerShell recovery
command and the separate safe stop command. `upgrade list --format json`
returned one valid schema-versioned document containing all eight currently
published releases from `0.2.7` through `0.2.0`, newest first, with current,
latest-stable, publication, channel, and compatible-asset data.

The first smoke command targeted `bin/runx.exe` after the complete matrix build
had replaced the local single-binary output with named assets. It failed only
because that path no longer existed; the smoke was immediately rerun against
the actual baseline matrix asset and passed.

## Superseded XDocs Gate

Using an equivalent validation-only XDocs configuration with agent file writes
disabled:

- strict metadata passed for `docs/plans`, `docs/reviews/plans`, `docs/todo`,
  `source`, `devops`, and `skills/guiho-s-runx`;
- `xdocs doctor --warnings-as-errors --format json` returned `valid: true`,
  zero errors, and zero warnings; and
- `xdocs tree` returned one connected RunX tree without duplicate subjects or
  orphans.

## Release Boundary

The branch preserves package version `0.2.7`. A future merge/release owner must
use Mirror to plan and apply any requested patch release; this implementation
does not publish or tag automatically.

## Required Revalidation

- `bun run typecheck`
- `bun test`, including the controlled installer/recovery suite on Ubuntu and Windows
- `bun run build`, `bun run binary`, and `bun run binaries`
- compiled CLI text/JSON smokes, including failed and rolled-back envelopes
- strict XDocs metadata, doctor warnings-as-errors, and tree validation

The current isolated Windows worktree could not begin that rerun because
`bun install --frozen-lockfile` was denied while opening `bun.lock` with
`EPERM`. This is an environment validation blocker, not a passing result.

## Review-Correction Checks

| Command | Result |
| --- | --- |
| Cached TypeScript 6 compiler, `-p . --noEmit` | Passed after the review corrections. |
| PowerShell parser for `devops/install.ps1` | Passed. |
| Git Bash `bash -n devops/install.sh` | Passed. |
| `git diff --check` | Passed. |
| Focused Bun tests | Blocked before test execution by `EPERM` reading the isolated worktree and `tsconfig.json`; zero tests passed in that attempt. |
