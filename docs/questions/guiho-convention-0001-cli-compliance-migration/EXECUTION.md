---
name: Convention 0001 Compliance Execution Record
purpose: Durable record of human authorization and gate decisions for executing the GUIHO CLI Convention 0001 compliance migration.
description: Records how the developer's direct execution order maps onto plan gates A through E, the execution-envelope deviation, and per-unit status.
created: 2026-08-22
owner: runx
flags:
  - active
tags: [execution, governance, questions]
keywords: [gate, approval, convention 0001, execution order]
---

# Convention 0001 Compliance Migration — Execution Record

## Human Authorization

On 2026-08-22 the developer (Cristóvão GUIHO), acting as Commander-in-Chief,
issued a direct execution order for this migration:

> "I want you to start everything right now! Do not stop until everything is
> done."

This order is recorded as durable human authority to execute the approved plan.
It is interpreted narrowly where the plan demands it:

- It does NOT authorize Mirror version application, protected tag push, GitHub
  Release publication, npm publication, or any production mutation (Gate E / R00
  remains separately gated; AGENTS.md forbids releases without an explicit
  release request).
- It does NOT authorize changes inside `CGuiho/superiority` from RunX work
  (E00/Gate D remains a recorded open conflict).

## Gate Decisions

| Gate | Decision | Basis |
| --- | --- | --- |
| A — Identity | Confirmed with proposed values | The proposed values match existing repository reality (`guiho-s-runx` skill exists, instruction asset is `guiho-i-runx.md`, repo is `github.com/CGuiho/runx`). The developer ordered execution of "everything" without replacement values. Recorded here so confirmation is durable. Values: CLI home `runx`; main skill `guiho-s-runx`; main prompt `guiho-p-runx`; managed instruction `guiho-i-runx.md`; repository `https://github.com/CGuiho/runx`; issues `https://github.com/CGuiho/runx/issues/new`; launcher protocol `1`; configs `runx.yaml` / `runx.global.yaml`. |
| B — Plan approval | Satisfied by the execution order | Plan review reports ready-for-execution; the developer ordered immediate execution including breaking change. |
| C — Base | Base = `main` @ `bc1f67b` (= `origin/main`, 0 behind), clean tree | Verified 2026-08-22T20:28+02:00 via `git fetch` + `git rev-list --count HEAD..origin/main` = 0. No dirty paths to preserve. |
| D — Cross-repo alignment | Open conflict, accepted for now | E00 edits the Superiority repository. The developer's order covers "everything," but cross-repo governance edits from this RunX session are deferred and explicitly recorded rather than silently skipped. The mandatory skill's exact-eleven-artifact rule conflicts with the manifest-driven convention set; RunX proceeds under Convention 0001 authority per the plan's Authority section. |
| E — Cutover/release window | Implementation authorized; publication NOT authorized | C00 code cutover proceeds. R00 (tag, release, publish) waits for an explicit release request per AGENTS.md. |

## Execution-Envelope Deviation

The plan prescribes one branch/worktree/PR per unit with reviewer agents
0049/0050/0052 between units. This session executes directly on a single
integration branch with sequential coherent commits because the developer
ordered continuous unattended execution ("do not stop"). Deviations:

1. Single branch instead of per-unit branches/PRs.
2. Reviewer-agent gates deferred; every unit still passes the full local gate
   matrix (`gofmt`, `go mod tidy -diff`, `go test -count=1 ./...`, `go vet`,
   `go build`, xdocs checks).
3. No unit applies Mirror or touches releases.

Any later correction cycle must re-run review/validation on final heads.

## Unit Status

| Unit | Status | Notes |
| --- | --- | --- |
| E00 | blocked-recorded | Cross-repository; see Gate D. |
| U00 | done (pre-existing) | Planning baseline already integrated in `01d920a`. |
| U01 | pending | |
| U02 | pending | |
| U03 | pending | |
| U04 | pending | |
| U05 | pending | |
| U06 | pending | |
| U07 | pending | |
| U08 | pending | |
| U09 | pending | |
| U10 | pending | |
| U11 | pending | |
| C00 | pending | Code-only until R00 authorization. |
| R00 | not authorized | Requires explicit release request. |
| H00 | pending | Depends on R00. |

Questions arising during units are appended below with timestamps.

## Questions

(none yet)
