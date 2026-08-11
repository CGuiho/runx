---
name: RunX Reveal Implementation Note
purpose: Record execution progress, decisions, validation evidence, and delivery state for TODO task 21.
description: Tracks the RX-REVEAL-1 implementation of the non-executing reveal command and its documentation and release handoff.
created: "2026-08-11"
flags:
  - completed
tags:
  - todo
  - implementation
  - cli
keywords:
  - runx reveal
  - RX-REVEAL-1
  - issue 47
owner: runx-todo-done
---

# RunX Reveal Implementation Note

## Summary

RX-REVEAL-1 adds the public `runx reveal <uid-or-selector-or-index>` Cobra
command. It reuses the existing catalog loader and resolver, prints the stored
command exactly once on stdout, and never invokes execution or confirmation.

## Links

- Task: [runx-reveal.md](runx-reveal.md)
- Plan: [../../plans/runx-reveal.md](../../plans/runx-reveal.md)
- Plan review: [../../reviews/plans/runx-reveal-review.md](../../reviews/plans/runx-reveal-review.md)
- Question ledger: [../../questions/runx-reveal/plan-unit-1.md](../../questions/runx-reveal/plan-unit-1.md)

## Decisions

- Keep the first public contract to one selector plus `--cwd`, `--config`, and
  `--verbose`; omit output-format, child-argument, confirmation, and dry-run
  flags.
- Print `selected.Command` verbatim with one trailing newline so output can be
  copied into the caller's shell without guessing shell quoting.
- Preserve existing selector precedence by calling `Catalog.Resolve` directly.

## Progress

- 2026-08-11T07:12:25Z — Claimed RX-REVEAL-1 on `codex/runx-reveal` at planning
  head `22f029560d322995aaaacbfbf1152d1ebca68b98`; loaded the approved plan,
  review, TODO, question ledger, and required CLI/documentation skills.
- 2026-08-11T07:20:00Z — Moved TODO task 21 to `testing` before final
  formatting, test, build, XDocs, and Mirror-plan validation.
- 2026-08-11 — The initial plain `git push` attempt failed with
  `schannel: AcquireCredentialsHandle failed: SEC_E_NO_CREDENTIALS`; the root
  context retried with its authorized GitHub credentials and published
  `codex/runx-reveal` through head `9b3db9a`, then corrected the delivery note
  and published head `6cbe047`.
- 2026-08-11 — After PR #50 exposed an advanced `main`, fetched
  `origin/main` at `00789a4` and merged it without rewriting this branch. The
  merge preserves PR #46 confirmation-policy implementation, review, and
  release-preparation evidence while retaining the reveal command and tests.
  TODO task 21 records reveal delivery because main already owns tasks 19 and
  20 for confirmation policy and its patch release.
- 2026-08-11 - Committed the merge as `6e400a0` (parents `6cbe047` and
  `00789a4`) and indexed the merged PR46 validation descriptor in `bc4375a`.
- 2026-08-11 - Addressed independent review findings in `8ada893` (plan and
  review history) and `ed6b96a` (unsupported reveal option/argument tests).
- 2026-08-11 - Root context completed the final no-ff history-only merge of
  `origin/main` `c84f643` (public `@guiho/runx/v0.11.1`) as `da614c3`, with
  parents including prior reveal head `9d3f2ff`; the tree is unchanged.
- 2026-08-11 - Merged the subsequently archived 0.11.1 evidence from
  `origin/main` `eb04196` as `844c3f2`, preserving completed tasks 19 and 20,
  reveal task 21, and all reveal descriptors.

## Validation Evidence

- `gofmt -l main.go cmd pkg embed devops` — clean; `git diff --check` — clean.
- `go mod tidy` — `go.mod` and `go.sum` unchanged.
- `go test -count=1 ./...` — passed.
- `go vet ./...` and `go build ./...` — passed.
- `go run devops/build-binaries.go --version 0.12.0 --commit
  22f029560d322995aaaacbfbf1152d1ebca68b98 --build-date
  2026-08-11T00:00:00Z` — eight binaries and three supporting artifacts.
- `go run devops/verify-release-assets.go` — exactly 11 assets and every
  checksum verified.
- `xdocs meta` strict checks for `cmd`, `skills/guiho-s-runx`, `prompts`, and
  `docs/todo` — passed; `xdocs tree` — complete; `xdocs doctor
  --warnings-as-errors` — zero errors and warnings.
- `mirror config check` — passed; `mirror version plan minor` — `0.11.0` to
  `0.12.0`, tag `@guiho/runx/v0.12.0`.

Post-merge exact-head validation at `6e400a0`:

- `gofmt -l main.go cmd pkg embed devops` - clean; `git diff --check` - clean.
- `go mod tidy` followed by `git diff --exit-code -- go.mod go.sum` - clean.
- `go test -count=1 ./cmd` - passed; `go test -count=1 ./...` - passed.
- `go vet ./...` - passed; `go build ./...` - passed.
- `go run devops/build-binaries.go --version 0.12.0 --commit
  6e400a0008c307dd5901026a73afebeb06f3e22b --build-date
  2026-08-11T00:00:00Z` - eight binaries and three supporting artifacts.
- `go run devops/verify-release-assets.go` - exactly 11 assets and every
  SHA-256 checksum verified.
- Strict XDocs metadata for `cmd`, `skills/guiho-s-runx`, `prompts`,
  `docs/todo`, and `.` - passed; `xdocs tree` - complete; `xdocs doctor
  --warnings-as-errors` - zero errors and warnings.
- `mirror config check` - passed; `mirror version plan minor` - `0.11.0` to
  `0.12.0`, tag `@guiho/runx/v0.12.0`.

No production mutation or Mirror apply occurred.

Review-correction validation at `ed6b96a`:

- `go test -count=1 ./cmd -run TestReveal` - passed, including rejection of
  `--format`, `--yes`, `--dry-run`, and extra child arguments.
- `go test -count=1 ./...` - passed; `go vet ./...` - passed; `go build ./...`
  - passed.
- `gofmt -l main.go cmd pkg embed devops`, `git diff --check`, and
  `go mod tidy` followed by the `go.mod`/`go.sum` clean-diff check - clean.
- Strict XDocs metadata for plans, plan reviews, `cmd`, `docs/todo`, and `.`
  - passed; `xdocs tree` - complete; `xdocs doctor --warnings-as-errors` -
  zero errors and warnings.

Final integration validation at `da614c3`:

- `go test -count=1 ./cmd -run TestReveal` - passed; `go test -count=1 ./...`
  - passed.
- `go vet ./...` - passed; `go build ./...` - passed.
- `gofmt -l main.go cmd pkg embed devops`, `git diff --check`, and
  `go mod tidy` followed by the `go.mod`/`go.sum` clean-diff check - clean.
- Strict XDocs metadata for `cmd`, `skills/guiho-s-runx`, `prompts`,
  `docs/plans`, `docs/reviews/plans`, `docs/todo`, and `.` - passed;
  `xdocs tree` - complete; `xdocs doctor --warnings-as-errors` - zero errors
  and warnings.

Final archive-integration validation at `844c3f2`:

- `go test -count=1 ./cmd -run TestReveal` - passed; `go test -count=1 ./...`
  - passed.
- `go vet ./...` - passed; `go build ./...` - passed.
- `gofmt -l main.go cmd pkg embed devops`, `git diff --check`, and
  `go mod tidy` followed by the `go.mod`/`go.sum` clean-diff check - clean.
- Strict XDocs metadata for `cmd`, `skills/guiho-s-runx`, `prompts`,
  `docs/plans`, `docs/reviews/plans`, `docs/todo`, `docs/validation`, and `.`
  - passed; `xdocs tree` - complete; `xdocs doctor --warnings-as-errors` -
  zero errors and warnings.

## Handoff

Implementation review and validation must bind their evidence to the exact
PR head. Ready PR #50 targets `main`, references issue 47 without closing it,
and now includes the no-rewrite merge of origin/main `00789a4`. Mirror
application, merge, release, and worktree cleanup remain integration-gated.

## Final Acceptance

Exact-head review and validation accepted the Luna implementation, PR #50 and
evidence PR #53 reached main, and RunX 0.12.0 published successfully with the
standard 11 assets. Public Windows AMD64 checksum, version, and reveal smoke
passed. Issue 47 remained open only for the separately completed shell fix.
