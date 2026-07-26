---
subject: runx-github-workflows
description: Go-only CI and protected-tag publication workflows for the canonical 11-artifact RunX release.
parent: runx-github
children: []
files:
  ci.yml: Runs formatting, module integrity, tests, vet, native smokes, all target builds, asset verification, and Node bootstrap syntax checks.
  publish.yml: Validates Go, builds 11 assets, publishes a protected GitHub Release, verifies the exact remote set and installer, then publishes the native npm bootstrap.
documents: {}
tags:
  - github-actions
  - go
  - release
keywords:
  - CI
  - publish
  - 11 artifacts
flags: []
status: stable
---

Neither workflow invokes the legacy Bun/TypeScript domain implementation.
