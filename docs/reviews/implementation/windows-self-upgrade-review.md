---
name: RunX Windows Self-Upgrade Implementation Review
purpose: Review the implemented Windows self-upgrade fix against its accepted decision and plan.
description: Confirms synchronous replacement, target verification, rollback, cleanup, tests, CI coverage, documentation, and release readiness.
created: 2026-07-14
flags:
  - accepted
tags:
  - reviews
  - implementation
  - windows
keywords:
  - runx upgrade
  - executable replacement
  - rollback
  - windows ci
owner: runx-implementation-reviews
---

# RunX Windows Self-Upgrade Implementation Review

## Verdict

Accepted. No blocker, high, medium, or low findings remain.

## Findings

No actionable findings. The implementation matches the approved synchronous
replacement and recovery design without changing manifest execution, Citty
routing, uninstall behavior, or non-Windows upgrade behavior.

## Acceptance Criteria

- `upgradeSelf()` renames the running Windows executable, installs the
  downloaded target at the canonical path, and awaits `<path> --version`.
- Successful Windows upgrade results report `scheduled: false`.
- Verification failures remove the invalid replacement and restore `.old`.
- `.new` is absent after success and rollback.
- Cleanup retries deletion of the locked `.old` image after the old process
  releases it.
- A Windows-only test holds a real copied executable open during replacement,
  verifies the new binary at the canonical path before completion, and observes
  backup cleanup.
- A second Windows test proves rollback after invalid-binary verification.
- Existing Citty tests continue to prove `-h` and `--help` outside a manifest.
- The GitHub workflow now runs typecheck, tests, build, and native compilation
  on `windows-latest` in addition to the existing Linux release-matrix gate.

## Verification Evidence

- `bun run typecheck`: passed.
- `bun test`: 15 passed, 0 failed, 134 assertions.
- `bun run build`: passed.
- `bun run binary`: passed.
- `bun run binaries`: twelve native assets built and verified.
- Windows x64 native `-h` outside a manifest: passed.
- Windows x64 native `--version`: reported `0.2.5` before the release bump.
- CI workflow YAML parse: passed.
- `xdocs doctor . --warnings-as-errors --format json`: 0 errors, 0 warnings.
- `git diff --check`: passed.

## Documentation and TODO

`DOCS.md`, source/root/workflow XDocs descriptors, the accepted decision, plan,
plan review, this implementation review, and validation evidence are aligned.
No TODO entry is needed for this completed single-session fix.

## Residual Risk

The new `windows-latest` job cannot execute until the committed workflow reaches
GitHub. Local Windows coverage exercises the same live-file replacement and
cleanup behavior. Public upgrade from 0.2.0 to the new patch can only be verified
after the patch release assets exist.

## References

- [Decision](../../decisions/windows-self-upgrade.md)
- [Plan](../../plans/windows-self-upgrade.md)
- [Plan Review](../plans/windows-self-upgrade-review.md)
