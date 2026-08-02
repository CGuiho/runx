---
name: Release All Remaining RunX Changes
purpose: Preserve the completed audit, validation, publication, and verification evidence for every RunX change after 0.10.0.
description: Records the live release boundary, justified patch transition, full Go and XDocs validation, Mirror-managed publication, and independent 11-artifact verification.
created: 2026-08-03
flags:
  - completed
  - release
tags:
  - todo
  - release
  - mirror
keywords:
  - RunX 0.10.1
  - release everything
  - 11 artifacts
owner: runx-todo-done
---

# Release All Remaining RunX Changes

## Todo Index

- Task: `16. Release All Remaining RunX Changes`
- Status: completed
- Index: [TODO.md](../../../TODO.md)
- Validation: [runx-0.10.1-release.md](../../validation/runx-0.10.1-release.md)

## Outcome

Every commit remaining after canonical RunX 0.10.0 was published in verified,
Mirror-managed RunX 0.10.1, with no duplicate or omitted content and no
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

- `@guiho/runx/v0.10.0` was confirmed as the prior canonical release and every
  commit after its peeled target was enumerated from synchronized `main`.
- GitHub reported no pull request merged after the prior release.
- Documentation-only and release-evidence scope justified patch 0.10.1.
- The prior validation report links to the archived task path.
- Full local Go CLI, installer, XDocs, release-note, cross-build, checksum, and
  Windows AMD64 smoke checks passed.
- Mirror created only canonical annotated tag `@guiho/runx/v0.10.1` after the
  preparation source was synchronized and CI-green.
- Publish succeeded and the public non-draft/non-prerelease release contains
  exactly eight binaries, `guiho-s-runx.zip`, `guiho-i-runx.md`, and
  `checksums.txt`.
- All downloaded checksum-listed payloads and GitHub asset digests verified,
  the skill ZIP was complete, npm latest was 0.10.1, and downloaded Windows
  AMD64 reported 0.10.1 and passed help-tree smoke.
- No production action occurred.

## Completion Evidence

- Release: [RunX 0.10.1](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.1).
- Release source: `15f17b3a01628d861ba6f83e47240304feeb0891`.
- Annotated tag object: `bf074f6a3dd0c1b30484313ebd622090e07cfaa5`.
- Branch CI: [run 30771543747](https://github.com/CGuiho/runx/actions/runs/30771543747), successful.
- Publication: [run 30771632490](https://github.com/CGuiho/runx/actions/runs/30771632490), successful.
- Public release: non-draft, non-prerelease, exactly 11 authored assets.
- Independent download: all ten checksum-listed payloads and all 11 GitHub
  asset digests matched; the skill ZIP had its three canonical files; Windows
  AMD64 reported 0.10.1 and passed root/list/describe/run help; npm latest
  reported 0.10.1.
- Production boundary: no production action occurred.

## Dependencies And Context

- Prior release: [RunX 0.10.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.0).
- Existing release evidence: [RunX 0.10.0 validation](../../validation/runx-0.10.0-release.md).
- Release behavior is owned by `mirror.yaml` and
  `.github/workflows/publish.yml`.
- Executed plan unit: release-only patch lifecycle. A separate architecture,
  product-requirements, or feature plan was unnecessary because no product or
  runtime contract changed; this specification owned the complete bounded unit.

## Watch-outs

- The audit compared against the peeled 0.10.0 tag target, not only its release
  title.
- The protected tag used the GitHub `production` environment solely as a
  source/package publication gate. The inspected workflow did not deploy or
  promote application infrastructure.
- Cross-build success is build-only evidence for foreign targets; only matching
  runners or hardware establish runtime execution.

## References

- [TODO.md](../../../TODO.md)
- [RunX Go rewrite RFC](../../rfc/runx-go-rewrite-rfc.md)
- [Release validation](../../validation/runx-0.10.1-release.md)
