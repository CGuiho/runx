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
| todo | 1 |
| in progress | 0 |
| testing | 3 |
| stopped | 0 |
| completed | 23 |

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

### 19. Make RunX Confirmation Opt-In In The Agent Skill

- Status: completed
- Created: `2026-08-11T08:40:53+02:00`
- Updated: `2026-08-11T10:00:00+02:00`
- Outcome: The RunX agent skill adds manifest confirmation only when the user explicitly requests `confirm: never` or `confirm: always` for that specific command, with omission resolving to `never`.
- Spec: [docs/todo/done/confirm-opt-in-skill.md](docs/todo/done/confirm-opt-in-skill.md)
- Related files:
  - [skills/guiho-s-runx/SKILL.md](skills/guiho-s-runx/SKILL.md) - Canonical bundled agent policy.
  - [embed/skills/guiho-s-runx.SKILL.md](embed/skills/guiho-s-runx.SKILL.md) - Embedded runtime copy.
  - [skills/guiho-s-runx/guiho-s-runx.xdocs.md](skills/guiho-s-runx/guiho-s-runx.xdocs.md) - Skill descriptor semantics.
- Implementation: [docs/todo/done/confirm-opt-in-skill-implementation.md](docs/todo/done/confirm-opt-in-skill-implementation.md)
- Review: [docs/reviews/implementation/runx-0.11.1-confirm-opt-in-review.md](docs/reviews/implementation/runx-0.11.1-confirm-opt-in-review.md)
- Validation: [docs/validation/runx-0.11.1-confirm-opt-in.md](docs/validation/runx-0.11.1-confirm-opt-in.md)
- Integration: merge commit [37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8](https://github.com/CGuiho/runx/commit/37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8) reached `main`; immutable review/validation evidence was materialized by protected docs PR #48 merge [93dc2f201045003b64360e2b214a6f7c86c176e6](https://github.com/CGuiho/runx/commit/93dc2f201045003b64360e2b214a6f7c86c176e6).
- Mirror: Patch `0.11.1` applied with annotated tag object `7d17153b44185f4215079f2e792effee87f36e9f`, peeled source `c84f6432d1e4a1169a356600265d0c583515371b`; final release evidence is recorded in [docs/validation/runx-0.11.1-release.md](docs/validation/runx-0.11.1-release.md).

### 20. Publish RunX 0.11.1

- Status: completed
- Priority: highest
- Created: `2026-08-11`
- Updated: `2026-08-11T10:00:00+02:00`
- Outcome: Publish the compatible confirmation-policy correction as the Mirror-managed `0.11.1` patch release after clean-clone validation, exact 11-asset publication, independent checksums/digests, native Windows AMD64 smoke, and npm latest verification.
- Spec: [docs/todo/done/release-runx-0.11.1.md](docs/todo/done/release-runx-0.11.1.md)
- Review: [docs/reviews/implementation/runx-0.11.1-confirm-opt-in-review.md](docs/reviews/implementation/runx-0.11.1-confirm-opt-in-review.md)
- Validation: [docs/validation/runx-0.11.1-release.md](docs/validation/runx-0.11.1-release.md)
- Mirror: Patch plan `@guiho/runx/v0.11.1` passed; only `mirror version apply patch --yes` was run, producing the annotated tag and successful [Publish workflow 31469968479](https://github.com/CGuiho/runx/actions/runs/31469968479).
- Release: [`@guiho/runx/v0.11.1`](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.11.1) is non-draft/non-prerelease with exactly 11 assets; independent checksums/digests, downloaded Windows AMD64 smoke, and npm latest `0.11.1` passed.
- Production boundary: only authorized repository source/package publication used `production`; no deployment, promotion, traffic, DNS, database, secret, or other infrastructure mutation occurred.

### 21. Add RunX Reveal

- Status: completed
- Created: `2026-08-11T07:01:35Z`
- Updated: `2026-08-11T07:20:00Z`
- Outcome: RunX resolves a catalog command through the same selector identities as `runx run` and prints its exact command for copy-and-paste without executing it.
- Spec: [docs/todo/done/runx-reveal.md](docs/todo/done/runx-reveal.md)
- Related files:
  - [docs/brainstorm/runx-reveal.md](docs/brainstorm/runx-reveal.md) - Captured user intent and current Windows shell-boundary evidence.
  - [docs/plans/runx-reveal.md](docs/plans/runx-reveal.md) - Approved implementation, validation, integration, and minor-release plan.
  - [docs/reviews/plans/runx-reveal-review.md](docs/reviews/plans/runx-reveal-review.md) - Ready-for-execution plan review.
- Implementation: [docs/todo/done/runx-reveal-implementation.md](docs/todo/done/runx-reveal-implementation.md)
- Review: [docs/reviews/implementation/runx-reveal-review.md](docs/reviews/implementation/runx-reveal-review.md)
- Validation: [docs/validation/runx-reveal.md](docs/validation/runx-reveal.md)
- External: GitHub issue [CGuiho/runx#47](https://github.com/CGuiho/runx/issues/47)
- Integration: PR [#50](https://github.com/CGuiho/runx/pull/50) merged as [`cccc87348da130df90665e954f48bad9bd652ceb`](https://github.com/CGuiho/runx/commit/cccc87348da130df90665e954f48bad9bd652ceb); evidence PR [#53](https://github.com/CGuiho/runx/pull/53) merged as [`7baa214c5514c39fdc5af8118925ffde015be5ae`](https://github.com/CGuiho/runx/commit/7baa214c5514c39fdc5af8118925ffde015be5ae).
- Release: [`@guiho/runx/v0.12.0`](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.12.0), successful Publish workflow [31474091754](https://github.com/CGuiho/runx/actions/runs/31474091754), exactly 11 assets.

### 22. Respect Git Bash For Windows Automatic Shell Execution

- Status: completed
- Priority: highest
- Created: `2026-08-11`
- Updated: `2026-08-11`
- Outcome: On Windows, `shell: auto` preserves a verified Git Bash/MSYS caller's path semantics while retaining `cmd.exe` as the safe fallback and keeping every explicit shell authoritative.
- Spec: [docs/todo/done/windows-git-bash-auto-shell.md](docs/todo/done/windows-git-bash-auto-shell.md)
- Plan: [docs/plans/windows-git-bash-auto-shell.md](docs/plans/windows-git-bash-auto-shell.md)
- Implementation: [docs/todo/done/windows-git-bash-auto-shell-implementation.md](docs/todo/done/windows-git-bash-auto-shell-implementation.md)
- External: GitHub issue [CGuiho/runx#47](https://github.com/CGuiho/runx/issues/47)
- Review: [docs/reviews/implementation/runx-windows-git-bash-auto-shell-review.md](docs/reviews/implementation/runx-windows-git-bash-auto-shell-review.md)
- Validation: [docs/validation/runx-0.12.1-shell-auto-release.md](docs/validation/runx-0.12.1-shell-auto-release.md)
- Integration: PR [#55](https://github.com/CGuiho/runx/pull/55) merged as [`a297aa25b94d83d11d361a0fab0b5415bdd1ba20`](https://github.com/CGuiho/runx/commit/a297aa25b94d83d11d361a0fab0b5415bdd1ba20).
- Release: [`@guiho/runx/v0.12.1`](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.12.1), successful Publish workflow [31476981185](https://github.com/CGuiho/runx/actions/runs/31476981185), exactly 11 assets, npm latest `0.12.1`, and downloaded Windows Git Bash smoke passed.

### 23. Make RunX Comply With GUIHO CLI Convention 0001

- Status: testing
- Priority: highest
- Created: `2026-08-16T00:00:00+02:00`
- Updated: `2026-08-22T20:28:23+02:00`
- Execution note: The developer issued a direct, durable execution order ("start everything right now! Do not stop until everything is done") on 2026-08-22. Gates are being satisfied as recorded in `docs/questions/guiho-convention-0001-cli-compliance-migration/EXECUTION.md`; implementation proceeds unit by unit from the clean protected-main base `bc1f67b`.
- Outcome: RunX uses a stable launcher, immutable payloads, manifest-owned transactional lifecycle operations, separate global/project configuration, compliant agent resources, and a release-safe protocol-v1 cutover.
- Spec: [docs/todo/guiho-convention-0001-cli-compliance-migration.md](docs/todo/guiho-convention-0001-cli-compliance-migration.md)
- Related:
  - [Compliance audit](docs/reviews/guiho-convention-0001-cli-compliance-audit.md)
  - [Architecture](docs/architecture/guiho-convention-0001-cli-compliance-migration.md)
  - [Architecture review](docs/reviews/architecture/guiho-convention-0001-cli-compliance-migration-review.md)
  - [Implementation plan](docs/plans/guiho-convention-0001-cli-compliance-migration.md)
  - [Plan review](docs/reviews/plans/guiho-convention-0001-cli-compliance-migration-review.md)

### 24. Fix Windows Installer Pointer Encoding

- Status: completed
- Priority: highest
- Created: `2026-08-23T12:18:30+02:00`
- Updated: `2026-08-23T13:38:20+02:00`
- Outcome: Windows PowerShell installation, rollback, reinstall, different-version whole-release upgrade, and same-version idempotence now pass native and public gates with BOM-free pointers and standard checksum manifests.
- Spec: [docs/todo/done/windows-installer-pointer-encoding.md](docs/todo/done/windows-installer-pointer-encoding.md)
- Review: [docs/reviews/implementation/windows-installer-pointer-encoding-review.md](docs/reviews/implementation/windows-installer-pointer-encoding-review.md)
- Validation: [docs/validation/windows-installer-pointer-encoding.md](docs/validation/windows-installer-pointer-encoding.md)
- Pull requests: [#59](https://github.com/CGuiho/runx/pull/59), [#60](https://github.com/CGuiho/runx/pull/60), [#61](https://github.com/CGuiho/runx/pull/61)
- Release: [runx/v0.14.8](https://github.com/CGuiho/runx/releases/tag/runx%2Fv0.14.8), workflow [32636794191](https://github.com/CGuiho/runx/actions/runs/32636794191).
- Parent: GUIHO CLI Convention 0001 compliance task 23.

### 25. Guarantee Windows Installer Staging Cleanup

- Status: completed
- Priority: highest
- Created: `2026-08-23T13:42:00+02:00`
- Updated: `2026-08-23T13:56:00+02:00`
- Outcome: Every Windows installer outcome removes only its unique RunX staging directory; public failure, recovery install, real upgrade, and zero-leftover checks passed.
- Spec: [docs/todo/done/windows-installer-staging-cleanup.md](docs/todo/done/windows-installer-staging-cleanup.md)
- Pull request: [#62](https://github.com/CGuiho/runx/pull/62)
- Release: [runx/v0.14.9](https://github.com/CGuiho/runx/releases/tag/runx%2Fv0.14.9), workflow [32637716640](https://github.com/CGuiho/runx/actions/runs/32637716640).
- Parent: GUIHO CLI Convention 0001 compliance task 23.

### 26. RX Short Alias and User-Only Commands

- Status: in progress
- Priority: highest
- Created: `2026-08-30T16:15:00Z`
- Updated: `2026-08-30T16:40:00Z`
- Outcome: `rx` lists on bare invocation and delegates `rx <selector>` to `runx run` with version/help parity; per-command `userOnly` guard refuses agent execution without spawn while `runx reveal` remains allowed; both binaries install/upgrade/uninstall transactionally as one minor release.
- Spec: [docs/todo/2026-08-30-rx-and-user-only-commands.md](docs/todo/2026-08-30-rx-and-user-only-commands.md)
- Brainstorm: [docs/brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md](docs/brainstorm/2026-08-30-rx-short-alias-and-user-only-commands.md)
- Requirements: [docs/requirements/2026-08-30-rx-and-user-only-commands.md](docs/requirements/2026-08-30-rx-and-user-only-commands.md)
- Architecture: [docs/architecture/2026-08-30-rx-and-user-only-commands.md](docs/architecture/2026-08-30-rx-and-user-only-commands.md)
- Plan: [docs/plans/2026-08-30-rx-and-user-only-commands.md](docs/plans/2026-08-30-rx-and-user-only-commands.md)
- Validation: `gofmt`, `go vet`, `go test -count=1 ./...`, `go build`, `build-binaries`, `verify-release-assets`, strict XDocs, `mirror config check` + `mirror version plan minor`
