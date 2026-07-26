---
subject: runx-cmd
description: Fresh, dependency-injected Cobra command tree for RunX domain, help, agent, lifecycle, upgrade, and uninstall behavior.
parent: runx
children: []
files:
  root.go: Constructs fresh Cobra trees, maps exits, bootstraps bare-invocation agent integration, injects dependencies and build information, and schedules hidden lifecycle workers.
  catalog.go: Implements real manifest-backed check, list, describe, run, and init commands with deterministic text and JSON.
  help.go: Traverses live Cobra commands to produce Unicode tree and stable Markdown help.
  agent.go: Implements embedded skill, instruction, and prompt commands with global or local scope.
  upgrade.go: Implements release checks, complete listing, verified self-upgrade reporting, and uninstall behavior.
  root_test.go: Exercises bootstrap/exclusions, version, welcome, every help scope, aliases, real manifests, no-spawn inspection, argument forwarding, full agent actions, and initialization.
documents: {}
tags:
  - cli
  - go
  - cobra
keywords:
  - runx
  - command tree
  - deterministic JSON
flags: []
status: stable
---

All public routing originates in `NewRootCommand`; there is no package-global
mutable Cobra singleton or second parser.
