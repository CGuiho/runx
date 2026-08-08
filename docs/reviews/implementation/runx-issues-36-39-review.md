---
name: RunX PR #40 Implementation Integration Review
purpose: Materialize the accepted implementation-review evidence for the merged RunX issues 36 and 39 pull request.
description: Records the immutable review verdict, exact reviewed head, scope, checks, and integration boundary for PR #40.
created: 2026-08-09
flags:
  - accepted
  - integrated
tags:
  - reviews
  - implementation
  - integration
keywords:
  - runx
  - pull request 40
  - issue 36
  - issue 39
  - exact head
  - implementation review
owner: runx-implementation-reviews
---

# RunX PR #40 Implementation Integration Review

## Verdict

Accepted for validation at exact head
64ab9ea0643763123dbd374b83f566b42172159e against base
364fdb8d8038bafd984a2cc978631dc07488a7b0. This verdict was bound to that
commit and was reobserved before merge.

## Immutable Review Evidence

- Pull request: [CGuiho/runx#40](https://github.com/CGuiho/runx/pull/40).
- Review record: [implementation review](https://github.com/CGuiho/runx/pull/40#pullrequestreview-4889711530).
- Review submitted: 2026-08-08T20:48:44Z.
- Reviewed branch/base: codex/issues-36-39 -> main.
- Exact reviewed head: [64ab9ea0643763123dbd374b83f566b42172159e](https://github.com/CGuiho/runx/commit/64ab9ea0643763123dbd374b83f566b42172159e).
- Scope: 25 files, 663 additions, 66 deletions; current Go implementation/tests, current documentation, embedded/canonical resources, task/question/XDocs records, and changelog. Legacy source/ TypeScript remained untouched.
- Findings: no unresolved blocker, high, medium, or low implementation finding.
- Review threads: none.

## Acceptance Criteria

- Issue #36: global UID lookup is distinct from canonical selector and unique-ID shorthand lookup. A UID may equal another command's scoped ID, exact UID lookup wins, true collisions remain rejected, and ambiguous unqualified IDs do not select arbitrarily.
- Resolution precedence is exact global UID, canonical selector, unique unqualified ID shorthand, then numeric index.
- Issue #39: current Go runx init emits scripts.directory: ".scripts"; explicit configured directories remain unchanged.
- The implementation, focused regressions, current docs, canonical/embedded resources, task/question records, XDocs metadata, and changelog are in scope.
- No legacy source/ TypeScript behavior was changed.

## Verification References

- CI: [run 31277425590](https://github.com/CGuiho/runx/actions/runs/31277425590), with go (windows-latest), go (ubuntu-latest), and release-contract jobs completed successfully.
- Exact-head validation: [PR #40 validation comment](https://github.com/CGuiho/runx/pull/40#issuecomment-5228132822).
- Issue #36: https://github.com/CGuiho/runx/issues/36.
- Issue #39: https://github.com/CGuiho/runx/issues/39.

## Integration Boundary

This record was materialized after merge commit
[e3aed59680c4eea4d18597c3ccd8754daa77830e](https://github.com/CGuiho/runx/commit/e3aed59680c4eea4d18597c3ccd8754daa77830e)
reached main. No implementation correction was needed. No version tag or
GitHub Release was created under merge-only authorization.
