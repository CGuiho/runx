---
name: Bound The RunX Update Worker
purpose: Define the CPU-safety outcome and completion signals for task 8 in TODO.md.
description: Requires RunX update checks to remain nonblocking, coalesced, time-bounded, ownership-safe, and self-recovering after interrupted workers.
created: 2026-07-21
flags:
  - testing
tags:
  - cli
  - reliability
  - performance
keywords:
  - update worker
  - CPU usage
  - process accumulation
  - cache lease
  - XDocs issue 14
owner: runx-todo
---

# Bound The RunX Update Worker

## Summary

The XDocs CPU incident showed that detached update checks can accumulate until
they saturate a machine. RunX must preserve its nonblocking cached update notice
without allowing repeated foreground commands, stalled network requests, or
suspended workers to create an unbounded process population.

## Todo Index

- Task: `8. Bound The RunX Update Worker`
- Status: testing
- Index: [TODO.md](../../TODO.md)

## Outcome

Ordinary RunX commands use a fresh four-hour cache when available and otherwise
schedule at most one finite hidden update check per global RunX cache directory.

## Scope

### In scope

- Hidden update-worker routing and recursion prevention.
- Four-hour cache freshness.
- Atomic, ownership-safe coalescing before process creation.
- Fifteen-second remote-check deadline.
- Thirty-second stale and orphaned lease recovery.
- Stress coverage for concurrent schedules, reclaimers, suspended workers, and
  malformed lease metadata.

### Out of scope

- RunX agent-maintenance behavior, which is a separate finite local worker.
- Unrelated RunX GitHub issues #22 and #23.
- Changes to release discovery, upgrade selection, or installer behavior.

## Acceptance Signals

- Sixty-four concurrent foreground scheduling attempts create exactly one
  detached update worker.
- Thirty-two concurrent stale-lease reclaimers create exactly one successor.
- A suspended old worker cannot remove the successor's lease.
- Missing or malformed lease metadata is reclaimed only after 30 seconds.
- A stalled release request ends at the 15-second production deadline and
  releases its lease.
- Fresh cache data suppresses worker creation for four hours.
- Worker scheduling failures never alter foreground output or exit behavior.
- The full RunX test suite, TypeScript typecheck, native build, release-asset
  verification, and XDocs checks pass.

## Watch-outs

- Lease removal must verify ownership while holding the mutation guard.
- The hidden worker must return before ordinary command routing and must never
  schedule another worker.
- Cache and lease files stay under `~/.guiho/runx/`.

## References

- [XDocs issue #14](https://github.com/CGuiho/xdocs/issues/14)
- [TODO.md](../../TODO.md)
