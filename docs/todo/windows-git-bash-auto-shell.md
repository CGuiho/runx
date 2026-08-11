---
name: Respect Git Bash For Windows Automatic Shell Execution
purpose: Fix the Windows shell boundary that changes MSYS path arguments when a RunX command is delegated through cmd.exe.
description: Defines issue 47 follow-up behavior for caller-aware shell:auto resolution, explicit-shell precedence, safe fallback, tests, and patch delivery.
created: "2026-08-11"
flags:
  - testing
tags:
  - todo
  - windows
  - execution
keywords:
  - Git Bash
  - shell auto
  - MSYSTEM
  - cmd.exe
  - issue 47
owner: runx-todo
---

# Respect Git Bash For Windows Automatic Shell Execution

## Problem

RunX 0.12.0 resolves `shell: auto` to `cmd.exe` solely because the binary is
running on Windows. When RunX itself was launched from Git Bash, that changes
the command boundary: Git Bash would translate MSYS-style native-process path
arguments such as `/c/GUIHO` to `C:/GUIHO`, while `cmd.exe` passes the original
token unchanged.

The issue was reproduced without secrets using native Node as an argument
probe. The real GUIHO Core catalog was inspected only through `runx reveal`;
no environment file, key file, or environment value was read.

## Outcome

On Windows, `shell: auto` preserves a verifiable Git Bash/MSYS caller context
by selecting the caller's resolvable Bash executable. Other Windows caller
contexts keep the existing `cmd.exe` default. Explicit `shell: cmd`,
`powershell`, `bash`, and `sh` values remain authoritative.

## Acceptance

- Automatic resolution is factored so Windows behavior is deterministic and
  testable on every CI platform.
- A recognized MSYS/Git Bash caller marker plus a resolvable Bash selects Bash.
- Missing or unverifiable Bash falls back to `cmd.exe`; automatic resolution
  never selects the Windows WSL launcher merely because it is named `bash`.
- An ordinary Windows context without the caller marker remains on `cmd.exe`.
- Explicit shell values bypass automatic caller inference.
- Child-argument transport, exit-code propagation, dry runs, confirmation, and
  `runx reveal` remain unchanged.
- Behavior-based tests cover the `/c/GUIHO` path boundary without using any
  secret-bearing file.
- Public docs, bundled/embedded RunX skill guidance, changelog, TODO state, and
  XDocs metadata describe the exact contract.
- Exact-head review, CI, protected integration, and a Mirror-managed patch
  release complete before issue 47 is closed.

## External

- [GitHub issue 47](https://github.com/CGuiho/runx/issues/47)
- [RunX 0.12.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.12.0)
