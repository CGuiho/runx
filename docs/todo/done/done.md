---
name: RunX Completed Tasks
purpose: Preserve the permanent completion ledger for archived RunX work.
description: Links completed RunX task specifications to decisions, plans, reviews, validation, releases, and closed issues.
created: 2026-07-23
flags:
  - completed
owner: runx-todo-done
tags:
  - todo
  - done
keywords:
  - runx
  - completion ledger
---

# RunX Completed Tasks

## 2026-07-23 - Manifest V2 Composition

- Specification: [manifest-v2-composition.md](manifest-v2-composition.md)
- Requirements: [manifest-v2-composition.md](../../requirements/manifest-v2-composition.md)
- Decision: [manifest-v2-composition.md](../../decisions/manifest-v2-composition.md)
- Plan: [manifest-v2-composition.md](../../plans/manifest-v2-composition.md)
- Review: [manifest-v2-composition-review.md](../../reviews/implementation/manifest-v2-composition-review.md)
- Validation: [manifest-v2-composition.md](../../validation/manifest-v2-composition.md)
- Final acceptance release: [RunX 0.7.2](https://github.com/CGuiho/runx/releases/tag/%40guiho%2Frunx%400.7.2)
- Issue: [CGuiho/runx#26](https://github.com/CGuiho/runx/issues/26), closed.

## 2026-08-03 - RunX 0.10.0 Release

- Specification: [release-runx-0.10.0.md](release-runx-0.10.0.md)
- Validation: [runx-0.10.0-release.md](../../validation/runx-0.10.0-release.md)
- Accepted pull requests:
  - [CGuiho/runx#29](https://github.com/CGuiho/runx/pull/29) - Fix Windows installer PATH setup.
  - [CGuiho/runx#31](https://github.com/CGuiho/runx/pull/31) - Align list output and resolve numeric indexes.
- Final acceptance release: [RunX 0.10.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.0)
- Workflows: branch [CI 30769897784](https://github.com/CGuiho/runx/actions/runs/30769897784) and tag-triggered [Publish 30770011393](https://github.com/CGuiho/runx/actions/runs/30770011393), successful.
- Acceptance: exactly 11 authored assets, all checksums and GitHub digests verified, downloaded Windows AMD64 smoke passed, npm latest 0.10.0, and no production action occurred.

## 2026-08-03 - RunX 0.10.1 Release

- Specification: [release-runx-0.10.1.md](release-runx-0.10.1.md)
- Validation: [runx-0.10.1-release.md](../../validation/runx-0.10.1-release.md)
- Unreleased boundary: two post-0.10.0 documentation and release-evidence commits; no pull request merged after the prior release.
- Final acceptance release: [RunX 0.10.1](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.1)
- Workflows: branch [CI 30771543747](https://github.com/CGuiho/runx/actions/runs/30771543747) and tag-triggered [Publish 30771632490](https://github.com/CGuiho/runx/actions/runs/30771632490), successful.
- Acceptance: exactly 11 authored assets, all ten checksum-listed payloads and all 11 GitHub digests verified, canonical skill ZIP verified, downloaded Windows AMD64 smoke passed, npm latest 0.10.1, and no production action occurred.
