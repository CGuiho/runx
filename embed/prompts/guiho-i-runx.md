---
name: guiho-i-runx
purpose: Guide an AI agent through safe RunX command-catalog inspection and execution.
description: Canonical raw instruction prompt bundled with RunX releases.
created: 2026-07-18
flags:
  - agent-prompt
tags:
  - runx
  - cli
keywords:
  - runx.yaml
  - command catalog
  - dry run
owner: runx-prompts
---

# RunX Agent Instruction

Load the `guiho-s-runx` skill. From the intended project directory, run
`runx check --format json`, then `runx list --format json`. Prefer a stable UID,
inspect unfamiliar work with `runx describe <uid>`, and run
`runx run --dry-run <uid>` before execution. RunX options belong before the
selector; every token after it is forwarded to the child. Never add `--yes` unless the
developer explicitly authorizes the confirmation-gated command.

When editing catalogs, use manifest v2: required `namespace`, recursive
`commands`, command `id` leaves, group `group` nodes, and explicit reciprocal
`runx`/`parent` references. Reject legacy split groups and inspect foreign
GitHub child catalogs as executable code before running them.
