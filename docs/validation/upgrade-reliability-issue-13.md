---
name: RunX Upgrade Reliability Issue 13 Validation
purpose: Record current-main exact recovery and installer evidence for GitHub issue 13.
description: Captures outcome-matrix, exact target, executable installer, native text/JSON, full suite, assets, XDocs, and Git checks.
created: 2026-07-19
flags:
  - validated
tags:
  - validation
  - cli
  - reliability
keywords:
  - RunX
  - issue 13
  - exact recovery
  - installer verification
owner: runx-validation
---

# RunX Upgrade Reliability Issue 13 Validation

## Summary

All local completion gates for GitHub issue #13 passed. No validation blocker
remains.

## Commands Run

| Command or check | Result |
| --- | --- |
| `bun run typecheck` | Passed |
| focused recovery, transaction, reporting, and installer tests | Passed: 23 tests, 126 assertions |
| `bun test --timeout 30000` | Passed: 52 tests, 358 assertions |
| `bun run build` | Passed |
| `bun run binary` and `bun run binaries` | Passed: single local build and twelve native targets |
| `bun run verify-assets` | Passed: exactly fourteen unique assets |
| native Windows x64 already-current text smoke | Passed: exact 0.4.1 installer followed by separate stop command |
| native Windows x64 already-current JSON smoke | Passed: resolved 0.4.1 recovery object |
| strict XDocs metadata and doctor | Passed for source, devops, and docs: zero errors and zero warnings |
| `git diff --check` | Passed |

## Recovery Evidence

- The reporter test matrix covers `upgraded`, `up-to-date`, `dry-run`,
  `rolled-back`, and `failed`; every summary precedes the exact installer, and
  every installer precedes the separate process-stop command.
- Discovery failure retains a visibly labeled fallback-current repair command.
- Download and replacement failures retain the resolved recovery object in the
  shared upgrade envelope.
- JSON exposes `targetVersion`, `targetSource`, `installCommand`, and
  `stopProcessCommand` without mixing text into stdout.

## Installer Evidence

- Windows and POSIX recovery generators preserve full prerelease identifiers.
- The PowerShell installer functions execute exact stable and prerelease
  normalization, accept the real Bun executable when its reported version
  matches, and reject a controlled mismatch.
- Installer contract tests cover exact flags, native validation, transactional
  replacement, rollback language, and the separate stop workflow.

## Platform Boundary

No WSL distribution is installed in the current Windows environment, so the
POSIX installer was not executed locally. Its shell behavior remains covered
by repository contract tests and Ubuntu CI.

## References

- [Implementation review](../reviews/implementation/upgrade-reliability-issue-13-review.md)
- [Umbrella validation](./upgrade-reliability.md)
