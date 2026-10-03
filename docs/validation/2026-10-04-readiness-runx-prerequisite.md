---
name: RunX Readiness Prerequisite Validation
purpose: Preserve the verified policy and runtime evidence for independent parent review.
description: Meaningful filesystem regressions, native task-binary provenance, process traces, safe invocation semantics and unresolved delivery boundaries.
created: 2026-10-04
flags: [testing, independent-review-pending]
tags: [validation, runx, readiness]
keywords: [issue-63, issue-64, filesystem, lifecycle, strace, task-binary]
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Readiness Prerequisite Validation

## Identity and Scope

- Owner: `CGuiho/runx`, `/root/swe/runx`, `main`.
- Clean initial HEAD and remote-main baseline:
  `4a9f8b9042ed8d2288a790b85c43d6803d0f5693`; no pre-existing changes or held-ahead commits.
- Policy: [issue #63](https://github.com/CGuiho/runx/issues/63),
  [requirements](../todo/native-background-policy.md).
- Separate runtime prerequisite: [issue #64](https://github.com/CGuiho/runx/issues/64),
  [requirements](../todo/inspection-lifecycle-boundary.md).
- Mode: dnd. Parent expressly authorized the minimum scheduling reconciliation;
  the [decision record](../questions/readiness-tool-prerequisite/2026-10-03-inspection-default.md)
  preserves the bounded planning waiver. Parent independent review precedes push
  or any installed-binary replacement.

## Checks and Exact Exit Evidence

| Command / check | Exit | Evidence |
| --- | ---: | --- |
| `go test -count=1 ./cmd -run 'TestInspectionLifecycleFilesystemBoundary\|TestRealRunRetainsLifecycleAndChildExecution' -v` | 0 | `/tmp/opencode/runx-prerequisite-boundary-tests.log` |
| `go test -count=1 -overlay /tmp/opencode/runx-prerequisite-prior-scheduler-overlay.json ./cmd -run 'TestInspectionLifecycleFilesystemBoundary' -v` | 1, expected negative | `/tmp/opencode/runx-prerequisite-negative-prior-scheduler.log` |
| `go test -count=1 ./...` after bundled-resource changes | 0 | `/tmp/opencode/runx-prerequisite-final-go-test.log` |
| `go vet ./...` | 0 | `/tmp/opencode/runx-prerequisite-final-go-vet.log` |
| `gofmt -l cmd/root.go cmd/lifecycle_boundary_test.go` | 0, no output | Static-check JSON below |
| `go build ./...` | 0 | `/tmp/opencode/runx-prerequisite-go-build-all.json` |
| Identified task-payload build below | 0 | `/tmp/opencode/runx-prerequisite-build.json` |
| `python3 /tmp/opencode/runx-prerequisite-binary-verify.py` | 0 | `/tmp/opencode/runx-prerequisite-binary-verification.json` |
| `python3 /tmp/opencode/runx-prerequisite-static-check.py` | 0 | `/tmp/opencode/runx-prerequisite-static-check.json` |

The static check verifies current native policy, rejection of baseline policy,
byte-identical retained historical body, owned Markdown copyright/frontmatter and
new-record links, both skill metadata versions and maintenance-contract parity,
accurate TODO counts, unchanged catalog/config/version files and scoped diff checks.
Its initial phrase assertion expected a different future-policy wording; correcting
the checker to the actual accepted `if CG later authorizes` clause passed without
a repository change. This is checker correction, not a runtime failure.

## Meaningful Filesystem and Execution Oracles

The committed Go suite has **28 inspection cases**: existing versus missing
resources; text/JSON check/list/describe; reveal; parsed dry-run flags; gated and
missing-selector errors; help/version/docs. Every case asserts no scheduling or
worker-executable resolution and identical byte/file-set/mode/mtime snapshots.

At the actual scheduling seam, a worker sentinel writes synchronously instead
of returning from a no-op mock. Maintenance executes the real resource
reconciliation against isolated project/home fixtures; the update sentinel
replaces the isolated cache. The baseline-source overlay substitutes only the
old `cmd/root.go`, leaves the checkout untouched, and **fails** by detecting
worker scheduling, resource creation/rewrites and cache changes.

Three real-run positive controls execute the child sentinel and both workers:
ordinary run, `--dry-run=false` before the selector, and `--dry-run` after it.
They prove the snapshot and execution oracle can detect actual mutations and
that real execution retains lifecycle behavior. Existing full-suite tests also
cover explicit setup/agent actions and child exit codes.

The **actual task binary** passed **38 isolated-fixture invocations** plus
**five owning-catalog inspections**. Both worker toggles were `0` (enabled), not
used to hide scheduling. Linux `strace -f -e trace=process` recorded exactly one
`execve` per invocation: the task binary itself. No catalog child or lifecycle
worker executed. Before/after snapshots include hashes, file/directory sets,
modes and nanosecond mtimes. Stale CRLF AGENTS/CLAUDE blocks, configs, ordinary
docs, invalid cache, unrelated agent files, stale global skills and custom extras
remained identical. Missing resources remained missing.

Fixtures, snapshots, stdout/stderr and traces are retained under
`/tmp/opencode/runx-prerequisite-binary-evidence/`; exact fixture roots are in
the verification JSON. No sleeps or worker polling were used. Access times are
not the snapshot contract. Foreign HTTPS catalog resolution was not exercised
in this local fixture; it remains a documented read/network operation.

The real RunX catalog validated **six commands**. Check/list/describe/reveal and
dry-run of `runx-tooling-go-test` preserved isolated-home snapshots, guarded owner
file hashes and before/after tracked/untracked Git state. No catalog entry was
executed during this owner-catalog inspection.

## Task Binary Provenance

- Payload: `/tmp/opencode/runx-prerequisite-task/runx` (not installed).
- Clean source/build commit: `d9f1675b35b757b30cebbba5ef00eb7778a881aa`.
- Raw version: `0.0.0-dev.readiness.runx.d9f1675`.
- Build timestamp: `2026-10-03T22:02:32Z`.
- Build target: `runx-linux-amd64`; Go `go1.27.1`, `CGO_ENABLED=0`, AMD64 V1.
- SHA-256: `cd76c511fba3c3751ea535c2f77b2c4bdb5ff97ae1fdbca3cbe6def6ffdeef5a`.

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1 go build -trimpath -buildvcs=true -ldflags '-s -w -X main.version=0.0.0-dev.readiness.runx.d9f1675 -X main.commit=d9f1675b35b757b30cebbba5ef00eb7778a881aa -X main.buildDate=2026-10-03T22:02:32Z -X main.buildTarget=runx-linux-amd64' -o /tmp/opencode/runx-prerequisite-task/runx .
```

`go version -m` records the same VCS revision with `vcs.modified=false` in
`/tmp/opencode/runx-prerequisite-build-info.log`. Raw version and hidden
`__self-test` succeeded; self-test reports protocol `1`, coherent version/target
and readable embedded resources. Later commits contain only non-embedded docs
and task mirrors; static validation compares all executable/embedded source
inputs against the clean built commit rather than mislabeling the payload as a
build of the final documentation HEAD.

Existing PATH command `/usr/local/bin/runx` was not executed or replaced. Its
before/after SHA-256 remains
`35907a58fee5619bb240ff49c28187412fd9c89bfe7e6edc3fbd8c4eaa9fc252`.
Its installed source equivalence and safe-boundary rollout are not certified.

## Tracking Readback and Remaining Boundaries

At `2026-10-03T22:03:14Z`, both real issues were verified **OPEN**, on
Project #2 (`GUIHO`, `PVT_kwHOBUk1ds4AULOg`), Component **runx**, Status
**Testing**. Policy item: `PVTI_lAHOBUk1ds4AULOgzg-XUSw`; runtime item:
`PVTI_lAHOBUk1ds4AULOgzg-Xbjw`. Component option `580160ad`; Testing option
`f874e9fb`. Local task/spec mirrors agree. A pre-transition GitHub read returned
HTTP 503 before any mutation; retry after independent build work succeeded.
Raw/summary evidence: `/tmp/opencode/runx-prerequisite-testing-readback.json`
and `/tmp/opencode/runx-prerequisite-testing-summary.json`.

- XDocs CLI validation is **skipped**, because installed-build safety is
  independently owned and not accepted by this unit. `xdocs.yaml` has no
  descriptor/frontmatter grants, and no named descriptor edits were authorized.
  Configuration, descriptors and managed XDocs blocks were left untouched;
  descriptor synchronization remains a visible readiness gap.
- Platform-native smoke evidence is Linux AMD64 only; no Windows/macOS native
  runtime or installer transaction is claimed.
- No Mirror/package version change, changelog, tag, publication, release matrix,
  installed-binary replacement, global skill installation or family certification
  was performed. Skill artifact metadata alone moved compatibly from `0.4.1` to
  `0.4.2` under Convention 0002; this is not a Mirror-managed project bump.
- Current policy/static checks do not prove future MiMo availability or native
  worker permissions. This leaf used the task-authorized model and actual tools,
  with no recursive delegation or provider fallback.
- Historical task bindings and stale generated entry blocks remain separately
  owned readiness work; their statuses and bodies were preserved.
