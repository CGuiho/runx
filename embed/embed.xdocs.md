---
subject: runx-embed
description: Compile-time RunX skill and prompt resources used by agent commands and automatic maintenance.
parent: runx
children: []
files:
  embed.go: Exposes the compile-time embedded filesystem for the RunX skill and prompt copies.
documents: {}
tags:
  - go
  - embed
keywords:
  - go:embed
  - agent resources
flags: []
status: stable
---

Embedded copies are validated against their canonical top-level sources before
release.

