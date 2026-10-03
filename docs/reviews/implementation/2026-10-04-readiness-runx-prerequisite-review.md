---
name: RunX Readiness Prerequisite Technical Review
purpose: Map the assigned policy and minimum runtime prerequisite to independently inspectable evidence.
description: Scoped self-review of native operating policy, scheduling exclusions, meaningful regressions, documentation and task delivery gates.
created: 2026-10-04
flags: [testing, independent-review-pending]
tags: [review, runx, readiness]
keywords: [issue-63, issue-64, native-policy, lifecycle-boundary]
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Readiness Prerequisite Technical Review

## Scope and Verdict

Mode: dnd. Self-review of owned work in `CGuiho/runx` from clean main baseline
`4a9f8b9042ed8d2288a790b85c43d6803d0f5693`. Assigned contracts live at
`/tmp/opencode/2026-10-03-native-deployment-execution-contract.md` and
`/tmp/opencode/2026-10-03-readiness-tool-prerequisite-contract.md`.

The assigned policy and runtime work pass the recorded scoped checks and are
ready for **new independent parent review**. This is not independent acceptance,
installed-tool verification, human acceptance or family readiness. Both owning
issues remain OPEN/Testing; push and installed replacement are held.

## Acceptance Mapping

| Required outcome | Owned implementation | Evidence / conclusion |
| --- | --- | --- |
| Current OpenCode-only, always-native background policy | Root AGENTS Mandume registry and native completion contract | Current-policy static acceptance; baseline-policy rejection; actual shell/Git/GitHub/write capability in this authorized leaf |
| Actual task-capable child permissions and authorized model/variant | General versus read-only Explore distinction, explicit model/variant selection, two-role MiMo default with empty variants | Actual 0007/0011 read; no invented MiMo effort tier, active CLI escape or fallback; future capability remains unproven |
| DND human-off does not disable native completion | Mode comments and native completion loop | No recursive delegation, native worker polling/tailing or human wake-up request performed |
| Preserve historical transport without active contradiction | Prominent inactive MANDUME banner and pointer to current registry | Historical body byte-identical to baseline; old harness/fallback/tail clauses explicitly inactive |
| Distinct real policy and runtime tasks | Issues #63 and #64; root TODO and separate specs | Fresh owning repository, membership, Component `runx`, OPEN/Testing readbacks; local mirrors agree |
| Inspection and dry-run cannot schedule workers | Twelve-line gate in `cmd/root.go` before executable resolution | 28 meaningful source cases, old-source overlay fails, 38 actual payload fixture invocations trace zero workers/children |
| Parsed flag ownership remains correct | Existing Cobra parser and RunX-before-selector ownership retained | Three positive real-run controls: ordinary run, parsed false, post-selector child flag |
| Project/global files and missing resources preserved | Isolated fixture snapshots with real maintenance and worker/child sentinels | Bytes/hashes, file sets, modes and mtimes match; old scheduler demonstrably changes them |
| Public, instruction and bundled-skill semantics match source | README, DOCS, owning AGENTS notes, canonical and embedded skills | Both skill metadata versions `0.4.2`, shared maintenance contract, confirmation unchanged |
| Honest executable identity | Explicit Linux AMD64 task payload, not installed globally | Clean built commit, version/target/self-test, SHA-256 and Go VCS metadata; PATH binary hash preserved |

## Findings and Corrections

No unresolved defect in the scoped scheduler implementation was found. The gate
uses parsed command flags rather than raw argument scanning, returns before
worker executable resolution and preserves the existing real-run/explicit
mutation paths. Tests reject both worker scheduling and actual resource/cache
side effects, including missing-resource creation; their positive controls
prove the oracles are not vacuous.

Documentation now distinguishes read-only catalog paths from real execution,
bare bootstrap and explicit resource mutations. It documents that referenced
HTTPS catalogs may be read, gated dry runs may return `2`, parsed false is real
execution, and post-selector flags belong to the child. Confirmation code and
manifest values were not modified. The embedded skill's obsolete nearest-file
maintenance description was corrected to the existing repository-root contract.

The static checker initially expected different wording for the dormant future
harness clause; correcting that external checker to the actual accepted wording
passed without changing repository policy. A transient GitHub HTTP 503 recovered
after independent build work; no mutation occurred in the failed attempt.

## Validation, Deviations and Remaining Review

The [validation record](../../validation/2026-10-04-readiness-runx-prerequisite.md)
contains exact commands/exits, process traces, fixtures, binary provenance and
remote readback identity. Go suite, vet, changed-file formatting, all-package
build, payload self-test, actual payload inspection and scoped static checks pass.

The [default-boundary decision](../../questions/readiness-tool-prerequisite/2026-10-03-inspection-default.md)
records the expressly authorized minimal runtime waiver and remains human-review
pending. No schema/flag/global opt-out was added. Skill artifact patch metadata
is required by Convention 0002; no project/release version bump was applied.

XDocs validation and descriptor synchronization remain unverified: installed
runtime safety belongs to another prerequisite owner, grants are empty and no
named descriptor edits were authorized. Stale generated instruction blocks,
historical task normalization and complete family setup remain outside this
unit. Linux AMD64 is the only native smoke platform proven here.

Parent must review the full owned outgoing ancestry and verify live remote-main
before separately authorizing delivery or local installation. No unowned
worktree/index changes were observed. Formal `g0000` remains only an observed
runtime-family identifier, not a newly verified registry ID.
