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

## 0.4.1 - 2026-07-18

### Fixed

- Fixed `runx upgrade` so its required `--version`, `--arch`, `--variant`,
  `--dry-run`, and `--format` flags route to the upgrade action instead of
  being hidden by the command group or intercepted as the root version flag.
- Fixed no-argument startup so a decoded cached update notice is printed before
  the exact RunX banner without waiting for network work.
- Fixed `runx agent prompt list --names` text output so it prints raw prompt
  names, one per line, while JSON mode remains structured.

## 0.4.0 - 2026-07-18

### Added

- Added complete Developer Context help, the singular agent skill/instruction/prompt namespace, a detached TypeBox-validated update cache, and the canonical `guiho-i-runx` prompt.
- Added a Node-compatible npm bootstrap, complete dual-tool direct installers, and exact verification for twelve native binaries plus two agent assets.

### Changed

- Rebuilt RunX around the complete RFC 0034 Bun, strict ESM TypeScript, raw Citty, and TypeBox contract.
- Configuration now resolves only through `--config`, effective-cwd `runx.yaml`, then `~/.guiho/runx/runx.yaml`.
- The startup banner, help modes, output streams, exit codes, upgrade catalog, release names, bundled skill, documentation, and workflows now follow the canonical RFC contract.

### Removed

- Removed the run command alias, root selector shorthand, `--file`, parent-directory discovery, plural agent namespace, arbitrary short flags, legacy platform asset names, and Bun-dependent npm launcher.

## 0.3.0 - 2026-07-18

### Added

- Added transactional executable replacement and automatic rollback for both Windows and POSIX updates with a bounded 10-second verification timeout.
- Added a complete GitHub release catalog parser that paginates all releases and implements SemVer prerelease precedence.
- Added structured JSON progress events and equivalent console phase reporting.
- Added exact-version recovery commands and process-stop commands to every upgrade outcome.
- Added hardened direct installers (`install.ps1`/`install.sh`) that verify the installed version and support rollback.

## 0.2.7 - 2026-07-15

### Added

- Added an interactive `runx init` command that collects the project name and
  scripts directory, previews the generated catalog, and confirms before writing.
- Added safe cancellation and explicit overwrite confirmation so initialization
  never leaves a partial manifest or replaces an existing catalog silently.

### Changed

- Initialized catalogs now use the SemVer-compatible manifest version `1.0.0`,
  configure `scripts.directory`, include the required `public` group, and begin
  with an empty command list.
- `runx init` creates only `runx.yaml`; the configured `scripts` directory is
  created later with the first real script.

## 0.2.6 - 2026-07-14

### Added

- Added Windows CI coverage for replacing a running executable, verifying the
  installed target version, cleaning the old image, and rolling back failures.

### Fixed

- Fixed `runx upgrade` on Windows so it replaces and verifies `runx.exe` before
  reporting success instead of returning a detached `scheduled: true` state.
- Failed Windows replacements now restore the previous executable, report a
  non-zero error, and clean `.new` and `.old` update files safely.

## 0.2.4 - 2026-07-14

### Added

- Added Citty as the runtime CLI builder for typed commands, nested routing,
  aliases, and generated usage.
- Added end-to-end command-contract coverage for help, version, catalog,
  agents, self-management, selector shorthand, and native/package behavior.

### Changed

- Replaced RunX's handwritten argument parser and top-level router with the
  complete Citty command tree while preserving `runx r` and `runx <selector>`.
- Updated public documentation, architecture, AGENTS instructions, and the
  bundled RunX skill for the Citty-owned interface.

### Fixed

- Fixed `runx -v` and `runx -h` so they work outside configured projects
  without attempting to discover `runx.yaml`.
- Unknown options and missing required selectors now report relevant command
  usage instead of unrelated manifest errors.

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
