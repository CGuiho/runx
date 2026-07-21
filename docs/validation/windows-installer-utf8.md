---
name: Windows Installer UTF-8 Validation
purpose: Record evidence for strict UTF-8 and idempotent RunX instruction reconciliation.
description: Covers the live 0.5.3 failure, focused regressions, full tests, typecheck, XDocs, release, and public installer verification.
created: 2026-07-21
flags:
  - validated
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

- Mirror applied `0.5.4` and pushed commit
  `02e213872b0695b5799f106831d791cce5f9eb21` plus tag
  `@guiho/runx@0.5.4`.
- [CI run 29865746436](https://github.com/CGuiho/runx/actions/runs/29865746436)
  passed on attempt 2 for Linux and Windows. The initial Linux public-installer
  probe ran before the protected publish job was approved; it passed unchanged
  after publication.
- [Publish run 29865749514](https://github.com/CGuiho/runx/actions/runs/29865749514)
  completed successfully.
- The public release contains exactly twelve native binaries and the two `.md`
  agent assets. Its body contains only the `0.5.4` changelog section.
- From `C:\Users\crist`, the public `irm .../install.ps1 | iex` command installed
  and verified `0.5.4` twice.
- After both installs, `C:\Users\crist\AGENTS.md` remained 10,126 bytes with
  SHA-256 `ea7ced628c146d88cf7fba85a3560ba4184b8675279d07e269ddb78d52428e5d`,
  one canonical RunX marker, valid UTF-8, and no BOM.
- Zero RunX processes remained after the second installation.

## Readiness

Validated and publicly released as RunX `0.5.4`.

## References

- [Task specification](../todo/windows-installer-utf8.md)
