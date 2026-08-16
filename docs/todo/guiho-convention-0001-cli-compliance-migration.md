---
name: Make RunX Comply With GUIHO CLI Convention 0001
purpose: Track the complete breaking migration from the audited legacy lifecycle to the accepted convention-compliant architecture.
description: Durable scope, dependency, status, acceptance, and evidence ledger for every planning, implementation, cutover, release, and hardening unit.
created: 2026-08-16
flags:
  - proposed
  - breaking-change
tags:
  - todo
  - cli
  - compliance
keywords:
  - GUIHO CLI Convention 0001
  - stable launcher
  - protocol-v1
  - lifecycle migration
owner: runx-todo
---

# Make RunX Comply With GUIHO CLI Convention 0001

## Summary

Replace RunX's audited legacy direct-binary lifecycle with the accepted stable
launcher, immutable payload, artifact ownership, two-phase transaction,
global/project configuration, agent-evolution, and complete release design.

This task is a breaking pre-1.0 migration. It is not authorized for execution
until the identities, architecture, plan, current base, and cross-repository
skill alignment gates in the linked plan are approved.

The independent architecture review is Ready for planning and the independent
plan review is Ready for approval/execution. Both remain unintegrated planning
artifacts in this dirty, stale checkout and do not satisfy their human or
current-base gates by themselves.

## Todo Index

- Root task: `18. Make RunX Comply With GUIHO CLI Convention 0001`
- Status: todo
- Index: [TODO.md](../../TODO.md)

## Required Outcome

- All 24 audit finding groups have direct evidence on protected `main`.
- The stable launcher is the only command entrypoint.
- Protocol-v1 install/reinstall/repair/uninstall/upgrade transactions are safe,
  synchronous, and manifest-owned.
- Config, schemas, examples, init, agent policy/resources, help, version,
  release, docs, CI, Mirror, RunX, and XDocs comply.
- The temporary dual-shape installer is removed after the accepted first
  protocol-v1 release.
- No unrelated user work is overwritten or absorbed.

## Unit Ledger

| Unit | Repository | Status | Dependency | Deliverable |
| --- | --- | --- | --- | --- |
| E00 | Superiority | pending | human cross-repo authority | Canonical Go CLI skill aligned |
| U00 | RunX | pending | E00/accepted exception; approvals | Planning baseline on current main |
| U01 | RunX | pending | U00 | Root Mirror/RunX/XDocs tooling |
| U02 | RunX | pending | U01 | Config, policy, schemas, examples, init |
| U03 | RunX | pending | U02 | Canonical agent resources/commands |
| U04 | RunX | pending | U03 | Manifest, checksums, SemVer, selection |
| U05 | RunX | pending | U04 | Ownership state and two-phase lifecycle |
| U06 | RunX | pending | U05 | Stable launcher protocol |
| U07 | RunX | pending | U06 | Dormant protocol-v1 release toolchain |
| U08 | RunX | pending | U07 | Dormant next installers |
| U09 | RunX | pending | U08 | Dormant next uninstall |
| U10 | RunX | pending | U09 | Dormant next upgrade |
| U11 | RunX | pending | U10 | Help and raw version |
| C00 | RunX | pending | U11; release window | Transition cutover PR |
| R00 | RunX/external | awaiting authorization | C00 integrated | First protocol-v1 release plus evidence-only PR |
| H00 | RunX | pending | accepted R00 | Remove transition and prove compliance |

Allowed implementation states are:

```text
pending -> in progress -> testing -> implementation review
        -> validation -> integrated
```

R00 additionally uses `awaiting authorization` and `published-awaiting-smoke`.
Its Mirror/tag/GitHub Release operation is owned by the Pull Request Integrator;
after remote validation, a separate docs-only branch and PR persist the release,
validation, question-ledger, and TODO evidence bound to the exact release commit
and URLs. H00 cannot begin until that evidence PR is integrated and reachable.

## Approval Checklist

- [ ] Human confirms CLI home `runx`.
- [ ] Human confirms main skill `guiho-s-runx`.
- [ ] Human confirms main prompt `guiho-p-runx`.
- [ ] Human confirms instruction `guiho-i-runx.md`.
- [ ] Human confirms repository and issue URLs.
- [ ] Human approves the architecture.
- [ ] Independent architecture review is Ready for planning.
- [ ] Independent plan review is Ready for approval/execution.
- [ ] Human approves the full plan and breaking cutover.
- [ ] Current protected-main base and planning baseline are recorded.
- [ ] E00 is integrated or its conflict is explicitly accepted.

The ready architecture-review item may be checked only when the planning
baseline containing that review is approved; this document does not self-approve.

## Scope

### In scope

- Mandatory repository tooling and version authority.
- Strict config, policy, schemas, examples, and init.
- Canonical agent skill/prompt/definition/instruction and command tree.
- Complete release integrity, SemVer, channels, and target validation.
- Stable launcher, current pointer, installed ledger, locks, instances, and
  transaction journals.
- PowerShell/POSIX installers/uninstallers, Cobra uninstall, and whole-release
  upgrade.
- Legacy shadow migration and npm distribution retirement.
- Help/version, docs, CI, workflows, XDocs, validation, first v1 release, and
  final hardening.

### Out of scope

- Secret or environment-value management.
- Unrelated RunX domain changes.
- Signatures/attestations not required by Convention 0001.
- Real user-home/PATH/global-agent mutation during implementation validation.
- Any publication, production, deployment, or issue action without separate
  explicit authorization.

## Persistent And Disposable State

- Preserve during reinstall/upgrade: global/project config, declared data,
  databases, and user-created persistent content.
- Replace/retire by manifest: launcher, payloads, canonical resources,
  projections, instruction block, schemas, examples, metadata.
- Dispose: caches, stale instances, completed journals, owned staging, and
  legacy npm cache.
- Always preserve: shared `.guiho`/`bin`/`.temp` roots, shared PATH entry,
  foreign files, and non-RunX `AGENTS.md` bytes.

## Completion Evidence

Each integrated unit must record:

- approved base, branch, worktree, and PR;
- final commit and exact reviewed/validated head;
- 0049 verdict;
- 0050 validation report;
- merge reachability and cleanup;
- tests passed/failed/skipped;
- question ledger;
- affected audit findings; and
- Mirror/release/production status.

E00 defers version/distribution only to a separately authorized Superiority
release. U01-U11 and C00 consolidate their RunX release effect into R00. H00
runs and records its own patch plan: main-branch installer compatibility removal
plus docs/CI/validation may require no new release when versioned protocol-v1
assets are unchanged; changes to binaries, embedded resources, or versioned
release assets require a separately authorized patch release.

The overall task remains open when any critical destructive test is skipped,
any audit finding is residual, the transition installer remains, review and
validation reference different SHAs, or the integrated-main rerun is missing.

## Stop Conditions

Stop the active unit for stale/unapproved base, unexplained overlapping dirty
work, unresolved identity/architecture question, unsafe ownership target,
unverifiable process identity, incomplete rollback, real-home dependency,
publication outside R00, or any requirement to bypass protected-main gates.

## References

- [Implementation plan](../plans/guiho-convention-0001-cli-compliance-migration.md)
- [Plan review](../reviews/plans/guiho-convention-0001-cli-compliance-migration-review.md)
- [Architecture](../architecture/guiho-convention-0001-cli-compliance-migration.md)
- [Architecture review](../reviews/architecture/guiho-convention-0001-cli-compliance-migration-review.md)
- [Compliance audit](../reviews/guiho-convention-0001-cli-compliance-audit.md)
- [TODO.md](../../TODO.md)
