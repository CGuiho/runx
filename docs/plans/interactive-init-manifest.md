---
name: RunX Interactive Init Manifest Plan
purpose: Sequence implementation of the accepted runx init workflow and manifest contract.
description: Defines the schema, interactive terminal UI, validation, documentation, and pull-request work for RunX initialization.
created: 2026-07-14
flags:
  - approved
tags:
  - plans
  - cli
  - manifest
keywords:
  - runx init
  - runx.yaml
  - semantic versioning
  - public group
  - scripts directory
owner: runx-plans
---

# RunX Interactive Init Manifest Plan

## Sources

- [Interactive Init Manifest Decision](../decisions/interactive-init-manifest.md)
- [Alpha Command Catalog Requirements](../requirements/alpha-command-catalog.md)
- [CLI Architecture](../architecture/cli-architecture.md)
- GitHub issue [#10](https://github.com/CGuiho/runx/issues/10)

## Unit 1: Establish the Manifest Contract

- Goal: Make the current manifest validator enforce the accepted `1.x.x`
  Semantic Versioning, scripts directory, public group, and empty catalog
  contract.
- Owner: `C:\GUIHO\runx`.
- Dependencies: Accepted interactive-init decision.
- Files:
  - `source/manifest.ts`
  - `source/types.ts` when TypeScript inference exposes the new manifest shape
  - `source/guiho-runx.spec.ts`
  - `source/cli.spec.ts` fixtures that construct manifests
- Behavior:
  - Replace numeric `version: 1` with a Semantic Versioning string and accept
    only supported major version `1`.
  - Require `scripts.directory` as a non-empty relative path inside the
    manifest root.
  - Require the `public` group with a non-empty summary.
  - Permit an empty `commands` array while keeping every populated command
    explicit about its existing group.
  - Preserve strict unknown-property rejection and existing selector and
    execution safety behavior.
- Data, auth, permissions, and cache: Local YAML contract only. No remote data,
  authentication, authorization, cache, or secret changes.
- Tests:
  - Accept `"1.0.0"` and prerelease/build SemVer forms with major `1`.
  - Reject numeric versions, malformed SemVer, unsupported major versions,
    missing scripts configuration, escaping script paths, missing public, and
    commands that name unknown groups.
  - Prove empty initialized catalogs check and list without execution.
- Acceptance: `readManifest()` returns a strict manifest with `scripts` and
  `public`, and invalid configurations fail with actionable `RunXError`
  messages.
- Stop conditions: Stop if SemVer support would require a migration policy,
  scripts need more than one root, or the change weakens command safety.

## Unit 2: Implement the Interactive Initializer

- Goal: Add a beautiful, cancellable `runx init` Citty command that creates a
  valid empty `runx.yaml` without ever executing a configured command.
- Dependencies: Unit 1.
- Files:
  - `package.json`
  - `bun.lock`
  - `source/init.ts`
  - `source/init.spec.ts`
  - `source/cli.ts`
  - `source/cli.spec.ts`
  - `source/help.ts`
- Behavior:
  - Use a dependency-free Bun terminal prompt adapter with clear visual states,
    so the initializer remains portable across native and Windows builds.
  - Isolate prompt interactions behind an injected interface so prompt flow and
    filesystem behavior can be tested without a real terminal.
  - In an interactive terminal, welcome the user, collect a project name
    defaulting to the selected directory name, collect a scripts directory
    defaulting to `scripts`, preview the exact YAML, and request confirmation.
  - Write only `runx.yaml`; create `scripts/` later when the first real script
    is added.
  - Detect an existing `runx.yaml` in the selected directory and require a
    second explicit overwrite confirmation.
  - Use an atomic temporary-file write and leave no partial manifest on
    cancellation or failure.
  - Fail clearly outside an interactive terminal; do not hang and do not emit
    interactive ANSI frames into JSON output.
  - Accept `--cwd` as the target directory, reject `--file` because the
    initializer always creates `runx.yaml`, and reject `--format json` because
    initialization is an interactive terminal workflow.
  - Register `init` as a first-class Citty command, include it in help/home and
    command-tree output, and reserve the word from selector shorthand.
- Tests:
  - Successful creation exactly matches the accepted YAML shape.
  - Cancellation leaves no file.
  - Existing-file rejection and confirmed overwrite are safe.
  - Default and custom project/scripts values validate.
  - CLI help works outside a manifest and non-interactive invocation exits with
    a clear error.
  - `--file` and `--format json` report clear unsupported-option errors without
    creating files.
  - The generated file passes `runx check` and `runx list`.
- Acceptance: `runx init` produces a polished, deterministic creation flow and
  never executes a manifest command.
- Stop conditions: Stop if the terminal prompt adapter fails Bun tests or native
  compile, if Citty cannot route `init` without conflicting with selector shorthand, or
  if atomic replacement is not portable on Windows.

## Unit 3: Align Canonical Documentation and XDocs

- Goal: Make public and structured documentation accurately describe the new
  manifest and initializer without modifying the bundled agent skill.
- Dependencies: Units 1 and 2.
- Files:
  - `README.md`
  - `DOCS.md`
  - `source/source.xdocs.md`
  - `runx.xdocs.md` when its package description needs updating
  - `docs/plans/plans.xdocs.md`
- Behavior:
  - Document `runx init`, Semantic Versioning, `scripts.directory`, mandatory
    `public`, empty catalog initialization, and command-to-script examples.
  - Keep the bundled `skills/guiho-s-runx/SKILL.md` unchanged; its update belongs
    to the separate agent-integration delivery requested by the project owner.
  - Update the source descriptor for the initializer module and changed CLI /
    manifest responsibilities.
- Checks:
  - `xdocs meta docs\plans --documents --strict --format json`
  - `xdocs meta source --strict --format json`
  - `xdocs tree`
- Acceptance: User-facing docs, XDocs metadata, and implementation agree.
- Stop conditions: Stop if XDocs requires unrelated coverage cleanup; validate
  only the affected subtrees and report wider pre-existing findings.

## Unit 4: Validate, Review, Commit, Push, and Open a Pull Request

- Goal: Produce evidence that the implementation satisfies the accepted
  decision, then deliver it as a draft pull request from `codex/runx-init`.
- Dependencies: Units 1 through 3.
- Checks:
  - `bun run typecheck`
  - `bun test`
  - `bun run build`
  - `bun run binary`
  - `xdocs meta docs\plans --documents --strict --format json`
  - `xdocs meta source --strict --format json`
  - `xdocs tree`
  - `git diff --check`
- Delivery:
  - Write an implementation review and validation report with XDocs metadata.
  - Commit coherent code, tests, documentation, and validation artifacts.
  - Push only the feature branch and create a draft pull request that links
    issue #10 and states that agent-skill automation is intentionally tracked
    separately in issue #11.
- Acceptance: All applicable checks pass, the worktree is clean, the branch is
  pushed, and the draft pull request is reviewable without a release, tag, or
  package publication.
- Stop conditions: Stop before push/PR if validation fails, unrelated changes
  appear, GitHub rejects authentication, or the branch diverges unexpectedly.

## TODO Alignment

This is one approved, focused package implementation and does not need a new
long-running TODO entry. The decision, plan, review, and validation records are
the durable delivery trail.

## First Executable Unit

Unit 1: update `source/manifest.ts` and its fixtures/tests to enforce the
accepted manifest contract.
