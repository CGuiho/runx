---
name: Bounded Update Worker Validation
purpose: Record reproducible evidence that RunX cannot accumulate hidden update-check processes.
description: Captures focused concurrency stress, full tests, typecheck, native smoke, release assets, source audit, XDocs validation, CI, and public release checks.
created: 2026-07-21
flags:
  - testing
tags:
  - validation
  - cli
  - performance
keywords:
  - update worker
  - CPU usage
  - process count
  - lease recovery
  - RunX
owner: runx-validation
---

# Bounded Update Worker Validation

## Summary

The local implementation gate passes. RunX coalesces simultaneous update checks
before process creation, bounds the worker network lifetime, and recovers stale
or orphaned leases without allowing an old owner to delete a successor lease.

## Scope

- `source/update-cache.ts`
- `source/update-cache.spec.ts`
- Hidden worker routing in `source/cli.ts`
- Task, reference, review, and XDocs metadata
- Patch-release readiness

## Commands Run

| Command or inspection | Result |
| --- | --- |
| `bun test source/update-cache.spec.ts --timeout 30000` | Passed: 6 tests, 35 expectations. |
| `bun test --timeout 30000` | Passed: 58 tests, 402 expectations. |
| `bun run typecheck` | Passed. |
| `bun run binary` then `bin/runx.exe --version` | Passed: native Windows executable reported `0.5.2` before the release bump. |
| `bun run verify-assets` | Passed after removing only the ignored local smoke binary: exactly 14 assets. |
| Prohibited Node-import search in non-test core source | Passed: no matches. |
| `Get-Process -Name runx` | Passed: zero resident RunX processes at the audit snapshot. |

## Stress Results

| Scenario | Observed result |
| --- | --- |
| 64 simultaneous foreground schedules | Exactly one detached worker spawn. |
| Valid cache younger than four hours | Zero worker spawns. |
| 32 simultaneous stale-lease reclaimers | Exactly one successor spawn. |
| Suspended old worker finishes after successor acquisition | Successor lease remains owned and present. |
| Never-resolving release request with a 20 ms injected test deadline | Deadline rejects and releases the lease in under one second; production deadline is 15 seconds. |
| Missing and malformed stale `lease.json` | Both recover after the 30-second lease threshold and complete one worker. |
| Cache directory cannot be created | Scheduling returns `false`; foreground remains unaffected. |
| Cache path contains spaces | Atomic directory acquisition and cleanup pass on Windows. |

## Failures Or Blockers

No local implementation blocker remains.

The first asset verification correctly failed because the local native smoke
binary added a fifteenth ignored file. That generated file was removed, and the
unchanged verifier then passed with the exact fourteen release assets.

## Release Evidence

Patch version, CI, tag workflow, public asset, and public binary evidence will
be appended during the authorized release stage.

## Readiness

Ready for the Mirror patch plan and release stage.

## References

- [Task specification](../todo/bounded-update-worker.md)
- [Implementation review](../reviews/implementation/bounded-update-worker-review.md)
- [XDocs issue #14](https://github.com/CGuiho/xdocs/issues/14)
