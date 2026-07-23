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
| in progress | 0 |
| testing | 2 |
| stopped | 0 |
| completed | 11 |

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

- Status: completed
- Created: `2026-07-19`
- Updated: `2026-07-20`
- Outcome: Ordinary RunX commands non-blockingly reconcile current global skills and one nearest managed `AGENTS.md` block without changing foreground output or exit behavior.
- Spec: [docs/todo/automatic-agent-maintenance.md](docs/todo/automatic-agent-maintenance.md)
- External: GitHub issue [CGuiho/runx#11](https://github.com/CGuiho/runx/issues/11)
- Related files:
  - [docs/plans/automatic-agent-maintenance.md](docs/plans/automatic-agent-maintenance.md) - Approved implementation plan.
  - [docs/reviews/plans/automatic-agent-maintenance-review.md](docs/reviews/plans/automatic-agent-maintenance-review.md) - Ready-for-execution plan review.
- Implementation: [docs/todo/automatic-agent-maintenance-implementation.md](docs/todo/automatic-agent-maintenance-implementation.md)
- Review: [docs/reviews/implementation/automatic-agent-maintenance-review.md](docs/reviews/implementation/automatic-agent-maintenance-review.md)
- Validation: [docs/validation/automatic-agent-maintenance.md](docs/validation/automatic-agent-maintenance.md)

### 3. Complete RunX Upgrade Reliability

- Status: completed
- Created: `2026-07-15`
- Updated: `2026-07-19`
- Outcome: RunX completes synchronous verified replacement, complete release listing, and exact-version recovery for every upgrade outcome.
- Spec: [docs/todo/upgrade-reliability.md](docs/todo/upgrade-reliability.md)
- Plan: [docs/plans/upgrade-reliability-implementation.md](docs/plans/upgrade-reliability-implementation.md)
- Reviews:
  - [docs/reviews/implementation/upgrade-reliability-issue-12-review.md](docs/reviews/implementation/upgrade-reliability-issue-12-review.md)
  - [docs/reviews/implementation/upgrade-reliability-issue-13-review.md](docs/reviews/implementation/upgrade-reliability-issue-13-review.md)
- Validation: [docs/validation/upgrade-reliability.md](docs/validation/upgrade-reliability.md)

### 4. Use Bash For The RunX Installer

- Status: completed
- Created: `2026-07-19`
- Updated: `2026-07-19`
- Outcome: Every canonical Linux/macOS installer and recovery surface invokes Bash, and real Bash tests cover syntax, piping, exact versions, and executable verification.
- Spec: [docs/todo/bash-installer.md](docs/todo/bash-installer.md)
- Review: [docs/reviews/implementation/bash-installer-review.md](docs/reviews/implementation/bash-installer-review.md)
- Validation: [docs/validation/bash-installer.md](docs/validation/bash-installer.md)

### 5. Preserve The RunX Unicode Help Tree

- Status: completed
- Created: `2026-07-19`
- Updated: `2026-07-19`
- Outcome: RunX renders a Unicode, nested, description-aligned command tree and explicitly rejects the legacy ASCII form in regression tests.
- Spec: [docs/todo/unicode-help-tree.md](docs/todo/unicode-help-tree.md)
- Review: [docs/reviews/implementation/unicode-help-tree-review.md](docs/reviews/implementation/unicode-help-tree-review.md)
- Validation: [docs/validation/unicode-help-tree.md](docs/validation/unicode-help-tree.md)

### 6. Resolve The Latest RunX Bash Install

- Status: testing
- Created: `2026-07-20`
- Updated: `2026-07-20`
- Outcome: The Linux/macOS installer resolves latest release assets without parsing a scoped package tag from a redirect URL.
- Spec: [docs/todo/linux-installer-latest-release.md](docs/todo/linux-installer-latest-release.md)
- Plan: [docs/plans/linux-installer-latest-release.md](docs/plans/linux-installer-latest-release.md)
- External: GitHub issue [CGuiho/runx#20](https://github.com/CGuiho/runx/issues/20)

### 7. Use The Runtime Platform In The RunX Greeting

- Status: testing
- Created: `2026-07-20`
- Updated: `2026-07-20`
- Outcome: A no-argument RunX invocation reports Windows, Linux, or macOS according to the runtime platform.
- Spec: [docs/todo/platform-aware-startup-greeting.md](docs/todo/platform-aware-startup-greeting.md)
- Plan: [docs/plans/platform-aware-startup-greeting.md](docs/plans/platform-aware-startup-greeting.md)
- External: GitHub issue [CGuiho/runx#21](https://github.com/CGuiho/runx/issues/21)

### 8. Bound The RunX Update Worker

- Status: completed
- Created: `2026-07-21`
- Updated: `2026-07-21`
- Outcome: RunX schedules at most one finite background update check per cache directory without foreground failures or persistent process accumulation.
- Spec: [docs/todo/bounded-update-worker.md](docs/todo/bounded-update-worker.md)
- Related files:
  - [docs/reviews/implementation/bounded-update-worker-review.md](docs/reviews/implementation/bounded-update-worker-review.md) - Delivery-readiness review of worker coalescing, deadlines, lease ownership, and recovery.
  - [docs/validation/bounded-update-worker.md](docs/validation/bounded-update-worker.md) - Stress, suite, build, asset, XDocs, and release evidence.
- External: Cross-repository incident [CGuiho/xdocs#14](https://github.com/CGuiho/xdocs/issues/14)

### 9. Preserve UTF-8 During Windows Installation

- Status: completed
- Created: `2026-07-21`
- Updated: `2026-07-21`
- Outcome: The PowerShell installer preserves existing UTF-8 project instructions, converges damaged or duplicate RunX blocks, and does not race background reconciliation during version verification.
- Spec: [docs/todo/windows-installer-utf8.md](docs/todo/windows-installer-utf8.md)
- Validation: [docs/validation/windows-installer-utf8.md](docs/validation/windows-installer-utf8.md)

### 10. Implement A Beautiful RunX Welcome Window

- Status: completed
- Created: `2026-07-22`
- Updated: `2026-07-23`
- Outcome: Bare RunX invocation presents a deterministic, platform-aware welcome and an optional validated cached update notice without foreground network work.
- Spec: [docs/todo/beautiful-welcome-window.md](docs/todo/beautiful-welcome-window.md)
- Related files:
  - [docs/requirements/runx-0.6.0-cli-experience.md](docs/requirements/runx-0.6.0-cli-experience.md) - Approved combined release requirements.
  - [docs/plans/runx-0.6.0-cli-experience.md](docs/plans/runx-0.6.0-cli-experience.md) - Approved executable plan.
- External: GitHub issue [CGuiho/runx#23](https://github.com/CGuiho/runx/issues/23)

### 11. Use The Simplified RunX Installation Command

- Status: completed
- Created: `2026-07-22`
- Updated: `2026-07-23`
- Outcome: The public README uses the exact simplified curl bootstrap while installer integrity and verification remain intact.
- Spec: [docs/todo/simplified-install-command.md](docs/todo/simplified-install-command.md)
- Related files:
  - [docs/plans/runx-0.6.0-cli-experience.md](docs/plans/runx-0.6.0-cli-experience.md) - Approved executable plan.
- External: GitHub issue [CGuiho/runx#24](https://github.com/CGuiho/runx/issues/24)

### 12. Forward RunX Child Arguments And Subcommands

- Status: completed
- Created: `2026-07-22`
- Updated: `2026-07-23`
- Outcome: RunX forwards every post-selector child argument losslessly and safely without letting child flags alter RunX routing.
- Spec: [docs/todo/forward-command-arguments.md](docs/todo/forward-command-arguments.md)
- Related files:
  - [docs/decisions/run-argument-ownership.md](docs/decisions/run-argument-ownership.md) - Approved ownership and shell-safety decision.
  - [docs/plans/runx-0.6.0-cli-experience.md](docs/plans/runx-0.6.0-cli-experience.md) - Approved executable plan.
- External: GitHub issue [CGuiho/runx#25](https://github.com/CGuiho/runx/issues/25)
