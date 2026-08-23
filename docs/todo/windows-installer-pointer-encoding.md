---
name: Fix Windows Installer Pointer Encoding
purpose: Define the incident correction and acceptance contract for TODO task 24.
description: Requires PowerShell 5 and 7 installation paths to write a BOM-free current.json and pass real Windows launcher activation tests.
created: 2026-08-23
flags:
  - testing
tags:
  - installer
  - windows
  - powershell
keywords:
  - current.json
  - UTF-8 BOM
  - stable launcher
owner: runx-todo
---

# Fix Windows Installer Pointer Encoding

## Status

- State: testing
- Created: `2026-08-23T12:18:30+02:00`
- Updated: `2026-08-23T12:39:29+02:00`
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
  `checksum entry missing for checksums.txt`; this became the 0.14.7 follow-up
  unit rather than being misreported as successful.

## Validation

The installer gates and 0.14.6 publication passed, but the post-publication
whole-release upgrade smoke found the checksum-manifest defect. Completion is
blocked until 0.14.7 passes the new native whole-release integration gate and a
real published 0.14.7 exact-version whole-release upgrade demonstrates the corrected path.

## Release Decision

Patch required. Do not publish until both clean install and reinstall pass on
Windows PowerShell with the actual launcher and payload.
