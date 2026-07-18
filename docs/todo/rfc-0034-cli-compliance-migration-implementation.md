---
name: RunX RFC 0034 CLI Compliance Migration Implementation
purpose: Preserve the executed RX-01 through RX-16 implementation record.
description: Records the delivered module changes, breaking-contract decisions, test-hang diagnosis, and validation handoff for the RunX RFC migration.
created: 2026-07-18
flags:
  - completed
tags:
  - implementation
  - cli
keywords:
  - RFC 0034
  - RX-01
  - RX-16
  - RunX
owner: runx-todo
---

# RunX RFC 0034 CLI Compliance Migration Implementation

## Summary

RX-01 through RX-16 were implemented as one breaking pre-1.0 migration. RunX
now has one Citty tree, Bun-only core modules, TypeBox-decoded configuration,
cache and release boundaries, complete Developer Context help, the singular
agent namespace, transactional upgrades, complete installers, a Node-compatible
npm bootstrap, and the exact fourteen release assets.

## Implementation Map

- RX-02-RX-05: `path-utils.ts`, `storage.ts`, `configuration.ts`,
  `update-cache.ts`, release schemas, and shared types establish Bun-first I/O
  and TypeBox boundaries.
- RX-06-RX-08: `cli.ts` and `help.ts` own the exact Citty tree, exit mapping,
  startup lifecycle, Unicode tree help, depth validation, and Markdown help.
- RX-09: `agents.ts`, embedded resources, the bundled skill, and
  `prompts/guiho-i-runx.md` implement dual-tool skills, instruction blocks, and
  prompt discovery.
- RX-10-RX-11: CLI output, release pagination, upgrade selection, replacement,
  rollback, reconciliation, and cache behavior use stable contracts.
- RX-12-RX-14: both installers, `scripts/runx-bin.mjs`, the binary builder,
  asset verifier, and GitHub workflows enforce the installation and fourteen
  asset contracts.
- RX-15-RX-16: canonical docs, TODO, xdocs, review, validation, tests, builds,
  and release preparation were aligned to the executable.

## Test-Hang Diagnosis

The original unbounded suite hung in legacy process-oriented tests. The
individual CLI baseline passed 10/10, proving the CLI was not persistently
stuck. The obsolete tests were replaced or bounded; the final complete suite
runs 40 tests and exits normally in about 30 seconds.

## Independent Audit Corrections

The final independent RFC audit found three public-surface defects that the
first validation suite did not exercise:

- the `upgrade` group hid its apply flags and allowed root `--version` handling
  to intercept `runx upgrade --version <version>`;
- no-argument startup skipped the cached update notice;
- text-mode `agent prompt list --names` rendered a JSON array instead of raw
  names.

The correction keeps `upgrade check` and `upgrade list`, makes the public
`upgrade` command own and execute all RFC flags directly, gates root version
handling to a one-token root invocation, reads the local cache before the
no-argument banner, and emits one raw prompt name per line. Focused tests invoke
the real Citty routing and output paths. The ignored `bin/runx.exe` left by a
single-binary build was removed before exact-asset verification; the strict
fourteen-name comparison remains unchanged.

## References

- [Plan](../plans/rfc-0034-cli-compliance-migration.md)
- [Implementation review](../reviews/implementation/rfc-0034-cli-compliance-migration-review.md)
- [Validation](../validation/rfc-0034-cli-compliance-migration.md)
