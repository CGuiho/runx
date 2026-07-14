---
subject: runx-decisions
description: Durable product and technical decisions for RunX.
parent: runx-docs
children: []
files:
  alpha-boundaries.md: Captures accepted alpha scope, safety, distribution, and compatibility decisions.
  citty-cli-migration.md: Defines the accepted full migration from handwritten parsing and routing to Citty.
  mirror-automatic-push.md: Defines automatic Mirror release pushes with a synchronized protected-main gate.
  npm-trusted-publishing.md: Defines the accepted npm OIDC publishing workflow and protected-branch patch-release trial.
  windows-self-upgrade.md: Defines synchronous Windows executable replacement, verification, rollback, and cleanup.
documents:
  alpha-boundaries.md: Decision record for RunX alpha boundaries.
  citty-cli-migration.md: Decision record for the complete Citty CLI parser and router migration.
  mirror-automatic-push.md: Accepted release-policy decision for Mirror push=true and protected delivery ordering.
  npm-trusted-publishing.md: Decision record for GitHub Actions npm trusted publishing and the first automated patch release.
  windows-self-upgrade.md: Accepted design for synchronous and recoverable Windows native self-upgrade.
tags:
  - decisions
keywords:
  - runx
  - decisions
  - mirror push
  - windows self-upgrade
flags: []
status: stable
---

Accepted decisions that later RunX work must preserve or intentionally supersede.
