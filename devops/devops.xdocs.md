---
subject: runx-devops
description: Native release-asset compilation and direct RunX installer scripts.
parent: runx
children: []
files:
  build-binaries.ts: Cross-compiles twelve RFC-named native binaries and adds the two exact agent assets.
  verify-release-assets.ts: Fails on any missing, duplicate, extra, or legacy release asset name.
  installers.spec.ts: Verifies both installers expose progress, validation, dual-tool skills, instruction reconciliation, and no POSIX Bun dependency.
  install.ps1: Resolves an exact Windows release, downloads and validates a compatible asset, transactionally replaces the canonical executable, verifies its version, and rolls back on failure.
  install.sh: POSIX installer for Linux or Darwin with progress, validation, PATH, dual skills, instructions, verification, and rollback.
documents: {}
tags:
  - devops
  - installers
  - releases
keywords:
  - runx
  - windows
  - darwin
  - linux
  - native binary
flags: []
status: stable
---

These scripts build and verify the exact fourteen-asset release and support
direct installation without npm or Bun.
