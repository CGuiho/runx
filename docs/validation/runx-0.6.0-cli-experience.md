---
name: RunX 0.6.0 CLI Experience Validation
purpose: Record verification evidence for GitHub issues 23, 24, and 25
description: Tracks focused tests, complete package gates, native smokes, installation, release assets, public acceptance, and residual risk.
created: 2026-07-22
flags:
  - validated
  - public-accepted
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
| `mirror version plan patch` from `0.6.0` | Passed: only `package.json`, release commit, tag, and configured pushes were planned for `0.6.1`. |
| `git diff --check` | Passed before review. |

## Manual Checks

- Verified the publish workflow retains tag-only triggering, typecheck, tests,
  build, binary matrix, exact-version release-note extraction, exact assets,
  GitHub Release verification, and npm OIDC publishing with the `production`
  environment identity required by the npm trusted-publisher record.
- Verified the README contains the exact simplified curl command.
- Verified issue 22 and Go rewrite work are absent from the diff.
- Verified the native Windows baseline forwards `-v` and a spaced value to a
  child script without consuming them.

## Public Release Acceptance

- Mirror created and pushed the immutable `@guiho/runx@0.6.1` patch tag after
  `0.6.0` exposed the npm environment-identity mismatch.
- [Publish run 30030768351](https://github.com/CGuiho/runx/actions/runs/30030768351)
  passed typecheck, all tests, build, twelve native binaries, exact fourteen
  assets, GitHub Release publication, npm OIDC publication, and provenance.
- The [0.6.1 GitHub Release](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx%400.6.1)
  contains exactly twelve native binaries plus `guiho-i-runx.md` and
  `guiho-s-runx.md`. Its body contains only the `0.6.1` changelog section.
- The public npm registry reports `@guiho/runx@0.6.1` as latest with npm
  provenance and the expected tarball integrity metadata.
- An isolated public PowerShell installation resolved `0.6.1`, verified the
  bordered Windows welcome, installed both skill copies, and forwarded `-v`,
  `space value`, and `--help` unchanged to a catalog child command.
- Attempt 2 of [CI run 30030765505](https://github.com/CGuiho/runx/actions/runs/30030765505)
  passed the real public Linux installer after `0.6.1` became the latest
  release. Attempt 1 raced publication and correctly observed the then-latest
  `0.6.0` release.

## GitHub CI

The first post-push run, [29968115962](https://github.com/CGuiho/runx/actions/runs/29968115962),
passed Windows and every Linux source, test, build, matrix, and asset gate. Its
public-installer step initially required the unreleased welcome from the current
`0.5.4` public binary. The gate was corrected to accept that exact legacy public
startup or the new welcome while continuing to require the exact installed
version and both skill copies. Strict new-welcome behavior remains owned by the
source and compiled-native tests.

The corrected post-push run,
[30029991055](https://github.com/CGuiho/runx/actions/runs/30029991055), passed
both Linux and Windows jobs, including the public Linux installer.

## Readiness

Implementation and public distribution are accepted in `0.6.1`. GitHub issues
23, 24, and 25 satisfy their acceptance signals and are ready for closure.
