---
subject: runx-config
description: Strict separate project and global configuration contracts with the mandatory agent-evolution policy and version-pinned schema references.
parent: runx-packages
children: []
files:
  config.go: Defines strict ProjectConfig/GlobalConfig types, the agent.evolution policy enum, overlay resolution, and schema URL pinning.
  load.go: Strict YAML loading with unknown-field rejection, semantic validation, and schema-pinned writers.
  config_test.go: Covers defaults, per-field project-over-global precedence, enum rejection, unknown fields, round trips, and pinned schemas.
documents: {}
tags:
  - go
  - configuration
keywords:
  - runx.yaml
  - runx.global.yaml
  - agent evolution
  - always-ask
flags: []
status: stable
---

The project contract (`runx.yaml`) and the global baseline
(`runx.global.yaml` in the CLI home) are independent strict documents.
Unknown fields fail closed; every `agent.evolution` value must be exactly
`disabled`, `always-ask`, or `always-proceed`, defaulting to `always-ask`.
