#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# GUIHO RunX

RunX is a native Go/Cobra CLI for documented, language-agnostic command
catalogs. A project owns one explicit `runx.yaml`; RunX validates and describes
catalog commands without executing them, `runx reveal` prints one exact command
for copy-and-paste, and only `runx run` starts a configured command.

## Install

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
```

macOS and Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | sh -s --
```

Prompt for the AI (load the `runx-install` prompt):

```text
Load the runx-install prompt and follow it to install the GUIHO RunX CLI.
```

Verify the installation:

```bash
runx --version
```

Both installers accept `--version <semver>` / `-Version <semver>` or
`--channel <name>` / `-Channel <name>` (mutually exclusive; default latest
stable). They stage downloads under `$HOME/.guiho/.temp/`, verify every artifact
against `checksums.txt` and the release `artifacts.json` manifest, run the
payload's hidden self-test before activation, install immutable payloads into
`$HOME/.guiho/runx/versions/<version>/`, activate them through the stable
launcher at `$HOME/.guiho/bin/runx`, preserve configuration and data on
reinstall, and roll back completely on failure.

## Start

```bash
runx init
runx check --format json
runx list --format json
runx describe <uid-or-selector-or-index>
runx reveal <uid-or-selector-or-index>
runx run --dry-run <uid-or-selector-or-index>
runx run --yes <uid-or-selector-or-index> -- <child arguments...>
```

`runx reveal` accepts one selector plus the catalog location flags `--cwd`,
`--config`, and `--verbose`. It resolves the same UID, canonical selector,
unique shorthand, and numeric index as `runx run`, then writes only the exact
stored command and one trailing newline to stdout. It never executes or asks
for confirmation; diagnostics remain on stderr.

RunX options precede the selector. Every token after the selector belongs to
the child and is forwarded without reinterpretation. The numeric `IDX` printed
by `runx list` is convenient for interactive use; prefer stable UIDs for
automation because indexes belong to the current resolved listing.

Selectors resolve in a deterministic order: exact global UID, canonical
group-scoped selector, then an unqualified ID shorthand only when that ID has a
single owner. A UID may equal another command's ID; the UID still wins. Duplicate
unqualified IDs remain ambiguous and fail instead of selecting an arbitrary
command.

Commands marked `confirm: always` ask `Are you sure? [y/N]` in an interactive
terminal and show the exact command that skips the prompt, such as
`runx run --yes cli-compile-host`. Enter or any answer other than `y` or `yes`
declines. Noninteractive and JSON invocations never prompt; they fail closed
with the same exact retry command.

## Manifest v2

```yaml
version: "2.0.0"
namespace: "example"
scripts:
  directory: ".scripts"
commands:
  - uid: "test-command"
    id: "test"
    summary: "Run tests."
    description: "Run the complete project test suite."
    command: "go test ./..."
    confirm: "never"
```

`runx init` writes `.scripts` as the default scripts directory. An explicitly
configured `scripts.directory` value is preserved unchanged.

Manifests are strictly decoded: unknown fields, invalid identifiers, unsafe
paths, duplicate UIDs or canonical selectors, non-reciprocal child catalogs,
and unsupported manifest versions fail closed. Configuration precedence is:

1. `--config <path>`;
2. effective-cwd `runx.yaml`;
3. `~/.guiho/runx/runx.yaml`.

RunX never searches parent directories.

## Developer Context

Every command scope supports `-h`/`--help`, `--help-tree`,
`--help-tree-depth <positive-integer>`, `--help-tree-global-flags`, and
`--help-docs`. Root additionally supports `-v`/`--version`. By default the tree
shows global flags once at the root; `--help-tree-global-flags` repeats them
under every descendant. The tree and Markdown are generated from the live
Cobra commands.

## Agent And Upgrade Commands

```text
runx agent skill install|uninstall|update|list|show
runx agent instruction apply|remove|update|show
runx agent prompt list|show
runx upgrade
runx upgrade check
runx upgrade list
runx uninstall [--preserve-config] [--preserve-data] (--dry-run | --yes)
```

Every terminal upgrade outcome prints a reinstallation recovery command pinned
to the resolved target version. The uninstall command shares its contract with
the remote uninstallers: it prints a REMOVE/PRESERVE plan, fails closed for
noninteractive use without `--yes`, and never removes shared `$HOME/.guiho/`
infrastructure or another CLI's files.

Bare `runx` first installs its embedded skill in both global tool locations and
idempotently reconciles a bounded RunX instruction block in the repository root.
Both `AGENTS.md` and `CLAUDE.md` are updated when both exist; otherwise the one
that exists is used, or `AGENTS.md` is created. Existing content and line endings
are preserved, malformed markers fail safely, and no catalog command or network
request runs during bootstrap. Help, version, agent-management, uninstall, and
non-repository paths do not perform repository bootstrap. Other foreground
startup reads only the local cache and starts bounded detached workers where
appropriate.
Self-upgrades verify published checksums and preserve the embedded build target,
including ARMv6 versus ARMv7.

## Uninstall

Uninstallation removes every RunX-owned artifact by default: the stable
launcher, all versioned payloads, `$HOME/.guiho/runx/`, global configuration,
persistent data, managed agent-skill copies, the managed instruction block in
the project's `AGENTS.md`, and the project's `runx.yaml`. The shared
`$HOME/.guiho/` infrastructure, its `bin/` and `.temp/` directories, and the
user PATH entry are preserved.

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/uninstall.ps1 | iex
```

macOS and Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/uninstall.sh | sh -s --
```

Prompt for the AI (load the `runx-uninstall` prompt):

```text
Load the runx-uninstall prompt and follow it to uninstall the GUIHO RunX CLI.
```

Destructive default (removes everything):

```bash
runx uninstall --yes
```

Dry run (shows REMOVE/PRESERVE plan without changing anything):

```bash
runx uninstall --dry-run
```

Preserve configuration and persistent data:

```bash
runx uninstall --preserve-config --preserve-data --yes
```

Noninteractive invocations without `--yes` fail without changing anything.
