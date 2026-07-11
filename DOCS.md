# RunX Documentation

## Purpose

RunX makes project operations visible and documented. A `runx.yaml` file is a
single local command catalog; RunX discovers the nearest catalog, validates it,
lists its commands, and runs exactly one selected command.

## Commands

| Command | Purpose |
| --- | --- |
| `runx` | Show the home page and usage without loading a manifest. |
| `runx list` | List all manifest commands. Add `--format json` for structured output. |
| `runx describe <selector>` | Show one command's description and operational metadata. |
| `runx check` | Validate discovery and the manifest without execution. |
| `runx run <selector>` | Execute a selected command. |
| `runx r <selector>` | Short alias for `runx run`. |
| `runx <selector>` | Human shorthand for `runx run <selector>` when no built-in name conflicts. |
| `runx agents install <local|global>` | Install `guiho-s-runx`; add `--tool agents|claude|all` when needed. |
| `runx agents instructions` | Add or refresh the managed RunX section in `AGENTS.md`. |
| `runx upgrade [check|list]` | Check releases, list releases, or replace a native installed executable. |
| `runx uninstall [--dry-run]` | Remove a native installed executable. |

Global flags: `--cwd <path>`, `--file <path>`, `--format <text|json>`,
`--verbose`, `--version`, `--help`, `--help-tree`, and `--help-docs`.

## Manifest

RunX searches upward from the current directory for `runx.yaml`. Use `--file`
to select an explicit manifest. RunX does not merge files.

```yaml
version: 1
project:
  name: example
groups:
  development:
    summary: Local development work.
  release:
    summary: Explicit release actions.
commands:
  - uid: example-check
    id: check
    group: development
    summary: Type-check the project.
    description: Runs the project's static type checker without emitting build files.
    command: bun run typecheck
    cwd: .
    shell: auto
    tags: [validation]
  - uid: example-release
    id: release
    group: release
    summary: Build the release artifact.
    description: Builds the release artifact after the project has passed its required checks.
    command: bun run build
    confirm: always
```

Required fields are `uid`, `id`, `group`, `summary`, `description`, and
`command`. `uid` and `id` begin with a lowercase letter and use lowercase
letters, digits, and hyphens. `uid` is globally unique; the group and ID pair
is also unique.

Optional fields are `cwd`, `shell`, `tags`, and `confirm`. `cwd` is relative to
the manifest and cannot escape its directory. Shell values are `auto`, `bash`,
`sh`, `powershell`, and `cmd`. Confirmation values are `never` and `always`.

## Selectors and Safety

RunX resolves selectors in this order: exact UID, exact `group/id`, one-based
manifest index, then an unambiguous bare ID. Use UIDs in automation. An index is
only stable until commands are reordered, and an ambiguous bare ID fails rather
than selecting a command.

`list`, `describe`, `check`, and `run --dry-run` never execute a configured
command. A real run streams stdout/stderr and uses the child process exit code.
`confirm: always` requires the user to add `--yes`.

Manifests are executable local code. Do not place secrets in a manifest and do
not assume a group name determines safety.

## Agent Skill

The bundled `guiho-s-runx` skill directs agents to validate and list a catalog,
select a UID, inspect unfamiliar commands, use dry runs, and obtain explicit
authorization before `--yes` for high-impact commands. Install it with:

```text
runx agents install local
```

## Native Install, Upgrade, and Uninstall

`devops/install.ps1` installs Windows assets and `devops/install.sh` installs
macOS/Linux assets from GitHub Releases into the user's local bin directory.
`runx upgrade` uses the matching release asset and replaces a native executable;
on Windows replacement is scheduled after RunX exits. `runx uninstall` removes
the same native executable and supports `--dry-run`.

The alpha CI validates code and compiles a native executable but intentionally
does not publish to npm or GitHub Releases.
