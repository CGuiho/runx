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

## Follow-up Findings And Corrections

The post-release gates intentionally remained open and found two additional
issues before completion:

- PR 60 corrected whole-release verification so the checksum manifest does not
  have to hash itself while every ordinary downloaded artifact remains
  mandatory and verified.
- PR 61 made an explicitly requested already-active version return verified
  `up-to-date` without trying to overwrite the running Windows payload.

Both corrections added native Windows integration coverage and passed exact-head
Ubuntu, Windows, and release-contract checks.

## Acceptance Criteria

Accepted after public 0.14.8 installation, public 0.14.7 → 0.14.8 synchronous
upgrade, and public exact 0.14.8 idempotence all passed.

## Evidence

- PR 59 review: https://github.com/CGuiho/runx/pull/59#issuecomment-5385635421
- PR 60 review: https://github.com/CGuiho/runx/pull/60#issuecomment-5385707234
- PR 61 review/validation: https://github.com/CGuiho/runx/pull/61#issuecomment-5385759966
- Merge commits: `a8eb075304a5b1ad36fdf583b430786f560c0f5d`,
  `caec3f7769b37f73a97817652f166d994e9400c7`, and
  `ee6c458262912c2fdabd119f8138cf7abf8112a8`
