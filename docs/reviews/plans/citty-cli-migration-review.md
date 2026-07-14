---
name: RunX Citty CLI Migration Plan Review
purpose: Verify that the full Citty migration plan is safe, complete, testable, and ready for execution.
description: Reviews command compatibility, safety invariants, native packaging, protected delivery, documentation, and Mirror release sequencing.
created: 2026-07-14
flags:
  - approved
tags:
  - reviews
  - plans
  - cli
keywords:
  - runx
  - citty
  - plan review
  - native executable
owner: runx-plan-reviews
---

# RunX Citty CLI Migration Plan Review

## Verdict

Ready for execution.

## Findings

- No blocker or high-severity findings remain.
- The plan explicitly removes both handwritten parsing and manual top-level
  routing, so the result cannot quietly retain a second parser behind Citty.
- The root-selector shorthand is treated as a compatibility path with known
  Citty commands taking precedence, preventing command names from being
  reinterpreted as manifest selectors.
- `-h` and `-v` are tested outside a configured project, directly covering the
  failure that motivated the migration.
- Manifest parsing, selectors, subprocess execution, upgrades, installers, and
  JSON payloads stay within their current domain modules; the migration is
  bounded to CLI ownership.
- No data, authentication, authorization, or cache migrations are involved.

## Safety and Compatibility Review

- The no-spawn invariant covers listing, describing, checking, help, version,
  and dry runs.
- Missing arguments and unknown options are usage failures, while RunX domain
  failures preserve their existing error contract.
- Every public command, nested route, alias, positional selector, and applicable
  option has an explicit test requirement.
- Citty must bundle into both npm output and Bun native executables without a
  Node.js runtime dependency.

## Delivery and Release Review

- Complete Bun and native validation precedes protected-branch delivery.
- The implementation reaches `main` before Mirror plans and applies patch
  `0.2.4`.
- The Mirror-managed release commit reaches protected `main` before its tag is
  pushed.
- The plan preserves the `production` environment's manual approval gate and
  observes npm and GitHub distribution after approval.

## Documentation and TODO Review

- README, detailed docs, the bundled RunX skill, and affected XDocs descriptors
  are named explicitly.
- The work is authorized for continuous execution in this session, so a second
  TODO artifact would duplicate the decision, plan, and review without adding a
  durable scheduling boundary.

## First Executable Unit

Execute Unit 1 in
[RunX Citty CLI Migration Plan](../../plans/citty-cli-migration.md): add Citty
with Bun and verify its installed command API before replacing production
routing.

## References

- [RunX Citty CLI Migration Plan](../../plans/citty-cli-migration.md)
- [Full Citty CLI Migration](../../decisions/citty-cli-migration.md)
