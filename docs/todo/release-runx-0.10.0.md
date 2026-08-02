---
name: Release Unreleased RunX Pull Requests
purpose: Define the evidence and publication gates for the first RunX release containing merged pull requests 29 and 31.
description: Requires live ancestry proof, a justified semantic version, full Go and XDocs validation, Mirror-managed publication, and independent 11-artifact verification.
created: 2026-08-03
flags:
  - testing
  - release
tags:
  - todo
  - release
  - mirror
keywords:
  - RunX 0.10.0
  - pull request 29
  - pull request 31
  - 11 artifacts
owner: runx-todo
---

# Release Unreleased RunX Pull Requests

## Todo Index

- Task: `15. Release Unreleased RunX Pull Requests`
- Status: testing
- Index: [TODO.md](../../TODO.md)
- Validation: [runx-0.10.0-release.md](../validation/runx-0.10.0-release.md)

## Outcome

The accepted Windows installer PATH repair and numeric-index selector feature
are published in the first canonical RunX release that contains their merge
commits, with no duplicate release and no production deployment.

## Scope

### In scope

- Prove the current GitHub Release, tag target, merged pull requests, merge
  commits, and main ancestry from live sources.
- Select the SemVer transition from the actual merged behavior.
- Update the changelog and durable release evidence.
- Validate formatting, module integrity, tests, vet, build, installers, XDocs,
  all eight targets, checksums, and the exact 11-artifact set.
- Commit and push release preparation, inspect Mirror configuration and the
  complete plan, and apply the authorized Mirror transition.
- Verify branch and tag ancestry, GitHub Actions, GitHub Release metadata,
  checksums, and a downloaded native Windows AMD64 smoke.

### Out of scope

- Cloud deployment or promotion, Cloud Run traffic or container selection,
  DNS, production databases, secrets, or any other live production mutation.
- Manual version edits, manual tags, duplicate releases, or force pushes.

## Acceptance Signals

- Pull requests 29 and 31 are confirmed merged, and both merge commits are
  absent from `@guiho/runx/v0.9.0` but reachable from `origin/main`.
- The backward-compatible numeric selector capability justifies the minor
  transition from 0.9.0 to 0.10.0; the installer repair is compatible fix scope.
- The repository is clean and synchronized before Mirror applies the release.
- The full Go CLI validation and XDocs checks pass on the exact release source.
- Mirror creates and pushes the canonical `@guiho/runx/v0.10.0` tag only after
  its plan and tag-triggered workflow effects are inspected.
- The public GitHub Release is non-draft and non-prerelease and contains exactly
  eight binaries, `guiho-s-runx.zip`, `guiho-i-runx.md`, and `checksums.txt`.
- All checksum-listed payloads verify, and the downloaded Windows AMD64 binary
  reports 0.10.0 and passes native help/list smoke coverage.
- No production action occurs.

## Dependencies And Context

- Prior canonical release: [RunX 0.9.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.9.0).
- Merged installer fix: [PR #29](https://github.com/CGuiho/runx/pull/29),
  merge commit `cca0d564a95b33ca8238d931cea1757da72620c5`.
- Merged list/index feature: [PR #31](https://github.com/CGuiho/runx/pull/31),
  merge commit `17f8f086d66411b4184a04ea694eb3ff6845a9c7`.
- Release behavior is owned by `mirror.yaml` and
  `.github/workflows/publish.yml`.

## Watch-outs

- The 0.9.0 annotated tag resolves to commit
  `6af31f1f03fdad2a93d258a969671e44aab9aa31`; compare merge commits to the
  peeled tag commit, not to a stale local branch or a release title.
- The protected tag triggers source/package publication through the GitHub
  `production` environment. Its inspected workflow does not deploy or promote
  production infrastructure.
- Cross-build success is build-only evidence for foreign targets; only matching
  runners or hardware establish runtime execution.

## Before Starting

- Read repository and parent instructions, TODO/specs, release decisions,
  `mirror.yaml`, changelog, workflows, release tooling, and required skills.
- Fetch safely, preserve existing work, and stop on dirty overlap or unsafe
  divergence.

## While Working

- Use explicit-path commits, plain pushes, and Mirror for the version commit
  and tag.
- Keep the task in testing until remote publication and asset verification are
  complete.
- Record passed, failed, skipped, and build-only checks without overstating
  foreign runtime support.

## After Finishing

- Materialize final release evidence in the linked validation report.
- Archive this task only after the release, workflows, assets, checksums,
  native smoke, branch/tag ancestry, and no-production boundary are verified.

## References

- [TODO.md](../../TODO.md)
- [RunX Go rewrite RFC](../rfc/runx-go-rewrite-rfc.md)
- [Release validation](../validation/runx-0.10.0-release.md)
