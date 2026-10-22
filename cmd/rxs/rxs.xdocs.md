---
subject: runx-rx
description: Thin rx launcher that translates ergonomics and delegates to the active RunX payload.
parent: runx-cmd
children: []
files:
  main.go: Translates arguments (bare list, passthrough, run selector), reads the pointer, resolves active-or-fallback payload, delegates, and maps failures to exit code 5.
  main_test.go: Exercises argument translation (bare list, passthrough, run selector) for the rx launcher.
  delegate_unix.go: Unix process delegation with identical streams and exit-code propagation.
  delegate_windows.go: Windows process delegation with identical streams and exit-code propagation.
documents: {}
tags:
  - go
  - launcher
  - rx
keywords:
  - rx
  - launcher
flags: []
status: stable
---

The `rx` launcher translates ergonomics (`rx` -> `runx list`, `rx <selector>` -> `runx run <selector>`, version and help passthrough) and delegates directly to the active RunX payload.
