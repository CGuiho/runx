---
name: RunX Smart and rxs Launcher
purpose: Define the outcome, acceptance, and release delivery for the runx smart command and the rxs short launcher
description: Specifies runx smart list/run delegation, the rxs stable launcher and root-flag parity, the rx-to-rxs rename, dual-launcher install/upgrade/uninstall, and the minor release publication.
created: "2026-10-22T00:00:00Z"
flags: []
tags:
  - todo
  - cli
  - smart
  - rxs
keywords:
  - runx smart
  - rxs
  - rx
  - launcher
  - release
owner: runx-todo
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Smart and rxs Launcher

## Summary

CG directed RunX Smart on 2026-10-22: a first-class `runx smart` subcommand plus
a short `rxs` launcher that replaces the unshipped `rx` name because `rx`
collides with an unrelated Linux command. Both entry points expose the catalog
ergonomics already accepted in the RX planning record.

## Todo Index

- Task: `30. RunX Smart Command and rxs Launcher`
- Status: testing
- Index: [TODO.md](../../TODO.md)
- GitHub project item: [RunX #66](https://github.com/CGuiho/runx/issues/66)
- GitHub component: `runx`
- Prior record: [2026-08-30 RX short alias and user-only commands](2026-08-30-rx-and-user-only-commands.md)
- Brainstorm: [2026-08-30 RX brainstorm](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- Architecture: [2026-08-30 RX architecture](../architecture/2026-08-30-rx-and-user-only-commands.md)
- Plan: [2026-08-30 RX plan](../plans/2026-08-30-rx-and-user-only-commands.md)

## Outcome

`runx smart` lists the catalog when invoked bare and runs a catalog command when
given a selector, delegating to the same list and run implementations as the
canonical commands. `rxs` is installed at `$HOME/.guiho/bin/rxs` (`.exe` on
Windows), forwards to the active payload's `smart` command, and preserves
entry-point `-v`/`--version`, `-h`/`--help`, `--help-tree`, `--help-docs`, and
`--help-tree-depth` parity with `runx`. Both launchers install, upgrade, verify,
and uninstall transactionally, and the next minor release publishes them.

## Scope

### In scope

- `runx smart` Cobra subcommand sharing the canonical list and run execution
  paths, with full-name flags only and `-h`/`--help` scoped to the command.
- Rename `cmd/rx` to `cmd/rxs`, including diagnostics, tests, and xdocs.
- `rxs-launcher-<platform>` release assets (eight platforms), `artifacts.json`
  entries with the canonical installed path, and checksums.
- Installer, upgrade, uninstall, and `runx uninstall` coverage for `rxs`.
- Agent skill, embedded skill, instruction block, README, DOCS, catalog, and
  xdocs updates.
- Mirror minor release and publication.

### Out of scope

- Per-command `userOnly` enforcement (the field remains informational).
- Any `rx` binary, `rx` alias, or compatibility artifact for the unshipped name.
- Changes to unrelated catalog, execution, confirmation, or lifecycle contracts.

## Acceptance Signals

- `runx smart` output is byte-identical to `runx list` for the same
  `--cwd`/`--config`/`--format`.
- `runx smart <selector> [--] [child args...]` behaves exactly like
  `runx run <selector> [--] [child args...]`, including confirmation handling,
  dry-run JSON, child argument forwarding, and exit-code preservation.
- `runx smart -h/--help`, `--help-tree`, `--help-tree-depth`, and `--help-docs`
  describe the `smart` command; `-v`/`--version` stays top-level only.
- `rxs` with no arguments lists; `rxs <selector>` runs; `rxs -v`, `rxs -h`,
  `rxs --help-tree`, and `rxs --help-docs` match the `runx` root behavior.
- The release publishes eight `rxs-launcher-*` assets declared in
  `artifacts.json` and `checksums.txt`; both installers verify and roll back
  `runx` and `rxs`; both uninstallers and `runx uninstall` remove `rxs`.
- `gofmt`, `go test -count=1 ./...`, `go vet ./...`, `go build ./...`,
  `build-binaries.go`, `verify-release-assets.go`, strict XDocs, `runx check`,
  and `mirror config check` pass on the reviewed head.
- Mirror publishes `runx/v0.17.0` with the full asset set and installed `rxs`
  smoke parity.

## Dependencies and Context

- Baseline `main` at the reviewed head of the readiness policy work.
- Prior RX implementation: `cmd/rx`, CHANGELOG 0.16.0, ADRs 001-005 in the
  2026-08-30 architecture record.
- Publication is authorized by CG on 2026-10-22 for this minor release.

## Watch-outs

- `rx` is an unrelated Linux command; never ship or document `rx` as a RunX
  entry point.
- Do not duplicate list or run semantics in `smart`; share the canonical
  execution path so behavior cannot drift.
- Preserve `--help-tree`/`--help-docs` contracts and the full-name flag policy.
- `userOnly` must not leak into `artifacts.json`.
- Installers must validate both launchers before activation and roll back both.
- `verify-release-assets.go` hardcodes the artifact and checksum counts; update
  them together with the matrix.
- The Windows PowerShell whole-release upgrade test runs in publication; its
  fake release must declare the `rxs` asset.

## Before Starting

- Confirm the live GitHub Project item, Component `runx`, and `In Progress`
  status for this task.

## While Working

- Keep commits scoped per coherent unit and let each worker commit its own work
  without pushing.

## After Finishing

- Record validation evidence, handoff notes, Mirror plan/apply evidence, the
  published asset list, and installed `rxs` smoke output.

## Delivery

- Handoff and evidence:
  [docs/reviews/handoff/2026-10-22-runx-smart-and-rxs-launcher.md](../reviews/handoff/2026-10-22-runx-smart-and-rxs-launcher.md)
- Release: `runx/v0.17.0`, 48 assets, including eight `rxs-launcher-*` files;
  publish workflow run 39728177667 succeeded.
- Installed locally at `0.17.0` with both launchers verified.
- Residual risks (pre-existing Windows CI failure, `curl | sh` mismatch,
  whole-release launcher refresh) are recorded in the handoff for human
  review, not silently changed.

## Lifecycle Waiver

CG directed immediate implementation and publication on 2026-10-22 and waived
the interactive planning gate for this bounded task. The design continues the
accepted RX architecture record; no requirement, decision, or architecture
contradiction was introduced.

## References

- [TODO.md](../../TODO.md)
- [RunX #66](https://github.com/CGuiho/runx/issues/66)
