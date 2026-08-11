#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# GUIHO RunX

RunX is a native Go/Cobra CLI for documented, language-agnostic command
catalogs. A project owns one explicit `runx.yaml`; RunX validates and describes
catalog commands without executing them, `runx reveal` prints one exact command
for copy-and-paste, and only `runx run` starts a configured command.

## Install

Linux and macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
```

Both installers select the canonical target, download `checksums.txt`, verify
SHA-256 before replacement, install the bundled skill into both supported agent
locations, and verify `runx --version`. The Windows installer adds its directory
to the persistent user `Path`, the current PowerShell process, and Git Bash's
`~/.bashrc` without duplicating entries. An existing Git Bash session can load
the change immediately with `source ~/.bashrc`.

### Migrate From RunX 0.8

RunX 0.8 uses the retired Bun release contract and cannot discover current
native releases. If `runx upgrade` reports that 0.8 is already up to date, use
the unpinned installer instead of the pinned 0.8 recovery command.

Linux and macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash
hash -r
runx --version
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex
runx --version
```

Restart an existing shell if it still resolves the old executable. In Git Bash,
`source ~/.bashrc` loads the installer-managed path without restarting.

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
`--help-tree-depth <positive-integer>`, and `--help-docs`. Root additionally
supports `-v`/`--version`. The tree and Markdown are generated from the live
Cobra commands.

## Agent And Upgrade Commands

```text
runx agent skill install|uninstall|update|list|show
runx agent instruction apply|remove|update|show
runx agent prompt list|show
runx upgrade
runx upgrade check
runx upgrade list
runx uninstall --dry-run
```

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

## Build

```bash
go test ./...
go vet ./...
go build ./...
go run devops/build-binaries.go --version 0.8.0 --commit <commit> --build-date 2026-07-26T00:00:00Z
go run devops/verify-release-assets.go
```

The release contract is exactly 11 artifacts: eight pure-Go executables
(Linux AMD64/ARM64/ARMv7/ARMv6, Darwin AMD64/ARM64, Windows AMD64/ARM64),
`guiho-s-runx.zip`, `guiho-i-runx.md`, and `checksums.txt`. AMD64 V2/V3/V4
variants are not part of the contract.

See [DOCS.md](DOCS.md) for the complete behavior and safety reference.
