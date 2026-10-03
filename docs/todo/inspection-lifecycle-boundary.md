---
name: RunX Inspection Lifecycle Boundary
purpose: Define the separately authorized runtime prerequisite tracked by RunX issue 64.
description: Default non-mutating catalog inspection and dry-run scheduling contract, fixture evidence and build-review gate.
created: 2026-10-03
flags: [testing]
tags: [runx, runtime, readiness]
keywords: [issue-64, scheduling, filesystem, dry-run]
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Inspection Lifecycle Boundary

## Todo Index

- Index: [TODO.md](../../TODO.md), task 28.
- Status: testing
- GitHub project item: [RunX #64](https://github.com/CGuiho/runx/issues/64)
- Project: [GUIHO #2](https://github.com/users/CGuiho/projects/2)
- GitHub component: `runx`
- Separate native policy: [RunX #63](https://github.com/CGuiho/runx/issues/63)

## Outcome and Acceptance

`check`, `list`, `describe`, `reveal` and `run --dry-run` do not schedule an
update or agent-maintenance worker, execute a catalog child, or mutate project
instructions/configuration or global skill/agent projections. Guard scheduling
itself rather than rely on `RUNX_DISABLE_AGENT_MAINTENANCE_WORKER`, which formerly
only skipped bare-root foreground maintenance. Preserve real `run` execution,
child exit codes, opt-in confirmation and explicitly authorized setup/agent actions.
RunX-owned run options precede the selector; child flags are not reinterpreted.

Isolated project and user-home byte/file-set snapshots protect stale AGENTS,
CLAUDE, ordinary documents, configs, cache and global resources. Executable child
and side-effect-worker sentinels prove the oracle can observe prior scheduling
and real execution. Relevant Go tests/vet/format/build pass. An explicitly
identified task binary supplies source commit, build command/path, SHA-256,
version and target; installation waits for new independent review.

## Scope Waiver and Decisions

Mode: dnd. Parent expressly authorized the minimum RunX source scheduling fix
as a separately tracked readiness prerequisite in
`/tmp/opencode/2026-10-03-readiness-tool-prerequisite-contract.md`, not implicit
source work inside a metadata unit. Default inspection exclusion is the smallest
reversible interface: no new global opt-out, configuration schema or flag.
Only the parsed `run --dry-run` flag selects this boundary; `--dry-run=false`
and post-selector child flags retain existing real-run behavior. The
[question record](../questions/readiness-tool-prerequisite/2026-10-03-inspection-default.md)
records the decision and later human-review gate.

## Guardrails and Review

Native leaf General, task-authorized `openai/gpt-6.1-sol#xhigh`; no recursive
delegation, CLI worker, secret inspection, host provider/permission changes,
catalog/application refactor, publication or Mirror bump. Own commits only on
main; no push before parent independent review. Keep issue OPEN/Testing after
validation and mirror fresh Project/Component/Status readback. This unit does not
certify the tool installation or any family. XDocs data validation remains
skipped pending its independent runtime reconciliation; project YAML grants no
automatic descriptor/frontmatter writes.

Fresh GitHub readback at `2026-10-03T22:30:09Z` confirms issue #64 remains
OPEN on Project #2 (`GUIHO`), Component `runx`, Status `Testing`. A new independent
native General reviewer accepted the scheduler boundary, repeated the meaningful
old-source failure and actual filesystem/process checks, and plain-pushed the
implementation through `dccc8fe` under express parent authority. The checked
`0.0.0-dev.readiness.runx.dccc8fe` build is active at existing `/usr/local/bin/runx`;
old binary and both exact existing skill files are backed up. Installed-binary
inspection passes 43 invocations with identical snapshots and no worker/child or
mutating filesystem syscall. Three real-run controls retain both workers and child
argument forwarding. The implementer's held gates are satisfied; human acceptance
and separately owned XDocs/full-readiness setup remain pending.

## Evidence and Handoff

- [Technical self-review](../reviews/implementation/2026-10-04-readiness-runx-prerequisite-review.md) — acceptance mapping and scoped findings.
- [Validation record](../validation/2026-10-04-readiness-runx-prerequisite.md) — meaningful negative/positive regressions, payload filesystem/process evidence and provenance.
- [Independent-review handoff](../reviews/handoff/2026-10-04-readiness-runx-prerequisite.md) — exact safe commands, pending delivery and installed-build boundaries.
- [Independent acceptance and activation](../reviews/implementation/2026-10-04-readiness-runx-independent-review.md) — checked installed identity, hashes/backups, real runtime proof, reviewed delivery and remaining gaps.
