---
name: RunX Manifest V2 Composition Implementation Review
purpose: Review issue 26 for correctness, safety, coverage, and release readiness.
description: Audits recursive decoding, graph identity, parent-child reciprocity, foreign loading, execution bases, CLI integration, and migration documentation.
created: 2026-07-23
flags:
  - approved
tags:
  - reviews
  - implementation
keywords:
  - manifest v2
  - issue 26
owner: runx-implementation-reviews
---

# RunX Manifest V2 Composition Implementation Review

## Verdict

Approved for a minor release. No blocking correctness, execution-safety, or
distribution defect remains in the reviewed scope. Issue 22 is absent.

## Findings Closed During Review

- Added end-to-end CLI coverage proving check, list, describe, and dry-run all
  consume the same composed graph rather than configuration-only helpers.
- Added a real GitHub blob-to-raw mock path, foreign provenance, upstream parent
  reciprocity, local mount cwd behavior, and transport-to-exit-3 coverage.
- Wrapped foreign timeout/transport failures as configuration errors and
  revalidated final redirected URLs against the GitHub-only allowlist.
- Confirmed local child commands resolve cwd from their own file and foreign
  commands cannot synthesize a remote filesystem path.
- Rejected cross-command UID, selector, and nested-ID shorthand ownership
  collisions while preserving multiple aliases that resolve to the same leaf.
- Enforced child namespace identity before mount aliasing, complete parent graph
  semantics, exact 32/33 depth behavior, foreign root confinement, and safe
  script/cwd paths during composition rather than execution.
- Replaced whole-body foreign reads with a cancelable streamed one-megabyte
  limit and verified cancellation after the first over-limit chunk.
- Normalized generated namespaces, derived nested init namespaces from the
  target directory, and marked all superseded v1 records as historical.
- Corrected release acceptance so mutable latest smoke belongs to ordinary CI
  and the immutable exact-version installer belongs after publication.

## Contract Review

The recursive TypeBox source model and separate effective runtime model keep
legacy keys out while preserving one selector/execution pipeline. Sibling,
namespace, UID, selector, cwd, parent, cycle, depth, URL, timeout, and size
checks happen before execution. Empty groups are valid only when expressed as
`commands: []`; groups with both or neither child form fail closed.

## Residual Operational Risk

A GitHub URL that names a mutable branch intentionally follows that branch.
Because catalogs are executable code, automation should prefer commit-pinned
GitHub URLs and always inspect foreign entries with check/list/describe/dry-run.
RunX does not persist foreign content, so it cannot silently reuse a stale copy.

## Final Release Evidence

- Independent focused audit: 32 tests and 221 assertions on exact pushed HEAD
  `376404c`, with all eleven blocker categories cleared.
- Full suite: 90 tests and 571 assertions.
- RunX 0.7.2 CI and publish passed, including exact fourteen assets, the
  tag-pinned public installer, npm trusted publication, and version-only notes.
- Both Windows process tests use explicit fifteen-second budgets and pass in the
  full Windows CI suite without changing runtime behavior.
- GitHub issue 26 closed as completed; issue 22 remained excluded and untouched.
