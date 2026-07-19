---
name: RunX Upgrade Reliability Issue 13 Implementation Review
purpose: Review exact-version recovery and installer behavior against GitHub issue 13.
description: Findings-first review of recovery coverage, target pinning, separate stop commands, installer version verification, JSON, and delivery readiness.
created: 2026-07-19
flags:
  - accepted
tags:
  - review
  - cli
  - reliability
keywords:
  - RunX
  - issue 13
  - recovery install
  - exact version
owner: runx-implementation-reviews
---

# RunX Upgrade Reliability Issue 13 Implementation Review

## Verdict

Accepted.

## Findings

No blocker, high, medium, or low implementation finding remains.

The implementation already generated recovery from the resolved plan before
mutation and retained it through failure envelopes. This audit strengthened
the release gate with an explicit matrix for all five terminal outcomes and an
executable PowerShell installer test for stable/prerelease normalization,
successful `--version` verification, and mismatch rejection.

## Acceptance Criteria Check

- Recovery after success, failure, already current, dry run, and rollback:
  accepted.
- Exact stable and prerelease target preservation: accepted.
- Copyable installer command generated from supported installer flags:
  accepted.
- Separate Windows and POSIX process-stop commands: accepted.
- Fallback-current repair when discovery cannot resolve a target: accepted.
- Text ordering and JSON recovery parity: accepted.
- Direct installer version arguments, canonical verification, and nonzero
  mismatch behavior: accepted.
- Download, replacement, verification, and already-current regressions:
  accepted.

## Platform Note

The current Windows gate executes PowerShell normalization and version
verification directly. Git for Windows Bash executes the Linux/macOS installer
syntax, piped startup, exact-version normalization, and executable verification
locally; Ubuntu CI remains the native Linux gate.

## References

- [Approved design](../../superpowers/specs/2026-07-15-upgrade-reliability-design.md)
- [Implementation plan](../../plans/upgrade-reliability-implementation.md)
- [Validation](../../validation/upgrade-reliability-issue-13.md)
- [GitHub issue #13](https://github.com/CGuiho/runx/issues/13)
