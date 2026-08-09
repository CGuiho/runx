---
name: Publish RunX 0.11.0
purpose: Preserve the complete evidence and acceptance state for the RunX 0.11.0 publication.
description: Records the minor-release decision, Mirror plan and apply, full Go and XDocs validation, canonical assets, and independent publication verification.
created: 2026-08-09
flags:
  - release
  - completed
tags:
  - todo
  - done
  - release
  - mirror
keywords:
  - RunX 0.11.0
  - interactive confirmation
  - selector identities
  - init defaults
owner: runx-todo-done
---

# Publish RunX 0.11.0

## Todo Index

- Task: `18. Publish RunX 0.11.0`
- Status: completed
- Index: [TODO.md](../../../TODO.md)

## Decision

Publish a minor release, `0.11.0`, rather than a patch release. The unreleased
range after `0.10.1` includes the externally visible interactive terminal
confirmation capability for `confirm: always` commands. The compatible UID/
selector-resolution and `.scripts` initialization corrections from PR #40 are
included in the same release boundary.

## Scope

### In scope

- Move all unreleased 0.11.0 notes into one dated changelog section and leave
  an explicit empty Unreleased section.
- Validate the clean Go/Cobra source, XDocs metadata, Mirror configuration, and
  canonical 11-artifact release matrix.
- Apply only `mirror version apply minor --yes`, allowing Mirror to create and
  push the canonical annotated tag and its exact source ref.
- Verify the tag, tag-triggered publication workflow, GitHub release, checksums,
  downloaded Windows AMD64 smoke, and npm latest version.
- Archive this task with durable release and validation evidence.

### Out of scope

- Manual version-field edits, manual tag creation, duplicate releases, force
  pushes, or bypassing protected refs.
- Deployment, promotion, traffic, DNS, databases, secrets, or other
  infrastructure mutations.

## Acceptance Signals

- Release source `main` was clean and synchronized at
  [`e2b86336ebd95bc6bf25d395f518a1dce66132e1`](https://github.com/CGuiho/runx/commit/e2b86336ebd95bc6bf25d395f518a1dce66132e1) before the tag apply; the archived documentation commits follow that source.
- Prior tag `@guiho/runx/v0.10.1` is an ancestor and Mirror config is valid;
  `mirror version plan minor` selected `@guiho/runx/v0.11.0`.
- Go formatting, tidy, tests, vet, build, native smoke, XDocs, and exact
  11-artifact verification passed.
- The annotated tag object is `51edf2e95849ca636fb05a4a6cfe161176a02314`
  and peels to the release source commit above.
- Publication workflow [31295252704](https://github.com/CGuiho/runx/actions/runs/31295252704) succeeded after the authorized package/source environment approval.
- The [GitHub release](https://github.com/CGuiho/runx/releases/tag/%40guiho%2Frunx%2Fv0.11.0) is non-draft and non-prerelease with exactly 11 assets: eight binaries, the skill ZIP, instruction Markdown, and checksums.
- Independent downloaded SHA-256 values match `checksums.txt` and every GitHub asset digest; downloaded Windows AMD64 `--version` reported `0.11.0`, `--help` and `--help-tree` exited 0; npm latest is `0.11.0`.
- No deployment, promotion, traffic, DNS, database, secret, or other
  production-infrastructure mutation occurred.

## Validation Record

Preparation integration is complete at merge commit
[`4ed75572acb32d35be12d36e1e760671300f2733`](https://github.com/CGuiho/runx/commit/4ed75572acb32d35be12d36e1e760671300f2733).
The immutable implementation review is recorded at
[docs/reviews/implementation/runx-0.11.0-release-prep-review.md](../../reviews/implementation/runx-0.11.0-release-prep-review.md)
and the ready validation at
[docs/validation/runx-0.11.0-release-prep.md](../../validation/runx-0.11.0-release-prep.md).

The completion ledger records the release URL, workflow, tag object and peeled
source target, asset/checksum verification, npm evidence, and native smoke.
