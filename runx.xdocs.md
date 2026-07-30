---
subject: runx
description: Native Go/Cobra command-catalog CLI with strict manifest v2, safe execution, bare-invocation agent bootstrap, verified upgrades, installers, and release automation.
parent: null
children:
  - runx-cmd
  - runx-packages
  - runx-embed
  - runx-scripts
  - runx-skills
  - runx-prompts
  - runx-devops
  - runx-source
  - runx-docs
files:
  .gitignore: Ignores local dependencies, caches, and generated build outputs.
  bun.lock: Preserved lockfile for the legacy TypeScript reference implementation.
  go.mod: Pins the production Go toolchain and module dependencies.
  go.sum: Records production Go dependency checksums.
  main.go: Thin production entrypoint that supplies embedded build metadata and mapped exit codes to the Cobra tree.
  mirror.yaml: Mirror semantic version configuration using Git tags as the version source.
  package.json: npm native-bootstrap metadata and Go-backed contributor command aliases; it contains no RunX domain implementation.
  tsconfig.json: Preserved compiler configuration for the legacy TypeScript reference implementation.
  xdocs.config.toml: XDocs metadata configuration with automatic documentation-update mode.
documents:
  AGENTS.md: Repository rules for the Go CLI, XDocs, Mirror, validation, and protected release boundaries.
  CHANGELOG.md: Mirror-managed historical release record.
  CONTRIBUTING.md: Go validation, documentation, and protected release contribution workflow.
  DOCS.md: Complete production CLI, manifest, lifecycle, agent, upgrade, installer, and 11-artifact reference.
  LICENSE.md: MIT license.
  README.md: Public installation with concise verified bootstrap commands and cross-shell PATH guidance, plus manifest-v2, command, and Go release overviews.
  SECURITY.md: Vulnerability reporting and trusted-manifest boundary.
  TODO.md: Package-local task index and migration status.
tags:
  - cli
  - go
  - cobra
  - open source
keywords:
  - runx
  - command catalog
  - yaml
  - manifest v2
  - native binary
  - RFC 0034
flags: []
status: stable
---

RunX production behavior is owned by one testable Cobra command tree and focused
Go packages. Structured manifest validation, command execution, idempotent
bare-invocation agent bootstrap, cached lifecycle workers, agent resources,
verified upgrades, installers, and the standard
11-artifact release matrix are enforced by Go tests and workflows.
