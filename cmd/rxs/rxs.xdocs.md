---
subject: runx-rxs
description: Thin rxs launcher that forwards to the active RunX payload smart command with root entry-flag parity.
parent: runx-cmd
children: []
files:
  main.go: Translates arguments (bare smart, root passthrough, smart forwarding), reads the pointer, resolves active-or-fallback payload, delegates, and maps failures to exit code 5.
  main_test.go: Exercises argument translation (bare smart, root passthrough, smart forwarding) for the rxs launcher.
  delegate_unix.go: Unix process delegation with identical streams and exit-code propagation.
  delegate_windows.go: Windows process delegation with identical streams and exit-code propagation.
documents: {}
tags:
  - go
  - launcher
  - rxs
keywords:
  - rxs
  - launcher
flags: []
status: stable
---

The `rxs` launcher forwards to the payload `smart` command (`rxs` -> `runx smart`, `rxs <selector>` -> `runx smart <selector>`) while the root entry points (`-v`/`--version`, `-h`/`--help`, `--help-tree`, `--help-docs`, `--help-tree-depth`, `--help-tree-global-flags`, single `--color`) pass through unchanged; the payload owns the list-or-run translation.
