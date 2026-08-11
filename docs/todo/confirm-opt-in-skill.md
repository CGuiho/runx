---
name: Make RunX Confirmation Opt-In In The Agent Skill
purpose: Define the expected outcome, constraints, and completion signals for TODO task 19.
description: Records the explicit confirmation-policy correction for the canonical and embedded RunX agent skills and its patch-release handoff.
created: 2026-08-11
flags:
  - testing
tags:
  - todo
  - agents
  - runx
keywords:
  - confirm
  - opt-in
  - guiho-s-runx
  - 0.11.1
owner: runx-todo
---

# Make RunX Confirmation Opt-In In The Agent Skill

## Todo Index

- Task: `19. Make RunX confirmation opt-in in the agent skill`
- Status: testing
- Index: [TODO.md](../../TODO.md)
- Parent TODO: [GUIHO root TODO](../../../guiho/TODO.md)
- Implementation notes: [confirm-opt-in-skill-implementation.md](./confirm-opt-in-skill-implementation.md)

## Plan Unit

`runx-confirm-opt-in-skill` is the single execution unit. A separate planning
phase is unnecessary because the Commander-in-Chief supplied the exact narrow
policy, owned skill copies, patch version, and release boundary in the request.

## Outcome

The canonical and embedded `guiho-s-runx` skills omit the manifest `confirm`
field by default and add it only when the user explicitly requests confirmation
behavior for that specific command. Only `confirm: never` and `confirm: always`
are supported; omission resolves to `never`.

## Scope

### In scope

- Update `skills/guiho-s-runx/SKILL.md` and
  `embed/skills/guiho-s-runx.SKILL.md` with the explicit opt-in policy.
- Remove the proactive instruction that inferred `confirm: always` for
  destructive or production-impacting commands.
- Keep existing authorization guidance for commands that already declare
  `confirm: always`.
- Bump skill metadata from `0.4.0` to `0.4.1` in both copies.
- Add focused regression coverage for the embedded policy.
- Update the owning XDocs descriptor and Unreleased changelog entry.
- Prepare patch release `0.11.1` for the post-merge Mirror release workflow.

### Out of scope

- Changes to RunX manifest parsing or execution behavior.
- Changes to legacy `source/` TypeScript references.
- Mirror version application, tag creation, GitHub Release publication, or
  production mutation from this feature branch.

## Acceptance Signals

- Both skill copies state that `confirm` is omitted by default and only added
  with the user's explicit request for a supported value.
- Neither skill contains the old proactive `confirm: always` instruction.
- Both skill copies retain the explicit authorization boundary for existing
  `confirm: always` commands.
- Embedded-skill regression coverage fails if the opt-in wording disappears or
  the proactive instruction returns.
- Focused tests, the full Go suite, vet, build, strict XDocs checks, and
  `mirror version plan patch` support a 0.11.1 handoff.

## Watch-outs

- Keep canonical and embedded skill copies semantically aligned.
- Do not change the runtime meaning of manifest `confirm`; omission already
  resolves to `never` in the Go parser.
- Do not add confirmation to a command merely because it is destructive,
  release-related, deployment-related, migration-related, or production-impacting.
- Preserve unrelated user changes in the original dirty checkout.

## Before Starting

- Confirm the dedicated branch is based on synchronized `main` at `8beb28d`.
- Read `xdocs.config.toml` and respect its `ai.mode = "auto"` setting.
- Inspect `mirror.yaml` and plan, but do not apply, the patch version.

## While Working

- Use explicit-path, smallest-coherent commits and keep edits within the owned
  skill, test, documentation, and lifecycle paths.
- Record a question ledger only if material uncertainty appears; none is
  currently expected because the user supplied an exact contract.

## After Finishing

- Move task 19 to `testing` before final validation.
- Push this branch and open a non-draft PR targeting `main`.
- Hand the exact PR head to independent implementation review and validation.
- Leave merge, Mirror apply, tag creation, and publication to the integrator.

## Mirror Decision

This is a compatible documentation and agent-policy correction, so the
appropriate transition is a patch bump to `0.11.1`. Mirror application,
version commit, tag, and GitHub Release are deferred until this PR is merged
and independently validated.
