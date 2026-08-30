# Review Handoff — RX-3 — RX Thin Launcher — Sol

## Goal
Mastermind review of RX-3 delivered by Gemini (engineer-3). Verify thin rx launcher correctness, no Cobra duplication, pointer/fallback handling, delegation, and XDocs.

## Context
- Branch: feat/rx-and-user-only-commands (from bfa7e5e + d0ba320)
- Gemini delivered: cmd/rx/main.go, delegate_unix.go, delegate_windows.go, main_test.go, rx.xdocs.md, and updated cmd/cmd.xdocs.md
- Requirements RX-F1..F7, Architecture ADR-001/002/003
- Plan RX-3 acceptance: go build produces both binaries, bare→list, <selector>→run, version/help parity, child-arg forwarding

## Checklist
- main.go translateArgs: bare→list, version/help passthrough (-v/--version/-h/--help/--help-tree/etc.), list with only flags (--cwd/--config/--format) vs run with selector, -- delimiter handling, global flags before selector
- Pointer read via installstate.ReadPointer, payload via launcher.PayloadPath with FallbackPath, error handling exit 5
- delegate_unix.go / delegate_windows.go preserve stdio and exit code (like runx-launcher)
- No manifest parsing or Cobra in rx
- Tests: TestTranslateArgs covers bare, version/help, list options, run with selector, numeric index, child args, options before selector
- XDocs: cmd/rx/rx.xdocs.md valid, cmd/cmd.xdocs.md child list updated, xdocs meta passes
- go vet ./... and go build ./... pass (re-run)
- Verify no userOnly enforcement code in rx (grep)

## Output
- Verdict Ready / Needs fixes with severity
- Do NOT commit/push — report only. Fix blocker if needed with minimal edit and re-run vet/test.

## Worker
- Harness: pi, Provider: openai-codex, Model: gpt-5.6-sol, Thinking: xhigh
