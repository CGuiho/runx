---
owner: runx-todo
tags:
  - installer
  - linux
  - github
keywords:
  - issue 20
  - latest release
  - install.sh
---

# Resolve The Latest RunX Bash Install

## Outcome

The Bash installer installs the latest stable RunX release on Linux and macOS
without parsing a redirected release tag whose package name contains a slash.

## External Tracker

- [CGuiho/runx#20](https://github.com/CGuiho/runx/issues/20)

## Acceptance Signals

- `latest` assets resolve through GitHub's stable
  `releases/latest/download/<asset>` endpoint.
- Exact stable and prerelease versions continue to use encoded RunX tags.
- Binary, skill, and instruction assets use the same release selector.
- Bash syntax, piped startup, exact-version, and executable-verification tests
  pass.
- The public installer succeeds against the released patch before issue
  closure.

## Watch-Outs

- Do not recover a latest tag from the final redirect URL segment. The RunX tag
  contains `@guiho/runx@...`; the slash can become a path separator and discard
  the package scope.
- Preserve transactional replacement, Markdown validation, PATH setup, dual
  skill installation, and final executable verification.
