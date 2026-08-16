---
name: RunX GUIHO CLI Convention 0001 Compliance Architecture Review
purpose: Independently verify that the proposed breaking lifecycle architecture is safe, complete, internally consistent, and ready for planning.
description: Reviews artifact integrity, lifecycle authority, launcher concurrency, Windows uninstall, legacy migration, npm retirement, configuration separation, and first-release cutover.
created: 2026-08-16
flags:
  - ready-for-planning
tags:
  - review
  - architecture
  - cli
keywords:
  - GUIHO CLI Convention 0001
  - stable launcher
  - two-phase transaction
  - Windows uninstall
  - protocol-v1 cutover
owner: runx-architecture-reviews
---

# RunX GUIHO CLI Convention 0001 Compliance Architecture Review

## Verdict

Ready for planning.

No blocker or high-severity architecture finding remains after three review and
revision passes. Human identity confirmation remains a mandatory pre-execution
approval gate; it does not block writing or reviewing the implementation plan.

## Reviewed Sources

- [Compliance architecture](../../architecture/guiho-convention-0001-cli-compliance-migration.md)
- [Compliance audit](../guiho-convention-0001-cli-compliance-audit.md)
- GUIHO CLI Convention 0001
- Current RunX repository instructions and prior lifecycle decisions

## Findings And Resolutions

### Resolved — Windows Cobra uninstall handoff was impossible as first drafted

The first draft proposed a supervised helper while both launcher and payload
were locked. The original invocation could not both wait for the helper and
exit so the helper could remove those executables.

The architecture now uses verified Windows rename plus delete-on-close semantics
for Cobra uninstall. The payload owns output and status, the launcher continues
waiting, canonical paths are quarantined and verified absent before success,
and the invoking shell regains control only after both RunX processes exit.
Unsupported filesystems fail before mutation and print the external
`uninstall.ps1` recovery command. The crash table and black-box NTFS tests cover
every handoff boundary.

### Resolved — Artifact integrity was circular and underspecified

The first draft did not define how `artifacts.json` could record its own digest
and size. The final design uses discriminated integrity modes:

- ordinary assets carry SHA-256 and size;
- the manifest self-record is verified externally by `checksums.txt`; and
- the checksum root is intentionally excluded from self-coverage.

The architecture fixes the checksum grammar, validation order, archive-member
digests, and authority at staging, retirement, journal, and committed-ledger
stages.

### Resolved — Lifecycle mutation authority was fragmented

The final architecture makes the verified candidate launcher's non-Cobra
bootstrap mode the single Go filesystem transaction engine. PowerShell and
POSIX scripts remain lifecycle owners and adapters for selection, download,
initial verification, PATH, and reporting. Capability-bound
`prepare`/`activate`/`finalize`/`abort` operations provide a two-phase commit
boundary across Go state, PATH, and legacy-entrypoint changes. Adapter loss,
abort, and stale recovery retain snapshots until joint commit.

### Resolved — Launcher fallback exposed uncommitted or failed payloads

Live journal/lock state now takes precedence over pointer dispatch. Ordinary
invocation accepts only a committed generation; transaction verification needs
the matching token. Fallback activates only an independently verified retained
payload and never places the failed payload in `previous`. Fresh, corrupt,
concurrent, and interrupted states have deterministic behavior and tests.

### Resolved — npm was a second command-entrypoint and lifecycle owner

The final architecture retires `package.json.bin`, the npm downloader, npm
publication, and the npm payload cache. The only supported command entrypoint is
the stable launcher under `$HOME/.guiho/bin/`.

### Resolved — Legacy migration was limited to one assumed path

The installer now observes actual command resolution for applicable Windows,
POSIX, and Git Bash shells, proves legacy ownership by release checksum or
installed ownership record, and fails closed on a foreign/custom shadowing
entrypoint. Missing historical checksums never authorize deletion.

### Resolved — Multi-PR cutover created an installer outage

Additive pre-cutover units leave canonical public lifecycle surfaces unchanged.
The transition PR installs a dual-shape, non-mixing installer so the canonical
main URL continues to install the latest legacy release before protocol v1 is
published. The separately authorized first protocol-v1 release is then remotely
verified. A second reviewed hardening PR removes every legacy branch. Compliance
is claimed only after that hardening integration and main-branch remote smoke.

### Resolved — Exact release count conflicted with future resources

Twenty-five is the exact protocol-v1 set for the confirmed one-skill,
one-prompt, one-definition RunX catalog. The verifier derives its expected set
from a versioned resource contract. A future resource addition changes that
contract, manifest, verification, tests, and documentation together; only
undeclared extras are rejected.

## Identity Approval Gate

Before implementation, the human must explicitly confirm or replace:

- CLI home `runx`;
- main skill `guiho-s-runx`;
- main prompt `guiho-p-runx`;
- instruction asset `guiho-i-runx.md`;
- canonical repository and issue-create URLs; and
- protocol-v1 resource/release contract.

Until that record exists, a plan may be written and reviewed but no
implementation unit is executable.

## Planning Guidance

The implementation plan must preserve:

- additive/dormant pre-cutover units;
- one transition PR plus a separately authorized first protocol-v1 release;
- one post-release hardening PR before the compliance claim;
- capability-bound two-phase lifecycle transactions;
- exact-head independent implementation review and validation for every PR;
- isolated homes for every destructive test; and
- separate Mirror, publication, and production authorization.

## Handoff

```yaml
handoff:
  from: guiho-a-0045-architecture-reviewer
  to: guiho-a-0001-swe
  verdict: ready-for-planning
  next_agent: guiho-a-0046-plan-writer
  implementation_authorized: false
  release_authorized: false
  production_authorized: false
```

## References

- [Compliance architecture](../../architecture/guiho-convention-0001-cli-compliance-migration.md)
- [Compliance audit](../guiho-convention-0001-cli-compliance-audit.md)
- [Prior Windows self-upgrade decision](../../decisions/windows-self-upgrade.md)
- [Prior release-asset decision](../../decisions/markdown-release-assets-and-version-scoped-notes.md)
