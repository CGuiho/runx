---
name: RunX Upgrade Reliability Implementation Plan
purpose: Sequence the approved RunX upgrade reliability design into executable and independently verifiable units.
description: Implementation plan for release discovery, progress reporting, transactional replacement, recovery commands, hardened installers, documentation, and validation for GitHub issues 12 and 13.
created: 2026-07-15
flags:
  - approved
  - implemented
tags:
  - plans
  - cli
  - reliability
keywords:
  - runx upgrade
  - upgrade list
  - windows replacement
  - recovery install
owner: runx-plans
---

# RunX Upgrade Reliability Implementation Plan

## Source of Truth

- Approved design: [upgrade reliability design](../superpowers/specs/2026-07-15-upgrade-reliability-design.md)
- GitHub issues: `CGuiho/runx#12` and `CGuiho/runx#13`
- Task: [Upgrade Reliability](../todo/upgrade-reliability.md)

## Unit 1: Release Catalog and Recovery Contract

- Goal: Replace one-page release lookup and ad hoc version comparison with a complete, deterministic catalog and reusable platform recovery instructions.
- Owner: RunX package, `source/`.
- Dependencies: approved design only.
- Expected files: new focused catalog/recovery modules and colocated specs; `source/source.xdocs.md`.
- Data impact: introduces the schema-versioned catalog envelope and normalized release metadata; no persisted data.
- Auth/cache impact: unauthenticated GitHub API access remains; no upgrade cache is introduced.
- Checks: pagination with `Link`, malformed/non-2xx failures, SemVer precedence, channel labels, non-SemVer tail ordering, compatible assets, exact pinned recovery commands.
- Acceptance: every non-draft published release is present newest first, latest stable is computed independently, asset policy is deterministic, and commands are safe/copyable.
- Stop condition: stop if implementing the approved contract requires changing release tag or asset naming.

## Unit 2: Planner, Events, Transactional Installer, and Reporter

- Goal: Make one execution pipeline expose the immutable plan before download, stream ordered text phases, buffer equivalent JSON, and return only after verified canonical replacement or rollback.
- Owner: RunX package, `source/`.
- Dependencies: Unit 1.
- Expected files: self-management modules/specs, shared types, CLI integration/specs, `source/source.xdocs.md`.
- Data impact: adds schema-versioned upgrade envelope and typed phase events.
- Auth/cache impact: no change; `cache` is explicitly reported as skipped.
- Checks: unique temporary paths, native magic validation, download-before-replace ordering, Windows rename/swap/verify/rollback, POSIX swap/verify/rollback, cleanup semantics, already-current/dry-run/failure recovery output, single JSON document.
- Acceptance: the plan is visible before awaiting the body, every slow phase is announced first, success proves the canonical binary reports the exact target, and all outcomes include recovery instructions.
- Stop condition: stop if Citty must be bypassed or a second token router would be required.

## Unit 3: CLI Catalog Presentation and Error Boundary

- Goal: Wire `runx upgrade`, `check`, and `list` to the shared envelopes without changing Citty command ownership.
- Owner: RunX CLI, `source/cli.ts` and colocated tests.
- Dependencies: Units 1 and 2.
- Checks: complete aligned list output, current/latest/asset markers, exact JSON, output ordering, concise stderr failures plus stdout recovery, help descriptions.
- Acceptance: text and JSON expose equivalent facts; JSON remains one parseable document; failures exit nonzero without suppressing recovery.
- Stop condition: stop if global error handling would leak manifest behavior into self-management commands.

## Unit 4: Direct Installer Hardening

- Goal: Make PowerShell and POSIX installers honor exact versions, unique temporary state, native validation, verified replacement, rollback, and actionable failure messages.
- Owner: RunX package, `devops/`.
- Dependencies: Unit 1 asset and version contract.
- Expected files: `devops/install.ps1`, `devops/install.sh`, focused contract/smoke tests, `devops/devops.xdocs.md`.
- Checks: exact stable/prerelease version, candidate ordering, version verification, rollback, nonzero failure, cleanup, paths with spaces.
- Acceptance: installers never report success until the installed canonical binary reports the requested resolved version.
- Stop condition: if a platform cannot be safely exercised locally, retain contract tests and record the unavailable live check.

## Unit 5: User and Agent Documentation

- Goal: Document the real upgrade/list/recovery behavior and keep bundled agent guidance aligned.
- Owner: RunX package docs and skill.
- Dependencies: Units 1-4.
- Expected files: `README.md`, `DOCS.md`, `skills/guiho-s-runx/SKILL.md`, affected xdocs descriptors, task state.
- Checks: examples match actual command output and installer flags; xdocs metadata strictness/tree/doctor.
- Acceptance: users can identify versions/channels, understand phases, and copy the pinned recovery command after every attempt.

## Unit 6: Full Validation and Handoff

- Goal: prove the branch is release-ready without publishing.
- Dependencies: Units 1-5.
- Checks: `bun run typecheck`, `bun test`, `bun run build`, `bun run binary`, `bun run binaries` when feasible, focused installer smokes, `xdocs meta ... --strict`, `xdocs doctor`, and `xdocs tree`.
- Acceptance: all feasible checks pass, deviations are recorded, no generated output is committed, and task status is `testing` or `completed` only when evidence supports it.
- Stop condition: do not publish, tag, or apply a Mirror version.

## Execution Order and Checkpoints

Execute Units 1 through 6 in order. Commit coherent units with explicit file staging. Architecture is approved and the user explicitly authorized uninterrupted implementation; no additional approval gate exists unless a design contradiction, release mutation, or destructive action is discovered.

## First Executable Unit

Unit 1: implement and test the release catalog and recovery contract.
