---
name: Publish RunX 0.11.1
purpose: Preserve the patch-release scope, acceptance signals, Mirror decision, and protected publication boundary for RunX 0.11.1.
description: Defines the authorized patch release for the confirmation-policy skill correction after PR #46 integration.
created: 2026-08-11
flags:
  - testing
tags:
  - todo
  - release
  - mirror
  - runx
keywords:
  - runx
  - 0.11.1
  - confirm
  - patch
  - mirror
owner: runx-todo
---

# Publish RunX 0.11.1

## Todo Index

- Task: `20. Publish RunX 0.11.1`
- Status: testing
- Index: [TODO.md](../../TODO.md)
- Parent TODO: [GUIHO root TODO](../../../guiho/TODO.md)
- Implementation review: [runx-0.11.1-confirm-opt-in-review.md](../reviews/implementation/runx-0.11.1-confirm-opt-in-review.md)
- Integration validation: [runx-0.11.1-confirm-opt-in.md](../validation/runx-0.11.1-confirm-opt-in.md)

## Scope

Publish the compatible patch correction from merged PR #46 as the
Mirror-managed `@guiho/runx/v0.11.1` release.

### In scope

- Finalize `CHANGELOG.md` with an empty `Unreleased` section and the dated
  `## 0.11.1 - 2026-08-11` section.
- Run the clean-clone Go, XDocs, Mirror, exact 11-asset, checksum, native
  Windows AMD64, and npm latest validation gates.
- Apply only `mirror version apply patch --yes` after the release-prep commit
  is on protected `main`.
- Verify the annotated tag, tag/source/main ancestry, successful tag-only
  publish workflow, non-draft/non-prerelease GitHub Release, exact 11 assets,
  independent checksums and GitHub digests, Windows AMD64 smoke, and npm
  `@guiho/runx@latest` version `0.11.1`.
- Archive this release task and task 19 only after final publication evidence.

### Out of scope

- Runtime implementation changes after merged PR #46.
- Production deployment, promotion, traffic, DNS, database, or secret changes.
- Manual version edits or manually-created tags.

## Mirror Decision

This is a compatible agent-policy and documentation correction. The correct
transition is a patch from `0.11.0` to `0.11.1`, using the configured
canonical tag `@guiho/runx/v0.11.1`.

## Acceptance

- Protected release-prep PR is merged and main is clean/synchronized.
- `mirror version plan patch` reports exactly `0.11.1`.
- Only `mirror version apply patch --yes` creates the version commit and tag.
- Publish workflow succeeds and publishes source/package artifacts only.
- Final release and npm evidence are independently recorded before task archive.
