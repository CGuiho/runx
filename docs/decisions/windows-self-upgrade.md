---
name: RunX Windows Self-Upgrade
purpose: Define synchronous and recoverable replacement of the running RunX executable on Windows.
description: Records the accepted Windows executable-swap, verification, rollback, cleanup, testing, and patch-release behavior.
created: 2026-07-14
flags:
  - accepted
  - implemented
tags:
  - decisions
  - cli
  - windows
keywords:
  - runx upgrade
  - windows self-upgrade
  - executable replacement
  - rollback
owner: runx-decisions
---

# RunX Windows Self-Upgrade

## Context

RunX 0.2.5 downloads a Windows update to `runx.exe.new`, starts a detached
`cmd.exe` process, reports `scheduled: true`, and exits without confirming that
the installed executable was replaced. The detached move can fail silently and
leave later invocations on the previous version.

Windows permits the running executable to be renamed while its image remains
mapped. A new executable can then be installed at the original path before the
current RunX process exits. Only deletion of the renamed, still-running image
must wait until the process releases it.

## Decision

`runx upgrade` on Windows will complete the executable swap synchronously:

1. Download the selected release asset to `runx.exe.new`.
2. Remove any stale backup that is not in use.
3. Rename the installed `runx.exe` to `runx.exe.old`.
4. Rename `runx.exe.new` to the canonical installed path.
5. Execute the newly installed path with `--version` and require the expected
   target version and a zero exit code.
6. If installation or verification fails, remove the failed replacement and
   restore `runx.exe.old` to the canonical path.
7. After successful verification, delete the backup immediately when possible;
   when Windows still locks the running image, schedule only that backup cleanup
   for after the current process exits.

The upgrade result reports `scheduled: false`. Deferred deletion of the old,
renamed image is cleanup and does not defer or weaken successful replacement.

## Error Handling

- A replacement failure must throw a concise `RunXError` and produce a non-zero
  CLI exit code.
- Rollback must be attempted whenever the installed executable was renamed but
  the new executable was not verified.
- If rollback also fails, the error must report both the replacement and
  recovery failures and identify the affected installed path.
- The `.new` file must be removed after failures and must not remain after a
  successful rename.
- A locked stale `.old` file must stop a new upgrade before the current
  executable is modified.

## Testing

Automated coverage must prove that:

- the configured Windows self path contains the downloaded replacement before
  `upgradeSelf()` resolves;
- successful Windows upgrades report `scheduled: false`;
- no `.new` file remains after success;
- a failed replacement restores the original executable;
- the equal-version no-op remains unchanged;
- `runx -h` continues to work outside a manifest directory.

The full release gate is `bun run typecheck`, `bun test`, `bun run build`,
`bun run binary`, `bun run binaries`, and strict XDocs validation.

## Release

After implementation and validation, use Mirror to plan and apply the next
patch. Update the changelog before applying the Mirror-managed version. Do not
hand-edit the package version or create the release tag manually.

## Acceptance Criteria

- When `runx upgrade` succeeds, the executable at the installed path is already
  the target release.
- A same-terminal or new-terminal invocation immediately reports the target
  version.
- The command never reports `scheduled: true` for upgrade replacement.
- Replacement and verification failures are visible and recover the old binary.
- Temporary update files are cleaned up safely.
- Issue #9 has automated Windows regression coverage and issue #1 remains
  covered by the Citty help tests.
