---
name: RunX Automatic Agent Maintenance Plan
purpose: Define executable units for GitHub issue 11 automatic skill and AGENTS.md maintenance.
description: Sequences atomic resource reconciliation, a hidden worker, startup integration, regression coverage, documentation, review, and validation.
created: 2026-07-19
flags:
  - approved
  - implementation-ready
  - implemented
tags:
  - planning
  - cli
  - agents
keywords:
  - RunX
  - issue 11
  - background worker
  - atomic writes
  - AGENTS.md
owner: runx-plans
---

# RunX Automatic Agent Maintenance Plan

## Outcome

Every ordinary RunX invocation schedules a hidden, non-blocking,
failure-isolated worker that reconciles the bundled `guiho-s-runx` skill in
both global tool locations and one compact managed block in the nearest
applicable `AGENTS.md`. Current resources remain untouched, stale or legacy
resources converge to bundled content, and foreground text, JSON, help,
version, command execution, and exit behavior remain deterministic.

## Authority

- GitHub issue [#11](https://github.com/CGuiho/runx/issues/11) is the accepted
  behavior requirement.
- `guiho-a-0001-swe` coordinates execution.
- `guiho-s-0034-cli-engineer`, `guiho-s-0023-plan-executor`,
  `guiho-s-0015-bun`, `guiho-s-0019-typescript`, `guiho-s-0011-typebox`, and
  `guiho-s-xdocs` govern implementation.
- The approved RFC breaking boundary remains: the namespace is singular
  `runx agent`; no plural compatibility alias will be restored.

## Unit AM-01 - Atomic Reconciliation Primitives

- Goal: reconcile only missing or stale agent resources without partial writes.
- Expected files: `source/storage.ts`, `source/agents.ts`, colocated tests.
- Actions:
  1. Add same-directory temporary-file writes followed by atomic replacement.
  2. Compare both global skill targets with the bundled skill before writing.
  3. Find the nearest existing `AGENTS.md` from effective cwd; create one at
     effective cwd only when no ancestor contains it.
  4. Replace current or legacy RunX blocks without duplicating them.
  5. Preserve user-authored content outside recognized markers.
- Acceptance:
  - Missing and stale resources are repaired.
  - Current resources report no change and are not rewritten.
  - Concurrent reconciliations converge to identical content.
  - Failures do not leave a truncated destination.

## Unit AM-02 - Hidden Worker Lifecycle

- Goal: schedule maintenance without blocking or polluting foreground output.
- Expected files: new `source/agent-maintenance.ts`, `source/cli.ts`, tests.
- Dependencies: AM-01.
- Actions:
  1. Add one hidden exact worker route absent from public help/tree/docs.
  2. Pass the Citty-decoded effective cwd to the detached worker.
  3. Disable recursive scheduling in the worker process.
  4. Ignore worker stdout/stderr and isolate spawn/worker failures.
  5. Exclude explicit agent-resource and uninstall commands so explicit removal
     is not immediately reversed.
- Acceptance:
  - Plain RunX and ordinary catalog commands schedule the worker.
  - The foreground does not wait for reconciliation.
  - Text and JSON output and exit codes remain unchanged.
  - Hidden worker execution cannot recursively spawn itself.

## Unit AM-03 - Acceptance Regressions

- Goal: prove every issue criterion with temporary homes and projects.
- Expected files: `source/agent-maintenance.spec.ts`, `source/cli.spec.ts`.
- Dependencies: AM-01 and AM-02.
- Tests:
  - missing, current, and outdated skills in both global tool locations;
  - nearest `AGENTS.md` insertion and stale/legacy block replacement;
  - preservation of user-authored content;
  - two concurrent reconciliations;
  - worker spawn failure and maintenance failure isolation;
  - plain invocation scheduling and effective `--cwd` routing;
  - stable text/JSON output and hidden help tree.
- Acceptance:
  - The full suite exits cleanly with no real global mutation.

## Unit AM-04 - Documentation, Review, And Validation

- Goal: close durable state and provide issue-resolution evidence.
- Dependencies: AM-03.
- Actions:
  1. Document automatic behavior and explicit-command boundaries.
  2. Review implementation against every issue checkbox.
  3. Run typecheck, tests, build, native smoke, prohibited-import scan, and
     narrow XDocs checks.
  4. Commit one file at a time, push main, comment factual evidence on issue
     #11, and close only after GitHub contains the implementation.
- Acceptance:
  - No review finding remains.
  - Pushed commit and validation evidence precede closure.

## Release Boundary

This externally observable capability requires a final Mirror `minor` plan
after all open issues are resolved. Do not version or tag during this unit.

## References

- [Task specification](../todo/automatic-agent-maintenance.md)
- [CLI architecture](../architecture/cli-architecture.md)
- [RFC migration plan](./rfc-0034-cli-compliance-migration.md)
