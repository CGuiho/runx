---
subject: runx-source
description: Bun and TypeScript implementation of the RunX command-catalog CLI and native executable.
parent: runx
children: []
files:
  agent-maintenance.spec.ts: Proves missing/current/stale resources, legacy migration, concurrency convergence, hidden worker validation, and spawn failure isolation.
  agent-maintenance.ts: Validates and runs the hidden automatic agent worker and detaches failure-isolated maintenance with effective cwd.
  agents.ts: Implements explicit dual-tool actions plus atomic global skill and nearest AGENTS.md reconciliation.
  cli.spec.ts: Exercises the welcome, argument forwarding, Citty Unicode help/alignment, routing, errors, automatic maintenance, upgrade leaf isolation, output safety, and agent behavior.
  cli.ts: Defines the single RFC Citty tree, welcome lifecycle, argument ownership boundary, hidden workers, automatic maintenance, isolated upgrade leaf routing, and public output.
  configuration.ts: Resolves YAML by explicit/cwd/global precedence, TypeBox-decodes manifests, and resolves catalog selectors.
  embedded-resources.ts: Embeds the agent skill and canonical prompt in compiled native executables.
  errors.ts: Defines user-facing RunX errors and assertions.
  execution-arguments.spec.ts: Proves RunX option ownership, delimiter removal, and lossless hostile-looking child token preservation.
  execution-arguments.ts: TypeBox-validates the RunX-owned router prefix and post-selector child argument boundary.
  executor.spec.ts: Proves POSIX positional, PowerShell splat, and cmd environment-backed argument transport never interpolates child values into shell source.
  executor.ts: Maps configured shells to local execution and safely transports forwarded child argument arrays.
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
  storage.ts: Provides Bun-first global storage, directory operations, and direct or atomic text writes.
  self-management.spec.ts: Proves output/envelope contracts, downgrade prevention, synchronous replacement, second-rename failure state, rollback success/failure, and Windows mapped-image behavior.
  self-management.ts: Plans upgrades, prevents downgrade, preserves post-backup mutation state, classifies stable failures, verifies replacement, rolls back, and performs native uninstall operations.
  release-catalog.ts: Retrieves every GitHub release page, applies SemVer ordering and channel labels, selects compatible assets, and applies only explicit view pagination.
  release-catalog.spec.ts: Covers complete default output, explicit pagination, SemVer/prerelease ordering, malformed responses, and platform asset candidates.
  recovery.ts: Generates exact-version PowerShell or Bash installation and explicit process-stop recovery commands.
  recovery.spec.ts: Verifies pinned and separate recovery commands for Windows and Bash-based systems.
  upgrade-reporting.ts: Renders streamed human upgrade phases, complete release tables, final outcomes, and recovery instructions.
  upgrade-reporting.spec.ts: Verifies output ordering, aligned catalog metadata, and exact recovery after every terminal outcome.
  types.ts: Defines shared CLI, manifest, command, agent, and self-upgrade result types.
  update-cache.spec.ts: Covers foreground cache decoding, 64-way process coalescing, cache freshness, ownership-safe stale recovery, hard deadlines, malformed leases, and cache writes.
  update-cache.ts: Reads cached notices and schedules one ownership-leased, time-bounded detached update worker per global cache directory.
  welcome.spec.ts: Proves deterministic product, platform, architecture, version, help, and cached-update rendering.
  welcome.ts: Purely renders the stable no-argument RunX welcome window without ANSI or side effects.
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
