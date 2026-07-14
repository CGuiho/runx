---
name: RunX Interactive Init Manifest Plan Review
purpose: Verify that the interactive-init implementation plan is executable, safe, and traceable.
description: Reviews manifest constraints, terminal workflow boundaries, test coverage, XDocs work, and pull-request delivery for runx init.
created: 2026-07-14
flags:
  - approved
tags:
  - reviews
  - plans
  - cli
keywords:
  - runx init
  - runx.yaml
  - semantic versioning
  - public group
  - scripts directory
owner: runx-plan-reviews
---

# RunX Interactive Init Manifest Plan Review

## Verdict

Ready for execution.

## Findings

No blocker, high, medium, or low findings remain.

The plan is traceable to the accepted interactive-init decision, the existing
command-catalog requirement, and the Citty CLI architecture. The only review
finding was resolved in the plan: `runx init` is an interactive terminal
workflow, so `--format json` and `--file` must fail clearly rather than create
ambiguous output or target files.

## Sequencing and Safety

- Manifest schema changes precede the initializer, so generated YAML is
  validated by the same strict contract that later commands use.
- The initializer is isolated from the executor and is explicitly prohibited
  from spawning configured manifest commands.
- The prompt boundary is injectable for tests; an interactive terminal is
  required for the real user interface.
- Existing-file confirmation, cancellation, and atomic writes prevent silent
  replacement and partial manifests.
- The `scripts` directory is configured but not created until a real script is
  added, matching the accepted file-only initialization boundary.
- Manifest evolution is constrained to supported SemVer major `1`; no legacy
  numeric compatibility or migration policy is hidden in the implementation.

## Coverage Review

- Tests cover SemVer, public-group, scripts-directory, empty-catalog, command
  group-reference, creation, cancellation, overwrite, non-interactive, and
  unsupported-option behavior.
- The full Bun typecheck, test, library build, local native compilation, XDocs
  subtree validation, and diff hygiene are named.
- Native compilation of the terminal prompt dependency is an explicit stop
  condition.
- No API, database, authentication, authorization, cache, cloud, secret,
  deployment, release, tag, or package-publication work is in scope.

## Documentation and TODO Alignment

The plan updates public CLI documentation and the affected XDocs descriptors.
The bundled RunX agent skill is deliberately deferred to the separate
agent-integration feature tracked in issue #11, as directed by the project
owner. The work is one focused package session, so no additional long-running
TODO entry is required.

## First Executable Unit

Update `source/manifest.ts` and its tests/fixtures to establish the SemVer,
scripts-directory, mandatory-public-group, and empty-catalog contract.

## References

- [Plan](../../plans/interactive-init-manifest.md)
- [Decision](../../decisions/interactive-init-manifest.md)
- [Requirements](../../requirements/alpha-command-catalog.md)
- [CLI Architecture](../../architecture/cli-architecture.md)
