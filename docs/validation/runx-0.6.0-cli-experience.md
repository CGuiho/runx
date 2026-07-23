---
name: RunX 0.6.0 CLI Experience Validation
purpose: Record verification evidence for GitHub issues 23, 24, and 25
description: Tracks focused tests, complete package gates, native smokes, installation, release assets, public acceptance, and residual risk.
created: 2026-07-22
flags:
  - validated
  - release-pending
tags:
  - validation
  - cli
keywords:
  - RunX 0.6.0
  - issue 23
  - issue 24
  - issue 25
owner: runx-validation
---

# RunX 0.6.0 CLI Experience Validation

## Scope

Validation covers the welcome window, cached notice, child argument partition,
shell-safe forwarding, dry-run contracts, simple installer command, publish
workflow, exact assets, and public release acceptance. Issue 22 is excluded.

## Commands Run

| Check | Result |
| --- | --- |
| `bun run typecheck` | Passed. |
| `bun test` | Passed: 73 tests, 502 assertions. |
| `bun run build` | Passed. |
| `bun run binary` | Passed. |
| `bun run binaries` | Passed: twelve native targets. |
| `bun run verify-assets` | Passed: exactly fourteen assets. |
| Compiled Windows baseline welcome/version/forwarding smoke | Passed. |
| `xdocs meta . --documents --strict --format json` | Passed. |
| `xdocs doctor` | Passed: zero errors and zero warnings. |
| `mirror config check --config C:\GUIHO\runx\mirror.yaml` | Passed with the portable Mirror schema directive. |
| `mirror version plan 0.6.0 --config C:\GUIHO\runx\mirror.yaml` | Passed: only `package.json`, release commit, tag, and configured pushes are planned. |
| `git diff --check` | Passed before review. |

## Manual Checks

- Verified the publish workflow retains tag-only triggering, typecheck, tests,
  build, binary matrix, exact-version release-note extraction, exact assets,
  GitHub Release verification, and npm OIDC publishing after removal of the
  environment approval declaration.
- Verified the README contains the exact simplified curl command.
- Verified issue 22 and Go rewrite work are absent from the diff.
- Verified the native Windows baseline forwards `-v` and a spaced value to a
  child script without consuming them.

## Residual Risks

- Public installation targets the current public release until `0.6.0` is
  tagged; live acceptance cannot complete before that release exists.
- Real Linux execution remains a CI/public-release gate even though Linux
  binaries compiled and POSIX adapters passed the complete test suite.

## GitHub CI

The first post-push run, [29968115962](https://github.com/CGuiho/runx/actions/runs/29968115962),
passed Windows and every Linux source, test, build, matrix, and asset gate. Its
public-installer step initially required the unreleased welcome from the current
`0.5.4` public binary. The gate was corrected to accept that exact legacy public
startup or the new welcome while continuing to require the exact installed
version and both skill copies. Strict new-welcome behavior remains owned by the
source and compiled-native tests.

## Readiness

Implementation validated and ready for Mirror version planning. Public release
acceptance and issue closure remain pending.
