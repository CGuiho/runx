---
subject: runx-scripts
description: Package-manager launcher scripts for RunX.
parent: runx
children: []
files:
  runx-bin.ts: Executes a packaged native binary when present or the source CLI in a checkout.
documents: {}
tags:
  - scripts
  - cli
keywords:
  - runx
  - launcher
flags: []
status: stable
---

The package launcher preserves the `runx` bin command in development checkouts
and installed packages.
