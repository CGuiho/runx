---
name: RunX npm Trusted Publishing Release Plan Review
purpose: Verify that the trusted-publishing workflow and patch-release plan is safe, sequenced, testable, and ready for execution.
description: Reviews protected-branch delivery, Mirror ordering, npm OIDC permissions, failure handling, documentation, and validation gates.
created: 2026-07-14
flags:
  - approved
tags:
  - reviews
  - plans
  - release
keywords:
  - runx
  - trusted publishing
  - plan review
  - mirror
  - github actions
owner: runx-plan-reviews
---

# RunX npm Trusted Publishing Release Plan Review

## Verdict

Ready for execution.

## Findings

- No blocker or high-severity findings remain.
- A package inspection found that the npm launcher referenced an excluded
  `source/` fallback. Unit 1 now includes the bounded launcher correction and
  pack/launcher verification, resolving the issue before publication.
- The plan traces to the accepted npm trusted-publishing decision and preserves
  Bun as the RunX runtime and development toolchain.
- OIDC permission is limited to the publish workflow, and the plan explicitly
  requires the configured `production` environment.
- Native GitHub Release publication precedes npm publication, making a failed
  OIDC attempt safely rerunnable before npm accepts the version.

## Sequencing Review

- Workflow changes land on protected `main` before any release tag is pushed.
- Mirror plans the patch before changelog preparation and applies the version
  only after validation and a clean worktree.
- The release branch is merged before its existing Mirror tag is pushed, so the
  tagged release commit is already part of protected `main` history when the
  public workflow starts.

## Acceptance and Validation Review

- Workflow contents, permissions, trigger, environment, native asset count,
  npm command, and trusted-publisher identity are explicit.
- Bun typecheck, tests, build, native matrix, XDocs checks, GitHub Actions run,
  GitHub Release assets, and npm registry state all have verification steps.
- Failure stop conditions prevent version replacement, duplicate publication,
  or a second release when the first npm outcome is uncertain.

## TODO Alignment

The repository TODO currently contains completed package and protection tasks.
This release is a single approved execution session with a durable decision and
plan, so no additional long-running TODO task is required.

## First Executable Unit

Implement Unit 1 in
[RunX npm Trusted Publishing Release Plan](../../plans/npm-trusted-publishing-release.md):
update the two workflows and their XDocs/agent documentation.

## References

- [RunX npm Trusted Publishing Release Plan](../../plans/npm-trusted-publishing-release.md)
- [RunX npm Trusted Publishing](../../decisions/npm-trusted-publishing.md)
