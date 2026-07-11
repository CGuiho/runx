---
name: RunX Alpha Command Catalog Requirements
purpose: Define the approved first-release outcome and acceptance criteria for RunX.
description: Captures the user-visible command catalog, safety, agent, installer, and CI requirements.
created: 2026-07-12
flags:
  - approved
tags:
  - requirements
  - cli
keywords:
  - runx
  - runx.yaml
  - command catalog
owner: runx-requirements
---

# RunX Alpha Command Catalog Requirements

## Outcome

RunX must provide a language-agnostic local command catalog in `runx.yaml` so
project operations no longer need to be hidden in `package.json` scripts.

## Acceptance Criteria

- `runx` shows a home page and usage; `runx list` lists catalog commands.
- Each command has a stable UID, group-scoped ID, summary, description, and
  shell command; validation rejects malformed or ambiguous manifests.
- `runx run` and `runx r` execute a selector, while `--dry-run` is read-only.
- Text and JSON inspection output support humans and agents.
- `confirm: always` blocks execution until `--yes` is supplied.
- The CLI provides help, help tree, help docs, native upgrade, uninstall, and
  skill-install commands.
- Windows and macOS/Linux direct-install scripts, CI, MIT licensing, and a
  bundled `guiho-s-runx` skill are present.

## Non-goals

The alpha does not publish to npm, orchestrate task graphs, manage secrets,
run commands remotely, or perform automatic package-script migration.
