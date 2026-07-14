---
name: RunX Interactive Init Manifest
purpose: Preserve the accepted manifest structure and initialization behavior for runx init.
description: Defines the Semantic Versioning manifest field, mandatory public group, configurable scripts directory, and empty interactive initialization contract.
created: 2026-07-14
flags:
  - accepted
  - decision
tags:
  - decisions
  - cli
  - manifest
keywords:
  - runx init
  - runx.yaml
  - semantic versioning
  - public group
  - scripts directory
owner: runx-decisions
---

# RunX Interactive Init Manifest

## Status

Accepted by the project owner on 2026-07-14.

## Context

RunX needs an interactive `runx init` command so a project can establish its
single `runx.yaml` command catalog before any commands exist. The repository
already defines groups and documented commands as the core catalog model. The
initializer must preserve that model without inventing framework-specific
commands, creating a second configuration file, or forcing a placeholder
command into the catalog.

Projects may later store executable Bash, PowerShell, or other script files in
a manifest-owned scripts directory. The manifest must identify that directory
so agents and developers use one predictable location when adding scripts.

## Decision

### Initialized manifest

`runx init` creates this shape, using the current directory name as the default
project name:

```yaml
version: "1.0.0"

project:
  name: my-project

scripts:
  directory: scripts

groups:
  public:
    summary: Default public project commands.

commands: []
```

The initializer creates only `runx.yaml`. It does not create an empty `scripts`
directory because the directory is created when the first script is written.

### Manifest version

- `version` is a valid Semantic Versioning 2.0.0 string with a
  `MAJOR.MINOR.PATCH` core and optional prerelease or build metadata.
- The first supported manifest contract is `"1.0.0"`.
- RunX accepts manifest versions whose major version is `1` and rejects
  unsupported major versions.
- Numeric `version: 1` manifests are invalid. RunX is still under development
  and has no deployed numeric-version compatibility requirement.
- Future backward-incompatible manifest changes require a major version change.

### Scripts directory

- Every manifest contains `scripts.directory`; it configures the project script
  location.
- The initializer's default value is `scripts` relative to the directory
  containing `runx.yaml`.
- The path must remain inside the manifest directory and cannot escape through
  an absolute path or parent traversal.
- Existing and future command strings may invoke any suitable script type from
  this directory; RunX remains language-agnostic.
- The manifest command entry is the script descriptor. No sidecar RunX config
  or per-script descriptor file is introduced.

### Groups and commands

- `groups` is a first-class catalog concept.
- Every manifest contains a `public` group, analogous to PostgreSQL's default
  `public` schema.
- Additional groups may organize commands by domain or responsibility.
- Every command explicitly names an existing group. Missing and unknown group
  references are invalid.
- Authoring workflows place a new command in `public` unless the user requests
  or selects another group, but the stored command still contains
  `group: public` explicitly.
- An initialized catalog may contain no commands. Once present, every command
  retains the required `uid`, `id`, `group`, `summary`, `description`, and
  executable `command` fields.

### Interactive behavior

- `runx init` is a first-class Citty command that never requires an existing
  manifest.
- The interactive flow collects or confirms the project name and scripts
  directory, renders a YAML preview, and requests confirmation before writing.
- The interface should use a polished terminal presentation with clear
  progress, validation, cancellation, and success states.
- No framework-specific templates or automatic package-script imports are
  offered.
- An existing `runx.yaml` is never overwritten without explicit confirmation.
- Cancellation leaves no partial file, and non-interactive invocation fails
  clearly instead of hanging.

## Alternatives Considered

### Require a first command during initialization

Rejected. Initialization establishes the catalog so commands can be added later
from explicit developer requests. A fabricated starter command would make the
manifest less truthful.

### Offer framework-specific templates or scan package scripts

Rejected. RunX is language-agnostic, and the project owner explicitly excluded
framework-specific starter templates from the initializer.

### Keep numeric `version: 1` compatibility

Rejected. RunX has not been used outside development, so there is no legacy
manifest population to migrate. Supporting two version shapes would add
unnecessary schema and documentation complexity.

### Store the scripts directory in a separate config file

Rejected. `runx.yaml` remains the single source of truth for the command
catalog and its script location.

### Create the scripts directory eagerly

Rejected for initialization. An empty directory is not useful to the catalog
and is not preserved by Git. The directory should be created with the first
real script.

## Consequences

- The manifest TypeBox schema must accept only supported Semantic Versioning
  strings, require `scripts.directory`, and allow an empty command array.
- Semantic validation must require the `public` group and continue requiring
  every command to reference a declared group.
- Script-directory validation must enforce a relative path within the manifest
  root.
- Existing fixtures and documentation using numeric `version: 1` must move to
  `version: "1.0.0"` and include `public`.
- Listing and checking an empty initialized manifest must succeed without
  attempting command execution.
- The initializer needs an isolated prompt boundary so interactive behavior can
  be tested without depending on a real terminal.
- The bundled agent skill update is intentionally deferred to the process that
  delivers the complete agent-integration request.

## Reversal Or Revisit Conditions

Revisit this decision only if real-world manifest migration requires a
compatibility policy, the command catalog needs multiple script roots, or a
future major manifest version replaces the group model. Such a change must be
recorded as a new decision and reflected in the manifest's major version.

## Follow-up Work

- Implement the manifest schema and semantic validation changes.
- Implement and test the interactive `runx init` command.
- Update canonical CLI and manifest documentation.
- Update XDocs descriptors for every changed module or companion document.
- Leave agent-skill content changes for the complete agent-integration work.

## References

- [Alpha boundaries](./alpha-boundaries.md)
- [Alpha command-catalog requirements](../requirements/alpha-command-catalog.md)
- [CLI architecture](../architecture/cli-architecture.md)
- [GitHub issue #10](https://github.com/CGuiho/runx/issues/10)
