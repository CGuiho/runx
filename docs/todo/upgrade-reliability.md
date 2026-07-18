---
name: Upgrade Reliability
purpose: Define the expected outcome, constraints, and completion signals for task 2 in TODO.md.
description: Requirements for reliable self-upgrade, complete release listing, recovery commands, direct installers, and release-ready validation.
created: 2026-07-15
flags:
  - in-progress
  - needs-revalidation
tags:
  - cli
  - reliability
keywords:
  - runx upgrade
  - issue 12
  - issue 13
  - upgrade list
owner: runx-todo
---

# Upgrade Reliability

## Todo Index

- Task: `2. Make RunX Upgrades Reliable and Recoverable`
- Status: in progress
- Index: [TODO.md](../../TODO.md)
- Plan: [upgrade-reliability-implementation.md](../plans/upgrade-reliability-implementation.md)
- Design: [upgrade reliability design](../superpowers/specs/2026-07-15-upgrade-reliability-design.md)
- Validation: [upgrade-reliability.md](../validation/upgrade-reliability.md)

## Review Correction Status

An independent implementation review superseded the first validation pass.
Envelope fidelity, rollback state propagation, downgrade prevention, installer
failure classification, Ubuntu portability, and controlled executable installer
tests have been corrected. CLI and installer candidate verification is also
time-bounded so an invalid executable cannot hang an upgrade indefinitely.
The complete release gate still requires a rerun.

## Outcome

`runx upgrade` visibly plans, downloads, validates, replaces, and verifies the canonical native binary; `runx upgrade list` shows the complete SemVer-sorted release catalog; every result supplies exact-version recovery and safe stop commands; direct installers verify and roll back.

## Scope

### In scope

- GitHub release pagination and metadata normalization.
- Stable/prerelease SemVer ordering and compatible asset selection.
- Streaming text phases and equivalent buffered JSON.
- Transactional Windows and POSIX replacement with verification and rollback.
- Exact pinned recovery instructions after every attempt.
- Hardened PowerShell and POSIX installers.
- Focused automated tests, docs, bundled skill, and xdocs metadata.

### Out of scope

- Publishing, tagging, or changing Mirror-managed versions.
- Automatic process termination, background update daemons, package-manager integrations, or a new signing system.

## Acceptance Signals

- Success means the canonical executable already reports the selected exact version.
- Failure is explicit, exits nonzero, rolls back when mutation began, and prints recovery.
- Listing includes all pages, newest valid SemVer first, with channel/date/current/latest/asset data.
- Installers verify the exact requested or resolved version and restore the previous binary on failure.
- Repository checks and focused XDocs validation pass.

## External Trackers

- GitHub issue `CGuiho/runx#12`: upgrade reliability and complete version catalog.
- GitHub issue `CGuiho/runx#13`: exact-version recovery commands after every attempt.
