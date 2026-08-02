---
name: RunX 0.10.1 Release Validation
purpose: Preserve the live audit, validation, publication, workflow, asset, checksum, npm, and native smoke evidence for releasing every RunX change after 0.10.0.
description: Records why the remaining documentation and release-evidence commits required patch 0.10.1 and the completed publication and verification gates.
created: 2026-08-03
flags:
  - completed
  - release-validation
tags:
  - validation
  - release
  - go
keywords:
  - RunX 0.10.1
  - release everything
  - Mirror
  - 11 artifacts
owner: runx-validation
---

# RunX 0.10.1 Release Validation

## Summary

Live Git and GitHub evidence proves that synchronized RunX `main` contains two
documentation and release-evidence commits after canonical 0.10.0. No pull
request merged after that release. Compatible documentation-only maintenance
requires a patch transition to 0.10.1 so every current-main change is released.

## Scope

- Source branch: `main`
- Audited pre-preparation main: `edebb5dd19049633c4303a7e74527cf786bf3005`
- Prior release tag: `@guiho/runx/v0.10.0`
- Peeled prior tag commit: `6ab7dbff2195ce660c169a48499b52e0414c00b1`
- Release tag: `@guiho/runx/v0.10.1`
- Locally validated preparation source:
  `a709dcfe8ed0902a8260a38b4209cdf66ab08862`.
- Release source: `15f17b3a01628d861ba6f83e47240304feeb0891`.
- Annotated tag object: `bf074f6a3dd0c1b30484313ebd622090e07cfaa5`.

## Live Release Audit

- The GitHub public release API reported `@guiho/runx/v0.10.0` as latest,
  published at `2026-08-02T22:28:13Z`, non-draft and non-prerelease.
- Local `main`, its upstream, and `origin/main` all resolved to clean commit
  `edebb5dd19049633c4303a7e74527cf786bf3005` after a safe fetch.
- Exactly two commits followed the peeled 0.10.0 target:
  - `6edcae40b8fd424cc059de76d27c54c6a7f04c27` - Record RunX 0.10.0 release evidence.
  - `edebb5dd19049633c4303a7e74527cf786bf3005` - Archive RunX 0.10.0 release task.
- Their changed paths are `TODO.md`, the 0.10.0 release task move and XDocs
  metadata, the done ledger, and the 0.10.0 validation report.
- The GitHub closed-pull-request API reported no pull request with a merge time
  after the 0.10.0 publication time.
- XDocs task context exposed a stale pre-archive release-task link in the
  0.10.0 validation report; it is repaired in this patch preparation.

## SemVer Decision

Use `patch`: 0.10.0 to 0.10.1. The release contains compatible documentation,
validation evidence, TODO archival, XDocs metadata, and one broken-link repair.
It adds no CLI feature, command, configuration, package, or behavior change.

## Commands Run

- Safe Git fetch, branch/upstream/status, tag, log, and diff audit - passed.
- GitHub latest-release and post-release merged-PR audit - passed.
- XDocs context selection - passed.
- Plain `mirror` bootstrap - passed idempotently with Mirror 4.0.0.
- `mirror config check` - passed.
- `mirror version current` - reported 0.10.0.
- `mirror version plan patch` - resolved 0.10.1 and exact tag
  `@guiho/runx/v0.10.1`; planned only annotated-tag creation and exact-tag push
  with `commit=false` and `exact_tag=true`.
- `gofmt -l main.go cmd pkg embed devops` - passed with no output.
- `go mod tidy` - passed without changing `go.mod` or `go.sum`.
- `go test -count=1 ./...` - passed.
- `go vet ./...` - passed.
- `go build ./...` - passed.
- Bash and PowerShell installer parsing - passed.
- `xdocs meta . --documents --strict`, `xdocs tree`, and
  `xdocs doctor .` - passed.
- Exact 0.10.1 changelog-note extraction - passed.
- `go run devops/build-binaries.go --version 0.10.1 --commit
  a709dcfe8ed0902a8260a38b4209cdf66ab08862 --build-date
  2026-08-02T23:00:05Z` - passed for all eight targets.
- `go run devops/verify-release-assets.go` - verified exactly 11 authored
  assets and every SHA-256 checksum.
- Locally built Windows AMD64 `--version` reported exactly `0.10.1`; root and
  `run` help passed the `list`, `describe`, `run`, selector, and child-argument
  contract.
- Exact-source [CI run 30771543747](https://github.com/CGuiho/runx/actions/runs/30771543747)
  passed its Ubuntu, Windows, and release-contract jobs.
- A fresh fetch, clean/synchronized branch check, remote-tag absence check,
  plain `mirror`, configuration check, current-version query, and patch plan
  all passed immediately before apply.
- `mirror version apply patch --yes` created and pushed only annotated tag
  `@guiho/runx/v0.10.1`; local and remote peeled targets both resolved to
  `15f17b3a01628d861ba6f83e47240304feeb0891`.
- Tag-triggered [Publish run 30771632490](https://github.com/CGuiho/runx/actions/runs/30771632490)
  passed every test, vet, build, release, installer, and npm step.
- GitHub Release inspection confirmed a public, non-draft, non-prerelease
  release with exactly 11 authored assets and a SHA-256 digest for each.
- npm `latest` reported `0.10.1`.
- All 11 assets were independently downloaded. All ten checksum-manifest
  payloads and all 11 GitHub digests matched; the skill ZIP contained
  `SKILL.md`, `agents/openai.yaml`, and `guiho-s-runx.xdocs.md`.
- The downloaded Windows AMD64 binary reported exactly `0.10.1`; root, list,
  describe, and run help passed.

## Manual Checks

- Inspected `.github/workflows/publish.yml`: the protected tag runs Go tests and
  vet, builds/verifies exactly 11 assets, extracts exact-version changelog notes,
  creates the GitHub Release, verifies the public asset set and exact installer,
  and publishes npm through OIDC. It performs no application deployment,
  production promotion, traffic, DNS, database, or secret mutation.

## Results

- Live unreleased-content audit: passed.
- SemVer classification: patch.
- Changelog and broken-link preparation: complete.
- Repository and XDocs validation: passed on archived preparation source
  `a709dcfe8ed0902a8260a38b4209cdf66ab08862`.
- Eight-target compilation and exact local 11-artifact verification: passed.
- Windows AMD64 local native smoke: passed.
- Mirror apply and annotated local/remote tag ancestry: passed.
- Exact-source branch CI and tag-triggered Publish: passed.
- Public release state and exact 11-asset inventory: passed.
- Independent checksum manifest and GitHub digests: 10/10 and 11/11 passed.
- Canonical skill ZIP and downloaded Windows AMD64 smoke: passed.
- npm latest: 0.10.1.
- Production boundary: no production deployment, promotion, traffic, DNS,
  database, secret, or other live production action occurred.

## Failures Or Blockers

- None.

## Skipped Checks

- Foreign-target runtime execution remains build-only where no matching
  runner or hardware executes the binary.

## Residual Risks

- No known release blocker remains. Foreign native targets retain the stated
  build-only limitation.

## Readiness

Complete. RunX 0.10.1 contains every audited post-0.10.0 change and passed the
local, CI, Mirror, publication, asset, checksum, npm, and native-smoke gates.

## References

- [Release task](../todo/done/release-runx-0.10.1.md)
- [RunX 0.10.1](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.1)
- [RunX 0.10.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.0)
- [RunX 0.10.0 validation](runx-0.10.0-release.md)
