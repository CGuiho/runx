---
subject: runx-executor
description: Cross-platform configured-command execution with lossless child arguments and explicit shell transports.
parent: runx-packages
children: []
files:
  executor.go: Builds safe POSIX, PowerShell, or cmd process invocations and preserves child exit codes.
  executor_test.go: Covers no-spawn dry-run behavior.
documents: {}
tags:
  - go
  - execution
keywords:
  - child arguments
  - exit code
flags: []
status: stable
---

Forwarded values are positional parameters or environment-backed data, never
interpolated into shell source.

