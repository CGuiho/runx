#### &copy; 2026 [GUIHO](https://guiho.co) as represented by [Cristóvão GUIHO](https://guiho.co/cguiho) All Rights Reserved.

# Mandume

> **All agents run by default with full permission** — full workspace read/write/execute, no sandbox or approval prompts. Every worker `Command` is executed with the same permissions as the calling agent.

## Mode

```yaml
execution: dnd  # dnd | interruptible — orchestrator NEVER stops during execution/review
notifications: off  # off | on
```

> **DND / Night Mode is the default during execution and review.** The user is only present **before execution** (planning) and **after review** (human review). During execution and review the user may be sleeping and **will not send another message**. The orchestrator writes `execution: dnd` here on `main`, commits it via `guiho-s-0032-git-commit`, and **never stops, never asks, never waits for a wake-up**. Every hand-off brief includes `Mode: dnd` so even stateless workers know. See `skills/guiho-s-0446-dnd/SKILL.md` and `skills/guiho-s-0440-hand-off/SKILL.md` (Never-Stop Contract — tails like `tail -f`).

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

### Never-Stop / Tail Contract — The Orchestrator Never, Never, Never Stops

> **The user is only there before execution and after review. During execution and review the orchestrator never, never, never stops — it tails like `tail -f`.**

**Mandume execution is unattended.** The user launches execution (often at night, "do not stop until everything is done") and goes to sleep. They will not send another message until Phase 3 human review. The work must go **from the first plan unit until the last, through execution and self-review, without stopping, without stopping, without stopping**. Handing a job to a worker does not mean stopping and waiting for the human.

**The orchestrator is `tail -f` on the worker:**

1. **Delegate** — write brief to `docs/plans/<plan>/execution/handoffs/<date>-<unit>-<worker>.md`, launch worker with full permission streaming to `out-<worker>.log` (`> out-<worker>.log 2>&1 &`).
2. **Tail** — **immediately** `tail -n 80 -f out-<worker>.log &` and stay alive. Never exit after `&`. For parallel batches, tail each log. This is mandatory whenever the harness allows streaming (all current harnesses — `pi`, `agy`, `codex`, `opencode`, `cline` — do). If a harness truly cannot stream, document it in the brief, but still never stop — wait, reap, and continue via file state.
3. **Reap, verify, continue** — on worker exit, capture exit code + log path + timing, check `limit/quota/exhausted` and fallback (`5h-limited` → `free` → `paid-api-rate`, retrying `5h-limited` casually), diff against acceptance criteria, run `runx`/`xdocs`, commit on `main` via `guiho-s-0032-git-commit`, then **immediately pick the next unit/batch**. No gap, no human prompt.
4. **Loop until done** — loop unit-by-unit or batch-by-batch (up to 10 parallel) until **every plan unit is verified and committed**, then enter self-review. The **only** permitted halt is a fail-closed hard stop under `docs/issues/` (safety/data-loss, impossible spec conflict, missing security authorization) — everything else is ledger-and-continue.

**If you stop after handing off, you failed.** The human expected morning completion. You handle everything yourself — delegation, tailing, fallback, retry, ledger, commit, next unit — without the user. See `skills/guiho-s-0440-hand-off/SKILL.md` (Never-Stop Contract) and `skills/guiho-s-0446-dnd/SKILL.md`.

### Calling Workers -- Subagents strategy

**There are no native pi subagents. There is one single pi instance — the orchestrator — that manages every task.**

The orchestrator (this `pi` session) never spawns internal subagents. It delegates to **external CLI workers as subagents** via `guiho-s-0440-hand-off`. Each worker is a separate process (`pi` or `agy` with its model/thinking) that receives a self-contained brief, runs with full permission, and returns. The orchestrator stays alive to answer you and to manage the lifecycle — **by tailing, never by stopping**.

**How it works:**

1. **You ask / orchestrator plans** — `cto`/`senior` (mastermind) reads `docs/specifications/`, `docs/architecture/`, `MANDUME.md`, writes `docs/plans/<plan>/plan.md` and `todo.md`. You can ask at any time _"Why are you working?"_ — the orchestrator answers from `docs/plans/`, `todo.md`, `docs/questions/` without interrupting workers.

2. **Orchestrator delegates one unit at a time** — picks a worker by `Class` (`mastermind` for judgment, `workhorse` for labor) and `Usage` (`paid-api-rate` primary, `5h-limited`/`free` fallback to spend subscriptions). Writes a brief to `docs/plans/<plan>/execution/handoffs/<date>-<unit>-<worker>.md` with Goal, Context, Skills to load, Constraints, Output contract.

3. **Call the worker and tail it like `tail -f` — never stop after launch (mandatory live streaming):**

   ```bash
   # pi mastermind/workhorse — stream to file so orchestrator can tail while waiting, never stopping
   pi --model <provider>/<model>:<thinking> --tools read,bash,edit,write,grep,find,ls -p --no-session "$(cat brief.md)" > "out-<worker>.log" 2>&1 &
   tail -n 80 -f out-<worker>.log &  # orchestrator is tail -f — stays alive, proves not frozen
   wait  # still waits for exit, but you can answer "What are you doing?" from tail; then reap and IMMEDIATELY continue
   # agy workhorse (Opus 4.6 / Gemini 3.7 Flash High) — same tailing pattern, never stop
   agy --model claude-opus-4-6-thinking --dangerously-skip-permissions --print "$(cat brief.md)" > "out-<worker>.log" 2>&1 &
   tail -n 80 -f out-<worker>.log &
   wait  # reap, verify, commit, then loop to next unit without waiting for human
   agy --model gemini-3.7-flash-high --effort high --dangerously-skip-permissions --print "$(cat brief.md)" > "out-<worker>.log" 2>&1 &
   tail -n 80 -f out-<worker>.log &
   wait
   ```

   Never cut a worker with a short synchronous `timeout` (do not use `timeout: 600` to kill `xhigh` thinking). Stream to `out-<worker>.log` and `tail -n 80 -f` by default when the harness allows it (all current harnesses — `pi`, `agy`, `codex`, `opencode`, `cline` — do), so you can answer `"What are you doing?"` / `"Why are you working?"` from `todo.md`/`docs/plans/`/`docs/questions/` *plus* live `tail` while workers run. **After `wait`, the orchestrator never stops — it captures stdout/stderr log path, exit code, checks `limit/quota/exhausted` → fallback, verifies, commits, and immediately loops to the next unit/batch.** If a harness truly does not support streaming, document it in the hand-off brief and fall back to synchronous wait, but still tail via file state and never stop.

4. **Verify before integrating — then immediately continue** — orchestrator diffs files against acceptance criteria, runs `runx`/`xdocs` checks, commits via `guiho-s-0032-git-commit`. On failure, retries once with refined brief, then records to `docs/questions/<plan>/YYYY-MM-DD-<topic>.md` and `docs/issues/` and **continues to next unit — never stops**. The loop only ends when all units are done.

5. **Review and handoff — without stopping between phases** — after execution, a `mastermind` worker reviews (`docs/reviews/implementation/`, `docs/validation/`). The orchestrator tails that worker too, then hands to Phase 3 human review. No human message is needed between execution and review.

**Best strategy for us:** Single orchestrator + file-based hand-offs, **parallel by default, never-stop loop**. If the plan has independent units that touch different files/dirs, do them in parallel — up to **10 concurrent workers** (or more if safe). No parallel pi subagents, no hidden state — all state is durable in `docs/`, `todo.md`, `MANDUME.md`. The orchestrator is the only agent you talk to; workers are stateless executors you never talk to directly. **The orchestrator tails every worker and loops until done — like `tail -f` from first unit to last.**

### Scalability — When the project is large and you need many agents in parallel

**Yes, it scales — still with one orchestrator that never stops.** The orchestrator does not need native pi subagents; parallelism is achieved by running **multiple external CLI workers in parallel as subagents**, all managed by the single instance that **tails each one**.

**Branching — work on `main`:** Do not use different branches for parallel work — work can be forgotten. All parallel work must start and finish **visible on `main`**, ready for review. Review itself stays **sequential** after parallel execution. Worktrees can isolate file sets but are not needed to satisfy the `main`-branch requirement; if you use worktrees, merge to `main` immediately after each worker finishes and before review — otherwise prefer direct `main` work with file-owned units.

**How to achieve it:**

1. **Partition the plan into independent units** — `docs/plans/<plan>/plan.md` units that touch different files/dirs (frontend vs backend, independent features). The orchestrator groups them by `todo.md` and file ownership so workers never write the same file.

2. **Dispatch a batch in parallel (always if possible) with live streaming and tailing by default — never stop:** For each independent unit, write its brief (`docs/plans/<plan>/execution/handoffs/<date>-<unit>-<worker>.md`), then launch workers concurrently (up to **10**) **streaming to per-worker logs and tailing each log by default** (unless the harness does not allow it):

   ```bash
   # from orchestrator — launch N workers in background, each streaming to its own log, tail each one, never stop
   agy --model claude-opus-4-6-thinking --dangerously-skip-permissions --print "$(cat brief-engineer-2.md)" > out-2.log 2>&1 &
   agy --model gemini-3.7-flash-high --effort high --dangerously-skip-permissions --print "$(cat brief-engineer-3.md)" > out-3.log 2>&1 &
   pi --model zai/glm-5.3-flash:max --tools read,bash,edit,write,grep,find,ls -p --no-session "$(cat brief-engineer.md)" > out-1.log 2>&1 &
   tail -n 80 -f out-2.log &  # tail each — proves you're not frozen, you are tail -f
   tail -n 80 -f out-3.log &
   tail -n 80 -f out-1.log &
   wait  # orchestrator waits for all, stays responsive to "What are you doing?" via docs state + tails, then reaps all and IMMEDIATELY verifies and loops — never stops
   ```

   Use `Usage` to spread load and respect limits (see below). Any harness the user registers in `MANDUME.md` works the same — `pi`, `agy`, `codex`, `opencode`, `cline`, etc. — the instructions are identical: **full permission** (`--tools read,bash,edit,write,grep,find,ls` for `pi`, `--dangerously-skip-permissions` for `agy`/others).

3. **Usage order — spend the 5-hour limits first, orchestrator handles fallback itself:** Always start with **`5h-limited`** workers — you already paid for them, money is lost if not used 100%. When a `5h-limited` worker returns `limit/quota/exhausted` (check `stderr` for `limit`/`quota`/`exhausted`/`weekly`/`five hour`), the orchestrator itself falls back within same class without asking the human: first to **`free`** workers, then lastly to **`paid-api-rate`** (credit-card, pay-as-you-go — money is saved if not used). While on `paid-api-rate`, **casually retry the `5h-limited` workers** periodically to see if their windows reset, and switch back immediately when they are available again. **Never stop to ask which worker to use — you decide and continue.**

4. **Verify and merge sequentially on `main` — then loop, never stop:** Even if workers ran in parallel, the orchestrator **verifies and commits one by one on `main`** (diff vs acceptance criteria, `runx`/`xdocs`, `git add -- <paths>` → `git diff --cached` → `commit` via `guiho-s-0032-git-commit`). If a worker hit `limit/quota/exhausted`, fallback as above and retry once. No branch merges to forget — every parallel result ends visible on `main` before review. **Then immediately dispatch the next batch — no waiting for human.**

5. **Stay responsive while tailing — never frozen:** While workers run, the orchestrator can still answer _"Why are you working?"_ from `todo.md`/`docs/plans/`/`docs/questions/` **plus** `tail -n 80 out-<worker>.log` because all state is file-based, not in worker memory. The `tail -f` is the proof the orchestrator never stopped.

**Result:** One instance controls up to 10 parallel executors without pi subagents, always on `main`, **tailing every one like `tail -f` and looping without stopping until done**. For very large projects, split by domain (`frontend/`, `backend/`) or by `todo.md` groups, dispatch each group to a `workhorse` worker, keep `mastermind` (`cto`/`reviewer`) sequential for judgment. This is how Mandume scales — orchestrator is the brain that never stops, CLI workers are the hands it tails.
