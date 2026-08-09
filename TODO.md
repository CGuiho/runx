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
| completed | 17 |

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

### 13. Implement RunX Manifest V2 Composition

- Status: completed
- Created: `2026-07-23`
- Updated: `2026-07-23`
- Outcome: RunX colocates recursive commands and groups and composes explicit,
  reciprocal local or GitHub child catalogs under renameable namespace aliases.
- Spec: [docs/todo/done/manifest-v2-composition.md](docs/todo/done/manifest-v2-composition.md)
- Requirements: [docs/requirements/manifest-v2-composition.md](docs/requirements/manifest-v2-composition.md)
- Decision: [docs/decisions/manifest-v2-composition.md](docs/decisions/manifest-v2-composition.md)
- Plan: [docs/plans/manifest-v2-composition.md](docs/plans/manifest-v2-composition.md)
- Review: [docs/reviews/implementation/manifest-v2-composition-review.md](docs/reviews/implementation/manifest-v2-composition-review.md)
- Validation: [docs/validation/manifest-v2-composition.md](docs/validation/manifest-v2-composition.md)
- External: GitHub issue [CGuiho/runx#26](https://github.com/CGuiho/runx/issues/26)

### 14. RunX Complete Go CLI Rewrite

- Status: completed
- Created: `2026-07-24`
- Updated: `2026-07-26`
- Outcome: RunX production, CI, installers, and release publication are owned by Go/Cobra with strict typed YAML, real manifest-backed commands, live Developer Context help, synchronous idempotent bare-command agent bootstrap, bounded detached lifecycle workers, embedded agent resources, checksum-verifying upgrades, and the standard 11-artifact compatibility matrix. Legacy TypeScript remains reference-only.
- Spec: [docs/rfc/runx-go-rewrite-rfc.md](docs/rfc/runx-go-rewrite-rfc.md)
- Validation: [docs/validation/runx-go-migration.md](docs/validation/runx-go-migration.md)
- External: GitHub issue [CGuiho/runx#22](https://github.com/CGuiho/runx/issues/22)

### 15. Release Unreleased RunX Pull Requests

- Status: completed
- Created: `2026-08-03T00:03:50+02:00`
- Updated: `2026-08-03T00:33:38+02:00`
- Outcome: RunX 0.10.0 published the accepted Windows installer and numeric-index selector changes through the complete Mirror-managed 11-artifact release lifecycle, with all remote checks passing and no production action.
- Spec: [docs/todo/done/release-runx-0.10.0.md](docs/todo/done/release-runx-0.10.0.md)
- Related files:
  - [docs/validation/runx-0.10.0-release.md](docs/validation/runx-0.10.0-release.md) - Live release audit, validation, workflow, asset, and smoke evidence.
- External:
  - GitHub PR [CGuiho/runx#29](https://github.com/CGuiho/runx/pull/29) - Fix Windows installer PATH setup.
  - GitHub PR [CGuiho/runx#31](https://github.com/CGuiho/runx/pull/31) - Align list output and resolve numeric indexes.
- Release: [RunX 0.10.0](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.10.0)

### 16. Integrate RunX Pull Request 38

- Status: completed
- Created: `2026-08-08T12:56:02Z`
- Updated: `2026-08-08T12:56:02Z`
- Outcome: PR #38 merged its changelog-only 0.11.0 confirmation note into protected `main` after exact-head implementation review, validation, CI, mergeability, and reachability gates passed.
- Review: [docs/reviews/implementation/runx-0.11.0-confirmation-review.md](docs/reviews/implementation/runx-0.11.0-confirmation-review.md)
- Validation: [docs/validation/runx-0.11.0-confirmation.md](docs/validation/runx-0.11.0-confirmation.md)
- External:
  - GitHub PR [CGuiho/runx#38](https://github.com/CGuiho/runx/pull/38) - Record the interactive confirmation changelog.
  - GitHub issue [CGuiho/runx#37](https://github.com/CGuiho/runx/issues/37) - Closed by the merged PR linkage.
- Integration: merge commit [`3a140e1528c8edbbb5e2d5842c841a3b436df1dd`](https://github.com/CGuiho/runx/commit/3a140e1528c8edbbb5e2d5842c841a3b436df1dd); source branch deleted after main reachability was verified.
- Release: Deferred; no version tag or GitHub Release was created under the user's merge-only authorization.
- Question ledger: No repository question-ledger file exists and no unresolved questions were recorded for this changelog-only integration.

### 17. Fix RunX Issues 36 And 39

- Status: completed
- Created: `2026-08-08T15:00:00+02:00`
- Updated: `2026-08-08T22:50:59Z`
- Outcome: RunX now resolves exact global UIDs before canonical selectors and unique unqualified ID shorthands, preserves ambiguity failures and numeric-index fallback, and emits `.scripts` as the current Go `runx init` default while preserving explicit directories.
- Spec: [docs/todo/issues-36-39.md](docs/todo/issues-36-39.md)
- Implementation: [docs/todo/issues-36-39-implementation.md](docs/todo/issues-36-39-implementation.md)
- Review: [docs/reviews/implementation/runx-issues-36-39-review.md](docs/reviews/implementation/runx-issues-36-39-review.md)
- Validation: [docs/validation/runx-issues-36-39.md](docs/validation/runx-issues-36-39.md)
- Questions: [docs/questions/issues-36-39/plan-unit-1.md](docs/questions/issues-36-39/plan-unit-1.md)
- External:
  - GitHub PR [CGuiho/runx#40](https://github.com/CGuiho/runx/pull/40) - Fix RunX selector identities and init defaults.
  - GitHub issue [CGuiho/runx#36](https://github.com/CGuiho/runx/issues/36) - Closed by the merged PR linkage.
  - GitHub issue [CGuiho/runx#39](https://github.com/CGuiho/runx/issues/39) - Closed by the merged PR linkage.
- Integration: merge commit [`e3aed59680c4eea4d18597c3ccd8754daa77830e`](https://github.com/CGuiho/runx/commit/e3aed59680c4eea4d18597c3ccd8754daa77830e); main reachability was verified before cleanup.
- Release: Deferred; no version tag or GitHub Release was created under the user's merge-only authorization.

### 18. Publish RunX 0.11.0

- Status: completed
- Priority: highest
- Created: `2026-08-09T00:00:00+02:00`
- Updated: `2026-08-09T08:00:00+02:00`
- Outcome: Publish the externally visible interactive-confirmation capability and compatible selector/init corrections as the Mirror-managed 0.11.0 release.
- Spec: [docs/todo/done/release-runx-0.11.0.md](docs/todo/done/release-runx-0.11.0.md)
- Review: [docs/reviews/implementation/runx-0.11.0-release-prep-review.md](docs/reviews/implementation/runx-0.11.0-release-prep-review.md)
- Validation: [docs/validation/runx-0.11.0-release-prep.md](docs/validation/runx-0.11.0-release-prep.md)
- Integration: merge commit [`4ed75572acb32d35be12d36e1e760671300f2733`](https://github.com/CGuiho/runx/commit/4ed75572acb32d35be12d36e1e760671300f2733).
- Release: [`@guiho/runx/v0.11.0`](https://github.com/CGuiho/runx/releases/tag/%40guiho%2Frunx%2Fv0.11.0), sourced from [`e2b86336ebd95bc6bf25d395f518a1dce66132e1`](https://github.com/CGuiho/runx/commit/e2b86336ebd95bc6bf25d395f518a1dce66132e1); publish workflow [31295252704](https://github.com/CGuiho/runx/actions/runs/31295252704) succeeded.
- Acceptance: non-draft, non-prerelease release with exactly 11 assets; independent checksums and GitHub digests matched, downloaded Windows AMD64 version/help/help-tree smoke passed, and npm latest is `0.11.0`.
- Production boundary: no deployment, promotion, traffic, DNS, database, secret, or other production-infrastructure mutation occurred.
