---
name: RunX npm Trusted Publishing Release Plan
purpose: Sequence the workflow implementation, protected-branch delivery, Mirror patch release, and public publish verification for RunX 0.2.1.
description: Executable units for adding npm OIDC publishing without removing native releases, then triggering and observing the first trusted publication.
created: 2026-07-14
flags:
  - approved
tags:
  - plans
  - release
  - npm
keywords:
  - runx
  - trusted publishing
  - github actions
  - mirror patch
  - 0.2.1
owner: runx-plans
---

# RunX npm Trusted Publishing Release Plan

## Source Decision

- [RunX npm Trusted Publishing](../decisions/npm-trusted-publishing.md)

## Unit 1: Implement the Workflows

- Goal: Keep two GitHub workflows while adding npm OIDC publication to the
  existing tag-triggered publish job.
- Owner: RunX repository.
- Dependencies: Accepted trusted-publishing decision; npm trusted publisher for
  `CGuiho/runx`, `publish.yml`, and environment `production`.
- Files:
  - `.github/workflows/ci.yml`
  - `.github/workflows/publish.yml`
  - `.github/workflows/workflows.xdocs.md`
  - `scripts/runx-bin.ts`
  - `scripts/scripts.xdocs.md`
  - `AGENTS.md`
- Permissions: CI retains `contents: read`; Publish retains `contents: write`
  and gains `id-token: write` only for npm OIDC.
- Checks: Inspect YAML, run `git diff --check`, and validate the package with
  Bun before delivery.
- Acceptance:
  - CI validates install, typecheck, tests, build, and native builds.
  - Publish retains tag trigger `@guiho/runx@*`, environment `production`, and
    twelve-asset GitHub Release handling.
  - Publish uses Node 24 with npm 11.5.1 or newer and runs
    `npm publish --access public` without an npm token.
  - The published `runx` bin launcher prefers the included compiled
    `library/guiho-runx-bin.js` and retains the source launcher fallback only
    for repository checkouts.
  - `npm pack --dry-run --json` confirms the launcher and compiled library are
    both included, and the packaged launcher can display RunX help through Bun.
- Stop conditions: Do not continue if the repository URL, workflow filename,
  environment, OIDC permission, or public access does not match npm trusted
  publisher configuration.

## Unit 2: Validate and Deliver Through the Protected Branch

- Goal: Land the workflow and documentation changes on `main` without bypassing
  the new branch ruleset.
- Dependencies: Unit 1 complete.
- Checks:
  - `bun run typecheck`
  - `bun test`
  - `bun run build`
  - `bun run binaries`
  - `xdocs meta .github --strict --format json`
  - `xdocs tree`
  - targeted XDocs doctor checks
- Delivery:
  - Commit only understood RunX changes on a `codex/` feature branch.
  - Push the feature branch, open a pull request, and wait for required checks.
  - Merge without bypassing protections, then synchronize local `main`.
- Acceptance: Remote `main` contains the trusted-publishing workflow and its
  required documentation before a release tag is pushed.
- Stop conditions: Stop on failed CI, unresolved required review, unexpected
  worktree changes, or protection failures.

## Unit 3: Prepare and Apply Patch 0.2.1

- Goal: Use Mirror to create the `0.2.1` release commit and tag from the updated
  protected-branch history.
- Dependencies: Unit 2 merged and local `main` synchronized.
- Sequence:
  - Run `mirror config show` and confirm changelog and dirty-worktree policy.
  - Confirm a clean worktree.
  - Run typecheck, tests, build, and native binaries.
  - Run `mirror version plan patch` and confirm `nextVersion` is `0.2.1`.
  - Update `CHANGELOG.md` and its XDocs owner metadata if required.
  - Commit release preparation before applying Mirror.
  - Run `mirror version apply patch --yes` with the configured release commit.
- Acceptance: package version, Mirror release commit, and
  `@guiho/runx@0.2.1` tag agree.
- Stop conditions: Stop if Mirror plans another version, validation fails, or
  the worktree is dirty before apply.

## Unit 4: Push and Observe Trusted Publishing

- Goal: Trigger the public Publish workflow and verify npm trusted publishing.
- Dependencies: Unit 3 complete; explicit user authorization to test the public
  release path.
- Delivery:
  - Land the Mirror release commit through the protected branch when required.
  - Push `@guiho/runx@0.2.1` only after its commit is part of `main` history.
  - Monitor the public GitHub Actions run through completion.
  - Verify the GitHub Release has twelve native assets.
  - Verify npm reports `@guiho/runx@0.2.1` with trusted provenance.
- Acceptance: the public Publish run succeeds and both distribution channels
  expose version `0.2.1`.
- Stop conditions: If npm fails, capture the exact public job failure and do not
  create another version or replace the tag; correct trust configuration and
  rerun the same workflow only after confirming npm did not accept the version.

## First Executable Unit

Unit 1: implement the workflows and their XDocs/agent documentation updates.
