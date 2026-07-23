---
name: Use The Simplified RunX Installation Command
purpose: Define completion signals for TODO task 11 and GitHub issue 24
description: Captures the exact public curl command while retaining downstream installer integrity and verification.
created: 2026-07-22
flags:
  - testing
tags:
  - todo
  - installer
keywords:
  - curl
  - issue 24
owner: runx-todo
---

# Use The Simplified RunX Installation Command

## Todo Index

- Task: `11. Use The Simplified RunX Installation Command`
- Status: testing
- Index: [TODO.md](../../TODO.md)

## Outcome

The README publishes the exact `curl -fsSL ... | bash` bootstrap while the
installer continues to verify the selected binary and agent assets.

## Acceptance Signals

- A source test locks the canonical README command.
- The command succeeds in an isolated-home installation test.
- Internal HTTPS and integrity safeguards remain.

## External Trackers

- GitHub: [CGuiho/runx#24](https://github.com/CGuiho/runx/issues/24) - status: open
