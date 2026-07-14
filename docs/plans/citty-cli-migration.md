---
name: RunX Citty CLI Migration Plan
purpose: Execute the accepted full migration from RunX handwritten argument parsing and routing to Citty.
description: Sequences dependency verification, command-tree migration, compatibility tests, documentation, protected delivery, and the next Mirror-managed patch release.
created: 2026-07-14
flags:
  - approved
  - in-progress
tags:
  - plans
  - cli
  - citty
keywords:
  - runx
  - citty
  - argument parsing
  - command routing
  - patch release
owner: runx-plans
---

# RunX Citty CLI Migration Plan

## Source Decision

- [Full Citty CLI Migration](../decisions/citty-cli-migration.md)

## Unit 1: Add and Verify Citty

- Goal: Add Citty as the only CLI parser/router dependency and verify its Bun,
  TypeScript, and native-executable interfaces before replacing production
  routing.
- Dependencies: Accepted Citty migration decision.
- Files:
  - `package.json`
  - `bun.lock`
- Data, auth, and cache: No persisted data, authentication, or cache behavior
  changes.
- Checks:
  - Install `citty` with Bun as a runtime dependency.
  - Inspect its locally installed types for command definitions, aliases,
    positional arguments, raw argument execution, help, and version handling.
  - Run a bounded native compile after integration begins to confirm Citty is
    bundled without a Node.js runtime dependency.
- Acceptance: Citty is declared once as a production dependency and exposes
  the interfaces needed for a testable RunX command tree.
- Stop conditions: Stop if Citty requires Node.js at runtime, cannot be bundled
  by Bun, or cannot express the accepted command and alias contract.

## Unit 2: Replace Parsing and Routing

- Goal: Make Citty own all CLI parsing, help, version handling, subcommands,
  aliases, and dispatch while retaining RunX domain behavior.
- Dependencies: Unit 1 complete.
- Files:
  - `source/cli.ts`
  - `source/flags.ts` (remove)
  - supporting source modules only when separation materially improves the
    command definitions
- Command tree:
  - `runx`
  - `runx list`
  - `runx describe <selector>`
  - `runx run <selector>` and alias `runx r <selector>`
  - `runx check`
  - `runx agents install <local|global>`
  - `runx agents instructions`
  - `runx upgrade`, `runx upgrade check`, and `runx upgrade list`
  - `runx uninstall`
- Compatibility:
  - Preserve root `runx <selector>` as a compatibility shorthand after known
    Citty subcommands take precedence.
  - Preserve `--cwd`, `--file`, `--format`, `--verbose`, `--dry-run`, `--yes`,
    `--help-tree`, `--help-docs`, and `--tool` on their applicable paths.
  - Citty owns `-h`/`--help` and `-v`/`--version`; neither may discover a
    manifest.
  - Unknown options are usage failures; manifest and execution failures remain
    `RunXError` domain failures.
- Data, auth, and cache: Manifest parsing, selectors, subprocess behavior,
  installer destinations, upgrade metadata, and output payloads remain
  unchanged.
- Acceptance: No handwritten token parser or top-level manual command router
  remains, and all accepted command paths dispatch through Citty.
- Stop conditions: Stop if listing, describing, checking, help, version, or dry
  runs spawn a configured command, or if JSON output contracts drift.

## Unit 3: Prove and Document the Public Contract

- Goal: Lock the migrated command tree with tests and align every public
  reference with the Citty-owned interface.
- Dependencies: Unit 2 complete.
- Files:
  - CLI tests under `source/`
  - `README.md`
  - `DOCS.md`
  - `skills/guiho-s-runx/SKILL.md`
  - affected named `*.xdocs.md` descriptors
  - `AGENTS.md` only if its durable command contract needs correction
- Tests:
  - Root help, `-h`, `--help`, `-v`, and `--version` outside a RunX project.
  - Every command, nested command, alias, positional selector, applicable flag,
    and root-selector shorthand.
  - Unknown option, missing argument, manifest/domain error, JSON, dry-run, and
    no-spawn safety behavior.
  - Upgrade and agent paths retain their existing business behavior.
- Documentation: Explain the Citty command tree and generated help without
  exposing an obsolete parser implementation.
- Acceptance: Tests cover all public paths and aliases; the bundled skill and
  XDocs metadata describe the same contract as the executable.
- Stop conditions: Stop on contradictory docs, missing command coverage, or a
  safety regression.

## Unit 4: Validate and Deliver Through Protected Main

- Goal: Demonstrate package and native-executable correctness, then merge the
  implementation without bypassing repository protections.
- Dependencies: Unit 3 complete.
- Checks:
  - `bun run typecheck`
  - `bun test`
  - `bun run build`
  - `bun run binary`
  - `bun run binaries`
  - package dry-run and packed-launcher `-v`/help checks outside a manifest
  - native executable `-v`/help checks outside a manifest
  - targeted strict XDocs metadata and doctor checks
  - `git diff --check`
- Delivery:
  - Commit only understood migration changes on the current `codex/` branch.
  - Push, open a pull request, wait for required CI, and merge through protected
    `main`.
  - Synchronize local `main` after the merge.
- Acceptance: CI passes and remote `main` contains the complete migration.
- Stop conditions: Stop on validation failure, unexpected worktree changes,
  unresolved review, or branch-protection failure.

## Unit 5: Publish the Next Mirror-Managed Patch

- Goal: Release the merged migration as patch `0.2.4` and verify both npm and
  native GitHub distribution.
- Dependencies: Unit 4 merged; user authorization in this task to complete the
  patch release.
- Sequence:
  - Load the Mirror workflow and inspect `mirror.config.toml`.
  - Confirm a clean, synchronized base and run the full release validation.
  - Run `mirror version plan patch` and require `0.2.4`.
  - Prepare the allowed changelog/XDocs updates and commit them.
  - Apply the patch with Mirror; do not hand-edit managed versions.
  - Deliver the release commit through protected `main`, then push the protected
    tag only after the tagged commit is in `main` history.
  - Monitor the `production` workflow while the user performs its manual
    environment approval.
  - Verify the GitHub Release assets and npm version/provenance.
- Acceptance: `package.json`, the Mirror release commit, protected tag,
  successful Publish run, GitHub Release, and npm registry all agree on
  `0.2.4`.
- Stop conditions: Stop before tagging if Mirror plans another version or
  validation fails. If publication fails after tagging, preserve the tag and
  capture the exact failure before deciding any follow-up release.

## TODO Alignment

This is a single authorized execution session with a durable decision, plan,
and review. It does not require a separate long-running TODO entry.

## First Executable Unit

Unit 1: add Citty with Bun and verify the installed API against the accepted
RunX command contract.
