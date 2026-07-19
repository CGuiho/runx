---
name: Preserve The RunX Unicode Help Tree
purpose: Track the human-readable tree contract requested by GitHub issue 17.
description: Records Unicode branches, nested guides, aligned descriptions, legacy ASCII rejection, tests, and delivery evidence.
created: 2026-07-19
flags:
  - completed
tags:
  - todo
  - cli
  - help
keywords:
  - RunX
  - issue 17
  - help tree
  - Unicode
owner: runx-todo
---

# Preserve The RunX Unicode Help Tree

## Todo Index

- Task: `5. Preserve The RunX Unicode Help Tree`
- Status: completed
- Index: [TODO.md](../../TODO.md)
- External: [CGuiho/runx issue #17](https://github.com/CGuiho/runx/issues/17)

## Outcome

`runx --help-tree` renders a `COMMAND TREE` with Unicode branch glyphs, nested
vertical guides, aligned descriptions, commands, and flags. The legacy ASCII
`|-` tree cannot regress unnoticed.

## Completion Signals

- Root and every public scope render from the real Citty definitions.
- Tests assert `├──`, `└──`, `│`, aligned descriptions, and no `|-`.
- Depth limiting and hidden-command exclusion remain intact.

## References

- [Implementation review](../reviews/implementation/unicode-help-tree-review.md)
- [Validation](../validation/unicode-help-tree.md)
