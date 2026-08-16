---
name: RunX GUIHO CLI Convention Compliance Audit
purpose: Determine whether RunX obeys GUIHO CLI Convention 0001 and record every confirmed gap.
description: Evidence-backed repository, runtime, installer, configuration, agent, release, and documentation compliance review.
created: 2026-08-16
owner: runx-reviews
flags:
  - action-required
tags:
  - review
  - cli
  - compliance
  - go
keywords:
  - GUIHO CLI Convention
  - guiho-convention-0001-cli
  - RunX
  - stable launcher
  - immutable payload
  - artifacts manifest
  - agent evolution
---

# RunX Compliance Audit Against GUIHO CLI Convention 0001

## Verdict

**RunX does not obey GUIHO CLI Convention 0001.**

The repository has a strong implementation of the earlier GUIHO Go/Cobra and
11-artifact CLI contract, but the convention established on 2026-08-16 requires
a substantially different installation, configuration, agent-management,
upgrade, uninstall, release, and project-tooling architecture.

This is not a small set of cosmetic deviations. The repository is materially
noncompliant in the convention's central safety and lifecycle requirements:

- no stable launcher or immutable versioned payload layout;
- no artifacts ownership manifest or whole-release transaction;
- asynchronous Windows upgrade activation with the expressly forbidden
  scheduled outcome;
- incomplete and non-transactional installers;
- no compliant uninstallers;
- no global configuration, agent-evolution policy, schemas, or examples;
- no mandatory root RunX or current XDocs configuration;
- an obsolete agent command tree and incomplete agent artifacts; and
- release automation that intentionally rejects the convention's required
  artifact set.

The implementation should not be described as convention-compliant until every
Critical and High finding in this report is resolved and the complete
convention validation gate is rerun.

## Audit Scope

| Item | Audited value |
| --- | --- |
| Repository | C:\GUIHO\runx |
| Convention | C:\GUIHO\guiho\docs\conventions\guiho-convention-0001-cli.md |
| Checked-out branch | main |
| Checked-out head | 364fdb8d8038bafd984a2cc978631dc07488a7b0 |
| Local origin/main reference | 2c96e8aec1052b35c6fb4061eba5a6b54405bbcb |
| Checkout state | 104 commits behind local origin/main |
| Pre-existing user change | TODO.md modified; preserved and not included in this audit |
| Review mode | Read-only implementation review plus this requested report |
| Mutating lifecycle commands | Install, init, upgrade, uninstall, publish, version apply, and release were not run |

The verdict is bound to the checked-out worktree. A targeted comparison found
that the core distribution files responsible for the principal findings are
unchanged between the checked-out head and the local origin/main reference:
installers, updater, release builder and verifier, workflows, package.json, and
mirror.yaml. The mandatory missing root files and uninstall scripts are also
absent from the local origin/main tree. This was not a network refresh, and the
report does not claim that the local origin/main reference is the current
remote repository state.

The convention is newer than both observed commits. That history explains why
RunX still implements the previous contract, but it does not make the current
repository compliant.

## Authority And Interpretation

The convention was treated as the primary acceptance contract. Existing RunX
instructions and guiho-s-0035-cli-engineer-go still require the earlier exact
11-artifact release, direct replacement behavior, positive-integer help depth,
and agent update command names. Those older requirements conflict with the new
convention in several places. The conflict is itself a remediation dependency:
the canonical Go CLI engineering skill and repository instructions must be
aligned before implementation, or future agents will be told to recreate the
noncompliant design.

The review evaluated:

- every convention section and mandatory command;
- the live Cobra tree and actual flag behavior;
- Go entrypoints, command packages, manifest parsing, maintenance, update, and
  updater packages;
- lifecycle scripts and the npm delegator;
- release construction, checksums, CI, publication, and Mirror configuration;
- XDocs configuration, scan coverage, tree relationships, and doctor output;
- README, DOCS, bundled skill, prompt, embedded copies, and agent definition;
- existing tests and native build behavior; and
- absent required files and absent implementation concepts.

## Status Legend

| Status | Meaning |
| --- | --- |
| PASS | The inspected repository and behavior satisfy the convention clause. |
| PARTIAL | A useful implementation exists, but one or more mandatory details do not comply. |
| FAIL | A mandatory behavior, file, command, or safety property is absent or contradicts the convention. |
| NOT PROVEN | Repository evidence is insufficient to establish compliance. |
| N/A | The clause does not apply to the current public CLI surface. |

## Clause-By-Clause Compliance Matrix

| Convention area | Status | Audit result |
| --- | --- | --- |
| Placeholder CLI name | N/A | The cliname placeholder rule governs the generic convention examples, not RunX's public name. |
| Go technology stack | PASS | go.mod pins Go 1.26.5; production is Go. |
| Cobra command builder | PASS | cmd/root.go constructs one Cobra tree; no second public router was found. |
| Go standard tooling | PASS | Formatting, tidy diff, tests, vet, build, and cross-build tooling exist and passed in this audit. |
| Latest Go/Cobra at project creation | NOT PROVEN | Current versions are pinned, but historical latest-at-creation evidence is not durable. |
| Root mirror.yaml | PASS | Present; mirror config check passes. |
| Mirror as sole version authority | FAIL | package.json is a separate stale version field and publish.yml runs npm version. |
| Root runx.yaml | FAIL | File is absent. |
| Complete RunX command catalog | FAIL | Repository commands remain in package.json, CI, docs, and instructions without a root catalog. |
| Root xdocs.yaml | FAIL | File is absent; the repository uses legacy xdocs.config.toml. |
| XDocs directory coverage | FAIL | Three owned directories are uncovered and the GitHub descriptor is orphaned from the root tree. |
| Long kebab-case public Cobra flags | PASS | Current public custom flags use long kebab-case names. |
| Short-alias policy | PARTIAL | Cobra exposes only standard h and root v, but install.sh adds forbidden installer alias -v. |
| Space and equals flag values | PASS | Cobra supports both forms for value flags. |
| StringArray for list-valued flags | N/A | No current public flag accepts a list. |
| Raw top-level version output | PARTIAL | Injected release builds print raw SemVer; normal source builds print non-SemVer dev. |
| Help on every command | PASS | Cobra help exists at every inspected public scope. |
| help-docs on every command | PASS | Persistent help-docs is wired through the live tree. |
| Complete help-tree | PARTIAL | The live tree is generated, but global-flag rendering violates the required default. |
| help-tree-depth contract | FAIL | max is rejected and 1 is accepted. |
| help-tree-global-flags | FAIL | Flag is absent and inherited flags repeat by default. |
| Four lifecycle scripts | FAIL | devops/uninstall.sh and devops/uninstall.ps1 are absent. |
| No Cobra install command | PASS | No install command exists in the Cobra tree. |
| Installer exact-version selection | PARTIAL | Exact version exists, but latest is not resolved then exactly verified. |
| Installer channel selection | FAIL | No channel option or complete channel selection exists. |
| Complete release artifact set | FAIL | Builder and verifier enforce the obsolete exact-11 set. |
| artifacts.json ownership manifest | FAIL | No manifest exists or is consumed. |
| Transactional installer behavior | FAIL | Binary backup is discarded before agent-resource and PATH reconciliation. |
| Reinstall persistence and repair | FAIL | No manifest-based persistent/disposable classification or complete repair transaction exists. |
| Shared uninstall contract | FAIL | Scripts are absent; Cobra removes only os.Executable(). |
| Shared .guiho/.temp staging | FAIL | Installers use OS temp and upgrade stages beside the executable. |
| Shared .guiho/bin entrypoint | FAIL | Installers default to .local/bin and install the payload directly. |
| Stable launcher | FAIL | No conventional launcher implementation exists. |
| Immutable versions directory | FAIL | No versions/<version> payload layout exists. |
| Atomic current.json pointer | FAIL | No current.json contract exists. |
| Shared .guiho ownership boundary | FAIL | Ownership is inferred from fixed paths/names without artifacts.json. |
| CLI home canonical artifacts | FAIL | Most release artifacts have no canonical installed copy under .guiho/runx. |
| Separate global/project configuration | FAIL | One runx.yaml schema is used as project-or-global fallback. |
| Global-overridden-by-project inheritance | FAIL | Configuration files are selected, not merged. |
| Agent-evolution policy | FAIL | Policy fields and their three values do not exist. |
| Separate JSON Schemas | FAIL | runx.schema.json and runx.global.schema.json are absent. |
| Complete configuration examples | FAIL | Required global/project examples are absent. |
| Version-pinned schema references | FAIL | init writes no schema reference. |
| AI-managed lifecycle skill | FAIL | Main skill lacks installation, setup, policy, feedback, and post-upgrade lifecycle guidance. |
| One managed agent instruction | PARTIAL | A bounded instruction exists, but it is hard-coded and differs from the published Markdown asset. |
| Main agent skill | PARTIAL | The skill exists, but mandatory CLI Evolution and Feedback content is absent. |
| Main setup prompt | FAIL | The only prompt covers catalog use, not install, verify, init, and upgrade. |
| Required command roots | PARTIAL | init, agent, upgrade, and uninstall exist, but their mandatory behavior is incomplete. |
| agent skill upgrade | FAIL | The tree uses prohibited update. |
| agent instruction upgrade | FAIL | The tree uses prohibited update. |
| init reconciliation | FAIL | init only creates a new project runx.yaml and errors if it already exists. |
| Whole-release upgrade | FAIL | Upgrade downloads and replaces only the payload. |
| Mandatory recovery block | FAIL | It is not the first output and is omitted from successful, dry-run, and up-to-date outcomes. |
| Synchronous upgrade activation | FAIL | Windows uses a detached helper and returns scheduled. |
| Uninstall preservation flags | FAIL | preserve-config, preserve-data, and yes are absent. |
| README lifecycle contract | FAIL | No final Uninstall section or required examples exist. |

## Findings Summary

| Severity | Count | Meaning |
| --- | ---: | --- |
| Critical | 4 | Foundational lifecycle or ownership architecture contradicts the convention and cannot be made compliant by documentation alone. |
| High | 16 | Mandatory repository, configuration, agent, installer, release, or validation capability is absent or materially incomplete. |
| Medium | 4 | Public behavior, documentation, or tests encode a superseded contract. |

## Critical Findings

### C-01 — Stable launcher and immutable versioned payload architecture are absent

Convention clauses: lines 315-341, 418-514, and 802-832.

Evidence:

- devops/install.sh:6 defaults to $HOME/.local/bin.
- devops/install.sh:66-70 installs and replaces $INSTALL_DIR/runx directly.
- devops/install.ps1:3 defaults to $HOME\.local\bin.
- devops/install.ps1:154-155 installs and replaces runx.exe directly.
- pkg/updater/rollback.go:94-160 stages beside and renames over the current
  executable path.
- scripts/runx-bin.mjs:21-46 downloads and delegates to a raw payload under
  .guiho/runx/npm/<package-version>; it is not the required stable launcher.
- Repository search found no current.json, installed-artifacts.json,
  immutable versions/<version> activation implementation, launcher protocol,
  launcher fallback, or launcher source.

Impact:

- activation changes the executable instead of atomically changing a pointer;
- a missing or broken active payload cannot fall back through a stable launcher;
- the immediately previous verified version is not durably recorded;
- ordinary upgrades depend on replacing the running command path; and
- installer-driven repair and ordinary upgrade do not share one launcher
  protocol.

Required compliance outcome:

~~~text
$HOME/.guiho/bin/runx[.exe]                       stable launcher
$HOME/.guiho/runx/current.json                    atomic active/previous pointer
$HOME/.guiho/runx/installed-artifacts.json        installed ownership state
$HOME/.guiho/runx/versions/<version>/runx[.exe]   immutable payload
$HOME/.guiho/runx/versions/<version>/artifacts/   complete selected release
~~~

The launcher must preserve arguments and standard streams, wait for the
payload, and return the payload's exact exit code.

### C-02 — Upgrade is not a synchronous whole-release transaction

Convention clauses: lines 769-832.

Evidence:

- pkg/updater/upgrade.go:93-114 performs release-network work before the
  recovery message can be rendered.
- pkg/updater/upgrade.go:156-169 requires only a compatible payload and
  checksums.txt.
- pkg/updater/upgrade.go:219-314 downloads only those two artifacts.
- pkg/updater/upgrade.go:324-332 sends raw payload bytes into direct executable
  replacement.
- pkg/updater/upgrade.go:378-382 explicitly produces outcome scheduled on
  Windows.
- pkg/updater/stage_windows.go:14-27 launches a detached PowerShell helper and
  returns before replacement and verification complete.
- pkg/updater/rollback.go:65-75 accepts version output containing the target
  text rather than requiring exact raw output.
- No exclusive lock with an ownership token, stale-owner process recovery,
  instance registry, verified executable-path process termination, immutable
  version install, artifact snapshot, current.json activation, transaction
  journal, post-activation self-test, or launcher fallback exists.

The convention expressly prohibits scheduled as an upgrade result and forbids
a detached helper from being the authority for upgrade success.

The current upgrade changes only one binary. It does not remove retired
resources or replace the complete selected-release skill, instruction, prompt,
agent-definition, schema, example, manifest, metadata, and projection set.

### C-03 — Installer rollback is not transactional across the installation

Convention clauses: lines 315-370.

Evidence:

- devops/install.sh:58-63 downloads only payload, checksums.txt, and the skill
  ZIP.
- devops/install.sh:66-71 replaces the payload and deletes its backup before
  skill extraction, skill projection, instruction reconciliation, and PATH
  work at lines 73-76.
- devops/install.ps1:151-153 downloads the same partial set.
- devops/install.ps1:154-155 deletes the executable backup before skill,
  instruction, and PATH operations at lines 156-169.
- devops/install.sh:74 and devops/install.ps1:157 recursively delete existing
  skill projections by fixed name before copying replacements.
- There is no installed manifest to identify the previous release's complete
  owned set, no snapshot of replaceable projections, and no rollback of agent
  resources or PATH changes.

A skill extraction, projection, instruction, or PATH failure after binary
activation can leave a mixed installation. The previous release cannot be
restored as a whole.

### C-04 — Uninstallation is unsafe, incomplete, and missing two interfaces

Convention clauses: lines 238-262, 372-398, and 834-836.

Evidence:

- devops/uninstall.sh is absent.
- devops/uninstall.ps1 is absent.
- cmd/upgrade.go:144-178 implements the Cobra uninstall command.
- cmd/upgrade.go:151-164 resolves os.Executable() and removes only that path.
- The live uninstall help exposes only dry-run and format.
- There is no yes, preserve-config, or preserve-data option.
- There is no grouped REMOVE/PRESERVE plan.
- There is no interactive confirmation.
- There is no noninteractive fail-closed confirmation rule.
- There is no cleanup of launcher, versioned payloads, CLI home, project/global
  configuration, caches, persistent data, databases, skills, prompts, agent
  definitions, instruction block, schemas, examples, or operation-owned temp.
- There is no installed artifact manifest to prove ownership before deletion.

The current behavior is incomplete on Unix and cannot reliably remove a
running Windows executable. It also cannot meet the convention's requirement
to preserve shared .guiho, .guiho/bin, .guiho/.temp, shared PATH, other CLI
artifacts, and unmanaged AGENTS.md content.

## High Findings

### H-01 — Mandatory root project tooling is missing

Convention clauses: lines 99-159.

Observed root files:

~~~text
mirror.yaml          present
runx.yaml            absent
xdocs.yaml           absent
xdocs.config.toml    present legacy configuration
~~~

Consequences:

- the repository cannot pass runx check --format json or runx list --format
  json against a project-owned catalog;
- supported commands in package.json, CI, publication workflows, README,
  CONTRIBUTING, DOCS, and agent instructions are not cataloged;
- unfamiliar and high-impact operations cannot follow the mandatory RunX
  discovery and dry-run workflow; and
- the XDocs configuration filename and schema do not satisfy the convention.

mirror config check passes, so this is not a claim that all mandatory tooling
is absent. It is a failure of the required three-file baseline and the RunX and
XDocs portions of that baseline.

### H-02 — Mirror is not the sole authority over every version-bearing file

Convention clauses: lines 111-124.

Evidence:

- mirror.yaml:5-9 uses Git as the sole version source and output.
- package.json:4 contains version 0.8.0.
- CHANGELOG.md begins with a newer 0.11.0 section in this checkout.
- .github/workflows/publish.yml:67-70 runs npm version with
  --no-git-tag-version immediately before npm publication.

The workflow uses a prohibited package-manager version mutation because Mirror
does not describe package.json as a version-bearing project file.

mirror config check passes. mirror version current does not pass in this
checkout because the locally known 0.12.1 tag is not reachable from the
104-commit-behind head. That reachability failure is a stale-checkout
condition, separate from the valid Mirror syntax check.

### H-03 — XDocs coverage and tree relationships are incomplete

Convention clauses: lines 143-159.

xdocs scan reported:

~~~text
total directories: 32
covered directories: 29
uncovered directories: 3
~~~

The uncovered project-owned directories are:

~~~text
embed/prompts
embed/skills
skills/guiho-s-runx/agents
~~~

These are prompt, skill, and agent directories that the convention explicitly
classifies as project-owned.

There is also a tree inconsistency:

- .github/github.xdocs.md:2-6 declares subject runx-github with parent runx
  and child runx-github-workflows.
- runx.xdocs.md:5-14 does not list runx-github as a child.
- xdocs tree omits the .github branch entirely.

xdocs meta, tree, and doctor pass for the represented 29-descriptor legacy
tree. Those green results do not prove convention coverage because the three
owned directories and orphaned GitHub branch are outside that tree.

### H-04 — Installers and upgrade do not support release channels

Convention clauses: lines 264-293.

Evidence:

- devops/install.sh:17-29 supports version and install-dir, but not channel.
- devops/install.sh:23 also defines the forbidden installer short alias -v.
- devops/install.ps1:1-5 supports Version and InstallDir, but not Channel.
- cmd/upgrade.go:42-45 exposes version, dry-run, and format, but not channel.
- A live upgrade --channel stable probe failed with unknown flag.
- The installers use releases/latest/download for the default instead of
  exhausting the release catalog and selecting the highest valid stable
  SemVer.

The internal release catalog can display a channel label, but no installation
or upgrade selector implements the convention's channel contract.

### H-05 — Installer staging, destination, and ownership boundaries are wrong

Convention clauses: lines 400-514.

Evidence:

- devops/install.sh:55 stages below system TMPDIR or /tmp.
- devops/install.ps1:148 stages below the operating-system temp directory.
- pkg/updater/rollback.go:94-102 stages beside the current executable.
- The installers do not validate staging as a strict descendant of
  $HOME/.guiho/.temp/.
- Both installers default to .local/bin rather than $HOME/.guiho/bin/.
- Both allow an arbitrary install directory rather than enforcing the shared
  GUIHO command-entrypoint location.
- The PATH mutation targets the configured direct-payload directory, not the
  canonical shared GUIHO bin directory.
- No artifacts.json binds canonical source, installed path, managed
  projection, checksum, artifact ID, version, and ownership boundary.

### H-06 — Release construction enforces an obsolete and incomplete asset set

Convention clauses: lines 295-313.

Evidence:

- devops/build-binaries.go:24-29 correctly defines eight payload targets.
- devops/build-binaries.go:65-82 adds only a skill ZIP, one Markdown file, and
  checksums.txt, then requires exactly 11 files.
- devops/verify-release-assets.go:19-56 hard-codes the same set and rejects any
  additional file.
- .github/workflows/publish.yml:30-53 builds, uploads, and remotely compares
  that exact set.
- README.md:134-147 and DOCS.md document the old exact-11 contract.

Missing from the complete convention release unit are at least:

- platform launcher artifact or artifacts;
- artifacts.json;
- compliant canonical managed instruction source;
- compliant main setup prompt;
- separately declared artifact metadata for the bundled agent definition;
- runx.schema.json;
- runx.global.schema.json;
- complete project and global configuration examples;
- canonical installation and upgrade metadata; and
- ownership/projection declarations for every contained artifact.

The agent definition skills/guiho-s-runx/agents/openai.yaml is contained in the
skill ZIP. The violation is not that no agent definition exists; it is that no
artifacts.json enumerates its path, ID, version, checksum, canonical installed
path, managed projection, and ownership boundary.

### H-07 — Reinstall and installer-driven upgrade semantics are absent

Convention clauses: lines 343-370.

The installers perform direct replacement, not a classified reinstall:

- no installed-manifest inventory exists;
- no path is classified persistent or disposable;
- no global/project configuration preservation contract exists;
- no persistent-data or database preservation contract exists;
- no retired version/resource/projection cleanup exists;
- no cache/disposable-state removal contract exists;
- no selected-release canonical resource replacement exists;
- no complete rollback exists; and
- no post-install runx init reconciliation occurs.

Re-running an installer may replace the binary and current skill copies, but
that is not the convention's controlled, idempotent whole-install repair.

### H-08 — init implements only old catalog creation

Convention clauses: lines 710-737.

cmd/catalog.go:303-362:

- chooses a project runx.yaml path;
- fails if that path already exists at lines 321-325;
- writes an empty manifest at lines 326-349; and
- prints the created path.

It does not:

- resolve and validate the project root;
- ensure every global skill projection;
- create AGENTS.md when absent and reconcile the installed instruction;
- create and validate runx.global.yaml;
- validate or merge project configuration;
- establish agent.evolution defaults;
- explain and collect evolution-policy choices;
- ask remaining domain setup questions;
- preserve existing valid values;
- write version-pinned schema references;
- perform a final full reconciliation; or
- behave idempotently when setup is already present.

Its noninteractive behavior can create an incomplete configuration instead of
failing for unanswered mandatory choices.

### H-09 — Global/project configuration, inheritance, schemas, and examples are absent

Convention clauses: lines 516-600.

Evidence:

- pkg/manifest/types.go:5-12 defines one domain catalog schema.
- pkg/manifest/parser.go:28-45 strictly decodes that schema and would reject
  an agent-evolution field.
- pkg/manifest/composition.go:45-81 resolves project runx.yaml and
  $HOME/.guiho/runx/runx.yaml as alternatives.
- The implementation does not distinguish runx.yaml from runx.global.yaml.
- No baseline-plus-project overlay is performed.
- runx.schema.json is absent.
- runx.global.schema.json is absent.
- Complete global and project examples are absent.
- cmd/catalog.go:326-328 writes no version-pinned HTTPS schema reference.
- No release tooling publishes schemas or examples.

Strict decoding of the domain manifest is good, but it does not satisfy the
new separate global/project configuration contract.

### H-10 — The agent-evolution policy is completely absent

Convention clauses: lines 529-577.

Repository search found no implementation of:

~~~text
agent.evolution
disabled
always-ask
always-proceed
issues.bugs
issues.improvements
issues.reviews
~~~

Consequences:

- init cannot establish or explain the policy;
- upgrade check cannot honor disabled;
- upgrade cannot distinguish persistent preauthorization from per-action
  approval;
- issue creation cannot be governed independently by category;
- the skill cannot teach correct agent authority; and
- cmd/root.go:226-245 schedules update work without consulting an effective
  policy.

### H-11 — Required agent command names and show semantics do not comply

Convention clauses: lines 680-768.

Evidence:

- cmd/agent.go:24-27 registers agent skill update.
- cmd/agent.go:88-123 implements update instead of upgrade.
- cmd/agent.go:174-218 registers agent instruction update instead of upgrade.
- cmd/help.go:118-123 orders and documents the obsolete names.
- Live help shows update and does not show upgrade under either subtree.
- cmd/agent.go:51-69 makes agent skill show print metadata and a source path
  rather than the selected raw bundled skill.
- cmd/agent.go:158-171 selects CLAUDE.md alone when it is the only existing
  instruction file; mandatory init instead requires ensuring AGENTS.md.

The convention explicitly prohibits update in the agent skill and instruction
trees.

### H-12 — Required agent lifecycle artifacts and guidance are incomplete

Convention clauses: lines 602-678.

Main skill gaps in skills/guiho-s-runx/SKILL.md:

- no exact heading named CLI Evolution and Feedback;
- no canonical repository URL;
- no canonical issue-creation URL;
- no bug/improvement/review recognition guidance;
- no instruction to read effective agent.evolution;
- no behavior for disabled, always-ask, and always-proceed;
- no policy-governed upgrade flow;
- no install-from-README and version-verification lifecycle;
- no post-upgrade runx init and raw-version verification; and
- no requirement to return the created issue URL.

Main prompt gaps in prompts/guiho-i-runx.md:

- describes catalog inspection and execution only;
- does not explain installation;
- does not explain installation verification;
- does not explain init/setup; and
- does not explain upgrade.

Instruction identity is also split:

- pkg/maintenance/maintenance.go:25-36 hard-codes the runtime managed block.
- devops/build-binaries.go:70-74 publishes prompts/guiho-i-runx.md as
  guiho-i-runx.md.
- cmd/agent.go:219-220 shows the hard-coded block, not the published Markdown
  file.

No durable repository record proves that the user confirmed the CLI home name,
main skill ID, or main prompt ID required by the convention.

### H-13 — Upgrade recovery output violates both timing and coverage

Convention clauses: lines 777-800.

Evidence:

- cmd/upgrade.go:17-31 enters UpgradeSelf before printing recovery.
- pkg/updater/upgrade.go:93-114 can make a network request before any recovery
  text is rendered.
- cmd/upgrade.go:139-140 prints recovery only for failed or rolled-back
  outcomes.
- Success, up-to-date, and dry-run outcomes omit the recovery block.
- The text is rendered as recovery and stop labels, not the mandatory two-line
  reinstall message.
- Default and channel selector preservation is not implemented.

The recovery block must be the first user-visible output after local validation
and the final block for every terminal outcome, pinned to the exact resolved
version once known.

### H-14 — Candidate and release selection validation is weaker than required

Convention clauses: lines 264-341 and 802-832.

Evidence:

- pkg/update/catalog.go:150-188 compares prerelease text lexicographically
  rather than by SemVer identifier rules.
- pkg/update/catalog.go:190-206 ignores numeric parse errors and accepts missing
  components as zero.
- pkg/update/catalog.go:263-300 includes non-draft tags without rejecting
  malformed SemVer and identifies completeness using only payload plus
  checksums.
- pkg/updater/upgrade.go:121-135 selects from that permissive catalog.
- pkg/updater/rollback.go:65-75 uses substring version verification instead of
  exact raw output.
- devops/install.sh:69-70 and devops/install.ps1:155 do not resolve a concrete
  SemVer and compare exact output when the requested selector is latest.
- No candidate hidden self-test exists.
- No architecture/build-target verification accompanies native executable
  magic checks.
- PowerShell checksum lookup selects the first matching line and does not reject
  duplicate entries.
- Upgrade verifies only the selected payload checksum, not every release
  artifact.

Full GitHub pagination is implemented in pkg/update/catalog.go, which is a
positive result. The problem is validation and selection of the complete
release after pagination.

### H-15 — The npm command entrypoint bypasses the mandatory lifecycle

Convention clauses: lines 295-341 and 418-514.

package.json:6-8 publishes scripts/runx-bin.mjs as the runx command.

The shim:

- supports only x64 and ARM64 at scripts/runx-bin.mjs:15-17;
- stores a raw payload in .guiho/runx/npm/<version> at lines 21-22;
- trusts an existing cached payload based on file existence at lines 27-30;
- downloads only payload and checksums at lines 31-41;
- performs foreground network installation on first use;
- does not use the stable launcher/current.json protocol;
- does not obtain artifacts.json, schemas, examples, prompt, instruction, or
  canonical agent resources;
- does not run exact version verification or a hidden self-test; and
- does not reconcile init after installation.

The shim correctly forwards arguments, standard streams, signals, and exit
status at lines 46-50. That delegation quality does not make its installation
path convention-compliant.

### H-16 — Embedded runtime and released skill sources have drifted

Convention clauses: lines 295-313 and 602-678.

Evidence:

- embed/embed.go:7-10 embeds embed/skills/guiho-s-runx.SKILL.md.
- cmd/root.go:248-250 reads that embedded copy for automatic installation.
- devops/build-binaries.go:65-67 creates the release ZIP from
  skills/guiho-s-runx/.
- The current hashes differ:

~~~text
embed/skills/guiho-s-runx.SKILL.md
A15CDDC9E7694140C25E62264CCB1927BF256FE25D279AAD4D36EAA5938D4D6D

skills/guiho-s-runx/SKILL.md
D550DA42E600E55DFF5F17ADCC91C8512FCAEF581A433B1E60B8EBE2CC4FDAF1
~~~

An installer can install the release ZIP skill, after which an ordinary RunX
invocation can overwrite it with different embedded bytes. A selected release
therefore does not have one self-consistent canonical skill artifact.

The top-level and embedded guiho-i-runx prompt hashes match in this checkout;
the confirmed drift is the skill pair.

## Medium Findings

### M-01 — Help-tree depth and global-flag behavior violate the public contract

Convention clauses: lines 200-235.

Evidence:

- cmd/root.go:104-106 and 168-170 model depth as an integer with sentinel 0.
- cmd/root.go:113-116 accepts any positive integer, including 1.
- The convention permits max or an integer greater than 1.
- A live help-tree-depth max probe failed as invalid integer input.
- A live depth 1 probe succeeded.
- help-tree-global-flags is not registered and a live probe failed as unknown.
- cmd/help.go:50-91 includes inherited flags beneath descendants by default.
- Live depth-two output repeated inherited help flags under every command,
  contrary to the default false behavior.
- cmd/root_test.go:157-170 explicitly tests and preserves depth 1.

The live tree generation and Unicode rendering are useful and should be
retained when correcting the flag contract.

### M-02 — Development builds do not report a SemVer-compatible version

Convention clauses: lines 188-198.

Evidence:

- main.go:11-15 initializes version to dev.
- cmd/root.go:98-100 also falls back to dev.
- go run . --version printed exactly dev.

The output contains no labels, ANSI, CLI name, or v prefix, so its shape is
otherwise correct. Linker-injected release builds print raw SemVer; the defect
is the normal source/development build value and the unproven Go module-version
fallback.

### M-03 — README lifecycle documentation is noncompliant

Convention clauses: lines 238-262 and 372-398.

README.md:

- provides separate remote installer commands;
- does not show a final runx --version command in the primary install sequence;
- has no final section named Uninstall;
- has no remote PowerShell uninstaller command;
- has no remote Linux/macOS uninstaller command;
- does not disclose that default uninstall removes all CLI-owned data;
- has no destructive-default example;
- has no combined preserve-config plus preserve-data example;
- documents obsolete agent update command names; and
- documents the obsolete exact-11 release set.

The final operational section is Build at README.md:134.

### M-04 — CI and tests enforce the superseded contract

Evidence:

- cmd/root_test.go:157-170 expects help-tree depth 1.
- cmd/root_test.go:380-389 expects agent update commands.
- devops/verify-release-assets.go:19-56 rejects any release that is not the old
  exact-11 set.
- .github/workflows/ci.yml:41-43 and 54-59 smoke only positive-integer depth and
  the old release matrix.
- .github/workflows/publish.yml publishes and verifies the incomplete asset set.
- CI does not run mirror config check, root RunX check/list, convention-complete
  XDocs scan coverage, or strict XDocs gates against xdocs.yaml.

There are no tests for:

- launcher/current.json fallback;
- immutable payload activation;
- artifacts.json ownership;
- complete artifact install/retirement/rollback;
- version and channel selection;
- global/project configuration inheritance;
- agent-evolution policy;
- lifecycle uninstaller scripts;
- uninstall preservation and confirmation;
- synchronous Windows upgrade;
- upgrade lock and instance registry; or
- first/final recovery output across every outcome.

## Confirmed Conforming Or Reusable Foundations

The noncompliant verdict should not obscure the parts worth retaining:

- go.mod defines the current Go module and pinned toolchain.
- cmd/root.go:96-180 constructs one dependency-injected Cobra tree.
- main.go is a thin entrypoint and centralizes exit handling.
- pkg/manifest/parser.go:28-45 uses KnownFields true.
- pkg/manifest/parser.go:48-128 performs explicit semantic validation.
- Public Cobra custom flags use full kebab-case names.
- Only root v and Cobra help h are public short aliases.
- Cobra accepts both spaced and equals value forms.
- No Cobra install command exists.
- Every current public scope has help, help-tree, and help-docs.
- Release builds can print only a raw injected SemVer.
- The eight payload targets and CPU tuning are correct.
- CGO is disabled for the eight release payloads.
- Checksum generation is deterministic and excludes checksums.txt itself.
- GitHub release catalog retrieval exhausts pagination.
- Draft releases are excluded from catalog retrieval.
- Cache and maintenance paths already use .guiho/runx in several runtime areas.
- Instruction-block replacement is bounded, strict, idempotent, and preserves
  line endings in its current limited scope.
- The npm shim preserves child arguments, standard streams, signals, and exit
  status.
- The current Go suite, vet, formatting, tidy diff, build, old release matrix,
  and native Windows AMD64 smoke are green.

These foundations reduce implementation risk, but none overrides a mandatory
failed convention clause.

## Validation Evidence

### Passed

| Command or probe | Result |
| --- | --- |
| gofmt -l main.go cmd pkg embed devops | Clean; no files printed. |
| go mod tidy -diff | Clean; no diff. |
| go test -count=1 ./... | Passed all Go packages. |
| go vet ./... | Passed. |
| go build ./... | Exit 0; Windows module stat-cache writes emitted access-denied warnings but build succeeded. |
| mirror config check | Passed. |
| xdocs meta . --documents --strict | Passed for 29 represented descriptors. |
| xdocs tree | Passed for the represented tree. |
| xdocs doctor . | Reported valid with zero errors/warnings for the represented tree. |
| Eight-target build | Passed for version 0.0.0-audit with the standard eight Go targets. |
| Old release verifier | Passed exactly 11 assets and all old-set checksums. |
| Windows AMD64 native smoke | Raw version 0.0.0-audit and help tree passed. |

The temporary audit release directory was removed after verification.

### Failed Or Convention-Incomplete

| Command or probe | Result |
| --- | --- |
| Project runx check --format json | Cannot pass because root runx.yaml is absent. |
| Project runx list --format json | Cannot pass because root runx.yaml is absent. |
| xdocs scan | 32 owned directories, 29 covered, 3 uncovered. |
| go run . --version | Printed non-SemVer dev. |
| help-tree-depth max | Rejected as invalid integer. |
| help-tree-depth 1 | Incorrectly accepted. |
| help-tree-global-flags | Unknown flag. |
| upgrade --channel stable | Unknown flag. |
| agent skill/instruction help | Exposed update rather than upgrade. |
| uninstall help | Missing yes, preserve-config, and preserve-data. |
| mirror version current | Failed because the known 0.12.1 tag is not reachable from this stale checkout. |
| Convention release verification | Impossible with current verifier because it rejects assets beyond the old exact-11 set. |

### Not Executed

The following were intentionally not run because they would change installed
software, user PATH, global agent directories, repository instructions,
version/tag state, or remote publication:

- install.sh and install.ps1;
- missing uninstall scripts;
- mutating runx init;
- mutating runx upgrade;
- mutating runx uninstall;
- plain bootstrap invocations that reconcile global agent resources;
- Mirror version apply;
- tag or GitHub Release creation;
- npm publication; and
- production deployment.

Foreign binaries were cross-built but not executed. Only Windows AMD64 received
a native runtime smoke in this audit.

## Evidence Gaps And Non-Findings

- The audit did not prove whether Go 1.26.5 and Cobra 1.10.2 were the latest
  releases at the historical moment the project was created.
- No durable record was found for the required human confirmation of the CLI
  home name, main skill ID, or main prompt ID.
- Live remote GitHub Release assets were not fetched; release findings are
  based on the source builder, verifier, workflows, installers, and updater.
- The convention does not require shell completions or man pages. Their absence
  is not a finding.
- No current list-valued public flag exists, so StringArray behavior is not
  applicable.
- Signature, notarization, and attestation requirements are not present in this
  convention; their absence is not reported as a violation.

## Required Remediation Order

The gaps are interdependent. Implementing them as isolated patches would be
unsafe. The recommended order is:

1. **Align governing contracts.**
   Update the canonical Go CLI engineering skill, RunX AGENTS.md instructions,
   architecture, and approved plan so they no longer require the superseded
   exact-11/direct-replacement/update-command design.

2. **Define the complete configuration and artifact contracts.**
   Establish runx.yaml, runx.global.yaml, their separate schemas and examples,
   agent.evolution, artifact IDs, artifacts.json, persistent/disposable paths,
   canonical home layout, managed projections, and release completeness rules.

3. **Implement the stable launcher foundation.**
   Add the platform launcher, immutable version directories, strict current.json,
   active/previous fallback, instance registry, locks, transaction journal, and
   candidate/post-activation self-tests.

4. **Rebuild install and uninstall around ownership.**
   Implement both installers and both uninstaller scripts against .guiho/bin,
   .guiho/.temp, the CLI home, artifacts.json, complete rollback, preservation,
   confirmation, and shared-directory boundaries.

5. **Rebuild upgrade as a synchronous whole-release transaction.**
   Add exact-version/channel resolution, first/final recovery blocks, complete
   downloads/checksums, safe old-process handling, full artifact replacement,
   atomic pointer activation, rollback, verification, and post-upgrade init.

6. **Correct the agent namespace and resources.**
   Rename update to upgrade, make show return raw artifacts, create one canonical
   instruction, eliminate embedded/release drift, add the compliant setup
   prompt and CLI Evolution and Feedback section, and enforce policy-aware
   agent behavior.

7. **Complete repository tooling and delivery gates.**
   Add the root RunX and XDocs YAML files, catalog every workflow, cover every
   owned XDocs directory, repair the GitHub tree link, place all version-bearing
   files under Mirror, replace the old release verifier/workflow, update README
   and DOCS, and add convention-specific tests.

After those units are independently reviewed, rerun every validation gate
against one exact head, including native platform tests where runners or
hardware are available.

## Final Compliance Statement

RunX is a capable Go/Cobra CLI and its legacy validation suite is healthy. It
is nevertheless **not compliant** with GUIHO CLI Convention 0001.

The decisive blockers are architectural: the convention treats installation
as a manifest-owned, complete-release, stable-launcher transaction, while RunX
currently treats installation and upgrade primarily as verified replacement of
one directly invoked binary plus selected skill copies. Configuration,
agent-evolution authority, release completeness, and uninstall safety are also
missing.

Compliance requires an approved architecture and migration plan followed by a
fresh exact-head implementation review and validation report. Green legacy Go
tests and the old 11-artifact verifier must not be used as evidence that the new
convention is satisfied.
