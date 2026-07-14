---
name: guiho-s-runx
purpose: Guide agents through safe inspection and execution of RunX command catalogs.
description: Use when inspecting, creating, validating, organizing, documenting, or running a RunX `runx.yaml` command catalog. This includes listing available project commands, choosing a command by UID or selector, reviewing command descriptions, dry runs, confirmation-gated execution, and agent-safe project automation.
created: 2026-07-12
owner: guiho-s-runx
flags: []
tags:
  - agents
  - runx
keywords:
  - runx
  - citty
  - command catalog
  - dry run
---

# GUIHO RunX

Use RunX as the command catalog for a project. It documents executable local
commands but does not make them inherently safe: a manifest is trusted code.

## Inspect Before Executing

1. From the intended project directory, run `runx check --format json`.
2. Run `runx list --format json` to inspect groups, UIDs, summaries, tags, and
   confirmation requirements.
3. Choose a stable `uid` for automation. An index is only for the current list
   and an unqualified ID can be ambiguous.
4. Run `runx describe <uid>` before an unfamiliar or high-impact command.
5. Run `runx run <uid> --dry-run` before a real execution when the command has
   filesystem, database, release, deployment, or production effects.

## Execute Safely

- Use `runx run <uid>` or the short alias `runx r <uid>`.
- `runx <selector>` is allowed as a human shorthand only when it cannot be
  confused with a built-in command; agents should prefer explicit `runx run`.
- If `confirm: always` is set, RunX requires `--yes`.
- Add `--yes` only after the user explicitly authorizes the specific command.
- Do not treat a group name such as `development` or `operations` as a safety
  guarantee. Read the command description and command text.
- Do not add secrets to `runx.yaml`. Use the project's existing environment or
  secret-management workflow.
- Use `runx <command> --help` for command-specific Citty usage. `runx -h` and
  `runx -v` are global operations and do not require a manifest.

## Maintain Manifests

- Keep `uid`, `id`, `group`, `summary`, `description`, and `command` present.
- Keep UIDs stable; never reuse a UID for a materially different command.
- Use a group for organization, a concise summary for listings, and a complete
  description for prerequisites, effects, and risk.
- Set `confirm: always` for destructive, release, deployment, migration, and
  production-impacting commands when confirmation is useful.
- Run `runx check` after every manifest change. It rejects unknown fields,
  duplicate UIDs/selectors, missing descriptions, invalid working directories,
  and unsupported shell choices.

## Useful Commands

```text
runx                         Show the RunX home page and usage.
runx -h                      Show generated command help.
runx -v                      Show the installed version.
runx list                    List the nearest manifest.
runx describe <uid>          Explain one command without execution.
runx run <uid> --dry-run     Inspect the execution plan.
runx r <uid>                 Run through the short alias.
runx agents install local    Install this skill into .agents/skills.
runx --help-tree             Show the command hierarchy.
runx --help-docs             Show manifest documentation guidance.
```
