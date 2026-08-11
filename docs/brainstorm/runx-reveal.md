---
name: RunX Reveal Feature Brainstorm
purpose: Preserve the requested RunX reveal behavior and the execution failure that motivated it before implementation.
description: Captures the existing RunX run behavior, desired copyable command output, selector requirements, secret boundary, and deferred Windows shell investigation.
created: "2026-08-11"
flags:
  - accepted
tags:
  - brainstorm
  - cli
  - execution
keywords:
  - runx reveal
  - runx run
  - dotenvx
  - selector
  - numeric index
owner: runx-brainstorm
---

# RunX Reveal Feature Brainstorm

## Feature Intent And Problem

A command stored in `runx.yaml` can succeed when pasted into the user's current
Git Bash session yet fail when delegated through `runx run` on Windows. The
observed catalog command starts the application directly, while RunX reports a
Dotenvx decryption failure, a libuv closing-handle assertion, and child exit
code `3221226505`.

The immediate requested capability is an execution-free escape hatch:
`runx reveal <selector>` prints the complete configured command so the user can
copy it and run it in the current shell.

## Existing Behavior

- `runx run` resolves exact global UIDs, canonical selectors, unique
  unqualified ID shorthands, and then canonical positive-decimal indexes.
- `runx run` delegates execution through the configured shell adapter.
- `shell: auto` selects `cmd.exe` on Windows, whereas a command pasted into Git
  Bash is interpreted by the Git Bash/MSYS environment.
- `runx run --dry-run` exposes a structured plan but does not provide the
  requested one-line copyable command surface.

## Desired Behavior

- `runx reveal <selector>` uses the same catalog loading, composition, selector
  precedence, configuration, and effective-cwd rules as `runx run`.
- The selector may be a numeric index, global UID, canonical/full selector, or
  existing unique ID shorthand.
- Successful text output is the exact manifest command plus one trailing
  newline.
- Reveal never starts the configured process and never asks for confirmation.
- Reveal treats environment and key paths as opaque command text; it never
  reads, validates, decrypts, or exposes file contents.

## Confirmed Scope

- One new public Cobra command named `reveal`.
- Exact command-string output with no prefix, label, or execution metadata on
  stdout.
- Selector parity and no-spawn regression coverage.
- Live help, public documentation, bundled skill, prompt, XDocs, changelog,
  native builds, and minor-release coverage.

## Deferred Scope

- Fixing or changing `runx run` shell selection is a follow-up after the reveal
  release.
- Rewriting POSIX paths for `cmd.exe`, detecting the caller's interactive
  shell, or changing manifest `shell` defaults is not part of this feature.
- Rendering forwarded child arguments is deferred because a lossless copyable
  representation depends on the shell into which the user will paste it. The
  first version returns the exact configured command without transformation.

## Constraints And Risks

- A manifest is trusted executable code; reveal exposes it but must not run it.
- Stdout purity matters because the output is intended for direct copy and
  paste. Optional diagnostics may use stderr only.
- The existing dirty `main` checkout and open PR 46 are outside this branch;
  their changes must remain intact and any later overlap must be integrated
  without reverting them.
- Secret-bearing environment and key files remain human-only and must never be
  opened, listed, searched, validated, staged, or committed.

## Source References

- [GitHub issue 47](https://github.com/CGuiho/runx/issues/47)
- [Run argument ownership](../decisions/run-argument-ownership.md)
- [RunX task index](../../TODO.md)

## Handoff

The request is complete enough to skip separate product-requirements and
architecture phases: the user supplied the command name, selector contract,
output behavior, release type, and delivery sequence. The next durable artifact
is the focused implementation plan.
