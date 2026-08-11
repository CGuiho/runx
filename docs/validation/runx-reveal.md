---
name: RunX Reveal Integration Validation
purpose: Preserve exact-head validation, protected merge, main reachability, and minor-release readiness for RunX pull request 50.
description: Records Go, selector, no-spawn, XDocs, Mirror, cross-build, CI, merge, and safety evidence for the reveal command.
created: 2026-08-11
flags:
  - ready
  - integrated
tags:
  - validation
  - integration
  - pull-request
  - cli
keywords:
  - runx reveal
  - pull request 50
  - exact head
  - 0.12.0
  - 11 artifacts
owner: runx-validation
---

# RunX Reveal Integration Validation

## Verdict

Ready for minor release from merged `main`. Exact head
[0dd42b86b9283e8c8ea73da67d0acad5dcc29dba](https://github.com/CGuiho/runx/commit/0dd42b86b9283e8c8ea73da67d0acad5dcc29dba)
passed independent review, local validation, CI, and protected integration.

## Immutable Validation Evidence

- Pull request: [CGuiho/runx#50](https://github.com/CGuiho/runx/pull/50).
- Validation comment: [issue comment 5250673353](https://github.com/CGuiho/runx/pull/50#issuecomment-5250673353).
- Accepted review for the same head: [review 4904190170](https://github.com/CGuiho/runx/pull/50#pullrequestreview-4904190170).
- Validated branch/base: `codex/runx-reveal` at `0dd42b86b9283e8c8ea73da67d0acad5dcc29dba`
  against `main` at `eb04196beedc2651521375e9d8080cea03f3ac11`.
- CI: [run 31472102788](https://github.com/CGuiho/runx/actions/runs/31472102788)
  completed successfully for the validated head.

## Exact-Head Checks

- `gofmt -l main.go cmd pkg embed devops` returned no paths and
  `git diff --check` passed.
- `go mod tidy` left the SHA-256 hashes of `go.mod` and `go.sum` unchanged.
- Focused reveal regressions passed for selector parity, confirmation bypass,
  invalid selectors, no-spawn behavior, catalog flags, help, aliases, and
  rejection of unsupported flags and child arguments.
- `go test -count=1 ./...`, `go vet ./...`, and `go build ./...` passed.
- Strict XDocs metadata passed for all touched scopes; `xdocs tree` was
  complete; `xdocs doctor --warnings-as-errors` reported zero errors and zero
  warnings.
- `mirror config check` passed and `mirror version plan minor` resolved the
  released `0.11.1` baseline to `@guiho/runx/v0.12.0`.
- The release builder compiled Linux AMD64/ARM64/ARMv7/ARMv6, Darwin
  AMD64/ARM64, and Windows AMD64/ARM64 with the exact head embedded.
- The release verifier confirmed exactly 11 artifacts and every SHA-256
  checksum.

## Merge Evidence

- Merge result: [cccc87348da130df90665e954f48bad9bd652ceb](https://github.com/CGuiho/runx/commit/cccc87348da130df90665e954f48bad9bd652ceb).
- `git merge-base --is-ancestor 0dd42b86b9283e8c8ea73da67d0acad5dcc29dba origin/main`
  succeeded after refreshing `main`.
- GitHub reported the PR merged and closed; the reviewed head is reachable
  from protected `main`.

## Release And Safety Boundary

- The authorized target is the Mirror-managed minor release `0.12.0`.
- The already-prepared 0.11.1 patch was verified public, stable, and archived
  before reveal merged; npm latest was `0.11.1` at the release boundary.
- Final publication still requires a clean integrated-main Mirror plan/apply,
  successful tag workflow, exact 11 public assets, checksums/digests, native
  Windows AMD64 reveal smoke, and npm latest `0.12.0`.
- Environment and key paths remained opaque; no environment/key file was read,
  decrypted, or modified.
- No deployment, promotion, traffic, DNS, database, secret, or other
  production-infrastructure mutation is part of this release.
