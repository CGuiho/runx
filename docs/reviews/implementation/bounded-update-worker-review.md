---
name: Bounded Update Worker Implementation Review
purpose: Review RunX against the XDocs CPU incident and the task 8 acceptance signals.
description: Assesses hidden-worker routing, cache freshness, process coalescing, network deadlines, ownership safety, stale recovery, tests, and documentation.
created: 2026-07-21
flags:
  - accepted
tags:
  - implementation-review
  - cli
  - performance
keywords:
  - update worker
  - CPU saturation
  - process accumulation
  - cache lease
  - RunX
owner: runx-implementation-reviews
---

# Bounded Update Worker Implementation Review

## Verdict

Accepted. No blocker, high, medium, or low implementation finding remains.

## Findings

No open findings.

## Acceptance Criteria Check

| Criterion | Result | Evidence |
| --- | --- | --- |
| Hidden worker cannot recursively schedule itself | Accepted | `runCli` handles the exact internal flag and returns before cache reads, spawning, Citty routing, or agent maintenance. |
| Foreground never waits for update network activity | Accepted | Foreground work is limited to cache and lease files plus detached process creation. |
| Repeated invocations cannot accumulate workers | Accepted | A serialized, atomic lease is acquired before spawn; the 64-way stress test observes one spawn. |
| Network lifetime is bounded | Accepted | Worker fetch receives an abort signal and a hard 15-second race deadline. |
| Stale recovery is ownership-safe | Accepted | Lease removal checks the UUID token while holding the mutation guard; the suspended-worker and 32-reclaimer stress test passes. |
| Interrupted lease creation recovers | Accepted | Missing or malformed metadata is reclaimed only after the 30-second directory lease age. |
| Foreground failures remain isolated | Accepted | Scheduling catches cache, directory, lease, and spawn failures and returns `false`. |
| Cache behavior remains compatible | Accepted | TypeBox decoding and atomic cache replacement remain; a valid result suppresses checks for four hours. |

## Verification Evidence

- `bun test --timeout 30000`: 58 tests passed with 402 expectations.
- `bun run typecheck`: passed.
- Native `bin/runx.exe --version` smoke: `0.5.2` before the release bump.
- `bun run verify-assets`: exactly fourteen RFC 0034 assets.
- Core-source prohibited Node import search: no matches.

## Docs And TODO Check

Task 8, its task spec, the canonical runtime documentation, source descriptor,
implementation review, and validation record cover the changed behavior.

## Residual Risk

No code-level residual risk is open. CI and public release verification remain
release-stage evidence rather than implementation findings.

## References

- [Task specification](../../todo/bounded-update-worker.md)
- [Validation report](../../validation/bounded-update-worker.md)
- [XDocs issue #14](https://github.com/CGuiho/xdocs/issues/14)
