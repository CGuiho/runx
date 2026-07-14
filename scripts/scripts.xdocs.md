---
subject: runx-scripts
description: Package-manager launcher scripts for RunX.
parent: runx
children: []
files:
  runx-bin.ts: Executes a packaged native binary, the published compiled library launcher, or the source CLI fallback in a checkout.
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

The package launcher preserves the `runx` bin command in installed npm packages
and development checkouts without requiring unpublished source files.
