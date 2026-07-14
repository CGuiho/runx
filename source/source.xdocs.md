---
subject: runx-source
description: Bun and TypeScript implementation of the RunX command-catalog CLI and native executable.
parent: runx
children: []
files:
  agents.ts: Installs the bundled guiho-s-runx skill and managed AGENTS.md instructions.
  cli.ts: Parses commands and reports manifest, agent, and native self-management outcomes, including explicit up-to-date upgrades.
  embedded-resources.ts: Embeds the agent skill in compiled native executables.
  errors.ts: Defines user-facing RunX errors and assertions.
  executor.ts: Maps configured shell choices to local child-process execution.
  flags.ts: Parses positional arguments and global flags.
  guiho-runx-bin.ts: Bun development CLI entrypoint.
  guiho-runx-native-bin.ts: Native executable entrypoint that registers embedded resources.
  guiho-runx.spec.ts: Tests manifest discovery, selector handling, and argument parsing.
  guiho-runx.ts: Public package exports.
  help.ts: Provides the home page, command tree, documentation, and version output.
  manifest.ts: Discovers, parses, validates, and resolves runx.yaml manifests.
  render.ts: Renders text and JSON command catalog output.
  self-management.ts: Checks GitHub Releases, returns equal-version no-op upgrades, and performs native replacement or uninstall operations.
  types.ts: Defines shared CLI, manifest, command, agent, and self-upgrade result types.
documents: {}
tags:
  - bun
  - cli
  - typescript
keywords:
  - runx
  - manifest
  - command catalog
  - native binary
flags: []
status: stable
---

The RunX implementation validates a local manifest before any command is listed
or run. The native entrypoint embeds the bundled agent skill so `runx agents
install` remains available without Bun at runtime.
