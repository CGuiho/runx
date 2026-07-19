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

- Status: completed
- Priority: highest
- Created: `2026-07-12T02:10:04+02:00`
- Updated: `2026-07-14T15:47:04+02:00`
- Outcome: RunX branch protections and tag creation rules prevent unauthorized or unsafe changes to protected branches and release tags.
- Spec: [docs/todo/protect-branches-and-tag-creation.md](docs/todo/protect-branches-and-tag-creation.md)

### 1. Migrate RunX To Full RFC 0034 Compliance

- Status: completed
- Created: `2026-07-18T18:48:11+02:00`
- Updated: `2026-07-18`
- Outcome: RunX fully implements the breaking GUIHO RFC 0034 CLI contract, including the independently audited upgrade routing, cached no-argument notice, raw prompt-name output, `.md` agent release assets, version-scoped GitHub Release notes, installers, npm distribution, tests, and documentation.
- Spec: [docs/todo/rfc-0034-cli-compliance-migration.md](docs/todo/rfc-0034-cli-compliance-migration.md)
- Related files:
  - [docs/plans/rfc-0034-cli-compliance-migration.md](docs/plans/rfc-0034-cli-compliance-migration.md) - Approved step-by-step migration plan.
  - [docs/reviews/plans/rfc-0034-cli-compliance-migration-review.md](docs/reviews/plans/rfc-0034-cli-compliance-migration-review.md) - Ready-for-execution plan review.
- Implementation: [docs/todo/rfc-0034-cli-compliance-migration-implementation.md](docs/todo/rfc-0034-cli-compliance-migration-implementation.md)

### 2. Add Automatic RunX Agent Maintenance

- Status: in progress
- Created: `2026-07-19`
- Updated: `2026-07-19`
- Outcome: Ordinary RunX commands non-blockingly reconcile current global skills and one nearest managed `AGENTS.md` block without changing foreground output or exit behavior.
- Spec: [docs/todo/automatic-agent-maintenance.md](docs/todo/automatic-agent-maintenance.md)
- External: GitHub issue [CGuiho/runx#11](https://github.com/CGuiho/runx/issues/11)
- Related files:
  - [docs/plans/automatic-agent-maintenance.md](docs/plans/automatic-agent-maintenance.md) - Approved implementation plan.
  - [docs/reviews/plans/automatic-agent-maintenance-review.md](docs/reviews/plans/automatic-agent-maintenance-review.md) - Ready-for-execution plan review.
- Implementation: [docs/todo/automatic-agent-maintenance-implementation.md](docs/todo/automatic-agent-maintenance-implementation.md)
- Review: [docs/reviews/implementation/automatic-agent-maintenance-review.md](docs/reviews/implementation/automatic-agent-maintenance-review.md)
- Validation: [docs/validation/automatic-agent-maintenance.md](docs/validation/automatic-agent-maintenance.md)
