---
name: RunX 0.10.0 Release Validation
purpose: Preserve live audit, validation, publication, workflow, asset, checksum, and native smoke evidence for RunX 0.10.0.
description: Records why a minor release is required for merged pull requests 29 and 31 and the gates that must pass before completion.
created: 2026-08-03
flags:
  - completed
  - release-validation
tags:
  - validation
  - release
  - go
keywords:
  - RunX 0.10.0
  - pull request 29
  - pull request 31
  - Mirror
  - 11 artifacts
owner: runx-validation
---

# RunX 0.10.0 Release Validation

## Summary

Live Git and GitHub evidence proves that merged pull requests 29 and 31 are not
contained in the canonical RunX 0.9.0 release. PR 31 adds backward-compatible
numeric index selection and therefore requires a minor transition to 0.10.0;
PR 29 is a compatible installer repair included in the same release.

## Scope

- Source branch: `main`
- Audited main commit: `17f8f086d66411b4184a04ea694eb3ff6845a9c7`
- Prior release tag: `@guiho/runx/v0.9.0`
- Peeled prior tag commit: `6af31f1f03fdad2a93d258a969671e44aab9aa31`
- Intended release: `@guiho/runx/v0.10.0`
- Release source commit: `6ab7dbff2195ce660c169a48499b52e0414c00b1`
- Annotated tag object: `227a5661ac1f1df93fafb206add75ca8bdec7324`

## Live Release Audit

- The GitHub public release API reported `@guiho/runx/v0.9.0` as the latest
  release, published at `2026-07-26T00:14:35Z`, non-draft, non-prerelease, with
  exactly 11 assets.
- [PR #29 - Fix Windows installer PATH setup](https://github.com/CGuiho/runx/pull/29)
  was merged at `2026-07-30T11:50:09Z` as
  `cca0d564a95b33ca8238d931cea1757da72620c5`.
- [PR #31 - Align list output and resolve numeric indexes](https://github.com/CGuiho/runx/pull/31)
  was merged at `2026-07-30T12:16:14Z` as
  `17f8f086d66411b4184a04ea694eb3ff6845a9c7`.
- `git merge-base --is-ancestor <merge> <0.9.0-tag-commit>` returned 1 for both
  merge commits, proving neither is contained in 0.9.0.
- The same ancestry check against `origin/main` returned 0 for both commits.
- The clean local `main` fast-forwarded from
  `b87ae967d5d01693bd226503be7d1272a0f6d972` to
  `17f8f086d66411b4184a04ea694eb3ff6845a9c7` with no local divergence.

## SemVer Decision

Use `minor`: 0.9.0 to 0.10.0. Numeric `IDX` values are a newly supported
interactive selector across `describe` and `run`, with exact identity lookup
remaining authoritative. This is a backward-compatible user-visible feature.
The Windows PATH work is a backward-compatible fix and does not increase the
bump beyond minor.

## Commands Run

- `gofmt -l main.go cmd pkg embed devops` - passed with no reported files.
- `go mod tidy` with before/after SHA-256 comparison of `go.mod` and `go.sum` -
  passed with no module-graph changes in a clean Git snapshot.
- `go test -count=1 ./...` - passed for `cmd`, `devops`, `pkg/executor`,
  `pkg/maintenance`, `pkg/manifest`, `pkg/update`, and `pkg/updater`; root and
  `embed` correctly reported no test files.
- `go vet ./...` - passed.
- `go build ./...` - passed.
- Git Bash `bash -n devops/install.sh` - passed.
- PowerShell parser validation of `devops/install.ps1` - passed with zero
  syntax errors.
- `xdocs meta . --documents --strict` - passed for 29 descriptors and all
  listed companion documents.
- `xdocs tree` - passed with the complete RunX hierarchy.
- `xdocs doctor .` - passed with 0 errors and 0 warnings.
- `go run devops/extract-release-notes.go --version 0.10.0` - passed and
  extracted only the new changelog section.
- `go run devops/build-binaries.go --version 0.10.0 ...` - passed for all
  eight pure-Go targets and produced the three supporting artifacts.
- `go run devops/verify-release-assets.go --directory bundle` - passed with
  exactly 11 assets and all ten checksum entries verified.
- Locally built `runx-windows-amd64.exe --version` and
  `--help-tree-depth 2 --help-tree` - passed and reported 0.10.0 with the live
  numeric-index command signatures.
- `git diff --check` - passed.

The Go checks used a clean `git archive` snapshot of audited main because the
saved checkout intentionally retains legacy `node_modules`, which Go 1.26
would otherwise traverse during `go mod tidy`. A workspace-local writable Go
build cache avoided the sandboxed user-cache permission failure. Vet, build,
and cross-builds returned exit 0; the toolchain emitted only non-failing
read-only module-cache timestamp warnings.

## Manual Checks

- Inspected `.github/workflows/publish.yml`: the protected tag builds and
  verifies Go assets, creates/updates a GitHub Release, runs an exact-version
  installer, and publishes the npm bootstrap. It performs no cloud deployment,
  production promotion, traffic, DNS, database, or secret mutation.
- Ran plain `mirror`: Mirror 4.0.0 refreshed its global skill and idempotently
  reconciled the bounded managed block in `AGENTS.md`.
- `mirror config check` loaded `C:\GUIHO\runx\mirror.yaml` successfully.
- `mirror version current` reported `0.9.0`.
- `mirror version plan minor` resolved `0.10.0` and the exact tag
  `@guiho/runx/v0.10.0`. Its only planned effects are annotated-tag creation
  and an exact-tag push (`commit=false`, `exact_tag=true`); it plans no version
  commit or additional file mutation.
- Branch CI [run 30769771102](https://github.com/CGuiho/runx/actions/runs/30769771102)
  passed for the first release-evidence push. Final branch CI
  [run 30769897784](https://github.com/CGuiho/runx/actions/runs/30769897784)
  passed for exact release source `6ab7dbf`; its Ubuntu, Windows, and
  `release-contract` jobs all completed successfully.
- Mirror applied the inspected minor plan and created annotated tag object
  `227a5661ac1f1df93fafb206add75ca8bdec7324`, peeled to exact release source
  `6ab7dbff2195ce660c169a48499b52e0414c00b1`. The same object and peeled commit
  were independently observed on `origin`, and the release commit is reachable
  from `origin/main`.
- The protected GitHub `production` environment required its configured
  reviewer approval. The approval started only the already-inspected source
  publication job; the workflow contains no application deployment or live
  production mutation.
- Publish [run 30770011393](https://github.com/CGuiho/runx/actions/runs/30770011393)
  passed on exact tag `@guiho/runx/v0.10.0`. Tests, vet, the canonical matrix,
  scoped notes, GitHub Release creation, exact asset-set verification,
  exact-version installer smoke, and npm OIDC publication all passed.
- The public [RunX 0.10.0 release](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.0)
  was published at `2026-08-02T22:28:13Z` as non-draft and non-prerelease.
  It contains exactly 11 authored assets: eight binaries,
  `guiho-s-runx.zip`, `guiho-i-runx.md`, and `checksums.txt`.
- Downloaded all 11 public assets. The checksum file contained exactly ten
  payload entries, every payload matched its listed SHA-256, and all 11 files
  also matched GitHub's independent asset digests.
- The downloaded skill archive contained exactly
  `guiho-s-runx/SKILL.md`, `guiho-s-runx/agents/openai.yaml`, and
  `guiho-s-runx/guiho-s-runx.xdocs.md`.
- The downloaded Windows AMD64 executable returned exit 0 and version `0.10.0`;
  its depth-two help tree returned exit 0 and showed `list`,
  `describe <uid-or-selector-or-index>`, and
  `run [options] <uid-or-selector-or-index> [--] [child arguments...]`.
- The public npm `@guiho/runx` latest endpoint reported `0.10.0` after the
  Publish workflow completed.

## Results

- Live unreleased-content audit: passed.
- SemVer classification: minor.
- Repository validation: passed.
- Exact local release matrix and checksums: passed.
- Native Windows AMD64 smoke: passed.
- Mirror bootstrap, configuration, plan, and apply: passed.
- Branch CI and tag-triggered Publish workflow: passed.
- GitHub Release metadata and exact 11-asset set: passed.
- Downloaded checksum and GitHub digest verification: passed.
- Downloaded native Windows AMD64 smoke: passed.
- npm 0.10.0 publication: passed.
- Production-boundary verification: passed; no production action occurred.

## Failures Or Blockers

- The first bulk public-asset download exceeded its observation window and left
  one incomplete temporary file. A bounded resume fetched only missing or
  size-mismatched temporary targets. A concurrent completion briefly caused a
  Windows sharing violation during that resume, but the final independent
  file-size, checksum, and GitHub-digest checks all passed. This affected only
  the disposable verification directory, not any published asset or repository
  file.

## Skipped Checks

- Linux AMD64, Linux ARM64, Linux ARMv7, Linux ARMv6, Darwin AMD64, Darwin
  ARM64, and Windows ARM64 were cross-compiled successfully but not executed in
  this Windows validation environment; they remain build-only until matching
  runners or hardware provide runtime evidence.

## Residual Risks

- Foreign targets remain build-only in this Windows audit unless the matching
  workflow runner or target hardware executed them. The workflow's native
  exact-version installer smoke covered Linux AMD64; this audit independently
  executed Windows AMD64.

## Readiness

Complete. The accepted pull requests are contained in canonical RunX 0.10.0,
all repository, workflow, publication, asset, checksum, and supported native
smoke gates passed, and no production action occurred.

## References

- [Release task](../todo/done/release-runx-0.10.0.md)
- [RunX 0.9.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.9.0)
- [PR #29](https://github.com/CGuiho/runx/pull/29)
- [PR #31](https://github.com/CGuiho/runx/pull/31)
