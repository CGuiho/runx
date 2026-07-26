---
subject: runx-update
description: Typed release discovery, target selection, local update cache, lease-coalesced background checking, and platform detachment.
parent: runx-packages
children: []
files:
  cache.go: Strictly decodes and atomically writes the global update cache and renders validated notices.
  cache_test.go: Covers cache read, write, freshness, and notice behavior.
  catalog.go: Fetches complete typed GitHub release pages and selects the canonical embedded-target asset and checksums.
  catalog_test.go: Covers SemVer ordering, platform mapping, pagination data, and compatible assets.
  worker.go: Coalesces bounded remote checks with a stale-recoverable lease and starts detached worker processes.
  worker_posix.go: Configures detached worker processes on non-Windows systems.
  worker_windows.go: Configures hidden detached worker processes on Windows.
  worker_test.go: Covers worker release checks and cache persistence.
  replace_unix.go: Atomically renames cache files on Unix.
  replace_windows.go: Atomically replaces cache files with MoveFileExW on Windows.
documents: {}
tags:
  - go
  - updates
keywords:
  - cache
  - GitHub Releases
  - detached worker
flags: []
status: stable
---

Foreground commands never perform the remote request; they only start a bounded
hidden worker after local cache handling.

