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
