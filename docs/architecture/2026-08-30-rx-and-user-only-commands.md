---
name: RX and User-Only Commands Architecture
purpose: Define the technical design for the rx short alias and the informational userOnly flag
description: Locks repository boundaries, launcher delegation, installer transactions, manifest extension, and docs/versioning for both features (no CLI enforcement)
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
- **User-only flag (informational):** a new strict-YAML per-command field `userOnly` (bool, optional, default `false`). There is **no CLI guard** — if an agent runs a `userOnly` command it still runs. The agent skill (`guiho-s-runx`) teaches agents to warn the user (“Hey, this is user-only, only a user should run this; if you insist I can run it”).

Both reuse the Go/Cobra stack, convention 0001 lifecycle, and Mirror/XDocs obligations.

## Source Material

- [Requirements — RX and User-Only Commands](../requirements/2026-08-30-rx-and-user-only-commands.md)
- [Brainstorm — 2026-08-30](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- `guiho-convention-0001-cli.md` — stable launcher, versioned payloads, `~/.guiho/bin`, `~/.guiho/runx`, atomic `current.json`, shared `.temp`, checksums, `artifacts.json`
- Current: `cmd/runx-launcher/main.go`, `pkg/launcher/launcher.go`, `pkg/installstate/*`, `pkg/updater/*`, `cmd/root.go`, `pkg/manifest/*`, `devops/install.sh|ps1`, `devops/build-binaries.go`

## System Map

### Repositories and Components

- Single repository `C:/GUIHO/runx` — no new repository. All changes are in-place.
- Components touched:
  - `cmd/runx-launcher` (reference for `cmd/rx`), `cmd/rx` (new), `cmd/*` (help/version delegation is not duplicated — rx does not own its own Cobra tree).
  - `pkg/installstate` (pointer remains one file; helpers for rx launcher/payload names).
  - `pkg/launcher` (payload resolution for rx).
  - `pkg/manifest` (type + strict parser + composition + resolve).
  - `devops/` installers, `devops/build-binaries.go`, `devops/verify-release-assets.go`.
  - `skills/guiho-s-runx/SKILL.md`, `embed/skills/guiho-s-runx.SKILL.md`, `prompts/guiho-i-runx.md` (agent guidance).
  - XDocs descriptors: `cmd/rx/*.xdocs.md`, `pkg/manifest/*.xdocs.md`, etc.

### External Dependencies

- GitHub Releases API for channel/version resolution (unchanged).
- No new cloud provider, database, Valkey, or paid service.

## Decisions

### ADR-001: RX as a Thin Launcher That Delegates to the RunX Payload

- **Decision:** `rx` is a minimal launcher binary that (a) inspects its own invocation args, (b) translates bare → `list` and `<selector> …` → `run <selector> …`, and (c) execs the active RunX payload resolved from `~/.guiho/runx/current.json`. Version/help flags are passed through unchanged. `rx` does not embed its own Cobra tree or manifest parser.
- **Why:** Satisfies user's "just run RunX under the hood" instruction, eliminates drift, keeps verification cheap (one payload self-test covers both), and matches convention 0001 launcher discipline.
- **Alternatives:** Full second Cobra binary duplicating `cmd/root.go` (rejected); shell alias/wrapper script (rejected — not a native launcher).
- **Consequences:** `rx` binary is tiny; any `runx run` behavior is automatically inherited.
- **Reversal:** Could promote `rx` to its own payload later without changing the installed `~/.guiho/bin/rx` path.

### ADR-002: Single Pointer Governs Both Launchers

- **Decision:** Keep one `~/.guiho/runx/current.json` (`protocol:1, active, previous`) as the source of truth for both `runx` and `rx`.
- **Why:** One atomic activation for both binaries; rollback restores both to the same version.
- **Alternatives:** Separate `rx.json` pointer (rejected — allows version skew).

### ADR-003: Two Launchers, One Payload Family

- **Decision:** Release matrix produces two launcher assets (`runx-launcher-<platform>` and `rx-launcher-<platform>`) but one coherent payload family (`runx-payload-<platform>`) that both launchers execute.
- **Why:** Minimizes artifact count, keeps `artifacts.json` clear.
- **Alternatives:** Duplicate `rx-payload` (rejected).

### ADR-004: Per-Command `userOnly` Boolean (Informational)

- **Decision:** Add optional `userOnly?: boolean` (strict YAML, default `false` when omitted) to `pkg/manifest/types.go#Command`. Excluded from group nodes by validation. Propagated into `ResolvedCommand.UserOnly`. No CLI enforcement.
- **Why:** Readable, additive, backwards compatible. `allowAgent` is the inverse and more confusing.
- **Alternatives:** `allowAgent: boolean` (rejected), `executionPolicy` (rejected).
- **Validation:** Strict decoder rejects unknown values; type must be boolean; non-boolean or group usage is an error. No executor guard is added.
- **Schema:** `schemas/runx.schema.json` is a config schema (agent.evolution) — no change needed for the manifest flag. Manifest validation is Go strict YAML.

### ADR-005: Informational Only — Skill Teaches, CLI Does Not Block

- **Decision:** There is no execution-time refusal. `userOnly` is surfaced via `runx list --format json` / `describe` / `check` and documented in the skill. The skill script is: “Hey, this is user-only, only a user should run this. Tell the user; if the user insists after the warning, you may execute.”
- **Why:** Matches your clarification: “We don't need to prevent the agent from running. If the agent runs, we don't care. We will just inform the agent.”
- **Alternatives:** Exit 2 refusal (rejected — was in earlier draft, now removed per your guidance).

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
      args = append([]string{"run"}, args...)
  }
  resolve pointer via installstate.ReadPointer()
  payload := launcher.PayloadPath(pointer)
  delegate(payload, args...)
  ```
- **Verification:** `rx --version` and `rx --help` verified like `runx` post-activation. `rx __self-test` not required.

### Manifest and Resolve

- **Types:** `Command` gains `UserOnly *bool `yaml:"userOnly,omitempty" json:"userOnly,omitempty"``; `ResolvedCommand` gains `UserOnly bool`.
- **Parser:** strict decoder + validation: `userOnly` allowed only on leaf commands; groups with `userOnly` → error.
- **Composition:** `composition.go` propagates `userOnly` into the resolved catalog; no inheritance.
- **No guard:** `pkg/executor` and `cmd/catalog.go` remain unchanged — `userOnly` does not affect spawn.

### Installers and Updaters

- **Build:** `devops/build-binaries.go` builds `runx-launcher-*` + `rx-launcher-*` + `runx-payload-*` (CGO_ENABLED=0).
- **Staging:** verify both launchers + payload.
- **Activation:** atomically place both launchers, verify both `--version`.
- **Rollback:** restore both launchers + pointer on failure.
- **Uninstall:** remove both `runx(.exe)` and `rx(.exe)`, never delete shared parents.

### Docs and Release

- **Agent skill:** `skills/guiho-s-runx/SKILL.md` adds `rx` quick-start and documents informational `userOnly` with warning script vault. No CLI exit-code change.
- **XDocs:** new descriptor `cmd/rx/rx.xdocs.md`.
- **Mirror:** minor bump — new public CLI + new manifest flag.

## Data Model

- No database. `userOnly` is in-memory `ResolvedCommand.UserOnly` after composition.

## Risks and Mitigations

- **Risk:** `rx` arg translation mistakes global flags. Mitigation: delegate raw args with just `list`/`run` prefix.
- **Risk:** Installer partial failure leaves one launcher stale. Mitigation: transactional both-launchers verification.
- **Risk:** Help/version drift. Mitigation: `rx` has no version string — it prints payload's version.

## Self-Review

| Finding | Severity | Evidence | Recommendation | Resolved |
|---|---|---|---|---|
| Dual-launcher rollback must be atomic | blocker | Installer currently backs up one launcher | Back up both; verify both `--version` | design locks |
| `userOnly` on groups must be rejected | high | Groups and leaves share `Command` type | Leaf-only validation | locked |
| No CLI enforcement (informational) | high | Your clarification | Skill teaches warning, CLI still executes | locked |
| Shared `~/.guiho/.temp` and `~/.guiho/bin` must not be deleted | blocker | Convention 0001 | Uninstaller removes only `runx(.exe)` and `rx(.exe)` | locked |

## Open Architecture Questions

- None blocking. Should `artifacts.json` model two launcher artifacts explicitly per platform? Yes.

## Handoff to Planning

Next: `docs/plans/2026-08-30-rx-and-user-only-commands.md` — updated to informational RX-1 + RX-3 + skill units. Requires human validation before execution resumes.
