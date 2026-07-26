---
subject: runx-source
description: Preserved legacy Bun/TypeScript RunX implementation retained only for historical reference during the completed Go migration.
parent: runx
children: []
files:
  agent-maintenance.spec.ts: Historical Bun worker tests.
  agent-maintenance.ts: Historical Bun maintenance implementation.
  agents.ts: Historical Bun agent-resource implementation.
  cli.spec.ts: Historical Citty command tests.
  cli.ts: Historical Citty command tree; not production authority.
  configuration.spec.ts: Historical TypeBox manifest tests.
  configuration.ts: Historical TypeBox manifest implementation.
  embedded-resources.ts: Historical Bun resource embedding.
  errors.ts: Historical Bun error definitions.
  execution-arguments.spec.ts: Historical Bun argument-boundary tests.
  execution-arguments.ts: Historical Bun argument-boundary implementation.
  executor.spec.ts: Historical Bun executor tests.
  executor.ts: Historical Bun executor implementation.
  guiho-runx-bin.ts: Historical Bun development entrypoint.
  guiho-runx-native-bin.ts: Historical Bun compiled entrypoint.
  guiho-runx.spec.ts: Historical Bun integration tests.
  guiho-runx.ts: Historical TypeScript exports.
  help.ts: Historical Citty help generation.
  init.spec.ts: Historical Bun initializer tests.
  init.ts: Historical Bun initializer.
  manifest.ts: Historical TypeBox compatibility exports.
  path-utils.ts: Historical Bun path helpers.
  recovery.spec.ts: Historical recovery tests.
  recovery.ts: Historical recovery output.
  release-catalog.spec.ts: Historical Bun release-catalog tests.
  release-catalog.ts: Historical Bun release catalog.
  render.ts: Historical Bun renderer.
  self-management.spec.ts: Historical Bun self-upgrade tests.
  self-management.ts: Historical Bun self-upgrade implementation.
  storage.ts: Historical Bun storage helpers.
  types.ts: Historical TypeScript contracts.
  update-cache.spec.ts: Historical Bun update-cache tests.
  update-cache.ts: Historical Bun update cache.
  upgrade-reporting.spec.ts: Historical Bun upgrade-reporting tests.
  upgrade-reporting.ts: Historical Bun upgrade reporting.
  upgrade-types.ts: Historical TypeScript upgrade types.
  welcome.spec.ts: Historical Bun welcome tests.
  welcome.ts: Historical Bun welcome renderer.
documents: {}
tags:
  - legacy
  - bun
  - typescript
keywords:
  - runx
  - citty
  - migration history
flags:
  - deprecated
status: deprecated
---

Nothing in this directory is part of the production build, CI authority, native
release matrix, or publication path. It remains available for traceability until
a separately authorized cleanup removes the legacy reference.
