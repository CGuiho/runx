---
name: Fix Windows Installer Pointer Encoding
purpose: Define the incident correction and acceptance contract for TODO task 24.
description: Requires PowerShell 5 and 7 installation paths to write a BOM-free current.json and pass real Windows launcher activation tests.
created: 2026-08-23
flags:
  - in-progress
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

- State: in progress
- Created: `2026-08-23T12:18:30+02:00`
- Updated: `2026-08-23T12:18:30+02:00`
- Parent task: RunX TODO 23, GUIHO CLI Convention 0001 compliance

## Incident

Windows PowerShell 5.1 interprets `Set-Content -Encoding utf8` as UTF-8 with a
byte-order mark. The RunX installer used that command to stage `current.json`.
The stable launcher strictly decoded JSON and rejected the leading bytes with
`invalid character 'ï' looking for beginning of value`, so installation rolled
back after all release artifacts had already downloaded and verified.

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

## Validation

Pending implementation. Required checks include focused Windows tests, the
complete Go suite, vet, builds, release-matrix verification, XDocs validation,
and a clean-home Windows installer smoke using the published release candidate.

## Release Decision

Patch required. Do not publish until both clean install and reinstall pass on
Windows PowerShell with the actual launcher and payload.
