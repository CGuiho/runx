# Handoff — RunX Smart agent artifacts and docs (Unit 3)

Mode: dnd

## Task identity

- GitHub issue: https://github.com/CGuiho/runx/issues/66
- Task: `30. RunX Smart Command and rxs Launcher`
- Project: GUIHO #2, Component `runx`, Status `In Progress` (already set)
- Spec: `docs/todo/2026-10-22-runx-smart-and-rxs-launcher.md`

## Repository and required reading

- Repository: `/root/swe/runx` (Go + Cobra, branch `main`).
- Read `AGENTS.md`, then
  `/root/swe/guiho/conventions/guiho-convention-0001-cli.md` (README
  installation/uninstallation sections, required agent artifacts, instruction
  block boundaries) and `/root/swe/guiho/conventions/guiho-convention-0002-agent-artifacts.md`.
- Load skills: `guiho-s-0001-guiho`, `guiho-s-0004-working-with-cg`,
  `guiho-s-0032-git-commit`, `guiho-s-0016-docs`, `guiho-s-xdocs`,
  `guiho-s-runx`.

## Model decision

- Provider/model: `xiaomi/mimo-v2.6-pro`, provider-default thinking, no variant.

## Exclusive paths

You own only:

- `README.md`
- `DOCS.md`
- `CHANGELOG.md` (add the `0.17.0` section only; never edit released sections
  and never touch version fields)
- `skills/guiho-s-runx/SKILL.md`
- `embed/skills/guiho-s-runx.SKILL.md`
- `prompts/guiho-i-runx.md`
- `runx.yaml`
- Owned XDocs descriptors for those paths if their `documents` metadata changes
  (`runx.xdocs.md`, `skills/**/*.xdocs.md`, `prompts/**/*.xdocs.md`)

Do not touch `cmd/`, `devops/`, `docs/todo/`, `docs/plans/`, or any Go source.
A sibling worker owns `cmd/smart.go` and `cmd/rxs/`; a second sibling owns
`devops/` and the installers.

## Required content

- Document `runx smart` as the primary ergonomic entry point: bare
  `runx smart` lists, `runx smart <selector> ...` runs, flags before the
  selector are RunX options, everything after the selector belongs to the child.
- Document `rxs` as the short installed launcher: same list/run behavior, plus
  `-v`/`--version`, `-h`/`--help`, `--help-tree`, `--help-docs`, and
  `--help-tree-depth` entry-point parity with `runx`. State explicitly that
  `rx` is not a RunX command because it collides with an unrelated Linux tool.
- Add a concise quick-start example and mention the installed path
  `$HOME/.guiho/bin/rxs`.
- Keep the README installation and uninstallation sections exactly Compliant
  with Convention 0001 (three install commands, three uninstall commands,
  version verification). Do not restructure them.
- `skills/guiho-s-runx/SKILL.md` and `embed/skills/guiho-s-runx.SKILL.md` must
  stay aligned: the embedded copy is the runtime projection of the same skill
  and both are published. Update both with identical guidance; do not duplicate
  full command help.
- `prompts/guiho-i-runx.md`: add one bounded sentence that `runx smart`/`rxs`
  exist for catalog listing and execution; keep the block bounded and do not
  duplicate the skill.
- `runx.yaml`: update the `build-binaries` description to reflect eight payloads
  plus two launchers per platform; add a stable-UID catalog entry for
  `go run devops/verify-release-assets.go` (for example
  `runx-tooling-verify-release-assets`, id `verify-release-assets`) if it is
  not already catalogued. Keep `runx check` passing.
- `CHANGELOG.md`: add the `## 0.17.0 - 2026-10-22` section with `### Added` for
  `runx smart` and `rxs-launcher` release assets, and `### Changed` for the
  `rx` rename and dual-launcher install/upgrade/uninstall. Match the format the
  release-notes extractor expects (`go run devops/extract-release-notes.go`).

## Validation you must run and report

- `runx check --format json` (must pass with the catalog changes)
- `xdocs meta <scope> --documents --strict`, `xdocs tree`, `xdocs doctor` for
  every descriptor you change (use the installed `xdocs` CLI; if it is
  unavailable, state that exactly and validate the YAML/links manually)
- Markdown copyright header present at the earliest parser-safe position in
  every modified Markdown file
- `go run devops/extract-release-notes.go --version 0.17.0` writes
  `release-notes.md` with the new section (delete the generated file afterwards
  and do not commit it)

## Commit authority

- Commit each coherent unit with `guiho-s-0032-git-commit` (explicit paths).
- Never push. The parent reviews and pushes.
- If a Git index lock collides with a sibling worker, wait briefly and retry.
- Report commit hashes, commands run, and residual risk.
