---
name: RunX README
purpose: Introduce RunX and give users the shortest installation and command-catalog path.
description: Public overview, direct installation commands, quick manifest example, and documentation links.
created: 2026-07-12
flags: []
tags:
  - documentation
  - cli
keywords:
  - runx
  - readme
owner: runx
---

# RunX

RunX is a local, language-agnostic command catalog. Put documented commands in
`runx.yaml`; use `runx list` to discover them and `runx run` to execute one.
It works for Bun, Python, Go, shell, and mixed repositories because it executes
local shell commands rather than package-manager scripts.

## Install

RunX is distributed as a native executable from GitHub Releases. It is not
published to npm during the alpha.

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
```

```sh
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | sh
```

## Quick Start

Create `runx.yaml` in a project root:

```yaml
version: 1
groups:
  development:
    summary: Local development commands.
commands:
  - uid: app-dev
    id: dev
    group: development
    summary: Start the application locally.
    description: Starts the application in local development mode until interrupted.
    command: bun run dev
    shell: bash
    tags: [local, watch]
```

Then use RunX:

```text
runx                         Show the home page and usage.
runx list                    List commands in the nearest manifest.
runx describe app-dev        Explain one command without execution.
runx run app-dev --dry-run   Inspect the execution plan.
runx run app-dev             Execute the command.
runx r app-dev               Use the short run alias.
```

Use `confirm: always` for commands that should require an explicit `--yes`.
Use `runx list --format json` and stable UIDs for agent automation.

## Documentation

- [Full CLI and manifest reference](DOCS.md)
- [Contributing](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [Changelog](CHANGELOG.md)

## License

[MIT](LICENSE.md)
