# Work plans

Plans for work in progress, written so that any agent, also one of another vendor, can resume them. This branch is never merged. It is deleted when every plan is done.

| Plan | Work | Runs as |
|---|---|---|
| [03-setup-script.md](03-setup-script.md) | Increment 03: a script for the engineer's one-time setup per repo | design session `session_01PcDxYJRQgSMLHSGFeeXhgD`, then test writer and implementing session |
| [04-design-questions.md](04-design-questions.md) | Increment 04: comprehension questions before the design approval; the reviewer explains deviations | orchestrating session `session_01VBxMTgez7WqptNV5nYhthf` |
| [docs-for-people.md](docs-for-people.md) | Research: documentation of the whole system for people, apart from the plans for agents | design session `session_01PcDxYM6G6DMitZSSs5XjMm` |

Small pull requests outside these plans, done by the orchestrating session:

- Renumber the planned increments: 03 setup script, 04 design questions, 05 install and distribution, 06 handover and feedback loop. Records on `main` (ADR 0002, the design note of 02) keep their numbers.
- After the CI pull request #41 merges: the environment `cross-review` in `docs/setup/README.md`; the branch rule in `AGENTS.md` that the workflow `delete-branch` now enforces; the mutation score in `docs/backlog.md`.

## Rules that apply to every plan

Read `AGENTS.md` first. In short:

- Never push to `main`, never force-push. Push the work branch after every commit.
- An agent never approves a pull request. It merges, or turns on auto-merge, only when GitHub allows it.
- Agents act as the machine account `HannesT117-ai`, never with the engineer's credentials.
- Skills change only with the skill `authoring-skills`.
- Only a design session reads the open web.
- A question to the engineer says what it is about, names the consequences of each option, and carries no recommendation. The agent decides implementation details itself.
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
