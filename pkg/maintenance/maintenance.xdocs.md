---
subject: runx-maintenance
description: Bare-invocation bootstrap plus automatic and explicit RunX skill and repository instruction reconciliation with atomic cross-platform storage.
parent: runx-packages
children: []
files:
  maintenance.go: Resolves repository-root targets, reconciles both global skills and applicable instruction files, validates markers, and starts the maintenance worker.
  maintenance_test.go: Covers dual instruction targets, CRLF preservation, skill installation, idempotency, malformed markers, and non-repository behavior.
  storage.go: Provides text, directory, global-path, and atomic-write helpers.
  storage_test.go: Covers maintenance storage behavior.
  replace_unix.go: Atomically renames maintenance files on Unix.
  replace_windows.go: Atomically replaces maintenance files with MoveFileExW on Windows.
documents: {}
tags:
  - go
  - agents
keywords:
  - skill
  - AGENTS.md
  - atomic write
flags: []
status: stable
---

Maintenance is silent and failure-isolated during ordinary command startup;
explicit agent commands remain the repair and local-scope interface.
