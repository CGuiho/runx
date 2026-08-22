---
name: runx-install
purpose: Install the GUIHO RunX CLI.
description: Canonical installation prompt for RunX.
created: 2026-08-22
owner: runx-prompts
flags:
  - agent-prompt
tags:
  - runx
  - install
keywords:
  - runx-install
---

Install the GUIHO RunX CLI by running exactly one command, chosen by the host operating system:

- Windows (PowerShell):
  irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
- macOS or Linux (POSIX sh):
  curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | sh

To install an exact version or a release channel, add --version 'X.Y.Z' or --channel '<name>'
to the install.sh command, or -Version 'X.Y.Z' or -Channel '<name>' to the install.ps1 command.
Never install RunX through a package manager (npm, Bun, pip, or similar); no such package
exists or is supported.

When the installer finishes, verify with `runx --version`.

If the installer fails, fall back to a manual binary installation:
1. Download the asset matching the host platform from
   https://github.com/CGuiho/runx/releases/latest - for example
   `runx-windows-amd64.exe` or `runx-linux-amd64` - and verify its SHA-256
   digest against the release's `checksums.txt`.
2. Place the binary at `$HOME/.guiho/bin/runx` (`runx.exe` on Windows),
   creating the directory when needed.
3. On macOS and Linux, mark it executable with `chmod +x`. Ensure `$HOME/.guiho/bin`
   is on the `PATH`, then verify with `runx --version`.

If installation still fails, create an issue at https://github.com/CGuiho/runx/issues/new
with the host OS, architecture, attempted version or channel, installer output, and
`runx --version` output.

When the CLI is installed, finish by running `runx init`, which applies the permanent
agent instruction that keeps every future agent aware of the CLI.
