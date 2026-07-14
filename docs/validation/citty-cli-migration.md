---
name: RunX Citty CLI Migration Validation
purpose: Preserve verification evidence for the full Citty CLI migration before protected delivery and release.
description: Records TypeScript, tests, library, native binaries, npm packaging, CLI behavior, XDocs, failures, and residual distribution risk.
created: 2026-07-14
flags:
  - validated
tags:
  - validation
  - cli
  - release
keywords:
  - runx
  - citty
  - native matrix
  - npm pack
owner: runx-validation
---

# RunX Citty CLI Migration Validation

## Summary

The complete Citty migration passes the local merge and release gate. Source,
package, and native executable behavior agree on version `0.2.3` before the
planned Mirror patch to `0.2.4`.

## Scope

- Citty dependency and command definitions.
- Root, catalog, agent, upgrade, and uninstall routing.
- Help, version, errors, aliases, selector shorthand, options, output formats,
  confirmation, and dry-run safety.
- TypeScript library, local native executable, twelve release assets, npm
  tarball, packed launcher, bundled skill, and XDocs metadata.

## Commands Run

| Check | Result | Evidence |
| --- | --- | --- |
| `bun run typecheck` | passed | TypeScript emitted no diagnostics. |
| `bun test` | passed | 13 tests, 0 failures, 125 assertions. |
| `bun run build` | passed | TypeScript library output completed. |
| `bun run binary` | passed | Local Windows executable compiled with Citty bundled. |
| `bin/runx.exe -v` | passed | Printed `0.2.3` without manifest discovery. |
| `bin/runx.exe --help` | passed | Printed the Citty root command tree and options. |
| `bun run binaries` | passed | All 12 configured native assets compiled and were non-empty. |
| `npm pack --dry-run --json` | passed | Reported 78 package files. |
| `npm pack --json` | passed | Produced a 39,749-byte validation tarball with 78 files. |
| packed `scripts/runx-bin.ts -v` | passed | Printed `0.2.3` outside a manifest. |
| packed `scripts/runx-bin.ts --help` | passed | Printed Citty usage outside a manifest. |
| `xdocs tree` | passed | Included the complete RunX documentation tree. |
| `xdocs scan` | passed | Found 16 valid descriptors before adding this validation/review scope. |
| `xdocs doctor .` | passed | 0 errors and 0 warnings. |

## Package Checks

- Packed `package.json` declares `citty: ^0.2.2` under runtime dependencies.
- Packed `package.json` has no `postinstall` script.
- The launcher selected the compiled `library/guiho-runx-bin.js` path and
  returned successful version/help output.

## Manual Behavior Checks

- Source and native `-v`/help worked outside a configured project.
- Citty usage displayed command-specific arguments and only applicable flags.
- Unknown options and missing selectors returned usage errors rather than a
  missing-manifest error.
- A valid manifest selector named `constructor` executed through the root
  shorthand, proving the compatibility map is prototype-safe.

## Failures Or Blockers

No implementation blocker remains. Two environment-only validation attempts
failed before the successful package check: the sandbox denied `C:\tmp`, and
npm's default user cache was not writable. Re-running with repository-local
ignored `.temp` paths and npm cache completed successfully without changing the
artifact.

## Skipped Checks

No plan-required local check was skipped. Public CI, protected merge, Mirror
patch application, protected tag push, manual environment approval, GitHub
Release assets, and npm provenance are pending Unit 5.

## Readiness

Ready for protected implementation delivery and the Mirror-managed `0.2.4`
patch release.

## References

- [Citty CLI Migration Decision](../decisions/citty-cli-migration.md)
- [Citty CLI Migration Plan](../plans/citty-cli-migration.md)
- [Citty CLI Migration Implementation Review](../reviews/implementation/citty-cli-migration-review.md)
