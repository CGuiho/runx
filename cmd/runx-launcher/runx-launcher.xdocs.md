---
subject: runx-launcher-command
description: Platform launcher entrypoint that delegates to the committed active payload and exits with its code.
parent: runx-cmd
children:
  - runx-launcher-command
files:
  main.go: Reads the pointer, resolves active-or-fallback payload, delegates, and maps failures to exit code 5.
  delegate_unix.go: Unix process delegation with identical streams and exit-code propagation.
  delegate_windows.go: Windows process delegation with identical streams and exit-code propagation.
documents: {}
tags:
  - go
  - launcher
keywords:
  - stable launcher
flags: []
status: stable
---

The shell stays attached until the payload completes; activation is never
background work.
