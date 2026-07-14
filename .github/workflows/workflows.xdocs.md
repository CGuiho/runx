---
subject: runx-github-workflows
description: CI and release workflows that validate RunX and publish native assets plus the npm package.
parent: runx-github
children: []
files:
  ci.yml: Validates RunX on Linux and Windows, including Windows self-upgrade regression coverage and the native release matrix.
  publish.yml: Publishes native binaries to GitHub Releases and the public npm package through OIDC trusted publishing on version tags.
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
environment, publishes twelve native GitHub Release assets, and publishes the
public npm package through short-lived OIDC credentials without an npm token.
