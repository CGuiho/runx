---
name: guiho-s-runx
description: Use when inspecting, validating, documenting, or safely executing a RunX runx.yaml command catalog.
purpose: Teach agents the supported RunX catalog inspection and execution workflow.
created: 2026-07-18
flags:
  - bundled
tags:
  - runx
  - cli
keywords:
  - runx.yaml
  - command catalog
  - dry run
owner: guiho-s-runx
metadata:
  version: "0.3.0"
---

# GUIHO RunX

## Inspect Before Execution

1. Run `runx check --format json`.
2. Run `runx list --format json`.
3. Prefer a stable UID for automation.
4. Run `runx describe <uid>` before unfamiliar work.
5. Run `runx run --dry-run <uid>` before any mutation or high-impact command.

RunX manifests are trusted executable code. A group name is not a safety
boundary. Never add `--yes` unless the developer explicitly authorizes the
specific confirmation-gated command.

## Configuration

RunX resolves YAML only:

1. `--config <path>`;
2. effective-cwd `runx.yaml`;
3. `~/.guiho/runx/runx.yaml`.

It does not search parent directories.

## Execute

Use only:

```text
runx run [RunX options] <uid> [--] <child arguments...>
```

Listing, describing, checking, help, agent operations, and dry runs must never
execute a manifest command. Preserve the child command's exact exit code.
RunX options such as `--dry-run`, `--yes`, `--cwd`, and `--format` belong before
the selector. Every token after the selector is forwarded to the child without
being interpreted as a RunX flag.

## Maintain Catalogs

- Keep `uid`, `id`, `group`, `summary`, `description`, and `command`.
- Keep UIDs stable and never reuse one for materially different behavior.
- Use `confirm: always` for destructive, release, deployment, migration, and
  production-impacting operations.
- Do not place secrets in `runx.yaml`.
- Run `runx check` after every change.

## Agent And Upgrade Commands

```text
runx agent skill list
runx agent skill show guiho-s-runx
runx agent instruction show
runx agent prompt list --names
runx agent prompt show guiho-i-runx
runx upgrade check
runx upgrade list
runx upgrade --dry-run
```

Use `--help`, `--help-tree`, or `--help-docs` at any command scope for the
current executable contract.

## Automatic Agent Maintenance

Ordinary RunX commands schedule a silent, non-blocking worker that keeps the
bundled skill current in both global agent-tool directories and reconciles one
compact managed block in the nearest `AGENTS.md`. A current installation is not
rewritten. Automatic failures never fail or pollute the foreground command.

Explicit `runx agent ...` commands remain the manual repair and local-scope
interface. Explicit agent-resource removal and `runx uninstall` do not schedule
automatic reinstallation.
