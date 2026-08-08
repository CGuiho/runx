---
name: RunX Issues 36 And 39 Question Ledger
purpose: Record provisional answers for the unattended execution of plan unit runx-issues-36-39.
description: Preserves evidence, selected answers, rationale, confidence, reversibility, and pending human review for material execution choices.
created: 2026-08-08
flags:
  - in-progress
tags:
  - questions
  - execution
  - runx
keywords:
  - issue 36
  - issue 39
  - selector resolution
  - init default
owner: runx-questions-issues-36-39
---

# RunX Issues 36 And 39 Question Ledger

## Execution Context

- Plan unit: `runx-issues-36-39`
- Branch: `codex/issues-36-39`
- Worktree: `C:\GUIHO\runx\.temp\issues-36-39`
- Approved base: `364fdb8`
- Human review status: pending

## Q1. Which identity key wins when a UID equals another command's ID?

- Evidence: Issue #36 explicitly requires globally unique UIDs, group-scoped
  IDs, exact UID precedence, and preservation of ambiguous-ID errors. Existing
  manifest composition already retains canonical selectors and numeric IDX
  fallback.
- Candidate answers:
  - Reject any UID/ID equality across commands.
  - Allow UID-vs-ID equality across commands and resolve exact UID first.
- Chosen answer: Allow UID-vs-ID equality across different commands; resolve an
  exact global UID before canonical selector or unique unqualified ID lookup.
- Rationale: This is the stated issue contract and is deterministic while
  preserving true cross-command collisions and ambiguous ID errors.
- Confidence: high
- Reversibility: high; validation and lookup helpers remain localized.
- Action: Implement parser/composition/Resolve changes and focused regressions.
- Human review status: pending

## Q2. What default should `runx init` emit for `scripts.directory`?

- Evidence: Issue #39 explicitly changes the current Go initializer default
  from `scripts` to `.scripts`; explicit manifest values must remain unchanged.
- Candidate answers:
  - Keep `scripts` for compatibility.
  - Emit `.scripts` only when the user does not provide a directory.
- Chosen answer: Emit `.scripts` as the initializer default and preserve any
  explicit configured directory value verbatim.
- Rationale: This is the stated issue contract and does not alter existing
  manifests with explicit configuration.
- Confidence: high
- Reversibility: high; one default constant and its regression coverage.
- Action: Update current Go init path, tests, and current docs/examples.
- Human review status: pending

## Residual Questions

No additional material uncertainty was found after reading the approved issue
contracts, current Go implementation, repository instructions, and existing
manifest composition requirements.
