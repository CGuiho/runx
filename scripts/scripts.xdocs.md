---
subject: runx-scripts
description: Node-compatible npm bootstrap scripts that verify and delegate to the native Go RunX executable.
parent: runx
children: []
files:
  runx-bin.mjs: Downloads the package version's canonical native target and checksums, verifies SHA-256, caches it, and delegates the process.
  runx-bin.spec.ts: Historical Bun-hosted test for the Node bootstrap; current CI performs Node syntax validation and Go release-matrix verification.
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
