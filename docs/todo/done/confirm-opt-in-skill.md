---
name: Make RunX Confirmation Opt-In In The Agent Skill
purpose: Preserve the completed confirmation-policy correction and its release acceptance for TODO task 19.
description: Records the explicit opt-in policy for manifest confirmation and the completed 0.11.1 patch-release handoff.
created: 2026-08-11
flags:
  - completed
tags:
  - todo
  - done
  - agents
  - runx
keywords:
  - confirm
  - opt-in
  - guiho-s-runx
  - 0.11.1
owner: runx-todo-done
---

# Make RunX Confirmation Opt-In In The Agent Skill

## Completion

- Task: `19. Make RunX confirmation opt-in in the agent skill`
- Status: completed
- Implementation PR: [CGuiho/runx#46](https://github.com/CGuiho/runx/pull/46)
- Merge commit: [37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8](https://github.com/CGuiho/runx/commit/37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8)
- Protected evidence: [integration review](../../reviews/implementation/runx-0.11.1-confirm-opt-in-review.md), [integration validation](../../validation/runx-0.11.1-confirm-opt-in.md)
- Final release validation: [RunX 0.11.1 release](../../validation/runx-0.11.1-release.md)

## Outcome

The canonical and embedded `guiho-s-runx` skills omit the manifest `confirm`
field by default and add it only when the user explicitly requests
`confirm: never` or `confirm: always` for that specific command. Omission
resolves to `never`; existing commands that explicitly declare `confirm:
always` retain their authorization boundary.

The proactive instruction that inferred `confirm: always` for destructive,
release, deployment, migration, or production-impacting commands was removed.
Canonical and embedded skill metadata were bumped from `0.4.0` to `0.4.1`, and
focused embedded-policy regression coverage was added. Runtime manifest parsing
and execution behavior were unchanged.

## Validation and release

Focused and full Go tests, vet, build, strict XDocs checks, Mirror checks, and
the release-contract CI passed. The correction was published in the
Mirror-managed [0.11.1 release](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.11.1),
with exact 11-asset, checksum/digest, Windows AMD64, and npm latest evidence
recorded in the final release validation.

No deployment, promotion, traffic, DNS, database, secret, or other
production-infrastructure mutation occurred.
