---
name: RunX Catalog Formatting Preference
purpose: Define the outcome, constraints, and completion signals for task 31 in todo.md
description: Records CG's runx.yaml catalog layout preference in the RunX agent skill — one entry per line, at most two lines per command, uid/id/command first, and every key aligned in its own table column.
created: "2026-10-24T17:43:34Z"
flags: []
tags:
  - todo
  - skill
  - runx
  - formatting
keywords:
  - runx.yaml
  - catalog
  - formatting
  - alignment
  - guiho-s-runx
  - compact catalog layout
owner: runx-todo
---

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# RunX Catalog Formatting Preference

## Summary

CG wants `runx.yaml` command catalogs written and shown as a compact, table-like
grid: one command entry per line (at most two lines per command), keys ordered
`uid`, `id`, `command`, then the remaining keys (`summary`, `description`, ...),
and every key starting at a fixed column shared by all entries so each field
reads straight down its own lane. This is instruction level only: no schema
change, no manifest reformatting, no CLI behavior change.

## Todo Index

- Task: `31. Record CG runx.yaml Catalog Formatting Preference In The Agent Skill`
- Status: completed
- Index: [done.md](done.md)
- GitHub project item: [RunX #68](https://github.com/CGuiho/runx/issues/68), closed.
- GitHub component: `runx`
- Original index: [../../../todo.md](../../../todo.md)

## Outcome

The RunX agent skill carries a durable CG-preference section that tells agents
exactly how to lay out `runx.yaml` catalogs and catalog views, with a rendered
example whose columns are mechanically verified to align.

## Scope

### In scope

- [skills/guiho-s-runx/SKILL.md](../../../skills/guiho-s-runx/SKILL.md) - repository skill source.
- [embed/skills/guiho-s-runx.SKILL.md](../../../embed/skills/guiho-s-runx.SKILL.md) - runtime embedded copy, kept identical in the new section.
- The installed global projection `~/.agents/skills/guiho-s-runx/SKILL.md`,
  refreshed with the same section so the live skill matches.

### Out of scope

- Any `runx.yaml` manifest content or schema, including this repository's.
- CLI code, tests, or behavior.
- Reconciling the pre-existing content drift between the `skills/` and `embed/`
  copies outside the new section; that drift is reported, not merged.

## Acceptance Signals

- Both repository copies carry the identical new section, with the key order
  `uid`, `id`, `command`, then `summary`/`description` stated explicitly.
- The example block shows narrow and wide UIDs, and every `uid:`, `id:`,
  `command:`, `summary:` and `description:` key starts at the same column
  across example rows on its line, verified mechanically.
- The example parses as YAML, `runx check --format json` still passes, and no
  example value or quoting is altered.
- The one-command-per-line / at-most-two-lines budget and the two-line split
  rule are explicit in the section text.
- No tabs and no trailing whitespace anywhere in the edited files.

## Watch-outs

- The earlier "Compact Catalog View (CG preference, 2026-10-24)" wording exists
  only in OpenCode's session database and is not durable; it must be written
  into the repository sources, never assumed present.
- RunX automatic agent maintenance installs
  `embed/skills/guiho-s-runx.SKILL.md` into `~/.agents/skills/` on real
  executions, silently overwriting local edits to the installed copy. The
  durable home of the preference is the embed copy.
- Alignment is spaces-only padding outside quoted scalars; values, quoting and
  manifest semantics stay byte-exact.

## References

- [todo.md](../../../todo.md)
- [skills/guiho-s-runx/SKILL.md](../../../skills/guiho-s-runx/SKILL.md)
