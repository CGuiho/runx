---
name: RX Short Alias and User-Only Commands
purpose: Define the outcome, scope, and acceptance for the rx alias and the informational userOnly flag
description: Specifies the rx bare-list/run delegation with version/help parity, the per-command informational userOnly flag taught via the skill, dual-binary lifecycle, and minor-release acceptance
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

- Task: `26. RX short alias and user-only commands`
- Status: in progress
- Index: [TODO.md](../../TODO.md)
- Brainstorm: [2026-08-30 RX brainstorm](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- Requirements: [2026-08-30 requirements](../requirements/2026-08-30-rx-and-user-only-commands.md) (informational)
- Architecture: [2026-08-30 architecture](../architecture/2026-08-30-rx-and-user-only-commands.md) (informational)
- Plan: [2026-08-30 plan](../plans/2026-08-30-rx-and-user-only-commands.md)
- ADRs: ADRs 001–005 inside the architecture (rx launcher, single pointer, two launchers/one payload, `userOnly` informational, no CLI guard)

## Outcome

RunX ships a thin `rx` launcher (`rx` / `rx.exe`) where bare `rx` lists the catalog and `rx <selector> …` runs a catalog command, both by delegating to the active RunX payload with full flag/chld-arg fidelity and identical `-v/--version` and `-h/--help` behavior. Any catalog leaf may declare `userOnly: true` (informational — CLI still executes, skill teaches agents to warn “Hey, this is user-only, only a user should run this; if you insist you may execute”). Both launchers are installed, upgraded, and uninstalled transactionally as one minor release with a complete `checksums.txt` / `artifacts.json` manifest.

## Scope

### In scope

- `pkg/manifest/types.go` + composition + strict parser: optional `userOnly?: boolean` on leaf commands, default `false`, group rejection, propagation to `ResolvedCommand.UserOnly` — informational, no executor guard (DONE via Gemini, committed d0ba320)
- `cmd/rx/main.go` thin launcher: bare→`list`, `<selector>…`→`run <selector>…`, version/help passthrough, pointer read via `pkg/installstate`, delegate via `pkg/launcher` pattern (OS-specific) — NEXT
- `devops/build-binaries.go`: 8-platform `runx-launcher-*` + `rx-launcher-*` + `runx-payload-*` (CGO_ENABLED=0), version/commit/date embedding
- `devops/verify-release-assets.go`: verification aware of both launchers
- `devops/install.sh` / `devops/install.ps1`: staged dual-launcher download, SHA-256 + `artifacts.json` verification, staged payload self-test, atomic `current.json` activation, post-activation `runx --version` + `rx --version` verification, rollback of both launchers + pointer on failure, idempotent `~/.guiho/bin/` PATH handling
- `cmd/uninstall.go`, `devops/uninstall.sh|ps1`, `runx uninstall`: removal of both `runx(.exe)` and `rx(.exe)` launchers plus all versioned payloads and `~/.guiho/runx/`, preserving shared `~/.guiho/bin` and `~/.guiho/.temp` parents, correct `--dry-run`/`--preserve-*`/`--yes`
- `skills/guiho-s-runx/SKILL.md` and `embed/skills/guiho-s-runx.SKILL.md`: `rx` ergonomics + informational `userOnly` authoring example + warning script + “if user insists you may execute” policy
- `prompts/guiho-i-runx.md` bounded block: mention of `rx` without skill duplication
- `README.md` / `DOCS.md` / `CHANGELOG.md`: `rx` quick-start + `userOnly` catalog example
- XDocs descriptors: new `cmd/rx/rx.xdocs.md`, updated `cmd/cmd.xdocs.md`, `pkg/manifest/manifest.xdocs.md`, `devops/devops.xdocs.md`, strict `xdocs meta` / `xdocs tree` / `xdocs doctor`
- Mirror: `minor` bump (new public CLI + new manifest flag), `mirror config check` + `mirror version plan` + `mirror version apply --yes` only after exact-head review/validation, publication workflow produces exactly the expected asset count (both launchers + payloads + checksums + artifacts + skill + prompts + schemas), independent checksum/digest + Windows AMD64 smoke (`runx --version` / `rx --version` / `rx` bare)

### Out of scope

- New database, cache, cloud resource, or paid service
- Per-group or per-namespace `userOnly` inheritance (v1 is per-leaf only, informational)
- OS-level user identity check; `userOnly` is catalog-declared, informational (no principal proof)
- CLI enforcement / refusal for `userOnly` (intentionally omitted — skill teaches, CLI still executes)
- Changes to shell adapter selection, path translation, or manifest parent/child composition beyond the new flag
- Production deployment, traffic, DNS, or secret mutation (release publication is not deployment)

## Acceptance Signals

- `rx` with no args prints byte-identical output to `runx list` for the same `--cwd/--config/--format`.
- `rx <selector> [-- child]` runs exactly as `runx run <selector> [-- child]` (lossless child forwarding, exit-code preservation, confirmation handling).
- `rx -v` / `rx --version` and `rx -h` / `rx --help` are byte-identical to `runx` equivalents; `rx --help-tree` / `rx --help-docs` delegate faithfully.
- A catalog leaf with `userOnly: true` resolved by UID, canonical selector, shorthand, or numeric index, is visible via `runx list --format json` (`"userOnly": true`); `runx describe` / `check` surface it without execution; setting `userOnly` on a group is a parse error; omitting is `false`.
- No CLI refusal: `runx run <userOnly>` and `rx <userOnly>` still execute the stored command (verified — no guard in `cmd/` or `pkg/executor`).
- Fresh install places both `~/.guiho/bin/runx(.exe)` and `~/.guiho/bin/rx(.exe)` plus payloads, verifies both `--version` and `runx __self-test`, and atomically activates `~/.guiho/runx/current.json`; failure rolls back both launchers and the pointer.
- Upgrade and uninstall handle both binaries transactionally with correct `--dry-run` preservation grouping and never delete shared `~/.guiho/bin` / `~/.guiho/.temp`.
- Go formatting, tidy, `go test -count=1 ./...`, `go vet ./...`, `go build ./...`, cross-build (`build-binaries.go`), `verify-release-assets.go`, strict XDocs (`meta`/`tree`/`doctor`), and `mirror config check` + `mirror version plan minor` all pass on the exact review head.
- The merged `main` is published as the next minor with exactly the expected assets and the downloaded Windows AMD64 binary smokes (`runx --version` / `rx --version` / `rx`) match checksums and digests.

## Dependencies and Context

- Baseline: 805a1be (planning commit)
- Branch: feat/rx-and-user-only-commands (pushed d0ba320)
- RX-1 done via agy gemini-3.7-flash-high, reviewed by pi gpt-5.6-sol
- Convention 0001 compliance task 23 is in `testing` — this plan preserves that lineage

## Watch-outs

- Preserve the established `--help-tree` / `--help-docs` contracts when delegating through `rx`.
- `userOnly` must not leak into `artifacts.json` ownership — it is a manifest field, not a release artifact.
- Installers must validate both launchers against `checksums.txt` and `artifacts.json` before any activation; never partially activate.
- Skill and instruction updates must stay bounded and not duplicate the full command help. Skill script: “Hey, this is user-only, only a user should run this. If you insist after the warning, you may execute.”

## Lifecycle Waivers

- None. Requirements, architecture, plan, implementation review, validation, protected merge, and Mirror release are all required.

## After Finishing

- Move task state to `testing` before final exact-head review/validation without changing the review head.
- Merge only after independent review, validation, CI, protection, and mergeability gates pass for the same commit.
- From clean integrated `main`, apply the clean Mirror minor transition, verify the protected publication workflow, exact release assets, independent checksums/digests, and the downloaded Windows AMD64 smoke (including `rx`), then archive the task with durable evidence.

## Mirror Decision

New public CLI (`rx`) plus informational manifest flag (`userOnly`) is a new public capability, so the required target is `minor`. Mirror owns the transition; no manual version/tag edits.

## Final Acceptance

TBD — to be recorded after protected `main` merge and the next minor's published release (exactly the expected assets, verified checksums/digests, and native smoke).
