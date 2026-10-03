---
name: RunX Agent Instructions
purpose: Give coding agents the durable repository context and operating constraints for RunX.
description: Defines package ownership, command behavior, parent coordination, XDocs, and Mirror rules.
created: 2026-07-12
flags: []
tags:
  - agents
  - cli
keywords:
  - runx
  - xdocs
  - mirror
owner: runx
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Agent

## Repository Notes

- RunX is the open-source Go/Cobra CLI for a documented,
  language-agnostic `runx.yaml` command catalog.
- The executable entrypoint is `main.go`; `cmd/` owns the single Cobra tree and
  `pkg/` owns manifest, execution, update, maintenance, lifecycle, and upgrade behavior.
  npm is not an executable or lifecycle distribution; installation is owned by
  the native lifecycle scripts and the stable launcher layout under `.guiho`.
- Use the Go toolchain for formatting, tests, vetting, builds, and release
  compilation. TypeScript/Bun sources are retained as legacy reference only.
- `runx` without arguments first performs local-only, idempotent agent
  bootstrap: install the embedded skill in both global tool locations and
  reconcile the bounded RunX block in both existing repository instruction
  files, the single existing file, or a new root `AGENTS.md`. It then prints
  `Hello Windows - runx v<version>` and may append only a validated cached
  update notice. Help, version, agent-management, uninstall, and non-repository
  paths do not perform repository bootstrap.
  `runx list` lists a configuration; only
  `runx run <selector>` executes a catalog command.
- RunX-owned `run` options precede the selector; post-selector tokens are child
  arguments and must be forwarded without reinterpretation.
- Cobra owns argument parsing and command routing. Only `-h`/`--help` and root
  `-v`/`--version` have short aliases.
- Configuration resolves by `--config`, effective cwd `runx.yaml`, then
  `~/.guiho/runx/runx.yaml`; never search parent directories.
- Treat manifests as trusted executable code. Listing, describing, checking,
  and dry runs must never spawn a configured command.
- Keep the bundled `skills/guiho-s-runx/SKILL.md` aligned with the CLI contract.
- Generated `library/`, `bin/`, `bundle/`, and `vendor/` outputs are ignored;
  never edit them manually.
- Do not publish packages or push tags unless the user explicitly requests a
  release. CI validates pull requests and `main`; protected version tags run the
  `production` workflow that publishes native assets and npm through OIDC.

## GUIHO Project

### Identity

| Field | Value |
| --- | --- |
| GUIHO Project ID | `g0000` observed in current GUIHO runtime artifacts; confirm before formal registry use |
| GUIHO Subject ID | Formal subject ID not declared |
| GUIHO Subject Name | RunX |
| Project Family | `guiho` |
| Repository Directory | `C:\GUIHO\runx` |
| Repository Kind | shared package |
| Parent Project | GUIHO Root (`C:\GUIHO\guiho`) |
| Parent Component | GUIHO Root |

### Component Purpose

RunX owns the reusable command-catalog manifest, standalone CLI, bundled agent
skill, native installers, and package-local documentation.

### Parent Context

- Parent instructions: [../guiho/AGENTS.md](../guiho/AGENTS.md)
- Parent TODO: [../guiho/todo.md](../guiho/todo.md)
- Local TODO: [TODO.md](TODO.md)

### Commands

- Format: `gofmt -w main.go cmd pkg embed devops`
- Tests: `go test ./...`
- Vet: `go vet ./...`
- Build: `go build ./...`
- Compile release asset matrix: `go run devops/build-binaries.go --version <version> --commit <commit> --build-date <RFC3339>`
- Verify release assets: `go run devops/verify-release-assets.go`

## GUIHO Conventions

User conventions are stored under `conventions/` in repository `CGuiho/guiho`.

* **Local:** `C:/GUIHO/guiho/conventions` — read directly, e.g. `cat /c/GUIHO/guiho/conventions/<file>.md`
* **Remote:** if `C:/GUIHO/guiho` is not present locally, the repository is remote at `github.com/CGuiho/guiho` — read via GitHub CLI, e.g. `gh api repos/CGuiho/guiho/contents/conventions --jq ".[].name"` or `gh api repos/CGuiho/guiho/contents/conventions/<file>.md --jq .content | base64 -d`

Check local first, then remote via `gh` if missing.

## GUIHO Conventions for This Project

> **MANDATORY — READ BEFORE ANY WORK:** Every agent working on this project **MUST** read `guiho-convention-0001-cli.md` **before doing anything else** — planning, brainstorming, architecture, implementation, review, validation, or release.
>
> * **Local:** `C:/GUIHO/guiho/conventions/guiho-convention-0001-cli.md` — `cat /c/GUIHO/guiho/conventions/guiho-convention-0001-cli.md`
> * **Remote fallback (if `C:/GUIHO/guiho` missing):** `gh api repos/CGuiho/guiho/contents/conventions/guiho-convention-0001-cli.md --jq .content | base64 -d`
> * **Why this one:** RunX is the GUIHO Go/Cobra CLI (`main.go` + `cmd/` + `pkg/`). `guiho-convention-0001-cli.md` (GUIHO CLI Convention) is the sole authority for technology stack (Go + Cobra), mandatory tooling (`mirror.yaml` / `runx.yaml` / `xdocs.yaml`), flags, help tree, installation lifecycle (`devops/install.sh` / `install.ps1`), stable launcher + versioned payloads, CLI home (`~/.guiho/runx`), configuration, agent artifacts, startup update check, and release channel contracts. No CLI work in this repository is valid without it.
>
> **Do not proceed with any task until you have read and understood this convention.**

| Priority | Convention | Path | Applies To |
|---|---|---|---|
| **1 — MANDATORY** | GUIHO CLI Convention | `guiho-convention-0001-cli.md` | Every change in this repository |
| 2 — Reference | GUIHO SWE Convention | `guiho-convention-0000-swe.md` | Lifecycle/phase work when relevant |
| 2 — Reference | GUIHO Agent Artifacts Convention | `guiho-convention-0002-agent-artifacts.md` | Skills/prompts/instructions (`guiho-s-*`, `guiho-p-*`, `guiho-i-*`) |

## GUIHO Essentials

This project uses Essentials to handle its foundational work.
Essentials is a swarm of essential, reusable agent skills you can use in every
project you work on. To work with Essentials, read the agent skill
`guiho-s-essentials`. It catalogs every Essentials skill, what it helps you do,
and when to use it.

## GUIHO Mandume

This project uses Mandume to handle its SWE work.
Mandume is the software-engineering swarm. Load `guiho-s-mandume` and use the
`## Mandume` section below as the current worker authority. `MANDUME.md` retains
inactive historical guidance only. Read the actual parent
[Convention 0011](../guiho/conventions/guiho-convention-0011-agent-readiness.md)
before readiness work; a handoff summary is not a substitute.

## Semantic Project Versioning -- GUIHO Mirror

Use the `guiho-s-mirror` skill whenever versioning, tags, changelogs, or release
configuration changes. Read `mirror.yaml`, run `mirror version plan`
before applying a version, and never hand-edit Mirror-managed version fields.

<!-- BEGIN XDOCS — DO NOT EDIT THIS SECTION -->
## XDocs Structured Documentation

This project uses **xdocs** (`@guiho/xdocs`) for structured, machine-readable
documentation. The repository has one root `XDOCS.md` index (no frontmatter),
and each package/application has a root named `*.xdocs.md` descriptor file. Each
documented module has exactly one named `*.xdocs.md` descriptor in its directory
with YAML frontmatter (`subject`, `description`, `parent`, `children`,
`files`, `documents`, `tags`, `keywords`, `flags`). Same-directory plain
`*.md` files are companion documents and must be listed in the descriptor's
`documents` metadata map. Ordinary companion documents should also include
frontmatter with `owner`, `tags`, and `keywords` so agents can inspect
metadata before reading full Markdown bodies.

**Load the `guiho-s-xdocs` agent skill** for any documentation work:
creating, updating, regenerating, scanning, merging, or navigating xdocs descriptors.
The skill holds the full workflow, metadata schema, and CLI reference.

Before changing documentation, read `xdocs.config.toml` and respect `[ai].mode`:

- **prompt** — announce which xdocs descriptors need updating and wait for confirmation.
- **auto** — update the relevant xdocs descriptors immediately.

Use the installed xdocs CLI for operations. Prefer `xdocs context "<query>"
[path] --documents --files --format json` to get a task-specific reading set,
or `xdocs meta [path] --documents --format json` when you only need
frontmatter. Other commands: `xdocs scan`, `xdocs tree`, `xdocs generate`,
`xdocs list`, `xdocs doctor`, `xdocs merge`, `xdocs upgrade`, and
`xdocs uninstall --dry-run`.
<!-- END XDOCS -->

<!-- BEGIN GUIHO MIRROR - DO NOT EDIT THIS SECTION -->
## Semantic Project Versioning -- GUIHO Mirror

Invoke the guiho-s-mirror agent skill every time the user wants to bump, tag, release, plan, initialize, configure, or troubleshoot semantic project versioning with GUIHO Mirror.

Before editing release docs or changelogs, inspect mirror.config.toml. If [agents].write_changelog is false, skip changelog edits. If it is missing or true, changelog edits are allowed when the project has a changelog.

Use [agents].changelog_path as the changelog file path. If it is missing, use CHANGELOG.md in the project root.
<!-- END GUIHO MIRROR -->

<!-- BEGIN MIRROR — DO NOT EDIT THIS SECTION -->
## GUIHO Mirror Instruction Block

Run plain `mirror` once in a repository to verify the global Mirror skill and
this bounded instruction block. Repeated runs are idempotent.

Use `mirror version plan <target>` and `mirror version apply <target>` for semantic versioning.
`mirror init` defaults to `v{version}` tags and enables release commits and
pushes; explicit interactive or flag selections remain authoritative.

When `mirror.yaml` defines hooks, follow AI instructions only at the
agent-controlled everything, plan, and apply boundaries. Treat command hooks as
repository code: pass `--run-hooks` or `--skip-hooks` only with explicit
authorization, independently of `--yes`.
<!-- END MIRROR -->

<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->
## RunX Command Catalog

Load the `guiho-s-runx` skill whenever discovering commands, creating or
updating catalog entries, validating `runx.yaml`, inspecting command details,
or executing RunX commands.
Start with `runx check --format json` and `runx list --format json`, select
stable UIDs, use `runx describe <uid>`, and run
`runx run --dry-run <uid>` before unfamiliar or side-effecting work.
RunX options precede the selector; post-selector tokens belong to the child.
<!-- END RUNX -->
## Mandume

GUIHO RunX.

Managed by the GUIHO Mandume swarm ([CGuiho/mandume](https://github.com/CGuiho/mandume)); the full worker-registry example lives at `example/AGENTS.md` there.

**OpenCode is the only currently supported agent harness.** Always use its
built-in native subagent tool. Select an authorized agent whose actual child
permissions and capabilities cover the brief: General can perform broad
multi-step shell/write work; Explore is read-only by default. A parent's authority
does not prove the child's permissions. Never widen host permissions to bypass a
denial or infer that every native agent is read-only.

### Mode

```yaml
execution: dnd  # dnd | interruptible — orchestrator NEVER stops during execution/review
notifications: off  # human-facing only; native completion notifications remain enabled
harness: opencode  # only current harness; native background subagents always
tmux-session: runx  # orchestrator session on su-57; convention = this project's name
```

### Coordination

- GitHub repository: https://github.com/CGuiho/runx.git
- GitHub Project: https://github.com/users/CGuiho/projects/2 — GUIHO; authoritative for task lists, ownership, scope, priority, Status and Project fields over every local Markdown helper.
- GitHub component: `runx` — exact existing Component option of Project #2, confirmed by live readback. Repository issue ownership and Component are distinct.
- Every task is a real issue in its owning repository, attached to this Project with exactly one nonblank Component; a Project draft alone is insufficient. Read back issue repository/state, Project membership, Component and Status after every change, then mirror the accepted remote state locally.
- To-do file: `TODO.md` (established casing; repo root)
- Current policy task: https://github.com/CGuiho/runx/issues/63; [local requirements](docs/todo/native-background-policy.md).
- Reserved port: pending — reserve in `apps.md` (`CGuiho/guiho`)

### Workers

| Worker | Class | Model (exact OpenCode ID) | Thinking | Usage |
| --- | --- | --- | --- | --- |
| `mastermind` | mastermind | MiMo-V2.6-Pro (`xiaomi/mimo-v2.6-pro`, direct Xiaomi API) | provider default; no variant | CG-authorized default; availability subject to provider/account |
| `engineer` | workhorse | MiMo-V2.6-Pro (`xiaomi/mimo-v2.6-pro`, direct Xiaomi API) | provider default; no variant | CG-authorized default; availability subject to provider/account |

Both roles use the same current model; judgment/labor remain distinct roles.
The canonical authority is
[Convention 0007](../guiho/conventions/guiho-convention-0007-models.md).
The former Muse Spark 1.3 Contributor (`vercel/meta/muse-spark-1.3-contributor`),
GLM 5.3 Flash (`opencode/glm-5.3-flash`) and DeepSeek V4.1 Flash
(`opencode/deepseek-v4-flash`) rows are **history only, not fallback capacity**.
Do not silently substitute models. Ledger a MiMo provider/quota failure, fail
that dependent unit safely and continue independent units. Weekly-pool models
require usage ensured and explicit CG authorization for that run.

Pass the authorized exact provider/model and available reasoning variant
explicitly through the actual live native tool schema after checking the model
catalog and registry. MiMo's variants list is empty: omit the variant parameter
and record provider-default thinking. Thinking enabled is not a differentiated
maximum tier; never invent MiMo `max`/`xhigh`. Native General
`openai/gpt-6.1-sol#xhigh` is explicitly authorized for the current
rollout/readiness task only; it is not a future default or silent fallback.

### Native Background Completion Contract

1. Write a complete owned-unit brief under
   `docs/plans/<plan>/execution/handoffs/<date>-<unit>-<worker>.md`, including
   `Mode: dnd`, actual Convention 0011, owning instructions, task identity,
   model/variant decision, required capabilities, exclusive paths, acceptance,
   checks and commit/push authority. Launch a capable native child in the background.
2. Remain responsive through independent ready work or the harness completion
   event. Native polling, sleeping, shell waits, repeated status loops and
   per-worker file tailing are forbidden. Human notifications off does not turn
   off child completion events or require an artificial keep-alive.
3. On completion inspect actual child/session identity, result, errors/exit
   evidence and timing; review the full scoped diff, safe checks, GitHub
   issue/Project/Component/Status readbacks and owned commits against acceptance.
   Record exact provider/tool/permission failures, select an already-authorized
   capable native agent when possible and continue independent valid units.
4. Integrate coherent owned work on `main` under `guiho-s-0032-git-commit`,
   inspect the full outgoing ancestry before plain push and verify live remote
   equality. Child push requires explicit parent authorization. Immediately
   continue the next ready unit through technical self-review without waiting for CG.

**No active CLI-worker exception exists.** Capability gaps, permission denial,
provider failures or earlier quoted CLI requests do not authorize `opencode run`
workers or host configuration changes. Other harness guidance in `MANDUME.md`
is inactive history. A CLI-worker/tailing contract is dormant future policy only
if CG later authorizes a different harness that truly lacks native subagents;
no such harness is currently enabled. Ordinary Git/gh/Go/RunX CLI use and launching
a primary OpenCode session through its CLI are distinct from worker delegation.

### Contract

- The orchestrator coordinates on `main`, remains responsive, and integrates the native workers above through `guiho-s-0440-hand-off`; each child's actual permissions must cover its brief.
- DND execution and technical review never ask CG or wait for a wake-up. Resolve questions with the safest reversible evidence-grounded choice under `docs/questions/`, and continue. Record actual security/data-loss, impossible-specification or missing-security-authorization blockers under `docs/issues/`; continue independent units. Later human review is a distinct phase.
- Use the Mandume skills (`guiho-s-mandume` + lifecycle skills) and the Essentials skills (`guiho-s-0001-guiho`, `guiho-s-0004-working-with-cg`, `guiho-s-0040-explorer`, `guiho-s-0032-git-commit`). Conventions: `conventions/` in `CGuiho/guiho` (`apps.md` for ports).
