#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# Mandume

> **All agents run by default with full permission** — full workspace read/write/execute, no sandbox or approval prompts. Every worker `Command` is executed with the same permissions as the calling agent.

## Available Harnesses

| Harness | Purpose                                                                                                                        |
| ------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `pi`    | pi harness — also runs Codex models (`gpt-5.6-luna`, `gpt-5.6-sol`, `gpt-5.6-terra`, `gpt-5.3-codex-spark` via `openai-codex`) |
| `agy`   | Antigravity harness                                                                                                            |

## User Conventions

User conventions are stored under `conventions/` in repository `CGuiho/guiho`.

* **Local:** `C:/GUIHO/guiho/conventions` — read directly, e.g. `cat /c/GUIHO/guiho/conventions/<file>.md`
* **Remote:** if `C:/GUIHO/guiho` is not present locally, the repository is remote at `github.com/CGuiho/guiho` — read via GitHub CLI, e.g. `gh api repos/CGuiho/guiho/contents/conventions --jq ".[].name"` or `gh api repos/CGuiho/guiho/contents/conventions/<file>.md --jq .content | base64 -d`

Check local first, then remote via `gh` if missing.

## AGENTS

### Mastermind Agent

The Mastermind agent handles judgment. It is the strong model used for planning, architecture, spec review, plan review, self-review, and any ambiguous or high-risk decision during execution. It reads full repository context, resolves uncertainty with the safest reversible choice, and records the rationale to the question ledger (`docs/questions/`). Use a `mastermind` worker when correctness and reasoning matter more than speed or cost.

### Workhorse Agent

The Workhorse agent handles labor. It is the efficient model used for well-specified implementation: code generation, refactors, schema changes, tests, and validation runs. It follows the plan unit exactly, runs with full workspace permissions via `guiho-s-0440-hand-off`, and returns verified output for the executor to integrate. Use a `workhorse` worker when the task is mechanical and speed/cost matter more than deep reasoning.

## Workers

Each worker is one `Harness` + `Provider` + `Model` + `Thinking` + `Command`. The `Base Command` is just the base command to give an idea — you can specify permissions or not — the agent constructs the full command on the go, injecting the model, thinking effort, and instructions/brief at runtime.

> **Only `primary` matters for delivery.** `fallback` worker (AG) is included only to spend the subscription already paid for — money has been invested and needs to be used. Codex models (`openai-codex`) are run through `pi`, no separate Codex CLI.

### Mastermind Workers

| Worker     | Harness | Provider       | Model                             | Thinking | Base Command                                   | Usage           |
| ---------- | ------- | -------------- | --------------------------------- | -------- | ---------------------------------------------- | --------------- |
| `cto`      | `pi`    | `meta`         | `muse-spark-1.2-contributor`      | `xhigh`  | `pi --tools read,bash,edit,write,grep,find,ls` | `paid-api-rate` |
| `senior`   | `pi`    | `openai-codex` | `gpt-5.6-sol`                     | `xhigh`  | `pi --tools read,bash,edit,write,grep,find,ls` | `5h-limited`    |
| `senior-2` | `pi`    | `opencode`     | `muse-spark-1.2-contributor-free` | `xhigh`  | `pi --tools read,bash,edit,write,grep,find,ls` | `free`          |

### Workhorse Workers

| Worker       | Harness | Provider    | Model                             | Thinking       | Base Command                                   | Usage           |
| ------------ | ------- | ----------- | --------------------------------- | -------------- | ---------------------------------------------- | --------------- |
| `engineer`   | `pi`    | `zai`       | `glm-5.3-flash`                   | `max`          | `pi --tools read,bash,edit,write,grep,find,ls` | `paid-api-rate` |
| `engineer-2` | `agy`   | `anthropic` | `claude-opus-4-6-thinking`        | _`do not set`_ | `agy --dangerously-skip-permissions --print`   | `5h-limited`    |
| `engineer-3` | `agy`   | `google`    | `gemini-3.7-flash-high`           | `high`         | `agy --dangerously-skip-permissions --print`   | `5h-limited`    |
| `engineer-4` | `pi`    | `opencode`  | `muse-spark-1.2-contributor-free` | `xhigh`        | `pi --tools read,bash,edit,write,grep,find,ls` | `free`          |

### Calling Workers -- Subagents strategy

**There are no native pi subagents. There is one single pi instance — the orchestrator — that manages every task.**

The orchestrator (this `pi` session) never spawns internal subagents. It delegates to **external CLI workers as subagents** via `guiho-s-0440-hand-off`. Each worker is a separate process (`pi` or `agy` with its model/thinking) that receives a self-contained brief, runs with full permission, and returns. The orchestrator stays alive to answer you and to manage the lifecycle.

**How it works:**

1. **You ask / orchestrator plans** — `cto`/`senior` (mastermind) reads `docs/specifications/`, `docs/architecture/`, `MANDUME.md`, writes `docs/plans/<plan>/plan.md` and `todo.md`. You can ask at any time _"Why are you working?"_ — the orchestrator answers from `docs/plans/`, `todo.md`, `docs/questions/` without interrupting workers.

2. **Orchestrator delegates one unit at a time** — picks a worker by `Class` (`mastermind` for judgment, `workhorse` for labor) and `Usage` (`paid-api-rate` primary, `5h-limited`/`free` fallback to spend subscriptions). Writes a brief to `docs/plans/<plan>/execution/handoffs/<date>-<unit>-<worker>.md` with Goal, Context, Skills to load, Constraints, Output contract.

3. **Call the worker synchronously with full permission but stream live by default so you stay responsive (unless the harness does not allow it):**

   ```bash
   # pi mastermind/workhorse — stream to file so orchestrator can tail while waiting
   pi --model <provider>/<model>:<thinking> --tools read,bash,edit,write,grep,find,ls -p --no-session "$(cat brief.md)" > "out-<worker>.log" 2>&1 &
   tail -n 80 -f out-<worker>.log &  # live stream, does not block
   wait  # still waits for exit, but you can answer "What are you doing?" from tail
   # agy workhorse (Opus 4.6 / Gemini 3.7 Flash High) — same streaming pattern
   agy --model claude-opus-4-6-thinking --dangerously-skip-permissions --print "$(cat brief.md)" > "out-<worker>.log" 2>&1 &
   agy --model gemini-3.7-flash-high --effort high --dangerously-skip-permissions --print "$(cat brief.md)" > "out-<worker>.log" 2>&1 &
   tail -n 80 -f out-<worker>.log &
   wait
   ```

   Never cut a worker with a short synchronous `timeout` (do not use `timeout: 600` to kill `xhigh` thinking). Stream to `out-<worker>.log` and `tail` by default when the harness allows it (all current harnesses — `pi`, `agy`, `codex`, `opencode`, `cline` — do), so you can answer `"What are you doing?"` / `"Why are you working?"` from `todo.md`/`docs/plans/`/`docs/questions/` *plus* live `tail` while workers run. If a harness truly does not support streaming, document it in the hand-off brief and fall back to synchronous wait. After `wait`, captures stdout/stderr log path, exit code and checks `limit/quota/exhausted` → fallback.

4. **Verify before integrating** — orchestrator diffs files against acceptance criteria, runs `runx`/`xdocs` checks, commits via `guiho-s-0032-git-commit`. On failure, retries once with refined brief, then records to `docs/questions/<plan>/YYYY-MM-DD-<topic>.md` and `docs/issues/` and continues.

5. **Review and handoff** — after execution, a `mastermind` worker reviews (`docs/reviews/implementation/`, `docs/validation/`). Only then Phase 3 human review.

**Best strategy for us:** Single orchestrator + file-based hand-offs, **parallel by default**. If the plan has independent units that touch different files/dirs, do them in parallel — up to **10 concurrent workers** (or more if safe). No parallel pi subagents, no hidden state — all state is durable in `docs/`, `todo.md`, `MANDUME.md`. The orchestrator is the only agent you talk to; workers are stateless executors you never talk to directly.

### Scalability — When the project is large and you need many agents in parallel

**Yes, it scales — still with one orchestrator.** The orchestrator does not need native pi subagents; parallelism is achieved by running **multiple external CLI workers in parallel as subagents**, all managed by the single instance.

**Branching — work on `main`:** Do not use different branches for parallel work — work can be forgotten. All parallel work must start and finish **visible on `main`**, ready for review. Review itself stays **sequential** after parallel execution. Worktrees can isolate file sets but are not needed to satisfy the `main`-branch requirement; if you use worktrees, merge to `main` immediately after each worker finishes and before review — otherwise prefer direct `main` work with file-owned units.

**How to achieve it:**

1. **Partition the plan into independent units** — `docs/plans/<plan>/plan.md` units that touch different files/dirs (frontend vs backend, independent features). The orchestrator groups them by `todo.md` and file ownership so workers never write the same file.

2. **Dispatch a batch in parallel (always if possible) with live streaming by default** — for each independent unit, write its brief (`docs/plans/<plan>/execution/handoffs/<date>-<unit>-<worker>.md`), then launch workers concurrently (up to **10**) **streaming to per-worker logs by default** (unless the harness does not allow it):

   ```bash
   # from orchestrator — launch N workers in background, each streaming to its own log
   agy --model claude-opus-4-6-thinking --dangerously-skip-permissions --print "$(cat brief-engineer-2.md)" > out-2.log 2>&1 &
   agy --model gemini-3.7-flash-high --effort high --dangerously-skip-permissions --print "$(cat brief-engineer-3.md)" > out-3.log 2>&1 &
   pi --model zai/glm-5.3-flash:max --tools read,bash,edit,write,grep,find,ls -p --no-session "$(cat brief-engineer.md)" > out-1.log 2>&1 &
   tail -n 80 -f out-2.log &  # stream live while waiting — proves you're not frozen
   tail -n 80 -f out-3.log &
   tail -n 80 -f out-1.log &
   wait  # orchestrator waits for all, stays responsive to "What are you doing?" via docs state + tails
   ```

   Use `Usage` to spread load and respect limits (see below). Any harness the user registers in `MANDUME.md` works the same — `pi`, `agy`, `codex`, `opencode`, `cline`, etc. — the instructions are identical: **full permission** (`--tools read,bash,edit,write,grep,find,ls` for `pi`, `--dangerously-skip-permissions` for `agy`/others).

3. **Usage order — spend the 5-hour limits first:** Always start with **`5h-limited`** workers — you already paid for them, money is lost if not used 100%. When a `5h-limited` worker returns `limit/quota/exhausted` (check `stderr` for `limit`/`quota`/`exhausted`/`weekly`/`five hour`), fallback within same class: first to **`free`** workers, then lastly to **`paid-api-rate`** (credit-card, pay-as-you-go — money is saved if not used). While on `paid-api-rate`, **casually retry the `5h-limited` workers** periodically to see if their windows reset, and switch back immediately when they are available again.

4. **Verify and merge sequentially on `main`** — even if workers ran in parallel, the orchestrator **verifies and commits one by one on `main`** (diff vs acceptance criteria, `runx`/`xdocs`, `git add -- <paths>` → `git diff --cached` → `commit` via `guiho-s-0032-git-commit`). If a worker hit `limit/quota/exhausted`, fallback as above and retry once. No branch merges to forget — every parallel result ends visible on `main` before review.

5. **Stay responsive** — while workers run, the orchestrator can still answer _"Why are you working?"_ from `todo.md`/`docs/plans/`/`docs/questions/` because all state is file-based, not in worker memory.

**Result:** One instance controls up to 10 parallel executors without pi subagents, always on `main`. For very large projects, split by domain (`frontend/`, `backend/`) or by `todo.md` groups, dispatch each group to a `workhorse` worker, keep `mastermind` (`cto`/`reviewer`) sequential for judgment. This is how Mandume scales — orchestrator is the brain, CLI workers are the hands.
