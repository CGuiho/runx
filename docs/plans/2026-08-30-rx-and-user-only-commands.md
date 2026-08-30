---
name: RX and User-Only Commands Implementation Plan
purpose: Sequence the rx alias and user-only guard from manifest through launcher, installers, docs, and minor release
description: Five sealed units plus release — strictly ordered where catalog/executor/installer/docs dependencies exist, otherwise parallel-safe
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

Deliver `rx` (bare → `runx list`, `<selector> …` → `runx run <selector> …`, `-v/--version` and `-h/--help` parity, no duplicated Cobra) and the per-command `userOnly` guard (refuse without spawn, stable message, `reveal` still allowed). Both install/update/uninstall transactionally as one minor release. The Go/Cobra payload stays the source of truth; `rx` is a thin launcher.

## Traceability

- [Brainstorm](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- [Requirements](../requirements/2026-08-30-rx-and-user-only-commands.md)
- [Architecture](../architecture/2026-08-30-rx-and-user-only-commands.md) — ADRs 001–005, self-review resolved
- `guiho-convention-0001-cli.md` — stable launcher, immutable payloads, `artifacts.json` / `checksums.txt`, pointer, shared ownership
- `runx.yaml` — validation signals (UIDs: `runx-tooling-go-fmt`, `runx-tooling-go-test`, `runx-tooling-go-vet`, `runx-tooling-go-build`, `runx-tooling-tidy-check`, `runx-tooling-build-binaries`)

## Sealed Execution Statement

- All material product questions are answered in the requirements (field name `userOnly`, refusal wording, exit 2, defaults, `reveal` exempt).
- Remaining assumption (dry-run on guarded command = refuse) is approved for this plan; executor will implement it and the question ledger will record if the human later prefers the alternative.
- No secret-bearing env/key file is required. The question ledger root is `docs/questions/2026-08-30-rx-and-user-only-commands/`.
- First executable unit is **RX-1** (manifest & schema). Units RX-2 and RX-3 are parallel-safe after RX-1. RX-4 depends on RX-3. RX-5 depends on RX-1/RX-2/RX-3. RX-6 (release) depends on all code + docs units and on human approval of exact-head review/validation.

## Units

### RX-1 — Manifest and Schema: `userOnly` Field

**Goal:** Accept optional `userOnly: true` on leaf commands, propagate to `ResolvedCommand`, keep strict-YAML and backwards compatibility, update JSON schema.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** none (first unit).

**Expected files:**
- `pkg/manifest/types.go` — add `UserOnly *bool` to `Command`, `UserOnly bool` to `ResolvedCommand`
- `pkg/manifest/parser.go` or `composition.go` — leaf-only validation, propagation through composition (including child catalogs)
- `pkg/manifest/manifest_test.go` — happy path (leaf true/false/omitted), group rejection, composition propagation
- `schemas/runx.schema.json` — add `userOnly` boolean to command definition
- Any embedded schema copy under `embed/` if applicable
- XDocs touches: `pkg/manifest/manifest.xdocs.md` if descriptor needs updating

**Data/auth/cache:** catalog trusted code; no DB/cache; validation is strict parse time.

**Acceptance:**
- `go vet` passes that `userOnly: true` on a leaf validates and round-trips through `Catalog.Resolve`.
- `userOnly` on a group is a parse error.
- Omitting `userOnly` behaves as `false`.
- `runx check` passes for a fixture catalog with and without the field.

**Validation signals:**
- `runx run runx-tooling-go-vet`
- `runx run runx-tooling-go-test` (at least manifest tests)
- `runx check --format json` against a temporary catalog exercising the field

**Stop conditions:** Do not proceed to RX-2 if strict decoder or schema rejects the new field.

---

### RX-2 — Executor Guard: Refuse Guarded Commands Without Spawn

**Goal:** After resolution and before any shell construction, refuse `run` of a `userOnly` command with a stable, agent-readable message and exit 2; keep `reveal`/`list`/`describe`/`check` allowed, keep `reveal` printing exact command.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-1 (needs `ResolvedCommand.UserOnly`).

**Expected files:**
- `cmd/catalog.go` and/or `cmd/run.go` (where `run` invokes the executor) — guard check
- `cmd/reveal.go` — confirm it bypasses guard (no change, but covered by test)
- `pkg/executor/executor.go` — ensure no spawn occurs when guard fires (or guard lives above it — choose one layer and document)
- Tests: `cmd/*_test.go` — cases: UID, canonical selector, shorthand, numeric index, all with `userOnly:true` refuse; unguarded still runs; `reveal` still prints; `--dry-run` on guarded also refuses; stdout purity (message on stderr)
- XDocs: `pkg/executor/executor.xdocs.md` if needed

**Acceptance:**
- `runx run <guarded-selector>` prints exact refusal `This command is marked as user-only and must be run by the user. Agents should not run it.` plus `(selector: …, uid: …)` to stderr, stdout empty, exit 2, no marker file from executor.
- `runx reveal <guarded-selector>` prints exact stored command + `\n` to stdout, exit 0.
- Unguarded commands are unaffected; numeric indexes resolve correctly for both paths.
- Agent tests prove refusal does not call the executor's shell adapter.

**Validation signals:**
- `runx run runx-tooling-go-test`
- `runx run runx-tooling-go-vet`
- Manual exercises (described in unit, not as RunX UIDs): `go run . reveal <guarded>` vs `go run . run <guarded>`

---

### RX-3 — RX Thin Launcher

**Goal:** Build the `rx` launcher that translates `rx` ergonomics and delegates to the active RunX payload.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** none for code, but verification (RX-4) needs the binary built. Parallel-safe with RX-1/RX-2.

**Expected files:**
- `cmd/rx/main.go` — arg translation (bare→list, `<selector>…`→`run <selector>…`, `-v/--version/-h/--help/--help-tree/--help-docs` passthrough) + pointer read + delegate (reuse `pkg/installstate`, `pkg/launcher` helpers, OS-specific `delegate_*.go` pattern)
- `cmd/rx/rx.xdocs.md` — new descriptor (subject `runx-rx`)
- `cmd/cmd.xdocs.md` — parent links
- XDocs root: no new exclusion

**Acceptance:**
- `go build ./...` produces both `runx` and `rx` binaries (verified via `devops/build-binaries.go` after RX-4).
- `rx` with no args produces byte-identical output to `runx list` for the same `--cwd/--config`.
- `rx <selector>` produces same behavior as `runx run <selector>` (including userOnly refusal, confirmation, child-arg forwarding) — i.e., delegation is faithful.
- `rx -v` / `rx --version` and `rx -h` / `rx --help` are byte-identical to `runx` equivalents (modulo program name in help).
- `rx __self-test` is not required; `runx __self-test` covers payload.

**Validation signals:**
- `runx run runx-tooling-go-build`
- `runx run runx-tooling-go-vet`
- `go test -count=1 ./...` (after RX-1/RX-2)

---

### RX-4 — Build Matrix, Installers, and Lifecycle (Dual Binary)

**Goal:** Extend the 8-platform build, checksums, `artifacts.json`, staging, atomic activation, rollback, and uninstall to handle both launchers.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-3 (needs `cmd/rx`). Also benefits from RX-1 for accurate schema artifacts.

**Expected files:**
- `devops/build-binaries.go` — emit `runx-launcher-<platform>` + `rx-launcher-<platform>` + `runx-payload-<platform>` (keep `CGO_ENABLED=0`, `go vet` clean)
- `devops/verify-release-assets.go` — expect both launchers + payloads + checksums
- `devops/install.sh` — download/verify both launchers, verify `artifacts.json` digests for both, stage both, backup/restore both on rollback, atomically move both, verify `runx --version` and `rx --version`
- `devops/install.ps1` — same (PowerShell, BOM-free pointer)
- `cmd/uninstall.go`, `devops/uninstall.sh|ps1` — remove both `runx(.exe)` and `rx(.exe)`, never delete shared `~/.guiho/bin` or `~/.guiho/.temp` parents, correct `--dry-run`/`--preserve-*`/`--yes` reporting for both
- `pkg/installstate/*`, `pkg/launcher/*`, `pkg/updater/*` — helpers for `rx` launcher/payload names, pointer helpers remain single-pointer
- XDocs: `devops/devops.xdocs.md`

**Acceptance:**
- `go run devops/build-binaries.go --version 0.13.0 --commit <head> --build-date <RFC3339>` emits two launcher assets per platform plus payloads.
- `go run devops/verify-release-assets.go` passes with the expanded matrix.
- Staged installer dry-run lists both launchers as `REMOVE`/`PRESERVE` correctly; failure injection restores both launchers and the pointer.
- Post-activation `runx --version` == `rx --version` == staged version; `runx __self-test` passes.

**Validation signals:**
- `runx run runx-tooling-go-build`
- `runx run runx-tooling-build-binaries` (or direct `go run devops/build-binaries.go ...`)
- `go run devops/verify-release-assets.go` (from `runx.yaml` if mapped, else direct)

---

### RX-5 — Docs, Agent Skill, Instruction, and XDocs

**Goal:** Document both features, update the agent skill and instruction, and keep the XDocs tree healthy.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-1, RX-2, RX-3 (docs describe finalized behavior).

**Expected files:**
- `skills/guiho-s-runx/SKILL.md` — add `## RX — Short Alias` section and document `userOnly` with example, refusal snippet, and "agent must stop and report the message" guidance; keep `## CLI Evolution and Feedback` intact
- `embed/skills/guiho-s-runx.SKILL.md` — mirror the canonical skill
- `prompts/guiho-i-runx.md` — bounded instruction block: mention `rx` exists (no duplication of skill body)
- `README.md`, `DOCS.md`, `CHANGELOG.md` — add `rx` usage and `userOnly` catalog example
- `schemas/*.md` companion docs if needed
- XDocs: `skills/guiho-s-runx/guiho-s-runx.xdocs.md`, `cmd/rx/rx.xdocs.md`, other affected descriptors; run `xdocs scan` equivalents
- `TODO.md` and `docs/todo/*.md` task spec transitions to `testing`

**Acceptance:**
- `skills/guiho-s-runx/SKILL.md` contains the new `rx` and `userOnly` guidance and remains convention-compliant.
- `xdocs meta --strict` / `xdocs tree` / `xdocs doctor --warnings-as-errors` pass for touched scopes.
- Help docs (`runx --help-docs`, `rx --help-docs` delegated) render without drift.

**Validation signals:**
- `xdocs meta <scope> --documents --strict`
- `xdocs tree`
- `xdocs doctor --warnings-as-errors`

---

### RX-6 — Release (Minor) — Gates and Publication

**Goal:** Publish the combined `rx` + `userOnly` capability as the next Mirror-managed minor release with exactly the verified asset set.

**Owning repository:** `C:/GUIHO/runx`

**Dependencies:** RX-1 through RX-5 (clean `main`, exact review head, CI green).

**Expected files:**
- `mirror.yaml` — not hand-edited; `mirror version plan minor` / `apply` owns version fields
- `CHANGELOG.md` — section for the new minor
- Review/validation evidence under `docs/reviews/implementation/` and `docs/validation/`

**Acceptance (from clean integrated `main`):**
- `gofmt -l main.go cmd pkg embed devops` → no output
- `go mod tidy` leaves `go.mod`/`go.sum` unchanged
- `go test -count=1 ./...` + `go vet ./...` + `go build ./...` pass
- `go run devops/build-binaries.go` + `verify-release-assets.go` pass
- Strict XDocs passes (`meta`, `tree`, `doctor`)
- `mirror config check` and `mirror version plan minor` show the expected minor (e.g., `runx/v0.13.0` from the current baseline) with no dirty plan
- `mirror version apply minor --yes` creates the correct annotated tag; the protected `production` workflow publishes the release as non-draft, non-prerelease with exactly the expected assets (both launchers + payloads + zips + `checksums.txt` + `artifacts.json`); independent checksums/digests match; downloaded Windows AMD64 `runx --version` / `rx --version` / `rx` bare / guarded `runx run` refusal all smoke-pass; npm latest matches.

**Validation signals (pre-apply):**
- All of the above, each recorded as evidence
- Protected branch checks and mergeability reobserved at the exact review head

**Stop conditions for the whole plan:**
- Stop before merge for mismatched review/validation heads, failed checks, unresolved conflicts, or protection failure.
- Stop before `mirror version apply` if the worktree is dirty, the plan is unexpected, or asset verification is not clean.
- Authentication/connectivity failures are technical blockers — never bypass hooks or force-push.
- Production deployment, traffic, DNS, or secret rotation requires explicit human approval for the exact action — a GitHub Release is publication, not deployment.

## First Executable Unit

**RX-1 — Manifest and Schema: `userOnly` Field** is the only unit that may start first. RX-2 and RX-3 may then run in parallel. RX-4 waits on RX-3. RX-5 waits on RX-1–RX-3. RX-6 waits on all and on human validation of the exact-head review/validation.

## Execution Order

```
RX-1 ─┬─► RX-2 ─┐
      └─► RX-3 ─┴─► RX-4 ─► RX-5 ─► RX-6 (release)
```
RX-2 and RX-3 are parallel-safe (different file sets: `pkg/manifest` vs `cmd/rx`).

## Commit / Delivery Behavior

- Smallest coherent commits per unit with explicit path staging; `guiho-s-0032-git-commit` owns push.
- No production mutation without explicit human approval for the exact deployment action.
- Every unit updates its `docs/todo/<task>.md` milestone and the XDocs descriptors it touches.
