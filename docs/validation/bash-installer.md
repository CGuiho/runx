---
name: RunX Bash Installer Validation
purpose: Record direct verification evidence for GitHub issue 15.
description: Captures Bash syntax, piped startup, exact versions, recovery, executable verification, full suite, native assets, XDocs, and Git checks.
created: 2026-07-19
flags:
  - validated
tags:
  - validation
  - installer
  - bash
keywords:
  - RunX
  - issue 15
  - bash -n
  - 14 assets
owner: runx-validation
---

# RunX Bash Installer Validation

## Summary

The canonical RunX Linux/macOS installer and recovery command use Bash. No
validation blocker remains.

## Commands Run

| Command or check | Result |
| --- | --- |
| `bun run typecheck` | Passed |
| focused recovery and installer tests | Passed: 7 tests, 45 assertions |
| Git Bash `bash -n devops/install.sh` | Passed |
| piped source-only Bash fixture | Passed with exit 0 and empty stderr |
| exact stable/prerelease normalization fixture | Passed |
| Bash executable-version verification fixture | Passed |
| `bun test --timeout 30000` | Passed: 53 tests, 363 assertions |
| build, native matrix, exact assets | Passed: library, single binary, twelve native targets, exactly fourteen assets |
| canonical `sh` reference scan | Passed: zero matches |
| strict XDocs metadata and doctor | Passed for source, devops, and docs: zero errors and zero warnings |
| `git diff --check` | Passed |

## Contract Evidence

- The script declares `#!/usr/bin/env bash` and `set -euo pipefail`.
- README and generated recovery commands pipe to `bash`, never `sh`.
- Tests run with Git for Windows Bash locally and resolve the ordinary Bash
  executable on Linux/macOS CI.
- The Bash function fixture verifies full prerelease identifiers and executes
  the installed-version check against a real executable.

## References

- [Implementation review](../reviews/implementation/bash-installer-review.md)
- [Task specification](../todo/bash-installer.md)
