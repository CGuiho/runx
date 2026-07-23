---
name: Forward RunX Child Arguments And Subcommands
purpose: Define completion signals for TODO task 12 and GitHub issue 25
description: Captures argument ownership, shell-safe transport, dry-run visibility, exact exits, documentation, and regression requirements.
created: 2026-07-22
flags:
  - completed
tags:
  - todo
  - execution
keywords:
  - arguments
  - issue 25
owner: runx-todo
---

# Forward RunX Child Arguments And Subcommands

## Todo Index

- Task: `12. Forward RunX Child Arguments And Subcommands`
- Status: completed
- Index: [TODO.md](../../TODO.md)

## Outcome

`runx run <selector> ...` forwards every post-selector argument verbatim and in
order without letting child flags alter RunX routing or become shell source.

## Acceptance Signals

- `runx run cli-ts -v` forwards `-v`.
- RunX options work before the selector and child options work after it.
- Shell adapters pass spaces, empty values, Unicode, quotes, and shell
  metacharacters literally.
- Dry-run text and JSON expose the forwarded array and never spawn.
- Child exit codes are preserved.

## External Trackers

- GitHub: [CGuiho/runx#25](https://github.com/CGuiho/runx/issues/25) - status: closed
