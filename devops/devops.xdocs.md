---
subject: runx-devops
description: Go release builders, asset verification, release-note extraction, and checksum-verifying direct installers for RunX.
parent: runx
children: []
files:
  build-binaries.go: Builds the eight standard pure-Go targets, skill ZIP, instruction Markdown, and deterministic checksum manifest.
  verify-release-assets.go: Verifies the exact 11-file set, every SHA-256 entry, and the skill archive structure.
  extract-release-notes.go: Extracts one exact version section from CHANGELOG.md for protected-tag publication.
  install.ps1: Detects Windows AMD64 or ARM64, verifies checksums, replaces transactionally, installs skills, verifies the version, and configures Windows and Git Bash PATH discovery.
  install_ps1_test.go: Verifies idempotent Git Bash PATH configuration, content and line-ending preservation, and the concise public PowerShell bootstrap.
  install.sh: Detects Linux AMD64, ARM64, ARMv7, ARMv6 or Darwin targets, verifies checksums, replaces transactionally, installs skills, and verifies the version.
  build-binaries.ts: Legacy Bun 14-asset builder retained as migration history and not invoked by current scripts or workflows.
  verify-release-assets.ts: Legacy Bun asset verifier retained as migration history and not invoked by current scripts or workflows.
  extract-release-notes.ts: Legacy Bun release-note extractor retained as migration history and not invoked by current workflows.
  extract-release-notes.spec.ts: Historical Bun release-note tests.
  installers.spec.ts: Historical Bun installer tests for the retired asset contract.
  workflows.spec.ts: Historical Bun workflow tests for the retired release path.
documents: {}
tags:
  - devops
  - go
  - installers
  - releases
keywords:
  - runx
  - checksums
  - cross compilation
  - Raspberry Pi
  - 11 artifacts
flags: []
status: stable
---

The Go programs and shell installers are current release authority. They build
and verify eight executables plus `guiho-s-runx.zip`, `guiho-i-runx.md`, and
`checksums.txt`. The Windows installer persists the install directory for
Windows processes and Git Bash startup without duplicating entries. TypeScript
files in this directory are retained history only.
