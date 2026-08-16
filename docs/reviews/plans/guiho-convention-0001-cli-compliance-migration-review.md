---
name: RunX GUIHO CLI Convention 0001 Compliance Migration Plan Review
purpose: Independently verify that the complete breaking migration plan is safe, question-sealed, branch-aware, and executable after its explicit human gates.
description: Reviews dormant sequencing, ownership, command catalogs, validation, cutover availability, release integration, durable evidence, and final hardening.
created: 2026-08-16
flags:
  - ready-for-approval-execution
  - blocked-pending-human-approval
tags:
  - review
  - plan
  - cli
keywords:
  - GUIHO CLI Convention 0001
  - stable launcher
  - protocol-v1 cutover
  - release evidence
  - exact-head validation
owner: runx-plan-reviews
---

# RunX GUIHO CLI Convention 0001 Compliance Migration Plan Review

## Verdict

Ready for approval/execution.

No blocker or high-severity planning finding remains after four independent
review and revision passes. This is a plan-quality verdict, not implementation,
cross-repository, integration, release, publication, or production authority.
The identity, current-base, Superiority, execution, and release gates in the
plan remain mandatory.

## Reviewed Sources

- [Compliance migration plan](../../plans/guiho-convention-0001-cli-compliance-migration.md)
- [Detailed migration task](../../todo/guiho-convention-0001-cli-compliance-migration.md)
- [Accepted architecture](../../architecture/guiho-convention-0001-cli-compliance-migration.md)
- [Architecture review](../architecture/guiho-convention-0001-cli-compliance-migration-review.md)
- [Compliance audit](../guiho-convention-0001-cli-compliance-audit.md)
- GUIHO CLI Convention 0001
- Current RunX repository instructions and canonical GUIHO SWE lifecycle

## Resolved Findings

### Resolved — Pre-cutover units changed public behavior

U02, U03, and U11 now build private or unregistered implementations only. U03
keeps the active legacy embed/import/maintenance/CI graph byte-for-byte intact.
C00 alone owns the active init, configuration, agent, help, version, resource,
maintenance, npm, workflow, and lifecycle switch.

### Resolved — Cross-repository prerequisite was circular

E00 now depends on identity approval, explicit Superiority authorization, and a
recorded current Superiority protected-main SHA. Its accepted integration—or an
explicit human governance exception—satisfies Gate D; Gate D is no longer an
E00 prerequisite.

### Resolved — Complete init was assigned before its dependencies

U02 owns strict configuration, policy, schemas/examples, and a dormant injected
orchestrator. U03 injects the private canonical resources. C00 activates the
complete public init only after resource, ownership, launcher, and lifecycle
primitives are integrated.

### Resolved — Unattended execution was not sealed

The shared lifecycle now states exactly what one human execution approval
covers. Each executor records exact bases rather than asking for per-unit
reapproval, may use only non-material reversible provisional answers, and must
reject/requeue a materially contradictory unit. Every implementation PR has an
exact-head 0049 review, unchanged-head 0050 validation, and 0052 integration,
reachability, and cleanup.

### Resolved — Release operation and evidence ownership were incomplete

R00 assigns Mirror apply, release commit, push, protected tag, GitHub Release,
workflow observation, reachability, and cleanup to 0052. 0050 validates the
exact source, release commit/tag, and remote result. A separate evidence-only
branch and PR persists the question, release, validation, URL, checksum, runner,
TODO, and XDocs records bound to `RELEASE_COMMIT_SHA` before H00 may start.

### Resolved — Active cutover ownership was implicit

C00 enumerates the active Cobra, config resolution, embedded resources,
maintenance, lifecycle, installer/uninstaller, release tooling, npm, workflow,
README/DOCS, skill/prompt/instruction, catalog, and descriptor paths it owns.
It excludes new algorithm work and publication before merge.

### Resolved — Validation and catalog commands were prose or unowned

The plan expands exact `RG` and `RG+EMBED` gates and names every future helper,
black-box package, script parse, build, verification, source-contract, Mirror,
and clean-clone command. U02-U11 each own only their stable-UID, nonmutating
`runx.yaml` entries, with `check`/`list` parity and no cross-unit selector edits.

### Resolved — Version/release decisions did not fit unit scope

E00 defers only to a separately authorized Superiority release; U00 is
planning-only; U01-U11 and C00 consolidate into R00; R00 records exact release
state; and H00 owns its patch-plan decision without deferring backward. The
accepted main-branch installer-only hardening can require no second release
because the already accepted protocol-v1 assets remain unchanged.

### Resolved — Audited edge conditions could disappear inside broad units

The final plan explicitly covers installer short-alias removal, rejection of
arbitrary install directories, full init outcome groups and domain questions,
npm executable/shim/workflow/cache retirement, Windows rename/delete-on-close,
two-phase prepare/activate/finalize/abort, recovery output, and clause-level
audit traceability.

## Readiness Guarantees

The reviewed plan now supplies, for every E/U/C/H unit:

- exact dependency and protected-main base rules;
- branch, isolated worktree, question ledger, owned paths, and exclusions;
- data/config/auth/cache/docs impacts;
- exact actions and validation commands;
- acceptance and stop conditions;
- smallest-coherent commits, plain push, and a focused PR;
- exact-head independent review and validation;
- gated integration, reachability, and cleanup; and
- a scope-correct Mirror, GitHub Release, and production decision.

The C00 to separately authorized R00 to evidence PR to H00 sequence preserves
the canonical main-branch installer throughout the first protocol-v1 release
window and keeps publication authority separate from implementation approval.

## Remaining Human Gates

Before any execution, the human must confirm or replace the proposed CLI home,
skill, prompt, instruction, repository/issue URLs, launcher protocol, and
resource contract; approve the architecture and complete breaking plan; and
authorize E00 against an exact current Superiority base. RunX U00 also requires
a current exact RunX protected-main base and preservation of the dirty local
TODO work.

R00 still requires a separate approval naming the exact Mirror target, release
commit/push, protected tag, GitHub Release, assets, workflow/production boundary,
and remote smoke. Nothing in this review grants those actions.

## First Executable Units

- First cross-repository unit: E00, after Gate A, separate Superiority
  authorization, and a recorded current Superiority protected-main SHA.
- First RunX unit: U00, after Gate D and a recorded current RunX protected-main
  base.

## Handoff

```yaml
handoff:
  from: guiho-a-0047-plan-reviewer
  to: guiho-a-0001-swe
  verdict: ready-for-approval-execution
  first_unit: E00
  next_agent_after_human_execution_approval: guiho-a-0048-plan-executor
  implementation_authorized: false
  release_authorized: false
  production_authorized: false
```

## References

- [Compliance migration plan](../../plans/guiho-convention-0001-cli-compliance-migration.md)
- [Detailed migration task](../../todo/guiho-convention-0001-cli-compliance-migration.md)
- [Accepted architecture](../../architecture/guiho-convention-0001-cli-compliance-migration.md)
- [Architecture review](../architecture/guiho-convention-0001-cli-compliance-migration-review.md)
- [Compliance audit](../guiho-convention-0001-cli-compliance-audit.md)
