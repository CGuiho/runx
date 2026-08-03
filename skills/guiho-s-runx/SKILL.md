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
  version: "0.4.0"
---

# GUIHO RunX

## Inspect Before Execution

1. Run `runx check --format json`.
2. Run `runx list --format json`.
3. Prefer a stable UID for automation. The numeric `IDX` from the current
   `runx list` output is available for interactive use.
4. Run `runx describe <uid-or-selector-or-index>` before unfamiliar work.
5. Run `runx run --dry-run <uid-or-selector-or-index>` before any mutation or
   high-impact command.

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
runx run [RunX options] <uid-or-selector-or-index> [--] <child arguments...>
```

Listing, describing, checking, help, agent operations, and dry runs must never
execute a manifest command. Preserve the child command's exact exit code.
Exact UID, selector, or unique-ID matches take precedence over numeric `IDX`
fallback. Numeric indexes belong to the current resolved listing and must not
replace stable UIDs in automation.
RunX options such as `--dry-run`, `--yes`, `--cwd`, and `--format` belong before
the selector. Every token after the selector is forwarded to the child without
being interpreted as a RunX flag.

For `confirm: always`, an interactive human receives a safe-default prompt and
an exact retry such as `runx run --yes cli-compile-host`. Noninteractive and
JSON runs fail closed with that retry instead of waiting for input. Treat an
interactive `y` or `yes` response and a retry with `--yes` as equivalent
authorization boundaries; never provide either without the developer's
explicit authorization for that command.

## Maintain Catalogs

- Use manifest v2 with one top-level `namespace`; never add legacy `project` or
  top-level `groups`.
- Keep commands and groups in the recursive `commands` list. Command leaves use
  `uid` and `id`; groups use `group` and exactly one of nested `commands` or a
  `runx` child reference.
- Use only relative paths or full HTTPS GitHub blob/raw URLs for `runx` and
  `parent`. A mounted child must declare the exact reciprocal parent. The mount
  group name may rename the child namespace.
- Keep sibling command/group names, global UIDs, and full selectors unique.
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

Bare `runx` synchronously ensures this embedded skill is installed in both
global agent-tool directories before showing the welcome. Inside a Git
repository it also reconciles the bounded RunX instruction block at the
repository root: both `AGENTS.md` and `CLAUDE.md` when both exist, the one that
exists, or a new `AGENTS.md`. Unmanaged content and line endings are preserved;
malformed markers fail safely. Help, version, agent-management, uninstall, and
non-repository paths do not perform repository instruction bootstrap.

Ordinary RunX commands schedule a silent, non-blocking worker that keeps the
bundled skill current and reconciles the same repository-root instruction
targets. A current installation is not rewritten. Background failures never
fail or pollute the foreground command.

Explicit `runx agent ...` commands remain the manual repair and local-scope
interface. Explicit agent-resource removal and `runx uninstall` do not schedule
automatic reinstallation.
