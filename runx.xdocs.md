---
subject: runx
description: Native Go/Cobra command-catalog CLI with strict manifest v2, exact non-executing command revelation, safe execution, bare-invocation agent bootstrap, verified upgrades, installers, and release automation.
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
  - runx-schemas
  - runx-examples
files:
  .gitignore: Ignores local dependencies, caches, and generated build outputs.
  bun.lock: Preserved lockfile for the legacy TypeScript reference implementation.
  go.mod: Pins the production Go toolchain and module dependencies.
  go.sum: Records production Go dependency checksums.
  main.go: Thin production entrypoint that supplies embedded build metadata and mapped exit codes to the Cobra tree.
  mirror.yaml: Mirror semantic version configuration using Git tags as the version source.
  package.json: npm native-bootstrap metadata and Go-backed contributor command aliases; it contains no RunX domain implementation.
  tsconfig.json: Preserved compiler configuration for the legacy TypeScript reference implementation.
  xdocs.config.toml: Legacy XDocs metadata configuration retained until tooling parity with xdocs.yaml is verified.
  xdocs.yaml: XDocs metadata configuration with automatic documentation-update mode and owned-tree scan exclusions.
  runx.yaml: RunX command catalog for this repository's supported development and validation commands.
documents:
  AGENTS.md: Repository rules for the Go CLI, XDocs, Mirror, validation, and protected release boundaries.
  CHANGELOG.md: Mirror-managed historical release record.
  CONTRIBUTING.md: Go validation, documentation, and protected release contribution workflow.
  DOCS.md: Complete production CLI, manifest, execution confirmation, lifecycle, agent, upgrade, installer, legacy-to-native migration, and 11-artifact reference.
  LICENSE.md: MIT license.
  README.md: Public installation with concise verified bootstrap commands, legacy-to-native migration, and cross-shell PATH guidance, plus manifest-v2, stable UID, numeric-index, safe confirmation, command, and Go release overviews.
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
Go packages. Structured manifest validation, exact command revelation, command execution, idempotent
bare-invocation agent bootstrap, cached lifecycle workers, agent resources,
verified upgrades, installers, and the protocol-v1
release matrix are enforced by Go tests and workflows.
