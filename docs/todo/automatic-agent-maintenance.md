---
name: Add Automatic RunX Agent Maintenance
purpose: Define the outcome and acceptance signals for RunX GitHub issue 11.
description: Tracks non-blocking global skill repair, nearest AGENTS.md reconciliation, atomicity, isolation, tests, and delivery evidence.
created: 2026-07-19
flags:
  - approved
  - completed
tags:
  - todo
  - cli
  - agents
keywords:
  - RunX
  - GitHub issue 11
  - guiho-s-runx
  - AGENTS.md
owner: runx-todo
---

# Add Automatic RunX Agent Maintenance

## Todo Index

- Task: `2. Add Automatic RunX Agent Maintenance`
- Status: completed
- Index: [TODO.md](../../TODO.md)
- External: [CGuiho/runx issue #11](https://github.com/CGuiho/runx/issues/11)

## Outcome

Ordinary RunX commands automatically and silently converge the bundled global
skill and nearest project instruction block without delaying, polluting, or
failing the requested command.

## Scope

### In scope

- Both global RunX skill targets.
- Nearest applicable `AGENTS.md`.
- Current and legacy managed blocks.
- Atomic, idempotent, concurrent-safe writes.
- Hidden detached worker startup and failure isolation.
- Text/JSON/help/version/exit-code stability.
- Tests, docs, xdocs, review, validation, push, comment, and closure.

### Out of scope

- Reintroducing plural `runx agents` compatibility.
- Network-based skill discovery.
- Implicit XDocs or other CLI resource mutations.
- Versioning before the complete open-issue sequence finishes.

## Acceptance Signals

- Plain RunX schedules the hidden worker.
- Missing/stale global skills are repaired; current skills are not rewritten.
- The nearest `AGENTS.md` contains exactly one compact current RunX block.
- User-authored content remains unchanged outside managed markers.
- Worker scheduling and reconciliation failures cannot fail the foreground.
- Concurrent workers converge safely.
- Public help contains no hidden worker route.
- Text and JSON output remain deterministic.
- Tests cover every state named by issue #11.
- Pushed evidence precedes the resolution comment and closure.

## Watch-outs

- Explicit skill uninstall and RunX uninstall must not trigger reinstallation.
- Tests must use temporary HOME/USERPROFILE paths and never mutate real global
  skills or repository instructions.
- The worker must use Citty-decoded `--cwd`, not a second public flag parser.

## Related Files

- [Implementation plan](../plans/automatic-agent-maintenance.md)
- [Implementation notes](./automatic-agent-maintenance-implementation.md)
- [Implementation review](../reviews/implementation/automatic-agent-maintenance-review.md)
- [Validation](../validation/automatic-agent-maintenance.md)
- [CLI architecture](../architecture/cli-architecture.md)

## References

- [TODO.md](../../TODO.md)
- [GitHub issue #11](https://github.com/CGuiho/runx/issues/11)
