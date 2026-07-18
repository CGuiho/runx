---
name: RunX Upgrade Reliability Design
purpose: Define the implementation-ready reliability contract for RunX self-upgrade, version listing, and recovery installation.
description: Specifies release discovery, visible progress, transactional replacement, rollback, recovery commands, installer hardening, tests, ownership, and release gates for GitHub issues 12 and 13.
created: 2026-07-15
flags:
  - accepted
tags:
  - architecture
  - cli
  - releases
  - reliability
keywords:
  - runx upgrade
  - runx upgrade list
  - issue 12
  - issue 13
  - windows self-upgrade
  - recovery install command
  - transactional replacement
owner: runx-superpowers-specs
---

# RunX Upgrade Reliability Design

## Status and Scope

This is the approved design for `CGuiho/runx` issues #12 and #13. It replaces
the assumption that a scheduled Windows move is a successful upgrade. It also
extends the existing verified swap into a complete operator experience:

- an upgrade plan and visible progress before long operations;
- synchronous, verified replacement with rollback;
- an exact-version recovery installation command after every attempt;
- a separate process-stop command for replacement conflicts;
- a complete, semantically sorted catalog of all published releases; and
- hardened direct installers that verify what they install.

The work remains inside the Bun/TypeScript native CLI and its existing
PowerShell and POSIX installers. Citty continues to own command routing and
ordinary argument parsing. No Node.js runtime dependency is introduced.

## Problem Statement

### The 0.2.5 bootstrap trap

RunX 0.2.5 downloads an update to `runx.exe.new`, launches a detached
`cmd.exe` process that tries to move the new file directly over the running
`runx.exe`, discards the helper's output, and returns `scheduled: true` without
verifying the canonical path. Windows can reject the direct overwrite while
the original image is mapped. The detached helper then fails silently and the
next invocation still reports 0.2.5.

The updater that is currently running controls the replacement procedure.
Therefore a corrected 0.2.6 asset cannot retroactively fix the 0.2.5 updater
that is trying to install it. Users on 0.2.5 need the exact-version recovery
installer as the bootstrap bridge into the first reliable release.

### What 0.2.6 already proves

RunX 0.2.6 implements the correct Windows primitive: rename the mapped current
executable to a backup, place the new executable at the canonical path, execute
that canonical path with `--version`, roll back on failure, and defer only
deletion of the mapped backup. This design preserves that mechanism.

The remaining defects are broader than the swap itself:

- text output is produced only after the awaited upgrade operation returns;
- release lookup failures and missing assets can look like a no-op;
- download validation and temporary-file naming are insufficient;
- asset selection does not use the installers' baseline-first policy;
- `upgrade list` reads only one page of 20 releases and lacks metadata;
- the current comparator does not implement SemVer prerelease precedence;
- JSON output lacks phases, structured failures, and recovery instructions;
- direct installers do not verify the installed version or roll back; and
- generic CLI error handling cannot guarantee recovery output after failure.

## Goals

1. A successful `runx upgrade` means the canonical executable already runs as
   the selected target version when the command returns.
2. Human users see the plan before download and a message before every
   potentially slow or disruptive phase.
3. Every result - upgraded, already current, dry run, or failed - ends with an
   exact-version recovery install command and a separate process-stop command.
4. `runx upgrade list` returns every published RunX release, newest SemVer
   first, with clear channel and compatibility metadata.
5. Text and JSON expose equivalent facts without corrupting JSON with streamed
   text.
6. Direct installers apply the same target, platform, architecture, variant,
   verification, and failure semantics as in-process upgrade.
7. Automated tests make Windows replacement, rollback, output ordering,
   catalog completeness, and installer verification release blockers.

## Non-Goals

- RunX will not add a background update daemon or automatically upgrade on
  startup.
- RunX will not automatically kill other RunX processes. It prints an explicit
  process-stop command for the user to choose.
- npm, Homebrew, winget, Scoop, and other package-manager integrations are not
  added by this work.
- This work does not add a new release-signing or checksum publication system.
  It validates platform binary format and verifies the downloaded executable's
  reported version. A future signing design can strengthen provenance.
- Scheduled replacement is not retained as a fallback success state. Only
  cleanup of a renamed old executable may be deferred.

## Chosen Architecture

The implementation uses five bounded units.

### 1. Release catalog

The release catalog owns GitHub pagination, tag normalization, SemVer parsing,
channel labeling, publication metadata, and compatible-asset selection. It
returns normalized entries and never silently converts transport or schema
failure into "already current."

### 2. Upgrade planner

The planner combines the catalog, installed version, platform, architecture,
CPU variant policy, and canonical executable path into an immutable plan. The
in-process upgrader always uses the baseline-first x64 policy; direct installers
may retain their explicit variant override. No download begins before the plan
is available and rendered.

```text
UpgradePlan
  currentVersion
  targetVersion
  os
  arch
  assetName
  assetUrl
  executablePath
  recovery
```

The default target is the highest valid stable SemVer release. Prereleases are
always visible in `upgrade list` but are not selected by default. A future
explicit prerelease selector can reuse the catalog without changing listing
semantics.

### 3. Transactional installer

The transactional installer owns download, validation, backup, canonical
replacement, target verification, rollback, and cleanup. It emits typed phase
events but does not render terminal output.

### 4. Upgrade reporter

The reporter converts the same plan, events, result, error, and recovery data
into either streaming human text or one buffered JSON document. This keeps
execution independent from presentation and prevents text/JSON drift.

### 5. Recovery contract

The recovery contract generates the exact pinned installer and process-stop
commands for the detected platform. The CLI and installer contract tests use
the same version and asset-selection rules. The PowerShell and shell scripts
remain platform-native implementations, with cross-language contract tests to
prevent drift.

## Release Catalog Contract

### Retrieval

- Request `GET /repos/CGuiho/runx/releases?per_page=100&page=N`.
- Continue until the GitHub `Link` header has no `rel="next"` entry. A short
  page may be used as a secondary termination signal, never as a substitute
  when a next link exists.
- Include every published, non-draft release, including prereleases.
- Preserve the original tag and normalize project tags such as
  `@guiho/runx@0.3.0` and `v0.3.0` to `0.3.0` for comparison.
- Report non-2xx responses, rate limits, malformed payloads, and pagination
  failures explicitly.

### Ordering and channels

- Valid SemVer releases sort by SemVer precedence, descending.
- A stable release outranks its prereleases.
- `alpha`, `beta`, and `rc` identifiers are displayed exactly as those
  channels. Another textual prerelease identifier is displayed as itself.
- A numeric-only prerelease identifier is displayed as `prerelease`.
- Published tags that are not valid SemVer remain visible after valid SemVer
  entries, ordered by publication time descending, and use channel `other`.
- `latest stable` means the greatest valid non-prerelease SemVer, independent
  of GitHub's presentation order.

### Catalog entry

Each text row and JSON entry exposes:

- original tag and normalized version;
- channel and prerelease state;
- publication timestamp;
- whether it is the installed version;
- whether it is the latest stable release; and
- whether the release has a compatible asset for the current OS,
  architecture, and selected x64 variant policy.

Text output uses the existing aligned-table style:

```text
AVAILABLE RUNX VERSIONS

VERSION              CHANNEL   PUBLISHED    CURRENT   LATEST   ASSET
0.3.0-alpha.2        alpha     2026-07-15                       yes
0.2.7                stable    2026-07-15             yes       yes
0.2.6                stable    2026-07-14   yes                  yes
```

The actual catalog is not capped at the example's three rows.

JSON uses a complete command-specific wrapper rather than returning a bare
array:

```json
{
  "schemaVersion": 1,
  "command": "runx upgrade list",
  "currentVersion": "0.2.6",
  "latestStableVersion": "0.2.7",
  "releases": [
    {
      "tag": "@guiho/runx@0.2.7",
      "version": "0.2.7",
      "channel": "stable",
      "prerelease": false,
      "publishedAt": "2026-07-15T00:00:00Z",
      "current": false,
      "latestStable": true,
      "compatibleAsset": {
        "name": "runx-windows-x64-baseline.exe",
        "url": "https://github.com/CGuiho/runx/releases/download/%40guiho%2Frunx%400.2.7/runx-windows-x64-baseline.exe"
      }
    }
  ]
}
```

`releases` contains every normalized published release from all retrieved
pages, not only compatible or stable releases. `currentVersion` is always the
installed version. `latestStableVersion` is the greatest valid stable SemVer,
or null when no published stable SemVer exists.

### Asset selection

The CLI and installers use one ordered candidate policy:

- x64 baseline preference: `baseline`, default, `modern`;
- x64 default preference: default, `baseline`, `modern`;
- x64 modern preference: `modern`, default, `baseline`;
- ARM64: the single ARM64 asset for the detected OS.

The selected asset name and exact GitHub download URL are part of the plan. If
the latest stable release lacks a compatible asset, upgrade fails explicitly;
it does not silently skip to an older release. The recovery command remains
pinned to that target and the installer reports the missing asset clearly.

## Upgrade Event and Human Output Contract

### Event sequence

The execution layer emits ordered events with a monotonically increasing
sequence number. Every event uses one shared phase name and one shared status:

- phase: `plan`, `download`, `validate`, `replace`, `verify`, `cache`, or
  `cleanup`;
- status: `started`, `succeeded`, `skipped`, or `failed`.

A normal mutation emits `started` before a phase and `succeeded` after it.
Failure emits `failed` for that same phase; there is no separate `failed`
phase. An already-current or dry-run outcome completes `plan` and marks
mutation phases `skipped`. RunX has no upgrade cache yet, so `cache` is emitted
as `skipped` after `verify`; retaining the common phase keeps the JSON contract
aligned with the other GUIHO CLIs without inventing cache behavior.

In text mode, the human message derived from each `started` event is written
and flushed before the implementation awaits that operation. In particular,
`Downloading...` is visible before the binary body is awaited,
`Validating...` appears before binary validation, `Replacing...` appears only
after a complete validated download, and `Verifying...` appears before
launching the canonical executable. `cleanup: skipped` is valid only when the
verified upgrade has handed deletion of a locked old image to the detached
cleanup helper; it never means replacement was scheduled.

### Successful text output

```text
------------------------------------------------------------
  Upgrading the CLI
------------------------------------------------------------
  current : 0.2.6
  target  : 0.2.7
  os      : windows
  arch    : x64
  binary  : runx-windows-x64-baseline.exe
  path    : C:/Users/crist/.local/bin/runx.exe
  url     : https://github.com/CGuiho/runx/releases/download/%40guiho%2Frunx%400.2.7/runx-windows-x64-baseline.exe
------------------------------------------------------------
Downloading...
Validating...
Replacing...
Verifying...
Upgrade complete: 0.2.6 -> 0.2.7

If the new version is not active, install RunX 0.2.7 directly:
  powershell.exe -NoProfile -ExecutionPolicy Bypass -Command '& ([scriptblock]::Create((Invoke-RestMethod "https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1"))) -Version "0.2.7"'

If RunX is still running and blocks installation, stop it first:
  powershell.exe -NoProfile -Command "Get-Process runx -ErrorAction SilentlyContinue | Stop-Process -Force"
```

The same recovery block follows `Already up to date`, `Dry run complete`, and
failure results. For an already-current result, the command is pinned to the
current version. For a release-discovery failure where no target can be
truthfully resolved, recovery is pinned to the installed version and described
as a repair installation.

Failures name the phase and specific cause, exit nonzero, and never print
`Upgrade complete`. Recovery instructions are not hidden behind `--verbose`.
Plan and phase text, non-failure results, and recovery instructions go to
stdout. The concise failure line goes to stderr, while recovery remains on
stdout. Diagnostic stacks remain a verbose-only stderr concern.

## JSON Contract

`--format=json` emits exactly one valid JSON document. It does not stream human
phase lines. The reporter buffers events and produces this shared envelope:

```json
{
  "schemaVersion": 1,
  "command": "runx upgrade",
  "outcome": "upgraded",
  "plan": {
    "currentVersion": "0.2.6",
    "targetVersion": "0.2.7",
    "os": "windows",
    "arch": "x64",
    "assetName": "runx-windows-x64-baseline.exe",
    "assetUrl": "https://github.com/CGuiho/runx/releases/download/%40guiho%2Frunx%400.2.7/runx-windows-x64-baseline.exe",
    "executablePath": "C:/Users/crist/.local/bin/runx.exe"
  },
  "events": [
    { "sequence": 1, "phase": "plan", "status": "started" },
    { "sequence": 2, "phase": "plan", "status": "succeeded" },
    { "sequence": 3, "phase": "download", "status": "started" },
    { "sequence": 4, "phase": "download", "status": "succeeded" },
    { "sequence": 5, "phase": "validate", "status": "started" },
    { "sequence": 6, "phase": "validate", "status": "succeeded" },
    { "sequence": 7, "phase": "replace", "status": "started" },
    { "sequence": 8, "phase": "replace", "status": "succeeded" },
    { "sequence": 9, "phase": "verify", "status": "started" },
    { "sequence": 10, "phase": "verify", "status": "succeeded" },
    { "sequence": 11, "phase": "cache", "status": "skipped" },
    { "sequence": 12, "phase": "cleanup", "status": "started" },
    { "sequence": 13, "phase": "cleanup", "status": "skipped" }
  ],
  "result": {
    "installedVersion": "0.2.7",
    "cleanupDeferred": true
  },
  "recovery": {
    "targetVersion": "0.2.7",
    "targetSource": "resolved",
    "installCommand": "powershell.exe -NoProfile -ExecutionPolicy Bypass -Command '& ([scriptblock]::Create((Invoke-RestMethod \"https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1\"))) -Version \"0.2.7\"'",
    "stopProcessCommand": "powershell.exe -NoProfile -Command \"Get-Process runx -ErrorAction SilentlyContinue | Stop-Process -Force\""
  },
  "error": null
}
```

`outcome` is one of `upgraded`, `up-to-date`, `dry-run`, `rolled-back`, or
`failed`. A successful rollback uses `rolled-back`, retains the primary error,
and reports the restored installed version in `result`. `failed` covers a
failure before replacement, a failure for which rollback is unnecessary, or a
rollback failure. Failed and rolled-back JSON results exit nonzero.

The error is null on non-failure outcomes. Otherwise it contains `code`,
`phase`, and `message`. When release discovery fails before a plan exists,
`plan` is null and recovery uses `targetSource: "fallback-current"` with the
installed version. All other resolved plans use `targetSource: "resolved"`.
Every envelope includes `result`; it is null when no verified installed-version
result exists. The top-level CLI handler does not append a second generic text
error after the JSON document.

A discovery-failure recovery block is visibly labeled in text:

```text
Repair reinstall (target lookup failed; pinned to installed RunX 0.2.6):
  powershell.exe -NoProfile -ExecutionPolicy Bypass -Command '& ([scriptblock]::Create((Invoke-RestMethod "https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1"))) -Version "0.2.6"'
```

Stable error codes are:

- `release_lookup_failed`
- `release_payload_invalid`
- `no_compatible_asset`
- `download_failed`
- `download_invalid`
- `backup_failed`
- `replace_failed`
- `verification_failed`
- `rollback_failed`
- `installer_verification_failed`

## Transactional Replacement Contract

The state machine is:

```text
discovered -> planned -> downloaded -> validated -> backed_up
  -> installed -> verified -> complete
```

The downloaded file and backup use unique names in the executable's directory,
for example `runx.exe.new.<pid>.<nonce>` and `runx.exe.old.<pid>.<nonce>`. The
same directory preserves rename semantics and avoids cross-volume moves.

Before modifying the canonical path, the installer requires:

- a successful HTTP response;
- a non-empty body;
- an expected Windows PE, Linux ELF, or macOS Mach-O binary signature;
- a filename compatible with the plan; and
- executable permissions on Unix platforms.

Replacement then proceeds:

1. Rename the canonical executable to the unique backup path.
2. Rename the validated download to the canonical path.
3. Launch the canonical path with `--version`.
4. Require exit code zero and output exactly equal to the target version.
5. Mark the upgrade complete only after verification.

If any step after backup fails, remove a failed canonical replacement when
present and rename the backup to the canonical path. A rollback failure reports
both the primary and rollback errors and identifies every affected path.
Temporary downloads are removed on every terminal path.

On Windows, the mapped old executable may prevent immediate backup deletion.
A detached retry helper may remove that backup after the current process exits.
This is cleanup only: the target is already installed and verified at the
canonical path. No result field or message describes the upgrade itself as
scheduled.

## Recovery Command Contract

Recovery instructions are generated before mutation and retained through every
error. They are always printed after the final text result and always included
in JSON.

Windows uses the supported PowerShell installer with an explicit version. The
command invokes `powershell.exe` directly so it is usable from the user's Git
Bash session as well as a PowerShell terminal. Its separate stop command is:

```text
powershell.exe -NoProfile -Command "Get-Process runx -ErrorAction SilentlyContinue | Stop-Process -Force"
```

Linux and macOS use the supported shell installer with an explicit version:

```text
curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | sh -s -- --version '0.2.7'
```

Their separate stop command is:

```text
pkill -x runx
```

The real output substitutes the fully resolved target, including complete
prerelease identifiers such as `0.3.0-alpha.2`. It never uses `latest` in a
recovery command. Shell quoting is generated and tested rather than assembled
from unescaped user input.

The recovery object records where its target came from. `resolved` means
release discovery selected the target. `fallback-current` means discovery
failed before a target could be selected, so the visibly labeled repair
reinstall is pinned to the installed version instead of inventing a target.

The process-stop command is advisory and never executed automatically. If the
installer detects an active-image conflict, it preserves the old installation,
exits nonzero, and tells the user to run the stop command before retrying the
same pinned install command.

## Direct Installer Hardening

`devops/install.ps1` already accepts `-Version`; `devops/install.sh` already
accepts `--version`. Both must now:

1. preserve full stable or prerelease versions when resolving the release tag;
2. use the same ordered asset candidates as the in-process planner;
3. print the target, OS, architecture, asset, destination, and exact URL before
   download;
4. distinguish an unavailable candidate from a download, validation,
   replacement, permission, or verification failure;
5. download and validate before modifying the destination;
6. install transactionally with backup and rollback;
7. execute the canonical destination with `--version` and require the requested
   version; and
8. exit nonzero without claiming installation when any requirement fails.

Installer URL, tag, and candidate behavior is covered by shared fixtures and
cross-language contract tests. The scripts remain independently executable;
the native CLI does not invoke a remote script during normal self-upgrade.

## Error Handling Rules

- Network, GitHub API, malformed release, and missing-asset errors are never
  converted to an up-to-date result.
- An equal installed and target version is the only normal up-to-date result.
- A newer target without a compatible asset is an explicit release-integrity
  failure.
- The first failed mutation phase determines the primary error code.
- Rollback is mandatory after the canonical executable has been moved.
- A successful rollback still yields a failed upgrade and nonzero exit.
- A rollback failure includes both errors and retains recovery commands.
- Human output remains concise; `--verbose` adds technical cause chains without
  changing exit status or recovery visibility.
- Text and JSON are derived from the same structured outcome.

## File Ownership

| Path | Responsibility |
| --- | --- |
| `source/releases.ts` | GitHub pagination, release normalization, SemVer/channel ordering, and compatible-asset selection. |
| `source/recovery.ts` | Exact-version install and process-stop command generation with platform-safe quoting. |
| `source/self-management.ts` | Upgrade planning orchestration, download validation, transactional replacement, verification, rollback, and cleanup. |
| `source/upgrade-output.ts` | Human event streaming and buffered JSON outcome rendering. |
| `source/types.ts` | Catalog, plan, event, outcome, structured error, and recovery contracts. |
| `source/render.ts` | Aligned `upgrade list` table rendering alongside existing command tables. |
| `source/cli.ts` | Citty routing and upgrade-specific exit handling without a second token parser. |
| `source/releases.spec.ts` | Pagination, SemVer, channel, metadata, and asset-selection tests. |
| `source/self-management.spec.ts` | Platform transaction, validation, verification, rollback, and cleanup tests. |
| `source/upgrade-output.spec.ts` | Human ordering, failure recovery, and single-document JSON tests. |
| `source/guiho-runx.spec.ts` | Command integration for upgrade, check, list, dry-run, and JSON routing. |
| `devops/install.ps1` | Exact-version transactional Windows installation and verification. |
| `devops/install.sh` | Exact-version transactional Linux/macOS installation and verification. |
| `.github/workflows/ci.yml` | Linux and Windows release-blocking coverage, including installer contracts. |
| `DOCS.md`, `README.md`, `CHANGELOG.md` | Public command, recovery, and release behavior once implementation is complete. |
| `source/source.xdocs.md`, `devops/devops.xdocs.md`, and affected docs descriptors | Structured ownership for changed and new files. |

The implementation plan may keep a proposed new source file inside an existing
module if review proves the boundary too small to justify a file. It may not
merge release retrieval, terminal rendering, and filesystem mutation back into
one untestable operation.

## Test Strategy

### Release catalog

- multiple GitHub pages, including more than 100 releases;
- `Link`-header continuation and mid-pagination failure;
- stable, `rc`, `beta`, `alpha`, arbitrary, and numeric prerelease channels;
- correct SemVer precedence and invalid-tag fallback ordering;
- current/latest markers and publication timestamps;
- Windows, Linux, macOS, x64, ARM64, and CPU-variant asset selection; and
- missing compatible assets without fallback to an older release.

### Output and CLI

- plan text is captured before the binary fetch starts;
- every phase is emitted and flushed in exact order;
- already-current and dry-run paths omit mutation phases;
- every success and failure path ends with pinned recovery commands;
- release-discovery failure recovers to the installed version;
- stable and prerelease command quoting is exact;
- JSON is one parseable document with equivalent facts;
- failed JSON has a nonzero exit and no trailing generic text; and
- Citty help, `-v`, `--version`, and manifest-free behavior remain unchanged.

### Transactional replacement

- running Windows executable renamed and replaced before success;
- canonical path reports the target version before the promise resolves;
- invalid download rejected before the canonical path changes;
- stale files do not collide because names are unique;
- permission, backup, rename, and process-launch failures;
- version mismatch and nonzero verification exit;
- successful rollback after every post-backup failure;
- rollback failure reports both causes and paths;
- temporary download cleanup on every terminal path;
- immediate or deferred backup cleanup without scheduled-upgrade output; and
- equivalent verified transaction behavior on Linux and macOS.

### Direct installers

- explicit stable and prerelease versions;
- candidate order for every OS, architecture, and x64 variant;
- missing asset versus network and replacement errors;
- temporary-directory installation with no user PATH pollution in tests;
- exact canonical version verification;
- rollback after invalid or mismatched binary; and
- PowerShell and POSIX recovery commands executed against controlled release
  fixtures.

## Validation and Release Gates

Implementation is not complete until all of these pass:

1. `bun run typecheck`
2. `bun test`
3. `bun run build`
4. `bun run binary`
5. `bun run binaries`
6. Windows CI, including live executable replacement and PowerShell installer
   contract tests
7. Linux CI, including POSIX installer and catalog tests
8. strict XDocs metadata and whole-tree doctor validation
9. native text and JSON smoke tests outside a manifest directory
10. `git diff --check` and a reviewed task-only diff

Release preparation then follows the repository's Mirror workflow:

- update public docs and the configured changelog;
- run `mirror version plan patch` and review the exact next version;
- never hand-edit Mirror-managed version fields;
- apply and publish only with explicit user release authorization; and
- wait for the protected version-tag workflow to publish all native assets and
  npm successfully.

### Bootstrap and issue-closure gates

The public acceptance run must test both paths:

1. From 0.2.5 on Windows, use the printed or documented exact-version recovery
   installer to enter the new reliable release. Do not claim 0.2.5 can repair
   its own updater.
2. From a reliable installed release, run `runx upgrade` to the newly published
   patch and confirm the visible plan, phase order, synchronous replacement,
   recovery block, and immediate `runx --version` result.

Issue #12 closes only after the published native asset passes the real Windows
upgrade and complete-catalog acceptance checks. Issue #13 closes only after the
printed stable and prerelease recovery commands are copied into fresh test
shells, install the pinned versions, and verify those exact versions. Neither
issue closes on source tests alone or before the release assets exist.

## Rejected Alternatives

### Keep scheduled replacement and improve the message

Rejected because a schedule is not proof of replacement. It recreates the
0.2.5 failure mode and cannot satisfy immediate version verification.

### Invoke the remote installer script from `runx upgrade`

Rejected because normal native self-upgrade would gain a shell dependency,
duplicate download behavior, weaken JSON/event ownership, and execute a second
remote program. The scripts remain explicit recovery and direct-install tools.

### Patch only terminal output

Rejected because progress text cannot repair silent release lookup, incomplete
pagination, incorrect SemVer ordering, installer verification, or rollback
gaps. Output must describe a reliable structured operation rather than conceal
an unreliable one.
