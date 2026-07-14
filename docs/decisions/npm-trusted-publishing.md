---
name: RunX npm Trusted Publishing
purpose: Define the accepted GitHub Actions and Mirror release flow for publishing RunX to npm with OIDC.
description: Records the two-workflow design, native release preservation, npm trusted-publishing requirements, and protected-branch release path.
created: 2026-07-14
flags:
  - accepted
tags:
  - decisions
  - release
  - npm
keywords:
  - runx
  - trusted publishing
  - oidc
  - github actions
  - mirror
owner: runx-decisions
---

# RunX npm Trusted Publishing

## Context

RunX already has separate CI and tag-triggered publishing workflows. The
publish workflow validates the package, compiles twelve native binaries, and
publishes them to a GitHub Release. The original alpha decision excluded npm
publishing, but RunX is now ready to publish `@guiho/runx` through npm trusted
publishing without a long-lived npm token.

## Decision

- Keep exactly two workflows: `.github/workflows/ci.yml` and
  `.github/workflows/publish.yml`.
- Keep CI read-only and run it for pull requests and pushes to `main`.
- Keep the publish trigger scoped to tags matching `@guiho/runx@*`.
- Preserve the existing native-binary GitHub Release and asset-verification
  behavior.
- Add npm trusted publishing to the existing publish job after the idempotent
  GitHub Release steps.
- Grant `id-token: write` only to the publish workflow and retain the
  `production` GitHub environment.
- Use a GitHub-hosted runner, Node 24, npm 11.5.1 or newer, and
  `npm publish --access public` without an npm token. OIDC trusted publishing
  supplies authentication and provenance.
- Keep Bun as the package development, validation, build, and binary toolchain.
  Node/npm is used only for the registry publication step and does not become a
  RunX runtime dependency.

## Release Trial

The first automated npm publication trial is the Mirror-managed patch release
from `0.2.0` to `0.2.1`.

1. Land the workflow and documentation changes on `main` through a pull
   request, because the default branch is protected.
2. Run the full RunX validation and `mirror version plan patch` before applying
   the release.
3. Prepare the changelog from the planned version, commit release preparation,
   and apply the patch with Mirror so the release commit and
   `@guiho/runx@0.2.1` tag remain Mirror-managed.
4. Push the protected release tag only after the release commit is part of the
   protected-branch history.
5. Observe the Publish workflow and verify both the GitHub Release assets and
   npm package publication.

## Failure and Retry Behavior

- GitHub Release creation/upload remains before npm publication because those
  operations are idempotent and can be retried safely.
- If npm rejects OIDC configuration, the workflow fails at the npm step without
  attempting a second version or modifying the existing tag.
- After correcting npm trusted-publisher settings, rerun the failed workflow;
  the GitHub Release assets are updated in place before npm is retried.
- Never retry by republishing the same version after npm has accepted it;
  confirm registry state first.

## Validation

- Validate workflow syntax and inspect the resulting diff.
- Run `bun run typecheck`, `bun test`, `bun run build`, and
  `bun run binaries` before the release.
- Run XDocs strict metadata, tree, and doctor checks for every changed
  documentation module.
- Confirm the npm trusted publisher names repository `CGuiho/runx`, workflow
  `publish.yml`, environment `production`, and permission to run `npm publish`.

## Superseded Boundary

This decision supersedes only the alpha statement that RunX has CI without an
npm publishing job. The remaining safety, manifest, runtime, and scope
decisions in [RunX Alpha Boundaries](alpha-boundaries.md) remain in force.
