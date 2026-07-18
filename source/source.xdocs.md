---
subject: runx-source
description: Bun and TypeScript implementation of the RunX command-catalog CLI and native executable.
parent: runx
children: []
files:
  agents.ts: Implements dual-tool skill actions, idempotent instruction blocks, and typed bundled prompt discovery.
  cli.spec.ts: Exercises Citty help, root/nested version routing, upgrade flags, cached no-argument startup, aliases, errors, safety gates, and agent output.
  cli.ts: Defines the single RFC Citty tree, routes public upgrade flags without root-version interception, and owns startup, catalog, agent, and uninstall output.
  configuration.ts: Resolves YAML by explicit/cwd/global precedence, TypeBox-decodes manifests, and resolves catalog selectors.
  embedded-resources.ts: Embeds the agent skill and canonical prompt in compiled native executables.
  errors.ts: Defines user-facing RunX errors and assertions.
  executor.ts: Maps configured shell choices to local child-process execution.
  guiho-runx-bin.ts: Bun development CLI entrypoint.
  guiho-runx-native-bin.ts: Native executable entrypoint that registers embedded resources.
  guiho-runx.spec.ts: Tests manifest discovery, selector handling, execution safety, and native self-management routing.
  guiho-runx.ts: Public package exports.
  help.ts: Traverses actual Citty definitions to render Unicode tree help, depth-limited help, Markdown docs, and version output.
  init.spec.ts: Tests initializer creation, cancellation, overwrite safety, and exact manifest output without a terminal.
  init.ts: Creates a strict initial runx.yaml at the selected effective path.
  manifest.ts: Compatibility export surface for the schema-backed configuration module.
  path-utils.ts: Provides narrow Bun-first cross-platform path resolution without prohibited Node path imports.
  render.ts: Renders text and JSON command catalog output.
  storage.ts: Provides Bun and Bun-shell backed global storage, text I/O, directory, and removal operations.
  self-management.spec.ts: Proves output/envelope contracts, downgrade prevention, synchronous replacement, second-rename failure state, rollback success/failure, and Windows mapped-image behavior.
  self-management.ts: Plans upgrades, prevents downgrade, preserves post-backup mutation state, classifies stable failures, verifies replacement, rolls back, and performs native uninstall operations.
  release-catalog.ts: Retrieves every GitHub release page, applies SemVer ordering and channel labels, and selects compatible assets.
  release-catalog.spec.ts: Covers pagination, SemVer/prerelease ordering, malformed responses, and platform asset candidates.
  recovery.ts: Generates exact-version platform installation and explicit process-stop recovery commands.
  recovery.spec.ts: Verifies pinned and separate recovery commands for Windows and POSIX systems.
  upgrade-reporting.ts: Renders streamed human upgrade phases, complete release tables, final outcomes, and recovery instructions.
  upgrade-reporting.spec.ts: Verifies human output ordering, aligned catalog metadata, and recovery visibility.
  types.ts: Defines shared CLI, manifest, command, agent, and self-upgrade result types.
  update-cache.spec.ts: Covers foreground cache decoding and worker update/up-to-date writes.
  update-cache.ts: Reads cached notices, detaches update workers, validates releases, and atomically refreshes cache data.
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

The implementation uses one raw Citty tree and TypeBox at every structured
boundary. Bun-first configuration, storage, execution, update caching, and
self-management modules keep core source free of prohibited Node built-ins.
