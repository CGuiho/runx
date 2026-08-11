---
name: Add RunX Reveal
purpose: Define the expected outcome, constraints, and completion signals for TODO task 21.
description: Specifies exact command revelation, selector parity, no-execution safety, documentation, validation, and minor-release acceptance.
created: "2026-08-11"
flags:
  - completed
tags:
  - todo
  - cli
  - execution
keywords:
  - runx reveal
  - issue 47
  - selector
  - no spawn
owner: runx-todo-done
---

# Add RunX Reveal

## Todo Index

- Task: `21. Add RunX Reveal`
- Status: completed
- Index: [TODO.md](../../../TODO.md)
- Parent coordination: [GUIHO root TODO](../../../../guiho/TODO.md)

## Outcome

RunX provides `runx reveal <selector>` to resolve one catalog command exactly
as existing inspection/execution commands do and print its exact manifest
command for the user to copy and run, without spawning or confirming it.

## Scope

### In scope

- Public Cobra `reveal` command and live Developer Context help.
- Exact global UID, canonical/full selector, unique ID shorthand, and numeric
  index resolution through the existing catalog resolver.
- Existing `--config`, `--cwd`, and optional text diagnostic behavior needed to
  load the same catalog as `runx run` while keeping stdout copyable.
- Exact command plus trailing newline on successful stdout.
- No-spawn, no-confirmation, invalid-selector, help, and selector-parity tests.
- README, CLI reference, prompt, canonical/bundled skill, changelog, and XDocs
  alignment.
- Independent exact-head review, validation, protected integration, Mirror
  minor version application, 11-asset release publication, and release smoke.

### Out of scope

- Reading or validating environment/key files, changing secret material, or
  printing secret values.
- Changing `runx run` shell adapters or automatic Windows shell selection.
- Path conversion between Git Bash, `cmd.exe`, and PowerShell.
- Forwarded child-argument rendering in the first reveal contract.
- Production deployment, promotion, traffic, DNS, database, or secret changes.

## Acceptance Signals

- `runx reveal <uid>` prints exactly the stored command plus `\n`.
- Numeric index, canonical/full selector, and unique shorthand select the same
  command as `runx run --dry-run` without executing it.
- Reveal of a `confirm: always` command succeeds without input, prompting, or
  `--yes` because it is non-executing.
- Invalid selectors preserve the established configuration/selection exit
  behavior and do not spawn.
- Help-tree and Markdown help include the new command with no short alias.
- Inspection safety tests prove `reveal` does not create a configured marker.
- Go formatting, tidy cleanliness, tests, vet, build, eight cross-builds,
  11-asset verification, XDocs checks, Mirror configuration, and minor plan all
  pass.
- The merged release is published as the next minor after `0.11.0`, with the
  exact 11 assets and a downloaded Windows AMD64 reveal smoke.

## Related Files

- [Feature brainstorm](../../brainstorm/runx-reveal.md)
- [Implementation plan](../../plans/runx-reveal.md)
- [Plan review](../../reviews/plans/runx-reveal-review.md)
- [Question ledger](../../questions/runx-reveal/plan-unit-1.md)
- [Implementation note](runx-reveal-implementation.md)
- [Run argument ownership](../../decisions/run-argument-ownership.md)

## Dependencies And Context

- Approved base: `8beb28d1690c050bc7345cdbf77da2bb143909e9` from current `origin/main`.
- Dedicated branch: `codex/runx-reveal`.
- Isolated worktree: `C:\GUIHO\runx\.temp\runx-reveal-47`.
- GitHub issue: [CGuiho/runx#47](https://github.com/CGuiho/runx/issues/47).
- Open PR 46 may overlap skill, changelog, TODO, or XDocs paths; never revert its
  confirmation-policy work if it reaches `main` before integration.

## Watch-outs

- Preserve stdout as raw copyable command text; send optional diagnostics to
  stderr.
- Reuse catalog resolution rather than creating a second selector parser.
- Never call the executor, confirmation path, or shell adapter from reveal.
- Preserve existing `runx run`, dry-run, child-argument, and exit-code behavior.
- Leave legacy TypeScript reference files and generated outputs untouched.

## Lifecycle Waivers

- Product requirements are waived because the user supplied explicit behavior
  and acceptance criteria.
- Technical architecture and architecture review are waived because this is a
  focused addition to the existing Cobra/catalog boundary with no schema,
  authentication, cache, cloud, API, or cross-repository design change.
- The focused plan review remains required because execution is delegated.

## After Finishing

- Move canonical task state to `testing` before final exact-head review and
  validation.
- Merge only after review, validation, CI, protection, and mergeability gates
  pass for the same head.
- Apply the clean Mirror minor transition only from integrated `main`, verify
  the protected publication workflow and exact release assets, then archive the
  task with durable evidence.

## Mirror Decision

`runx reveal` is a new public capability, so the required target is `minor`.
Mirror, not manual version edits or tags, owns the transition.

## Final Acceptance

PR #50 and evidence PR #53 reached protected main. Mirror published
`@guiho/runx/v0.12.0` with exactly 11 assets through successful workflow
31474091754. The downloaded Windows AMD64 binary matched its checksum, reported
0.12.0, and revealed the exact GUIHO Core command without execution.
