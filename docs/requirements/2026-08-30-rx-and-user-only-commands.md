---
name: RX and User-Only Commands Requirements
purpose: Define the product scope and acceptance criteria for the rx short alias and the informational user-only flag
description: Specifies the rx bare-list / run delegation, version/help parity, dual-binary lifecycle, and the per-command informational userOnly flag taught via the agent skill
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
  - userOnly
  - installer
owner: runx-requirements
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RX and User-Only Commands Requirements

## Summary

RunX is the GUIHO Go/Cobra CLI for a `runx.yaml` command catalog. Two complementary features are requested:

1. **`rx` — a short, ergonomic alias for RunX.** Bare `rx` lists the catalog (`runx list`). `rx <selector> ...` runs a catalog command (`runx run <selector> ...`). `-v/--version` and `-h/--help` behave identically to RunX. Installation, upgrade, and uninstall manage both binaries transactionally.

2. **User-only commands — an informational per-command flag `userOnly` that marks certain catalog entries as user-only.** There is no CLI enforcement: if an agent runs a `userOnly` command it still runs. The agent skill teaches agents to recognize the flag, tell the user “Hey, this is user-only, only a user should run this,” and only execute if the user insists after that warning. `reveal` remains the human copy-paste hatch.

Both ship together as one minor release. No new database, cloud resource, or paid service is involved.

## Source Material

- [Brainstorm — 2026-08-30 RX Short Alias and User-Only Commands](../brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- `guiho-convention-0001-cli.md` (mandatory — read before any work)
- Current implementation: `main.go`, `cmd/root.go`, `cmd/runx-launcher/*`, `devops/install.sh|ps1`, `pkg/manifest/types.go`, `pkg/executor/*`

## Target Users

- **Primary:** Developers and operators who use RunX daily and want the shortest possible invocation (`rx`).
- **Secondary:** Teams where a `runx.yaml` contains sensitive, destructive, or credentialed commands (deploy, production migration, secret rotation) that should be marked as user-only so agents warn before running them.
- **Agent consumer:** AI agents using the bundled `guiho-s-runx` skill. Agents discover `userOnly` via `runx list/describe --format json` and the skill guidance.

## Roles and Permissions

- No new authentication or session model. File-system access to `runx.yaml` governs who can declare a user-only command.
- Catalog author (human) decides which commands are user-only by setting `userOnly: true`.
- Agent is the advisory consumer: the flag is informational. The CLI does not block execution. The skill teaches: see `userOnly: true` → tell the user it is user-only and you should not normally execute it; if the user explicitly insists after the warning, you may execute.

## MVP Scope

### In Scope — RX

- `rx` binary (`rx` on POSIX, `rx.exe` on Windows) installed alongside `runx` via the same transactional lifecycle.
- Bare `rx` → `runx list` with identical flags, output, and exit codes (including `--format json`, `--cwd`, `--config` where applicable).
- `rx <selector> [--] [child args...]` → `runx run <selector> [--] [child args...]` with lossless forwarding (same semantics as RunX, including the one-delimiter removal rule).
- `rx -v` / `rx --version` and `rx -h` / `rx --help` identical to `runx` equivalents. Help-tree and help-docs delegation is included for completeness (no second implementation).
- `rx` delegates by invoking the active RunX payload under the hood; it does not duplicate Cobra parsing, manifest loading, or execution logic.
- Installers (`install.sh`, `install.ps1`), upgrader, launcher activation, pointer (`~/.guiho/runx/current.json`), resource layout, `checksums.txt`, `artifacts.json`, and uninstaller handle both binaries atomically with rollback.

### In Scope — User-Only Flag (Informational)

- Per-command manifest field `userOnly?: boolean` on leaf commands. Optional, default `false` / omitted → backwards compatible.
- Strict YAML validation: leaf-only, groups cannot declare `userOnly`; type must be boolean; unknown fields remain rejected.
- Propagated to `ResolvedCommand.userOnly` and visible via `runx list --format json`, `runx describe`, `runx check` without execution.
- `runx reveal <selector>` remains allowed for `userOnly` commands (human copy-paste).
- Agent skill documents `userOnly` with example, discovery via JSON, and the warning script (“Hey, this is user-only, only a user should run this. If you insist I can run it.”). No CLI exit-code change.

### Dependencies

- RX depends on the existing RunX payload and `runx list` / `runx run` contracts.
- `userOnly` flag depends on strict manifest parsing and the existing selector resolution order.

## User Workflows

### Workflow 1 — Quick List with RX

1. User types `rx` in a project directory containing `runx.yaml`.
2. CLI prints the catalog list (same as `runx list`).
3. User scans available commands by index, UID, or selector.

### Workflow 2 — Quick Run with RX

1. User types `rx 2` or `rx deploy -- --dry-run`.
2. CLI resolves selector `2` / `deploy` and runs the stored command exactly as `runx run 2` would, forwarding child args losslessly and preserving exit code.

### Workflow 3 — Author Marks a Command User-Only

1. Author edits `runx.yaml` and sets `userOnly: true` on a sensitive command (e.g., `deploy-prod`).
2. `runx check` passes; `runx list --format json` shows `"userOnly": true`.
3. Agent later discovers the flag via `list/describe` or the skill before attempting `runx run deploy-prod` or `rx deploy-prod`.

### Workflow 4 — Agent Sees the Flag (Informational)

1. Agent resolves `deploy-prod` and sees `userOnly: true` (or reads the skill).
2. Agent tells the user: “Hey, this is user-only, only a user should run this.” It does not execute unless the user explicitly insists.
3. If the user insists, the agent may execute; otherwise the user runs it manually or copies via `runx reveal deploy-prod`.

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

### User-Only — Functional (Informational)

- **UO-F1:** A catalog leaf may declare `userOnly: true` via strict YAML (boolean). Groups cannot declare it.
- **UO-F2:** Omitting the field or setting `userOnly: false` preserves current behavior (no advisory).
- **UO-F3:** The flag is visible without execution: `runx list --format json` includes `userOnly` in each `ResolvedCommand`; `runx describe` / `runx check` surface it; no execution is needed to discover it.
- **UO-F4:** There is no CLI refusal or exit-code change. `runx run <guarded>` and `rx <guarded>` still execute the stored command (agent skill is the teaching layer).
- **UO-F5:** `runx reveal <selector>` for a `userOnly` command prints the exact stored command (stdout) as the human copy-paste hatch.
- **UO-F6:** The agent skill contains the informational warning script and the “if user insists, you may execute” policy.

### Lifecycle — Functional

- **LC-F1:** `devops/install.sh` and `devops/install.ps1` install both launchers and all payload resources with SHA-256 and `artifacts.json` verification, staged under `~/.guiho/.temp/runx-install-<id>/`, with rollback on any failure.
- **LC-F2:** Activation atomically updates `~/.guiho/runx/current.json` (protocol 1) and verifies both launchers with `--version`.
- **LC-F3:** `runx upgrade` (or re-running the installer for a newer version) upgrades both binaries together; explicit older-version install is honoured.
- **LC-F4:** `runx uninstall` and both `uninstall.sh|ps1` scripts remove both launchers and all versioned payloads; `--preserve-config`, `--preserve-data`, `--dry-run`, `--yes` behave identically for the combined installation.

## Non-Functional Requirements

- **NF-1 — Convention 0001 compliance:** Flags use full long form; only `-v/--version` and `-h/--help` have short aliases; help-tree, global-flags rules, CLI home, and stable launcher are preserved.
- **NF-2 — No foreground network:** Bare invocations (`runx`, `rx`) perform no network fetch; update checks remain background, single-flight, leased.
- **NF-3 — Determinism:** List, describe, check, reveal, and run delegation are deterministic given the same catalog and arguments.
- **NF-4 — Safety:** Listing, describing, checking, and revealing never spawn a configured command. Informational `userOnly` does not change spawn behavior.
- **NF-5 — Parity:** `runx list` vs `rx` bare, and `runx run` vs `rx <selector>`, are indistinguishable except for the invoked binary name.

## Data Owned By The Product

- No new database or persistent state.
- `runx.yaml` catalog is the source of truth; `userOnly` is part of that trusted executable code.
- `~/.guiho/runx/current.json` pointer now logically governs two launchers but remains one pointer (active + previous).
- `artifacts.json` and `checksums.txt` are the source of truth for ownership and integrity of both launchers/payloads.

## Integrations

- **GitHub Releases:** Each release publishes both `runx` and `rx` launchers + payloads as native assets plus `checksums.txt`, `artifacts.json`, skill ZIP, prompts, instruction, and schemas.
- **Agent skill (`guiho-s-runx`):** Updated to document `rx` ergonomics and the informational `userOnly` flag, including the warning script and the “if user insists” policy.
- **No external APIs** beyond the existing release catalog fetch for upgrades.

## Non-Goals

- No `rx` interactive TUI, shell completions v1, or command-specific help beyond delegation.
- No per-group or per-namespace `userOnly` inheritance in v1 — only per-command leaf flag.
- No OS-level user identity check — `userOnly` is catalog-declared, informational.
- No CLI enforcement / refusal for `userOnly` (intentionally informational; skill teaches, CLI still executes).
- No `--allow-agent` override (not needed without enforcement).
- No change to `runx run` shell adapter selection or path translation.
- No change to manifest parent/child composition beyond the new flag.
- No production deployment or DNS/traffic mutation.

## Open Questions (Resolved for MVP)

1. Manifest field name locked: `userOnly: true` (informational boolean).
2. `runx list --format json` includes `userOnly` in each `ResolvedCommand` — implemented.
3. No enforcement: `runx run --dry-run` on a `userOnly` command still renders the plan (because there is no guard).

## Architecture Inputs

- **Must be strict typed YAML:** New field goes through `pkg/manifest/types.go` and the strict decoder; unknown fields remain rejected.
- **Not a config JSON schema change:** `schemas/runx.schema.json` / `runx.global.schema.json` are config schemas (agent.evolution) — no change needed for the manifest `userOnly` flag.
- **Must be XDocs-aware:** New binary `cmd/rx` and manifest field require descriptor updates.
- **Must be Mirror-aware:** New public capability = `minor` bump; installers + release matrix change belongs to the same release.
- **Known risks:** Dual-binary drift (help/version divergence), installer rollback completeness for two launchers, manifest backward compatibility, and agent skill currency.
