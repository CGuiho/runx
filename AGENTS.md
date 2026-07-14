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

## Repository Notes

- RunX is the open-source `@guiho/runx` Bun/TypeScript CLI for a documented,
  language-agnostic `runx.yaml` command catalog.
- The executable entrypoint is `source/guiho-runx-bin.ts`; native builds use
  `source/guiho-runx-native-bin.ts`. Do not add a Node.js runtime dependency.
- Use Bun for installs, tests, typechecking, builds, and executable compilation.
- `runx` without arguments shows the home page and usage. `runx list` lists a
  manifest; `runx run` and `runx r` execute a selected command.
- Citty owns argument parsing, aliases, command routing, and ordinary usage.
  Keep `-h`/`--help` and `-v`/`--version` manifest-free, and do not add a
  second handwritten token parser or manual execution router.
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

- Typecheck: `bun run typecheck`
- Tests: `bun test`
- Build: `bun run build`
- Compile local executable: `bun run binary`
- Compile release asset matrix: `bun run binaries`

## Semantic Project Versioning -- GUIHO Mirror

Use the `guiho-s-mirror` skill whenever versioning, tags, changelogs, or release
configuration changes. Read `mirror.config.toml`, run `mirror version plan`
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
