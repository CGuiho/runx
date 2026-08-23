---
subject: runx-lifecycle
description: Protocol-v1 whole-release upgrade engine with staging, verification, self-test, immutable install, atomic activation, launcher verification, and rollback.
parent: runx-packages
children: []
files:
  upgrade.go: Detects protocol-v1 installations, selects releases, stages and verifies every checksummed artifact without requiring the checksum manifest to hash itself, self-tests the payload, installs an immutable version, atomically swaps current.json, verifies through the launcher, and rolls back on failure.
  upgrade_test.go: Covers installation detection, fail-closed behavior before network work, and the standard checksum-manifest rule that checksums.txt never hashes itself.
documents: {}
tags:
  - go
  - upgrade
  - lifecycle
keywords:
  - whole-release upgrade
  - protocol v1
  - checksums.txt
  - rollback
flags: []
status: stable
---

Legacy direct-binary installations are reported as ErrLegacyInstallation so the
caller can fall back to the verified legacy replacement path.
