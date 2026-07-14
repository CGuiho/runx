---
name: RunX Mirror Automatic Push Plan
purpose: Execute the accepted push=true release policy through protected main without creating another release.
description: Sequences configuration, safety instructions, XDocs validation, protected delivery, and post-merge verification.
created: 2026-07-14
flags:
  - approved
tags:
  - plans
  - release
  - mirror
keywords:
  - runx
  - mirror push
  - protected main
  - configuration
owner: runx-plans
---

# RunX Mirror Automatic Push Plan

## Source Decision

- [RunX Mirror Automatic Release Push](../decisions/mirror-automatic-push.md)

## Unit 1: Configure And Document Automatic Push

- Goal: Make automatic release pushing a repository-owned Mirror policy.
- Owner: RunX repository.
- Dependencies: Accepted automatic-push decision; `0.2.4` already applied and
  pushed, so this change must not apply another version.
- Files:
  - `mirror.config.toml`
  - `AGENTS.md`
  - `runx.xdocs.md`
- Data, auth, and cache: No application data, authentication, authorization, or
  cache impact. GitHub credentials and rulesets remain external.
- Changes:
  - Set `[git].push = true`.
  - Require clean synchronized `main` before `mirror version apply`.
  - Document that apply is an externally visible release mutation.
- Checks:
  - `mirror config check`
  - `mirror config show`
  - `mirror version plan patch --format json`
  - Verify the plan reports `pushEnabled: true` and next patch `0.2.5` without
    applying it.
  - Strict root XDocs metadata and doctor checks.
- Acceptance: The resolved config enables push; safety guidance prevents apply
  from a feature branch or stale main; no version, commit, or tag is created by
  validation.
- Stop conditions: Stop if Mirror cannot load the config, the read-only plan
  mutates state, or the worktree contains unrelated changes.

## Unit 2: Deliver Through Protected Main

- Goal: Merge the configuration policy without bypassing required CI.
- Dependencies: Unit 1 complete and committed.
- Delivery:
  - Push `codex/mirror-auto-push`.
  - Open a ready pull request to `main`.
  - Wait for required CI and merge through the repository ruleset.
  - Synchronize local `main` and confirm `mirror config show` reports
    `push: true`.
- Acceptance: `origin/main` owns the configuration and instructions; no new
  release tag was created.
- Stop conditions: Stop on CI failure, unresolved protection failure, or an
  unexpected version/tag mutation.

## TODO Alignment

This is a bounded owner-approved configuration change completed in one session;
it does not require a separate TODO entry.

## First Executable Unit

Unit 1: set and document the resolved Mirror automatic-push policy.
