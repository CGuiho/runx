---
name: Protect RunX Branches and Tag Creation
purpose: Define the expected outcome, constraints, and completion signals for the highest-priority RunX repository protection task in TODO.md.
description: Describes the branch protection and tag creation safeguards needed before release activity continues.
created: 2026-07-12
flags:
  - repository-security
tags:
  - todo
  - security
  - release
keywords:
  - runx
  - branch protection
  - protected branches
  - tag creation
  - release tags
owner: runx-todo
---

# Protect RunX Branches and Tag Creation

## Todo Index

- Task: `0. Protect RunX Branches and Tag Creation`
- Status: completed
- Priority: highest
- Index: [TODO.md](../../TODO.md)

## Outcome

RunX has repository rules that protect important branches and restrict tag creation so release history cannot be changed or created without the intended review and authorization controls.

## Scope

### In Scope

- Protect the default branch and any release-maintained branches that RunX uses.
- Restrict creation of release tags and other protected tag patterns.
- Confirm the protections match the RunX release and Mirror workflow before release work continues.

### Out of Scope

- Publishing packages, pushing tags, or changing release versions.
- Bypassing repository protections for local convenience.

## Acceptance Signals

- The RunX repository has branch protection or repository rulesets active for the default branch and any selected release branches.
- Tag creation is restricted for the relevant release tag patterns.
- Required maintainers or administrators can verify that the rules are active before the next release/tagging step.

## Active GitHub Rulesets

- [`Protect main`](https://github.com/CGuiho/runx/rules/18933230) - active branch ruleset `18933230` for `refs/heads/main` and `~DEFAULT_BRANCH`; blocks deletion and non-fast-forward updates, requires pull requests, and enables Copilot review on pushes and draft pull requests.
- [`Protect release tags`](https://github.com/CGuiho/runx/rules/18933243) - active tag ruleset `18933243` for `refs/tags/@guiho*/*@*`; restricts creation, updates, deletion, and non-fast-forward changes and requires linear history.

## Verification Evidence

- GitHub CLI authenticated as repository owner `CGuiho` with the required repository scope.
- Both rulesets were created with `gh api` using the same active conditions and rule parameters used by the Mirror and XDocs repositories.
- Both RunX rulesets were read back from the GitHub API after creation and confirmed active.
- No release tag was created or pushed, and no package was published while configuring the protections.

## Watch-outs

- Do not push tags or publish packages while configuring protections.
- Coordinate with GitHub permissions and GUIHO release ownership before locking down rules that could block maintainers.
- Keep Mirror-managed versioning and tag workflows intact; protections should guard the workflow, not replace it.

## Before Starting

- Confirm the GitHub repository, default branch, release branch names, tag patterns, and required approvers.
- Confirm whether protections should be implemented with classic branch protection rules, repository rulesets, or both.

## After Finishing

- Record the active rules, protected patterns, and any remaining manual verification needed.
- Update this task status only after the rules are confirmed in GitHub.

## References

- [TODO.md](../../TODO.md)
