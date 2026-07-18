---
name: Migrate RunX To Full RFC 0034 Compliance
purpose: Define the required outcome, scope, constraints, and completion signals for the RunX RFC 0034 migration task.
description: Captures what must be true after RunX adopts the complete GUIHO CLI contract and links the executable migration plan.
created: 2026-07-18
flags:
  - approved
  - breaking-change
  - completed
tags:
  - todo
  - cli
  - migration
keywords:
  - runx
  - RFC 0034
  - guiho-s-0034-cli-engineer
  - agent namespace
  - fourteen release assets
owner: runx-todo
---

# Migrate RunX To Full RFC 0034 Compliance

## Summary

Bring RunX into complete compliance with GUIHO RFC 0034. This is an approved
breaking migration for a pre-1.0 CLI. Compatibility with `runx r`, root selector
shorthand, `--file`, upward manifest discovery, `runx agents`, `--tool`, old
markers, `macos` asset names, and the Bun-dependent npm launcher is not required.

## Todo Index

- Task: `1. Migrate RunX To Full RFC 0034 Compliance`
- Status: completed
- Index: [TODO.md](../../TODO.md)

## Outcome

RunX uses the mandatory Bun/TypeScript/Citty/TypeBox stack, has a Bun-only core,
implements the exact startup/configuration/help/agent/upgrade/installer/output
contracts, ships a Node-compatible npm bootstrap, and publishes exactly twelve
RFC-named binaries plus `guiho-s-runx` and `guiho-i-runx`.

## Scope

### In scope

- Core CLI source, entrypoints, command tree, domain adapters, and errors.
- `runx.yaml` resolution and TypeBox decoding.
- update cache and background worker under `~/.guiho/runx/`.
- all Developer Context help modes at every scope.
- complete singular `agent skill`, `agent instruction`, and `agent prompt`
  namespaces.
- transactional upgrade, catalog pagination, and post-upgrade reconciliation.
- PowerShell and POSIX installers.
- Node-compatible npm bootstrap.
- exact fourteen-asset build and release enforcement.
- tests, CI, README, DOCS, bundled skill, TODO, changelog, decisions, and xdocs.

### Out of scope

- Workflow dependency graphs or a proprietary task language.
- Remote command execution.
- Secret management.
- Publishing, version bumping, tagging, or live installation without separate
  authorization.

## Acceptance Signals

- No prohibited Node built-ins exist in core CLI source.
- No arguments prints exactly `Hello Windows - runx v<version>`.
- `runx.yaml` resolves only by `--config`, effective cwd, then the standardized
  global path, and prints its absolute loaded path.
- TypeBox decodes configuration, cache, remote releases, structured flags, and
  stable output objects before use.
- Only `-h` and root `-v` short aliases exist.
- Every command scope supports standard help, tree help, positive tree depth,
  and redirect-safe Markdown help.
- The final command catalog in the plan is the only public catalog.
- Agent skill actions target both agent-tool directories; instruction actions
  handle zero, one, or both instruction files idempotently; prompt output obeys
  raw/names-only rules.
- Upgrade, list pagination, pre-release filtering, installers, and the npm
  wrapper pass isolated tests.
- Release verification finds exactly fourteen correctly named assets and no
  `macos` names.
- Typecheck, tests, safe builds, xdocs validation, implementation review, and
  validation reporting complete successfully.

## Dependencies And Context

- [Executable migration plan](../plans/rfc-0034-cli-compliance-migration.md)
- `guiho-a-0001-swe` is the coordinating Software Engineer/SWE agent.
- `guiho-s-0034-cli-engineer` is the mandatory specialist skill.
- The current Citty migration, upgrade work, and bundled RunX skill are inputs,
  not constraints that override RFC 0034.

## Watch-outs

- Inspection and dry-run commands must never spawn a manifest command.
- JSON stdout must remain one parseable document.
- The foreground startup path must never await a network request.
- Tests must isolate home directories, skill installation, instruction files,
  update caches, executable replacement, and release hosting.
- Generated `library/`, `bin/`, `bundle/`, and `vendor/` outputs remain
  unedited.
- Breaking changes are approved, but unrelated RunX domain behavior must not be
  removed accidentally.

## Before Starting

- Read the repository and parent instructions, this spec, and the full plan.
- Load the SWE agent and every skill named in the plan.
- Confirm the live branch, clean/understood worktree, and baseline checks.
- Inventory current agent behavior and prohibited imports again because source
  may have changed after this plan was written.

## While Working

- Execute one numbered plan unit at a time with `guiho-s-0023-plan-executor`.
- Update implementation notes, TODO state, docs, tests, and xdocs in the same
  unit.
- Do not preserve an obsolete interface merely to avoid a breaking change.
- Stop on unapproved publishing, tagging, pushing, or real global installation.

## After Finishing

- Run implementation review and validation reporting.
- Record every passed, failed, and skipped check.
- Keep the TODO in `testing` until all local acceptance signals are proven.
- Request separate authorization before Mirror versioning or release actions.

## Related Files

- [Implementation plan](../plans/rfc-0034-cli-compliance-migration.md) -
  Ordered file-level execution units and validation gates.
- [CLI architecture](../architecture/cli-architecture.md) - Current architecture
  that the migration must revise where it conflicts with RFC 0034.
- [Previous Citty decision](../decisions/citty-cli-migration.md) - Earlier
  compatibility contract that must be superseded where it conflicts.
- [Plan review](../reviews/plans/rfc-0034-cli-compliance-migration-review.md) -
  Ready-for-execution review of sequencing, coverage, and approval gates.

## References

- [TODO.md](../../TODO.md)
- [AGENTS.md](../../AGENTS.md)
- [Implementation notes](./rfc-0034-cli-compliance-migration-implementation.md)
- [Implementation review](../reviews/implementation/rfc-0034-cli-compliance-migration-review.md)
- [Validation report](../validation/rfc-0034-cli-compliance-migration.md)
