---
name: Windows Installer UTF-8 Validation
purpose: Record evidence for strict UTF-8 and idempotent RunX instruction reconciliation.
description: Covers the live 0.5.3 failure, focused regressions, full tests, typecheck, XDocs, release, and public installer verification.
created: 2026-07-21
flags:
  - testing
tags:
  - validation
  - installer
  - windows
keywords:
  - UTF-8
  - mojibake
  - AGENTS.md
  - PowerShell
owner: runx-validation
---

# Windows Installer UTF-8 Validation

## Root Cause

Windows PowerShell 5.1 used its legacy default encoding for `Get-Content` and
rewrote valid UTF-8 em dashes as mojibake. The installed-binary version probe
also started native background maintenance before installer reconciliation;
the native worker did not recognize the damaged marker and appended a second
canonical block.

## Local Evidence

| Check | Result |
| --- | --- |
| Live repository diff after the public 0.5.3 installer | Reproduced mojibake in existing XDocs guidance plus damaged and canonical RunX blocks. |
| Focused installer and agent-maintenance suites | Passed: 10 tests, 77 expectations. |
| Full `bun test --timeout 60000` | Passed: 60 tests, 413 expectations. |
| `bun run typecheck` | Passed. |
| Restored repository `AGENTS.md` | Byte-equivalent to committed intended state; no residual diff. |

## Regression Contract

- `Read-Utf8Text` rejects invalid UTF-8 and handles an optional BOM.
- `Write-Utf8Text` emits deterministic UTF-8 without a BOM.
- Reconciliation preserves Unicode and newline style and is byte-idempotent.
- Canonical, legacy, damaged, and duplicate RunX blocks converge to one.
- Installer version verification disables update and agent-maintenance workers.
- Native maintenance recognizes the damaged 0.5.3 marker.

## Release Evidence

Patch, CI, publish, exact fourteen assets, and public piped-installer evidence
will be appended during the authorized 0.5.4 release.

## Readiness

Ready for the Mirror patch release stage.

## References

- [Task specification](../todo/windows-installer-utf8.md)
