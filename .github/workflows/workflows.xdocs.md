---
subject: runx-github-workflows
description: CI and release workflows that validate RunX and publish native assets plus the npm package.
parent: runx-github
children: []
files:
  ci.yml: Validates tests, builds, npm bootstrap syntax, the exact release matrix, and the public one-line Linux installer on Linux and Windows.
  publish.yml: Creates or updates exact-version GitHub Release notes, publishes exactly fourteen assets, and publishes the public npm package through OIDC trusted publishing on version tags.
documents: {}
tags:
  - ci
  - github actions
keywords:
  - runx
  - bun
  - validation
  - trusted publishing
flags: []
status: stable
---

The CI workflow verifies RunX on Linux and Windows for every main-branch push
and pull request. Its Ubuntu job also executes the public `curl -fsSL ... | bash`
installer in an isolated home and project, then checks the published version,
Linux welcome, and both global skill files. The Publish workflow runs only for
`@guiho/runx@*` tags without an environment approval gate, publishes twelve
native binaries plus two `.md` agent assets,
creates or refreshes notes from only the matching changelog section, and
publishes the public npm package through short-lived OIDC credentials without
an npm token.
