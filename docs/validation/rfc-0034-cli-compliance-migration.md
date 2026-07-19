---
name: RunX RFC 0034 CLI Compliance Migration Validation
purpose: Record direct verification evidence for every migration completion gate.
description: Command results for typecheck, tests, builds, CLI contracts, npm bootstrap, installers, release assets, imports, xdocs, and Git integrity.
created: 2026-07-18
flags:
  - validated
tags:
  - validation
  - cli
keywords:
  - RFC 0034
  - 44 tests
  - 14 assets
  - RunX
owner: runx-validation
---

# RunX RFC 0034 CLI Compliance Migration Validation

## Summary

All local RFC completion gates passed after independent audit corrections for
upgrade routing, no-argument cached notices, and names-only prompt output. No
validation blocker remains.

## Commands Run

| Command or check | Result |
| --- | --- |
| `bun run typecheck` | Passed |
| `bun test --timeout 30000` | Passed: 44 tests, 288 assertions |
| `bun run build` | Passed |
| `bun run binary` | Passed |
| `bun run binaries` | Passed: twelve native targets |
| `bun run verify-assets` | Passed: exactly fourteen unique assets |
| `bun test devops/extract-release-notes.spec.ts` | Passed: exact heading boundaries, missing-section failure, and frontmatter/older-section exclusion |
| `node --check scripts/runx-bin.mjs` | Passed |
| packed npm bootstrap local-server smoke | Passed with Node and Bun removed from PATH |
| prohibited core Node-import scan | Passed: zero matches |
| CLI banner/help/config/agent/output/exit smoke tests | Passed in CLI suite |
| live `upgrade --help`, prompt names, source banner, and compiled version smokes | Passed |
| PowerShell and POSIX installer contract tests | Passed |
| disguised-PE Markdown installer regression | Passed: installer rejected `MZ` payload before any resource write |
| `xdocs doctor --warnings-as-errors` | Passed: zero errors and zero warnings |
| `git diff --check` | Passed |

## Independent Audit Regression Evidence

- `runx upgrade` now exposes and executes `--version`, `--arch`, `--variant`,
  `--dry-run`, and `--format` without root-version interception while retaining
  `upgrade check` and `upgrade list`.
- A valid `newVersionAvailable: true` cache produces the required notice before
  the no-argument banner with the update worker disabled.
- `agent prompt list --names` prints `guiho-i-runx` as a raw text line; JSON
  mode still returns a parseable names array.
- The stale ignored `bin/runx.exe` created by the single-target build was
  removed before matrix verification. The verifier itself remains strict and
  observed only the twelve RFC binaries plus the two named agent assets.

## Exact Release Asset Evidence

The verifier observed twelve `runx-*` assets using Linux, Darwin, and Windows
names plus `guiho-s-runx.md` and `guiho-i-runx.md`. It found no duplicate,
extra, missing, or legacy platform name.

The publish workflow writes a temporary notes file containing only the exact
version section from `CHANGELOG.md`. Existing releases receive the same notes
through `gh release edit` before their asset set is refreshed.

## Skipped Checks

- No npm publication.
- No GitHub Release creation.
- No live global RunX replacement or global agent installation.
- No production deployment.

These are release mutations rather than local validation requirements.

## Readiness

Ready for the authorized Mirror patch application to `0.4.1`, one-file commits,
main push, and Mirror-managed tag/ref push. Package publication and GitHub
Release creation remain intentionally unperformed as direct local actions.

## References

- [Implementation review](../reviews/implementation/rfc-0034-cli-compliance-migration-review.md)
- [Implementation notes](../todo/rfc-0034-cli-compliance-migration-implementation.md)
