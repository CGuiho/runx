---
name: RunX npm Trusted Publishing 0.2.2 Validation
purpose: Record the release checks and GitHub result for the 0.2.2 trusted-publishing retry.
description: Documents successful build validation and the tag-ruleset rejection that prevented the Publish workflow from starting.
created: 2026-07-14
flags: []
tags:
  - validation
  - release
keywords:
  - runx
  - npm
  - trusted publishing
owner: runx-validation
---

# npm Trusted Publishing 0.2.2 Validation

## Result

The `0.2.2` package bump was merged to protected `main`, but GitHub rejected
creation of the `@guiho/runx@0.2.2` tag. Because the tag was not created, the
Publish workflow did not run and npm remained on `0.2.0`.

## Passed checks

- `bun run typecheck`
- `bun test`: 4 passed, 0 failed
- `bun run build`
- `bun run binary`
- `bun run binaries`: 12 native assets built and verified
- `mirror version plan patch`: planned `0.2.1` to `0.2.2`
- `xdocs meta . --documents --strict --owner runx --format json`
- Pull request `#4` CI: passed

## Release evidence

- Protected-main merge: `57c7d81`
- Mirror release commit: `6c67be5`
- Local annotated tag: `@guiho/runx@0.2.2`
- Remote tag: not created
- Publish workflow run: not started
- npm latest version after the attempt: `0.2.0`

## Blocking rules

GitHub rejected the tag with `GH013` because the active tag ruleset both
requires linear history and restricts tag creation. The reported merge-commit
violation was `f445631`. The ruleset had no bypass actors, and the authenticated
repository owner could not bypass it.

Trusted publishing itself was not exercised by this attempt. The tag ruleset
must permit the release actor to create the version tag before npm OIDC can be
validated.
