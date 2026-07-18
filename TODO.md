---
name: RunX TODO List
purpose: Track package-owned RunX work and link the task specifications that define it.
description: Local task index for the RunX CLI package.
created: 2026-07-12
flags: []
tags:
  - todo
keywords:
  - runx
  - tasks
owner: runx
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Crist&oacute;v&atilde;o GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# GUIHO RunX TODO List

## Parent TODO

- Parent: [../guiho/TODO.md](../guiho/TODO.md)
- Local context: Open-source command-catalog CLI, bundled agent skill, native installers, and package documentation.

## Status Summary

| Status | Count |
| --- | ---: |
| todo | 0 |
| in progress | 1 |
| testing | 0 |
| stopped | 0 |
| completed | 2 |

## Tasks

### 0. Protect RunX Branches and Tag Creation

- Status: in progress
- Priority: highest
- Created: `2026-07-12T02:10:04+02:00`
- Updated: `2026-07-14T15:47:04+02:00`
- Outcome: RunX branch protections and tag creation rules prevent unauthorized or unsafe changes to protected branches and release tags.
- Spec: [docs/todo/protect-branches-and-tag-creation.md](docs/todo/protect-branches-and-tag-creation.md)

### 1. Implement RunX Alpha Command Catalog

- Status: completed
- Created: `2026-07-12T00:40:36+02:00`
- Updated: `2026-07-12T01:02:00+02:00`
- Outcome: RunX provides a documented, language-agnostic `runx.yaml` command catalog with local execution, agent support, installers, CI, and validation evidence.
- Spec: [docs/todo/implement-runx-alpha.md](docs/todo/implement-runx-alpha.md)

### 2. Make RunX Upgrades Reliable and Recoverable

- Status: completed
- Priority: highest
- Created: `2026-07-15`
- Updated: `2026-07-15`
- Outcome: RunX upgrades and direct installers replace and verify the selected version, list every release, and always provide pinned recovery commands.
- Spec: [docs/todo/upgrade-reliability.md](docs/todo/upgrade-reliability.md)
- Related files:
  - [docs/superpowers/specs/2026-07-15-upgrade-reliability-design.md](docs/superpowers/specs/2026-07-15-upgrade-reliability-design.md) - Approved design.
  - [docs/plans/upgrade-reliability-implementation.md](docs/plans/upgrade-reliability-implementation.md) - Executable implementation plan.
  - [docs/validation/upgrade-reliability.md](docs/validation/upgrade-reliability.md) - Full release-ready validation evidence.
- External: GitHub issues `CGuiho/runx#12` and `CGuiho/runx#13`
