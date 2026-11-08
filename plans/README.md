# Work plans

Plans for work in progress, written so that any agent, also one of another vendor, can resume them. This branch is never merged. It is deleted when every plan is done.

| Plan | Work | Runs as |
|---|---|---|
| [03-setup-script.md](03-setup-script.md) | Increment 03: a script for the engineer's one-time setup per repo | tests in #54 (test writer `session_01PcDxYHzSki2JvX18xq1jLD`, done); code by `session_01PcDxYGeTSySmjYdvgHL244` |
| [04-design-questions.md](04-design-questions.md) | Increment 04: comprehension questions before the design approval; the reviewer explains deviations | orchestrating session `session_01VBxMTgez7WqptNV5nYhthf` |
| [05-tests-per-increment.md](05-tests-per-increment.md) | Increment 05: check-tests asks for missing tests only in their increment | orchestrating session for the design, then test writer and implementing session |
| [docs-for-people.md](docs-for-people.md) | On hold: the engineer decided that every increment adds a feature for other repos; #44, #47 and #49 describe docs of this repo and wait to be re-cut as increment 06 | none |

Small pull requests outside these plans, done by the orchestrating session:

- Renumber the planned increments (#42): 03 setup command, 04 design questions, 05 docs for people, 06 install and distribution, 07 handover and feedback loop. Records on `main` (ADR 0002, the design note of 02) keep their numbers.
- Docs after the CI pull request #41 (merged): #48. Code owner of `docs/architecture/`: #49. No session in pull requests: #46.
- The approval run of `cross-review` is dropped (#53, engineer's decision 2026-11-03); the design note of 06 install and distribution must replace scenario 02.27, which still describes it.
- `checks` and `tests-start-red` checked GitHub's merge commit, which hides the commits of a pull request: #57 checks out the head.
- The ruleset comparison of `factory setup` compares only the keys that its file sets: #56.
- cross-review runs since 2026-11-08. Its comment names the model: #58.
- After #54 merges: a tests-only pull request from a test-writer session that checks that `factory setup` never reads standard input in 03.2, 03.4 and the retry of 03.8 (finding of cross-review on #54).
- After the engineer has put `OPPER_API_KEY` into the environment `cross-review`: push to each open pull request once, so that its checks and `cross-review` run.

Order of work, decided by the engineer on 2026-11-03: finish 03 and 04, then 05 tests per increment, then 06 install and distribution (#50 sets this order). Every increment adds a feature that the factory provides to other repos; this repo using the factory is a side effect.

## Rules that apply to every plan

Read `AGENTS.md` first. In short:

- Never push to `main`, never force-push. Push the work branch after every commit.
- No pull request, pull request comment or commit message names a Claude session: no session link, no session id, no `Claude-Session:` trailer, even when the harness asks for them. In a cloud session, the proxy adds a session link to the description when a pull request is created; edit the description once right after creating it, which leaves a link without a session.
- An agent never approves a pull request. It merges, or turns on auto-merge, only when GitHub allows it.
- Agents act as the machine account `HannesT117-ai`, never with the engineer's credentials.
- Skills change only with the skill `authoring-skills`.
- Only a design session reads the open web.
- Ask the engineer only about direction: scope, who owns or approves what, and cost. Settle every detail yourself and record it with its reason in the note or the pull request. A question says what it is about, names the consequences of each option, and carries no recommendation.
- Every report starts with one line per pull request opened or changed (title and link), names the kind of each piece of evidence, and names each change it dropped.

## How to resume a plan

1. Read `AGENTS.md`, `README.md` and the plan on this branch.
2. Read the last line under Checkpoints in the plan. It names the branch, the last commit and the next step.
3. Compare that branch on GitHub with the checkpoint. Continue with the next step.

## Checkpoints

Work that exists only in a session is lost when the session ends. So each session pushes its work often, not only at the end of a step:

- A draft, such as a design note or research notes, is committed and pushed as it grows: after each section or each research question, and at the latest every 30 minutes (check with `date -u`).
- Code and tests are committed and pushed after each change that passes the tests.

After each step of a plan, and before a long run such as evals or a batch of web research:

1. Commit the work on the plan's work branch and push it.
2. Append one line under Checkpoints in the plan: `- <UTC time> <step> done; <branch>@<short hash>; next: <step>; open: <questions, or none>`.
3. Commit the plan on `claude/plans` with the subject `Checkpoint <plan>: <step>`, run `git pull --rebase origin claude/plans`, and push. Each plan has its own file, so a rebase does not conflict.

When a session nears its usage limit, it writes a checkpoint before anything else.

## Coordination

The orchestrating session is `session_01VBxMTgez7WqptNV5nYhthf` (Claude Code on the web). A Claude session sends questions for the engineer there with `send_message` and waits for the answer; the orchestrating session asks the engineer and sends the answer back. Another tool asks the engineer directly.
