---
name: RunX README
purpose: Introduce RunX and provide the supported RFC 0034 installation and quick-start workflow.
description: Public overview of the RunX YAML command catalog, native CLI, agent integration, and release channels.
created: 2026-07-12
flags:
  - public
tags:
  - cli
  - runx
keywords:
  - runx.yaml
  - native binary
  - command catalog
owner: runx
---

# RunX

RunX is a language-agnostic command catalog backed by a strict `runx.yaml`
file. It makes project operations discoverable and keeps execution explicit:
only `runx run <selector>` spawns a configured command.

## Install

PowerShell:

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
```

POSIX:

```sh
curl --proto '=https' --tlsv1.2 -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | sh
```

The direct installers select and verify a native asset, install `runx` on PATH,
install `guiho-s-runx` into both global agent-tool directories, and reconcile
RunX instructions in the current project.

The npm package is a Node-compatible bootstrap that downloads and delegates to
the version-matched native binary. Bun is not required for npm installation or
first execution.

## Quick Start

```text
runx init
runx check
runx list
runx describe app-dev
runx run app-dev --dry-run
runx run app-dev
```

Configuration resolves in this exact order:

1. `--config <path>`;
2. `<effective-cwd>/runx.yaml`;
3. `~/.guiho/runx/runx.yaml`.

RunX never searches parent directories.

## Command Catalog

```text
runx
├── list
├── describe <selector>
├── run <selector>
├── check
├── init
├── agent
│   ├── skill install|uninstall|update|list|show
│   ├── instruction apply|remove|update|show
│   └── prompt list|show
├── upgrade
│   ├── check
│   └── list
└── uninstall
```

Every scope supports `--help`, `-h`, `--help-tree`,
`--help-tree-depth <positive-integer>`, and `--help-docs`. Only root
`--version` also has `-v`.

See [DOCS.md](DOCS.md) for configuration, output, exit codes, agent resources,
upgrade behavior, installers, npm bootstrap, and the exact release matrix.
