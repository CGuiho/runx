---
name: RunX Smart and rxs Launcher Handoff
purpose: Record the executed RunX Smart work, validation evidence, publication, and residual risks for human review
description: Self-review handoff for runx smart and the rxs launcher, including unit commits, safe checks, release evidence, installed smoke, and disclosed deviations.
created: "2026-10-22T05:40:00Z"
flags: []
tags:
  - review
  - handoff
  - release
keywords:
  - runx smart
  - rxs
  - launch release
  - self-review
owner: runx-reviews
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Smart and rxs Launcher Handoff

## Outcome

`runx smart` and `rxs` are implemented, validated, released as `runx/v0.17.0`,
and installed locally at `0.17.0`. GitHub issue
[RunX #66](https://github.com/CGuiho/runx/issues/66) carries the task; the
specification is
[docs/todo/2026-10-22-runx-smart-and-rxs-launcher.md](../todo/2026-10-22-runx-smart-and-rxs-launcher.md).

## Delivered contracts

- `runx smart` is a first-class Cobra command that shares the extracted
  `runCatalogList` and `runCatalogRun` bodies with `list` and `run`; bare
  `runx smart` is byte-identical to `runx list`, and
  `runx smart <selector> [--] [child args...]` is identical to
  `runx run <selector> ...`. `SetInterspersed(false)` stops RunX flag parsing at
  the first positional so child arguments are forwarded untouched.
- `runx smart` help is scoped to the command; `-v`/`--version` stays top-level
  only and is rejected with exit 2 on `smart`.
- The inspection boundary covers `smart`: bare `smart` and `smart --dry-run`
  schedule no update or maintenance worker; real `smart <selector>` runs
  schedule exactly like `run`.
- `cmd/rx` was renamed to `cmd/rxs` (never shipped, no compatibility alias).
  Diagnostics read `rxs launcher:`; bare `rxs` forwards `smart`; root entry
  points (`-v`, `--version`, `-h`, `--help`, `--help-tree`, `--help-docs`,
  `--help-tree-depth[=N]`, `--help-tree-global-flags`, single `--color`) pass
  through unchanged; every other argument list forwards as `smart <args>`.
- Release publishes eight `rxs-launcher-*` assets with `kind: launcher`,
  `installedPath: ./.guiho/bin/rxs[.exe]`, and no `rx` artifact.
- Both installers and the PowerShell whole-release test install, verify,
  back up, and roll back `runx` and `rxs` transactionally; `uninstall.sh`,
  `uninstall.ps1`, and `runx uninstall` remove both launchers while preserving
  the shared `~/.guiho/bin` directory.

## Unit commits

| Unit | Commits |
| --- | --- |
| Orchestration | `276f88a` task and briefs, `ba67e3a` xdocs descriptors, `171c10a` xdocs backfill |
| Core CLI | `b759177` shared bodies, `0fa8fed` smart command, `d11fc13` rxs retarget |
| Release plumbing | `f4edc20`, `901f46d`, `eb3ab64`, `f7fb3e8`, `a886793`, `3441da1`, `9f09329`, `bc99ff1`, `2f2972c` |
| Agent artifacts and docs | `27f6bd2`, `3b3789b`, `684a84d`, `7503610`, `ac3a727`, `34588df`, `07e1984`, `2cf2adc`, `6c75322` |
| Instruction surfaces | `f3e5205` embedded instruction sync, `95bd279` managed block sentence |

## Validation evidence

| Check | Result |
| --- | --- |
| `gofmt -l main.go cmd pkg embed devops` | clean |
| `go build ./...` | pass |
| `go vet ./...` | pass |
| `go test -count=1 ./...` | all packages pass |
| `go run devops/build-binaries.go --version 0.17.0 ...` | 46 declared artifacts plus manifest and checksums (48 files) |
| `go run devops/verify-release-assets.go` | verified protocol-v1 release v0.17.0 |
| `runx check --format json` | valid, 7 commands |
| `xdocs meta . --strict`, `xdocs doctor .` | valid, 0 errors, 0 warnings |
| `mirror config check` | ok |
| `mirror version plan minor` | 0.16.1 -> 0.17.0, tag `runx/v0.17.0` |

Behavioral evidence with the built payload:

- `smart` vs `list`: equal for empty, `--format json`, `--verbose`, and
  `--cwd/--format json`.
- `smart` vs `run`: equal for a selector, `--dry-run`, `--dry-run --format
  json`, and `--yes` on a confirmation-gated command.
- Child forwarding: `smart echoargs -- --yes --dry-run x` and
  `run echoargs -- --yes --dry-run x` produced identical output and exit code.
- Help scoping: `runx smart --help` shows the smart command,
  `runx smart --help-tree --help-tree-depth 1` shows only smart and its flags,
  `rxs --help` shows root help, `rxs --help-tree` lists `smart`.
- Local installer transaction from a bundle-backed server: checksum and
  manifest verification, both launchers installed, both `--version` outputs
  `0.17.0`, hidden self-test ok, pointer `{"protocol":1,"active":"0.17.0"}`,
  and `runx uninstall --yes` removed both launchers and the version tree while
  preserving `~/.guiho/bin`.

## Release evidence

- Tag `runx/v0.17.0` pushed; release published 2026-10-22T05:32:18Z with 48
  assets, including all eight `rxs-launcher-*` files.
- Publish workflow run 39728177667 succeeded end to end: tests, vet, matrix
  build and verification, release-notes extraction, asset publication,
  exact-asset verification, and the exact-version installer check that
  requires `rxs --version` to equal `0.17.0`.
- Local installation upgraded to `0.17.0`: `/root/.guiho/bin/runx` and
  `/root/.guiho/bin/rxs` both report `0.17.0`; `rxs` lists the repository
  catalog, `rxs -v` prints `0.17.0`, and `rxs --help-tree` lists `smart`.

## Deviations and disclosed risks

1. **Mixed commit `f3e5205`.** The pure `git mv cmd/rx cmd/rxs` moves were
   already staged in the shared index when the parent committed the embedded
   instruction sync, so that commit also contains five R100 renames. Content is
   correct and the semantic rename landed in `d11fc13`; history was not
   rewritten.
2. **Pre-existing Windows CI failure.** `go (windows-latest)` fails
   `TestInspectionLifecycleFilesystemBoundary` at
   `cmd/lifecycle_boundary_test.go:163`, both at `5390d20` before this work and
   at `171c10a`. It is unrelated to RunX Smart and was not fixed or hidden.
   The publish workflow's Windows job runs only the PowerShell installer test,
   which passed.
3. **`curl | sh` mismatch (pre-existing).** `devops/install.sh` is a bash script
   using process substitution, while the README and install prompt advertise
   POSIX `sh`. On Debian/Ubuntu `curl ... | sh` fails with
   `Syntax error: "(" unexpected`; `| bash` works. This predates the task and
   needs an explicit human decision: make the script POSIX or document `bash`.
4. **Whole-release upgrade does not refresh launchers (pre-existing).**
   `pkg/lifecycle.UpgradeWholeRelease` downloads the new payload but reinstalls
   neither launcher; `rxs` keeps working through the shared pointer and is
   verified by the PowerShell upgrade test.
5. **Cosmetic and test-scope notes.** `smart` sorts in the unranked help group
   (`cmd/help.go`), and `cmd/lifecycle_boundary_test.go` was not extended with
   `smart` rows because equivalent boundary coverage exists in
   `cmd/smart_test.go` `TestSmartInspectionBoundary`.
6. **XDocs backfill.** Ten pre-existing undocumented documents across seven
   descriptors were added to make strict XDocs pass repository-wide; metadata
   only, no content changes.

## Changes outside the repository

The published installer reconciled its managed RunX block in `/root/AGENTS.md`,
which now contains the bounded instruction including the `runx smart`/`rxs`
sentence. That is expected installer behavior.

## Human review

Verify `runx smart` and `rxs` on the installed `0.17.0`: bare listing, selector
execution with child arguments, `--dry-run` and `--yes`, root entry-point
parity, and the reminder text for confirmation-gated commands. Task status is
`Testing` until that review accepts the work.
