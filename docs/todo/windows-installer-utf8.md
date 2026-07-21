---
name: Preserve UTF-8 During Windows Installation
purpose: Define the completion contract for TODO task 9.
description: Requires the PowerShell installer and native maintenance path to preserve Unicode and converge damaged or duplicate RunX instruction blocks.
created: 2026-07-21
flags:
  - testing
tags:
  - installer
  - windows
  - reliability
keywords:
  - UTF-8
  - mojibake
  - AGENTS.md
  - idempotence
owner: runx-todo
---

# Preserve UTF-8 During Windows Installation

## Outcome

Installing RunX from Windows PowerShell preserves all existing UTF-8 guidance
outside the managed RunX block and leaves exactly one canonical RunX block.

## Acceptance Signals

- Existing Unicode text is decoded as strict UTF-8 and written without a BOM.
- A second reconciliation produces byte-identical output.
- Canonical, legacy, mojibake, and duplicate RunX blocks converge to one block.
- The installed-binary version probe cannot spawn either background worker.
- Native agent maintenance also repairs the mojibake marker emitted by 0.5.3.
- Focused tests, the full suite, typecheck, XDocs, CI, release assets, and a
  public `irm ... | iex` installation pass.

## Exclusions

- Reinterpreting arbitrary user-authored mojibake outside RunX managed blocks.
- Changes to CPU-worker TTL, lease, deadline, or recovery behavior.

## References

- [TODO.md](../../TODO.md)
- [Validation](../validation/windows-installer-utf8.md)
