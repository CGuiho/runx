---
name: Complete RunX Upgrade Reliability
purpose: Track the accepted outcomes for GitHub issues 12 and 13.
description: Records immediate verified replacement, complete listing, exact recovery, direct installers, reviews, validation, and delivery evidence.
created: 2026-07-15
updated: 2026-07-19
flags:
  - completed
tags:
  - todo
  - cli
  - reliability
keywords:
  - RunX
  - issue 12
  - issue 13
  - self-upgrade
owner: runx-todo
---

# Complete RunX Upgrade Reliability

## Todo Index

- Task: `3. Complete RunX Upgrade Reliability`
- Status: completed
- Index: [TODO.md](../../TODO.md)
- External: [issue #12](https://github.com/CGuiho/runx/issues/12) and
  [issue #13](https://github.com/CGuiho/runx/issues/13)

## Outcome

RunX replaces and verifies the canonical executable before success, lists the
complete published catalog, and exposes exact-version recovery instructions
after every upgrade outcome in text and JSON.

## Completion Signals

- Transactional Windows replacement and rollback are tested.
- Catalog pagination, sorting, channels, assets, and single-leaf routing are
  tested.
- Every terminal outcome includes a pinned installer and separate stop command.
- Direct installers normalize exact stable/prerelease versions and reject
  installed-version mismatches.
- Current-main reviews and validation contain pushed evidence for both issues.

## References

- [Approved design](../superpowers/specs/2026-07-15-upgrade-reliability-design.md)
- [Implementation plan](../plans/upgrade-reliability-implementation.md)
- [Umbrella validation](../validation/upgrade-reliability.md)
