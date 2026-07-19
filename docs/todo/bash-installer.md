---
name: Use Bash For The RunX Installer
purpose: Track the canonical shell correction requested by GitHub issue 15.
description: Records Bash shebang, strict mode, documented invocation, recovery generation, executable tests, and delivery evidence.
created: 2026-07-19
flags:
  - completed
tags:
  - todo
  - installer
  - bash
keywords:
  - RunX
  - issue 15
  - install.sh
  - pipefail
owner: runx-todo
---

# Use Bash For The RunX Installer

## Todo Index

- Task: `4. Use Bash For The RunX Installer`
- Status: completed
- Index: [TODO.md](../../TODO.md)
- External: [CGuiho/runx issue #15](https://github.com/CGuiho/runx/issues/15)

## Outcome

RunX declares and invokes the Linux/macOS installer as Bash everywhere. The
installer uses Bash strict mode and passes syntax, piped startup, exact-version,
and executable-version checks under a real Bash runtime.

## Completion Signals

- No canonical installer or recovery command pipes `install.sh` to `sh`.
- The script starts with `#!/usr/bin/env bash` and `set -euo pipefail`.
- Shell-profile detection remains safe when `SHELL` is unset.
- Git Bash and Ubuntu-compatible tests exercise the script, not only its text.

## References

- [Implementation review](../reviews/implementation/bash-installer-review.md)
- [Validation](../validation/bash-installer.md)
