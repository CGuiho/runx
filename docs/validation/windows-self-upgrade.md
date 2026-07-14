---
name: RunX Windows Self-Upgrade Validation
purpose: Record verification evidence for the Windows self-upgrade fix and issue closure.
description: Captures Windows replacement and rollback tests, complete Bun build checks, native CLI checks, CI syntax, XDocs health, and residual release risk.
created: 2026-07-14
flags:
  - validated
tags:
  - validation
  - windows
  - cli
keywords:
  - runx upgrade
  - windows executable
  - rollback
  - issue 9
  - issue 1
owner: runx-validation
---

# RunX Windows Self-Upgrade Validation

## Summary

The implementation is ready for commit, issue closure, and the requested
Mirror-managed patch. All local checks passed on Windows.

## Scope

- Synchronous Windows replacement of a running native executable.
- Target-version verification before successful completion.
- Rollback after invalid replacement verification.
- Deferred cleanup of the locked old executable only.
- Existing manifest-free `-h` behavior from issue #1.
- Windows CI coverage and structured documentation health.

## Commands and Results

| Check | Result | Evidence |
| --- | --- | --- |
| `bun run typecheck` | passed | TypeScript completed without diagnostics. |
| `bun test` | passed | 15 tests, 0 failures, 134 assertions. |
| `bun run build` | passed | Library output compiled. |
| `bun run binary` | passed | Local Windows executable compiled. |
| `bun run binaries` | passed | All twelve release assets built and verified. |
| Native `runx-windows-x64.exe -h` outside a manifest | passed | Citty usage rendered with exit code 0. |
| Native `runx-windows-x64.exe --version` | passed | Reported `0.2.5` before the patch bump. |
| `Bun.YAML.parse(.github/workflows/ci.yml)` | passed | Workflow YAML parsed successfully. |
| `xdocs meta . --documents --strict --owner runx --format json` | passed | Root metadata and companions valid. |
| `xdocs doctor . --warnings-as-errors --format json` | passed | 0 errors and 0 warnings. |
| `git diff --check` | passed | No whitespace errors. |

## Behavioral Evidence

`source/self-management.spec.ts` runs only on Windows for the platform-specific
cases. Its success case launches a copied PowerShell executable from the
configured self path, keeps that executable running, upgrades the same path to
the current Bun executable, and requires the replacement to report the expected
version before `upgradeSelf()` resolves. It also observes removal of `.new` and
eventual removal of `.old` after the running process exits.

The rollback case downloads invalid executable bytes, requires verification to
fail, and confirms that the original bytes are restored with neither `.new` nor
`.old` left behind.

The existing `source/cli.spec.ts` root-help test proves both `-h` and `--help`
work without a manifest, satisfying issue #1.

## Skipped or Blocked Checks

- The new GitHub-hosted `windows-latest` job is not locally runnable; it will
  execute after the commit reaches GitHub.
- An actual 0.2.0-to-new-patch public upgrade is blocked until the new release
  assets are published.

## Readiness

Ready for commit, closing GitHub issues #9 and #1, changelog preparation, and
the requested Mirror patch plan/apply sequence.

## References

- [Implementation Review](../reviews/implementation/windows-self-upgrade-review.md)
- [Plan](../plans/windows-self-upgrade.md)
- [Decision](../decisions/windows-self-upgrade.md)
