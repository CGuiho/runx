---
name: RunX Automatic Agent Maintenance Implementation
purpose: Preserve the implemented GitHub issue 11 automatic agent-resource maintenance record.
description: Records the atomic reconciler, hidden detached worker, command boundaries, regression coverage, and delivery handoff.
created: 2026-07-19
flags:
  - completed
tags:
  - implementation
  - cli
  - agents
keywords:
  - RunX
  - issue 11
  - guiho-s-runx
  - AGENTS.md
owner: runx-todo
---

# RunX Automatic Agent Maintenance Implementation

## Summary

Ordinary RunX invocations now schedule a silent detached worker that compares
the bundled `guiho-s-runx` skill with both global tool installations and
reconciles one nearest applicable `AGENTS.md`. Foreground command output, JSON,
help, delegated exit behavior, and startup latency remain independent of the
maintenance result.

## Implementation Map

- `source/storage.ts` provides same-directory temporary writes followed by
  replacement and cleanup on failure.
- `source/agents.ts` compares current skill bytes, repairs missing or stale
  global skills, finds the nearest ancestor `AGENTS.md`, migrates legacy
  markers, preserves user content, and emits exactly one current managed block.
- `source/agent-maintenance.ts` TypeBox-decodes the hidden worker invocation,
  disables recursion, detaches standard streams, and isolates spawn failures.
- `source/cli.ts` schedules maintenance after Citty execution with its decoded
  effective cwd while excluding explicit `agent` and `uninstall` operations.
- `source/agent-maintenance.spec.ts` and `source/cli.spec.ts` cover missing,
  current, stale, legacy, concurrent, hidden, failure, output, cwd, and explicit
  removal behavior without touching real user locations.

## Delivery Boundary

This implementation completes GitHub issue #11. Versioning remains deferred to
the single final Mirror minor release after the complete open-issue sequence.

## References

- [Task specification](./automatic-agent-maintenance.md)
- [Implementation plan](../plans/automatic-agent-maintenance.md)
- [Implementation review](../reviews/implementation/automatic-agent-maintenance-review.md)
- [Validation](../validation/automatic-agent-maintenance.md)
