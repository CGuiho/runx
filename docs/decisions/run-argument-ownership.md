---
name: Run Command Argument Ownership
purpose: Record how RunX separates its own options from child command arguments
description: Defines the approved selector boundary, optional delimiter, shell-safe forwarding, dry-run contract, and confirmation behavior.
created: 2026-07-22
flags:
  - decision
  - approved
tags:
  - decision
  - cli
keywords:
  - runx run
  - arguments
  - shell safety
owner: runx-decisions
---

# Run Command Argument Ownership

## Decision

For `runx run`, RunX owns tokens before the selector and the selected command
owns tokens after it. An immediate `--` after the selector is a delimiter and
is not forwarded. RunX never reinterprets a child `-v`, `--help`, `--yes`,
`--dry-run`, `--format`, or other leading-dash value.

Forwarded arguments travel as an immutable array. Shell adapters use positional
arguments or fixed environment-backed expansion; they never concatenate raw
child values into executable source.

## Consequences

- RunX options such as `--yes` and `--dry-run` must precede the selector.
- Existing examples using post-selector RunX options are updated without an
  alias or compatibility parser.
- Text and JSON dry runs expose the forwarded array separately from the trusted
  manifest command.
