---
subject: runx-source
description: Bun and TypeScript implementation of the RunX command-catalog CLI and native executable.
parent: runx
children: []
files:
  agents.ts: Installs the bundled guiho-s-runx skill and managed AGENTS.md instructions.
  cli.spec.ts: Exercises Citty help, version, command routing, aliases, global options, errors, safety gates, and agent paths.
  cli.ts: Defines the Citty command tree, including interactive initialization, and reports manifest, agent, and native self-management outcomes.
  embedded-resources.ts: Embeds the agent skill in compiled native executables.
  errors.ts: Defines user-facing RunX errors and assertions.
  executor.ts: Maps configured shell choices to local child-process execution.
  guiho-runx-bin.ts: Bun development CLI entrypoint.
  guiho-runx-native-bin.ts: Native executable entrypoint that registers embedded resources.
  guiho-runx.spec.ts: Tests manifest discovery, selector handling, execution safety, and native self-management routing.
  guiho-runx.ts: Public package exports.
  help.ts: Provides the home page, command tree, upgrade/list guidance, documentation, and version output.
  init.spec.ts: Tests initializer creation, cancellation, overwrite safety, and exact manifest output without a terminal.
  init.ts: Provides the interactive terminal workflow, preview, confirmation, and atomic write for a new runx.yaml manifest.
  manifest.ts: Discovers, parses, validates, and resolves strict SemVer, scripts-directory, and public-group runx.yaml manifests.
  render.ts: Renders text and JSON command catalog output.
  self-management.spec.ts: Proves output/envelope contracts, downgrade prevention, synchronous replacement, second-rename failure state, rollback success/failure, and Windows mapped-image behavior.
  self-management.ts: Plans upgrades, prevents downgrade, preserves post-backup mutation state, classifies stable failures, verifies replacement, rolls back, and performs native uninstall operations.
  release-catalog.ts: Retrieves every GitHub release page, applies SemVer ordering and channel labels, and selects compatible assets.
  release-catalog.spec.ts: Covers pagination, SemVer/prerelease ordering, malformed responses, and platform asset candidates.
  recovery.ts: Generates exact-version platform installation and explicit process-stop recovery commands.
  recovery.spec.ts: Verifies pinned and separate recovery commands for Windows and POSIX systems.
  upgrade-reporting.ts: Renders streamed human upgrade phases, complete release tables, final outcomes, and recovery instructions.
  upgrade-reporting.spec.ts: Verifies human output ordering, aligned catalog metadata, and recovery visibility.
  types.ts: Defines shared CLI, manifest, command, agent, and self-upgrade result types.
  upgrade-types.ts: Defines the exact release catalog plus nullable-plan, result, stable error, rolled-back outcome, and target-sourced recovery contracts.
documents: {}
tags:
  - bun
  - cli
  - typescript
keywords:
  - runx
  - citty
  - manifest
  - initialization
  - command catalog
  - native binary
flags: []
status: stable
---

The RunX implementation uses Citty for command parsing, aliases, routing, and
ordinary usage, then validates a local manifest before any catalog command is
listed or run. The native entrypoint embeds the bundled agent skill so `runx
agents install` remains available without Bun at runtime. `init.ts` provides a
dependency-free interactive terminal workflow that produces the same strict
empty manifest contract that later commands validate.
