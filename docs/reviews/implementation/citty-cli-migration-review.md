---
name: RunX Citty CLI Migration Implementation Review
purpose: Review the completed Citty migration against its accepted decision, executable plan, architecture, and public contract.
description: Records findings, acceptance checks, verification evidence, documentation alignment, and residual risk before protected delivery.
created: 2026-07-14
flags:
  - accepted
tags:
  - reviews
  - implementation
  - cli
keywords:
  - runx
  - citty
  - command routing
  - native executable
owner: runx-implementation-reviews
---

# RunX Citty CLI Migration Implementation Review

## Verdict

Accepted for protected delivery.

## Findings

No open blocker, high-severity, or medium-severity findings remain.

Two review findings were corrected before this verdict:

- The selector-compatible command map now has a null prototype, so valid
  selector names such as `constructor` cannot collide with inherited object
  properties.
- Root, catalog, run, agent-install, and self-management argument definitions
  are separated, so Citty help shows only applicable flags while the root still
  accepts options before or after commands.

## Acceptance Criteria Check

- Citty is the only argument parser, alias registry, command router, and
  ordinary usage renderer.
- `source/flags.ts` is removed and no RunX-owned token parsing loop remains.
- `-h`, `--help`, `-v`, and `--version` work outside a configured project.
- All catalog, nested agent, nested upgrade, uninstall, alias, and selector
  shorthand paths dispatch through Citty command definitions.
- Unknown options and missing positionals produce relevant usage without
  manifest discovery.
- Manifest, selector, execution, agent, upgrade, uninstall, text, and JSON
  domain behavior remains in the existing modules.
- Inspection, help, version, and dry-run paths do not execute a configured
  command.
- Citty bundles into the local native executable and all twelve release targets
  without adding a Node.js runtime requirement.
- The npm package declares Citty as a runtime dependency and has no
  `postinstall` script.

## Verification Evidence

- Typecheck, thirteen tests, library build, local native compile, twelve-target
  native matrix, npm pack dry-run, packed launcher, and strict XDocs checks
  passed.
- Tests cover root and nested help, short and long version/help flags, all
  public command routes, aliases, options before and after commands, JSON,
  safety gates, usage errors, agent installation, self-management routing, and
  selector prototype-name compatibility.

## Docs And TODO Check

- README, canonical docs, architecture, AGENTS instructions, bundled skill,
  decision, plan, and affected XDocs descriptors align with the implementation.
- No standalone TODO was needed because the accepted plan was executed in this
  authorized session.

## Residual Risk

The remaining risk is distribution-only: GitHub CI, protected merge, tag
creation, manual `production` environment approval, and public npm/GitHub
publication must still complete. Unit 5 owns that evidence.

## References

- [Citty CLI Migration Decision](../../decisions/citty-cli-migration.md)
- [Citty CLI Migration Plan](../../plans/citty-cli-migration.md)
- [Citty CLI Migration Validation](../../validation/citty-cli-migration.md)
