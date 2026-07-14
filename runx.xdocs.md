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
  package.json: Package metadata, scripts, native CLI entrypoint, and Citty runtime dependency.
  tsconfig.json: Strict TypeScript compiler configuration for library output and type checking.
  xdocs.config.toml: XDocs metadata and agent automation configuration.
documents:
  AGENTS.md: Repository instructions for coding agents, including Citty CLI ownership and protected release boundaries.
  CHANGELOG.md: Mirror-managed release history, including the 0.2.2 trusted-publishing retry and native distribution changes.
  CONTRIBUTING.md: Contribution, validation, and protected release workflow guide.
  DOCS.md: Canonical Citty-backed CLI, manifest, distribution, and release reference.
  LICENSE.md: MIT license.
  README.md: Public project overview, npm/native installation paths, Citty help aliases, and quick start.
  SECURITY.md: Vulnerability reporting and manifest trust-boundary policy.
  TODO.md: Package-local task index.
tags:
  - cli
  - bun
  - open source
keywords:
  - runx
  - citty
  - command catalog
  - yaml
  - native binary
flags: []
status: stable
---

RunX is a standalone command-catalog CLI. It keeps project operations in a
strict, documented `runx.yaml` manifest and supports local human and agent use
without requiring the target project to use JavaScript.
