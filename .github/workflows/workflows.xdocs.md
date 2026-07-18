---
subject: runx-github-workflows
description: CI and release workflows that validate RunX and publish native assets plus the npm package.
parent: runx-github
children: []
files:
  ci.yml: Validates tests, builds, npm bootstrap syntax, and the exact release matrix on Linux and Windows.
  publish.yml: Publishes exactly fourteen GitHub assets and the public npm package through OIDC trusted publishing on version tags.
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
and pull request. The Publish workflow runs only for `@guiho/runx@*` tags, retains the `production`
environment, publishes twelve native binaries plus two agent assets, and publishes the
public npm package through short-lived OIDC credentials without an npm token.
