---
name: RX Short Alias and User-Only Commands
purpose: Define the outcome, scope, and acceptance for the rx alias and the user-only execution guard
description: Specifies the rx bare-list/run delegation with version/help parity, the per-command userOnly guard with agent refusal, dual-binary lifecycle, and minor-release acceptance
created: "2026-08-30T16:15:00Z"
flags:
  - draft
tags:
  - todo
  - cli
  - rx
keywords:
  - rx
  - runx
  - userOnly
  - launcher
owner: runx-todo
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RX Short Alias and User-Only Commands

## Todo Index

- Task: `27. RX short alias and user-only commands` (next index — confirm against `TODO.md`)
- Status: draft (awaiting human validation of requirements → architecture → plan)
- Index: [TODO.md](../../TODO.md)
- Brainstorm: [2026-08-30 RX brainstorm](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- Requirements: [2026-08-30 requirements](../requirements/2026-08-30-rx-and-user-only-commands.md)
- Architecture: [2026-08-30 architecture](../architecture/2026-08-30-rx-and-user-only-commands.md)
- Plan: [2026-08-30 plan](../plans/2026-08-30-rx-and-user-only-commands.md)
- ADRs: ADRs 001–005 inside the architecture document (rx launcher, single pointer, two launchers/one payload, `userOnly` field, execution-time refusal)

## Outcome

RunX ships a thin `rx` launcher (`rx` / `rx.exe`) where bare `rx` lists the catalog and `rx <selector> …` runs a catalog command, both by delegating to the active RunX payload with full flag/chld-arg fidelity and identical `-v/--version` and `-h/--help` behavior. Any catalog leaf may declare `userOnly: true`; a guarded `runx run` (and therefore `rx <selector>`) refuses without spawning and prints a stable, agent-readable message, while `runx reveal` remains the human copy-paste hatch. Both launchers are installed, upgraded, and uninstalled transactionally as one minor release with a complete `checksums.txt` / `artifacts.json` ownership manifest.

## Scope

### In scope

- `pkg/manifest/types.go` + composition + strict parser: optional `userOnly?: boolean` on leaf commands, default `false`, group rejection, propagation to `ResolvedCommand.UserOnly`
- `schemas/runx.schema.json` and any embedded schema copy for the new field
- `cmd/catalog.go` / `cmd/run.go` / `pkg/executor`: guard between resolve and spawn — no shell construction, no confirm prompt, no env/key handling when guarded; exit 2 with message `This command is marked as user-only and must be run by the user. Agents should not run it.` plus selector/uid diagnostics on stderr; `reveal` exempt; `list`/`describe`/`check` surface without execution
- `cmd/rx/main.go` thin launcher: bare→`list`, `<selector>…`→`run <selector>…`, version/help passthrough, pointer read via `pkg/installstate`, delegate via `pkg/launcher` pattern (OS-specific)
- `devops/build-binaries.go`: 8-platform `runx-launcher-*` + `rx-launcher-*` + `runx-payload-*` (CGO_ENABLED=0), version/commit/date embedding
- `devops/verify-release-assets.go`: verification aware of both launchers
- `devops/install.sh` / `devops/install.ps1`: staged dual-launcher download, SHA-256 + `artifacts.json` verification, staged payload self-test, atomic `current.json` activation, post-activation `runx --version` + `rx --version` verification, rollback of both launchers + pointer on failure, idempotent `~/.guiho/bin/` PATH handling, BOM-free pointer on PowerShell
- `cmd/uninstall.go`, `devops/uninstall.sh|ps1`, `runx uninstall`: removal of both `runx(.exe)` and `rx(.exe)` launchers plus all versioned payloads and `~/.guiho/runx/`, preserving shared `~/.guiho/bin` and `~/.guiho/.temp` parents, correct `--dry-run`/`--preserve-*`/`--yes` reporting for the combined install
- `skills/guiho-s-runx/SKILL.md` and `embed/skills/guiho-s-runx.SKILL.md`: `rx` ergonomics + `userOnly` authoring example + refusal handling + evolution policy preserved
- `prompts/guiho-i-runx.md` bounded block: mention of `rx` without skill duplication
- `README.md` / `DOCS.md` / `CHANGELOG.md`: `rx` quick-start + `userOnly` example
- XDocs descriptors: new `cmd/rx/rx.xdocs.md`, updated `cmd/cmd.xdocs.md`, `pkg/manifest/manifest.xdocs.md`, `pkg/executor/executor.xdocs.md`, `devops/devops.xdocs.md`, strict `xdocs meta` / `xdocs tree` / `xdocs doctor --warnings-as-errors`
- Mirror: `minor` bump (new public CLI + new manifest capability), `mirror config check` + `mirror version plan` + `mirror version apply --yes` only after exact-head review/validation, publication workflow produces exactly the expected asset count (both launchers + payloads + checksums + artifacts + skill + prompts + schemas), independent checksum/digest + Windows AMD64 smoke (`runx --version` / `rx --version` / `rx` bare / guarded refusal)

### Out of scope

- New database, cache, cloud resource, or paid service
- Per-group or per-namespace `userOnly` inheritance (v1 is per-leaf only)
- OS-level user identity check; `userOnly` is catalog-declared, not a principal proof
- `--allow-agent` override (deferred, addable backwards-compatibly)
- Changes to shell adapter selection, path translation, or manifest parent/child composition beyond the new field
- New Cobra commands on `rx` beyond delegation (no `rx list` subcommand tree)
- Production deployment, traffic, DNS, or secret mutation (release publication is not deployment)

## Acceptance Signals

- `rx` with no args prints byte-identical output to `runx list` for the same `--cwd/--config/--format`.
- `rx <selector> [-- child]` runs exactly as `runx run <selector> [-- child]` (lossless child forwarding, exit-code preservation, confirmation handling).
- `rx -v` / `rx --version` and `rx -h` / `rx --help` are byte-identical to `runx` equivalents; `rx --help-tree` / `rx --help-docs` delegate faithfully.
- A catalog leaf with `userOnly: true` resolved by UID, canonical selector, shorthand, or numeric index, when targeted by `runx run` or `rx <selector>`, is refused without spawning: stderr contains the stable refusal, stdout is empty, exit is 2, no executor marker file.
- `runx reveal <guarded-selector>` prints the exact stored command + `\n` to stdout, exit 0, no refusal.
- Omitting `userOnly` or setting `userOnly: false` leaves behavior unchanged; setting `userOnly` on a group is a parse error.
- Fresh install places both `~/.guiho/bin/runx(.exe)` and `~/.guiho/bin/rx(.exe)` plus payloads, verifies both `--version` and `runx __self-test`, and atomically activates `~/.guiho/runx/current.json`; failure rolls back both launchers and the pointer, leaving a usable previous install.
- Upgrade and uninstall handle both binaries transactionally with correct `--dry-run` preservation grouping and never delete shared `~/.guiho/bin` / `~/.guiho/.temp`.
- Go formatting, tidy, `go test -count=1 ./...`, `go vet ./...`, `go build ./...`, cross-build (`build-binaries.go`), `verify-release-assets.go`, strict XDocs (`meta`/`tree`/`doctor`), and `mirror config check` + `mirror version plan minor` all pass on the exact review head.
- The merged `main` is published as the next minor with exactly the expected assets and the downloaded Windows AMD64 binary smokes (`runx --version` / `rx --version` / `rx` / guarded refusal) match checksums and digests.

## Dependencies and Context

- Baseline: current protected `main` (verify with `git rev-parse HEAD` before planning commit).
- Dedicated feature branch and worktree are to be created at execution start; no branch yet (planning only).
- GitHub issues/PRs for this initiative will be created at execution start if desired; none exists yet.
- Convention 0001 compliance task 23 is in `testing` — this plan assumes that lineage is preserved, not re-litigated.

## Watch-outs

- Preserve the established `--help-tree` / `--help-docs` contracts when delegating through `rx`.
- Guard must sit above the executor's shell construction, not inside the shell adapter, to guarantee no-spawn.
- `userOnly` must not leak into `artifacts.json` ownership — it is a manifest field, not a release artifact.
- Installers must validate both launchers against `checksums.txt` and `artifacts.json` before any activation; never partially activate.
- Skill and instruction updates must stay bounded and not duplicate the full command help.

## Lifecycle Waivers

- None. Requirements, architecture, plan, implementation review, validation, protected merge, and Mirror release are all required.

## After Finishing

- Move task state to `testing` before final exact-head review/validation without changing the review head.
- Merge only after independent review, validation, CI, protection, and mergeability gates pass for the same commit.
- From clean integrated `main`, apply the clean Mirror minor transition, verify the protected publication workflow, exact release assets, independent checksums/digests, and the downloaded Windows AMD64 smoke (including `rx` and guarded refusal), then archive the task with durable evidence.

## Mirror Decision

New public CLI (`rx`) plus new manifest capability (`userOnly`) is a new public capability, so the required target is `minor`. Mirror owns the transition; no manual version/tag edits.

## Final Acceptance

TBD — to be recorded after protected `main` merge and the next minor's published release (exactly the expected assets, verified checksums/digests, and native smoke).
