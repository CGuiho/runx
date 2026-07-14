---
name: RunX Mirror Automatic Release Push
purpose: Record the accepted policy that Mirror pushes RunX release commits and tags automatically.
description: Defines push=true, the synchronized-main precondition, protected GitHub delivery, failure behavior, and rollback boundary.
created: 2026-07-14
flags:
  - accepted
tags:
  - decisions
  - release
  - mirror
keywords:
  - runx
  - mirror push
  - protected main
  - release tags
owner: runx-decisions
---

# RunX Mirror Automatic Release Push

## Status

Accepted by the repository owner on 2026-07-14.

## Context

RunX previously configured Mirror with `push = false`. Each release applied the
version commit and tag locally, then required separate branch delivery and tag
push commands. The owner wants RunX to match the automatic Mirror release flow
used by the related GUIHO repositories.

RunX also protects `main` and `@guiho/runx@*` tags. Automatic pushing must not
allow a release tag to reference a commit that has not reached protected
`main`.

## Decision

- Set `[git].push = true` in `mirror.config.toml`.
- Run future `mirror version apply` operations only from a clean, synchronized
  local `main` whose HEAD equals `origin/main`.
- Complete implementation changes through protected pull requests before
  applying Mirror.
- Treat `mirror version apply` as the single release mutation that writes the
  version, commits, tags, and pushes the configured refs.
- Keep package publication behind the GitHub `production` environment approval
  and npm OIDC trusted publishing.

## Alternatives Considered

- Keep `push = false`: rejected because it preserves the manual release step the
  owner explicitly asked to remove.
- Pass `--push` per release: rejected because it makes release behavior depend
  on operator memory instead of repository configuration.
- Apply Mirror on a release branch and merge later: rejected because the tag
  could publish a commit before that commit is part of protected `main`.

## Consequences

- Successful Mirror applies will immediately push release refs and can trigger
  the public Publish workflow.
- Operators must treat `mirror version apply` as an externally visible release,
  not a local preparation step.
- The synchronized-main gate becomes mandatory; failures stop before apply.
- Existing GitHub rulesets, environment review, and OIDC remain the enforcement
  layers for branch, tag, and package publication.

## Reversal Or Revisit Conditions

Return to `push = false` if Mirror cannot push through repository protections
reliably, if releases need a separate staging phase, or if automatic tag
publication creates unacceptable operational risk.

## Follow-up Work

- Update repository instructions and XDocs metadata.
- Validate resolved Mirror configuration and a read-only patch plan.
- Deliver the configuration through protected `main` without applying another
  release.

## References

- [Citty CLI Migration Plan](../plans/citty-cli-migration.md)
- [npm Trusted Publishing Decision](npm-trusted-publishing.md)
