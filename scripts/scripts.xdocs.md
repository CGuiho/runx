---
subject: runx-scripts
description: Package-manager launcher scripts for RunX.
parent: runx
children: []
files:
  runx-bin.mjs: Node-compatible npm bootstrap that downloads, caches, and delegates to the exact native release binary.
  runx-bin.spec.ts: Packs the npm package and proves Node-only bootstrap download/delegation with Bun absent from PATH.
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

The npm launcher contains no RunX domain logic and requires no Bun installation.
