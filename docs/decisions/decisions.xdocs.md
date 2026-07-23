---
subject: runx-decisions
description: Durable product and technical decisions for RunX.
parent: runx-docs
children: []
files:
  alpha-boundaries.md: Captures accepted alpha scope, safety, distribution, and compatibility decisions.
  citty-cli-migration.md: Defines the accepted full migration from handwritten parsing and routing to Citty.
  interactive-init-manifest.md: Historical v1 initializer decision whose manifest shape is superseded by manifest-v2-composition.md.
  markdown-release-assets-and-version-scoped-notes.md: Defines .md agent release filenames, downloaded-resource validation, and exact-version idempotent GitHub Release notes.
  mirror-automatic-push.md: Defines automatic Mirror release pushes with a synchronized protected-main gate.
  npm-trusted-publishing.md: Defines the accepted npm OIDC publishing workflow and protected-branch patch-release trial.
  no-argument-welcome-window.md: Replaces the legacy one-line greeting with a deterministic welcome and validated cached notice.
  run-argument-ownership.md: Defines the selector boundary, child-token ownership, dry-run behavior, and shell-safe transport.
  manifest-v2-composition.md: Replaces split groups and project.name with recursive commands, namespaces, reciprocal mounts, and bounded GitHub loading.
  windows-self-upgrade.md: Defines synchronous Windows executable replacement, verification, rollback, and cleanup.
documents:
  alpha-boundaries.md: Decision record for RunX alpha boundaries.
  citty-cli-migration.md: Decision record for the complete Citty CLI parser and router migration.
  interactive-init-manifest.md: Historical initializer decision; do not use its superseded v1 manifest shape.
  markdown-release-assets-and-version-scoped-notes.md: Accepted public asset filename, validation, and release-description policy.
  mirror-automatic-push.md: Accepted release-policy decision for Mirror push=true and protected delivery ordering.
  npm-trusted-publishing.md: Decision record for GitHub Actions npm trusted publishing and the first automated patch release.
  no-argument-welcome-window.md: Approved RunX bare-invocation welcome decision.
  run-argument-ownership.md: Approved RunX run argument-boundary decision.
  manifest-v2-composition.md: Accepted RunX manifest v2 composition decision.
  windows-self-upgrade.md: Accepted design for synchronous and recoverable Windows native self-upgrade.
tags:
  - decisions
keywords:
  - runx
  - decisions
  - runx init
  - scripts directory
  - mirror push
  - markdown release assets
  - version-scoped release notes
  - windows self-upgrade
flags: []
status: stable
---

Accepted decisions that later RunX work must preserve or intentionally supersede.
