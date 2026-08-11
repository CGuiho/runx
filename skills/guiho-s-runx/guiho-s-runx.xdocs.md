---
subject: guiho-s-runx
description: Agent workflow for safe RunX catalog work, deterministic UID/canonical/ID selector resolution, exact non-executing command revelation, caller-aware Windows automatic shell selection, interactive numeric-index selection, user-requested confirmation, bare-command bootstrap, automatic resource maintenance, and verified native upgrade/list/recovery operations.
parent: runx-skills
children: []
files:
  SKILL.md: RFC-aligned RunX catalog inspection and exact-command revelation, deterministic UID/canonical/ID and numeric selector guidance, opt-in confirmation rules for explicitly requested never/always values, safe-default prompt and explicit confirmation rules, execution safety, YAML precedence, agent, init, and upgrade guidance.
documents:
  SKILL.md: RunX-specific command execution and verified native self-upgrade workflow.
tags:
  - agents
  - runx
keywords:
  - runx
  - cobra
  - command catalog
  - dry run
  - Git Bash
  - shell auto
  - runx upgrade
  - recovery install
  - automatic maintenance
flags: []
status: stable
---

The skill makes manifest inspection, stable UID or interactive numeric-index
selection, dry runs, and user-requested prompt-or-flag confirmation the default
agent workflow for RunX commands. On Windows, automatic shell selection follows
verified Git Bash/MSYS caller evidence and rejects the System32/WSL Bash
launcher before falling back to `cmd.exe`. Agents omit `confirm` unless the user asks
for one of the supported values, while existing `confirm: always` commands
retain their explicit authorization boundary. The skill also teaches agents
how to interpret transactional native upgrades and their recovery contract.
Ordinary RunX invocations also reconcile the current global skill and nearest
managed project instruction block through a silent detached worker.
