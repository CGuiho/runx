---
subject: runx-launcher
description: Stable-launcher payload dispatch with committed-pointer resolution and retained-fallback selection.
parent: runx-packages
children: []
files:
  launcher.go: Resolves the active payload from a validated pointer and the previous verified payload as fallback, failing closed on missing or unsafe targets.
documents: {}
tags:
  - go
  - launcher
keywords:
  - delegation
  - fallback
  - current.json
flags: []
status: stable
---

The launcher never interprets domain arguments; it delegates to the committed
active payload with identical arguments, environment, working directory, and
standard streams, then returns the payload's exact exit code. Fallback applies
only when the active payload cannot start.
