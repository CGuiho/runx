---
name: RunX Bash Installer Implementation Review
purpose: Review the GitHub issue 15 Bash conversion against all invocation and test surfaces.
description: Findings-first review of shebang, strict mode, recovery commands, docs, source-only piping, exact versions, and executable verification.
created: 2026-07-19
flags:
  - accepted
tags:
  - review
  - installer
  - bash
keywords:
  - RunX
  - issue 15
  - install.sh
owner: runx-implementation-reviews
---

# RunX Bash Installer Implementation Review

## Verdict

Accepted.

## Findings

The current installer no longer contained the historical invalid
`set -o pipefail` under `sh`, but the public shebang, README command, and
recovery generator still selected `sh`. That did not satisfy the issue title.

All canonical surfaces now select Bash. Strict mode is enabled, unset `SHELL`
is handled safely, and a real Bash runtime validates syntax, piped source-only
startup, exact stable/prerelease normalization, and executable version
verification.

No blocker, high, medium, or low finding remains.

## References

- [Task specification](../../todo/bash-installer.md)
- [Validation](../../validation/bash-installer.md)
- [GitHub issue #15](https://github.com/CGuiho/runx/issues/15)
