---
name: RunX GUIHO CLI Convention 0001 Compliance Migration Plan
purpose: Implement the accepted convention-compliance architecture through ordered, independently reviewed and validated work units.
description: Branch-aware plan for stable launchers, immutable payloads, strict artifact ownership, separate configuration, canonical agents, transactional lifecycle operations, and a release-safe breaking cutover.
created: 2026-08-16
flags:
  - proposed
  - breaking-change
  - ready-for-human-approval
  - blocked-pending-human-approval
tags:
  - plan
  - cli
  - compliance
keywords:
  - GUIHO CLI Convention 0001
  - stable launcher
  - artifacts.json
  - current.json
  - transactional lifecycle
  - protocol-v1 cutover
owner: runx-plans
---

# RunX GUIHO CLI Convention 0001 Compliance Migration Plan

## Status

Complete and independently reviewed; ready for human approval.

This plan is not executable until its human approval gates are satisfied. It
does not authorize implementation, branch or worktree creation, commit, push,
installation, PATH/global-agent mutation, issue creation, version application,
tag, GitHub Release, publication, deployment, or production action.

## Outcomes

After every unit, cutover, release, and hardening gate is complete:

- the only supported command entrypoint is a stable launcher under
  `$HOME/.guiho/bin/`;
- payloads and canonical release resources are immutable and versioned;
- `current.json` activates only a committed verified generation and retains a
  verified fallback;
- `artifacts.json` and the installed ledger constrain every lifecycle write and
  deletion;
- installers, reinstall, repair, upgrade, and uninstall are complete,
  synchronous, two-phase, and rollback-safe;
- project `runx.yaml` and global `runx.global.yaml` are distinct strict
  configuration contracts;
- `runx init` performs complete idempotent reconciliation and agent-evolution
  setup;
- canonical skill, prompt, definition, and instruction bytes cannot drift;
- release selection implements full SemVer and channel pagination;
- RunX releases the protocol-v1 resource catalog, initially 25 assets;
- npm is no longer an executable or lifecycle distribution;
- Mirror, RunX, XDocs, CI, release tooling, README, DOCS, skill, and tests
  enforce Convention 0001; and
- every audit finding has exact-head review and validation evidence reachable
  from protected `main`.

## Authority

1. GUIHO CLI Convention 0001.
2. Current GUIHO parent and RunX repository instructions.
3. [Accepted target architecture](../architecture/guiho-convention-0001-cli-compliance-migration.md).
4. [Ready-for-planning architecture review](../reviews/architecture/guiho-convention-0001-cli-compliance-migration-review.md).
5. [Repository compliance audit](../reviews/guiho-convention-0001-cli-compliance-audit.md).
6. This plan after independent review and human approval.
7. Earlier records only where they do not conflict with the preceding sources.

The old direct executable, old global catalog fallback, exact-eleven and
exact-fourteen releases, background replacement, agent `update` names, bare-run
resource mutation, and npm downloader are historical behavior, not constraints.

## Human Approval Gates

### Gate A — Identity And Architecture Confirmation

The human must explicitly confirm or replace all values:

| Decision | Proposed value |
| --- | --- |
| CLI home | `runx` |
| Main skill | `guiho-s-runx` |
| Main setup prompt | `guiho-p-runx` |
| Managed instruction | `guiho-i-runx.md` |
| Repository | `https://github.com/CGuiho/runx` |
| Issue creation | `https://github.com/CGuiho/runx/issues/new` |
| Launcher protocol | `1` |
| Project/global config | `runx.yaml` / `runx.global.yaml` |
| Initial resource contract | one skill, one prompt, one definition, one instruction, two schemas, two examples |

Approval must be durable. Approval of “the plan” alone does not silently confirm
these values.

### Gate B — Plan Approval

The independent plan review at
[`docs/reviews/plans/guiho-convention-0001-cli-compliance-migration-review.md`](../reviews/plans/guiho-convention-0001-cli-compliance-migration-review.md)
must report Ready for approval/execution, and the human must approve this plan
and its breaking/cutover strategy.

### Gate C — Current Approved Base

The audited checkout is `main` at `364fdb8`, 104 commits behind the local
`origin/main` observation, with user-owned `TODO.md` edits and planning files.
No implementation may use that stale, dirty checkout as its base.

Before execution:

1. inventory and preserve every dirty path;
2. refresh remote state read-only;
3. select an exact protected-main commit;
4. recheck audited surfaces that changed upstream;
5. merge the approved planning baseline;
6. create each implementation branch from the exact post-dependency `main`;
7. use an isolated worktree; and
8. never reset, overwrite, absorb, stage, or commit unrelated user work.

### Gate D — Cross-Repository Contract Alignment

Gate D is the result of E00, not an E00 prerequisite. It is satisfied when the
mandatory `guiho-s-0035-cli-engineer-go` alignment is integrated in
`CGuiho/superiority`, or when the human explicitly accepts the governance
conflict for this migration. RunX must not “fix” that other repository from a
RunX worktree.

### Gate E — Cutover And Release Window

Transition cutover C00 may not merge without separate authorization for the
first protocol-v1 Mirror version, protected tag/push, GitHub Release, and remote
smoke in the same controlled window. Publication is not implied by
implementation approval.

## Shared Unit Lifecycle

Every E/U/C/H unit must include this execution envelope:

1. Verify dependencies are reachable from protected `main`.
2. Record exact base SHA, branch, worktree, owned paths, exclusions, and
   question-ledger path before editing.
3. Use one declared non-protected `codex/` branch and isolated worktree.
4. Set its detailed TODO row to `in progress` before implementation and
   `testing` before final checks.
5. Resolve new questions from the plan's accepted assumptions. Use a safe,
   reversible provisional answer only when it cannot change public behavior,
   ownership, security, data, or release scope; record it in the ledger.
   Otherwise reject/requeue the unit to the Plan Writer with the exact
   contradiction. An unattended executor does not pause for interactive input.
6. Modify only owned paths and include their tests/docs/XDocs in the same unit.
7. Use isolated homes, PATHs, caches, fake release servers, and processes.
8. Create the smallest coherent commit or commit sequence, push the declared
   branch plainly, and open the declared PR automatically when the human has
   approved the execution envelope below.
9. Target `main`.
10. `guiho-a-0049-implementation-reviewer` reviews the exact PR head without
    fixing it.
11. `guiho-a-0050-validation-reporter` validates that same unchanged head.
12. Any correction invalidates both gates and requires both again.
13. `guiho-a-0052-pull-request-integrator` merges only after accepted review,
    accepted validation, live CI/mergeability re-observation, and integration
    authority.
14. Verify merge reachability from protected `main`.
15. Remove only the integrated branch and its worktree after reachability.
16. Rebase or re-plan downstream units against the new protected-main state.

Human approval to execute must explicitly cover U00-U11, C00, and H00 branch,
worktree, commit, push, PR, review, validation, integration, and cleanup actions.
Under that envelope, exact bases are reobserved and recorded, not re-approved
interactively per unit. E00 still needs separate cross-repository authorization,
and R00 still needs separate release/publication authorization. A rejected unit
returns to planning; it does not broaden its authority.

Every unit records a scope-correct terminal version/release decision:

- E00 records `mirror: deferred to a separately authorized Superiority release`,
  `github_release: none`, and `production_mutation: none`; RunX R00 cannot version
  or distribute the Superiority skill.
- U00 records `mirror: not applicable; planning-only`, `github_release: none`,
  and `production_mutation: none`.
- U01-U11 and C00 record `mirror: consolidated into separately authorized RunX
  R00 because this unit is not independently releasable`, `github_release:
  none`, and `production_mutation: none`.
- R00 uses its release-operation envelope and records the exact applied version,
  tag, release URL, and every approved production mutation.
- H00 records an exact `mirror version plan patch` result. It may record
  `deferred: no release required` when its diff is limited to the main-branch
  remote installer's post-v1 compatibility removal plus docs, CI, and validation
  hardening, because the accepted protocol-v1 release assets remain unchanged.
  Any binary, embedded resource, or versioned release-asset change instead
  creates a separately authorized follow-up patch-release operation. H00 never
  defers backward to R00.

Each question ledger lives at
`docs/questions/guiho-convention-0001-cli-compliance-migration/<unit>.md` and
must explicitly say “none” when no question arose.

## Global Safety And Quality Invariants

- No destructive target comes from a prefix or filename alone.
- Removal authority is installed ledger + immutable manifest ownership + hard
  allowed-root intersection.
- Symlink, junction, traversal, alternate-separator, and foreign-root escapes
  fail closed.
- Persistent configuration/data never rolls back during reinstall/upgrade.
- No lifecycle validation touches the real home, PATH, skills, repository
  instruction, installation, or live RunX processes.
- Ordinary upgrade never replaces its launcher.
- No process is terminated by filename; path, user, PID, start identity, and
  instance token must agree.
- An uncommitted pointer generation never handles an ordinary invocation.
- Only inactive locked-file cleanup may be deferred; success is never
  `scheduled`.
- Source and generated schemas/examples/resource contracts cannot drift.
- No unit applies Mirror, tags, releases, publishes, or mutates production
  except separately authorized R00.

## Plan Tree

```text
E00  Align the canonical GUIHO Go CLI skill (Superiority prerequisite)
U00  Integrate the approved planning baseline on current RunX main
U01  Establish mandatory RunX, XDocs, and Mirror tooling
U02  Implement separate configuration, policy, schemas, and init
U03  Canonicalize agent resources and commands
U04  Implement release manifest, checksum, SemVer, and selection contracts
U05  Implement ownership state and two-phase lifecycle primitives
U06  Implement the stable launcher and fallback protocol
U07  Build the dormant protocol-v1 release toolchain
U08  Build dormant next-generation installers and legacy migration
U09  Build dormant next-generation uninstall interfaces
U10  Build dormant synchronous whole-release upgrade
U11  Correct help-tree and raw-version behavior
C00  Atomically wire the transition public surface
R00  Publish and verify the first protocol-v1 release (separate authorization)
H00  Remove transition support and prove final compliance
```

Units are sequential. A later plan revision may parallelize only work with
disjoint ownership and no shared lifecycle/API dependency.

## Unit Command Matrix

These are minimum exact commands. A unit may add narrower tests but may not
replace these with prose. Commands run from that unit's isolated worktree with
workspace-local caches where required. `RG` below means this exact repository
gate, in this order:

```text
gofmt -l main.go cmd pkg devops
go mod tidy -diff
go test -count=1 ./...
go vet ./...
go build ./...
mirror config check
runx check --format json
runx list --format json
xdocs meta . --documents --strict
xdocs scan
xdocs tree
xdocs doctor .
git diff --check
```

`RG+EMBED` means `gofmt -l resources_embed.go` followed by every `RG` command.
The named gates are literal expansions, not discretion to substitute a smaller
check. A command that is unavailable or fails is a failed gate, not a skip.

| Unit | Required commands |
| --- | --- |
| E00 | `bun -e "JSON.parse(await Bun.file('skills/guiho-s-0035-cli-engineer-go/evals/evals.json').text())"`; `xdocs meta skills/guiho-s-0035-cli-engineer-go --documents --strict`; `xdocs tree`; `xdocs doctor skills/guiho-s-0035-cli-engineer-go`; `git diff --check` |
| U00 | `xdocs meta docs/architecture --documents --strict`; `xdocs meta docs/plans --documents --strict`; `xdocs meta docs/todo --documents --strict`; `xdocs meta docs/reviews --documents --strict`; `xdocs tree`; `xdocs doctor docs`; `git diff --check` |
| U01 | `mirror config check`; `runx check --format json`; `runx list --format json`; `xdocs meta . --documents --strict`; `xdocs scan`; `xdocs tree`; `xdocs doctor .`; `go test -count=1 ./...`; `git diff --check` |
| U02 | `go test -count=1 ./pkg/config/... ./pkg/manifest/... ./cmd/conventionv1/...`; `go vet ./pkg/config/... ./pkg/manifest/... ./cmd/conventionv1/...`; `go run ./devops/conventionv1/validate-schemas.go`; `RG` |
| U03 | `gofmt -l resources_embed.go`; `go test -count=1 ./pkg/agent/... ./cmd/conventionv1/...`; `go vet ./pkg/agent/... ./cmd/conventionv1/...`; `go run ./devops/conventionv1/verify-resources.go`; `RG+EMBED` |
| U04 | `go test -count=1 ./pkg/release/...`; `go vet ./pkg/release/...`; `go run ./devops/conventionv1/verify-release-fixtures.go`; `RG+EMBED` |
| U05 | `go test -count=1 ./pkg/installstate/... ./pkg/lifecycle/...`; `go vet ./pkg/installstate/... ./pkg/lifecycle/...`; `RG+EMBED` |
| U06 | `go test -count=1 ./pkg/launcher/... ./pkg/installstate/... ./devops/conventionv1/test/launcher/...`; `go vet ./pkg/launcher/...`; `go run ./devops/conventionv1/build-launchers.go --version 0.0.0-test --commit test --build-date 2000-01-01T00:00:00Z`; `go run ./devops/conventionv1/smoke-launcher.go`; `RG+EMBED` |
| U07 | `go run ./devops/conventionv1/build-release.go --version 0.0.0-test --commit test --build-date 2000-01-01T00:00:00Z`; `go run ./devops/conventionv1/verify-release.go`; `go run ./devops/conventionv1/verify-source-contract.go --phase dormant`; `go test -count=1 ./devops/conventionv1/...`; `RG+EMBED` |
| U08 | `go test -count=1 ./pkg/lifecycle/... ./devops/conventionv1/test/installers/...`; `bash -n devops/conventionv1/install.sh`; `powershell.exe -NoProfile -Command "[void][scriptblock]::Create((Get-Content -Raw 'devops/conventionv1/install.ps1'))"`; `RG+EMBED` |
| U09 | `go test -count=1 ./pkg/lifecycle/... ./devops/conventionv1/test/uninstallers/... ./devops/conventionv1/test/windows/...`; `bash -n devops/conventionv1/uninstall.sh`; `powershell.exe -NoProfile -Command "[void][scriptblock]::Create((Get-Content -Raw 'devops/conventionv1/uninstall.ps1'))"`; `RG+EMBED` |
| U10 | `go test -count=1 ./pkg/lifecycle/... ./pkg/release/... ./cmd/conventionv1/... ./devops/conventionv1/test/upgrade/...`; `go vet ./pkg/lifecycle/... ./pkg/release/... ./cmd/conventionv1/...`; `RG+EMBED` |
| U11 | `go test -count=1 ./cmd/conventionv1/... ./pkg/version/...`; `go vet ./cmd/conventionv1/... ./pkg/version/...`; `RG+EMBED` |
| C00 | `bash -n devops/install.sh`; `bash -n devops/uninstall.sh`; `powershell.exe -NoProfile -Command "[void][scriptblock]::Create((Get-Content -Raw 'devops/install.ps1')); [void][scriptblock]::Create((Get-Content -Raw 'devops/uninstall.ps1'))"`; `go test -count=1 ./devops/conventionv1/test/transition/... ./devops/conventionv1/test/installers/... ./devops/conventionv1/test/uninstallers/... ./devops/conventionv1/test/upgrade/...`; `go run ./devops/conventionv1/verify-source-contract.go --phase transition`; `RG+EMBED` |
| H00 | `bash -n devops/install.sh`; `bash -n devops/uninstall.sh`; `powershell.exe -NoProfile -Command "[void][scriptblock]::Create((Get-Content -Raw 'devops/install.ps1')); [void][scriptblock]::Create((Get-Content -Raw 'devops/uninstall.ps1'))"`; `go run ./devops/conventionv1/build-release.go --version 0.0.0-test --commit test --build-date 2000-01-01T00:00:00Z`; `go run ./devops/conventionv1/verify-release.go`; `go run ./devops/conventionv1/verify-source-contract.go --phase final`; `go test -count=1 ./devops/conventionv1/test/launcher/... ./devops/conventionv1/test/installers/... ./devops/conventionv1/test/uninstallers/... ./devops/conventionv1/test/upgrade/... ./devops/conventionv1/test/windows/...`; `mirror version plan patch`; `RG+EMBED`; repeat `RG+EMBED` from a fresh clean clone of the integrated protected-main SHA |

Paths named here are expected deliverables of their owning earlier unit. If a
path is renamed, that is a material plan change and must be revised before its
consumer unit begins.

---

## E00 — Align The Canonical GUIHO Go CLI Skill

**Goal:** Remove the mandatory skill's conflict with Convention 0001.

**Repository/base:** `C:\GUIHO\superiority` from an explicitly approved current
protected-main SHA. Branch `codex/go-cli-convention-0001-contract`; isolated
worktree `C:\GUIHO\worktrees\superiority\go-cli-convention-0001-contract`; PR to
`main`. Question ledger
`docs/questions/guiho-s-0035-cli-convention-0001/E00.md`.

**Dependencies:** Gate A, explicit cross-repository execution authorization,
and a recorded current Superiority `main` SHA. Successful integration satisfies
Gate D. This is the first cross-repository unit.

**Owned:** `skills/guiho-s-0035-cli-engineer-go/**`;
`TODO.md` only for a new conflict-safe dispatch;
`docs/todo/guiho-s-0035-cli-convention-0001.md`;
`docs/plans/guiho-s-0035-cli-convention-0001.md`;
`docs/reviews/plans/guiho-s-0035-cli-convention-0001-review.md`;
`docs/reviews/implementation/guiho-s-0035-cli-convention-0001-review.md`;
`docs/validation/guiho-s-0035-cli-convention-0001.md`; the question ledger and
their direct XDocs descriptors. **Excluded:** every RunX path including its
dirty `TODO.md` and planning files, unrelated Superiority skills/agents, installed
global copies, version application, release/publication.

**Actions:**

1. Replace exact-eleven/direct-binary/background-replacement rules with stable
   launcher, immutable store, complete manifest, synchronous transaction, full
   uninstall, configuration, init, and agent-evolution rules.
2. Replace agent `update` with `upgrade` and remove bare-run mutation mandates.
3. Make resource counts derive from a versioned resource contract.
4. Update templates, examples, eval expectations, and XDocs.
5. Search the complete skill for obsolete mandatory text.

**Impacts:** Governance/docs/evals only; no user data, configuration, auth,
cache, installation, or production.

**Tests/acceptance:** `xdocs meta skills/guiho-s-0035-cli-engineer-go --documents
--strict`; `xdocs tree`; `xdocs doctor
skills/guiho-s-0035-cli-engineer-go`; `bun -e
"JSON.parse(await Bun.file('skills/guiho-s-0035-cli-engineer-go/evals/evals.json').text())"`;
exact obsolete-text search; and `git diff --check`. Accept when the skill no
longer contradicts the convention,
0049/0050 accept one exact head, and the merge is reachable.

**Stop:** unapproved cross-repository work, unrelated policy impact, dirty owned
paths, or unresolved convention authority.

**TODO/delivery/Mirror:** dedicated Superiority task progresses through the
shared lifecycle. Focused PR and cleanup as above. No Mirror apply/release; any
skill distribution is a later separate decision.

---

## U00 — Integrate The Approved Planning Baseline

**Goal:** Put the audit, accepted architecture/review, plan/review, and task
ledger onto a current clean RunX `main` without product implementation.

**Base/branch:** exact refreshed protected-main SHA; branch
`codex/runx-cli-convention-u00-planning-baseline`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u00-planning-baseline`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U00.md`; PR to `main`.

**Dependencies:** E00 integrated or conflict explicitly accepted; Gates A-C.

**Owned:** this migration's audit, architecture/review, plan/review, TODO spec,
one conflict-safe root TODO entry, affected descriptors, narrow supersession
frontmatter. **Excluded:** product/workflow/scripts and unrelated TODO/history.

**Actions:** preserve the user TODO patch byte-for-byte; refresh/re-audit changed
surfaces; transplant only reviewed planning documents; durably record identity
approval; add the concise TODO dispatch; mark only directly contradictory active
records superseded; update descriptors; prove no implementation entered.

**Impacts:** documentation/governance only.

**Tests/acceptance:** XDocs strict metadata/tree/doctor, link/frontmatter check,
`git diff --check`, diff proof of preserved user TODO. Accept when planning truth
is on current `main` and the task is ready to execute.

**Stop:** TODO overlap, upstream material drift, missing approval, review
rejection, or unexplained dirty owned path.

**TODO/delivery/Mirror:** shared lifecycle, 0049/0050 exact-head docs review,
0052 integration. No version/release.

---

## U01 — Establish Mandatory Project Tooling

**Goal:** Make Mirror, RunX, and XDocs validate the repository's supported
commands, version surfaces, and complete owned tree.

**Base/branch:** post-U00 `main`; branch
`codex/runx-cli-convention-u01-project-tooling`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u01-project-tooling`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U01.md`; PR to `main`.

**Dependencies:** U00.

**Owned:** root `runx.yaml`, `xdocs.yaml`, `mirror.yaml`, retirement of
`xdocs.config.toml` after parity, `XDOCS.md`, missing/root descriptors, tooling
tests. **Excluded:** lifecycle/product implementations, manual version values,
workflows except a read-only validation hook if required.

**Actions:** catalog all current supported commands with stable UIDs; introduce
root XDocs YAML; cover all owned directories and repair parent/child links;
configure Mirror for every current version-bearing file; add nonmutating tooling
checks; require later units to extend catalog/descriptors in their PR.

**Impacts:** repository config/docs only; no secrets, user data, live cache, or
global state.

**Tests/acceptance:** `runx check/list --format json`, XDocs strict
meta/scan/tree/doctor, `mirror config check`, Go tests, `git diff --check`.
Accept with no XDocs ownership gap and no version authority outside Mirror.

**Stop:** metadata loss, unsupported Mirror target, mutating catalog check, or
unresolved command identity.

**TODO/delivery/Mirror:** row U01; shared exact-head delivery. Mirror config
check only.

For U02-U11, narrow `runx.yaml` ownership means only adding or updating that
unit's own stable-UID, nonmutating development/validation commands and their
descriptions. A unit may not change another unit's UID/body or public product
selector through this allowance. Each PR must prove `runx check` and `runx list`
parity; its direct XDocs descriptors remain owned with the files they describe.
C00 may retire or rename dormant-only catalog entries only while atomically
moving their underlying helpers to canonical paths.

---

## U02 — Configuration, Evolution Policy, Schemas, And Init

**Goal:** Implement distinct strict global/project contracts, inheritance,
complete policy, version-pinned schemas/examples, and the dormant reconciliation
engine that U03 will connect to canonical agent resources and public `init`.

**Base/branch:** post-U01 `main`; branch
`codex/runx-cli-convention-u02-config-init`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u02-config-init`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U02.md`; PR to `main`.

**Dependencies:** U01 and confirmed identities.

**Owned:** `pkg/config/**`, needed `pkg/manifest/**` revisions, private
`cmd/conventionv1/**` init/config adapters and tests, an internal
dependency-injected init orchestrator, `schemas/**`, `examples/**`, focused
`devops/conventionv1/validate-schemas.go`, tests/fixtures/descriptors, and narrow
`runx.yaml` additions for this unit's stable validation commands.
**Excluded:** the public init-command switch,
canonical agent resources/commands, launcher/installstate, lifecycle scripts,
upgrade/uninstall, release workflow, real config/home.

**Actions:** define strict types and overlays; model catalogs as project-only
without changing current public resolution; implement four agent-evolution
leaves and exact enums; validate schemas/examples; create pinned schema comments;
build the private init orchestrator as validate-plan-apply-revalidate; preserve
existing catalog/values; plan bounded `AGENTS.md` operations; fail test-only
noninteractive missing answers without writes; model every remaining RunX
domain question; produce Created/Upgraded/Verified/Unchanged summary groups with
absolute paths. Do not register the orchestrator, remove public fallback, or
change bare-run/background behavior before C00.

**Impacts:** config contract changes; no database; caches remain disposable;
policy governs agent authority but grants no new technical permission.

**Tests/acceptance:** unknown fields/enums, full overlay matrix, no-parent/no-
fallback, schemas/examples, interactive choices, noninteractive zero-write,
catalog preservation, orchestrator idempotency, rollback injection, bare-run no
mutation. Accept when both configs are independent/strict and the orchestrator
is ready for U03 resource injection without weakening the current public path.

**Stop:** data/catalog loss, Go/schema drift, ambiguous policy, or real-home
dependency.

**TODO/delivery/Mirror:** row U02; shared lifecycle. No version application.

---

## U03 — Canonical Agent Resources And Commands

**Goal:** Build the private byte-identical protocol-v1 source for skill,
definition, prompt, and instruction plus convention command names/raw show,
while retaining the active legacy compatibility graph until C00.

**Base/branch:** post-U02 `main`; branch
`codex/runx-cli-convention-u03-agent-resources`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u03-agent-resources`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U03.md`; PR to `main`.

**Dependencies:** U02, Gate A.

**Owned:** `pkg/agent/**`, unregistered convention-v1 command constructors and
test adapters, `devops/conventionv1/verify-resources.go`, convention-v1 canonical
`skills/guiho-s-runx/**`, `prompts/guiho-p-runx.md`,
`instructions/guiho-i-runx.md`, additive private `resources_embed.go`,
tests/descriptors, and narrow `runx.yaml` additions for this unit's stable
validation commands. **Excluded:** active legacy `embed/**`, `cmd/root.go`,
`cmd/agent.go`, `pkg/maintenance/**`, live projections/instruction, release
builder, lifecycle state/scripts, and `.github/workflows/ci.yml`.

**Actions:** build the top-level protocol-v1 resources and private embedding;
add exact
`CLI Evolution and Feedback` section, URLs, policies, issue guidance, install,
verify, init, upgrade steps; create main setup prompt; replace hard-coded
instruction only inside the private protocol-v1 resource graph; implement
`upgrade`-named private handlers, raw show, bundled
listing, and bounded `AGENTS.md` behavior without registering them in the public
Cobra tree; inject canonical resources into the U02 orchestrator; expose
projection hashes to the later ownership model. Retain the active legacy embed
package, imports, maintenance implementation, bytes, and CI checks unchanged.
C00 alone switches active imports/maintenance, retires legacy `embed/**`,
replaces the public agent/init wiring, and removes the old aliases.

**Impacts:** agent docs/commands/projections; no external issue or global
mutation during tests.

**Tests/acceptance:** protocol-v1 source/private-embed/package byte identity,
skill validation, raw
goldens, `upgrade`-only help, block/line-ending/idempotency/malformed-marker
tests on the unregistered tree, policy guidance, complete interactive/
noninteractive orchestrator tests, XDocs/search, and a probe proving current
public agent/init behavior is unchanged. Accept when drift is structurally
impossible inside the protocol-v1 graph, the legacy active graph remains
byte-for-byte unchanged, and C00 can atomically retire that compatibility graph
and activate a complete idempotent init without new logic.

**Stop:** unconfirmed identity, second source copy, or unprovable byte
preservation.

**TODO/delivery/Mirror:** row U03; shared lifecycle. No release.

---

## U04 — Release Integrity, SemVer, Channels, And Selection

**Goal:** Strict serializable `artifacts.json`/checksum contracts and complete
exact/channel/stable selection across every release page.

**Base/branch:** post-U03 `main`; branch
`codex/runx-cli-convention-u04-release-contract`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u04-release-contract`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U04.md`; PR to `main`.

**Dependencies:** U03.

**Owned:** `pkg/release/**`, retirement/migration of pure selection code under
`pkg/update`, `devops/conventionv1/verify-release-fixtures.go`, release
contract fixtures/descriptors, and narrow `runx.yaml` additions for this unit's
stable validation commands. **Excluded:** filesystem
mutation, public builder/workflow/script changes, live release writes.

**Actions:** strict manifest/integrity union; checksum grammar/order; archive
member validation; installed path/projection/ownership semantics; strict SemVer
and numeric prerelease ordering; full pagination/draft rejection; mutually
exclusive selectors; validate protocol-v1 complete resource contract and target
matrix before selecting; freeze selection; candidate target/version/self-test
requirements.

**Impacts:** JSON/release contracts and disposable catalog cache only; no auth or
user data.

**Tests/acceptance:** SemVer table, channels, pages, malformed/incomplete
release, self/checksum/duplicates, traversal/ownership, all target maps. Accept
when no partial/malformed release can be selected.

**Stop:** unowned artifact, ambiguous selector, circular serialization, or
consumer-specific selection behavior.

**TODO/delivery/Mirror:** row U04; shared lifecycle. No publication.

---

## U05 — Install State And Two-Phase Lifecycle Primitives

**Goal:** Safe dependency-injected paths, pointer/ledger, locks, instances,
journals, ownership, and capability-bound prepare/activate/finalize/abort.

**Base/branch:** post-U04 `main`; branch
`codex/runx-cli-convention-u05-install-state`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u05-install-state`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U05.md`; PR to `main`.

**Dependencies:** U04.

**Owned:** `pkg/installstate/**`, foundational `pkg/lifecycle/**`, platform
atomic/process helpers, fixtures/descriptors, and narrow `runx.yaml` additions
for this unit's stable validation commands. **Excluded:** public Cobra wiring,
canonical scripts, launcher entrypoint, network, real processes/home/PATH.

**Actions:** native home/roots; strict pointer and installed ledger; atomic
durable replacements; symlink/junction-safe allowed roots; token/PID/start/path
locks; instance registry; journal states and pointer generations; persistence
classes; removal intersection; two-phase capability operations; adapter/Go
journal joint commit and abandoned-adapter abort; inject every external seam.

**Impacts:** state formats only in isolated tests; no real auth/config/data.

**Tests/acceptance:** strict decode, atomic failures, traversal/junction/symlink,
token/stale/PID reuse/path mismatch, child exclusion, all journal/pointer
transitions, PATH snapshot recovery fixtures, persistence preservation. Accept
when every destructive target and process is proven.

**Stop:** prefix ownership, unverifiable process identity, platform contract
divergence, uncommitted dispatch, or incomplete abort.

**TODO/delivery/Mirror:** row U05; shared lifecycle; critical 0049 ownership
review. No release.

---

## U06 — Stable Launcher And Fallback

**Goal:** Minimal separate launcher with exact delegation, committed-generation
dispatch, safe verified fallback, and capability bootstrap host.

**Base/branch:** post-U05 `main`; branch
`codex/runx-cli-convention-u06-stable-launcher`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u06-stable-launcher`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U06.md`; PR to `main`.

**Dependencies:** U05.

**Owned:** `cmd/runx-launcher/**`, `pkg/launcher/**`, narrow installstate
extensions, `devops/conventionv1/build-launchers.go`,
`devops/conventionv1/smoke-launcher.go`, build/smoke fixtures/descriptors, and
narrow `runx.yaml` additions for this unit's stable validation commands.
**Excluded:** public installer,
Cobra domain, live bin/home, publication.

**Actions:** strict pointer/protocol; dispatch committed active only; args/env/
cwd/streams/signals/exit preservation; instances; fallback only on start
failure; lock-aware bounded wait; never retain failed pointer; local journal
repair; capability-token transaction dispatch; eight pure-Go launcher builds.

**Impacts:** process/state only in isolated layouts.

**Tests/acceptance:** exact delegation, no fallback on domain error, every
pointer/lock/journal state, fallback repair/null previous, all builds, native
smokes, target/protocol/self-test. Accept when launcher is domain-independent
and protocol 1 is frozen.

**Stop:** imperfect delegation, unverified fallback, protocol question, or real
installation dependency.

**TODO/delivery/Mirror:** row U06; shared lifecycle. No published launcher yet.

---

## U07 — Dormant Protocol-v1 Release Toolchain

**Goal:** Build/verify the full new release without changing the currently
published legacy contract before cutover.

**Base/branch:** post-U06 `main`; branch
`codex/runx-cli-convention-u07-release-toolchain`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u07-release-toolchain`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U07.md`; PR to `main`.

**Dependencies:** U06.

**Owned:** additive `devops/conventionv1/**` builder/verifier/contract fixtures,
including `verify-source-contract.go` and `test/launcher/**`; release-contract
source, tests/descriptors, narrow ignore rules, and narrow `runx.yaml` additions
for this unit's stable validation commands. **Excluded:**
canonical legacy builder/verifier, publish workflow, canonical asset names,
generated output, manual versions.

**Actions:** build 8 renamed `runx-payload-*` and 8 launchers; package canonical
resources; deterministic manifest/checksums; derive initial 25 assets from the
resource contract; reject declared-set errors; native smoke; exact release-note
extraction reuse. Do not make main-branch release upload use it yet.

**Impacts:** build cache/output only; generated bundles untracked.

**Tests/acceptance:** 16 cross-builds, exact derived set, reproducibility,
archive members, native smoke, protocol/target/version, Mirror config. Accept
when C00 can wire the toolchain without algorithm changes.

**Stop:** nondeterminism, missing ownership, current release behavior change, or
generated output in Git.

**TODO/delivery/Mirror:** row U07; shared lifecycle. No tag/upload.

---

## U08 — Dormant Next-Generation Installers

**Goal:** Complete PowerShell/POSIX install, reinstall, repair, two-phase PATH,
and shell-aware legacy migration without replacing canonical remote scripts.

**Base/branch:** post-U07 `main`; branch
`codex/runx-cli-convention-u08-next-installers`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u08-next-installers`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U08.md`; PR to `main`.

**Dependencies:** U07.

**Owned:** `devops/conventionv1/install.sh`,
`devops/conventionv1/install.ps1`, adapter fixtures, install/repair lifecycle
parts, direct tests/descriptors, and narrow `runx.yaml` additions for this unit's
stable validation commands. **Excluded:** canonical `devops/install.*`,
publish workflow, real home/PATH, uninstall/upgrade.

**Actions:** exact/channel/stable selection; complete target/common download;
checksums/target/version/self-test; protocol capability; prepare/activate; exact
PATH snapshot and shell resolution; proof-only legacy migration; finalize/abort;
active-version repair; complete rollback; persistent preservation; init/version
follow-up output. Accept only full `--version`/`--channel` and
`-Version`/`-Channel`, with no short `-v` and no arbitrary install-directory
override. Include transition dual-shape mode as an isolated fixture, never mix
shapes.

**Impacts:** isolated filesystem/PATH/projections; config/data preserved.

**Tests/acceptance:** real script black boxes on Windows/POSIX runners, same fake
release/final tree, all fault boundaries, exact PATH restoration, legacy owned/
foreign/missing checksums/custom shadow, staging descendant, re-run repair.
Accept when both adapters are logically identical and ready for path wiring.

**Stop:** incomplete rollback, shell ambiguity without fail-closed result,
platform drift, direct payload on PATH, or real-home requirement.

**TODO/delivery/Mirror:** row U08; shared lifecycle and critical 0049 deletion/
rollback review. No canonical installer change/release.

---

## U09 — Dormant Next-Generation Uninstall

**Goal:** Equivalent safe uninstall plans for Cobra engine, PowerShell, and
POSIX, including Windows delete-on-close semantics.

**Base/branch:** post-U08 `main`; branch
`codex/runx-cli-convention-u09-next-uninstall`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u09-next-uninstall`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U09.md`; PR to `main`.

**Dependencies:** U08.

**Owned:** uninstall lifecycle packages, nonpublic Cobra adapter exercised only
by tests, `devops/conventionv1/uninstall.*`, fixtures/descriptors, and narrow
`runx.yaml` additions for this unit's stable validation commands. **Excluded:**
canonical Cobra handler/scripts, shared GUIHO roots/PATH, real install.

**Actions:** exact REMOVE/PRESERVE plan; `--yes`/dry-run/preserve flags;
noninteractive fail; ledger/manifest/root intersection; managed block only;
external Windows path; Cobra NTFS rename/quarantine/delete-on-close handoff;
unsupported capability zero-write fallback; crash journal recovery; verify
removal.

**Impacts:** default data/config destruction only inside isolated homes;
preservation behavior explicit.

**Tests/acceptance:** default/each/combined preserve, dry-run, terminal cases,
foreign/shared preservation, traversal, real launcher/payload Windows black box,
all crash boundaries, unsupported filesystem. Accept when three interfaces
produce equivalent results and no detached success exists.

**Stop:** ambiguous/unowned target, shared root deletion, confirmation bypass,
unverified Windows removal, or helper-based eventual success.

**TODO/delivery/Mirror:** row U09; shared lifecycle and critical deletion review.
No canonical exposure/release.

---

## U10 — Dormant Synchronous Whole-Release Upgrade

**Goal:** Complete whole-release upgrade implementation and recovery output
without switching legacy public installations before cutover.

**Base/branch:** post-U09 `main`; branch
`codex/runx-cli-convention-u10-next-upgrade`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u10-next-upgrade`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U10.md`; PR to `main`.

**Dependencies:** U09.

**Owned:** new lifecycle upgrade engine, private/test-only command adapter,
retirement-ready updater interfaces, `devops/conventionv1/test/upgrade/**`,
fixtures/descriptors, and narrow `runx.yaml` additions for this unit's stable
validation commands. **Excluded:** active public `cmd/upgrade.go` switch,
launcher replacement, real catalog/processes.

**Actions:** first recovery block before network; lock/selection; full stage/
verify; immutable install; exact other-instance termination; projection snapshot;
two-phase activation; launcher/self-test/init/postverify; ledger commit/rollback;
previous retention; final pinned recovery block for every result; remove
scheduled/direct-replacement assumptions from new engine.

**Impacts:** isolated state/process/config; cache advisory/disposable.

**Tests/acceptance:** exact/channel/stable/up-to-date/dry/failure/rollback,
first/final recovery, interruption at each state, PID/path/child checks,
projection rollback, no launcher replacement, exact version. Accept when C00 can
switch public wiring without algorithm changes.

**Stop:** async activation, partial release, filename kill, missing recovery,
uncommitted dispatch, or real process termination.

**TODO/delivery/Mirror:** row U10; shared exact-head lifecycle. No live upgrade.

---

## U11 — Help-Tree And Raw Version

**Goal:** Build and test the convention help/version implementation without
changing public behavior before cutover.

**Base/branch:** post-U10 `main`; branch
`codex/runx-cli-convention-u11-help-version`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-u11-help-version`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/U11.md`; PR to `main`.

**Dependencies:** U10.

**Owned:** additive unregistered help-tree/version helpers and private test
adapters under `cmd/conventionv1/**` and `pkg/version/**`, tests/descriptors, and
narrow `runx.yaml` additions for this unit's stable validation commands.
**Excluded:** active `cmd/help.go`, `cmd/root.go`, `main.go`, npm, scripts,
release wiring, lifecycle packages, version bump.

**Actions:** implement and privately test depth `max` or integer >1, default
max, rejection of 1, false-by-default global-flags behavior, single-root versus
descendant flags, source `0.0.0-dev`, and raw release SemVer. Do not register the
flags or change active version output before C00.

**Impacts:** public output only.

**Tests/acceptance:** private adapter/golden max/2/equals/space/invalid/1,
global false/true, every scope, source/release raw versions, no pollution, plus a
probe proving current public output remains unchanged. Accept when C00 can switch
the one Cobra tree without algorithm changes.

**Stop:** global flag leakage, nonraw version, or duplicated command tree.

**TODO/delivery/Mirror:** row U11; shared lifecycle. No Mirror/version apply.

---

## C00 — Transition Public Cutover

**Goal:** Wire all prevalidated dormant components into one coherent public
transition surface that remains installable before and after the first
protocol-v1 release.

**Base/branch:** post-U11 `main`; branch
`codex/runx-cli-convention-c00-transition-cutover`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-c00-transition-cutover`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/C00.md`; PR to `main`.

**Dependencies:** E00 and U00-U11 integrated; Gate E; approved release window.

**Owned:** active entry and Cobra surfaces `main.go`, `cmd/root.go`,
`cmd/catalog.go`, `cmd/agent.go`, `cmd/help.go`, `cmd/upgrade.go`, new
`cmd/uninstall.go`, `cmd/root_test.go`, and `cmd/cmd.xdocs.md`; active resolution
and maintenance surfaces `pkg/manifest/**`, `pkg/update/**`, `pkg/updater/**`,
`pkg/maintenance/**`, and their direct descriptors/tests; narrow public wiring
in `pkg/config/**`, `pkg/agent/**`, `pkg/lifecycle/**`, `pkg/installstate/**`,
`pkg/release/**`, `pkg/launcher/**`, and `pkg/version/**`; canonical
`resources_embed.go`; active legacy `embed/**` retirement; `devops/install.sh`,
`devops/install.ps1`, new
`devops/uninstall.sh`, new `devops/uninstall.ps1`, `devops/build-binaries.go`,
`devops/verify-release-assets.go`, their direct tests, and
`devops/devops.xdocs.md`; dormant-to-canonical moves/removals under
`cmd/conventionv1/**` and `devops/conventionv1/**`;
`.github/workflows/ci.yml`, `.github/workflows/publish.yml`, and
`.github/workflows/workflows.xdocs.md`; `package.json`,
`scripts/runx-bin.mjs`, `scripts/runx-bin.spec.ts`, and
`scripts/scripts.xdocs.md`; `README.md`, `DOCS.md`, `runx.yaml`, `xdocs.yaml`,
`mirror.yaml`; `skills/guiho-s-runx/**`, `skills/skills.xdocs.md`,
`prompts/guiho-p-runx.md`, `prompts/prompts.xdocs.md`,
`instructions/guiho-i-runx.md`, new `instructions/instructions.xdocs.md`, and
the direct root/package XDocs descriptors affected by those exact paths.
**Excluded:** changes to algorithms already accepted in U02-U11 except public
wiring and canonical path moves; unrelated domain packages, workflows,
documentation, TODO rows, generated release outputs, or release/publication
actions before merge.

**Actions:**

1. Wire only already reviewed engines and fixtures.
2. Make the transition installer execute either verified legacy shape before
   protocol v1 exists or protocol v1 once compatible for the selector; never
   mix or fall back from available v1 to legacy.
3. Switch builder/verifier/workflow to renamed payloads, launchers, full
   resources, manifest, and checksums.
4. Switch public upgrade/uninstall and required flags.
5. Switch public project/global config resolution, complete init orchestration,
   agent `upgrade` commands/raw show, bare-run mutation removal, disabled-policy
   background behavior, help-tree flags, and raw version.
6. Switch active resource imports and maintenance calls to the accepted
   protocol-v1 graph, retire legacy `embed/**`, and update CI's old embed path,
   asset-name, help-probe, and npm-shim assumptions in the same PR.
7. Retire npm executable/downloader/cache ownership and npm publication:
   remove `package.json.bin`, delete `scripts/runx-bin.mjs`, remove the workflow
   publish step, and declare legacy npm cache removal in the ownership contract.
8. Remove installer short `-v` and every arbitrary install-directory override;
   canonical scripts accept only full version/channel names and canonical roots.
9. Make init report Created/Upgraded/Verified/Unchanged groups with absolute
   paths and ask all remaining RunX domain questions.
10. Update recovery text and transition documentation.
11. Verify legacy installer path still works in isolated home before R00.
12. Keep explicit transition noncompliance documented until H00.

**Impacts:** breaking public lifecycle/release surface; user data still isolated
in tests; release credentials remain unused until R00.

**Tests/acceptance:** entire U01-U11 suite plus legacy-before-v1 transition,
protocol-v1 fake release, no shape mixing, canonical script URLs, workflow
dry/local validation, exact-head 0049/0050. Accept only when merge and R00 are
authorized in one window.

**Stop:** no release authorization, any dormant-to-public algorithm change,
legacy install outage before v1, workflow could publish incomplete assets, or
any exact-head mismatch.

**TODO/delivery/Mirror:** row C00; shared PR lifecycle. Merge is authorized only
as part of the separately approved R00 window. Mirror plan may be inspected;
apply/tag/push/release remain R00.

---

## R00 — First Protocol-v1 Release

R00 is an external-state release operation, not an implementation PR.

**Coordinator/integrator/validators:** `guiho-a-0001-swe` coordinates and binds
the approved envelope; `guiho-a-0052-pull-request-integrator` owns the release
integration itself—Mirror apply, release commit, push, protected tag, GitHub
Release/workflow observation, reachability, and cleanup—because these are
post-merge integration operations. `guiho-a-0050-validation-reporter` validates
the exact source SHA before apply and the unchanged release commit/tag plus
remote result afterward. `guiho-a-0048-plan-executor` may prepare a read-only
preflight record but does not apply Mirror, push, tag, or publish R00.

**Source binding:** record `C00_MERGE_SHA` after proving it is reachable from
protected `main`. Use a fresh isolated clean clone checked out to that exact
`main` state. The release operation has no feature branch/worktree because
Mirror, when separately authorized and operated by 0052, owns the release
commit, protected tag, and push. After apply, record `RELEASE_COMMIT_SHA` and
require the protected tag to point exactly to it. The tree difference from
`C00_MERGE_SHA` must contain only the reviewed Mirror-managed
version/changelog/release mutations.

R00's durable evidence uses a separate post-release docs-only delivery path:
branch `codex/runx-cli-convention-r00-release-evidence`, isolated worktree
`C:\GUIHO\worktrees\runx\cli-convention-r00-release-evidence`, based on the
refreshed protected `main` containing `RELEASE_COMMIT_SHA`, and a PR to `main`.
It may contain only the ledgers/evidence paths below, their direct XDocs
descriptors, and the R00 row in the detailed migration TODO. It must not modify
Mirror-managed version/changelog fields, product code, workflows, or assets.

**Ledgers/evidence:**

- question ledger:
  `docs/questions/guiho-convention-0001-cli-compliance-migration/R00.md`;
- operation record:
  `docs/releases/guiho-convention-0001-first-protocol-v1.md`;
- validation:
  `docs/validation/guiho-convention-0001-first-protocol-v1.md`; and
- remote artifact/smoke evidence attached to the operation record before H00.

Each record names `C00_MERGE_SHA`, `RELEASE_COMMIT_SHA`, protected tag, release
URL, workflow run URLs, exact asset URLs/checksums, runner/platform, command,
exit status, and timestamp. After remote validation, the coordinator creates the
evidence branch from protected `main`; 0049 reviews its exact head, 0050 confirms
that its assertions match the already validated release state, and 0052 merges
it after live CI/mergeability re-observation. H00 remains blocked until this
evidence PR is merged and reachable.

**Preconditions:** C00 merged/reachable; 0049/0050 accepted its exact head; CI
green; transition legacy smoke passes; exact release-contract build verifies;
clean clone; exact target plan reviewed; explicit authority names version
mutations, release commit, push, protected tag, GitHub Release assets, workflow
`production` environment if triggered, and remote installer/upgrade smoke.

**Required pre-apply commands:**

```text
gofmt -l main.go resources_embed.go cmd pkg devops
go mod tidy -diff
go test -count=1 ./...
go vet ./...
go build ./...
mirror config check
runx check --format json
runx list --format json
xdocs meta . --documents --strict
xdocs scan
xdocs tree
xdocs doctor .
go run ./devops/build-binaries.go --version <approved-target> --commit <C00_MERGE_SHA> --build-date <approved-RFC3339>
go run ./devops/verify-release-assets.go
mirror version plan <approved-target>
```

Before apply, re-read `mirror.yaml`, `.github/workflows/publish.yml`, every
workflow triggered by the planned tag, required environment/approval rules, and
live branch/tag protection. Record whether any workflow deploys or promotes
production. Any unapproved new trigger stops R00.

**Apply/publish commands:** 0052 alone runs these only after the human approves
the exact Mirror plan and release envelope:

```text
mirror version apply <approved-target> --yes
git rev-parse HEAD
git rev-parse <protected-tag>^{}
gh release view <protected-tag> --json url,tagName,isDraft,isPrerelease,publishedAt,assets
```

Do not hand-edit versions, create/repoint tags, rerun `npm version`, or publish an
npm executable/package. The protected workflow builds from the recorded release
commit. Verify every remote asset name, size, SHA-256, embedded commit/target,
archive member, manifest record, and release-note version section.

**Remote smoke:** from isolated homes on the available Windows and POSIX native
runners, execute the canonical main-branch exact-version installer, default
stable installer after catalog visibility, stable-channel installer, launcher
`--version`, `init` with scripted answers, `upgrade check`, dry-run upgrade,
recovery first/final blocks, reinstall repair, rollback fault fixture, and both
uninstaller dry-runs. Never target the real user installation.

**Failure handoff:** before protocol-v1 catalog visibility, leave the transition
installer serving the verified legacy shape and record publication unknown/
failed. After any immutable v1 asset is visible, never mutate, replace, retag, or
reuse it. Halt acceptance, record exact remote state/URLs, keep H00 blocked, and
create a new reviewed corrective-release unit. Do not claim latest success.

**Acceptance/cleanup:** source and release SHAs/tag are bound; workflow and
production triggers match approval; remote assets/checksums and all required
smokes pass; direct release URL is recorded; 0050 marks the exact remote result
Ready; 0052 verifies `main` and tag reachability and removes only the isolated
release clone/cache. The evidence-only PR is then reviewed, validated, merged,
and proven reachable through the lifecycle above; only that PR sets the R00 TODO
row to integrated/accepted and unblocks H00. Its worktree/branch are removed by
0052 after protected-main reachability.

---

## H00 — Remove Transition Support And Prove Final Compliance

**Goal:** Remove every legacy release-selection branch, finalize public docs/CI,
and close the audit against one protected-main head.

**Base/branch:** post-R00 protected `main`; branch
`codex/runx-cli-convention-h00-hardening`; worktree
`C:\GUIHO\worktrees\runx\cli-convention-h00-hardening`; ledger
`docs/questions/guiho-convention-0001-cli-compliance-migration/H00.md`; PR to `main`.

**Dependencies:** accepted R00.

**Owned:** canonical installer legacy branch removal; final README/DOCS; CI;
release workflow hardening; `devops/conventionv1/verify-source-contract.go` final
assertions and the exact `devops/conventionv1/test/**` suites named in the command
matrix; supersession metadata; root catalogs/descriptors; audit closure,
implementation review, validation report, migration TODO.
**Excluded:** new product algorithms; version apply/tag/release; unrelated TODO.

**Actions:** reject all legacy/incomplete releases; final README installation
ends in version command; final operational `Uninstall` section with remote
scripts/destructive default/dry/preserve examples; remove obsolete docs/tests;
catalog every final command; strict complete XDocs; CI for all lifecycle,
resource, schema, build, help, Mirror/RunX/XDocs gates; repository searches;
assert `package.json.bin` is absent, `scripts/runx-bin.mjs` is deleted, no npm
publication step remains, and the legacy npm cache is removable only through
declared ownership; clause-by-clause re-audit; clean-clone integrated validation
after merge.

**Impacts:** removes temporary compatibility; docs/CI truth final.

**Tests/acceptance:** run every exact H00 command in the Unit Command Matrix.
The release builder must produce all 16 target-specific launcher/payload builds;
the verifier must derive and validate the initial 25-asset set; the named Go
black-box suites cover native launcher smoke, schemas/examples, installer/
reinstall/repair/uninstall/upgrade faults, Windows behavior, help/version probes,
and the source-contract verifier covers npm absence and obsolete-contract
searches. Require 0049/0050 on the same head and the stated post-merge clean-main
`RG+EMBED` rerun. A missing platform runner or helper is a failure, not an
implicit skip.

Accept only when every audit row has evidence and no active transition or
superseded assertion remains. Any residual creates a focused repair unit; the
task remains testing.

**Stop:** failed/skipped critical gate, incomplete platform coverage without
human risk acceptance, publication mutation, residual legacy selection, or
review/validation SHA mismatch.

**TODO/delivery/Mirror:** row H00, then overall task may complete only after
protected-main reachability and integrated validation. Record the exact
`mirror version plan patch` result. A diff limited to the main-branch remote
installer's post-v1 compatibility removal plus docs, CI, and validation may
record that no new release is required because the accepted protocol-v1 release
assets remain unchanged. Any binary, embedded resource, or versioned release
asset change requires a new separately authorized patch-release operation.
Never defer H00 backward to R00.

## Traceability Matrix

| Audit finding | Closing units |
| --- | --- |
| C-01 stable launcher/immutable payload absent | U04-U08, C00, H00 |
| C-02 upgrade not synchronous whole-release | U04-U06, U10, C00, H00 |
| C-03 installer rollback incomplete | U05, U08, C00, H00 |
| C-04 uninstall unsafe/incomplete | U05, U09, C00, H00 |
| H-01 root tooling absent | U01, H00 |
| H-02 Mirror not sole version authority | E00, U01, C00, R00, H00 |
| H-03 XDocs gaps | U00-U11 descriptors, C00, H00 |
| H-04 channels absent | U04, U08, U10, C00 |
| H-05 wrong staging/destination/ownership | U05, U08, C00 |
| H-06 obsolete release set | E00, U03-U07, C00, R00 |
| H-07 reinstall semantics absent | U05, U08, C00 |
| H-08 init is catalog-only | U02, U03, C00, H00 |
| H-09 config/inheritance/schema/examples absent | U02, U07, C00 |
| H-10 evolution policy absent | U02, U03, U10 |
| H-11 agent names/show wrong | E00, U03, C00 |
| H-12 lifecycle artifacts/guidance incomplete | U03, U07, C00, H00 |
| H-13 recovery output wrong | U10, C00, H00 |
| H-14 release/candidate validation weak | U04-U08, U10, C00 |
| H-15 npm lifecycle bypass | C00, H00 |
| H-16 embedded/released drift | U03, U07, C00 |
| M-01 help-tree contract | U11, H00 |
| M-02 development version | U11, U07, H00 |
| M-03 README lifecycle | C00, H00 |
| M-04 tests/CI superseded | every unit tests; H00 final gate |

## TODO State Model

The detailed task specification stores one row per unit:

```text
pending -> in progress -> testing -> implementation review
        -> validation -> integrated
```

E00 is tracked in Superiority. R00 is `awaiting authorization` until explicitly
approved. The root RunX TODO is an index and must not duplicate the unit detail.
It remains `todo` until execution starts, `testing` throughout C00/R00/H00, and
`completed` only after H00 integrated-main validation.

## First Executable Unit

First cross-repository unit: E00, after cross-repository authorization.

First RunX unit: U00, after Gates A-D and selection of a current protected-main
base. No product unit is executable from the current stale, dirty checkout.

## Completion Definition

Compliance is complete only when E00 or its accepted governance exception is
recorded, U00-U11 and C00 are integrated, R00 is accepted, H00 is integrated,
all audit findings have direct protected-main evidence, and no temporary legacy
selection or superseded active contract remains.

## References

- [Detailed task specification](../todo/guiho-convention-0001-cli-compliance-migration.md)
- [Compliance architecture](../architecture/guiho-convention-0001-cli-compliance-migration.md)
- [Architecture review](../reviews/architecture/guiho-convention-0001-cli-compliance-migration-review.md)
- [Compliance audit](../reviews/guiho-convention-0001-cli-compliance-audit.md)
- [Plan review](../reviews/plans/guiho-convention-0001-cli-compliance-migration-review.md)
- [RunX TODO](../../TODO.md)
