---
subject: runx-devops
description: Native release-asset compilation and direct RunX installer scripts.
parent: runx
children: []
files:
  build-binaries.ts: Cross-compiles supported Bun native executable assets.
  installers.spec.ts: Compiles a controlled native fixture and executes stable/prerelease recovery, 404 fallback, network/corrupt failures, space paths, exact verification, and rollback on the current platform while keeping PowerShell checks portable on Ubuntu.
  install.ps1: Resolves an exact Windows release, downloads and validates a compatible asset, transactionally replaces the canonical executable, verifies its version, and rolls back on failure.
  install.sh: POSIX-sh installer that resolves an exact macOS or Linux release, distinguishes 404 fallback from download/corrupt failures, transactionally installs, verifies, and rolls back.
documents: {}
tags:
  - devops
  - installers
  - releases
keywords:
  - runx
  - windows
  - macos
  - linux
  - native binary
flags: []
status: stable
---

These scripts support direct native installation without an npm publication.
