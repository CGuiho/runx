---
name: RunX RFC 0034 CLI Compliance Migration Plan
purpose: Provide an executable, breaking-change migration plan that brings the RunX CLI into complete compliance with GUIHO RFC 0034.
description: Defines the final command catalog, Bun-only runtime migration, TypeBox validation, YAML configuration, startup lifecycle, Developer Context help, agent integration, upgrades, installers, npm bootstrap, release assets, documentation, and verification.
created: 2026-07-18
flags:
  - approved
  - breaking-change
  - implementation-ready
  - implemented
tags:
  - planning
  - cli
  - migration
  - rfc-0034
keywords:
  - runx
  - guiho-s-0034-cli-engineer
  - Bun
  - Citty
  - TypeBox
  - runx.yaml
  - agent namespace
  - fourteen release assets
owner: runx-plans
---

# RunX RFC 0034 CLI Compliance Migration Plan

## Outcome

Migrate `@guiho/runx` from its current partial RFC implementation to the complete
GUIHO RFC 0034 CLI contract. The migration is intentionally breaking. RunX is
pre-1.0, and the developer has approved removing aliases, compatibility
shorthands, legacy configuration discovery, legacy asset names, legacy agent
commands, and any other interface that conflicts with RFC 0034.

The completed CLI must use Bun, strict ESM TypeScript, raw Citty, and TypeBox;
contain no prohibited Node built-ins in core CLI source; expose the exact
Developer Context and agent command contracts; use the standardized startup and
storage lifecycle; and publish exactly fourteen release assets.

## Authority And Required Execution Roles

Every implementation session must load and follow:

- Agent: `guiho-a-0001-swe`, the GUIHO Software Engineer/SWE coordinator at
  `C:\GUIHO\superiority\agents\guiho-a-0001-swe.AGENTS.md`.
- Mandatory specialist skill: `guiho-s-0034-cli-engineer`.
- Implementation controller: `guiho-s-0023-plan-executor`.
- Runtime and language skills: `guiho-s-0015-bun` and
  `guiho-s-0019-typescript`.
- Runtime-contract skill: `guiho-s-0011-typebox`.
- Documentation skills: `guiho-s-0016-writing-docs`,
  `guiho-s-0017-todo`, and `guiho-s-xdocs`.
- Release skills: `guiho-s-mirror` and, when GitHub release automation changes,
  `guiho-s-0020-cloud-computing`.
- Review and validation skills: `guiho-s-0029-implementation-reviewer` and
  `guiho-s-0030-validation-reporter`.

`guiho-s-0034-cli-engineer` is a skill, not an agent. The SWE agent coordinates
the work and the CLI engineer skill supplies the mandatory CLI contract.

## Approved Breaking-Change Boundary

The implementation must prefer a clean RFC 0034 design over compatibility.
Specifically:

- Remove the `runx r` command alias.
- Remove the root `runx <selector>` compatibility shorthand.
- Replace `--file` with RFC configuration selection through `--config`.
- Stop searching parent directories for `runx.yaml`.
- Replace `runx agents ...` with the singular RFC `runx agent ...` namespace.
- Remove `--tool agents|claude|all` and positional `local|global` installation
  syntax in favor of global-by-default operations and `--local`.
- Replace `macos` release asset names with `darwin`.
- Replace the Bun-dependent npm launcher with a Node-compatible native-binary
  bootstrap.
- Do not add deprecated aliases, compatibility adapters, dual configuration
  loaders, dual release names, or migration warnings that keep the old contract
  alive.

The RunX manifest remains YAML and keeps its command-catalog domain purpose.
Its schema may change where necessary to satisfy strict TypeBox decoding, but
the plan does not add workflow graphs, remote execution, or secret management.

## Final Public Command Catalog

The final Citty tree must be:

```text
runx
├── list
├── describe <selector>
├── run <selector>
├── check
├── init
├── agent
│   ├── skill
│   │   ├── install
│   │   ├── uninstall
│   │   ├── update
│   │   ├── list
│   │   └── show <id>
│   ├── instruction
│   │   ├── apply
│   │   ├── remove
│   │   ├── update
│   │   └── show
│   └── prompt
│       ├── list
│       └── show <id>
├── upgrade
│   ├── check
│   └── list
└── uninstall
```

Rules:

- Citty is the only parser and router.
- Only `-h` for `--help` and root `-v` for `--version` exist.
- Every command and subcommand supports `--help`, `-h`, `--help-tree`,
  `--help-tree-depth <positive-integer>`, and `--help-docs`.
- Domain flags remain descriptive long kebab-case flags, including
  `--dry-run`, `--yes`, `--format`, `--cwd`, and `--config` where applicable.
- `runx run` is the only execution path. Inspection, help, configuration,
  agent, upgrade, and dry-run paths must not execute manifest commands.

## Skill Routing By Unit

| Units | Required skills |
| --- | --- |
| RX-01 | `guiho-a-0001-swe`, `guiho-s-0023-plan-executor`, `guiho-s-0034-cli-engineer` |
| RX-02-RX-03 | Add `guiho-s-0015-bun`, `guiho-s-0019-typescript`, and `guiho-s-xdocs` |
| RX-04-RX-05 | Add `guiho-s-0011-typebox`; keep Bun, TypeScript, CLI engineer, and xdocs loaded |
| RX-06-RX-10 | Keep CLI engineer, Bun, TypeScript, TypeBox, and xdocs loaded |
| RX-11-RX-14 | Add `guiho-s-0020-cloud-computing` for release hosting/CI; use `guiho-s-mirror` only for release configuration inspection, not version application |
| RX-15 | Add `guiho-s-0016-writing-docs` and `guiho-s-0017-todo` |
| RX-16 | Add `guiho-s-0029-implementation-reviewer` and `guiho-s-0030-validation-reporter` |

Every row also inherits the SWE agent and plan executor. Reload a skill in a new
execution session rather than assuming it remains active from an earlier unit.

## Execution Sequence

### Unit RX-01 - Freeze The Baseline And Create The Migration Ledger

- Goal: capture the current public behavior, test baseline, file inventory, and
  known RFC violations before source changes begin.
- Owner: RunX repository root.
- Dependencies: none.
- Read first:
  - `AGENTS.md`
  - `TODO.md`
  - this plan and its linked task specification
  - `package.json`, `bun.lock`, `tsconfig.json`
  - `source/cli.ts`, `source/help.ts`, `source/agents.ts`,
    `source/manifest.ts`, `source/self-management.ts`
  - entrypoints, launchers, installers, binary builder, workflows, bundled
    skill, descriptors, README, and DOCS.
- Actions:
  1. Confirm `main`, record `git status --short`, and preserve unrelated work.
  2. Run the current typecheck and tests.
  3. Record current command/output snapshots for root, every command, agent
     operations, upgrade operations, and invalid input.
  4. Search `source/` for prohibited Node imports and save the exact result in
     implementation notes.
  5. Inventory current release filenames and current GitHub upload globs.
- Acceptance:
  - Baseline successes and existing failures are distinguishable from migration
    regressions.
  - No source or generated output changes occur in this unit.
- Stop conditions:
  - Stop if the worktree contains overlapping unexplained changes.

### Unit RX-02 - Normalize The Project Layout And Dependencies

- Goal: establish the RFC module boundaries and mandatory dependencies without
  changing domain behavior.
- Dependencies: RX-01.
- Expected files:
  - `package.json`, `bun.lock`, `tsconfig.json`
  - `source/commands/agent.ts`
  - `source/commands/catalog.ts`
  - `source/commands/upgrade.ts`
  - `source/configuration.ts`
  - `source/storage.ts`
  - `source/update-cache.ts`
  - `source/output.ts`
  - existing source modules moved only when ownership becomes clearer
  - affected `*.xdocs.md` descriptors.
- Actions:
  1. Keep `citty` and `@sinclair/typebox` as runtime dependencies.
  2. Keep strict ESM TypeScript and Bun types.
  3. Separate platform behavior from RunX catalog behavior.
  4. Keep `source/cli.ts` responsible for assembling and running one command
     tree, not domain implementation.
  5. Keep entrypoints thin and make the hidden worker test-callable.
  6. Update source descriptors for every new, moved, or removed module.
- Acceptance:
  - The final ownership of CLI, catalog, configuration, storage, help, agents,
    upgrades, output, and errors is explicit.
  - No duplicated parser or router exists.

### Unit RX-03 - Remove Prohibited Node Built-ins From Core Source

- Goal: make core RunX Bun-only.
- Dependencies: RX-02.
- Expected files:
  - `source/agents.ts`
  - `source/cli.ts`
  - `source/executor.ts`
  - `source/init.ts`
  - `source/manifest.ts`
  - `source/self-management.ts`
  - new Bun-first path/storage modules.
- Actions:
  1. Replace `node:fs` and `node:fs/promises` operations with `Bun.file`,
     `Bun.write`, and approved Bun-native filesystem operations.
  2. Replace `node:path` joins/resolution with URL-based resource paths and a
     narrow platform-aware path utility.
  3. Resolve home from `Bun.env.HOME` or `Bun.env.USERPROFILE`; fail clearly
     when neither exists.
  4. Keep subprocesses on `Bun.spawn`/`Bun.spawnSync`.
  5. Permit Node imports only in the isolated npm bootstrap wrapper.
  6. Add a static test that fails when prohibited imports appear in core
     source.
- Acceptance:
  - Static search finds no `node:fs`, `node:fs/promises`,
    `node:child_process`, `node:path`, or `node:os` in core CLI source.
  - RunX catalog reads and executions still work on Windows, Linux, and Darwin
    path forms.

### Unit RX-04 - Implement The Exact YAML Configuration Contract

- Goal: make `runx.yaml` the schema-backed RFC configuration.
- Dependencies: RX-03.
- Expected files:
  - `source/configuration.ts`
  - `source/manifest.ts`
  - `source/types.ts`
  - `source/init.ts`
  - CLI adapters and tests.
- Resolution order:
  1. `--config <path>`
  2. `<effective-cwd>/runx.yaml`
  3. `~/.guiho/runx/runx.yaml`
- Actions:
  1. Remove `--file`.
  2. Remove upward parent-directory discovery.
  3. Resolve `--cwd` first, then apply the exact configuration precedence.
  4. Parse with `Bun.YAML.parse`.
  5. Decode the complete manifest with TypeBox.
  6. Print `configuration file loaded: <absolute-path>` whenever a manifest is
     loaded.
  7. Make `runx init` create the selected/effective `runx.yaml` and never create
     TOML or JSON configuration.
  8. Reject missing, malformed, or schema-invalid input with exit code `3`.
- Acceptance:
  - Precedence, absolute-path reporting, invalid YAML, invalid schema, and
    absent-config behavior have focused tests.
  - No fallback silently changes an invalid object.

### Unit RX-05 - Expand TypeBox To Every Structured Boundary

- Goal: use TypeBox for all untrusted or structured runtime input, not only the
  manifest.
- Dependencies: RX-04.
- Schemas:
  - RunX manifest and nested command/group records.
  - cached update object.
  - GitHub release and asset API responses.
  - release catalog entries and upgrade selections.
  - enum-like flags such as format, architecture, and variant.
  - positive integers such as help depth, page, and per-page.
  - structured JSON output envelopes where stability matters.
- Actions:
  1. Define schemas next to the owning data boundary.
  2. Derive static types from schemas.
  3. Decode unknown values before access.
  4. Map validation failures to usage/config/remote exit codes.
- Acceptance:
  - No `as` cast substitutes for runtime decoding of external data.
  - Invalid cache and remote payloads cannot reach selection logic.

### Unit RX-06 - Rebuild The Citty Tree To The Final Catalog

- Goal: remove compatibility routing and define the final command tree exactly
  once.
- Dependencies: RX-05.
- Expected files: `source/cli.ts`, `source/commands/*.ts`, CLI tests.
- Actions:
  1. Remove `runx r`.
  2. Remove root selector fallback.
  3. Define every final command/subcommand with `defineCommand`.
  4. Route with Citty `runMain` or `runCommand`.
  5. Use common flag factories without pre-scanning ordinary arguments.
  6. Preserve the child process exit code for an actual `runx run`.
  7. Enforce long kebab-case flags and reject arbitrary short forms.
- Acceptance:
  - Unknown commands are usage errors, never manifest selectors.
  - One Citty tree owns routing, usage metadata, help traversal, and docs
    traversal.

### Unit RX-07 - Implement The Standard Startup Lifecycle

- Goal: make startup deterministic, local-first, and non-blocking.
- Dependencies: RX-05 and RX-06.
- Exact sequence:
  1. Read and TypeBox-decode `~/.guiho/runx/cache.json` synchronously.
  2. If `newVersionAvailable` is true, print exactly
     `New version available. Run this command to upgrade: runx upgrade`.
  3. Resolve, decode, and report configuration only for commands that need it.
  4. Spawn the hidden detached update worker without awaiting network work.
  5. Route the command.
  6. For no arguments, print exactly `Hello Windows - runx v<version>`.
- Actions:
  - Add a recursion-safe hidden worker path.
  - Validate GitHub release responses, compare SemVer, and atomically write the
    cache object.
  - Ignore corrupt foreground cache with verbose diagnostics and allow the
    worker to replace it.
  - Ensure help, version, and help-docs output are not polluted by notices or
    progress when their contract requires clean redirection.
- Acceptance:
  - Tests prove the foreground never calls `fetch`.
  - The worker writes both update-available and up-to-date cache shapes.

### Unit RX-08 - Generate Complete Developer Context Help

- Goal: derive all help forms from the actual Citty tree.
- Dependencies: RX-06.
- Actions:
  1. Delete the handwritten drifting command list in `source/help.ts`.
  2. Make standard help include usage, description, positionals, flags, and
     practical examples.
  3. Implement `--help-tree` at every scope, beginning with `COMMAND TREE`.
  4. Render Unicode `├──`, `└──`, and `│   ` branches with flags nested under
     their owning command.
  5. Implement and TypeBox-validate
     `--help-tree-depth <positive-integer>`.
  6. Implement `--help-docs` as redirect-safe Markdown for the selected scope.
  7. Align description columns deterministically.
- Acceptance:
  - No ASCII `|-` tree remains.
  - Every root, group, and leaf has all help forms and focused snapshot tests.

### Unit RX-09 - Replace Agent Automation With The Full RFC Namespace

- Goal: implement every required agent action against embedded, versioned
  resources.
- Dependencies: RX-03, RX-05, and RX-06.
- Expected files:
  - `source/commands/agent.ts`
  - `source/agents.ts`
  - `source/embedded-resources.ts`
  - `skills/guiho-s-runx/SKILL.md`
  - `prompts/guiho-i-runx.md`
  - tests and descriptors.
- Skill actions:
  - `install`, `update`, and `uninstall` always target both
    `.agents/skills/guiho-s-runx` and `.claude/skills/guiho-s-runx`.
  - Default scope is global; `--local` selects project-local scope.
  - `list [--filter <keyword>]` enumerates bundled skills.
  - `show <id>` prints path, description, and metadata.
- Instruction actions:
  - Resolve existing `AGENTS.md`, existing `CLAUDE.md`, both when both exist,
    and create `AGENTS.md` when neither exists.
  - Use exact markers:
    `<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->` and
    `<!-- END RUNX -->`.
  - `apply` appends or replaces; `remove` removes; `update` replaces stale
    content; `show` prints the raw template.
  - Every action is idempotent.
- Prompt actions:
  - Create `prompts/guiho-i-runx.md` as the canonical bundled RunX instruction
    prompt asset.
  - `list` prints names and descriptions.
  - `list --names` prints names only.
  - `show <id>` prints only the raw prompt body.
- Removal:
  - Delete `runx agents install`, `runx agents instructions`, `--tool`, and the
    old marker format.
- Acceptance:
  - Tests cover global/local paths, both tool directories, zero/one/two
    instruction files, repeated actions, stale content, names-only output, and
    unknown IDs without touching the user's real home.

### Unit RX-10 - Standardize Output And Exit Codes

- Goal: keep text human-readable and JSON deterministic.
- Dependencies: RX-05 through RX-09.
- Exit map:
  - `0` success
  - `1` unexpected operational failure
  - `2` usage or TypeBox validation failure
  - `3` configuration resolution/decoding failure
  - `4` release/network failure
  - `5` installation/upgrade/filesystem mutation failure
  - `130` interruption
- Actions:
  - Send ordinary results to stdout and diagnostics to stderr.
  - Emit one valid JSON document for JSON mode.
  - Never mix progress or ANSI decoration into JSON stdout.
  - Preserve the exact delegated command exit code for `runx run`.
- Acceptance:
  - CLI-level tests assert stdout, stderr, JSON parsing, and exit codes.

### Unit RX-11 - Complete Upgrade, Catalog, And Post-Upgrade Reconciliation

- Goal: retain transactional replacement while adding the exact RFC interface.
- Dependencies: RX-05, RX-07, RX-09, and RX-10.
- `runx upgrade` flags:
  - `--version <version>`
  - `--arch <x64|arm64>`
  - `--variant <baseline|default|modern>`
  - `--dry-run`
  - `--format <text|json>`
- `runx upgrade list` flags:
  - `--page <positive-integer>`
  - `--per-page <positive-integer>`
  - `--pre-releases`
  - applicable format/compatibility flags.
- Actions:
  1. Default x64 variant to `baseline`.
  2. List stable releases by default, latest first.
  3. Show real-time text progress while keeping JSON clean.
  4. Replace and verify the canonical executable transactionally.
  5. After success, update the global skill in both tool locations.
  6. When inside a project, reconcile instruction blocks.
  7. Commit the standardized cache only after verified replacement.
- Acceptance:
  - Dry-run, exact version, architecture, variant, pagination,
    pre-release listing, rollback, cache warning, and reconciliation tests pass.

### Unit RX-12 - Rebuild Both Direct Installers

- Goal: make `devops/install.sh` and `devops/install.ps1` complete RFC
  installers.
- Dependencies: RX-09 and RX-11.
- Actions:
  1. Select the exact `darwin`, Linux, or Windows asset.
  2. Print the required sequence heading, target version, architecture,
     variant, and source URL before download.
  3. Display real-time download progress; remove silent download flags.
  4. Validate integrity and architecture.
  5. Transactionally install and verify the binary.
  6. Add the global binary directory to PATH when missing.
  7. Download/extract `guiho-s-runx` into both global skill locations.
  8. Download `guiho-i-runx`, discover instruction files, and reconcile their
     managed blocks.
  9. Log every action and final `runx --version` verification.
- Acceptance:
  - Isolated Windows and POSIX installer tests cover success, fallback,
    corruption, rollback, PATH, both skills, instructions, and progress.

### Unit RX-13 - Replace The Npm Launcher With A Node-Compatible Bootstrap

- Goal: allow `@guiho/runx` installation and first execution without Bun.
- Dependencies: RX-11 and final asset names.
- Expected files: package metadata and `scripts/runx-bin.mjs` as the isolated
  launcher/bootstrap source.
- Actions:
  1. Keep all RunX domain logic in native binaries.
  2. Replace the package `bin` target with a small Node-compatible ESM bootstrap
     at `scripts/runx-bin.mjs`; this is the sole exception to the Node-import
     ban.
  3. Detect platform/architecture and choose the exact RFC asset, defaulting
     x64 to baseline.
  4. Cache the versioned native binary, download only when absent, apply Unix
     executable permissions, pass args/stdio/env unchanged, and preserve the
     native exit code.
  5. Remove Bun source and library fallbacks from the published launcher.
- Acceptance:
  - Packed-package tests run the bootstrap with Node and no Bun in PATH.

### Unit RX-14 - Publish Exactly Fourteen Release Assets

- Goal: make the build and GitHub release workflows enforce the exact release
  contract.
- Dependencies: RX-09, RX-12, and RX-13.
- Binary names:
  - `runx-linux-arm64`
  - `runx-linux-x64`
  - `runx-linux-x64-baseline`
  - `runx-linux-x64-modern`
  - `runx-darwin-arm64`
  - `runx-darwin-x64`
  - `runx-darwin-x64-baseline`
  - `runx-darwin-x64-modern`
  - `runx-windows-arm64.exe`
  - `runx-windows-x64.exe`
  - `runx-windows-x64-baseline.exe`
  - `runx-windows-x64-modern.exe`
- Agent asset names:
  - `guiho-s-runx`
  - `guiho-i-runx`
- Actions:
  1. Replace every `macos` name and selector with `darwin`.
  2. Use Bun APIs in the binary builder rather than Node filesystem/path
     imports.
  3. Package the skill directory and instruction/prompt artifact reproducibly.
  4. Upload only the fourteen expected assets.
  5. Make CI fail for missing, duplicate, extra, wrongly suffixed, or legacy
     assets.
- Acceptance:
  - Automated verification observes exactly fourteen unique names.

### Unit RX-15 - Align Documentation, Skill, TODO, And XDocs

- Goal: remove the obsolete pre-RFC public contract from every durable surface.
- Dependencies: RX-01 through RX-14 behavior stable.
- Expected files:
  - `README.md`, `DOCS.md`, `CHANGELOG.md`
  - `AGENTS.md`
  - `skills/guiho-s-runx/SKILL.md`
  - `TODO.md`, task spec, implementation notes
  - architecture/requirements/decision records that still promise aliases,
    upward discovery, or a home page
  - affected descriptors.
- Actions:
  1. Document the final command catalog and exact configuration precedence.
  2. Mark incompatible older RunX command/compatibility decisions as
     superseded; do not leave contradictory approved documents.
  3. Document agent commands, markers, targets, prompt output, storage, cache,
     installers, wrapper, exit codes, and fourteen assets.
  4. Update the bundled skill to match the executable.
  5. Maintain descriptor `files` and `documents` maps and tree links.
  6. Keep the TODO task open until full validation passes.
- Acceptance:
  - Repository search finds no public promise for `runx r`, root selector
    shorthand, `--file`, `runx agents`, `macos` asset names, or the old startup
    home page.

### Unit RX-16 - Full Validation And Delivery Handoff

- Goal: prove RFC compliance without publishing or mutating real installations.
- Dependencies: RX-15.
- Checks:
  1. `bun run typecheck`
  2. `bun test`
  3. `bun run build`
  4. `bun run binary`
  5. `bun run binaries`
  6. all CLI contract cases from RFC 0034
  7. Node-only packed npm bootstrap smoke tests
  8. isolated PowerShell and POSIX installer tests
  9. static prohibited-import scan
  10. exact fourteen-asset verification
  11. narrow strict xdocs metadata, tree, scan, and doctor checks
  12. `git diff --check` and final scoped status.
- Evidence:
  - Write a validation report under `docs/validation/`.
  - Run an implementation review and a plan-conformance review.
- Acceptance:
  - Every RFC completion-gate item has direct evidence.
  - Failures or skipped checks keep the TODO in `testing` or `todo` with exact
    residual risk.
- Approval gates:
  - Do not apply a Mirror version, tag, publish, push, install globally, or
    overwrite the user's live RunX executable without separate explicit
    authorization.

## First Executable Unit

Start with RX-01. Establish clean baseline evidence and the violation inventory,
then execute RX-02 and RX-03 before changing configuration or command behavior.
This prevents Bun-only and module-boundary work from being entangled with the
public command migration.

## Completion Definition

The task is complete only when all sixteen units pass, the final command catalog
is the only documented catalog, core source is Bun-only, TypeBox covers every
structured boundary, both installers and the npm wrapper work in isolation, and
the release workflow proves exactly twelve binaries plus the two named agent
assets.

## References

- [RunX TODO](../../TODO.md)
- [RFC 0034 task specification](../todo/rfc-0034-cli-compliance-migration.md)
- [RunX CLI architecture](../architecture/cli-architecture.md)
- [Previous Citty migration plan](./citty-cli-migration.md)
- [Canonical package documentation](../../DOCS.md)
