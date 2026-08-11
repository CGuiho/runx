---
name: RunX Windows Git Bash Automatic Shell Plan
purpose: Define the sealed implementation and patch-delivery unit for issue 47's Windows shell boundary.
description: Sequences caller-aware shell:auto resolution, deterministic regression tests, documentation, exact-head gates, protected integration, and a Mirror-managed patch release.
created: "2026-08-11"
flags:
  - approved
  - implementation-ready
tags:
  - plan
  - cli
  - windows
keywords:
  - Git Bash
  - shell auto
  - issue 47
  - gpt-5.6-luna
  - patch release
owner: runx-plans
---

# RunX Windows Git Bash Automatic Shell Plan

## Summary

Correct one Windows execution-boundary bug after the successful 0.12.0 reveal
release. GPT-5.6 Luna at maximum reasoning effort owns all executable-code,
regression-test, and implementation-side documentation changes. The root agent
owns this plan, independent exact-head review, validation, integration, and
release decisions.

## Evidence And Scope

- [Issue 47](https://github.com/CGuiho/runx/issues/47) contains the original
  Dotenvx symptom and the secret-safe reproduction.
- Released `runx reveal dev` prints the exact stored GUIHO Core command.
- A non-secret native Node probe receives `C:/GUIHO` through direct Git Bash,
  `/c/GUIHO` through RunX `shell: auto`, and `C:/GUIHO` through RunX
  `shell: bash` when RunX is launched from Git Bash.
- `pkg/executor.BuildShellExecution` currently maps every Windows `auto` value
  to `cmd` before considering caller context.

This unit changes only automatic shell selection. It must not parse, normalize,
or rewrite configured command text or paths, and it must not access any
environment/key file.

## Unit RX-SHELL-1: Caller-Aware Windows Auto Resolution

### Ownership And Isolation

- Base: `origin/main` at `107f3f69a71c05a877a9854fee063ddce2ecd0ce`.
- Branch: `codex/runx-git-bash-auto`.
- Worktree: `C:\GUIHO\runx\.temp\runx-git-bash-auto-47`.
- Implementation agent: GPT-5.6 Luna, maximum reasoning effort.
- PR target: protected `main`.

### Required Design

1. Keep `BuildShellExecution` as the public production entry and factor the
   platform, environment lookup, and executable lookup boundary into a private
   deterministic helper so Windows cases can be tested on all CI platforms.
2. For `shell: auto` on Windows, select Bash only when the inherited process
   environment proves an MSYS/Git Bash caller and executable lookup resolves
   that caller's Bash. Use the resolved executable for execution.
3. Reject or fall back from the Windows/WSL system `bash.exe` launcher; a stale
   marker must not redirect auto execution into WSL.
4. Fall back to the current `cmd.exe` transport whenever caller evidence or
   Bash resolution is absent. Non-Windows `auto` remains `sh`.
5. Explicit `cmd`, `powershell`, `bash`, and `sh` remain unchanged and bypass
   caller inference.
6. Do not change manifest schema, configured command strings, selector
   behavior, confirmation, dry runs, child-argument transport, exit mapping,
   or reveal output.

### Required Tests

- Windows auto plus a recognized Git Bash/MSYS environment and resolved Git
  Bash executable produces the Bash transport with that executable.
- Windows auto without the marker remains `cmd.exe`.
- Windows auto with a marker but missing Bash remains `cmd.exe`.
- Windows auto with a marker resolving only the Windows/WSL launcher remains
  `cmd.exe`.
- Explicit shell values remain authoritative under the same caller markers.
- Non-Windows auto remains `sh`.
- A behavior-based `/c/GUIHO` regression proves the intended transport without
  reading or constructing any secret-bearing path.
- Existing executor and full repository tests remain green.

### Owned Paths

- `pkg/executor/executor.go` and focused executor tests.
- `pkg/executor/executor.xdocs.md`.
- `DOCS.md` and `CHANGELOG.md` for the public shell contract and unreleased
  patch note.
- Canonical and embedded `guiho-s-runx` skill copies and their descriptor only
  if shell guidance changes.
- Task/implementation evidence and directly owning XDocs descriptors.

Do not edit generated `library/`, `bin/`, `bundle/`, or `vendor` outputs; legacy
TypeScript under `source/`; environment/key files; Mirror-managed version
fields; or the dirty original `main` checkout.

## Validation And Delivery

1. Format, tidy-check, focused/full tests, vet, and build in the isolated
   worktree using repository-local Go caches when Windows sandbox caches are
   unavailable.
2. Run strict XDocs metadata, tree, and doctor with warnings as errors.
3. Build and verify the standard 8-target/11-asset matrix for the next patch
   candidate.
4. Push a ready PR referencing issue 47 without closing it before released
   acceptance.
5. Root reviews and validates the exact PR head. Findings return to Luna and
   invalidate prior review/validation evidence.
6. Merge only after exact-head CI, review, validation, mergeability, and
   protection agree; then verify main reachability.
7. Inspect `mirror.yaml`, run `mirror version plan patch`, and apply only the
   expected `0.12.1` protected tag from a clean integrated main.
8. Verify the successful publication workflow, non-draft/non-prerelease GitHub
   Release, exactly 11 assets, independent checksum/native smoke, npm latest,
   and the non-secret Git Bash path probe before closing issue 47.

## Stop Conditions

- Stop before merge for a head mismatch, failed check, unresolved conflict,
  missing caller-detection test, or changed explicit-shell behavior.
- Stop before release for a dirty tree, unexpected Mirror plan, absent
  version-scoped changelog notes, or incomplete 11-asset verification.
- Never bypass protection, force-push, inspect secrets, or normalize command
  strings as a substitute for correct shell transport.
