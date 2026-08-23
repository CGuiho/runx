---
subject: runx-lifecycle
description: Protocol-v1 whole-release upgrade engine with staging, verification, self-test, immutable install, atomic activation, launcher verification, and rollback.
parent: runx-packages
children: []
files:
  upgrade.go: Detects protocol-v1 installations, returns verified up-to-date for an already-active exact or default target, and otherwise stages and verifies every checksummed artifact without requiring the checksum manifest to hash itself before immutable activation and rollback-safe launcher verification.
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
