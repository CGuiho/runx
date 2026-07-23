---
subject: runx-github-workflows
description: CI and release workflows that validate RunX and publish native assets plus the npm package.
parent: runx-github
children: []
files:
  ci.yml: Validates tests, builds, npm bootstrap syntax, the exact release matrix, and a self-consistent latest public Linux installation without coupling unreleased source to the public latest version.
  publish.yml: Creates or updates exact-version GitHub Release notes, publishes exactly fourteen assets, accepts the exact tagged public installer, and publishes the public npm package through OIDC trusted publishing.
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
and pull request. Its Ubuntu job executes the public `curl -fsSL ... | bash`
installer in an isolated home and project, then verifies that the installed
latest version, startup output, and both global skill files agree. CI never
requires that public latest equal an unreleased or concurrently publishing
`package.json` version.

The Publish workflow runs only for `@guiho/runx@*` tags. It publishes twelve
native binaries plus two `.md` agent assets, verifies the exact asset set, and
then executes the tagged public installer with `--version` for that release.
Only after exact-version public installation succeeds may npm OIDC publication
continue. Release notes contain only the matching changelog section.
