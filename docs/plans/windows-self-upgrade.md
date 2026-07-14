---
name: RunX Windows Self-Upgrade Plan
purpose: Execute the accepted synchronous and recoverable Windows self-upgrade design.
description: Sequences executable replacement, regression coverage, documentation, validation, issue closure, and the next Mirror patch.
created: 2026-07-14
flags:
  - approved
  - executed
tags:
  - plans
  - cli
  - windows
keywords:
  - runx upgrade
  - executable replacement
  - rollback
  - patch release
owner: runx-plans
---

# RunX Windows Self-Upgrade Plan

## Sources

- [Windows Self-Upgrade Decision](../decisions/windows-self-upgrade.md)
- [Alpha Command Catalog Requirements](../requirements/alpha-command-catalog.md)
- [CLI Architecture](../architecture/cli-architecture.md)
- GitHub issues `CGuiho/runx#9` and `CGuiho/runx#1`

## Unit 1: Replace the Windows Executable Synchronously

- Goal: Make `upgradeSelf()` return successfully only after the canonical
  Windows executable path contains and runs the target release.
- Owner: `C:\GUIHO\runx`.
- Dependencies: Accepted Windows self-upgrade decision.
- Files:
  - `source/self-management.ts`
  - `source/types.ts` only if the result contract must change
- Behavior:
  - Download to `.new`, rename the installed executable to `.old`, install the
    new file at the original path, and verify `<self-path> --version`.
  - Return `scheduled: false` after verified replacement.
  - On failure, remove the failed target, restore `.old`, remove `.new`, and
    return a non-zero `RunXError` with recovery context.
  - Delete `.old` immediately when possible; defer only deletion of the locked
    old image after the current process exits.
- Data/schema, auth, permissions, and cache: No application data, schema,
  authentication, authorization, or cache impact. Filesystem write permission
  to the installed executable remains required.
- Acceptance: The installed path is the verified target before
  `upgradeSelf()` resolves and no `.new` file remains.
- Stop conditions: Stop if Windows cannot rename the running executable, the
  target cannot be verified, rollback cannot restore the installed path, or
  unrelated self-management behavior changes.

## Unit 2: Add Regression Coverage

- Goal: Prove the fixed upgrade contract and preserve the already-correct help
  behavior from issue #1.
- Dependencies: Unit 1 complete.
- Files:
  - `source/guiho-runx.spec.ts` or a focused colocated self-management spec
  - existing `source/cli.spec.ts` only if help coverage needs correction
- Checks:
  - Mock the GitHub latest-release response and downloaded Windows binary.
  - Configure `RUNX_SELF_PATH` to an isolated executable path.
  - Assert the replacement exists at that path before the promise resolves,
    reports `scheduled: false`, runs as the expected version, and leaves no
    `.new` file.
  - Force verification failure and assert restoration of the original file.
  - Run the existing `-h`/`--help` tests outside a manifest directory.
- Acceptance: Windows CI exercises success and rollback; issue #1 remains
  covered without changing its implemented behavior.
- Stop conditions: Stop on nondeterministic process cleanup, leaked test
  processes, platform-dependent assertions on non-Windows runners, or failure
  to restore global fetch/environment state.

## Unit 3: Align Documentation and Validate

- Goal: Make public and structured documentation describe completed rather than
  scheduled Windows replacement, then pass the full release gate.
- Dependencies: Units 1 and 2 complete.
- Files:
  - `DOCS.md`
  - `source/source.xdocs.md`
  - `runx.xdocs.md`
  - accepted decision and plan lifecycle flags when implementation completes
- Checks:
  - `bun run typecheck`
  - `bun test`
  - `bun run build`
  - `bun run binary`
  - `bun run binaries`
  - `xdocs meta . --documents --strict --owner runx --format json`
  - `xdocs doctor . --warnings-as-errors --format json`
  - `git diff --check`
- Acceptance: All checks pass, docs agree with behavior, and the worktree
  contains only understood task changes.
- Stop conditions: Stop on any failed validation or unexpected worktree change.

## Unit 4: Close Issues and Apply the Patch

- Goal: Close both resolved GitHub issues and publish the next
  Mirror-managed patch state.
- Dependencies: Unit 3 committed and the worktree clean.
- Sequence:
  - Confirm issue #9 acceptance evidence and issue #1 Citty regression evidence.
  - Close issues #9 and #1 with concise resolution comments.
  - Run `mirror config show` and `mirror version plan patch`.
  - Update the configured changelog with the planned next version and commit
    release preparation.
  - Run the full typecheck/test gate again if release preparation changes any
    executable inputs.
  - Apply `mirror version apply patch --yes`; allow configured commit/push/tag
    behavior and do not hand-edit Mirror-managed fields.
- Acceptance: Both issues are closed, the changelog records the fix, and
  package version, release commit, tag, and configured remote push agree.
- Stop conditions: Stop if an issue acceptance criterion is unproved, Mirror
  plans an unexpected version, validation fails, or protected push is rejected.

## TODO Alignment

This is one approved, focused execution session. No additional durable TODO
entry or delegated component task is required.

## First Executable Unit

Unit 1: implement synchronous Windows replacement and rollback in
`source/self-management.ts`.
