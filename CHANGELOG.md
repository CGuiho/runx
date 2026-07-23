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

## 0.7.2 - 2026-07-23

### Fixed

- Gave the Windows child-process forwarding and PowerShell installer
  verification tests explicit fifteen-second budgets so normal runner startup
  variance no longer causes false five-second CI failures.

## 0.7.1 - 2026-07-23

### Fixed

- Removed the release race that made ordinary CI compare unreleased source to
  the concurrently changing latest public release.
- Moved exact-version public Linux installer acceptance into Publish after the
  tagged GitHub Release contains exactly fourteen assets and before npm
  publication.

## 0.7.0 - 2026-07-23

### Added

- Added recursive colocated command groups and explicit local or GitHub child
  catalog mounts with renameable namespace aliases and reciprocal parents.
- Added canonical nested selectors plus local/foreign catalog provenance to
  list, describe, check, dry-run, and execution data.

### Changed

- Replaced manifest v1 `project.name` and split top-level `groups` with required
  manifest v2 `namespace` and one recursive `commands` tree.
- Changed `runx init` to create strict manifest v2 catalogs.

### Security

- Restricted foreign catalogs to HTTPS GitHub blob/raw URLs with redirect
  validation, ten-second timeouts, one-MiB limits, cycle/depth protection, no
  persistent cache, and configuration-error mapping.

## 0.6.1 - 2026-07-23

### Fixed

- Restored the `production` GitHub Environment identity required by the npm
  trusted-publisher record so tag releases authenticate through OIDC and reach
  npm after the GitHub Release and exact fourteen-asset checks succeed.

## 0.6.0 - 2026-07-22

### Added

- Added a deterministic bordered RunX welcome with product, GUIHO, platform,
  architecture, version, help, and validated cached-update information.
- Added lossless post-selector child argument and subcommand forwarding across
  Bash, sh, PowerShell, cmd, and automatic shell selection.
- Added forwarded argument arrays to text and JSON dry-run output.

### Changed

- Simplified the canonical POSIX installation command to
  `curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash`
  while retaining installer-internal validation and transactional replacement.
- Removed the GitHub Environment approval gate from tag publishing while
  retaining typecheck, test, build, release-note, npm OIDC, and exact
  fourteen-asset verification gates.
- Established that RunX-owned options precede the selector and every
  post-selector token belongs to the child command.

### Fixed

- Prevented child flags such as `-v`, `--help`, `--yes`, and `--dry-run` from
  being consumed or dropped by RunX routing.
- Revalidated cached latest versions against the running SemVer before showing
  an update warning.

## 0.5.4 - 2026-07-21

### Fixed

- Preserved existing UTF-8 project instructions during Windows installation
  and emitted deterministic UTF-8 without a byte-order mark.
- Prevented installer verification from racing background reconciliation and
  recovered damaged, legacy, or duplicate RunX managed blocks idempotently.

## 0.5.3 - 2026-07-21

### Fixed

- Coalesced automatic update checks behind a per-user lease so concurrent RunX
  invocations launch at most one detached update worker.
- Bounded the detached update worker with stale-lease recovery and a hard
  network deadline, preventing orphaned workers from accumulating CPU or
  blocking foreground CLI commands.

## 0.5.2 - 2026-07-20

### Fixed

- Changed the Linux/macOS installer to download latest native and agent assets
  through GitHub's stable `releases/latest/download` aliases, avoiding scoped
  package-tag truncation during redirect parsing.
- Changed the no-argument greeting to report Windows, Linux, or macOS according
  to the runtime operating system.

### Changed

- Migrated RunX release configuration from legacy `mirror.config.toml` to
  `mirror.yaml` while preserving package, commit, tag, push, and changelog
  behavior.

## 0.5.1 - 2026-07-19

### Fixed

- Enforced LF line endings for Bash scripts so Windows checkouts preserve
  portable installer syntax.
- Gave the detached agent-maintenance integration test an explicit
  platform-appropriate timeout so slower Windows CI runners validate the full
  workflow without changing runtime behavior.

## 0.5.0 - 2026-07-19

### Added

- Added silent detached maintenance that keeps both global `guiho-s-runx`
  skill installations current and reconciles one nearest managed `AGENTS.md`
  block without changing foreground output or exit behavior.
- Added explicit regression coverage for every upgrade recovery outcome,
  executable stable/prerelease installer verification, and the Unicode
  description-aligned command tree.

### Changed

- `runx upgrade list` now returns the complete published catalog, including
  labeled prereleases, unless pagination is explicitly requested.
- The Linux/macOS installer and generated recovery commands now use Bash with
  `set -euo pipefail`; canonical documentation no longer invokes `sh`.
- Release delivery now validates the exact twelve native binaries and two
  Markdown agent assets and uses only the matching changelog version section
  for GitHub Release notes.

### Fixed

- Fixed `runx upgrade list` and `runx upgrade check` so Citty does not execute
  the parent upgrade action after the selected leaf or append a second JSON
  document.
- Fixed automatic agent-resource updates to use idempotent atomic writes,
  preserve user-authored instructions, migrate legacy markers, and remain safe
  under concurrent invocations.

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
