---
name: Bounded Update Worker Validation
purpose: Record reproducible evidence that RunX cannot accumulate hidden update-check processes.
description: Captures focused concurrency stress, full tests, typecheck, native smoke, release assets, source audit, XDocs validation, CI, and public release checks.
created: 2026-07-21
flags:
  - validated
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
| `bun run binary` then `bin/runx.exe --version` | Passed before the release bump. The public Windows baseline binary later reported `0.5.3`. |
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

A post-release local full-suite rerun encountered two timing-only failures in
the pre-existing help and agent-maintenance tests while the shared Windows host
was heavily loaded. The update-worker suite still passed all 35 expectations.
GitHub CI attempt 2 then passed both Linux and Windows jobs on clean runners.

## Release Evidence

- Mirror applied patch version `0.5.3` and pushed commit
  `c7dfdbf75fbbdb5d4b9c317f386c7f48ac87d151` plus tag
  `@guiho/runx@0.5.3`.
- [GitHub CI run 29859187581](https://github.com/CGuiho/runx/actions/runs/29859187581)
  passed on attempt 2 for Linux and Windows. The first Linux public-installer
  probe ran before the concurrently triggered release existed; rerunning after
  publication passed without a code change.
- [Publish run 29859191198](https://github.com/CGuiho/runx/actions/runs/29859191198)
  completed successfully.
- The public release contains exactly fourteen assets: twelve native binaries
  plus `guiho-i-runx.md` and `guiho-s-runx.md`.
- The public release body contains only the `0.5.3` changelog section.
- The downloaded 97,947,648-byte Windows x64 baseline asset matched published
  SHA-256 `edc3f15fc66d5785dd9fdbfa1767f5d108e2a9f8836154d8c293f117d1a6fe86`
  and reported version `0.5.3`.
- The public-binary probe's finite detached maintenance process exited, leaving
  zero RunX processes and allowing deterministic cleanup of the downloaded
  executable.

## Readiness

Validated and publicly released as RunX `0.5.3`.

## References

- [Task specification](../todo/bounded-update-worker.md)
- [Implementation review](../reviews/implementation/bounded-update-worker-review.md)
- [XDocs issue #14](https://github.com/CGuiho/xdocs/issues/14)
