---
name: RunX Changelog
purpose: Record externally meaningful RunX release changes.
description: Tracks release-level changes managed with Mirror.
created: 2026-07-12
flags: []
tags:
  - releases
  - changelog
keywords:
  - runx
  - changelog
owner: runx
---

# Changelog

## 0.2.0 - 2026-07-12

### Added

- Added `publish.yml` GitHub Actions workflow to compile and upload release binaries to GitHub releases on version tag pushes.

### Changed

- Updated `devops/build-binaries.ts` to compile the full 12 binary target matrix in parallel using optimized build flags.
- Replaced `devops/install.sh` and `devops/install.ps1` with robust, parameter-driven installers.

## 0.1.1 - 2026-07-12

### Changed

- Finalized task completion and validation evidence for the initial alpha release.

## 0.1.0 - 2026-07-12

### Added

- First alpha implementation of the RunX command-catalog CLI.
