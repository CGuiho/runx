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

## 2026-08-09 - RunX 0.11.0 Release

- Specification: [release-runx-0.11.0.md](release-runx-0.11.0.md)
- Review: [runx-0.11.0-release-prep-review.md](../../reviews/implementation/runx-0.11.0-release-prep-review.md)
- Validation: [runx-0.11.0-release-prep.md](../../validation/runx-0.11.0-release-prep.md)
- Preparation and integration: [CGuiho/runx#41](https://github.com/CGuiho/runx/pull/41), accepted review [4890485277](https://github.com/CGuiho/runx/pull/41#pullrequestreview-4890485277), ready validation [5229738249](https://github.com/CGuiho/runx/pull/41#issuecomment-5229738249), merge commit [`4ed75572acb32d35be12d36e1e760671300f2733`](https://github.com/CGuiho/runx/commit/4ed75572acb32d35be12d36e1e760671300f2733).
- Final acceptance release: [`@guiho/runx/v0.11.0`](https://github.com/CGuiho/runx/releases/tag/%40guiho%2Frunx%2Fv0.11.0), annotated tag object `51edf2e95849ca636fb05a4a6cfe161176a02314`, peeled source [`e2b86336ebd95bc6bf25d395f518a1dce66132e1`](https://github.com/CGuiho/runx/commit/e2b86336ebd95bc6bf25d395f518a1dce66132e1).
- Workflow: [Publish 31295252704](https://github.com/CGuiho/runx/actions/runs/31295252704), successful after the authorized production-environment approval for package/source publication.
- Acceptance: release is non-draft and non-prerelease with exactly 11 assets; independent downloaded SHA-256 values matched `checksums.txt` and each GitHub digest; Windows AMD64 `--version` reported `0.11.0`, `--help` and `--help-tree` exited 0; npm latest is `0.11.0`.
- Production boundary: no deployment, promotion, traffic, DNS, database, secret, or other production-infrastructure mutation occurred.

## 2026-08-11 - RunX Reveal And Windows Git Bash Automatic Shell

- Reveal specification: [runx-reveal.md](runx-reveal.md)
- Reveal implementation: [runx-reveal-implementation.md](runx-reveal-implementation.md)
- Shell-fix specification: [windows-git-bash-auto-shell.md](windows-git-bash-auto-shell.md)
- Shell-fix implementation: [windows-git-bash-auto-shell-implementation.md](windows-git-bash-auto-shell-implementation.md)
- Review: [runx-windows-git-bash-auto-shell-review.md](../../reviews/implementation/runx-windows-git-bash-auto-shell-review.md)
- Validation: [runx-0.12.1-shell-auto-release.md](../../validation/runx-0.12.1-shell-auto-release.md)
- Integration: [PR #50](https://github.com/CGuiho/runx/pull/50), [PR #53](https://github.com/CGuiho/runx/pull/53), [PR #55](https://github.com/CGuiho/runx/pull/55), and [PR #56](https://github.com/CGuiho/runx/pull/56).
- Final releases: [RunX 0.12.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.12.0) and [RunX 0.12.1](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.12.1).
- Acceptance: exactly 11 assets per release, public checksums/native smoke, npm latest 0.12.1, exact reveal output, corrected non-secret Git Bash path probe, and closed issue [#47](https://github.com/CGuiho/runx/issues/47).
- Production boundary: source/package publication only; no application or infrastructure mutation.

## 2026-08-11 - RunX Confirmation Opt-In Skill

- Specification: [confirm-opt-in-skill.md](confirm-opt-in-skill.md)
- Implementation: [confirm-opt-in-skill-implementation.md](confirm-opt-in-skill-implementation.md)
- Review: [runx-0.11.1-confirm-opt-in-review.md](../../reviews/implementation/runx-0.11.1-confirm-opt-in-review.md)
- Validation: [runx-0.11.1-confirm-opt-in.md](../../validation/runx-0.11.1-confirm-opt-in.md)
- Integration: [CGuiho/runx#46](https://github.com/CGuiho/runx/pull/46), merge commit [37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8](https://github.com/CGuiho/runx/commit/37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8).
- Final acceptance release: [@guiho/runx/v0.11.1](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.11.1), with final evidence in [runx-0.11.1-release.md](../../validation/runx-0.11.1-release.md).
- Production boundary: only authorized source/package publication occurred; no infrastructure mutation.

## 2026-08-11 - RunX 0.11.1 Release

- Specification: [release-runx-0.11.1.md](release-runx-0.11.1.md)
- Preparation: [CGuiho/runx#49](https://github.com/CGuiho/runx/pull/49), merge commit [00789a4a8225bdae8fd4e78f418834af8d3662db](https://github.com/CGuiho/runx/commit/00789a4a8225bdae8fd4e78f418834af8d3662db).
- Validation: [runx-0.11.1-release.md](../../validation/runx-0.11.1-release.md).
- Final acceptance release: [@guiho/runx/v0.11.1](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.11.1), sourced from [c84f6432d1e4a1169a356600265d0c583515371b](https://github.com/CGuiho/runx/commit/c84f6432d1e4a1169a356600265d0c583515371b).
- Workflow: [Publish 31469968479](https://github.com/CGuiho/runx/actions/runs/31469968479), successful after the authorized package/source environment approval.
- Acceptance: non-draft/non-prerelease release with exactly 11 assets, independent checksums and GitHub digests, downloaded Windows AMD64 version/help/help-tree smoke, and npm latest `0.11.1`.
- Production boundary: no deployment, promotion, traffic, DNS, database, secret, or other production-infrastructure mutation occurred.
