---
name: RunX 0.10.1 Release Validation
purpose: Preserve the live audit, validation, publication, workflow, asset, checksum, npm, and native smoke evidence for releasing every RunX change after 0.10.0.
description: Records why the remaining documentation and release-evidence commits require patch 0.10.1 and the gates that must pass before completion.
created: 2026-08-03
flags:
  - testing
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
- Intended release: `@guiho/runx/v0.10.1`
- Intended release source: pending preparation commit and CI.

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
- Full repository, XDocs, release matrix, checksum, and native checks - pending.

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
- Repository validation: pending.
- Mirror apply and remote publication verification: pending.

## Failures Or Blockers

- None at preparation time.

## Skipped Checks

- Foreign-target runtime execution will remain build-only where no matching
  runner or hardware executes the binary.

## Residual Risks

- Tag creation triggers public GitHub/npm source publication after the protected
  environment gate is approved.

## Readiness

Not ready for Mirror apply until full validation passes, preparation is committed
and pushed, exact-source CI succeeds, and `main` is clean and synchronized.

## References

- [Release task](../todo/release-runx-0.10.1.md)
- [RunX 0.10.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.0)
- [RunX 0.10.0 validation](runx-0.10.0-release.md)
