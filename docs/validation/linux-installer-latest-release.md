---
name: Latest RunX Bash Installer Validation
purpose: Record verification evidence for GitHub issue 20.
description: Captures Bash syntax, installer regression tests, full package checks, release assets, public latest-alias resolution, and remaining live-install verification.
created: 2026-07-20
flags:
  - release-ready
tags:
  - validation
  - installer
keywords:
  - issue 20
  - latest download
  - fourteen assets
owner: runx-validation
---

# Latest RunX Bash Installer Validation

## Summary

The implementation is ready for patch release. Public installation against the
new patch remains the final closure gate.

## Commands Run

| Check | Result |
| --- | --- |
| `bash -n devops/install.sh` | Passed |
| `bun test devops/installers.spec.ts` | Passed: 5 tests, 39 expectations |
| `bun run typecheck` | Passed |
| `bun test` | Passed: 54 tests, 373 expectations |
| `bun run build` | Passed |
| `bun run binary` | Passed |
| `bun run binaries` | Passed: 12 native binaries |
| `bun run verify-assets` | Passed: exactly 14 assets |
| Public `releases/latest/download/runx-linux-x64-baseline` headers | Passed: redirects to the current RunX release asset and returns HTTP 200 |
| XDocs strict metadata, tree, and doctor | Passed |

## Regression Evidence

- The script contains the stable latest-download endpoint.
- The script does not inspect `url_effective`.
- Sourced Bash tests prove latest and exact encoded-tag URL shapes.
- Existing payload, Markdown, and executable checks remain green.

## Remaining Gate

After publication, execute the public Bash installer on Linux, verify the
reported version, and record the result in issue 20 before closure.
