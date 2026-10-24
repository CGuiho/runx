# Hand-off Task Brief — unit `skill-section` (2026-10-24)

#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

## Goal

Add one durable CG-preference section to the RunX agent skill `guiho-s-runx`
teaching agents the exact `runx.yaml` catalog layout CG wants, in both
repository copies and in the installed global projection.

## Context

- Repository and working directory: `/root/swe/runx`, branch `main`.
- Task identity: GitHub issue <https://github.com/CGuiho/runx/issues/68> on
  GUIHO Project #2 (`CGuiho` owner), Component `runx`, Status `In Progress` —
  created and read back by the parent. Do not create task records; do not edit
  `TODO.md`.
- Durable source of truth: [docs/todo/runx-catalog-formatting-preference.md](../../../todo/runx-catalog-formatting-preference.md).
- Files this unit owns (write ownership is exactly these three):
  1. `skills/guiho-s-runx/SKILL.md`
  2. `embed/skills/guiho-s-runx.SKILL.md`
  3. `~/.agents/skills/guiho-s-runx/SKILL.md` (installed projection)
- Placement: new section `## Compact Catalog Layout (CG preference, 2026-10-24)`
  immediately after `## Maintain Catalogs` and before `## Agent And Upgrade
  Commands`. No preference section exists in the repository copies today: add
  new text, never merge or rewrite unrelated text.
- Pre-existing content drift between the `skills/` and `embed/` copies in other
  sections is OUT of scope. Do not reconcile it and do not report it as your
  change.

## Skills To Load

- Tier-1 Essentials, always: `guiho-s-0001-guiho`, `guiho-s-0004-working-with-cg`,
  `guiho-s-0040-explorer`, `guiho-s-0032-git-commit` (in `~/.agents/skills/`).
- `guiho-s-0016-docs` (Markdown, copyright, frontmatter rules) and
  `guiho-s-runx` (the artifact being edited).
- GUIHO conventions in `/root/swe/guiho/conventions/`, at minimum
  `guiho-convention-0001-cli.md` (MANDATORY before any work in this
  repository), `guiho-convention-0002-agent-artifacts.md` (skill authoring),
  `guiho-convention-0007-models.md`.
- This repository's root `AGENTS.md`, `## Mandume` section.

## Requirements (CG's 2026-10-24 preference, transcription-reconciled)

1. One command entry per line; at most two lines per command.
2. Key order inside each entry: `uid`, `id`, `command`, then the remaining keys
   (`summary`, `description`, then any others).
3. Every key starts at a fixed column shared by all entries on its line — uid
   lane, id lane, command lane, summary lane, description lane — padded with
   spaces so the columns align like table columns and each field reads straight
   down ("all the UIDs in the same column").
4. Two-line split rule: when an entry carries long `summary`/`description`
   values, split at the same key boundary for every entry of that block —
   `uid`, `id`, `command` on line 1; `summary`, `description` on line 2 — with
   each line's lanes aligned and the continuation indented five spaces under
   the entry. Entries without long trailing keys stay on one line.
5. Spaces only, never tabs; no trailing whitespace; the last key on a line is
   ragged right and unpadded.
6. Never change values, quoting semantics, keys, or the manifest schema. Every
   example must parse as YAML and `runx check --format json` must still pass.
7. The section states that the layout governs how agents write and update
   `runx.yaml` catalogs AND how they show catalogs to CG, and that the same
   lane discipline helps other small structured YAML lists CG reads.

### Required example block (canonical, two-line form)

```yaml
commands:
  - {uid: dev,               id: dev,               command: "bun run dev",
     summary: Inspect development startup,         description: Existing dev includes cleanup and opaque configuration injection.}
  - {uid: typecheck,         id: typecheck,         command: "bun run typecheck",
     summary: Inspect route typegen and typecheck, description: Existing typecheck generates route types before TypeScript checking.}
  - {uid: compile-linux-x64, id: compile-linux-x64, command: "bun run elysia:compile-linux-x64",
     summary: Inspect Linux x64 compilation,      description: Existing Bun Linux x64 server compilation script.}
```

Lane targets (1-based columns): line 1 `uid:` col 6, `id:` col 30, `command:`
col 53; line 2 `summary:` col 6, `description:` col 52. Verify mechanically and
correct padding if any column is off.

### Required one-line form (short entries)

```yaml
commands:
  - {uid: dev, id: dev, command: "bun run dev"}
  - {uid: tc,  id: tc,  command: "bun run _tc"}
```

## Constraints

- Mode: dnd. Never ask questions. Resolve doubts with the safest reversible
  choice and record them in the result note.
- Model: `vercel/xiaomi/mimo-v2.6-pro#high` (CG conversation order 2026-10-24,
  live-catalog variant `high`). You are already running on it; do not switch.
- Write ownership: the three named files plus your result note. Nothing else.
  No `runx.yaml` manifest edits (including this repository's), no code, tests,
  release or versioning changes.
- Secrets boundary: never read, print, or commit environment files or key
  material.
- No production mutation. No `git push` — commit your owned unit locally and
  leave delivery to the parent.
- Keep every existing GUIHO copyright header intact; if you create Markdown,
  put the header at the earliest parser-safe position.

## Output Contract

- Modify: `skills/guiho-s-runx/SKILL.md`, `embed/skills/guiho-s-runx.SKILL.md`,
  `~/.agents/skills/guiho-s-runx/SKILL.md` — the new section text identical in
  all three.
- Result note: `docs/plans/2026-10-24-runx-catalog-formatting/execution/results/2026-10-24-skill-section-general.md`
  containing the verification transcript (mechanical column check, YAML parse
  of both example blocks, `runx check --format json` result), your commit hash,
  and every interpretation or deviation.
- Acceptance: the signals in
  [docs/todo/runx-catalog-formatting-preference.md](../../../todo/runx-catalog-formatting-preference.md).
