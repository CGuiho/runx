---
subject: runx-devops
description: Go release builders, asset verification, release-note extraction, Convention 0001 validation helpers, and checksum-verifying lifecycle installers and uninstallers for RunX.
parent: runx
children:
  - runx-conventionv1-devops
files:
  build-binaries.go: Builds the eight immutable runx-payload targets, eight stable launchers, skill ZIP, instruction Markdown, both configuration schemas, artifacts.json ownership manifest, and deterministic checksum manifest.
  verify-release-assets.go: Verifies the exact protocol-v1 declared artifact set, manifest digests, every SHA-256 entry, and the skill archive structure.
  extract-release-notes.go: Extracts one exact version section from CHANGELOG.md for protected-tag publication.
  install.ps1: Stages under $HOME/.guiho/.temp/, verifies checksums and the release manifest, self-tests the payload before activation, installs immutable payloads plus the stable launcher into the canonical .guiho layout, activates via atomic current.json, preserves configuration and data, rolls back on failure, and configures user PATH idempotently.
  install_ps1_test.go: Verifies canonical layout usage, full-name-only flags, shared uninstall contract options, hidden self-test invocation, and PowerShell/bash syntax.
  install.sh: Detects Linux AMD64, ARM64, ARMv7, ARMv6 or Darwin targets, supports --version/--channel selection across full release pagination, verifies checksums and manifest digests, self-tests before activation, installs the canonical layout with rollback, and updates PATH idempotently.
  uninstall.ps1: Shares the uninstallation contract with REMOVE/PRESERVE planning, preservation options, dry-run, fail-closed confirmation, managed-block-only instruction removal, and Windows delete-on-close quarantine semantics.
  uninstall.sh: Shares the uninstallation contract with REMOVE/PRESERVE planning, preservation options, dry-run, fail-closed confirmation, owned-path-only removal, and managed-block-only instruction removal.
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
