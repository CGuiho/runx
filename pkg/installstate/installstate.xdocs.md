---
subject: runx-installstate
description: Canonical .guiho installation layout ownership — stable launcher path, immutable version directories, atomic current.json pointer, installed-artifacts ledger, and CLI-owned staging.
parent: runx-packages
children: []
files:
  paths.go: Resolves $HOME/.guiho roots via os.UserHomeDir plus filepath.Join for the launcher, versions, resources, pointer, ledger, and staging locations.
  pointer.go: Strict protocol-1 pointer validation, SemVer sanitization, and traversal rejection.
  ledger.go: Strict pointer and ownership-ledger decoding, including a narrow legacy UTF-8 BOM compatibility boundary for PowerShell-generated current.json files.
  atomic.go: Atomic file replacement and unique staging directories under the shared temp root.
  installstate_test.go: Covers canonical layout, pointer/ledger validation, legacy PowerShell BOM pointers, foreign-path rejection, atomic writes, and staging cleanup.
documents: {}
tags:
  - go
  - installation
keywords:
  - current.json
  - UTF-8 BOM compatibility
  - installed-artifacts.json
  - stable launcher
  - immutable payloads
flags: []
status: stable
---

Every lifecycle write is constrained to `$HOME/.guiho/runx/`, the single
launcher path under `$HOME/.guiho/bin/`, or a managed skill projection. Shared
`.guiho/`, `.guiho/bin/`, and `.guiho/.temp/` infrastructure is never removed
or written beyond RunX's own entries.
