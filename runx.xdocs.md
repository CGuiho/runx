---
subject: runx
description: Open-source Bun CLI package for documented, language-agnostic local command catalogs.
parent: null
children:
  - runx-source
  - runx-scripts
  - runx-skills
  - runx-devops
  - runx-docs
files:
  .gitignore: Ignores local dependencies and generated build outputs.
  bun.lock: Locks Bun dependencies for reproducible installs.
  mirror.config.toml: Mirror semantic versioning configuration.
  package.json: Package metadata, scripts, native CLI entrypoint, and dependencies.
  tsconfig.json: Strict TypeScript compiler configuration for library output and type checking.
  xdocs.config.toml: XDocs metadata and agent automation configuration.
documents:
  AGENTS.md: Repository instructions for coding agents.
  CHANGELOG.md: Release-level change history.
  CONTRIBUTING.md: Contribution and validation guide.
  DOCS.md: Canonical user-facing CLI and manifest reference.
  LICENSE.md: MIT license.
  README.md: Public project overview and quick start.
  SECURITY.md: Vulnerability reporting and manifest trust-boundary policy.
  TODO.md: Package-local task index.
tags:
  - cli
  - bun
  - open source
keywords:
  - runx
  - command catalog
  - yaml
  - native binary
flags: []
status: stable
---

RunX is a standalone command-catalog CLI. It keeps project operations in a
strict, documented `runx.yaml` manifest and supports local human and agent use
without requiring the target project to use JavaScript.
