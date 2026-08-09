---
name: Publish RunX 0.11.0
purpose: Preserve the complete evidence and acceptance state for the RunX 0.11.0 publication.
description: Records the minor-release decision, Mirror plan and apply, full Go and XDocs validation, canonical assets, and independent publication verification.
created: 2026-08-09
flags:
  - release
  - testing
tags:
  - todo
  - release
  - mirror
keywords:
  - RunX 0.11.0
  - interactive confirmation
  - selector identities
  - init defaults
owner: runx-todo
---

# Publish RunX 0.11.0

## Todo Index

- Task: `18. Publish RunX 0.11.0`
- Status: testing
- Index: [TODO.md](../../TODO.md)

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

- `main` is clean and synchronized at `71117e5715c7c5095429dd482e6ad4f809853b1d`.
- Prior tag `@guiho/runx/v0.10.1` is an ancestor and no `0.11.0` tag/release exists.
- Mirror config is valid and `mirror version plan minor` selects
  `@guiho/runx/v0.11.0`.
- Go formatting, tidy, tests, vet, build, native smoke, XDocs, and exact
  11-artifact verification pass.
- Publication completes with a non-draft, non-prerelease release containing
  exactly eight binaries, the skill ZIP, instruction Markdown, and checksums.
- Independent asset digests and checksum downloads match; npm latest is 0.11.0.

## Validation Record

To be completed after the tag-triggered workflow and independent download
verification. Record the final main SHA, annotated tag object and peeled source
target, workflow URLs/status, release URL/assets, npm evidence, Windows AMD64
version/help smoke, and production-boundary result here.
