---
name: guiho-p-runx-uninstall
purpose: Uninstall the GUIHO RunX CLI.
description: Canonical uninstallation prompt for RunX.
created: 2026-08-22
owner: runx-prompts
flags:
  - agent-prompt
tags:
  - runx
  - uninstall
keywords:
  - guiho-p-runx-uninstall
---

Uninstall the GUIHO RunX CLI by running exactly one command, chosen by the host operating system:

- Windows (PowerShell):
  irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/uninstall.ps1 | iex
- macOS or Linux (POSIX sh):
  curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/uninstall.sh | sh

By default the uninstaller removes every RunX-owned artifact, including the stable
launcher, all versioned payloads, `$HOME/.guiho/runx/`, global configuration,
persistent data, managed agent-skill copies, the managed instruction block in the
project's `AGENTS.md`, and the project's `runx.yaml`. The shared `$HOME/.guiho/`
infrastructure, its `bin/` and `.temp/` directories, and the user `PATH` entry
are preserved.

To preview the plan without changing anything:

- Windows: `runx uninstall --dry-run` or the remote script with `-DryRun`
- macOS/Linux: `runx uninstall --dry-run` or `bash devops/uninstall.sh --dry-run`

To preserve configuration and persistent data:

```powershell
runx uninstall --preserve-config --preserve-data --yes
```

or

```bash
runx uninstall --preserve-config --preserve-data --yes
```

Noninteractive invocations without `--yes` fail without changing anything.
If uninstallation fails, create an issue at https://github.com/CGuiho/runx/issues/new
with the host OS, architecture, and uninstaller output.
