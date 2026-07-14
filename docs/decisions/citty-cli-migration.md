---
name: RunX Citty CLI Migration
purpose: Define the accepted replacement of RunX's handwritten argument parser and router with Citty.
description: Records the full Citty command-tree migration, compatibility requirements, error behavior, validation scope, and patch-release boundary.
created: 2026-07-14
flags:
  - accepted
  - needs-implementation
tags:
  - decisions
  - cli
  - typescript
keywords:
  - runx
  - citty
  - argument parsing
  - command routing
  - short flags
owner: runx-decisions
---

# RunX Citty CLI Migration

## Context

RunX currently parses arguments and routes commands with repository-owned loops
and conditionals. The parser recognizes only long options beginning with `--`,
so conventional aliases such as `runx -v` and `runx -h` are treated as command
selectors. That incorrect fallback causes global CLI operations to search for a
`runx.yaml` manifest.

Mirror and XDocs provide useful command semantics and help conventions, but
both also maintain handwritten parsers. RunX will instead use Citty as the
single parser, command router, alias registry, and usage generator.

## Decision

- Add `citty` as a runtime dependency and bundle it into native RunX builds.
- Replace the handwritten parser in `source/flags.ts` and the manual routing in
  `source/cli.ts` with a complete Citty command tree.
- Do not introduce Commander, yargs, oclif, Clipanion, or a second parsing
  abstraction around Citty.
- Keep Bun as the development, test, build, and native-compilation toolchain.
- Do not add a Node.js runtime dependency or an npm `postinstall` hook.
- Release the completed migration as the next Mirror-managed patch after all
  validation passes.

## Command Tree

The Citty tree must preserve the current public commands:

```text
runx
|- list
|- describe <selector>
|- run <selector>
|- r <selector>
|- check
|- agents
|  |- install <local|global>
|  `- instructions
|- upgrade
|  |- check
|  `- list
`- uninstall
```

`r` remains an alias of `run`. The human shorthand `runx <selector>` remains
supported through the root command's positional selector compatibility path.
Known Citty subcommands always take precedence over selector fallback.

## Arguments and Help

- Citty owns `-h`/`--help` and `-v`/`--version`.
- Version and help must return without reading or discovering `runx.yaml`.
- Preserve `--cwd`, `--file`, `--format`, `--verbose`, `--dry-run`, `--yes`,
  `--help-tree`, and `--help-docs` where they apply.
- Preserve `--tool agents|claude|all` for agent-skill installation.
- Preserve JSON output contracts for successful commands and self-management
  results.
- Use Citty metadata and argument definitions as the source of ordinary command
  usage text. Keep RunX's extended help tree and documentation guidance as
  explicit custom outputs rather than duplicating Citty's ordinary help.

## Implementation Boundaries

- `source/cli.ts` owns the root Citty command and CLI entry function.
- Command definitions may be grouped into focused modules for catalog,
  agent, and self-management commands when that keeps `source/cli.ts` small.
- Existing manifest parsing, selector resolution, execution, agent
  installation, upgrade, and uninstall functions remain the business-logic
  implementations. The migration changes their CLI adapters, not their domain
  behavior.
- Remove `source/flags.ts` after all consumers migrate to Citty.
- Keep entrypoints thin: they register embedded native resources when needed
  and invoke the Citty-backed CLI.

## Error Behavior

- Unknown options fail as CLI usage errors and never trigger manifest discovery.
- Unknown explicit subcommands show relevant Citty usage.
- A root positional that is not a known command is treated as the intentional
  selector shorthand and may discover a manifest.
- Missing command arguments show the specific command's usage.
- Domain failures continue to use `RunXError` messages and non-zero exit codes.
- Verbose mode may expose additional diagnostics, but normal errors remain
  concise and actionable.

## Testing

Automated coverage must include:

- `runx -v` and `runx --version` outside a configured project.
- `runx -h` and `runx --help` outside a configured project.
- Help and missing-argument behavior for each command group.
- `runx run <selector>`, `runx r <selector>`, and `runx <selector>` compatibility.
- Unknown flags and unknown explicit commands.
- Existing global flags, JSON output, confirmation gates, and dry runs.
- Agent install/instructions routing.
- Upgrade check/list/apply and uninstall routing.
- Typecheck, unit tests, library build, local native build, and the twelve-target
  native binary matrix.

## Documentation and Release

Update public command documentation, the bundled `guiho-s-runx` skill, source
descriptors, and release notes to match the Citty command tree. Validate the
affected XDocs scopes before delivery. Use Mirror to plan and apply the next
patch; never hand-edit the package version or create the release tag manually.

## Acceptance Criteria

- No RunX-owned argument-token loop or manual top-level command router remains.
- `-v`, `--version`, `-h`, and `--help` work without a manifest.
- Current commands and selector shorthand remain compatible.
- Invalid CLI input produces usage errors rather than unrelated manifest errors.
- The npm package and native executables expose the same Citty-backed behavior.
- The implementation adds no `postinstall` hook and no Node.js runtime
  requirement.
