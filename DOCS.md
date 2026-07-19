---
name: RunX Documentation
purpose: Define the canonical supported RunX CLI contract.
description: Complete reference for command routing, YAML configuration, startup, help, output, agents, upgrades, installation, and release assets.
created: 2026-07-12
flags:
  - canonical
tags:
  - cli
  - reference
keywords:
  - RFC 0034
  - runx.yaml
  - agent namespace
  - release assets
owner: runx
---

# RunX Documentation

## Runtime Contract

RunX uses Bun, strict ESM TypeScript, raw Citty, and TypeBox. Core source is
Bun-only. The isolated npm bootstrap is Node-compatible because it must work
before Bun is installed.

No arguments prints exactly:

```text
Hello Windows - runx v<version>
```

The foreground reads `~/.guiho/runx/cache.json` and never waits for network
work. A detached worker validates GitHub release data and atomically refreshes
that cache. When a decoded cache announces an update, RunX prints:

```text
New version available. Run this command to upgrade: runx upgrade
```

## Configuration

RunX loads YAML only. Resolution is:

1. `--config <path>`;
2. `<effective-cwd>/runx.yaml`;
3. `~/.guiho/runx/runx.yaml`.

Whenever loaded, the absolute path is reported as:

```text
configuration file loaded: <absolute-path>
```

The complete manifest and nested records are TypeBox-decoded. Invalid or absent
configuration exits `3`.

## Commands And Help

The public command tree is the catalog shown in [README.md](README.md). Citty is
the only parser and router. The only short flags are `-h` and root `-v`.

Every root, group, and leaf supports:

- standard help with usage, description, positionals, flags, and examples;
- `--help-tree`, using Unicode box-drawing branches;
- `--help-tree-depth <positive-integer>`;
- redirect-safe Markdown through `--help-docs`.

Only `runx run <selector>` executes catalog code. Listing, describing, checking,
help, initialization, agent operations, upgrade inspection, and dry runs do not.

## Output And Exit Codes

Text results use stdout and diagnostics use stderr. JSON mode emits one valid
JSON document on stdout; configuration reports and diagnostics remain on
stderr.

| Code | Meaning |
| ---: | --- |
| `0` | Success |
| `1` | Unexpected operational failure |
| `2` | Usage or structured flag validation failure |
| `3` | Configuration resolution or decoding failure |
| `4` | Release or network failure |
| `5` | Installation, upgrade, or filesystem mutation failure |
| `130` | Interruption |

An executed catalog command preserves its exact delegated exit code.

## Agent Integration

Skill install, update, and uninstall default to global scope and target both:

```text
~/.agents/skills/guiho-s-runx
~/.claude/skills/guiho-s-runx
```

Use `--local` for the corresponding project-local paths.

Ordinary RunX commands schedule a hidden detached worker that compares the
bundled skill with both global copies and rewrites only missing or stale
targets. The worker finds the nearest existing `AGENTS.md` from effective
`--cwd`; when none exists in an ancestor, it creates `AGENTS.md` at effective
cwd. Reconciliation is atomic, migrates legacy RunX markers, preserves content
outside managed markers, and is silent when no change is required. Worker
failures never change foreground output or exit status.

Instruction actions manage both `AGENTS.md` and `CLAUDE.md` when both exist,
the existing one when only one exists, and create `AGENTS.md` when neither
exists. The exact idempotent boundaries are:

```text
<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->
<!-- END RUNX -->
```

The canonical prompt is `guiho-i-runx`. `agent prompt list --names` prints raw
names; `agent prompt show <id>` prints only the raw prompt.

## Upgrade And Installation

`runx upgrade` accepts `--version`, `--arch`, `--variant`, `--dry-run`, and
`--format`. x64 defaults to `baseline`. `upgrade list` supports positive
`--page` and `--per-page` plus `--pre-releases`.

Replacement is transactional: download, native-format validation, backup,
replacement, version verification, rollback on failure, agent-skill refresh,
instruction reconciliation, then cleanup.

Both direct installers show target metadata and download progress, configure
PATH, install both global skill copies, reconcile project instructions, and
verify the final version.

## Npm Bootstrap

`scripts/runx-bin.mjs` detects platform and architecture, chooses the exact
versioned native asset, caches it under `~/.guiho/runx/npm/<version>/`, applies
Unix execute permissions, delegates args/stdin/stdout/stderr/environment, and
preserves the native exit code. It contains no RunX domain logic.

## Release Assets

Every release contains exactly fourteen assets: twelve native binaries and two
agent assets.

```text
runx-linux-arm64
runx-linux-x64
runx-linux-x64-baseline
runx-linux-x64-modern
runx-darwin-arm64
runx-darwin-x64
runx-darwin-x64-baseline
runx-darwin-x64-modern
runx-windows-arm64.exe
runx-windows-x64.exe
runx-windows-x64-baseline.exe
runx-windows-x64-modern.exe
guiho-s-runx.md
guiho-i-runx.md
```

The Markdown suffix is part of the public GitHub Release filename contract.
Installers retain the standard installed skill filename `SKILL.md` and use the
instruction asset contents when reconciling managed instruction blocks. Before
either write, installers reject empty, executable, binary, invalid UTF-8, or
misidentified Markdown payloads.

Release descriptions contain only the matching `## <version> - <date>`
changelog section. The publish workflow fails closed when that exact section is
missing and refreshes the notes when rerunning an existing release.
