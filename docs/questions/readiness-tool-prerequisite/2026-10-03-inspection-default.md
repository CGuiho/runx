---
name: Default RunX Inspection Scheduling Decision
purpose: Record the reversible safe boundary selected for the authorized runtime prerequisite.
description: Source-grounded choice to bypass lifecycle scheduling for inspection and parsed dry runs while retaining real execution and explicit mutations.
created: 2026-10-03
flags: [human-review-pending]
tags: [questions, runx, readiness]
keywords: [scheduling, issue-64, safe-boundary]
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# Default RunX Inspection Scheduling Decision

- Human review: pending
- Task: [RunX #64](https://github.com/CGuiho/runx/issues/64)
- Mode: dnd

**Question:** use a new opt-out interface or make the existing inspection
commands safe by default?

**Evidence:** `cmd/root.go` schedules both lifecycle workers from pre-run before
catalog operations. The named maintenance environment toggle guards only bare
foreground bootstrap; it cannot prove no scheduling. Cobra parses RunX run flags
only before the selector (`SetInterspersed(false)`), so inspecting the actual
parsed Boolean distinguishes a real dry run from a child argument or false flag.

**Decision:** return from `scheduleLifecycle` before executable resolution for
check/list/describe/reveal and parsed dry-run=true. No new schema, flag, global
setting or user confirmation policy. Real run and explicit setup/agent operations
retain their existing behavior. High confidence; reversing this small scheduler
gate would restore the old contract, but would also restore the readiness gap.

**Authority:** parent expressly authorized this minimum owning-source fix and a
separate runtime issue. Planning-phase waiver is bounded to that outcome. Commit
and independent-review evidence are required; no global binary replacement or
push by this leaf. Existing XDocs grants are empty: do not broaden them or invoke
unsafe tooling implicitly. Descriptor synchronization remains a visible separate
gap pending named authorization, rather than silently granting the entire root.
