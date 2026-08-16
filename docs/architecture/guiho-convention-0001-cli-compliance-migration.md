---
name: RunX GUIHO CLI Convention 0001 Compliance Architecture
purpose: Define the target architecture that replaces RunX's direct-binary lifecycle with convention-compliant installation, configuration, agent, upgrade, and release systems.
description: Specifies the stable launcher protocol, immutable payload store, artifact ownership model, configuration inheritance, agent resources, lifecycle transactions, migration boundaries, and verification seams required by GUIHO CLI Convention 0001.
created: 2026-08-16
flags:
  - proposed
  - breaking-change
tags:
  - architecture
  - cli
  - compliance
keywords:
  - GUIHO CLI Convention 0001
  - stable launcher
  - current.json
  - artifacts.json
  - immutable payloads
  - runx.global.yaml
  - agent evolution
  - synchronous upgrade
owner: runx-architecture
---

# RunX GUIHO CLI Convention 0001 Compliance Architecture

## Status

Proposed architecture for human approval and independent architecture review.
It is planning input, not implementation authorization.

This architecture deliberately breaks the superseded direct-executable,
single-configuration, eleven-asset, background-replacement, and agent `update`
contracts. RunX is pre-1.0, and the repository instructions permit breaking
changes when an approved migration requires them.

## Authority And Inputs

The authority order for this migration is:

1. [GUIHO CLI Convention 0001](../../../guiho/docs/conventions/guiho-convention-0001-cli.md).
2. Current RunX repository and parent instructions.
3. [Convention compliance audit](../reviews/guiho-convention-0001-cli-compliance-audit.md).
4. This architecture after human approval.
5. Earlier requirements, decisions, plans, and implementation records only
   where they do not conflict with items 1 through 4.

The audit found 4 critical, 16 high, and 4 medium violation groups. The target
architecture closes all of them as one lifecycle design rather than as isolated
patches.

## Decision Status

The following names are recommended from current product identity and existing
repository history, but Convention 0001 requires human confirmation. Approval
of this architecture must explicitly confirm or replace them before source
implementation starts.

| Decision | Recommended value | Status |
| --- | --- | --- |
| CLI home-directory name | `runx` | Proposed; human confirmation required |
| Main agent skill ID | `guiho-s-runx` | Proposed; human confirmation required |
| Main setup prompt ID | `guiho-p-runx` | Proposed; human confirmation required |
| Managed instruction asset | `guiho-i-runx.md` | Proposed; preserves the accepted instruction filename |
| Canonical repository | `https://github.com/CGuiho/runx` | Observed repository identity |
| Canonical issue-create URL | `https://github.com/CGuiho/runx/issues/new` | Derived from observed repository identity |

The accepted decision
[`markdown-release-assets-and-version-scoped-notes.md`](../decisions/markdown-release-assets-and-version-scoped-notes.md)
continues to govern version-scoped release notes and the `.md` suffix for the
instruction. Its exact fourteen-asset decision is superseded by the convention's
complete artifact requirement. The existing Windows self-replacement decision
is superseded by stable-launcher activation; only its principle that success is
synchronous remains.

## Architectural Goals

1. One selected release is a complete, checksummed installation unit.
2. The command on `PATH` is a stable launcher, never the running application
   payload.
3. Payload versions and canonical release resources are immutable after
   verification.
4. Install, reinstall, upgrade, repair, and uninstall act only on
   manifest-proven RunX ownership.
5. Activation is one atomic pointer replacement with an immediately usable
   verified fallback.
6. Every mutating lifecycle operation is synchronous, journaled, recoverable,
   and testable inside an isolated home.
7. Project and global configuration are different contracts with explicit
   inheritance.
8. `runx init` is the sole public project reconciliation entrypoint.
9. A single canonical source exists for every bundled skill, prompt,
   definition, and instruction.
10. Release selection implements complete SemVer ordering and validates the
    complete compatible artifact set before mutation.
11. Mirror, RunX, and XDocs describe and validate the same supported project
    truth.

## Non-Goals

- Preserve direct installation in `$HOME/.local/bin/`.
- Preserve the independent npm payload downloader/cache.
- Preserve `agent skill update` or `agent instruction update` aliases.
- Preserve `$HOME/.guiho/runx/runx.yaml` as a fallback project catalog.
- Add a public or hidden Cobra `install` command.
- Auto-publish, tag, modify the live user installation, or mutate production as
  part of implementation or validation.
- Add signatures or attestations not required by Convention 0001. SHA-256
  verification remains mandatory and complete.

## System Context

```text
GitHub Release catalog
  -> release resolver
  -> checksums.txt + artifacts.json
  -> platform launcher + payload + common resources
  -> isolated staging under $HOME/.guiho/.temp/
  -> transaction and ownership engine
  -> immutable version directory
  -> managed agent projections and AGENTS.md block
  -> atomic current.json activation
  -> stable launcher verification

Shell or npm proxy
  -> $HOME/.guiho/bin/runx[.exe] stable launcher
  -> strict current.json
  -> immutable active payload
  -> original args/stdin/stdout/stderr
  -> exact payload exit code
```

The external installer scripts own fresh installation and installer-driven
reinstallation. They download and verify a candidate launcher, then invoke its
capability-token-protected bootstrap mode. That non-Cobra mode is the single Go
transaction engine for canonical files, projection snapshots, journaling,
activation, verification, and rollback. It is not reachable as `runx install`,
does not appear in help, and refuses any plan or capability that was not created
inside the current RunX-owned staging directory. Scripts remain the lifecycle
owners: they resolve the release, verify the bootstrap inputs, authorize the
one-time capability, coordinate PATH, and report the result.

The Go payload uses the same transaction package for `upgrade`, repair, agent
projection, and Cobra uninstall planning. All paths consume the same versioned
formats and golden conformance fixtures; no implementation may infer ownership
from a filename prefix alone.

## Target Source Boundaries

The current single Cobra tree remains. New and revised packages have these
responsibilities:

| Path | Responsibility |
| --- | --- |
| `main.go` | Thin payload entrypoint with SemVer-compatible defaults and injected release resources |
| `cmd/runx-launcher/main.go` | Separate stable-launcher entrypoint and capability-protected installer bootstrap host |
| `cmd/` | One public Cobra tree, lifecycle command adapters, help rendering, and no domain filesystem logic |
| `pkg/config/` | Separate global/project types, strict decoding, merge semantics, policy resolution, schema validation, and init question model |
| `pkg/manifest/` | RunX command-catalog model, strict decoding, semantic validation, composition, and execution-safe inspection |
| `pkg/release/` | Strict SemVer, channel selection, complete catalog pagination, release compatibility, checksum, and `artifacts.json` decoding |
| `pkg/installstate/` | Native paths, `current.json`, installed manifest, locks, instances, journals, persistence classes, safe descendant checks, and atomic files |
| `pkg/launcher/` | Strict active/previous selection, fallback, process delegation, and exact exit-code forwarding |
| `pkg/lifecycle/` | Single Go transaction state machine for launcher bootstrap installation, upgrade, repair, projection reconciliation, rollback, and uninstall planning/execution |
| `pkg/agent/` | Canonical resource catalog, global projections, bounded instruction operations, and raw list/show behavior |
| `pkg/updater/` | Thin compatibility facade removed or reduced to the new release/lifecycle packages |
| `resources_embed.go` | Root-package `go:embed` of the canonical top-level resource sources |
| `skills/guiho-s-runx/` | One canonical main skill plus its agent definition |
| `prompts/guiho-p-runx.md` | One canonical main setup prompt |
| `instructions/guiho-i-runx.md` | One canonical managed instruction source |
| `schemas/` | Project/global JSON Schemas |
| `examples/` | Complete valid project/global examples |
| `devops/install.*` | External bootstrap, selection, download, verification, transaction orchestration, PATH, and legacy migration |
| `devops/uninstall.*` | External uninstallation with the same ownership and preservation contract as Cobra |

`resources_embed.go` embeds the top-level canonical files directly. The
duplicate `embed/skills` and `embed/prompts` copies are removed. Release
construction and runtime embedding therefore consume the same bytes.

## Installed Filesystem Model

With the proposed `runx` home name, the normative layout is:

```text
$HOME/.guiho/
├── bin/
│   └── runx[.exe]                         stable launcher
├── .temp/
│   └── runx-<operation>-<unique-id>/      one owned transaction only
└── runx/
    ├── current.json                       active and previous pointer
    ├── installed-artifacts.json           active ownership/projection ledger
    ├── runx.global.yaml                   persistent global configuration
    ├── versions/
    │   └── <semver>/
    │       ├── runx[.exe]                 immutable payload
    │       ├── artifacts.json             immutable selected release manifest
    │       └── artifacts/
    │           ├── skills/
    │           ├── prompts/
    │           ├── instructions/
    │           ├── schemas/
    │           ├── examples/
    │           └── agent-definitions/
    ├── state/
    │   ├── upgrade.lock
    │   ├── transactions/<id>.json
    │   └── instances/<pid>-<token>.json
    ├── cache/                              disposable
    └── data/                               persistent when introduced/declared
```

The RunX home, shared `bin`, and shared `.temp` paths are resolved with
`os.UserHomeDir` and `filepath.Join` in Go. Scripts use the platform-native home
API and normalize/validate every destructive target.

### Persistence Classification

| Path or projection | Class | Reinstall/upgrade | Default uninstall |
| --- | --- | --- | --- |
| `runx.global.yaml` | Persistent configuration | Preserve | Remove unless `--preserve-config` |
| project `runx.yaml` | Persistent configuration and catalog | Preserve | Remove unless `--preserve-config` |
| `data/` and declared databases | Persistent data | Preserve | Remove unless `--preserve-data` |
| launcher, payloads, canonical resources | Replaceable | Replace/retire by manifest | Remove |
| global skill projections | Replaceable projection | Reconcile | Remove |
| managed `AGENTS.md` block | Replaceable projection | Reconcile | Remove block only |
| cache, instance files, completed journals | Disposable | Remove | Remove |
| shared `.guiho`, `bin`, `.temp`, PATH entry | Shared infrastructure | Preserve | Preserve |

Future persistent paths must be added to both schemas and the artifact/installed
ownership contract before code writes user-created content there.

## Release Artifact Contract

Every RunX release publishes exactly these 25 assets:

### Platform payloads

1. `runx-payload-linux-amd64`
2. `runx-payload-linux-arm64`
3. `runx-payload-linux-armv7`
4. `runx-payload-linux-armv6`
5. `runx-payload-darwin-amd64`
6. `runx-payload-darwin-arm64`
7. `runx-payload-windows-amd64.exe`
8. `runx-payload-windows-arm64.exe`

### Platform launchers

9. `runx-launcher-linux-amd64`
10. `runx-launcher-linux-arm64`
11. `runx-launcher-linux-armv7`
12. `runx-launcher-linux-armv6`
13. `runx-launcher-darwin-amd64`
14. `runx-launcher-darwin-arm64`
15. `runx-launcher-windows-amd64.exe`
16. `runx-launcher-windows-arm64.exe`

### Common resources and metadata

17. `guiho-s-runx.zip`
18. `guiho-p-runx.md`
19. `guiho-i-runx.md`
20. `runx.schema.json`
21. `runx.global.schema.json`
22. `runx.example.yaml`
23. `runx.global.example.yaml`
24. `artifacts.json`
25. `checksums.txt`

The `runx-payload-*` names intentionally differ from the legacy direct-binary
names. An old installer or old self-upgrader therefore fails release
compatibility instead of placing a new payload directly on `PATH`.

`guiho-s-runx.zip` contains `SKILL.md` and
`agents/openai.yaml`. `artifacts.json` enumerates the archive and every contained
path separately, including its content digest, artifact ID, canonical installed
path, managed projections, and ownership boundary. The release-level checksum
file covers every asset except itself, including `artifacts.json`.

The release verifier requires all 25 unique assets and rejects missing,
duplicate, extra, malformed, or legacy assets. A local installation downloads
the selected platform payload and launcher plus every common asset. The release
resolver still verifies that the release as published contains the complete
eight-target payload and launcher matrix before selecting it.

Twenty-five is the exact RunX convention-protocol-v1 set for the approved
one-skill, one-prompt, one-definition resource catalog. The verifier derives the
expected set from a versioned repository release-contract source rather than a
hard-coded numeric constant. Adding a future skill, prompt, agent definition,
schema, example, or permitted metadata asset requires one reviewed contract
change that updates the canonical resource catalog, generated manifest, release
verifier, installers, tests, and documentation. Undeclared extras are rejected;
declared additions change the expected count intentionally.

## `artifacts.json` Contract

The strict manifest has:

- schema/protocol version;
- release SemVer, Git tag, source commit, and build date;
- launcher protocol compatibility range;
- complete published asset records;
- target applicability (`goos`, `goarch`, tuning, executable format);
- asset integrity mode, checksum reference, and size;
- archive-member records and member digests;
- canonical relative installed paths;
- managed projection destinations and projection kind;
- ownership class (`replaceable`, `persistent-config`, `persistent-data`, or
  `disposable`);
- retention/retirement behavior;
- self-test requirements; and
- required-versus-optional status.

Every installed path is relative to an explicit ownership root. Strict
validation rejects absolute paths, empty segments, `.` or `..`, alternate path
separators that escape a root, duplicate IDs, duplicate canonical paths,
overlapping destructive roots, foreign `.guiho` locations, and projections
outside the supported global skill directories or bounded `AGENTS.md` block.

The manifest includes records for itself and `checksums.txt` through a
discriminated integrity union:

- ordinary assets use `sha256`, with a lowercase 64-hex digest and exact byte
  size in `artifacts.json`;
- the `artifacts.json` self-record uses `external-sha256`, omits digest and size,
  and names `checksums.txt` as its integrity authority; and
- the `checksums.txt` record uses `checksum-root`, omits its own digest, and is
  explicitly excluded from checksum coverage.

`checksums.txt` is UTF-8 without BOM, LF-terminated, and uses exactly
`<64 lowercase hex><two ASCII spaces><asset filename>`. Filenames are NFC,
repository-declared root asset names without separators, control characters, or
leading/trailing whitespace. The parser rejects duplicate names, malformed
lines, missing or extra undeclared names, and any checksum entry for
`checksums.txt`.

Validation order is fixed:

1. download `checksums.txt` and `artifacts.json`;
2. parse the checksum grammar and reject duplicates/malformed names;
3. verify the downloaded `artifacts.json` bytes from its checksum entry;
4. strictly decode and semantically validate the manifest;
5. derive the required target-specific plus common download set;
6. verify every downloaded asset's bytes and size;
7. extract archives without traversal; and
8. verify each archive member's SHA-256 over its uncompressed bytes against its
   manifest member record.

Authority changes by transaction stage:

- the verified release manifest plus checksum root authorizes candidate source
  bytes and intended new paths;
- the previous installed ledger authorizes retirement or rollback of existing
  paths, always intersected with hard allowed roots;
- the transaction journal owns only paths it created or snapshotted during the
  current operation; and
- after commit, `installed-artifacts.json` is the active normalized ownership
  ledger and records the immutable release-manifest digest that produced it.

`installed-artifacts.json` records the actual canonical and projection paths
installed by the completed transaction. It is atomically replaced only after
verification succeeds. Destructive plans require the intersection of the
installed ledger, the referenced immutable manifest's ownership class, and the
operation's hard allowed roots; no single source alone is sufficient authority.

## Stable Launcher Protocol

The launcher is a separate small Go program without Cobra or RunX catalog
domain logic. It:

1. resolves the native RunX home;
2. strictly decodes `current.json`;
3. validates active and previous payload paths as relative descendants of
   `versions/`;
4. registers the selected payload instance through inherited metadata;
5. starts the active payload with unchanged arguments, environment, standard
   streams, and working directory;
6. waits for completion and returns the exact exit code; and
7. falls back only when the active payload is missing, invalid, or cannot be
   started, never because the payload command returned a domain error.

`current.json` contains a protocol version, an active record, and a nullable
previous record. Fresh installation writes `"previous": null`:

```json
{
  "protocolVersion": 1,
  "active": {
    "version": "1.2.3",
    "payload": "versions/1.2.3/runx.exe"
  },
  "previous": {
    "version": "1.2.2",
    "payload": "versions/1.2.2/runx.exe"
  }
}
```

The platform-specific extension is omitted on macOS/Linux. The pointer is
written to a same-directory unique file, flushed, and atomically renamed. If
active startup fails and previous startup succeeds, the launcher repairs the
pointer under the lifecycle lock so the verified fallback becomes active. It
never falls back on a nonzero child exit.

Launcher dispatch and recovery follow this transition table:

| State | Lock state | Deterministic action |
| --- | --- | --- |
| Live transaction journal/lock exists | Live owner | Give the journal precedence; dispatch only the last committed generation, or fail if fresh install has none |
| Active valid/startable and generation committed | No live transaction | Dispatch active; do not inspect or mutate previous |
| Active unstartable, previous valid/startable | Lock free | Acquire lock, re-read pointer, activate the verified fallback only if unchanged, set previous to another independently verified retained payload or null, then dispatch |
| Active unstartable, previous valid/startable | Live lock held | Bounded wait for the lock and re-read; on timeout dispatch previous without pointer mutation and emit one local repair diagnostic to stderr |
| Active and previous invalid/unstartable | Any | Fail with exit 5 and exact external reinstall command; never guess a version |
| Incomplete journal before ordinary read-only/domain dispatch | Lock free | Run local non-network repair limited to pointer/journal consistency, then dispatch |
| Incomplete journal with live owner | Live lock held | Bounded wait; after timeout fail with exit 5 rather than race activation |
| Pointer changes during fallback acquisition | Lock acquired after change | Discard stale decision, re-read, and restart selection once |

Installer activation, upgrade activation, and fallback repair use the same
exclusive lock and compare the pointer generation/digest before replacement.
Launcher protocol compatibility is checked three ways: the release manifest
declares its supported launcher protocol range, the payload embeds the same
range, and the launcher rejects an active record outside its own supported
range. A fresh install has no fallback; a corrupt fresh pointer fails closed.

Every pointer generation has one explicit state: `prepared`,
`active-uncommitted`, `committed`, or `rolled-back`. Ordinary launcher
dispatch accepts only `committed`. Transaction verification may dispatch an
`active-uncommitted` payload only with the matching transaction ID and
capability token. Finalization atomically replaces it with the committed
generation; abort atomically restores the last committed generation and marks
the candidate journal rolled back.

Fallback repair never swaps the failed payload into `previous`. It sets active
to the independently verified fallback and selects `previous` only from another
retained payload whose installed ledger and self-test prove it verified;
otherwise it writes `previous: null`. Concurrency tests cover every pointer
state against prepare, activation, finalization, abort, live-lock wait, stale
recovery, and fallback.

Launcher protocol version 1 is frozen for all payloads retained as rollback
candidates. An external installer may replace the launcher only with one that
declares compatibility with both the candidate and retained previous payload.
Ordinary `runx upgrade` never replaces its running launcher.

## Lifecycle Coordination State

### Exclusive Lock

Mutating install-state operations use one RunX-owned lock containing a random
ownership token, PID, process start identity, verified executable path,
operation, and creation time. A process releases only its own token. Stale
recovery is allowed only after platform process inspection proves that the
recorded process identity is no longer active.

### Instance Registry

Every payload creates an instance record after startup and removes it at normal
exit. Records contain the PID, random token, exact payload path, process start
identity, and start time. Upgrade correlates registry entries with OS process
inspection and terminates only another current-user process whose live
executable path equals the previous verified RunX payload. Child commands are
never registered as RunX instances and are never selected by filename.

### Transaction Journal

Every mutation journals these states:

```text
planned -> staged -> verified -> installed -> projected -> activated
        -> post-verified -> committed -> cleaned
```

Rollback states record pointer restoration, canonical resource restoration,
projection restoration, candidate removal, and residual cleanup. A launcher or
payload encountering an incomplete journal runs a local repair check before an
ordinary mutating lifecycle command. Read-only catalog commands never perform
remote or destructive repair.

Only cleanup of an inactive, OS-locked payload may be deferred. Deferred cleanup
cannot change the active version and is never reported as `scheduled` upgrade
success.

## Configuration Architecture

### Separate Contracts

`runx.yaml` is the project configuration. It retains the RunX command catalog
(`namespace`, commands, groups, imports, and catalog metadata) and adds optional
project-level configuration fields, including partial `agent.evolution`
overrides.

`$HOME/.guiho/runx/runx.global.yaml` is a distinct global configuration. It
contains user-wide RunX settings and the complete agent-evolution baseline. It
does not contain or substitute for a project command catalog.

Catalog selection remains:

1. explicit `--config` project path;
2. `runx.yaml` in the effective working directory;
3. no catalog.

RunX never searches parents and never treats `runx.global.yaml` as a catalog.
Global configuration is loaded independently, then matching project settings
overlay it recursively. Omitted project fields inherit; an explicitly provided
field replaces the matching scalar. Catalog-only fields have no global
counterpart and never merge from global state.

### Schemas And Runtime Validation

`runx.schema.json` and `runx.global.schema.json` are separate, strict schemas.
The project schema covers both the existing catalog and permitted project
settings; the global schema excludes catalog fields. Both accept exactly
`disabled`, `always-ask`, or `always-proceed` for every evolution policy leaf.
Unknown fields and invalid partial structures fail validation.

Runtime uses the embedded equivalent Go contract without network access. Files
created by an installed release begin with a schema comment pinned to that
exact GitHub Release version. Complete examples exercise every optional field
and are verified against the corresponding schema in CI.

## `runx init` Architecture

`init` resolves the project root as the parent of explicit `--config`, otherwise
the effective working directory; it does not search ancestors. It executes one
idempotent reconciliation plan:

1. validate RunX-owned roots and installed canonical resources;
2. verify or repair all global main-skill projections from the active canonical
   artifact;
3. create `AGENTS.md` when absent and reconcile only the RunX bounded block;
4. create or validate `runx.global.yaml`;
5. create or validate project `runx.yaml` without discarding an existing
   catalog;
6. compute effective agent-evolution values;
7. explain all three policies, recommend and offer all-`always-proceed`, then
   ask only unanswered global/project questions;
8. fail without writes when mandatory answers are missing in a noninteractive
   session;
9. validate a complete proposed write set before any configuration replacement;
10. atomically apply configuration and resource changes;
11. revalidate every common and RunX-specific check; and
12. print created, upgraded, verified, and unchanged results with absolute
    paths.

Bare `runx` no longer performs project or global resource mutations. It retains
the greeting and any separately authorized, local-cache-only update notice.
Explicit `init`, external installation, and approved upgrade own reconciliation.

The evolution policy governs AI agent behavior. The main skill must read the
effective policy before checking, upgrading, or creating issues. Conservatively,
the built-in background upgrade check is disabled when the effective upgrade
policy is `disabled`; it never performs an upgrade automatically.

## Agent Resource Architecture

- `skills/guiho-s-runx/SKILL.md` is the single canonical skill body.
- `skills/guiho-s-runx/agents/openai.yaml` is the canonical agent definition.
- `prompts/guiho-p-runx.md` is the main installation/setup prompt.
- `instructions/guiho-i-runx.md` is the one canonical managed instruction.
- Build embedding and release packaging read those exact files.
- Canonical installed copies live under the active immutable version.
- Global skill copies and the project instruction are projections recorded in
  the installed ownership ledger.

The public tree uses `agent skill upgrade` and `agent instruction upgrade` only.
`show` prints raw bundled content without mutation. `list` reports bundled IDs.
Instruction commands manage `AGENTS.md` only, preserve all bytes outside the
bounded block, reject duplicate/malformed markers, and never delete the file.

The main skill contains the exact `## CLI Evolution and Feedback` section,
repository and issue-create URLs, RunX-specific bug/improvement/review guidance,
all policy behavior, upgrade checking, `init` after upgrade, raw version
verification, and direct issue-URL reporting. The setup prompt covers product
purpose, installation, verification, initialization, and upgrade.

## Release Resolution

One strict SemVer implementation is shared by catalog, installer fixtures, and
upgrade tests. It rejects malformed versions and implements numeric prerelease
identifier ordering. The selector is exactly one of:

- exact raw SemVer;
- channel, derived from the first prerelease identifier; or
- default stable.

Selection paginates the full release catalog, excludes drafts, validates the
tag, rejects an incomplete or checksum-invalid release, verifies the complete
target matrix and common resources, and then chooses the highest compatible
version. Exact version and channel are mutually exclusive. The selected target
is frozen before any mutation.

## Installer And Reinstall Transaction

The PowerShell and POSIX scripts are bootstrap adapters around lifecycle
protocol v1. The verified candidate launcher is the single canonical filesystem
transaction engine. The adapters are verified with common black-box fixtures:

1. resolve exact selector, platform, native home, RunX home, launcher, and
   payload;
2. print version, OS, architecture, asset, source, and destination;
3. create and validate one strict descendant of `.guiho/.temp`;
4. download the selected platform pair and every common asset;
5. reject missing, duplicate, malformed, extra-required, or mismatched checksums;
6. verify executable format, OS, architecture, embedded build target, exact raw
   `--version`, and hidden `__self-test`;
7. write a strict bootstrap plan and random capability inside staging;
8. call `prepare` with that capability; the Go engine snapshots prior state,
   installs the immutable candidate, records the adapter PID/token and prior
   PATH/legacy-entrypoint plan, but exposes no new committed generation;
9. call `activate`; the engine installs/reconciles replaceable resources,
   launcher, and projections, writes an `active-uncommitted` pointer, and
   verifies the candidate through a transaction-authorized launcher dispatch;
10. snapshot the exact prior user PATH value in the adapter journal, add the
    canonical shared entry idempotently, verify shell resolution, and perform
    only ownership-proven legacy-entrypoint migration;
11. call `finalize` with the same transaction ID and capability; the engine
    repeats post-activation verification, atomically commits the pointer and
    installed ledger, and jointly commits the Go and adapter journals; or
12. on any adapter or engine failure, restore the exact prior PATH/legacy state
    and call `abort`, which restores pointer, launcher, resources, projections,
    and ledger; and
13. remove staging/backups only after `finalize` confirms joint commit.

Installing the already-active version is a repair transaction, not a no-op.
Configuration and persistent data are preserved. Any failure restores the
complete previous working state and never rolls back or deletes persistent
content.

The launcher bootstrap returns a structured result containing transaction ID,
previous/new versions, installed paths, projections, activation digest, and
rollback status. Both adapters must produce byte-equivalent final installed
trees and ledgers from the same fake release, and fault-injection suites compare
their rollback trees and exact prior PATH restoration at every adapter boundary.
Native Windows and POSIX CI runs execute the real scripts, not only parse their
fixtures.

`prepare`, `activate`, `finalize`, and `abort` are capability-bound
protocol operations. No operation accepts another transaction's token, and
terminal `finalize` or `abort` calls are idempotent. The lock and both journals
remain recoverable until finalization. If the adapter process disappears after
activation, ordinary launcher dispatch uses the last committed generation,
never the candidate. A later invocation reports the exact installer recovery
command; rerunning the same platform installer detects the abandoned adapter
PID, restores the prior PATH and ownership-proven legacy entrypoint from the
adapter journal, invokes `abort`, and only then begins a new transaction.
Backups are never deleted by stale recovery before joint commit.

After every successful reinstall or installer-driven version change, the
installer prints the exact `runx init` and `runx --version` follow-up. The main
skill defines these as mandatory steps for the managing AI agent; an
agent-managed lifecycle operation is not complete until both succeed. The
installer does not silently create an unrelated project configuration when run
outside the intended project. When it is explicitly given a project root, it
may offer to run `init` there; required interactive answers remain governed by
the `init` contract.

### Legacy Direct-Executable Migration

The installer inventories the actual command resolution before and after the
proposed user PATH change for PowerShell, `cmd.exe`, the current POSIX shell, and
Git Bash when present. It examines the historical default
`$HOME/.local/bin/runx[.exe]` and any resolved custom entrypoint, but it never
recursively searches arbitrary PATH directories.

Ownership proof requires an exact raw version plus a digest match against the
corresponding historical release checksum, or an already installed legacy
RunX ownership record. A proven legacy file is backed up and removed only after
the new launcher passes activation verification; rollback restores it. If the
historical release/checksum is unavailable, proof fails. An unproven file is
never changed. If any shell would still resolve an unproven/custom entrypoint
before the canonical launcher, installation fails closed with the exact path,
shell, PATH order, and manual remediation; it does not claim success or delete
the file.

## Uninstall Architecture

The two scripts and Cobra command consume the installed ownership ledger and
the same preservation model. They first print every absolute target under
`REMOVE` or `PRESERVE`. Without `--yes`, an interactive terminal must confirm;
noninteractive input fails closed. `--dry-run` performs no write.

Default removal includes the launcher, all payloads/resources/state, global
configuration, persistent data/databases, managed skill projections, the RunX
block, and project `runx.yaml`. `--preserve-config` and `--preserve-data` apply
exactly as Convention 0001 defines. Shared GUIHO directories, shared PATH,
foreign files, and all non-RunX `AGENTS.md` bytes are always preserved.

On Windows, the external PowerShell uninstaller is the universally supported
path because it does not run from either locked RunX executable. Cobra uninstall
uses a precise in-process handoff only on filesystems that support verified
rename plus delete-on-close semantics:

```text
payload plans/confirms -> acquire lock -> quarantine ordinary targets
  -> rename canonical launcher and running payload into owned temp quarantine
  -> mark both handles delete-on-close with Windows disposition semantics
  -> verify canonical paths and RunX home are absent
  -> print final result -> payload exits -> launcher exits
  -> Windows closes both image sections and deletes quarantine files
  -> invoking shell regains control only after launcher exit
```

The payload retains stdout/stderr and the uninstall exit status. The launcher
continues waiting for the payload and returns that status; it performs no new
dispatch after its canonical path has been quarantined. Before renaming, the
transaction snapshots all reversible targets. If either rename, delete-pending
mark, or canonical-path verification fails, it clears/reverses every possible
quarantine move, restores the snapshot, returns exit 5, and prints the exact
`uninstall.ps1` recovery command. It never starts a detached helper or reports
eventual success.

The journal state is held in the quarantined transaction directory and each
transition is flushed before the next rename. Crash recovery is:

| Crash boundary | Recovery |
| --- | --- |
| Before launcher/payload quarantine | Next lifecycle operation restores or completes from journal |
| After one executable quarantine | External `uninstall.ps1` reads the token-bound journal and either completes removal or restores the canonical name |
| After both delete-pending marks | Process exit completes kernel deletion; external uninstaller removes any residual empty transaction directory |
| Unsupported filesystem/API | Cobra makes no mutation and directs the user to `uninstall.ps1` |

Black-box Windows tests invoke Cobra through the real launcher on NTFS, verify
the shell does not regain control before both canonical names disappear, assert
the final exit code/output, inject failure at every handoff boundary, and prove
that an unsupported delete-on-close capability fails before mutation. Filename-
only process or file selection is forbidden.

## Synchronous Upgrade Architecture

`runx upgrade` prints the platform recovery command before network access and
again as its final block for every outcome. The final block pins the resolved
exact version.

Upgrade then:

1. acquires the exclusive lock;
2. resolves a complete exact/channel/stable release;
3. stages and verifies every selected artifact;
4. executes exact version and `__self-test` checks;
5. installs the immutable version and canonical resources;
6. terminates only verified other instances of the previous payload;
7. snapshots and replaces all manifest-owned projections/resources;
8. atomically activates `current.json`;
9. verifies launcher version, post-activation self-test, and every projection;
10. commits the installed ledger and journal;
11. rolls back pointer/resources/projections on any failure; and
12. reports previous/new versions, launcher/payload paths, changed artifact IDs,
    verification, and the final recovery block.

Ordinary upgrade does not replace the launcher. Installer repair is the only
launcher replacement path. The immediately previous verified payload is
retained; older inactive payloads are retired by manifest-owned garbage
collection.

## npm Distribution

The npm executable and independent distribution path are retired. The final
source contract removes `package.json.bin`, deletes `scripts/runx-bin.mjs`,
stops npm publication, and never creates `node_modules/.bin/runx`. If
`package.json` remains for retained legacy-reference checks, it is private and
has no CLI entrypoint; any remaining version field is either removed or governed
by Mirror. The publish workflow never runs `npm version`.

The `.guiho/runx/npm` cache is declared disposable legacy RunX state and is
removed only by a manifest/legacy-aware reinstall or default uninstall. An npm
package deprecation notice and README migration command may be published only
with separate external/publication authorization. `npx runx` has no supported
post-cutover version promise; the canonical entrypoint is exclusively
`$HOME/.guiho/bin/runx[.exe]`.

## Help And Version

- Root `-v`/`--version` prints one raw SemVer-compatible value. Source builds
  use `0.0.0-dev`, not `dev`.
- `--help-tree-depth` is a string contract accepting `max` or an integer greater
  than 1; default is `max`.
- `--help-tree-global-flags` is a false-by-default Boolean presence flag.
- Without it, global flags appear once at tree root; with it, inherited flags
  appear at each applicable command.
- Every public command retains help, help docs, and help tree from the one Cobra
  tree.

## Repository Tooling Architecture

- Root `runx.yaml` catalogs every supported development, validation,
  documentation, schema, build, install-test, upgrade-test, package, and release
  preparation command with stable UIDs.
- Root `xdocs.yaml` replaces the legacy TOML authority and covers every owned
  directory, including agent resource and review subtrees.
- `mirror.yaml` governs Git version identity and `package.json`; no workflow
  mutates version fields with npm.
- CI calls cataloged RunX commands and verifies RunX, XDocs, Mirror, schemas,
  release assets, and lifecycle fixtures.

## Verification Seams

All filesystem, home, process, terminal, clock, token, GitHub, downloader, PATH,
and launcher operations are dependency-injected. Tests never use the real home,
PATH, global agent directories, running installation, or release.

Required suites include:

- strict `artifacts.json`, `current.json`, installed ledger, lock, instance, and
  journal decoding;
- path traversal, symlink/junction escape, foreign ownership, and duplicate
  target rejection;
- launcher exact delegation, fallback, pointer repair, and exit codes;
- complete SemVer/channel pagination and incomplete-release exclusion;
- all eight payload and launcher cross-builds plus native-platform smoke;
- candidate format/architecture/build-target/version/self-test validation;
- install/reinstall/repair rollback at every journal boundary on POSIX and
  Windows fixtures;
- default/preserved/dry-run/noninteractive uninstall plans;
- other-instance path verification and child-process preservation;
- interrupted-upgrade recovery and post-activation rollback;
- config inheritance, strict schemas, version-pinned references, interactive
  init, noninteractive failure, and idempotency;
- exact agent projections and byte-preserving instruction blocks;
- npm no-download delegation;
- help/version contract probes;
- exact 25-asset build, checksum, and workflow verification; and
- Mirror, RunX, XDocs, Go format/test/vet/build, Git, and documentation gates.

## Migration And Compatibility

This is a one-way public lifecycle migration, with transactional downgrade by
explicit installer version still supported. The first conforming external
installer migrates ownership-proven legacy direct installations. Old payloads
cannot self-upgrade into the new layout safely; their `runx upgrade` recovery
message and public documentation must direct users to the new installer.

Historical requirements, plans, reviews, and decisions remain as evidence but
must be marked superseded wherever they claim current direct replacement,
exact-eleven/exact-fourteen assets, background bootstrap mutation, old global
catalog fallback, `update` command names, or npm download ownership.

### Compatibility Matrix

| Combination | Required behavior |
| --- | --- |
| Old installer with first convention release | Fails compatible-asset lookup because legacy payload names are absent; leaves old installation unchanged |
| Transition installer with legacy default/channel release before protocol v1 exists | Executes only the preserved checksum-verified legacy transaction; never mixes shapes |
| Transition installer with explicit legacy version | Executes only the preserved checksum-verified legacy transaction while transition support is active |
| Transition installer with protocol-v1 release | Selects protocol v1 and never falls back to legacy for that selector |
| Hardened installer with old release | Rejects missing protocol-v1 manifest, launcher, and resources before mutation |
| Hardened installer with protocol-v1 release | Performs the full protocol-v1 transaction |
| Old payload `upgrade` against convention release | Finds no legacy payload asset and fails safely with external installer guidance |
| New payload in legacy direct layout | Domain/read-only commands may run; lifecycle mutation refuses and prints the canonical installer recovery command |
| New launcher with old payload | Bootstrap rejects missing launcher protocol/self-test compatibility; never activates |
| New installer with convention release | Performs full protocol-v1 transaction and legacy shadow check |

### Multi-PR Cutover

Implementation units before cutover may merge only additive or dormant pieces:
formats, schemas, canonical resources, internal packages, launcher/lifecycle
code, `.next` installer fixtures, tests, and repository tooling. They must not
change the canonical remote installer paths, public release verifier, current
workflow uploads, or old payload asset names while the latest release remains
legacy.

One transition PR owns the inseparable protocol publication surface:

- a canonical dual-shape transition installer;
- payload asset rename and launcher publication;
- protocol-v1 release contract, builder, verifier, and publish workflow;
- public upgrade compatibility/failure guidance;
- npm executable/publication retirement;
- transition README/DOCS/skill/instruction/prompt guidance; and
- removal of development-only `.next` paths.

Before the first protocol-v1 release exists, the transition installer recognizes
the legacy release shape and executes only the preserved, checksum-verified
legacy transaction. It does not mix legacy and protocol-v1 assets. Once any
compatible protocol-v1 release is published for the requested selector, it
always selects protocol v1 and never falls back to legacy for that selector.
This keeps the documented `main/devops/install.*` URLs functional between merge,
tag creation, asset upload, and release visibility.

The transition PR may not merge unless the same approved change window
separately authorizes the Mirror-managed version, tag/push, first protocol-v1
GitHub Release, and immediate remote lifecycle smoke. The release sequence is:

1. merge the reviewed transition PR;
2. verify the canonical installer can still install the latest legacy release
   into an isolated home;
3. apply the separately approved Mirror version and protected tag;
4. publish and verify the complete protocol-v1 assets without marking success
   early;
5. verify exact-version and default/channel protocol-v1 installation,
   upgrade-recovery output, launcher dispatch, and rollback remotely; and
6. if publication fails, keep the transition installer on its isolated legacy
   path and do not claim protocol-v1 availability.

After the first protocol-v1 release passes remote acceptance, a second reviewed
hardening PR removes every legacy-shape branch from the canonical installers,
updates final documentation/recovery text, and makes incomplete legacy releases
unconditionally ineligible. Existing protocol-v1 assets make this hardening
merge safe without another release. Full Convention 0001 compliance is claimed
only after this hardening PR is integrated and its main-branch installer smokes
against the already published protocol-v1 release.

The first protocol-v1 release is the artifact rollback boundary. Failure before
tag publication leaves the dual-shape installer serving the legacy release.
Failure after publication but before acceptance halts promotion, preserves the
last accepted legacy release, publishes no false success claim, and requires a
reviewed corrective protocol-v1 release rather than reusing or mutating released
assets. Old installers/upgraders fail closed against the renamed assets
throughout. The dual-shape transition is temporary and is itself a recorded
noncompliance until the hardening PR removes it.

## Risks And Mitigations

| Risk | Mitigation |
| --- | --- |
| Launcher protocol bricks all versions | Freeze protocol v1, cross-version fixtures, retain previous verified pointer |
| Ownership manifest enables unsafe deletion | Strict relative paths, allowed roots, installed-ledger intersection, symlink/junction tests |
| Windows locks prevent synchronous completion | Immutable payload activation; defer only inactive cleanup, never activation |
| Old `.local/bin` file shadows launcher | Check exact legacy path, prove checksum ownership, migrate transactionally, fail closed if foreign |
| Two installer implementations drift | Shared JSON fixtures, identical state assertions, platform CI |
| Agent copies drift | One canonical source embedded and packaged; projection hashes recorded |
| Config migration discards catalogs | Project schema retains catalog; init overlays policy without recreating valid catalog |
| Partial release gets selected | Full pagination plus exact 25-asset/completeness validation before mutation |
| Stale plan executes on stale main | Every unit records an approved base and re-plans after upstream integration |
| Old docs/tests reassert obsolete behavior | Dedicated supersession and final repository-search gates |

## Approval And Production Boundary

Approval of this architecture authorizes planning only. It does not authorize
implementation, branch creation, commit, push, issue creation, live install,
PATH or global-agent mutation, version application, tag, GitHub Release, npm
publication, deployment, or production promotion.

The human approval must explicitly confirm the proposed identity table. Every
implementation unit then requires its own approved base, branch/worktree,
focused PR, exact-head implementation review, exact-head validation, gated
integration, and post-merge reachability check. Mirror application and release
remain separate explicit decisions after integrated validation.

## References

- [Compliance audit](../reviews/guiho-convention-0001-cli-compliance-audit.md)
- [Current CLI architecture](./cli-architecture.md)
- [Prior RFC 0034 plan](../plans/rfc-0034-cli-compliance-migration.md)
- [Prior RFC 0034 implementation record](../todo/rfc-0034-cli-compliance-migration-implementation.md)
- [Prior RFC 0034 validation](../validation/rfc-0034-cli-compliance-migration.md)
- [RunX TODO](../../TODO.md)
