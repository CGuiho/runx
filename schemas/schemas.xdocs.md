---
subject: runx-schemas
description: JSON Schemas for the strict project and global RunX configuration contracts, published with every release.
parent: runx
children: []
files:
  runx.schema.json: Strict schema for runx.yaml including the agent.evolution policy enum.
  runx.global.schema.json: Strict schema for runx.global.yaml including the agent.evolution policy enum.
documents: {}
tags:
  - configuration
  - schemas
keywords:
  - json schema
  - yaml-language-server
flags: []
status: stable
---

Generated configuration files reference these schemas through version-pinned
release URLs; runtime validation uses the embedded Go contract and never
requires network access.
