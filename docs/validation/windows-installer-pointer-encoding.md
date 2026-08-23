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

## Corrected Gate Finding

The first PR CI run caught that GitHub's Windows PowerShell did not auto-load
`Get-FileHash`. The installer was corrected to use the module-independent .NET
SHA-256 API; workflow 32635016462 then passed on Windows.

## Readiness

Ready for the Mirror-managed 0.14.6 patch release. After publication, validate
a real 0.14.3–0.14.5 `runx upgrade` through the compatibility aliases before
marking the task complete.

## Evidence

- Validation comment: https://github.com/CGuiho/runx/pull/59#issuecomment-5385635476
- CI: https://github.com/CGuiho/runx/actions/runs/32635016462
- PR: https://github.com/CGuiho/runx/pull/59

No production deployment, promotion, traffic, DNS, database, or secret mutation
occurred.
