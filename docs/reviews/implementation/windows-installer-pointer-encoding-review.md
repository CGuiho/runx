---
name: RunX Windows Installer And Upgrade Review
purpose: Materialize the exact-head implementation review for the Windows installer and protocol-v1 upgrade corrections.
description: Accepts BOM-free activation pointers, canonical payload lookup, compatibility aliases, rollback coverage, and the native Windows release gate.
created: 2026-08-23
flags:
  - accepted
tags:
  - review
  - implementation
  - windows
keywords:
  - PR 59
  - current.json
  - protocol-v1 upgrade
owner: runx-implementation-reviews
---

# RunX Windows Installer And Upgrade Review

## Verdict

Accepted with no blocker, high, or medium findings at implementation head
`538476c9176f8ffe047a0df771498a5c2ba3e2fa` in PR
[CGuiho/runx#59](https://github.com/CGuiho/runx/pull/59).

## Findings

- Windows PowerShell 5.1 now writes `current.json` as explicit BOM-free UTF-8.
- The pointer reader accepts only the legacy UTF-8 BOM before applying strict
  JSON field and pointer validation.
- Protocol-v1 release lookup removes the existing `runx-` prefix before
  constructing canonical payload names.
- Compatibility payload and launcher aliases are manifest-declared,
  checksummed, and required to be byte-identical to canonical artifacts.
- Same-version payload state is backed up and restored when activation fails.
- Release publication is blocked on the native Windows installer integration
  test.

## Acceptance Criteria

The implementation satisfies the task criteria. The one post-merge acceptance
boundary is a live upgrade from a released 0.14.3–0.14.5 client after the
compatibility aliases become public.

## Evidence

- Review comment: https://github.com/CGuiho/runx/pull/59#issuecomment-5385635421
- Validated PR head: `538476c9176f8ffe047a0df771498a5c2ba3e2fa`
- Merge commit: `a8eb075304a5b1ad36fdf583b430786f560c0f5d`
