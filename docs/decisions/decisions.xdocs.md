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
documents:
  alpha-boundaries.md: Decision record for RunX alpha boundaries.
  citty-cli-migration.md: Decision record for the complete Citty CLI parser and router migration.
  mirror-automatic-push.md: Accepted release-policy decision for Mirror push=true and protected delivery ordering.
  npm-trusted-publishing.md: Decision record for GitHub Actions npm trusted publishing and the first automated patch release.
tags:
  - decisions
keywords:
  - runx
  - decisions
  - mirror push
flags: []
status: stable
---

Accepted decisions that later RunX work must preserve or intentionally supersede.
