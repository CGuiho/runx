---
name: RunX RFC 0034 CLI Compliance Migration Plan Review
purpose: Verify that the RunX RFC 0034 migration plan is sequenced, explicit, testable, and safe to execute.
description: Reviews the breaking migration plan against the CLI contract, current RunX architecture, TODO tracking, documentation duties, and release approval gates.
created: 2026-07-18
flags:
  - approved
  - ready-for-execution
tags:
  - review
  - plan
  - cli
keywords:
  - runx
  - RFC 0034
  - plan readiness
  - breaking migration
owner: runx-plan-reviews
---

# RunX RFC 0034 CLI Compliance Migration Plan Review

## Verdict

Ready for execution.

## Findings

No blocker or high-severity finding remains.

- Medium, resolved: Earlier RunX requirements and decisions require `runx r`,
  root selector shorthand, `--file`, upward discovery, and a home page. The new
  plan records the developer-approved breaking boundary and requires conflicting
  durable documents to be superseded.
- Medium, resolved: Current agent code only installs/updates one selected tool,
  writes AGENTS only, uses non-RFC markers, and has no prompt catalog. RX-09
  explicitly covers both skill paths, all actions, instruction resolution,
  idempotency, prompts, embedding, and isolated tests.
- Medium, resolved: Current release work builds twelve `macos` binaries but no
  agent assets. RX-12 through RX-14 sequence installer, bootstrap, naming,
  packaging, workflow, and exact-set verification.
- Medium, resolved: The npm launcher requirement could conflict with Bun-only
  core source. RX-13 names `scripts/runx-bin.mjs` as the isolated Node exception.
- Low, resolved: Creating plans/specs/reviews requires xdocs companion metadata.
  The affected descriptors are included in this planning unit.

## Sequencing Risks

The Bun-only source foundation and TypeBox boundaries correctly precede
configuration, startup, agent, and upgrade behavior. The command tree precedes
generated help. Agent resources precede installers and release packaging.
Documentation follows stable behavior and precedes final validation.

## Acceptance Criteria Review

Every unit states a goal, dependencies, files or modules, actions, acceptance
signals, and stop/approval conditions where risk exists. The validation unit
covers source, CLI, native binaries, npm bootstrap, installers, prohibited
imports, fourteen assets, xdocs, review, and release boundaries.

## TODO Alignment

RunX TODO task `1` links the task specification and executable plan. It remains
`todo` and cannot become completed until validation evidence proves every RFC
completion-gate item.

## First Executable Unit

RX-01: record baseline commands, checks, prohibited imports, agent behavior, and
release assets in an understood worktree.

## Recommended Next Skill

Use `guiho-s-0023-plan-executor` with `guiho-s-0034-cli-engineer`.

## References

- [Migration plan](../../plans/rfc-0034-cli-compliance-migration.md)
- [Task specification](../../todo/rfc-0034-cli-compliance-migration.md)
- [TODO.md](../../../TODO.md)
