---
subject: runx-updater
description: Checksum-verifying target-preserving self-upgrade, rollback, recovery, and staged Windows replacement.
parent: runx-packages
children: []
files:
  types.go: Defines deterministic upgrade plan, event, result, recovery, error, and dependency contracts.
  upgrade.go: Resolves releases, downloads binary and checksums, verifies SHA-256 and native format, and orchestrates replacement.
  upgrade_test.go: Covers up-to-date, dry-run, checksum-backed upgrade, and native format behavior.
  fileops.go: Implements platform-aware filesystem operations and deferred Windows backup cleanup.
  rollback.go: Performs same-filesystem staging, verification, rollback, and running-binary detection.
  rollback_test.go: Covers replacement success and rollback failure modes.
  stage_windows.go: Starts the hidden post-exit Windows replacement and verification helper.
  stage_other.go: Rejects Windows staging on foreign platforms.
documents: {}
tags:
  - go
  - upgrade
keywords:
  - SHA-256
  - rollback
  - Windows replacement
flags: []
status: stable
---

No downloaded binary reaches the installed path before checksum and native-file
validation succeed.

