---
name: RunX Independent Prerequisite Review and Activation
purpose: Record independent acceptance, reviewed source delivery, and verified activation of the existing local RunX tool.
description: Full scheduler and operating-policy review, meaningful runtime regressions, exact binary and skill provenance, task readbacks, and remaining readiness gaps.
created: 2026-10-04
flags: [technically-validated, human-review-pending]
tags: [review, runx, readiness, activation]
keywords: [issue-63, issue-64, native-background, filesystem-boundary, plain-push]
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Independent Prerequisite Review and Activation

## Verdict and Authority

**The bounded native-policy and inspection-lifecycle prerequisite is independently
accepted, delivered, and locally activated.** This is technical acceptance only;
both real issues remain OPEN/Testing for later human acceptance. It does not
certify RunX's complete setup or any GUIHO family.

Mode: dnd. A new native General reviewer used expressly authorized
`openai/gpt-6.1-sol#xhigh`, read the actual full conventions 0011/0002/0007/CLI
0001, owning instructions, complete implementation contracts/results, and actual
delta. No recursive delegation or CLI worker was used. Canonical future defaults
remain both-role MiMo, provider-default thinking and no variant.

Parent authority in `/tmp/opencode/2026-10-04-runx-prerequisite-review-brief.md`
expressly permits reviewed plain push, bounded corrective commits, replacement
of the existing `/usr/local/bin/runx`, and refresh of only the two existing changed
first-party skill files. It supersedes the implementer's earlier NOPUSH and
no-install gate for this independent reviewer.

## Acceptance Matrix

| Signal | Independent evidence and verdict |
| --- | --- |
| Full current operating policy | Read all active AGENTS clauses and inactive MANDUME framing. OpenCode-only/native-always, capable actual child permissions, explicit authorized model/available variant, background/completion events, human-off only, no polling/tailing/shell-wait, no CLI escape/fallback, and scoped push rules agree with conventions. |
| Meaningful policy negatives | Four positive and twelve negative bounded connected-clause probes pass, including old policy, CLI launch, tail/poll, permission escape, other harness, fabricated variant, silent fallback, unrelated-negation masking, primary-label worker, shell-wait, and completion-off. These are bounded textual checks, not provider runtime proof. |
| Real production scheduler boundary | The twelve-line gate returns before executable resolution and both spawn calls. Cobra parses the actual Boolean before pre-run; `SetInterspersed(false)` preserves post-selector child flags. Catalog load paths read files/HTTPS and do not persist them. |
| No actual filesystem or child effects | Repeated 28 source inspection cases and 38 payload fixtures plus five owning-catalog inspections pass. Actual installed-binary repetition also passes all 43 invocations. Snapshots preserve hashes, file/directory sets, modes and nanosecond mtimes, including protected legacy fixture indices, stale resources, missing resources and extras. |
| Process and transient-write proof | Installed and rebuilt payload traces show exactly one `execve` per inspection; no lifecycle worker/catalog child, mutating filesystem syscall, or network connect/send in these local-reference fixtures. |
| Non-vacuous regression | Independently regenerated baseline-root Go overlay exits 1 as expected, detecting both workers, actual maintenance/resource writes and cache replacement. Three source positive controls pass. |
| Intended real execution | Three installed-payload controls for ordinary run, parsed false and a post-selector dry-run flag execute both workers and the child. The forwarded child flag is written exactly. Fresh isolated cache prevents release-network checks; real login-shell profile helpers and local journal connections remain real-run behavior. |
| Existing compatibility | Full Go tests, vet, formatting and all-package build pass. Confirmation and child exit handling are unchanged; gated noninteractive/JSON dry runs still return 2 without scheduling. |
| Docs and artifact contract | README/DOCS/owning notes and both skill boundary sections agree. Both source artifact versions are quoted `0.4.2`; embedded and canonical maintenance sections match. No project release version was changed. |
| Tracking and delivery | Both actual owning issues are OPEN/Testing on GUIHO Project #2, exact Component `runx`. Full 13-commit ancestry matches the completed worker's declared paths; no held/unowned ancestor exists. Plain push succeeds and source refs match. |

## Checks and Evidence

Evidence root: `/tmp/opencode/runx-review/`. Exact commands/exits are in
`checks.json`, `build.json`, task/installed `*-verification.json`, `policy.json`,
`positive-runtime.json`, and `implementation-prepush.json` / `implementation-delivery.json`.

- `go test -count=1 ./cmd -run 'TestInspectionLifecycleFilesystemBoundary|TestRealRunRetainsLifecycleAndChildExecution' -v`: **0**.
- `go test -count=1 -overlay /tmp/opencode/runx-review/prior-overlay.json ./cmd -run TestInspectionLifecycleFilesystemBoundary -v`: **1, expected failure**.
- `go test -count=1 ./...`, `go vet ./...`, `go build ./...`: **0** each.
- `gofmt -l cmd/root.go cmd/lifecycle_boundary_test.go`: **0**, no output.
- Full owned `git diff --check`, implementer static checks, independent policy probes and retained original-payload evidence verification: **0**.
- Rebuilt and installed payload fixture checks, strengthened process/file/network syscall checks, three actual installed positive controls: **0**.

External verification scaffolding initially had a string/Path substitution error,
an incomplete positive manifest, and an overstrict real-run process/network oracle.
They were corrected outside the repository. The latter exposed existing `sh -lc`
profile helper processes and AF_UNIX journal connections, not inspection behavior.
All final checks pass; no scheduler correction was needed during review.

## Exact Local Activation and Safe Interface

- Checked clean source commit: `dccc8fe5aadb962d1160417e004979c7b19a7d29`.
- Raw identity: `0.0.0-dev.readiness.runx.dccc8fe`; Go `go1.27.1`, Linux AMD64 V1,
  CGO disabled, trimpath, VCS metadata `vcs.modified=false`.
- Build date: `2026-10-03T22:23:41Z` (UTC; record named for the session's Oct 4 date).
- Built payload: `/tmp/opencode/runx-review/runx`; installed atomically at existing
  `/usr/local/bin/runx`, same SHA-256
  `25ca8207c221dfde0c4472e126ab6239776b490e5dbdfd95059071aa79fbcdd8`.
- Preserved old binary: `/tmp/opencode/runx-review/backups/runx-before-35907a58fee5619b`,
  SHA-256 `35907a58fee5619bb240ff49c28187412fd9c89bfe7e6edc3fbd8c4eaa9fc252`.
- Self-test reports `ok`, protocol 1, target `runx-linux-amd64`, skill bytes 7221,
  instruction bytes 1293. This is an honest local dev build, not an upstream release.
- Only `/root/.agents/skills/guiho-s-runx/SKILL.md` and
  `/root/.claude/skills/guiho-s-runx/SKILL.md` were refreshed. Both previously
  exactly matched baseline embedded source; both now exactly match accepted
  embedded source, SHA-256 `6cf5e20f4db1e557f57dc0c57824ee3c316a06cf2b4dc3fb2d9354372d37f91c`.
  Old copies are backed up in the same task backup directory. Extras, their modes/
  mtimes, existing path set and symlinks are preserved; no custom body was found.

```sh
/usr/local/bin/runx --version
/usr/local/bin/runx __self-test
/usr/local/bin/runx check --cwd /root/swe/runx --format json
/usr/local/bin/runx list --cwd /root/swe/runx --format json
/usr/local/bin/runx describe --cwd /root/swe/runx --format json runx-tooling-go-test
/usr/local/bin/runx reveal --cwd /root/swe/runx runx-tooling-go-test
/usr/local/bin/runx run --cwd /root/swe/runx --dry-run --format json runx-tooling-go-test
```

Options precede selectors. `--dry-run=false` and post-selector `--dry-run` are
real execution. Do not add unauthorized `--yes` to inspect confirmation-gated
entries; describe/reveal do not require confirmation.

## Delivery, Tasks and Remaining Gaps

Reviewed implementation range: baseline
`4a9f8b9042ed8d2288a790b85c43d6803d0f5693` through
`dccc8fe5aadb962d1160417e004979c7b19a7d29`, thirteen commits/fifteen paths.
Plain `git push` exits 0. HEAD, origin/main and live main were verified equal at
the delivered endpoint, with a clean index/worktree. GitHub printed its existing
server-side PR-rule bypass notice for the authenticated actor; no bypass flag,
hook bypass, force or history rewrite was used. This review record and linked task
mirror corrections are a subsequent owned documentation unit; final delivery/ref
proof and executable/embedded input equivalence are in the parent result packet.

Fresh readback at `2026-10-03T22:30:09Z` confirms [policy #63](https://github.com/CGuiho/runx/issues/63)
and [runtime #64](https://github.com/CGuiho/runx/issues/64) remain OPEN/Testing,
Component `runx`, Project #2. Membership/field pagination is exhausted. Existing
option IDs, colors and task item identities are preserved. [TODO.md](../../../TODO.md)
and the two owning specs remain Testing. No issue was closed or archived.

The existing catalog/config/descriptors, historical MANDUME body and unrelated
tasks retain their bytes. XDocs CLI validation is skipped: installed-safe XDocs
activation is separately owned and not accepted by this reviewer. RunX's missing
explicit descriptor grants and stale/missing source metadata remain separate
readiness setup work; no grants were broadened or automatic maintenance invoked.
Static checks cover this record's frontmatter/copyright/links. Linux AMD64 is the
only native smoke platform; Windows/macOS installer/runtime transactions and
foreign HTTPS catalogs were not tested. Complete docs/wiki/family setup, older
task binding reconciliation, and formal `g0000` confirmation remain open.

Prior [self-review](2026-10-04-readiness-runx-prerequisite-review.md),
[validation](../../validation/2026-10-04-readiness-runx-prerequisite.md), and
[handoff](../handoff/2026-10-04-readiness-runx-prerequisite.md) retain their historical
implementation-stage evidence. This independent review resolves their held
delivery/installation gates. Final machine-local return:
`/tmp/opencode/2026-10-04-readiness-runx-review-result.md` and matching JSON.
