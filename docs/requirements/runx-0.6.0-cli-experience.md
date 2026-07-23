---
name: RunX 0.6.0 CLI Experience Requirements
purpose: Define the approved user-visible requirements for GitHub issues 23, 24, and 25
description: Specifies the welcome window, simplified public installation command, safe child argument forwarding, output behavior, and delivery boundaries for RunX 0.6.0.
created: 2026-07-22
flags:
  - approved
tags:
  - requirements
  - cli
keywords:
  - welcome window
  - argument forwarding
  - installer
owner: runx-requirements
---

# RunX 0.6.0 CLI Experience Requirements

## Scope

This release implements [issue 23](https://github.com/CGuiho/runx/issues/23),
[issue 24](https://github.com/CGuiho/runx/issues/24), and
[issue 25](https://github.com/CGuiho/runx/issues/25). The proposed Go rewrite
in issue 22 is explicitly out of scope.

## Welcome Window

- Bare `runx` renders a deterministic bordered welcome containing product,
  organization, platform, architecture, version, and `runx --help` guidance.
- Help, version, JSON, Markdown, and command output remain clean.
- A valid cached newer stable release adds a two-line warning after the welcome.
- Startup performs no foreground network request.

## Installation

The canonical POSIX README command is exactly:

```sh
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash
```

The simpler bootstrap does not weaken downstream HTTPS, checksum, asset,
transactional replacement, or executable verification.

## Execution Arguments

The public execution grammar is:

```text
runx run [RunX options] <selector> [--] <child arguments...>
```

- RunX options occur before the selector.
- Every token after the selector belongs to the child.
- One immediate delimiter is removed; later delimiters are preserved.
- Forwarded values preserve order and cannot become unintended shell syntax.
- `--dry-run` reports a lossless argument array and never spawns.
- Executed commands preserve the child exit code.

## Delivery

- The publish workflow has no GitHub Environment approval gate.
- All build, test, release-note, and exact-fourteen-asset gates remain.
- Release notes contain only the `0.6.0` changelog section.
- Issues close only after public release acceptance.
