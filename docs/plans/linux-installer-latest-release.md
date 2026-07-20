---
owner: runx-plans
tags:
  - installer
  - linux
keywords:
  - issue 20
  - latest download
  - regression test
---

# Latest RunX Bash Installer Plan

## Unit 1: Replace Fragile Latest-Tag Discovery

1. Keep exact-version normalization for explicit stable and prerelease inputs.
2. Route `latest` binary and agent assets through
   `releases/latest/download/<asset>`, matching the proven XDocs installer.
3. Verify the installed executable's semantic version and report the resolved
   version after installation.

## Unit 2: Prove The Regression

1. Assert that latest URLs use the stable GitHub download alias.
2. Assert that redirect `url_effective` parsing is absent.
3. Preserve exact encoded-tag URL coverage, Bash syntax, piping, native payload
   validation, agent-asset validation, and executable version checks.
4. Run focused tests, then the full RunX typecheck, test, build, binary, asset,
   and XDocs gates.

## Unit 3: Release And Close

1. Prepare a Mirror patch release with version-only changelog notes.
2. Publish exactly fourteen assets.
3. Execute the public Linux installer path against the release.
4. Post evidence to issue 20 and close it only after the public install passes.
