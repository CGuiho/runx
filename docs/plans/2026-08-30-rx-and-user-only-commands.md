---
name: RX and User-Only Commands Implementation Plan
purpose: Sequence the rx alias and informational userOnly flag from manifest through launcher, installers, docs, and minor release
description: Four sealed units plus release — strictly ordered where catalog/installer/docs dependencies exist, otherwise parallel-safe (no enforcement)
created: "2026-08-30T16:00:00Z"
flags:
  - approved
tags:
  - plan
  - cli
  - rx
  - manifest
keywords:
  - rx
  - userOnly
  - launcher
  - installers
  - minor release
owner: runx-plans
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RX and User-Only Commands Implementation Plan

## Summary

Deliver `rx` (bare → `runx list`, `<selector> …` → `runx run <selector> …`, `-v/--version` and `-h/--help` parity, no duplicated Cobra) and the informational per-command `userOnly` flag (visible in `list --format json` / `describe`, taught via the skill — CLI still executes). Both install/update/uninstall transactionally as one minor release. The Go/Cobra payload stays the source of truth; `rx` is a thin launcher.

## Traceability

- [Brainstorm](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- [Requirements](../requirements/2026-08-30-rx-and-user-only-commands.md) (informational)
- [Architecture](../architecture/2026-08-30-rx-and-user-only-commands.md) — ADRs 001–005, self-review resolved
- `guiho-convention-0001-cli.md` — stable launcher, immutable payloads, `artifacts.json` / `checksums.txt`
- `runx.yaml` — validation signals (UIDs: `runx-tooling-go-fmt`, `runx-tooling-go-test`, etc.)

## Sealed Execution Statement

- Scope is informational: `userOnly` is a flag surfaced via JSON, taught via the skill (“Hey, this is user-only, only a user should run this; if the user insists you may execute”). No CLI refusal or exit-code change.
- No secret-bearing env/key file is required. Question ledger: `docs/questions/2026-08-30-rx-and-user-only-commands/`.
- First executable unit is **RX-1** (manifest). RX-3 (launcher) is parallel-safe after RX-1. RX-4 depends on RX-3. RX-5 depends on RX-1/RX-3.

## Units

### RX-1 — Manifest: `userOnly` Field (Informational) — DONE

**Goal:** Accept optional `userOnly: true` on leaf commands, propagate to `ResolvedCommand`, keep strict-YAML and backwards compatibility.

**Files:** `pkg/manifest/types.go` (+ `UserOnly`), `parser.go` (group cannot declare), `composition.go` (propagation), `manifest_test.go` (TestUserOnlyField), `TODO.md`

**Acceptance:** `userOnly: true` validates and round-trips; group with `userOnly` is parse error; omitted = false; `go vet`/`go test` pass. Delivered via `agy` gemini-3.7-flash-high, reviewed by `pi` gpt-5.6-sol, committed `d0ba320`.

**Status:** Done — validated by you.

---

### RX-3 — RX Thin Launcher

**Goal:** Build the `rx` launcher that translates `rx` ergonomics and delegates to the active RunX payload.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-1 (needs manifest). Parallel-safe with docs.

**Expected files:**
- `cmd/rx/main.go` — arg translation (bare→list, `<selector>…`→`run <selector>…`, `-v/--version/-h/--help/--help-tree/--help-docs` passthrough) + pointer read + delegate (reuse `pkg/installstate`, `pkg/launcher` helpers, OS-specific delegate)
- `cmd/rx/rx.xdocs.md` — new descriptor (subject `runx-rx`)
- `cmd/cmd.xdocs.md` — parent links

**Acceptance:**
- `go build ./...` produces both `runx` and `rx` binaries.
- `rx` with no args produces byte-identical output to `runx list` for the same `--cwd/--config`.
- `rx <selector>` produces same behavior as `runx run <selector>` (including child-arg forwarding) — delegation is faithful.
- `rx -v` / `--version` and `rx -h` / `--help` are byte-identical to `runx` equivalents.

**Validation signals:**
- `runx run runx-tooling-go-build`
- `runx run runx-tooling-go-vet`
- `go test -count=1 ./...`

---

### RX-4 — Build Matrix, Installers, and Lifecycle (Dual Binary)

**Goal:** Extend the 8-platform build, checksums, `artifacts.json`, staging, atomic activation, rollback, and uninstall to handle both launchers.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-3 (needs `cmd/rx`).

**Expected files:**
- `devops/build-binaries.go` — emit `runx-launcher-*` + `rx-launcher-*` + `runx-payload-*` (CGO_ENABLED=0)
- `devops/verify-release-assets.go` — expect both launchers + payloads
- `devops/install.sh` / `install.ps1` — download/verify both launchers, verify `artifacts.json` digests for both, stage both, backup/restore both, atomically move both, verify `runx --version` and `rx --version`
- `cmd/uninstall.go`, `devops/uninstall.sh|ps1` — remove both `runx(.exe)` and `rx(.exe)`, never delete shared parents, correct `--dry-run`/`--preserve-*`/`--yes`
- `pkg/installstate/*`, `pkg/launcher/*` — helpers for `rx` launcher/payload names

**Acceptance:**
- `go run devops/build-binaries.go --version 0.13.0 …` emits two launcher assets per platform plus payloads.
- `verify-release-assets.go` passes with expanded matrix.
- Staged installer dry-run lists both launchers correctly; failure injection restores both launchers + pointer.
- Post-activation `runx --version` == `rx --version` == staged version; `runx __self-test` passes.

---

### RX-5 — Docs, Agent Skill, Instruction, and XDocs

**Goal:** Document `rx` and the informational `userOnly` flag, update the agent skill and instruction, keep XDocs healthy.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-1, RX-3.

**Expected files:**
- `skills/guiho-s-runx/SKILL.md` — add `## RX — Short Alias` and document informational `userOnly` with warning script (“Hey, this is user-only… if you insist you may execute”); keep `## CLI Evolution and Feedback` intact
- `embed/skills/guiho-s-runx.SKILL.md` — mirror canonical skill
- `prompts/guiho-i-runx.md` — bounded block: mention `rx` exists
- `README.md`, `DOCS.md`, `CHANGELOG.md` — add `rx` usage and `userOnly` catalog example
- XDocs: `skills/guiho-s-runx/guiho-s-runx.xdocs.md`, `cmd/rx/rx.xdocs.md`

**Acceptance:**
- Skill contains new `rx` and `userOnly` informational guidance and remains convention-compliant.
- `xdocs meta --strict` / `xdocs tree` / `xdocs doctor --warnings-as-errors` pass.

---

### RX-6 — Release (Minor) — Gates and Publication

**Goal:** Publish `rx` + informational `userOnly` as the next Mirror-managed minor release.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-1, RX-3, RX-4, RX-5, exact review head, CI green.

**Acceptance (from clean integrated `main`):**
- `gofmt`, `go vet`, `go test`, `go build`, `build-binaries`, `verify-release-assets`, strict XDocs, `mirror config check` + `mirror version plan minor` all pass
- `mirror version apply minor --yes` creates tag; protected `production` workflow publishes release as non-draft with exactly the expected assets (both launchers + payloads); checksums/digests match; Windows AMD64 `runx --version` / `rx --version` / `rx` bare all smoke-pass

## First Executable Unit (Remaining)

**RX-3 — RX Thin Launcher** is next. RX-4 waits on RX-3. RX-5 waits on RX-1/RX-3. RX-6 waits on all plus human validation of exact-head review/validation.

## Execution Order (Updated)

```
RX-1 (done) ─► RX-3 ─► RX-4 ─► RX-5 ─► RX-6 (release)
```
RX-5 docs can start in parallel with RX-3/RX-4 where file sets do not conflict, but plan keeps it after RX-4 for the skill's `rx` example to be final.

## Commit / Delivery Behavior

- Smallest coherent commits per unit with explicit path staging; `guiho-s-0032-git-commit` owns push.
- No production mutation without explicit human approval.
- Every unit updates its `docs/todo/<task>.md` milestone and XDocs descriptors it touches.
