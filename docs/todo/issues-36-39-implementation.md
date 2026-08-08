---
name: Resolve RunX Issues 36 And 39 Implementation Notes
purpose: Preserve execution progress, validation evidence, and delivery handoff for TODO task 17.
description: Records the implementation choices, changed behavior, checks, Mirror plan, and exact-head handoff for the combined RunX issues 36 and 39 unit.
created: 2026-08-08
flags:
  - testing
tags:
  - todo
  - implementation
  - validation
keywords:
  - issue 36
  - issue 39
  - selector resolution
  - scripts directory
owner: runx-todo
---

# Resolve RunX Issues 36 And 39 Implementation Notes

## Summary

The dedicated `codex/issues-36-39` worktree implements the approved combined
manifest selector and initializer-default changes from base `364fdb8`.

## Decisions

- Track UID and canonical-selector namespaces separately from group-scoped ID
  owners so exact UID lookup cannot be overwritten by an ID shorthand.
- Permit UID-vs-ID equality across different commands; reject duplicate UIDs,
  UID/canonical-selector collisions, canonical-selector/ID collisions, and
  arbitrary resolution of ambiguous duplicate IDs.
- Keep numeric index lookup after textual identity resolution.
- Emit `.scripts` from `runx init`; parsing and loading preserve explicit
  `scripts.directory` values.
- Update current README, CLI reference, architecture/requirements/decision/RFC
  guidance, prompts, canonical skill, embedded skill, and XDocs metadata. Leave
  legacy `source/` references and historical initializer docs unchanged.

## Progress Log

- `in progress`: task spec and question ledger created; feature branch/worktree
  verified from `364fdb8`.
- `testing`: implementation, focused regressions, current docs, changelog, and
  XDocs descriptors updated.

## Validation Evidence

- `gofmt -l main.go cmd pkg embed devops` — passed with no output.
- Focused manifest and init tests — passed.
- `go test ./...` — passed.
- `go vet ./...` — passed.
- `go build ./...` — passed after the required elevated retry for Go VCS
  metadata access in the isolated worktree.
- `xdocs scan`, `xdocs tree`, and `xdocs doctor --warnings-as-errors` — passed.
- `xdocs meta docs/todo --documents --strict`,
  `xdocs meta docs/questions --documents --strict`, and
  `xdocs meta skills/guiho-s-runx --documents --strict` — passed.
- `mirror config check` — passed.
- `mirror version plan minor` — planned `@guiho/runx/v0.11.0`; no apply/tag/
  release action performed.

## Handoff

The branch will be pushed and one pull request opened against `main` with
`Closes #36` and `Closes #39`. Independent implementation review and validation
must bind to the exact PR head before integration. No production mutation has
occurred.
