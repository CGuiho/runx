---
name: Run Command Argument Ownership
purpose: Record how RunX separates its own options from child command arguments
description: Defines exact and numeric selector resolution, the approved selector boundary, optional delimiter, shell-safe forwarding, dry-run contract, and confirmation behavior.
created: 2026-07-22
flags:
  - decision
  - approved
tags:
  - decision
  - cli
keywords:
  - runx run
  - numeric index
  - runx list
  - arguments
  - shell safety
owner: runx-decisions
---

# Run Command Argument Ownership

## Decision

For `runx run`, RunX owns tokens before the selector and the selected command
owns tokens after it. A selector may be an exact UID, canonical selector,
unique ID shorthand, or the canonical positive-decimal `IDX` from the current
resolved listing. Exact textual identities resolve before numeric-index
fallback. Numeric indexes are an interactive convenience; stable UIDs remain
the automation contract.

An immediate `--` after the selector is a delimiter and is not forwarded. RunX
never reinterprets a child `-v`, `--help`, `--yes`, `--dry-run`, `--format`, or
other leading-dash value.

Forwarded arguments travel as an immutable array. Shell adapters use positional
arguments or fixed environment-backed expansion; they never concatenate raw
child values into executable source.

A `confirm: always` command prompts only for interactive text input. The prompt
defaults to no and accepts only `y` or `yes`, case-insensitively. It also prints
an exact retry with `--yes` in RunX-owned position before the selector.
Noninteractive and JSON invocations never read a response or block; they fail
closed and report that same retry. The retry retains explicitly supplied RunX
options and uses `--` before preserved child arguments.

## Consequences

- RunX options such as `--yes` and `--dry-run` must precede the selector.
- Numeric index selection preserves the selected command's confirmation policy
  and never bypasses confirmation.
- `--yes` bypasses the prompt only when it appears before the selector.
- Enter, `n`, `no`, EOF, and unrecognized responses decline without spawning.
- Existing examples using post-selector RunX options are updated without an
  alias or compatibility parser.
- Text and JSON dry runs expose the forwarded array separately from the trusted
  manifest command.
