---
subject: runx
description: Open-source Bun CLI package for documented, language-agnostic local command catalogs with platform-aware startup and direct native installation.
parent: null
children:
  - runx-source
  - runx-scripts
  - runx-skills
  - runx-prompts
  - runx-devops
  - runx-docs
files:
  .gitignore: Ignores local dependencies and generated build outputs.
  bun.lock: Locks Bun dependencies for reproducible installs.
  mirror.config.toml: Mirror semantic versioning configuration.
  package.json: Package metadata, RFC scripts, Node-compatible npm bootstrap entrypoint, and runtime dependencies.
  tsconfig.json: Strict TypeScript compiler configuration for library output and type checking.
  xdocs.config.toml: XDocs metadata and agent automation configuration.
documents:
  AGENTS.md: Repository instructions for coding agents, including the mandatory SWE agent, CLI engineer skill, Citty ownership, breaking RFC migration, and protected release boundaries.
  CHANGELOG.md: Mirror-managed release history, including the 0.4.0 full RFC 0034 migration.
  CONTRIBUTING.md: Contribution, validation, and protected release workflow guide.
  DOCS.md: Canonical Citty-backed CLI, manifest, complete release catalog, transactional upgrade/recovery, direct installer, distribution, and release reference.
  LICENSE.md: MIT license.
  README.md: Public RFC 0034 overview, installation, command catalog, YAML precedence, and help entrypoints.
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
  - RFC 0034
  - cli engineer
flags: []
status: stable
---

RunX is a standalone RFC 0034 command-catalog CLI with one Citty tree, TypeBox
boundaries, Bun-only core source, platform-aware startup, complete agent
integration, transactional native installation and upgrades, a Node-compatible
npm bootstrap, and fourteen release assets.
