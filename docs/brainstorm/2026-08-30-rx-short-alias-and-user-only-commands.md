---
name: RX Short Alias and User-Only Commands Brainstorm
purpose: Preserve the complete brainstorm for the rx alias and user-only command restriction features before planning
description: Exhaustive record of the rx short CLI for RunX and the user-only execution guard requested for agent safety
created: "2026-08-30T15:00:00Z"
flags:
  - complete
tags:
  - brainstorm
  - cli
  - runx
  - rx
  - alias
  - user-only
  - agent-safety
keywords:
  - rx
  - runx
  - runx list
  - runx run
  - short alias
  - launcher
  - installer
  - user-only
  - agent guard
owner: runx-brainstorm
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RX Short Alias and User-Only Commands Brainstorm

## Context

- Project: RunX — open-source Go/Cobra CLI for a documented, language-agnostic `runx.yaml` command catalog. Executable entrypoint `main.go`, `cmd/` owns Cobra tree, `pkg/` owns manifest/execution/update, native lifecycle via `devops/install.sh` / `devops/install.ps1` and stable launcher layout under `~/.guiho`.
- User has already added the `## GUIHO Conventions` section to `AGENTS.md` and required the agent to read `guiho-convention-0001-cli.md` before any work. That gate is now satisfied and captured in `C:/GUIHO/runx/AGENTS.md`.
- User wants two new features in the same initiative:
  1. `rx` — a short CLI alias for RunX.
  2. A per-command opt-in that marks certain catalog entries (described as "URLs that must be run only by the user" / scripts) as human-only, with an agent-friendly refusal message.
- User explicitly said brainstorming is done ("We are done with the conventions. We are done with the brainstorming.") and wants the material written to a dated brainstorm file, then wants the agent to write the specification/architecture/plan (described as "schema") autonomously for later human approval.

### References Inspected

- `C:/GUIHO/guiho/conventions/guiho-convention-0001-cli.md` — GUIHO CLI Convention (mandatory for RunX)
- `C:/GUIHO/runx/AGENTS.md` — current RunX agent instructions + newly added conventions section
- `C:/GUIHO/runx/main.go`, `cmd/root.go`, `cmd/runx-launcher/main.go`, `pkg/launcher/launcher.go`, `pkg/installstate/*` — current launcher/payload model
- `C:/GUIHO/runx/devops/install.sh`, `C:/GUIHO/runx/devops/install.ps1` — current transactional installers
- `C:/GUIHO/runx/runx.yaml`, `pkg/manifest/types.go` — current manifest v2 catalog shape

## Detailed Brainstorm

### Feature 1 — RX Short Alias

#### Core Idea (Verbatim Intent)

- RX will be a short CLI for RunX.
- RunX CLI does a lot of things, but RX CLI will just be: when I run RX with nothing, absolutely nothing, it will just list. It is equivalent to `runx list`.
- When I run RX with another argument, it is equivalent to `runx run`.
- For the sake of being compliant, RX will also support the commands `-v` and `--version`, and also `-h` and `--help`, and will be exactly similar to the equivalent commands on the RunX side. It doesn't even need to implement that stuff again. It can just run RunX under the hood, and that's all.
- That's really what this needs to do, so it is kind of an alias to make RunX even easier.
- If I run without an argument, it is just to list everything, and if I run with an argument, it is just to run anything. That's it, the easiest way.
- When installing RunX, those two binaries will be installed and added to the path, of course, because the directory is also already on the path. Both RunX and RX will be installed when installing, upgrading, and uninstalling.

#### Behavior — Restated Exhaustively

- `rx` with no arguments, no flags, no extra tokens → equivalent to `runx list`. User said "when I run RunX with nothing, absolutely nothing, it will just list. It is equivalent to `RunX list`" and "If I run without an argument, it is just to list everything" — interpreted as `rx` bare invocation lists the catalog.
- `rx <selector> [child args...]` → equivalent to `runx run <selector> [child args...]`. User said "when I run RX with another argument, it is equivalent to RunX run" and "if I run with an argument, it is just to run anything." That includes forwarding of post-selector child arguments exactly as `runx run` does today (lossless, shell-safe, with `--` delimiter handling).
- `rx -v` / `rx --version` → exactly similar to `runx -v` / `runx --version` (raw SemVer, no extra text).
- `rx -h` / `rx --help` → exactly similar to `runx -h` / `runx --help` (help for the current command and nothing else).
- User explicitly said RX doesn't need to reimplement those; it can just run RunX under the hood. So implementation may be a thin delegate/wrapper that execs the real `runx` payload with translated arguments.
- The user also said RX should support those version/help commands "for the sake of being compliant" — implying convention 0001 compliance (`-v/--version` top-level, `-h/--help` on every command, `--help-tree`/`--help-docs` inherited behavior is already handled by Cobra root, but RX's help should mirror RunX).
- Edge considerations preserved without deciding: Should `rx --help-tree`, `rx --help-docs`, `rx --color`, `rx --help-tree-depth` also delegate? User only explicitly named `-v/--version` and `-h/--help`, but said "exactly similar to the equivalent commands on the RunX side" — could be interpreted as full parity. Keep as open question for spec.
- Edge: `rx --version` with bare invocation? Today `runx --version` is top-level flag handled in `PersistentPreRunE`. Delegation should preserve that.
- Edge: `rx` bare should still perform the same welcome/update bootstrap that `runx` bare does? Or should it strictly list? User said `rx` bare = `runx list`, not `runx` welcome. That is a deliberate divergence for ergonomics. Preserve both statements.
- User wants RX to be "the easiest way" — minimal learning curve, no extra flags to remember.

#### Distribution and Lifecycle — Restated Exhaustively

- When installing RunX, those two binaries will be installed and added to the path, because the directory is already on the path.
- Both RunX and RX (user said "Rex" once, interpreted as RX) will be installed when installing, upgrading, and uninstalling.
- Implications preserved without solving: `~/.guiho/bin/` already on PATH, so installing a second launcher `rx` / `rx.exe` there satisfies PATH automatically. Need to consider:
  - `devops/install.sh` and `devops/install.ps1` must download/verify/install/activate both payload + launcher pairs.
  - `pkg/installstate`, `pkg/launcher`, `pkg/updater` transactional model must handle both.
  - `cmd/runx-launcher` currently builds `runx` launcher; need parallel `rx` launcher (or shared launcher that inspects `os.Args[0]`).
  - Upgrades must keep both launchers backwards-compatible, atomically activated, verified with `--version`, rolled back together on failure.
  - Uninstall must remove both launchers, both versioned payloads, and preserve shared `~/.guiho` correctly (not delete shared `bin/` or `.temp`).
  - `devops/build-binaries.go` must produce doubled matrix.
  - `artifacts.json` / `checksums.txt` must declare both launchers and both payloads.
  - `mirror.yaml` / `xdocs.yaml` must remain valid.

#### Naming

- User called it "rx" and once "Rex" (voice-to-text variant of rx). Canonical is `rx` (lowercase, two letters). Binary names: `rx` on POSIX, `rx.exe` on Windows. Keep consistent with `runx` launcher naming.

### Feature 2 — User-Only Commands (Agent Guard)

#### Core Idea (Verbatim Intent)

- The second feature will be the ability for the users to specify the URLs that must be run only by the user.
- If an agent attempts, just tell the agent: "Hey, there are some URLs that must be run only by the user, and you do not run" — there are some scripts.
- So, those two figures — interpreted as "those two features".

#### Behavior — Restated Exhaustively

- Users can mark certain catalog entries (user described as URLs, interpreted as commands/scripts in `runx.yaml`) as requiring direct human execution.
- If an agent attempts to run such a command, RunX should refuse and tell the agent: "Hey, there are some URLs that must be run only by the user, and you do not run" (exact wording to be finalized in spec; preserve the intent: friendly agent-facing refusal, not a stack trace).
- Scope preserved without solving: Should this apply to `runx run`, `rx` (which delegates to `runx run`), and `runx reveal`? Probably only `run`, but keep as open question.
- Manifest design preserved without solving: New per-command boolean/field such as `agent: false`, `humanOnly: true`, `runBy: user`, or `allowAgent: false` — to be decided in spec/architecture. Must be strict YAML, validated, defaulting to agent-allowed to preserve backwards compatibility.
- Agent detection: How does RunX know it's an agent vs human? Environment hint? Or simply that any `runx run` invocation that hits a guarded command should explain the guard, and the agent skill (`guiho-s-runx`) teaches agents to respect it? User's message implies the guard is for agents, not a hard OS user check.
- Preserve user's phrase "There are some scripts" — suggests some commands are side-effecting / credentialed / destructive and should not be auto-run by AI.
- User wants spec/schema written autonomously after brainstorm; that spec will lock the exact field name, default, error message, exit code, and interaction with `confirm` and `rx`.

### Combined Initiative Notes

- User said "We are done with the conventions. We are done with the brainstorming." — signals planning can begin.
- User explicitly requested: prefix the brainstorm file with the date, write everything down, then go on to planning phase / writing the schema. "Write the schema yourself. Write everything, and then tell me. I will approve."
- That is an autonomous drafting instruction: agent should produce the full planning artifacts (spec, architecture, plan) without blocking for questions, then present for explicit human validation before execution (Phase 1 gate).
- No production deployment, package publishing, or tag pushing until explicitly approved.

### Open Ideas and Unresolved Thoughts (Preserved, Not Decided)

- Should `rx --help-tree` / `rx --help-docs` delegate identically to RunX, or is `-h/--help` and `-v/--version` sufficient for v1?
- Should `rx` bare also print the welcome window + cached update notice, or strictly `runx list`? User said strictly list.
- Should `rx` support `rx --cwd`, `rx --config`, `rx --format json`, `rx --dry-run`? Today `runx list` and `runx run` accept those; delegation can pass them through, but bare `rx` ergonomics may want minimal flags.
- Should user-only guard be enforced in `runx` only, or also in `rx` (which is just a delegate)?
- Exact refusal message, exit code, and whether refusal should be on stdout/stderr; whether `--yes` bypasses guard (likely not).
- Whether `runx reveal` should also be blocked for guarded commands or still allowed to show the command text for human copy-paste.
- Whether guard should be per-command, per-group, or per-catalog (namespace-level).

## Chronological Capture

### Session — 2026-08-30

**User message 1 (verbatim, lightly repaired for readability):**

> Okay, we are going to implement a new feature which will be called "rx". And then the second feature will be the ability for the users to specify the URLs that must be run only by the user. If an agent attempts, just tell the agent: "Hey, there are some URLs that must be run only by the user, and you do not run" There are some scripts. So, those two figures Let's start with brainstorming. The first one, RX, will be a short CLI for RunX. RunX CLI does a lot of things, but RunX CLI will just be: when I run RunX with nothing, absolutely nothing, it will just list. It is equivalent to `RunX list`. And then, when I run `RX` with another argument, it is equivalent to `RunX run`. For the sake of being compliant, Rex will also support the commands `-v` and `--version`, and also `-h` and `--help`, and will be exactly similar to the equivalent commands on the RunX side. It doesn't even need to implement that stuff again. It can just run RunX under the hood, and that's all. That's really what this needs to do, so it is kind of an alias to make RunX even easier.
>
> If I run without an argument, it is just to list everything, and if I run with an argument, it is just to run anything. That's it, the easiest way. When installing RunX, those two binaries will be installed and added to the path, of course, because the directory is also already on the path. Both RunX and Rex will be installed when installing, upgrading, and uninstalling. We are done with the conventions. We are done with the brainstorming. Now write all that down to a brainstorming file. Prefix the file with the date. I want this to have some date, and then let's go on to the planning phase or to writing the schema. Write the schema yourself. Write everything, and then tell me. I will approve.

**Agent actions in this session:**

- Scanned `C:/GUIHO/guiho/conventions` and `C:/GUIHO/runx/AGENTS.md`, selected and applied `guiho-convention-0001-cli.md` as mandatory, patched `## GUIHO Conventions for This Project`.
- Created this dated brainstorm file `docs/brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md` preserving both features verbatim, organized without loss.
- Next steps per user instruction: autonomously draft specification, architecture, and implementation plan for both features, then present for explicit human approval before any Phase 2 execution.

## References

- `C:/GUIHO/guiho/conventions/guiho-convention-0001-cli.md` — inspected, applied as mandatory convention
- `C:/GUIHO/runx/AGENTS.md` — updated with `## GUIHO Conventions for This Project` mandatory block
- `C:/GUIHO/runx/docs/brainstorm/brainstorm.xdocs.md` — descriptor for brainstorm docs
- User-supplied feature names: `rx` (also transcribed as "Rex") and "user-only URLs/scripts"
- `main.go`, `cmd/root.go`, `cmd/runx-launcher/*`, `devops/install.sh`, `devops/install.ps1` — current launcher/installer baseline for rx distribution design

## Questions and Answers

No questioning round was requested. Per `guiho-s-0010-brainstorm`, no questions were asked.

