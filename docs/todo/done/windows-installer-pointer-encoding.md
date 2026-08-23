---
name: Fix Windows Installer Pointer Encoding
purpose: Define the incident correction and acceptance contract for TODO task 24.
description: Requires PowerShell 5 and 7 installation paths to write a BOM-free current.json and pass real Windows launcher activation tests.
created: 2026-08-23
flags:
  - completed
tags:
  - installer
  - windows
  - powershell
keywords:
  - current.json
  - UTF-8 BOM
  - stable launcher
owner: runx-todo-done
---

# Fix Windows Installer Pointer Encoding

## Status

- State: completed
- Created: `2026-08-23T12:18:30+02:00`
- Updated: `2026-08-23T13:38:20+02:00`
- Parent task: RunX TODO 23, GUIHO CLI Convention 0001 compliance

## Incident

Windows PowerShell 5.1 interprets `Set-Content -Encoding utf8` as UTF-8 with a
byte-order mark. The RunX installer used that command to stage `current.json`.
The stable launcher strictly decoded JSON and rejected the leading bytes with
`invalid character 'ï' looking for beginning of value`, so installation rolled
back after all release artifacts had already downloaded and verified.

Independent reproduction also found that release selection prefixed the
already-prefixed target `runx-windows-amd64` with `runx-payload-`, yielding the
nonexistent `runx-payload-runx-windows-amd64.exe`. The post-0.14.6 live smoke
then exposed a second upgrade defect: the installed 0.14.3–0.14.6 lifecycle
engine requires `checksums.txt` to contain its own SHA-256 digest, which no
valid checksum manifest can do.

## Plan Unit

`windows-installer-pointer-encoding`: make pointer creation explicitly UTF-8
without a BOM and prove the actual PowerShell installer can activate the real
Windows launcher.

A separate architecture or implementation plan is unnecessary because this is
an isolated correction to the already-approved stable-launcher transaction in
the Convention 0001 migration. Requirements, architecture, and plan-writing
phases are waived in favor of the existing Convention 0001 contract and task
23 artifacts.

## Acceptance Criteria

1. `install.ps1` writes staged `current.json` as BOM-free UTF-8 under Windows
   PowerShell 5.1 and PowerShell 7.
2. The launcher can strictly decode the resulting pointer and report the exact
   installed SemVer.
3. A Windows-native integration test executes `install.ps1` against a complete
   local release fixture rather than merely parsing or searching script text.
4. The test verifies a clean install and reinstall/repair of the same version,
   including pointer bytes, launcher dispatch, and raw `--version` output.
5. CI and release publication run the Windows installer integration test before
   any release is considered successful.
6. Existing installation failure leaves the prior active pointer and launcher
   usable.

## Implementation Milestone

- `install.ps1` writes pointer JSON with `System.Text.UTF8Encoding($false)` and
  hashes files with the module-independent .NET SHA-256 API.
- Pointer readers tolerate only the legacy UTF-8 BOM before strict JSON decode.
- Protocol-v1 asset selection removes the existing `runx-` prefix before
  constructing the payload name.
- Release builds include checksummed compatibility payload/launcher aliases for
  every supported target, allowing defective 0.14.3–0.14.5 clients to enter the
  corrected release.
- The Windows integration test passed clean installation, injected
  post-activation rollback, and same-version reinstall with native binaries.
- A real-network installer-driven upgrade from 0.14.3 to published 0.14.5
  passed in an isolated Windows home and produced BOM-free pointer bytes.
- The mandatory post-0.14.6 live `runx upgrade` smoke failed closed with
  `checksum entry missing for checksums.txt`; this became the 0.14.7 follow-up.
- The public 0.14.7 clean install passed, then its exact same-version upgrade
  smoke found a Windows lock: the transaction tried to overwrite the active
  immutable payload. The 0.14.8 follow-up returns verified `up-to-date` without
  mutation for an already-active exact target.

## Validation

- Public clean installation of 0.14.8 succeeded with verified SHA-256 digests,
  BOM-free `current.json`, stable launcher dispatch, and raw `0.14.8` output.
- Public `runx upgrade` from installed 0.14.7 to 0.14.8 completed synchronously,
  activated the new immutable payload, and verified through the launcher.
- Public `runx upgrade --version 0.14.8` returned verified `up-to-date` without
  mutating or locking the active payload.
- Release workflow `32636794191` passed the native Windows install, rollback,
  reinstall, different-version upgrade, same-version idempotence, Go, release
  matrix, exact remote assets, and public installer gates.

## Release Decision

Completed through Mirror-managed patches 0.14.6, 0.14.7, and final corrective
release 0.14.8. The canonical installer is the one-time recovery path for
protocol-v1 0.14.3–0.14.6 because their defect exists inside already-installed
binaries.
