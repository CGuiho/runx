---
name: RX and User-Only Commands Architecture
purpose: Define the technical design for the rx short alias and the user-only command guard
description: Locks repository boundaries, launcher delegation, installer transactions, manifest extension, executor guard, and docs/versioning for both features
created: "2026-08-30T15:45:00Z"
flags:
  - approved
tags:
  - architecture
  - cli
  - rx
  - manifest
  - launcher
keywords:
  - rx
  - runx-launcher
  - installstate
  - updater
  - userOnly
owner: runx-architecture
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RX and User-Only Commands Architecture

## Summary

This architecture delivers two features as one minor release:

- **RX:** a thin `rx` launcher (`rx` / `rx.exe`) that translates ergonomics (`rx` → `runx list`, `rx <args>` → `runx run <args>`, version/help parity) and delegates to the active RunX payload. It reuses the existing pointer, payload, and verification infrastructure so no Cobra logic is duplicated.
- **User-only guard:** a new strict-YAML per-command field `userOnly` (bool, optional, default `false`) enforced at execution time. Guarded runs refuse without spawning and emit a stable, agent-readable message. `reveal` remains allowed as the human copy-paste hatch.

Both reuse the Go/Cobra stack, convention 0001 lifecycle, and Mirror/XDocs obligations.

## Source Material

- [Requirements — RX and User-Only Commands](../requirements/2026-08-30-rx-and-user-only-commands.md)
- [Brainstorm — 2026-08-30](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- `guiho-convention-0001-cli.md` — stable launcher, versioned payloads, `~/.guiho/bin`, `~/.guiho/runx`, atomic `current.json`, shared `.temp`, checksums, `artifacts.json`
- Current: `cmd/runx-launcher/main.go`, `pkg/launcher/launcher.go`, `pkg/installstate/*`, `pkg/updater/*`, `cmd/root.go`, `pkg/manifest/*`, `pkg/executor/*`, `devops/install.sh|ps1`, `devops/build-binaries.go`

## System Map

### Repositories and Components

- Single repository `C:/GUIHO/runx` — no new repository. All changes are in-place.
- Components touched:
  - `cmd/runx-launcher` (reference for `cmd/rx`), `cmd/rx` (new), `cmd/*` (help/version delegation is not duplicated — rx does not own its own Cobra tree).
  - `pkg/installstate` (pointer remains one file; helpers for rx launcher/payload names).
  - `pkg/launcher` (payload resolution for rx).
  - `pkg/manifest` (type + strict parser + composition + resolve).
  - `pkg/executor` (guard check before spawn).
  - `devops/` installers, `devops/build-binaries.go`, `devops/verify-release-assets.go`.
  - `schemas/runx.schema.json` (manifest schema) and embedded schema copies.
  - `skills/guiho-s-runx/SKILL.md`, `embed/skills/guiho-s-runx.SKILL.md`, `prompts/guiho-i-runx.md` (agent guidance).
  - XDocs descriptors: `cmd/rx/*.xdocs.md`, `pkg/manifest/*.xdocs.md`, `pkg/executor/*.xdocs.md`, etc.

### External Dependencies

- GitHub Releases API for channel/version resolution (unchanged).
- No new cloud provider, database, Valkey, or paid service.

## Decisions

### ADR-001: RX as a Thin Launcher That Delegates to the RunX Payload

- **Decision:** `rx` is a minimal launcher binary that (a) inspects its own invocation args, (b) translates bare → `list` and `<selector> …` → `run <selector> …`, and (c) execs the active RunX payload resolved from `~/.guiho/runx/current.json`. Version/help flags are passed through unchanged. `rx` does not embed its own Cobra tree or manifest parser.
- **Why:** Satisfies user's "just run RunX under the hood" instruction, eliminates drift between `runx` and `rx` help/version trees, keeps verification cheap (one payload self-test covers both), and matches convention 0001 launcher discipline.
- **Alternatives:** Full second Cobra binary duplicating `cmd/root.go` (rejected — drift, double tests); shell alias/wrapper script (rejected — not a native launcher, no checksum ownership).
- **Consequences:** `rx` binary is tiny; any `runx run` behavior change (including the user-only guard) is automatically inherited by `rx`. Installers must manage two launchers but only one payload family needs verification of business logic.
- **Reversal:** Could promote `rx` to its own payload later without changing the installed `~/.guiho/bin/rx` path.

### ADR-002: Single Pointer Governs Both Launchers

- **Decision:** Keep one `~/.guiho/runx/current.json` (`protocol:1, active, previous`) as the source of truth for both `runx` and `rx`. Both launchers read the same pointer. `previous` is still one rollback payload.
- **Why:** One atomic activation for both binaries; rollback restores both to the same version. Simpler than dual pointers with cross-version skew.
- **Alternatives:** Separate `rx.json` pointer (rejected — allows version skew and doubles activation failure modes).

### ADR-003: Two Launchers, One Payload Family

- **Decision:** Release matrix produces two launcher assets (`runx-launcher-<platform>` and `rx-launcher-<platform>`) but one coherent payload family (`runx-payload-<platform>`) that both launchers execute. `rx` launcher execs the `runx` payload path. Installers place launchers at `~/.guiho/bin/runx(.exe)` and `~/.guiho/bin/rx(.exe)`; payloads at `~/.guiho/runx/versions/<version>/runx-payload(.exe)`. `rx` has no separate versioned payload file.
- **Why:** Minimizes native artifact count, keeps build matrix small (8 platforms × 2 launchers + 8 payloads = 16 binaries plus zips/schemas), and keeps `artifacts.json` ownership clear.
- **Alternatives:** Duplicate `rx-payload` identical to `runx-payload` (rejected — doubles build/publish without value). If needed later, a separate payload can be introduced as a new artifact.

### ADR-004: Per-Command `userOnly` Boolean

- **Decision:** Add optional `userOnly?: boolean` (strict YAML, default `false` when omitted, explicit `false` preserves behavior) to `pkg/manifest/types.go#Command`. Excluded from group nodes (`commands` with children) by validation. Propagated into `ResolvedCommand.UserOnly`.
- **Why:** Readable, additive, backwards compatible (omitting stays runnable). `allowAgent` is the inverse and more confusing for authors ("allowAgent: false" double-negative in review).
- **Alternatives:** `allowAgent: boolean` (rejected), `executionPolicy: "agent" | "user"` (rejected — heavier for a v1 bool).
- **Validation:** Strict decoder rejects unknown values; type must be boolean; non-boolean or non-leaf usage is an error.
- **Schema:** `schemas/runx.schema.json` adds `userOnly: { type: "boolean" }` to the command definition, not to the group definition.

### ADR-005: Execution-Time Refusal with Stable Message, No Spawn

- **Decision:** Guard is enforced in `cmd/catalog.go` / `pkg/executor` boundary after selector resolution and before any shell construction or `confirm` prompting. On `userOnly == true`, write to stderr:

  ```
  This command is marked as user-only and must be run by the user. Agents should not run it.
  ```

  Optionally include `uid` / `selector` on the next line for diagnostics. Exit `2` (validation/usage error family). No `confirm` interaction, no shell spawn, no env/key handling. `runx reveal` and `runx check` are explicitly exempt — they never spawn.

- **Why:** Deterministic, no-spawn, matches user's requested phrasing while being grammatically stable and machine-testable. Exit 2 reuses the existing "unknown flag / validation" bucket the root already maps.
- **Alternatives:** Exit 3 custom code (rejected — not a known bucket), silent success (rejected — agent would think it ran).

## Detailed Design

### RX Launcher

- **Location:** `cmd/rx/main.go` (new package `rx`). Mirrors `cmd/runx-launcher/main.go` structure.
- **Arg translation (pseudocode):**
  ```
  args := os.Args[1:]
  if len(args) == 0 {
      args = []string{"list"}
  } else if args[0] == "-v" || args[0] == "--version" || args[0] == "-h" || args[0] == "--help" || args[0] == "--help-tree" || args[0] == "--help-docs" || args[0] == "--help-tree-depth" || args[0] == "--help-tree-global-flags" || args[0] == "--color" {
      // version/help passthrough — execute runx with identical args
  } else {
      // first positional is selector + child args → prefix with "run"
      args = append([]string{"run"}, args...)
  }
  resolve pointer via installstate.ReadPointer()
  payload := launcher.PayloadPath(pointer) // reuse existing helper
  delegate(payload, args...) // same syscall strategy as runx-launcher (delegate_unix.go / delegate_windows.go)
  ```
- **Edge:** flags like `--cwd` / `--config` / `--format` are valid for `list` and `run` — they pass through naturally because translation is just prefixing `list`/`run`.
- **Verification:** `rx --version` and `rx --help` are verified exactly like `runx` post-activation. `rx __self-test` is not required — `runx __self-test` already covers the payload.

### Manifest and Resolve

- **Types:** `Command` gains `UserOnly *bool `yaml:"userOnly,omitempty" json:"userOnly,omitempty"``; `ResolvedCommand` gains `UserOnly bool`.
- **Parser:** `pkg/manifest/parser.go` strict decoder automatically accepts the new field once present in the struct; add validation:
  - `userOnly` must be boolean when present (YAML strictness already enforces).
  - `userOnly` is allowed only on leaf commands (not on groups where `Commands != nil`).
  - Groups ignore `userOnly`; setting it on a group is an error.
- **Composition:** `composition.go` must propagate `userOnly` into the resolved catalog (including child catalogs if they declare it). No inheritance — each leaf carries its own flag.
- **Resolve:** no change to selector precedence; guard is orthogonal.
- **Schemas:** update `schemas/runx.schema.json`, `embed/schemas/*` copies if any, and `schemas/runx.global.schema.json` is not affected.

### Execution Guard

- **Place:** between `catalog.Resolve` success and `executor.ExecuteCommand` call in `cmd/catalog.go` (or `cmd/run.go`). Pseudocode:
  ```
  cmd := catalog.Resolve(selector)
  if cmd.UserOnly {
      fmt.Fprintf(os.Stderr, "This command is marked as user-only and must be run by the user. Agents should not run it. (selector: %s, uid: %s)\n", cmd.Selector, cmd.UID)
      return exitCode 2
  }
  ```
- **Scope:** applies to `runx run` and by delegation to `rx <selector>` . Does not apply to `list`, `describe`, `check`, `reveal`, `init`, `agent`, `upgrade`, `uninstall`.
- **`reveal`:** explicitly allowed; it prints `cmd.Command` verbatim and exits 0 without guard.
- **`--dry-run`:** dry-run is still a `run` intention — guard applies (refuse). Architecture locks this as refusal; discuss in validation.
- **Tests:** add `executor` / `manifest` / `cmd` tests for guarded vs unguarded selectors (UID, canonical, shorthand, numeric), refusal message, exit code, and no-spawn guarantee.

### Installers and Updaters

- **Build:** `devops/build-binaries.go` now builds 16 native binaries (8 runx-launcher + 8 rx-launcher + 8 runx-payload) or at minimum 16 assets: currently 8 platforms. Logic: build `cmd/runx-launcher` → `runx-launcher-<platform>` and `cmd/rx` → `rx-launcher-<platform>` and `cmd` payload → `runx-payload-<platform>`. Embed `version/commit/date` consistently.
- **Staging:** download and verify both launchers + payload + `checksums.txt` + `artifacts.json` + skill zip + instruction + prompts + schemas using SHA-256 and `artifacts.json` digests.
- **Activation:** place launchers atomically in `~/.guiho/bin/`, payload in `versions/<version>/`, update `current.json` once (with `previous`), verify both launchers: `runx --version` and `rx --version` both equal resolved SemVer, plus `runx __self-test`.
- **Rollback:** on any failure, restore previous `current.json`, both launchers (from `.previous` copies), and remove the new version dir. Staging directory naming stays `runx-install-<id>` but now owns two launchers.
- **Updater:** `pkg/updater` and `pkg/lifecycle` paths remain; incremental `runx upgrade` should reuse installer semantics (transactional, both launchers).
- **Uninstall:** `cmd/uninstall.go`, `devops/uninstall.sh|ps1`, and `runx uninstall` remove both `~/.guiho/bin/runx(.exe)` and `~/.guiho/bin/rx(.exe)`, all version dirs, `~/.guiho/runx/`, skill copies. Shared `~/.guiho/bin` and `~/.guiho/.temp` parents are never deleted. Flags `--preserve-config`, `--preserve-data`, `--dry-run`, `--yes` apply to the combined install.

### Docs and Release

- **Agent skill:** `skills/guiho-s-runx/SKILL.md` (and embedded copy) adds `rx` quick-start section and documents `userOnly` with an example and the refusal pattern. Agents are instructed to stop and report the URL when they see the refusal.
- **Instruction:** `prompts/guiho-i-runx.md` managed block updated to mention `rx` existence (bounded block, no duplication).
- **XDocs:** new descriptor `cmd/rx/rx.xdocs.md` (subject `runx-rx`), updates to `cmd/cmd.xdocs.md`, `pkg/manifest/manifest.xdocs.md`, `pkg/executor/executor.xdocs.md`, `devops/devops.xdocs.md` as needed. `xdocs.yaml` excludes remain unchanged.
- **Mirror:** minor bump (`0.12.x → 0.13.0` or next minor from current baseline) — new public capability + new CLI. `mirror config check` and `mirror version plan minor` are gates.
- **Help:** `rx --help` delegates to the payload, so no new Cobra help text is authored for `rx`; the payload's help already includes the updated manifest docs for `userOnly`.

## Data Model

- No database. Catalog is YAML on disk.
- `userOnly` is not persisted elsewhere; it is in-memory `ResolvedCommand.UserOnly` after composition.

## Cache and Background Work

- No new cache. Existing `~/.guiho/runx/cache.json` and lease logic unchanged.
- Bare `rx` does not schedule extra workers beyond what the delegated payload schedules via `scheduleLifecycle`.

## AuthN / AuthZ

- Catalog file write access is the authz boundary.
- The guard is a policy flag, not an identity proof. Agents are told to respect it. No principal detection is introduced in v1.

## Risks and Mitigations

- **Risk:** `rx` arg translation mistakes global flag handling. Mitigation: delegate raw args with just `list`/`run` prefix; never parse manifest in `rx`.
- **Risk:** Installer partial failure leaves one launcher updated and the other stale. Mitigation: transactional staging — both launchers downloaded and verified before any activation; rollback restores both from backups.
- **Risk:** Manifest backward compat breaks old catalogs. Mitigation: `userOnly` is optional, default `false`; strict decoder still rejects unknown fields but this one becomes known.
- **Risk:** Help/version drift if `rx` hardcodes version. Mitigation: `rx` has no version string — it prints whatever the payload reports; both launchers report the same `version`.

## Self-Review

| Finding | Severity | Evidence | Recommendation | Resolved |
|---|---|---|---|---|
| Dual-launcher rollback must be atomic for both binaries | blocker | Installer currently backs up one launcher | Back up and restore both launchers; verify both `--version` after activation | design locks, plan will enforce |
| `userOnly` on groups must be rejected | high | Groups and leaves share `Command` type | Add explicit leaf-only validation | locked |
| `reveal` must not be blocked | high | User wants copy-paste hatch for guarded commands | Exempt `reveal`, `list`, `describe`, `check` from guard | locked |
| Shared `~/.guiho/.temp` and `~/.guiho/bin` must not be deleted on uninstall | blocker | Convention 0001 shared ownership | Uninstaller removes only `runx(.exe)` and `rx(.exe)` plus `~/.guiho/runx/` | locked |
| Build matrix must remain pure-Go, `CGO_ENABLED=0` | medium | Current `build-binaries.go` does 8 payloads | Extend to 8 rx-launchers with same flags; keep verify script aware | plan will cover |

Blockers and highs are resolved in this architecture before human validation. Remaining medium item is covered by the plan's validation signals.

## Open Architecture Questions

- Should `--dry-run` on a guarded command refuse (proposed yes) or render the plan? Leaning refuse — finalize in plan review.
- Should `artifacts.json` model two launcher artifacts explicitly per platform? Yes — each launcher is an artifact with its own `file` and `sha256`.

## Handoff to Planning

Next: `docs/plans/2026-08-30-rx-and-user-only-commands.md` — sequenced units for manifest, executor, rx launcher, installers/build, docs/skill, and release. Requires human validation of this architecture before planning execution is considered valid.
