---
name: RunX Upgrade Reliability Issue 12 Implementation Review
purpose: Review current-main self-upgrade and complete release listing behavior against GitHub issue 12.
description: Findings-first review of synchronous replacement, rollback, output ordering, complete catalog routing, JSON, tests, and delivery readiness.
created: 2026-07-19
flags:
  - accepted
tags:
  - review
  - cli
  - reliability
keywords:
  - RunX
  - issue 12
  - upgrade list
  - transactional replacement
owner: runx-implementation-reviews
---

# RunX Upgrade Reliability Issue 12 Implementation Review

## Verdict

Accepted after correction.

## Findings

Two public-list findings were corrected before acceptance:

- the default view truncated the complete release catalog to twenty entries and
  hid prereleases unless `--pre-releases` was supplied;
- Citty executed the selected `upgrade list` leaf and then also executed the
  parent `upgrade` action, appending a second JSON document and failure result.

The default list now returns every normalized published release, including
labeled prereleases. Pagination occurs only when the user explicitly supplies
`--page` or `--per-page`. The parent action checks the resolved Citty command
identity and does not execute after `upgrade list` or `upgrade check`.

No blocker, high, medium, or low finding remains.

## Acceptance Criteria Check

- Immediate canonical replacement before success: accepted.
- Exact canonical `--version` verification: accepted.
- Backup rollback on replacement or verification failure: accepted.
- Mapped Windows executable replacement and deferred cleanup only: accepted.
- Visible plan and ordered download/validate/replace/verify phases: accepted.
- Every GitHub release page, SemVer order, channels, dates, current/latest, and
  missing compatible assets: accepted.
- Complete prerelease-inclusive default list: accepted after correction.
- One parseable JSON document with equivalent metadata: accepted after
  correction.
- Windows swap, permission/rename failure, rollback, version mismatch,
  pagination, labels, and missing-asset regressions: accepted.

## References

- [Approved design](../../superpowers/specs/2026-07-15-upgrade-reliability-design.md)
- [Implementation plan](../../plans/upgrade-reliability-implementation.md)
- [Validation](../../validation/upgrade-reliability-issue-12.md)
- [GitHub issue #12](https://github.com/CGuiho/runx/issues/12)
