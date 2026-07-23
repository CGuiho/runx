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

No arguments print a deterministic bordered welcome containing RunX, its
purpose, GUIHO identity, platform, architecture, version, and help guidance:

```text
╔════════════════════════════════════════════════════════╗
║  RUNX                                                  ║
║  Documented command catalog                            ║
║  GUIHO · Cristóvão GUIHO                               ║
╚════════════════════════════════════════════════════════╝
```

The platform label matches the runtime operating system.

The foreground reads `~/.guiho/runx/cache.json` and never waits for network
work. A valid cache remains fresh for four hours. When another check is needed,
RunX atomically acquires an ownership-safe lease before it detaches exactly one
worker for the global cache directory. The worker validates GitHub release data,
atomically refreshes the cache, stops remote work after 15 seconds, and releases
its lease on every outcome. A 30-second stale lease can be reclaimed; malformed
or missing lease metadata uses the same delayed recovery. When a decoded cache
announces an update and its validated latest SemVer is newer than the running
version, RunX appends:

```text
  ⚠ New version available: v<version>
    Run `runx upgrade` to update.
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

### Manifest V2

Every catalog uses manifest major version 2 with required `namespace`,
`scripts.directory`, optional `parent`, and recursive `commands`. `project`, a
top-level `groups` map, and `group` on command leaves are invalid legacy keys.

Each `commands` entry is either an executable command (`uid`, `id`, summaries,
and command fields) or a group (`group`, `summary`) with exactly one of nested
`commands` or `runx`. Sibling command IDs, group names, and mounted namespaces
share one collision domain. The catalog namespace cannot collide with a
top-level entry, and composed UIDs/selectors are globally unique.

`runx` and `parent` accept relative paths or full HTTPS GitHub blob/raw URLs.
The group name aliases the mounted child namespace. A child must point back to
the exact parent identity within the same source kind. A local working copy may
mount a foreign child only after RunX loads and validates the declared upstream
foreign parent and confirms it mounts that child. Local children execute relative to their own
catalog directory; foreign children are marked `foreign` and execute relative
to the local mount root. GitHub fetches have a ten-second timeout, a one-MiB
limit, no persistent cache, cycle detection, and a maximum depth of 32.
Foreign-relative references cannot escape their GitHub owner/repository/ref;
use a full URL for an intentional cross-root mount.

Canonical selectors join nested group and mount names, for example
`worker-alias/build/compile`. Unique UIDs remain preferred for automation.
Loading a child directly validates its declared parent but does not implicitly
replace the child with the parent. RunX follows only explicit graph edges and
does not search ancestor directories.

UIDs, canonical selectors, and usable unique-ID shorthands share one collision
domain. Same-command aliases may coincide, but different leaves cannot shadow
one another. Scripts directories and every effective command cwd are validated
during graph composition, before check/list can report a catalog as valid.

`runx init` derives the namespace from the directory that will contain the
target configuration (including nested `--config` paths), lowercases it,
collapses invalid runs to hyphens, trims edges, prefixes leading digits with
`n-`, and falls back to `runx`. It validates the complete rendered v2 document
before an atomic write.

## Commands And Help

The public command tree is the catalog shown in [README.md](README.md). Citty is
the only parser and router. The only short flags are `-h` and root `-v`.

Every root, group, and leaf supports:

- standard help with usage, description, positionals, flags, and examples;
- `--help-tree`, using Unicode box-drawing branches;
- `--help-tree-depth <positive-integer>`;
- redirect-safe Markdown through `--help-docs`.

Only `runx run <selector>` executes catalog code. RunX options belong before the
selector; every token after the selector belongs to the child command. One
immediate `--` after the selector acts as a delimiter and is removed. Listing,
describing, checking, help, initialization, agent operations, upgrade
inspection, and dry runs do not execute catalog code.

Shell adapters keep forwarded values as data. Bash and sh use positional
parameters, PowerShell uses JSON-backed array splatting, and cmd uses a
short-lived delayed-expansion-disabled wrapper with environment-backed values.
Text dry runs show an indexed argument list and JSON dry runs expose a separate
`arguments` array.

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
`--format`. x64 defaults to `baseline`. `upgrade list` returns the complete
published catalog, including labeled prereleases, by default. Positive
`--page` and `--per-page` values request an explicit view; `--pre-releases` is
accepted but is unnecessary because prereleases are never hidden.

Replacement is transactional: download, native-format validation, backup,
replacement, version verification, rollback on failure, agent-skill refresh,
instruction reconciliation, then cleanup.

Both direct installers show target metadata and download progress, configure
PATH, install both global skill copies, reconcile project instructions, and
verify the final version.

The canonical POSIX bootstrap is:

```sh
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash
```

The PowerShell installer reads and writes instruction files as strict UTF-8
without a byte-order mark. Reconciliation preserves the file's newline style,
removes canonical, legacy, damaged, or duplicate RunX managed blocks, and then
writes exactly one canonical block. Its verification probe disables both
background workers so installation cannot race a second instruction write.

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
