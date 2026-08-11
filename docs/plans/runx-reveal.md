---
name: RunX Reveal Implementation Plan
purpose: Define the sealed unit for implementing, reviewing, integrating, and publishing RunX reveal.
description: Sequences the public Cobra command, selector/no-spawn regressions, documentation, exact-head gates, protected merge, and Mirror-managed minor release.
created: "2026-08-11"
flags:
  - approved
tags:
  - plan
  - cli
  - release
keywords:
  - runx reveal
  - issue 47
  - gpt-5.6-luna
  - minor release
owner: runx-plans
---

# RunX Reveal Implementation Plan

## Summary

Implement one public, non-executing `runx reveal <selector>` command, deliver it
through a dedicated branch and reviewed PR, merge it into protected `main`, and
publish the next Mirror-managed minor release. GPT-5.6 Luna at maximum reasoning
effort owns all implementation code and implementation-side documentation.

## Traceability

- [Feature brainstorm](../brainstorm/runx-reveal.md)
- [TODO task 19 specification](../todo/runx-reveal.md)
- [GitHub issue 47](https://github.com/CGuiho/runx/issues/47)
- [Run argument ownership decision](../decisions/run-argument-ownership.md)

## Unit RX-REVEAL-1: Implement And Deliver Reveal

### Goal

Add a public Cobra command that loads and resolves a catalog command through the
existing manifest path and writes only the exact stored command plus a newline,
without any configured-command execution or confirmation path.

### Ownership And Isolation

- Repository: `C:\GUIHO\runx`.
- Approved base: `8beb28d1690c050bc7345cdbf77da2bb143909e9`.
- Branch: `codex/runx-reveal`.
- Worktree: `C:\GUIHO\runx\.temp\runx-reveal-47`.
- Implementation agent: GPT-5.6 Luna, maximum reasoning effort.
- PR target: `main`.
- Question ledger: [plan-unit-1.md](../questions/runx-reveal/plan-unit-1.md).

### Owned Paths

- `cmd/catalog.go`, `cmd/root.go`, and focused `cmd/root_test.go` changes.
- `cmd/cmd.xdocs.md` and any directly affected current XDocs descriptor.
- `README.md`, `DOCS.md`, `CHANGELOG.md`, and `CONTRIBUTING.md` only where the
  public command contract requires it.
- `skills/guiho-s-runx/SKILL.md`, `embed/skills/guiho-s-runx.SKILL.md`,
  `prompts/guiho-i-runx.md`, and their owning descriptors.
- `docs/todo/runx-reveal-implementation.md` plus task/descriptor milestone
  updates assigned by this plan.

### Shared-File Exclusions And Integration Order

- Do not touch the dirty original `main` checkout or its existing `TODO.md` and
  `.review-pr46-c479bb6/` state.
- Do not alter environment or key files, generated `library/`, `bin/`, `bundle/`,
  or `vendor/`, or legacy `source/` TypeScript reference files.
- PR 46 owns confirmation-policy changes. If it merges first, integrate current
  `main` and preserve those changes before final review. If this PR merges first,
  do not mutate PR 46; its owner must integrate the new main independently.

### Implementation Contract

1. Add `newRevealCommand` to the live Cobra tree with no alias.
2. Reuse `loadCatalog` and `Catalog.Resolve`; do not add selector parsing.
3. Accept exactly one selector. Support catalog `--cwd`, `--config`, and
   diagnostic behavior without adding JSON or child-argument rendering to the
   first public contract.
4. Print `selected.Command` verbatim with exactly one trailing newline.
5. Never call `executor.ExecuteCommand`, confirmation helpers, or shell
   construction from reveal.
6. Preserve established configuration/selector exit mappings.
7. Update live help/docs, public docs, prompt, canonical and embedded skill,
   changelog, and accurate XDocs metadata.

### Required Tests

- Exact stdout for UID, canonical/full selector, unique shorthand, and numeric
  index.
- A `confirm: always` selection reveals without input or prompt.
- Invalid selector exit behavior.
- No-spawn marker coverage includes reveal.
- Help tree/docs expose reveal and alias policy remains unchanged.
- Existing run/dry-run/confirmation/argument tests stay green.

### Validation

- `gofmt -l main.go cmd pkg embed devops` returns no paths.
- `go mod tidy` leaves `go.mod` and `go.sum` unchanged.
- `go test -count=1 ./...`.
- `go vet ./...`.
- `go build ./...`.
- `go run devops/build-binaries.go --version 0.12.0 --commit <head> --build-date <RFC3339>`.
- `go run devops/verify-release-assets.go`.
- Strict XDocs metadata for touched scopes, `xdocs tree`, and
  `xdocs doctor --warnings-as-errors`.
- `mirror config check` and `mirror version plan minor`, expecting
  `@guiho/runx/v0.12.0` from the released `0.11.0` baseline.

### Delivery And Gates

1. Set task 19 to `testing` before final validation without changing the exact
   review head afterward.
2. Commit smallest coherent changes with explicit paths and push the branch.
3. Open a ready PR against `main`, referencing issue 47 without closing the
   deferred Windows shell investigation prematurely.
4. Review the exact PR head without implementation edits. Any blocker/high
   finding returns to Luna and invalidates the prior verdict.
5. Validate the same accepted head and persist review/validation evidence in PR
   comments or reviews.
6. Reobserve required checks, protection, mergeability, and head identity; merge
   only when every gate targets the same commit.
7. Verify the merged commit is reachable from `main` and remove only the merged
   feature branch/worktree after release evidence is complete.

### Release

1. From clean integrated `main`, run plain `mirror`, repository checks,
   `mirror config check`, and `mirror version plan minor`.
2. Update the allowed changelog through the normal implementation/integration
   flow; never hand-edit Mirror-managed version fields.
3. Apply `mirror version apply minor --yes` only after the inspected clean plan.
4. Verify the protected production-named publication workflow, GitHub Release,
   canonical tag, and exactly 11 assets. This workflow publishes artifacts and
   npm; repository evidence shows it does not deploy an application.
5. Download the Windows AMD64 asset and smoke `--version`, `--help`,
   `--help-tree`, and `reveal` against a temporary non-secret manifest.

### Acceptance Criteria

- Every acceptance signal in the task spec is satisfied.
- The implementation was authored by the required Luna/max sub-agent.
- Independent root review and validation accept the exact PR head.
- The change is merged to `main` and the next minor release is publicly
  available with verified assets.
- No secret-bearing file was accessed and no production mutation occurred.

### Stop Conditions

- Stop before merge for mismatched review/validation heads, failed checks,
  unresolved conflicts, or protection/mergeability failure.
- Stop before Mirror apply if the worktree is dirty, the plan is unexpected,
  release automation performs production deployment, or asset verification is
  not ready.
- Authentication/connectivity failures are technical blockers; never bypass
  hooks, protection, checks, or force-push.

## First Executable Unit

RX-REVEAL-1 is the only implementation unit and is ready for unattended Luna
execution after this planning baseline is committed and pushed.
