---
name: RunX Windows Installer And Upgrade Validation
purpose: Record exact-head and CI evidence for the Windows installer and protocol-v1 upgrade corrections.
description: Captures native PowerShell install, rollback, reinstall, Go, release-matrix, XDocs, and PR gate results for RunX PR 59.
created: 2026-08-23
flags:
  - passed
tags:
  - validation
  - windows
  - release
keywords:
  - PR 59
  - PowerShell 5.1
  - compatibility aliases
owner: runx-validation
---

# RunX Windows Installer And Upgrade Validation

## Scope

Validated implementation head
`538476c9176f8ffe047a0df771498a5c2ba3e2fa`, merged through PR
[CGuiho/runx#59](https://github.com/CGuiho/runx/pull/59) as
`a8eb075304a5b1ad36fdf583b430786f560c0f5d`.

## Results

| Check | Result |
| --- | --- |
| `go test -count=1 ./...` | Passed locally and on Ubuntu/Windows CI |
| Native Windows PowerShell clean install | Passed |
| Injected post-activation rollback | Passed; launcher and pointer remained usable |
| Same-version Windows reinstall | Passed |
| BOM-free pointer bytes and strict launcher dispatch | Passed |
| `go vet ./...` and `go build ./...` | Passed locally and on CI |
| Release matrix | 40 files; 38 manifest artifacts; all checksums, digests, and alias byte-equality checks passed |
| `xdocs meta . --documents --strict` | Passed |
| `xdocs tree` | Passed |
| `xdocs doctor .` | Passed with zero warnings/errors |
| `runx check --format json` | Passed |
| PR workflow 32635016462 | Ubuntu, Windows, and release-contract jobs passed |
| Real-network installer-driven 0.14.3 → 0.14.5 | Passed in an isolated Windows home with BOM-free pointer |

## Corrected Gate Findings

1. PR 59 CI caught that GitHub's Windows PowerShell did not auto-load
   `Get-FileHash`; the installer moved to the module-independent .NET SHA-256
   API.
2. The post-0.14.6 public upgrade smoke caught an impossible self-checksum
   requirement; PR 60 excludes only `checksums.txt` itself and adds a native
   different-version whole-release test.
3. The public 0.14.7 exact-version smoke caught a Windows lock on the active
   payload; PR 61 returns verified `up-to-date` before mutation.

## Public Acceptance

- Public 0.14.8 clean install: passed.
- Public 0.14.7 → 0.14.8 `runx upgrade`: outcome `upgraded`, launcher
  verification `ok`, active pointer `0.14.8`, previous pointer `0.14.7`.
- Public `runx upgrade --version 0.14.8`: outcome `up-to-date`, verification
  `ok`, no payload mutation.
- Actual user installation repaired through the canonical installer and raw
  `runx --version` returned `0.14.8`.
- Public 0.14.9 injected installer failure rolled back with nonzero status and
  zero leftover staging; immediate recovery install succeeded with zero
  leftover staging.
- Actual user installation upgraded synchronously 0.14.8 → 0.14.9; exact
  0.14.9 returned verified `up-to-date`.

## Readiness

Passed and complete at release `runx/v0.14.9` with 40 non-draft release assets.
Protocol-v1 0.14.3–0.14.6 installations require the canonical installer once
because their checksum defect is inside the already-installed executable.

## Evidence

- PR 59 validation: https://github.com/CGuiho/runx/pull/59#issuecomment-5385635476
- PR 60 validation: https://github.com/CGuiho/runx/pull/60#issuecomment-5385707301
- PR 61 gate: https://github.com/CGuiho/runx/pull/61#issuecomment-5385759966
- 0.14.8 workflow: https://github.com/CGuiho/runx/actions/runs/32636794191
- Final cleanup workflow: https://github.com/CGuiho/runx/actions/runs/32637716640
- Final release: https://github.com/CGuiho/runx/releases/tag/runx%2Fv0.14.9

No production deployment, promotion, traffic, DNS, database, or secret mutation
occurred.
