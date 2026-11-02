# Working on the software factory

Read `README.md` first. It is the entry point and shows the status. The factory follows its own rules from `docs/principles.md`. This file holds only rules that no check enforces; `factory` and the hooks enforce the others and name the fix in their output.

## Workflow: every increment follows these steps

An increment is one unit of work in two pull requests, in the folder `docs/increments/<nn>-<slug>/`. The formats and the path groups are in the design note of increment 02, section Format.

1. **Design pull request.** Write `design.md` with the skill `designing-increments`. Turn on auto-merge. Stop. The engineer approves as code owner.
2. **Tests.** The agent starts a separate session as `acceptance-test-writer`. It commits one `testscript` file per automated scenario to a new work branch, and stops.
3. **Code, on the same branch.** The implementing agent writes the code until the tests pass. Stop and ask if a test looks wrong.
4. **Review and demonstration.** When the stop hook blocks, start the subagent `increment-reviewer`. It writes `evidence/<nn>-<slug>/review.md` with the skill `reviewing-increments`, and turns on auto-merge only when it found nothing. The engineer answers its comprehension questions in the conversation before the reviewer explains; neither questions nor answers go into the repo. The same session runs the demonstrated scenarios and puts their recordings into `evidence/`.
5. **Manual scenarios.** The engineer runs them.
6. **Retro.** Write `retro.md` in the same pull request: metrics, agent version, corrections, and every point that concerns the factory. From increment 04 on, `factory feedback` prints a prefilled issue link for each factory point, and the engineer submits it. Then a summary of at most 10 lines.

Verification kinds: **automated** is a `testscript` file from step 2; **demonstrated** is run and recorded by the review session, and the engineer views the recording; **manual** is run by the engineer.

## Sessions and subagents

The agent starts subagents and separate sessions on its own whenever a step needs them, such as the test writer, the reviewer, or a fresh context that explains scenarios. It passes prompts and results between sessions itself, and turns to the engineer only for approvals, decisions and answers to comprehension questions.

## Skills

Skills are created and changed only with the skill `authoring-skills`. Where it differs from `skill-creator` or another skill for skills, `authoring-skills` applies.

## The implementing agent never

- writes code that detects a test run, or that hard-codes test cases
- uses the engineer's own credentials, or reads `~/.ssh`, `~/.config/gh` or the SSH agent. Agents act as the machine account.
- approves a pull request. It may merge one, or turn on auto-merge, once GitHub allows the merge. It does not merge while the engineer has switched off the ruleset.
- counts its own unit tests as evidence. They are scaffolding.
- reads an issue in the factory repo that does not have the label `accepted`

## Reports

Every report starts with one line for each pull request that the agent opened or changed since the last report: its title and its link. Before building a change, the agent weighs its benefit against its cost. It drops a change with little benefit, also one that a review suggested, and names what it dropped in the report. Every report names the kind of each piece of evidence: approved scenario, demonstrated, manual, review, ADR confirmation, agent test (scaffolding), or not verified.

The agent asks the engineer directly, never through a section in a document. A question says what it is about, so that a reader without context understands it, and names the consequences of each option. The agent decides implementation details itself and does not ask about them.

## Pull requests and commits

Pull requests follow the skill `writing-pull-requests`, commits the skill `writing-commit-messages`.

## Source of truth and push

- The engineer's clone at `~/dev/software-factory` is the source of truth. Before each change, compare GitHub with the local clone. If the engineer works locally without pushing, ask for a push or a copy and build on that state.
- Push the work branch after every commit, fast-forward only, never with force. After the merge, the branch is deleted. A branch whose content moved to other branches is deleted too.

## Docs

- Docs state general rules. They contain no "you" and no "I". The roles are engineer and agent.
- Decision records use MADR, in `docs/decisions/<category>/`, with the categories `architecture` and `process`. A process decision gets a record only when it had real alternatives; its rule in `docs/principles.md` links to the record. Numbers count per category. Context and decision drivers are value-neutral facts. A record on `main` is in force, because the engineer merged it, and is never rewritten. A new record that replaces it, in the same category, adds `status: superseded by ADR-<nnnn>` to the old one, in the same pull request. A record that stops applying without a replacement gets `status: deprecated` and a one-line reason under More Information.
- Docs describe the current state. How a doc came about stays in git and in the pull request. Research notes keep rejected options with their reason, not the course of the research.
- A doc points to a config file, such as a ruleset or `CODEOWNERS`, instead of restating its content, so that it does not go stale when the file changes.
- A doc that follows a named format (MADR, the Agent Skills specification, the design-note format) uses only the elements of that format. An addition names its source in the pull request.
- Process and policy are rules in `docs/principles.md`. Each rule names its source, what it depends on, how it is enforced, and when to re-check it.
- `README.md` stays one page. Update its status in every increment.

## Code

- The CLI is Go. It uses only the standard library at runtime. Test-only dependencies are allowed.
- `core/` and the tool-neutral packages name no specific agent tool. Tool-specific parts live in `adapters/<tool>/`.
- Output is short and grep-able. Each line starts with the command's prefix, for example `verify:`.
