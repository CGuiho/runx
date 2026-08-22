---
subject: runx-cmd
description: Fresh, dependency-injected Cobra command tree for RunX domain, help, agent, lifecycle, upgrade, and uninstall behavior.
parent: runx
children:
  - runx-launcher-command
files:
  root.go: Constructs fresh Cobra trees, maps exits, bootstraps bare-invocation agent integration, injects dependencies including terminal detection and build information, schedules hidden lifecycle workers, exposes the hidden installation self-test, and owns the help-tree flags including --help-tree-global-flags.
  catalog.go: Implements real manifest-backed check, aligned list, numeric-index-aware describe/reveal/run, and safe-default interactive confirmation with exact retries.
  help.go: Traverses live Cobra commands to produce Unicode tree output that shows global flags once by default or under every descendant with --help-tree-global-flags, plus stable Markdown help.
  agent.go: Implements embedded skill, instruction, and prompt commands with global or local scope.
  init.go: Implements the Convention 0001 initialization sequence with Created/Upgraded/Verified/Unchanged reporting, interactive evolution-policy questions stored in the global configuration, and fail-closed noninteractive behavior.
  upgrade.go: Implements release checks, complete listing, verified self-upgrade reporting with mandatory pre-network and final pinned reinstallation recovery blocks, and protocol-v1 whole-release upgrade dispatch.
  uninstall.go: Implements the shared Convention 0001 uninstallation contract with REMOVE/PRESERVE planning, preservation flags, fail-closed noninteractive confirmation, and managed-block-only instruction removal.
  root_test.go: Exercises bootstrap/exclusions, version, welcome, every help scope, aliases, aligned lists, numeric index selection, reveal output, confirmations, real manifests, argument forwarding, full agent actions, and initialization.
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
