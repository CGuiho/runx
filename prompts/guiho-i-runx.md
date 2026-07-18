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
`runx run <uid> --dry-run` before execution. Never add `--yes` unless the
developer explicitly authorizes the confirmation-gated command.
