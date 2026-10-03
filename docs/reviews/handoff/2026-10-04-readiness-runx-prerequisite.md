---
name: RunX Readiness Prerequisite Review Handoff
purpose: Give the parent an executable entry point for new independent review and subsequent delivery decisions.
description: Exact safe task-binary commands, acceptance and evidence navigation, pending review choices and no-push state.
created: 2026-10-04
flags: [testing, independent-review-pending]
tags: [handoff, runx, readiness]
keywords: [issue-63, issue-64, safe-invocations, no-push]
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Readiness Prerequisite Review Handoff

## What Was Built

Owned `CGuiho/runx` main commits deploy the current native background operating
policy and the separately authorized minimum runtime prerequisite: `check`,
`list`, `describe`, `reveal`, and parsed `run --dry-run=true` return before
lifecycle scheduling. Documentation and both bundled skill contracts agree.
Real execution, flag ownership, opt-in confirmation and explicit mutations retain
their existing behavior. Mode: dnd; no recursive workers were launched.

## How to Inspect and Verify

Use the identified task payload, not an unidentified PATH installation:

```sh
task_binary=/tmp/opencode/runx-prerequisite-task/runx
sha256sum "$task_binary"
go version -m "$task_binary"
"$task_binary" --version
"$task_binary" __self-test
"$task_binary" check --cwd /root/swe/runx --format json
"$task_binary" list --cwd /root/swe/runx --format json
"$task_binary" describe --cwd /root/swe/runx --format json runx-tooling-go-test
"$task_binary" reveal --cwd /root/swe/runx runx-tooling-go-test
"$task_binary" run --cwd /root/swe/runx --dry-run --format json runx-tooling-go-test
```

Expected SHA-256:
`cd76c511fba3c3751ea535c2f77b2c4bdb5ff97ae1fdbca3cbe6def6ffdeef5a`.
Expected raw version: `0.0.0-dev.readiness.runx.d9f1675`.
Self-test reports protocol `1`, target `runx-linux-amd64` and readable resources.
These inspection invocations return `0`, schedule no worker and execute no child.
No environment toggle is needed. Referenced HTTPS catalogs may still be read.

RunX-owned flags belong **before** the selector. Parsed `--dry-run=false` or a
post-selector `--dry-run` is real execution. A noninteractive or JSON dry run of
an unapproved `confirm: always` command returns `2`; inspect it with `describe`
or `reveal` rather than adding an unauthorized `--yes`.

The retained fixture verifier, if repetition is needed for independent review:

```sh
python3 /tmp/opencode/runx-prerequisite-binary-verify.py
python3 /tmp/opencode/runx-prerequisite-static-check.py
```

## Acceptance, Validation and Questions

- [Technical self-review](../implementation/2026-10-04-readiness-runx-prerequisite-review.md)
  maps every assigned policy/runtime acceptance signal and scoped finding.
- [Validation record](../../validation/2026-10-04-readiness-runx-prerequisite.md)
  gives commands/exits, 28 source inspection cases, three real-run controls,
  the expected-failing old-source overlay, Go checks, 38 payload fixture cases,
  five owner-catalog inspections, process traces, snapshots and build provenance.
- [Default inspection decision](../../questions/readiness-tool-prerequisite/2026-10-03-inspection-default.md)
  records the parent-authorized minimal source waiver, safest reversible answer
  and pending later human review.
- XDocs CLI validation and descriptor synchronization are explicitly skipped:
  independently owned installed-build safety is not accepted here; current
  write grants are empty and no named descriptor edits were authorized.
- Native platform smoke coverage is Linux AMD64. Existing PATH binary and global
  skill installations were not replaced; no project version bump or release ran.

## Task and Delivery State

[Issue #63](https://github.com/CGuiho/runx/issues/63) and
[issue #64](https://github.com/CGuiho/runx/issues/64) remain **OPEN / Testing**
on [GUIHO Project #2](https://github.com/users/CGuiho/projects/2), exact
Component `runx`, with verified local mirrors in [TODO.md](../../../TODO.md).
Historical tasks are preserved. This work does not certify any family.

All task commits remain local on main under the parent's explicit **NOPUSH**
gate. No unowned changes or held-ahead ancestry were observed. Parent must
perform new independent review, inspect live remote-main and the entire outgoing
range, and separately authorize plain push and any local installed replacement.
Those gates are pending acceptance, not completed by this self-review packet.

Machine-local complete result packet, with final hashes and task readbacks:
`/tmp/opencode/2026-10-03-readiness-runx-prerequisite-result.md` and matching JSON.
