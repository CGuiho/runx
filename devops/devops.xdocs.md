---
subject: runx-devops
description: Native release-asset compilation and direct RunX installer scripts.
parent: runx
children: []
files:
  build-binaries.ts: Cross-compiles twelve RFC-named native binaries and adds the two exact .md agent assets.
  extract-release-notes.ts: Extracts one exact version section from CHANGELOG.md and fails closed when the heading or notes are missing.
  extract-release-notes.spec.ts: Verifies exact heading boundaries and exclusion of frontmatter and other release sections.
  verify-release-assets.ts: Fails on missing, duplicate, extra, legacy, empty, binary, or misidentified release assets.
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

These scripts build and verify the exact fourteen-asset release, validate the
downloaded Markdown resource contract, create exact version-only GitHub Release
notes, and support direct installation without npm or Bun.
