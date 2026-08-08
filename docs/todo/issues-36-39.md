---
name: Resolve RunX Issues 36 And 39
purpose: Define completion signals for TODO task 17 and GitHub issues 36 and 39.
description: Captures the manifest identity-resolution contract, the init scripts-directory default, focused regressions, documentation scope, and delivery boundaries.
created: 2026-08-08
flags:
  - testing
tags:
  - todo
  - cli
  - manifest
keywords:
  - issue 36
  - issue 39
  - selector
  - uid
  - scripts directory
owner: runx-todo
---

# Resolve RunX Issues 36 And 39

## Todo Index

- Task: `17. Resolve RunX issues 36 and 39`
- Status: testing
- Index: [TODO.md](../../TODO.md)
- Parent coordination: [GUIHO root TODO](../../../guiho/TODO.md)

## Plan Unit

`runx-issues-36-39` is the approved combined execution unit. A separate
planning phase is waived because the Commander-in-Chief supplied the explicit
issue contracts, base commit, owned paths, validation commands, and delivery
boundary in the execution request.

## Outcome

RunX accepts a globally unique UID that equals another command's group-scoped
ID, resolves exact UIDs before canonical selectors or unique ID shorthands, and
generates `.scripts` as the default `runx init` scripts directory while
preserving explicitly configured directories.

## Scope

### In scope

- Manifest parser, composition, and selector resolution behavior for issue 36.
- `runx init` default `scripts.directory` behavior for issue 39.
- Focused Go regressions and current user-facing/bundled skill documentation.
- Current xdocs descriptors, strict metadata, tree, and doctor validation.
- One feature branch and one pull request targeting `main`.

### Out of scope

- Legacy `source/` TypeScript reference files.
- Release version application, tag creation, GitHub Release publication, or
  production mutation.
- Unrelated historical migration-document rewrites.

## Acceptance Signals

- A manifest containing `cli/test` (`uid: cli-test`, `id: test`) and
  `go/test` (`uid: test`, `id: test`) parses and composes successfully.
- Exact UID lookup wins over canonical selector or ID shorthand; ambiguous
  duplicate unqualified IDs and true cross-command UID/selector collisions
  remain validation errors.
- Numeric index selection still occurs only after identity resolution.
- `runx init` writes `scripts.directory: .scripts` by default and leaves an
  explicit configured directory unchanged.
- Focused tests, `go test ./...`, `go vet ./...`, `go build ./...`, xdocs
  validation, and `mirror version plan` provide evidence for the handoff.

## External Trackers

- GitHub: [CGuiho/runx#36](https://github.com/CGuiho/runx/issues/36) - open
- GitHub: [CGuiho/runx#39](https://github.com/CGuiho/runx/issues/39) - open

## Related Files

- [docs/questions/issues-36-39/plan-unit-1.md](../questions/issues-36-39/plan-unit-1.md) - Provisional execution answers and question ledger.
- [docs/todo/issues-36-39-implementation.md](issues-36-39-implementation.md) - Execution progress, validation evidence, and handoff.
- [docs/plans/manifest-v2-composition.md](../plans/manifest-v2-composition.md) - Existing manifest identity context.
- [skills/guiho-s-runx/SKILL.md](../../skills/guiho-s-runx/SKILL.md) - Bundled current CLI contract.

## Watch-outs

- UID uniqueness remains global; only UID-vs-ID equality across different
  commands is relaxed.
- Canonical selectors and true UID collisions must not be silently shadowed.
- Explicit `scripts.directory` values remain untouched by initializer defaults.
- Do not modify generated outputs or legacy TypeScript reference files.

## Before Starting

- Confirm the dedicated worktree is based on `364fdb8` and is not `main`.
- Keep the coordination TODO entry on `main` as the source of shared lifecycle
  status; this branch owns only the task spec, ledger, implementation, and
  evidence files.

## While Working

- Preserve unrelated user changes and use explicit-path, smallest-coherent
  commits.
- Record material implementation uncertainty in the question ledger and use a
  reversible contract-consistent answer without pausing execution.

## After Finishing

- Move the task to `testing` before final validation through the coordination
  owner, then leave completion/archive state for the integrator after review
  and validation.
- Return the exact branch/worktree/base/head, changed paths, checks, PR, Mirror
  decision, and residual risks to `guiho-a-0001-swe`.

## Mirror Decision

The compatible feature/default-behavior unit is a `minor` candidate. Mirror
planned `@guiho/runx/v0.11.0` from the current `0.10.1` Git baseline. Version
application, tag creation, push, and release publication are deferred to the
pull-request integrator after merge and independent validation.
