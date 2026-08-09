---
name: RunX 0.11.0 Release Preparation Integration Validation
purpose: Preserve exact-head validation, protected merge, main reachability, and final publication evidence for RunX 0.11.0.
description: Records the ready validation gate, successful merge, reachable source state, canonical tag, successful publication workflow, asset verification, and durable evidence for the RunX 0.11.0 release.
created: 2026-08-09
flags:
  - ready
  - integrated
  - release
tags:
  - validation
  - integration
  - release
  - pull-request
keywords:
  - runx
  - 0.11.0
  - pull request 41
  - exact head
  - merge commit
  - mirror
owner: runx-validation
---

# RunX 0.11.0 Release Preparation and Publication Validation

## Verdict

The preparation gate was ready for the pull-request integrator at exact head
`5126294b3ecc74b6d606c8b4159b03d01ea9de6f`; the gate was reobserved before
merge and remained valid. Final publication verification is complete below.

## Immutable Validation Evidence

- Validation record: [PR #41 validation comment](https://github.com/CGuiho/runx/pull/41#issuecomment-5229738249).
- Pull request: [CGuiho/runx#41](https://github.com/CGuiho/runx/pull/41),
  `codex/release-0.11.0-prep` -> `main`.
- Validated head: [`5126294b3ecc74b6d606c8b4159b03d01ea9de6f`](https://github.com/CGuiho/runx/commit/5126294b3ecc74b6d606c8b4159b03d01ea9de6f).
- Accepted implementation review bound to the same head: [review-4890485277](https://github.com/CGuiho/runx/pull/41#pullrequestreview-4890485277).
- CI: [run 31293614915](https://github.com/CGuiho/runx/actions/runs/31293614915),
  successful for Ubuntu Go, Windows Go, and the canonical release-contract
  job.
- Exact-head checks passed: `gofmt`, clean `go mod tidy` diff,
  `go test -count=1 ./...`, `go vet ./...`, `go build ./...`, strict XDocs
  scan/meta/tree/doctor, Mirror config/current/plan, exact 11-artifact
  checksum verification, and native Windows AMD64 `--version`, `--help`, and
  `--help-tree` smoke.

## Merge Evidence

- Merge result: [`4ed75572acb32d35be12d36e1e760671300f2733`](https://github.com/CGuiho/runx/commit/4ed75572acb32d35be12d36e1e760671300f2733).
- The PR merge commit is an ancestor of the release source `main` commit
  `e2b86336ebd95bc6bf25d395f518a1dce66132e1`.
- The merged PR is closed, no open PR remains, and issues #36 and #39 are
  closed with completed state reasons.
- No implementation correction or production mutation occurred.

## Historical Protected-Credential Attempt

- An earlier attempt from the pre-credential clean clone reached the protected
  Mirror tag push but could not persist `wincredman` credentials. Its local tag
  was removed and no remote tag or release remained. This historical blocker is
  retained in PR #43's merged evidence; it did not bypass the canonical path.
- No deployment, promotion, traffic, DNS, database, secret, or other
  production-infrastructure mutation occurred.

## Final Publication Verification

- Clean final clone was synchronized to release source `main`
  `e2b86336ebd95bc6bf25d395f518a1dce66132e1`; `mirror config check` passed,
  `mirror version current` was `0.10.1`, and `mirror version plan minor`
  selected `@guiho/runx/v0.11.0`.
- The only apply command was `mirror version apply minor --yes`; Mirror reported
  `applied: true` and created/pushed annotated tag object
  `51edf2e95849ca636fb05a4a6cfe161176a02314`, peeling to
  `e2b86336ebd95bc6bf25d395f518a1dce66132e1`.
- Tag-only workflow [31295252704](https://github.com/CGuiho/runx/actions/runs/31295252704) completed successfully. The [GitHub release](https://github.com/CGuiho/runx/releases/tag/%40guiho%2Frunx%2Fv0.11.0) is non-draft and non-prerelease with exactly 11 assets.
- Independent download contained the exact 11 expected names; all ten
  payload hashes matched `checksums.txt` and each corresponding GitHub asset
  digest. Windows AMD64 download reported `0.11.0` for `--version`, while
  `--help` and `--help-tree` exited 0. `npm.cmd view @guiho/runx@latest version`
  returned `0.11.0`.
- No unmerged or arbitrary branch was touched. Merged release-prep,
  integration, and blocker branches were deleted after main reachability was
  verified.
- No deployment, promotion, traffic, DNS, database, secret, or other
  production-infrastructure mutation occurred.

## Handoff

This record hands the immutable review, validation, merge, publication, and
production-boundary evidence to the completion ledger. No unresolved
question-ledger item exists for this release.
