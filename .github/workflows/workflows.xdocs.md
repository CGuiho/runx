---
subject: runx-github-workflows
description: Go-only CI and protected-tag publication workflows for the protocol-v1 RunX release and native installer gates.
parent: runx-github
children: []
files:
  ci.yml: Runs formatting, module integrity, tests, vet, native smokes, all target builds, and exact protocol-v1 asset verification on Linux and Windows.
  publish.yml: Blocks publication on a real Windows PowerShell clean-install/reinstall gate, then validates Go, builds the complete release, publishes it, verifies the exact remote asset set, and smoke-tests the public installer.
documents: {}
tags:
  - github-actions
  - go
  - release
keywords:
  - CI
  - publish
  - Windows PowerShell installer
  - protocol-v1 artifacts
flags: []
status: stable
---

Neither workflow invokes the legacy Bun/TypeScript domain implementation.
