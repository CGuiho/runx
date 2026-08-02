---
name: Release All Remaining RunX Changes
purpose: Define the audit, validation, publication, and verification gates for releasing every change after RunX 0.10.0.
description: Requires live release-boundary proof, a justified patch transition, full Go and XDocs validation, Mirror-managed publication, and independent 11-artifact verification.
created: 2026-08-03
flags:
  - testing
  - release
tags:
  - todo
  - release
  - mirror
keywords:
  - RunX 0.10.1
  - release everything
  - 11 artifacts
owner: runx-todo
---

# Release All Remaining RunX Changes

## Todo Index

- Task: `16. Release All Remaining RunX Changes`
- Status: testing
- Index: [TODO.md](../../TODO.md)
- Validation: [runx-0.10.1-release.md](../validation/runx-0.10.1-release.md)

## Outcome

Every commit remaining after canonical RunX 0.10.0 is contained in a verified
Mirror-managed RunX 0.10.1 release, with no duplicate or omitted content and no
production deployment.

## Scope

### In scope

- Prove the current release/tag, synchronized main, post-release commits,
  changed files, and merged-PR boundary from live Git and GitHub sources.
- Repair the archived-task link discovered in the prior validation report.
- Publish the compatible documentation and release-evidence changes as patch
  version 0.10.1.
- Update the changelog, task state, XDocs metadata, and validation evidence.
- Validate formatting, module integrity, tests, vet, build, installers, XDocs,
  all eight targets, checksums, and the exact 11-artifact set.
- Commit and push preparation, gate Mirror apply on green CI, and independently
  verify the tag, workflows, public release, npm package, checksums, skill ZIP,
  and downloaded Windows AMD64 binary.

### Out of scope

- New RunX product behavior or compatibility changes.
- Cloud deployment or promotion, Cloud Run traffic or container selection,
  DNS, production databases, secrets, or any other live production mutation.
- Manual version edits, manual tags, duplicate releases, or force pushes.

## Acceptance Signals

- `@guiho/runx/v0.10.0` is confirmed as the prior canonical release and every
  commit after its peeled target is enumerated from synchronized `main`.
- No merged pull request or remote-main commit is omitted from the audit.
- Documentation-only and release-evidence scope justifies patch 0.10.1.
- The prior validation report links to the archived task path.
- Full local Go CLI, installer, XDocs, release-note, cross-build, checksum, and
  Windows AMD64 smoke checks pass on the exact intended release source.
- Mirror plans and creates only canonical tag `@guiho/runx/v0.10.1` after the
  preparation source is clean, synchronized, and CI-green.
- Publish succeeds and the public non-draft/non-prerelease release contains
  exactly eight binaries, `guiho-s-runx.zip`, `guiho-i-runx.md`, and
  `checksums.txt`.
- All downloaded checksum-listed payloads verify, all GitHub asset digests
  match, the skill ZIP is complete, npm latest is 0.10.1, and downloaded
  Windows AMD64 reports 0.10.1 and passes help-tree smoke.
- No production action occurs.

## Dependencies And Context

- Prior release: [RunX 0.10.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.0).
- Existing release evidence: [RunX 0.10.0 validation](../validation/runx-0.10.0-release.md).
- Release behavior is owned by `mirror.yaml` and
  `.github/workflows/publish.yml`.
- Executing plan unit: release-only patch lifecycle. A separate architecture,
  product-requirements, or feature plan is unnecessary because no product or
  runtime contract changes; this specification owns the complete bounded unit.

## Watch-outs

- Compare against the peeled 0.10.0 tag target, not only its release title.
- The protected tag triggers source/package publication through the GitHub
  `production` environment. The inspected workflow does not deploy or promote
  application infrastructure.
- Cross-build success is build-only evidence for foreign targets; only matching
  runners or hardware establish runtime execution.

## Before Starting

- Read repository/SWE instructions, TODO/release evidence, Mirror/XDocs config,
  changelog, workflows, release tooling, and required release skills.
- Fetch safely and stop on dirty overlap or unsafe branch divergence.

## While Working

- Use explicit-path commits, plain pushes, and Mirror for the version tag.
- Keep the task active until remote publication and independent asset checks
  are complete.
- Record passed, failed, skipped, and build-only checks without overstating
  foreign runtime support.

## After Finishing

- Materialize final evidence in the linked validation report.
- Archive this task only after tag, workflow, release, asset, checksum, npm,
  native-smoke, ancestry, and no-production gates are verified.

## References

- [TODO.md](../../TODO.md)
- [RunX Go rewrite RFC](../rfc/runx-go-rewrite-rfc.md)
- [RunX 0.10.0 validation](../validation/runx-0.10.0-release.md)
