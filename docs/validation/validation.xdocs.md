---
subject: runx-validation
description: Validation evidence and implementation summaries for RunX releases.
parent: runx-docs
children: []
files:
  alpha-implementation-summary.md: Records the completed alpha implementation, checks, and release boundaries.
  citty-cli-migration.md: Records the complete local validation gate for the Citty command-tree migration.
  interactive-init-manifest.md: Records validation evidence for the interactive initializer and strict manifest contract.
  windows-self-upgrade.md: Records Windows replacement, rollback, cleanup, native build, CI, and XDocs validation evidence.
  upgrade-reliability.md: Tracks the superseded first pass and required revalidation after independent upgrade reliability review corrections.
  rfc-0034-cli-compliance-migration.md: Records passing RFC tests, live upgrade/startup/prompt regressions, builds, bootstrap, installers, import scan, and fourteen-asset evidence.
  automatic-agent-maintenance.md: Records focused and full tests, native builds, exact assets, import scan, XDocs, and Git evidence for issue 11.
documents:
  alpha-implementation-summary.md: Validation summary for the first RunX implementation.
  citty-cli-migration.md: Validation evidence for TypeScript, tests, native assets, npm packaging, CLI behavior, and XDocs.
  interactive-init-manifest.md: Validation evidence for the RunX interactive init manifest feature.
  npm-trusted-publishing-0.2.2.md: Validation evidence for the blocked 0.2.2 npm trusted-publishing retry.
  windows-self-upgrade.md: Validation evidence for GitHub issues #9 and #1 and the Windows self-upgrade patch.
  upgrade-reliability.md: Revalidation record for GitHub issues 12 and 13; currently marked needs-revalidation.
  rfc-0034-cli-compliance-migration.md: Complete RFC 0034 correction validation and patch-release readiness report.
  automatic-agent-maintenance.md: Complete validation report for automatic RunX agent maintenance.
tags:
  - validation
keywords:
  - runx
  - citty
  - runx init
  - scripts directory
  - tests
  - summary
  - windows self-upgrade
  - upgrade reliability
flags: []
status: stable
---

Validation records document exact checks, results, exclusions, and manual
release boundaries for RunX.
