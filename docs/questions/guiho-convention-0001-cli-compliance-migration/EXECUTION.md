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
| U01 | done (core) | Root `runx.yaml` catalog + `xdocs.yaml`; toml retirement deferred until installed xdocs CLI parity is verified. |
| U02 | done (core) | Strict config contracts, policy enum, schemas, examples, validator; full interactive init orchestrator delivered with the init command. |
| U03 | partial | Manifest carries optional strict agent.evolution; canonical byte-pinned resource graph remains open. |
| U04 | partial | Release manifest contract shipped via build/verify tooling; pkg/release selection package remains open. |
| U05 | done (core) | installstate paths, pointer, ledger, atomic writes, staging; lock/journal/instance-registry primitives remain open. |
| U06 | done (core) | Stable launcher binary with delegation and fallback; capability-token dispatch remains open. |
| U07 | done | Protocol-v1 release matrix: 8 payloads + 8 launchers + resources + artifacts.json + checksums. |
| U08 | done (core) | Both installers rewritten to the canonical layout with channels, self-test, and rollback. |
| U09 | done (core) | Shared uninstallation contract across Cobra and both scripts. |
| U10 | done (core) | Recovery blocks plus protocol-v1 whole-release engine with legacy fallback. |
| U11 | done | --help-tree-global-flags implemented per convention. |
| C00 | partial | Public surface wired for uninstall/upgrade/init/help-tree; npm retirement and DOCS.md refresh remain open. |
| R00 | **done** | 0.13.0 released 2026-08-22; see docs/releases/guiho-convention-0001-first-protocol-v1.md. |
| H00 | unblocked | Transition-support removal now possible against the published protocol-v1 release. |

Questions arising during units are appended below with timestamps.

## Questions

(none yet)
