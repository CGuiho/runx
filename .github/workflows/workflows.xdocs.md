---
subject: runx-github-workflows
description: CI workflows that validate RunX without publishing packages or releases.
parent: runx-github
children: []
files:
  ci.yml: Installs Bun dependencies, type-checks, tests, builds, and compiles the CLI.
documents: {}
tags:
  - ci
  - github actions
keywords:
  - runx
  - bun
  - validation
flags: []
status: stable
---

The CI workflow verifies RunX on every main-branch push and pull request. It
intentionally does not publish to npm or GitHub Releases.
