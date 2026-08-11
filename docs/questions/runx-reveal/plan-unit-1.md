---
name: RunX Reveal Plan Unit 1 Question Ledger
purpose: Preserve the selected unattended answer for the first reveal output contract.
description: Records why the first release prints the exact configured command without child-argument or shell-path rewriting.
created: "2026-08-11"
flags:
  - pending-human-review
tags:
  - questions
  - execution
  - cli
keywords:
  - runx reveal
  - child arguments
  - shell quoting
owner: runx-questions-reveal
---

# RunX Reveal Plan Unit 1 Question Ledger

## Question RX-REVEAL-1-Q1

Should the first `runx reveal` release append forwarded child arguments or
rewrite the configured command for a particular shell?

### Evidence

- The user explicitly defined selector parity by numeric index, UID, and full
  command path and asked for the full configured command to copy and run.
- The motivating use case depends on pasting into the user's current Git Bash
  session instead of the Windows `cmd.exe` adapter selected by `shell: auto`.
- RunX cannot reliably infer the shell into which stdout will later be pasted.
- Lossless quoting differs among Bash, `cmd.exe`, and PowerShell, so transformed
  output could change literal child values or recreate the original mismatch.

### Candidate Answers

1. Print the exact stored manifest command only.
2. Append child arguments with one guessed quoting convention.
3. Print the internal shell-adapter invocation and transport details.

### Selected Answer

Candidate 1. The first release accepts exactly one selector and prints the
stored command verbatim plus one newline. Child arguments and shell/path
rewriting remain explicit follow-up scope.

### Rationale

This is the narrowest faithful implementation of the user's enumerated
selection contract, preserves the catalog author's exact command text, keeps
stdout directly copyable in the caller's chosen shell, and never invents a
cross-shell transformation.

### Confidence And Reversibility

- Confidence: high.
- Reversibility: high; a later additive release can introduce an explicit
  target-shell or argument-rendering contract without changing exact-command
  output for the simple form.

### Resulting Action

Implement `runx reveal [--cwd ...] [--config ...] <selector>` as a no-spawn,
no-confirmation, exact-text command.

### Human Review Status

`pending` — execution proceeds under the user-approved unattended-work
instruction and this reversible scope boundary.
