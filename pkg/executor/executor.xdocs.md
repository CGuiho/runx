---
subject: runx-executor
description: Cross-platform configured-command execution with lossless child arguments, explicit shell transports, and caller-aware Windows automatic shell selection.
parent: runx-packages
children: []
files:
  executor.go: Builds safe POSIX, PowerShell, or cmd process invocations, resolves verified Git Bash callers on Windows, and preserves child exit codes.
  executor_test.go: Covers no-spawn dry-run behavior.
  shell_resolution_test.go: Covers deterministic caller-aware automatic shell resolution and safe fallbacks.
documents: {}
tags:
  - go
  - execution
keywords:
  - child arguments
  - exit code
  - Git Bash
  - shell auto
flags: []
status: stable
---

Forwarded values are positional parameters or environment-backed data, never
interpolated into shell source. On Windows, automatic selection uses Bash only
for a recognized MSYS/Git Bash caller with a resolved non-System32 executable;
all other automatic cases use `cmd.exe`.
