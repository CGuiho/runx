---
name: RunX
purpose: Introduce the native RunX command-catalog CLI and its safe operating contract.
description: Installation, manifest-v2 workflow, command discovery, execution, agent resources, upgrades, and Go release targets.
created: 2026-07-12
owner: runx
flags: []
tags:
  - cli
  - go
  - cobra
keywords:
  - runx
  - runx.yaml
  - command catalog
---

# RunX

RunX is a native Go/Cobra CLI for documented, language-agnostic command
catalogs. A project owns one explicit `runx.yaml`; RunX validates and describes
catalog commands without executing them, and only `runx run` starts a configured
command.

## Install

Linux and macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash
```

Windows PowerShell:

```powershell
& ([scriptblock]::Create((Invoke-RestMethod 'https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1')))
```

Both installers select the canonical target, download `checksums.txt`, verify
SHA-256 before replacement, install the bundled skill into both supported agent
locations, and verify `runx --version`.

## Start

```bash
runx init
runx check --format json
runx list --format json
runx describe <uid>
runx run --dry-run <uid>
runx run --yes <uid> -- <child arguments...>
```

RunX options precede the selector. Every token after the selector belongs to
the child and is forwarded without reinterpretation.

## Manifest v2

```yaml
version: "2.0.0"
namespace: "example"
scripts:
  directory: "scripts"
commands:
  - uid: "test-command"
    id: "test"
    summary: "Run tests."
    description: "Run the complete project test suite."
    command: "go test ./..."
    confirm: "never"
```

Manifests are strictly decoded: unknown fields, invalid identifiers, unsafe
paths, invalid shells, duplicate identities, non-reciprocal child catalogs,
and unsupported manifest versions fail closed. Configuration precedence is:

1. `--config <path>`;
2. effective-cwd `runx.yaml`;
3. `~/.guiho/runx/runx.yaml`.

RunX never searches parent directories.

## Developer Context

Every command scope supports `-h`/`--help`, `--help-tree`,
`--help-tree-depth <positive-integer>`, and `--help-docs`. Root additionally
supports `-v`/`--version`. The tree and Markdown are generated from the live
Cobra commands.

## Agent And Upgrade Commands

```text
runx agent skill install|uninstall|update|list|show
runx agent instruction apply|remove|update|show
runx agent prompt list|show
runx upgrade
runx upgrade check
runx upgrade list
runx uninstall --dry-run
```

Bare `runx` first installs its embedded skill in both global tool locations and
idempotently reconciles a bounded RunX instruction block in the repository root.
Both `AGENTS.md` and `CLAUDE.md` are updated when both exist; otherwise the one
that exists is used, or `AGENTS.md` is created. Existing content and line endings
are preserved, malformed markers fail safely, and no catalog command or network
request runs during bootstrap. Help, version, agent-management, uninstall, and
non-repository paths do not perform repository bootstrap. Other foreground
startup reads only the local cache and starts bounded detached workers where
appropriate.
Self-upgrades verify published checksums and preserve the embedded build target,
including ARMv6 versus ARMv7.

## Build

```bash
go test ./...
go vet ./...
go build ./...
go run devops/build-binaries.go --version 0.8.0 --commit <commit> --build-date 2026-07-26T00:00:00Z
go run devops/verify-release-assets.go
```

The release contract is exactly 11 artifacts: eight pure-Go executables
(Linux AMD64/ARM64/ARMv7/ARMv6, Darwin AMD64/ARM64, Windows AMD64/ARM64),
`guiho-s-runx.zip`, `guiho-i-runx.md`, and `checksums.txt`. AMD64 V2/V3/V4
variants are not part of the contract.

See [DOCS.md](DOCS.md) for the complete behavior and safety reference.
