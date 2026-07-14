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

## 0.2.2 - 2026-07-14

### Changed

- Retried the npm trusted-publishing release after the initial `0.2.1` tag could not be created under the repository tag ruleset.

## 0.2.1 - 2026-07-14

### Added

- Added npm OIDC trusted publishing to the protected `production` workflow for `@guiho/runx@*` tags.
- Expanded CI to validate the complete 12-target native binary matrix.
- Added active GitHub rulesets protecting `main` and `@guiho/runx@*` release tags.

### Fixed

- Fixed the published `runx` launcher to use the compiled library entrypoint included in the npm package instead of relying on excluded source files.

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
