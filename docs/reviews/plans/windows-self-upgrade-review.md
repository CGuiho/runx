---
name: RunX Windows Self-Upgrade Plan Review
purpose: Verify that the Windows self-upgrade plan is safe, complete, and executable.
description: Reviews replacement sequencing, rollback, tests, documentation, issue closure, and Mirror patch delivery.
created: 2026-07-14
flags:
  - approved
tags:
  - reviews
  - plans
  - windows
keywords:
  - runx upgrade
  - executable replacement
  - rollback
  - mirror patch
owner: runx-plan-reviews
---

# RunX Windows Self-Upgrade Plan Review

## Verdict

Ready for execution.

## Findings

No blocker, high, medium, or low findings remain. The plan traces to the accepted
Windows self-upgrade decision and existing native-upgrade requirement. Each unit
has a bounded outcome, named files, dependencies, acceptance criteria, checks,
and explicit stop conditions.

## Sequencing and Failure Safety

- Replacement precedes verification, and successful completion requires the
  canonical path to execute the expected target version.
- Rollback occurs before any failure is surfaced after the old executable was
  renamed.
- Deferred work is limited to deleting the locked old image; replacement itself
  is complete before success.
- Issue closure and Mirror application follow implementation validation and a
  clean committed worktree.
- The plan recognizes configured Mirror commit, push, and tag side effects and
  stops on protected-push rejection.

## Coverage Review

- Success, failure, rollback, temporary-file cleanup, and platform isolation are
  explicit.
- Existing `-h` coverage directly supports closing issue #1.
- Typecheck, unit tests, library build, local binary, the twelve-target matrix,
  strict XDocs checks, and diff hygiene are included.
- No route, API, UI, data schema, authentication, authorization, cache, cloud,
  secret, or database work is in scope.

## TODO and Documentation Alignment

The work is a single focused package session and does not need a long-running
TODO entry. Public `DOCS.md`, source/root XDocs descriptors, lifecycle records,
and the changelog are included in the correct sequence.

## First Executable Unit

Implement synchronous replacement, target-version verification, rollback, and
old-image cleanup in `source/self-management.ts`.

## References

- [Plan](../../plans/windows-self-upgrade.md)
- [Decision](../../decisions/windows-self-upgrade.md)
- [CLI Architecture](../../architecture/cli-architecture.md)
