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

# GUIHO RunX

RunX is a local, language-agnostic command catalog. Put documented commands in
`runx.yaml`; use `runx list` to discover them and `runx run` to execute one.
It works for Bun, Python, Go, shell, and mixed repositories because it executes
local shell commands rather than package-manager scripts.

## Install

RunX is distributed as a native executable from GitHub Releases and as the
public `@guiho/runx` package on npm. The npm launcher requires Bun.

```sh
bun add --global @guiho/runx
```

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
```

```sh
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | sh
```

## Quick Start

Start with the interactive initializer:

```text
runx init
```

It writes only `runx.yaml`—not an empty `scripts/` directory. The resulting
catalog starts with a SemVer `1.x` manifest version, a configurable scripts
directory, the required `public` group, and no commands. Add commands directly
or ask an agent to add a script under the configured directory and its command
entry together.

For example, a catalog with a development command can look like this:

```yaml
version: "1.0.0"
project:
  name: example
scripts:
  directory: scripts
groups:
  public:
    summary: Default public project commands.
  development:
    summary: Local development commands.
commands:
  - uid: app-dev
    id: dev
    group: development
    summary: Start the application locally.
    description: Starts the application in local development mode until interrupted.
    command: bash scripts/dev.sh
    shell: bash
    tags: [local, watch]
```

Then use RunX:

```text
runx                         Show the home page and usage.
runx -h                      Show Citty-generated command help.
runx -v                      Show the installed version without loading a manifest.
runx init                    Interactively create an empty runx.yaml catalog.
runx list                    List commands in the nearest manifest.
runx describe app-dev        Explain one command without execution.
runx run app-dev --dry-run   Inspect the execution plan.
runx run app-dev             Execute the command.
runx r app-dev               Use the short run alias.
runx app-dev                 Use the human selector shorthand.
```

Use `confirm: always` for commands that should require an explicit `--yes`.
Use `runx list --format json` and stable UIDs for agent automation.
RunX uses Citty for argument parsing, aliases, command routing, and ordinary
usage text, so conventional `-h`/`--help` and `-v`/`--version` work outside a
configured project.

## Upgrade and Recovery

`runx upgrade` prints the selected current/target versions, operating system,
architecture, binary, canonical path, and exact download URL before it starts
the download. It then streams `Downloading`, `Validating`, `Replacing`, and
`Verifying`. Success means the canonical executable already reports the exact
target version; a failed verification restores the previous executable.

Every outcome—including already current, dry run, and failure—prints a
copyable direct-install command pinned to the selected version and a separate
process-stop command. Use `runx upgrade --format json` for the same plan,
ordered phase events, nullable plan, verified result, stable error code, and
recovery data as one JSON document. A successful automatic rollback uses the
nonzero `rolled-back` outcome and reports the restored installed version.

```text
runx upgrade check          Check the latest stable release.
runx upgrade --dry-run      Print the plan and pinned recovery without mutation.
runx upgrade                Install and verify the latest stable release.
runx upgrade list           List every stable and prerelease version, newest first.
runx upgrade list --format json
```

## Documentation

- [Full CLI and manifest reference](DOCS.md)
- [Contributing](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [Changelog](CHANGELOG.md)

## License

[MIT](LICENSE.md)
