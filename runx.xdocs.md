---
subject: runx
description: Open-source Bun CLI package for documented, language-agnostic local command catalogs with interactive initialization.
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
  package.json: Package metadata, scripts, native CLI entrypoint, and runtime dependencies.
  tsconfig.json: Strict TypeScript compiler configuration for library output and type checking.
  xdocs.config.toml: XDocs metadata and agent automation configuration.
documents:
  AGENTS.md: Repository instructions for coding agents, including Citty CLI ownership and protected release boundaries.
  CHANGELOG.md: Mirror-managed release history, including the 0.2.6 verified Windows self-upgrade fix.
  CONTRIBUTING.md: Contribution, validation, and protected release workflow guide.
  DOCS.md: Canonical Citty-backed CLI, manifest, complete release catalog, transactional upgrade/recovery, direct installer, distribution, and release reference.
  LICENSE.md: MIT license.
  README.md: Public project overview, installation, interactive quick start, Citty aliases, and verified upgrade/list/recovery commands.
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
  - initialization
  - native binary
flags: []
status: stable
---

RunX is a standalone command-catalog CLI. It keeps project operations in a
strict, documented `runx.yaml` manifest and supports local human and agent use
without requiring the target project to use JavaScript. `runx init` creates an
empty, strict catalog through a guided terminal workflow before commands are
added.
