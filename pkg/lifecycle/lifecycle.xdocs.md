---
subject: runx-lifecycle
description: Protocol-v1 whole-release upgrade engine with staging, verification, self-test, immutable install, atomic activation, launcher verification, and rollback.
parent: runx-packages
children: []
files:
  upgrade.go: Detects protocol-v1 installations, selects releases, stages and verifies every artifact under the shared temp root, self-tests the staged payload, installs an immutable version directory, atomically swaps current.json, verifies through the stable launcher, and rolls back on failure.
  upgrade_test.go: Covers installation detection and fail-closed behavior before any network work.
documents: {}
tags:
  - go
  - upgrade
  - lifecycle
keywords:
  - whole-release upgrade
  - protocol v1
  - rollback
flags: []
status: stable
---

Legacy direct-binary installations are reported as ErrLegacyInstallation so the
caller can fall back to the verified legacy replacement path.
