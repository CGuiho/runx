---
name: RX and User-Only Commands Requirements
purpose: Define the product scope and acceptance criteria for the rx short alias and the user-only command guard
description: Specifies the rx bare-list / run delegation, version/help parity, dual-binary lifecycle, and the per-command human-only restriction with agent refusal
created: "2026-08-30T15:30:00Z"
flags:
  - approved
tags:
  - requirements
  - cli
  - rx
  - agent-safety
keywords:
  - rx
  - runx list
  - runx run
  - user-only
  - allowAgent
  - installer
owner: runx-requirements
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RX and User-Only Commands Requirements

## Summary

RunX is the GUIHO Go/Cobra CLI for a `runx.yaml` command catalog. Two complementary features are requested:

1. **`rx` — a short, ergonomic alias for RunX.** Bare `rx` lists the catalog (`runx list`). `rx <selector> ...` runs a catalog command (`runx run <selector> ...`). `-v/--version` and `-h/--help` behave identically to RunX. Installation, upgrade, and uninstall manage both binaries transactionally.

2. **User-only commands — a per-command guard that marks certain catalog entries as human-only.** When an agent attempts to run a guarded command, RunX refuses with a friendly, agent-readable message explaining that the command must be run by the user.

Both ship together as one minor release. No new database, cloud resource, or paid service is involved.

## Source Material

- [Brainstorm — 2026-08-30 RX Short Alias and User-Only Commands](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- `guiho-convention-0001-cli.md` (mandatory — read before any work)
- Current implementation: `main.go`, `cmd/root.go`, `cmd/runx-launcher/*`, `devops/install.sh|ps1`, `pkg/manifest/types.go`, `pkg/executor/*`

## Target Users

- **Primary:** Developers and operators who use RunX daily and want the shortest possible invocation (`rx`).
- **Secondary:** Teams where a `runx.yaml` contains sensitive, destructive, or credentialed commands (deploy, production migration, secret rotation) that must never be auto-executed by an AI agent.
- **Agent consumer:** AI agents using the bundled `guiho-s-runx` skill. Agents must be able to discover the guard and respect it without extra configuration.

## Roles and Permissions

- No new authentication or session model. File-system access to `runx.yaml` governs who can declare a user-only command.
- Catalog author (human) decides which commands are user-only.
- Agent is the restricted actor. The guard is advisory-enforced by the CLI: any invocation hitting a guarded command receives the refusal, regardless of caller. The skill teaches agents to treat the message as a hard stop, while a human can still run the command directly in their shell after seeing the printed command (e.g., via `runx reveal`).

## MVP Scope

### In Scope — RX

- `rx` binary (`rx` on POSIX, `rx.exe` on Windows) installed alongside `runx` via the same transactional lifecycle.
- Bare `rx` → `runx list` with identical flags, output, and exit codes (including `--format json`, `--cwd`, `--config` where applicable).
- `rx <selector> [--] [child args...]` → `runx run <selector> [--] [child args...]` with lossless forwarding (same semantics as RunX, including the one-delimiter removal rule).
- `rx -v` / `rx --version` and `rx -h` / `rx --help` identical to `runx` equivalents. Help-tree and help-docs delegation is included for completeness (no second implementation).
- `rx` delegates by invoking the active RunX payload under the hood; it does not duplicate Cobra parsing, manifest loading, or execution logic.
- Installers (`install.sh`, `install.ps1`), upgrader, launcher activation, pointer (`~/.guiho/runx/current.json`), resource layout, `checksums.txt`, `artifacts.json`, and uninstaller handle both binaries atomically with rollback.

### In Scope — User-Only Guard

- Per-command manifest field that marks a command as user-only (field name `userOnly` proposed; architecture will lock the final name).
- Default is `false` / unset → backwards compatible: existing catalogs without the field are agent-runnable as before.
- When a guarded command is selected for `runx run` (and therefore `rx` delegation), the CLI does **not** spawn the configured command. It prints a friendly agent-readable refusal to stderr and exits with a deterministic non-zero code (proposed `2` — usage/validation error — to be locked in architecture).
- Refusal message must contain the user-requested guidance: that the command is marked as user-only and must be run by the user, and that the agent should not run it. Exact wording to be finalized but must be stable and documented.
- Guard applies to direct UIDs, canonical selectors, shorthand IDs, and numeric indexes — every resolution path that `runx run` supports.
- `runx reveal <selector>` remains allowed for guarded commands so a human can copy the command text without executing it. `runx list` and `runx describe` surface the guard without executing.

### Dependencies

- RX depends on the existing RunX payload and `runx list` / `runx run` contracts.
- User-only guard depends on strict manifest parsing and the existing selector resolution order.

## User Workflows

### Workflow 1 — Quick List with RX

1. User types `rx` in a project directory containing `runx.yaml`.
2. CLI prints the catalog list (same as `runx list`).
3. User scans available commands by index, UID, or selector.

### Workflow 2 — Quick Run with RX

1. User types `rx 2` or `rx deploy -- --dry-run`.
2. CLI resolves selector `2` / `deploy` and runs the stored command exactly as `runx run 2` would, forwarding child args losslessly and preserving exit code.

### Workflow 3 — Author Marks a Command Human-Only

1. Author edits `runx.yaml` and sets the guard field on a sensitive command (e.g., `deploy-prod`).
2. `runx check` passes; `runx list` shows the command but does not run it.
3. Agent later attempts `runx run deploy-prod` or `rx deploy-prod`.

### Workflow 4 — Agent Hits the Guard

1. Agent runs `runx run deploy-prod` (or `rx deploy-prod`).
2. CLI refuses without spawning, prints the user-only message to stderr.
3. Agent reports the refusal to the user and stops; user runs the command manually if desired (or copies via `runx reveal deploy-prod`).

### Workflow 5 — Install / Upgrade / Uninstall

1. Fresh install places both `runx` and `rx` launchers in `~/.guiho/bin` and payload(s) in `~/.guiho/runx/versions/<version>/`, verifies both with `--version` and `__self-test`, atomically activates the pointer.
2. Upgrade replaces both launchers and payloads transactionally; failure rolls back to previous `runx` + `rx`.
3. Uninstall removes both launchers, all versioned payloads, `~/.guiho/runx/`, and installed skills; preserves shared `~/.guiho/bin` and `.temp` parents.

## Functional Requirements

### RX — Functional

- **RX-F1:** `rx` with zero non-flag arguments executes `runx list` semantics and produces byte-identical stdout/stderr/exit behavior.
- **RX-F2:** `rx <selector> [args...]` executes `runx run <selector> [args...]` semantics, including `--cwd`, `--config`, `--format`, `--yes`, `--dry-run`, and child-arg forwarding.
- **RX-F3:** `rx -v` / `rx --version` prints the raw SemVer (same as `runx --version`) and nothing else.
- **RX-F4:** `rx -h` / `rx --help` and `rx <subcommand> --help` render the same help as the equivalent `runx` invocation for the delegated target.
- **RX-F5:** `rx --help-tree` and `rx --help-docs` delegate identically (no divergent tree).
- **RX-F6:** RX does not introduce new top-level Cobra commands beyond delegation; it reuses the RunX payload so help/version output never drifts.
- **RX-F7:** Unknown flags or selectors produce the same exit code and error style as RunX (exit 2 for usage/unknown-flag, 1 for runtime failure).

### User-Only — Functional

- **UO-F1:** A catalog leaf may declare `userOnly: true` (or final field name) via strict YAML.
- **UO-F2:** When a guarded command is resolved for execution, the CLI refuses without spawning, prints the agent-readable refusal, and exits non-zero.
- **UO-F3:** The refusal message states that the command is marked as user-only / must be run by the user and that the agent should not run it.
- **UO-F4:** The guard is visible without execution: `runx list` and `runx describe` / `runx check` surface that the command is user-only (format to be locked in architecture — e.g., JSON field, tag, or flag).
- **UO-F5:** `runx reveal <selector>` for a guarded command prints the exact stored command (stdout) and does not refuse — it is the human copy-paste escape hatch.
- **UO-F6:** Omitting the field or setting `userOnly: false` preserves current behavior (runnable by anyone).

### Lifecycle — Functional

- **LC-F1:** `devops/install.sh` and `devops/install.ps1` install both launchers and all payload resources with SHA-256 and `artifacts.json` verification, staged under `~/.guiho/.temp/runx-install-<id>/`, with rollback on any failure.
- **LC-F2:** Activation atomically updates `~/.guiho/runx/current.json` (protocol 1) and verifies both launchers with `--version`.
- **LC-F3:** `runx upgrade` (or re-running the installer for a newer version) upgrades both binaries together; explicit older-version install is honoured.
- **LC-F4:** `runx uninstall` and both `uninstall.sh|ps1` scripts remove both launchers and all versioned payloads; `--preserve-config`, `--preserve-data`, `--dry-run`, `--yes` behave identically for the combined installation.

## Non-Functional Requirements

- **NF-1 — Convention 0001 compliance:** Flags use full long form; only `-v/--version` and `-h/--help` have short aliases; help-tree, global-flags rules, CLI home, and stable launcher are preserved.
- **NF-2 — No foreground network:** Bare invocations (`runx`, `rx`) perform no network fetch; update checks remain background, single-flight, leased.
- **NF-3 — Determinism:** List, describe, check, reveal, run delegation, and refusal are deterministic given the same catalog and arguments.
- **NF-4 — Safety:** Listing, describing, checking, and revealing never spawn a configured command; guarded runs never spawn.
- **NF-5 — Parity:** `runx list` vs `rx` bare, and `runx run` vs `rx <selector>`, are indistinguishable except for the invoked binary name.

## Data Owned By The Product

- No new database or persistent state.
- `runx.yaml` catalog is the source of truth; the new per-command guard field is part of that trusted executable code.
- `~/.guiho/runx/current.json` pointer now logically governs two launchers but remains one pointer (active + previous).
- `artifacts.json` and `checksums.txt` are the source of truth for ownership and integrity of both launchers/payloads.

## Integrations

- **GitHub Releases:** Each release publishes both `runx` and `rx` launchers + payloads as native assets plus `checksums.txt`, `artifacts.json`, skill ZIP, prompts, instruction, and schemas.
- **Agent skill (`guiho-s-runx`):** Updated to document `rx` ergonomics and the user-only guard, including how agents should handle the refusal.
- **No external APIs** beyond the existing release catalog fetch for upgrades.

## Non-Goals

- No `rx` interactive TUI, shell completions v1, or command-specific help beyond delegation.
- No per-group or per-namespace guard in v1 — only per-command leaf guard.
- No OS-level user identity check — guard is catalog-declared, CLI-enforced.
- No `--allow-agent` override in v1 (deferred; could be added backwards-compatibly if needed).
- No change to `runx run` shell adapter selection or path translation.
- No change to manifest parent/child composition beyond the new guard field.
- No production deployment or DNS/traffic mutation.

## Open Questions (Recorded, Not Blocking MVP)

1. Exact manifest field name: `userOnly` vs `allowAgent` vs `humanOnly`. Prefer `userOnly: true` for readability — architecture will accept.
2. Exact refusal exit code: `2` (usage/validation) vs `3` (policy). `2` keeps it in the usage-error family.
3. Exact refusal wording — finalize to stable, testable string containing "user-only" and "must be run by the user" and "agent should not run it".
4. Whether `runx list --format json` should include a `userOnly` boolean in each `ResolvedCommand` — strongly recommended for discoverability.
5. Whether guarded `runx run --dry-run` should also refuse or should still render the dry-run plan (propose: dry-run also refuses, because even a plan implies intent to run).

## Architecture Inputs

- **Must be strict typed YAML:** New field goes through `pkg/manifest/types.go` and the strict decoder; unknown fields remain rejected.
- **Must update JSON schemas:** `schemas/runx.schema.json` and embedded schema copies must accept the new field.
- **Must be XDocs-aware:** New binary `cmd/rx` and manifest field require descriptor updates.
- **Must be Mirror-aware:** New public capability = `minor` bump; installers + release matrix change belongs to the same release.
- **Known risks:** Dual-binary drift (help/version divergence), installer rollback completeness for two launchers, manifest backward compatibility, and agent skill currency.
