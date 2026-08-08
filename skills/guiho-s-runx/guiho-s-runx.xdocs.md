---
subject: guiho-s-runx
description: Agent workflow for safe RunX catalog work, deterministic UID/canonical/ID selector resolution, interactive numeric-index selection, explicit confirmation, bare-command bootstrap, automatic resource maintenance, and verified native upgrade/list/recovery operations.
parent: runx-skills
children: []
files:
  SKILL.md: RFC-aligned RunX catalog inspection, deterministic UID/canonical/ID and numeric selector guidance, safe-default prompt and explicit confirmation rules, execution safety, YAML precedence, agent, init, and upgrade guidance.
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
  - runx upgrade
  - recovery install
  - automatic maintenance
flags: []
status: stable
---

The skill makes manifest inspection, stable UID or interactive numeric-index
selection, dry runs, and explicit prompt-or-flag confirmation the default agent workflow for
RunX commands and teaches agents how to interpret transactional native upgrades
and their recovery contract.
Ordinary RunX invocations also reconcile the current global skill and nearest
managed project instruction block through a silent detached worker.
