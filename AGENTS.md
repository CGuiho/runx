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

# RunX Agent

## Agent

Always read this: /c/GUIHO/superiority/agents/guiho-a-0001-swe.AGENTS.md (C:\GUIHO\superiority\agents\guiho-a-0001-swe.AGENTS.md)
Stop if you can not find it.

## Required CLI Engineering

- Use `guiho-a-0001-swe` as the coordinating GUIHO Software Engineer/SWE agent
  for CLI architecture, planning, execution, review, validation, and release
  work.
- Load and follow the `guiho-s-0035-cli-engineer-go` agent skill whenever creating,
  upgrading, refactoring, reviewing, testing, packaging, installing, or
  releasing the RunX CLI.
- `guiho-s-0035-cli-engineer-go` is a skill, not an agent. It supplements the SWE
  agent and does not replace the lifecycle controller required by that agent.
- During Go CLI work, also load the XDocs and Mirror skills at the boundaries
  defined by the Go CLI engineering skill and the approved migration RFC.
- The approved RFC 0034 migration may make breaking changes. RunX is pre-1.0;
  do not preserve legacy aliases, configuration discovery, command shapes, or
  release names when they conflict with the approved migration plan.


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
- Parent TODO: [../guiho/TODO.md](../guiho/TODO.md)
- Local TODO: [TODO.md](TODO.md)

### Commands

- Format: `gofmt -w main.go cmd pkg embed devops`
- Tests: `go test ./...`
- Vet: `go vet ./...`
- Build: `go build ./...`
- Compile release asset matrix: `go run devops/build-binaries.go --version <version> --commit <commit> --build-date <RFC3339>`
- Verify release assets: `go run devops/verify-release-assets.go`

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

Before editing release docs or changelogs, inspect mirror.yaml. If agents.write_changelog is false, skip changelog edits. If it is missing or true, changelog edits are allowed when the project has a changelog.

Use [agents].changelog_path as the changelog file path. If it is missing, use CHANGELOG.md in the project root.
<!-- END GUIHO MIRROR -->

<!-- BEGIN MIRROR — DO NOT EDIT THIS SECTION -->
---
name: guiho-i-mirror
description: Mirror agent instruction block.
purpose: Provide the canonical managed project instruction for Mirror versioning.
created: 2026-07-18
owner: mirror-embed-prompts
flags: []
tags: [mirror, instruction, agents]
keywords: [version plan, version apply]
---

# GUIHO Mirror Instruction Block

Run plain `mirror` once in a repository to verify the global Mirror skill and
this bounded instruction block. Repeated runs are idempotent.

Use `mirror version plan <target>` and `mirror version apply <target>` for semantic versioning.
`mirror init` defaults to `v{version}` tags and enables release commits and
pushes; explicit interactive or flag selections remain authoritative.
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
