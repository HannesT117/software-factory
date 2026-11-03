# Plan: increment 04, comprehension questions before the design approval

- **Folder:** `docs/increments/04-design-questions/`
- **Target pull requests:**
  1. Design pull request with `design.md`, from the branch `claude/04-design-questions-design`.
  2. Code pull request from `claude/04-design-questions`, after the design note is on `main`. The changes are skills and docs, verified by evals and by the engineer. If the note has no automated scenario, no test files come first.
- **Runs as:** the orchestrating session. No step needs the open web.

## Context

- Today the reviewer asks the engineer 2 or 3 comprehension questions after the code, in the conversation. At least one is about debugging. The engineer answers first, and the reviewer explains the gaps after that. Neither questions nor answers go into the repo. Sources:
  - `AGENTS.md`, step 4
  - `docs/principles.md`, rule 19 (source: Anthropic RCT, largest gap in debugging)
  - `core/skills/reviewing-increments/SKILL.md`, section Comprehension questions, and its eval `asks-questions-outside-the-repo`
  - `adapters/claude-code/agents/increment-reviewer.md`
- `core/skills/designing-increments/SKILL.md` already has a fresh context explain each scenario before the design pull request.
- The reviewer measures against the design note on `main`, and the stop hook accepts a review result only in `evidence/<nn>-<slug>/review.md`. That is why this change is an increment: a decision of the agent.

## Decisions of the engineer, 2026-11-03

- **Design pull request:** after it is open, a fresh session reads only the design note and asks the engineer 2 or 3 comprehension questions, before the approval. The design session starts that session and passes the questions and answers between it and the engineer.
- **Debugging:** at least one question is about debugging. It rests on the outputs and failure signals that the note fixes, in Interface and Rollout. The engineer answers first; then the session explains the gaps. Questions and answers stay out of the repo.
- **After the code:** the reviewer asks no comprehension questions. In the conversation, it lists and explains each place where the code deviates from the note, or decides something that the note left open.
- **Diagram after the code:** a sketch of about 7 elements that answers one question, only where the built flow deviates from the diagram of the note.
- The engineer added that one diagram per design note may be too little. That question belongs to the plan [docs-for-people.md](docs-for-people.md), not to this increment.

## Changes in the code pull request

- `core/skills/designing-increments/SKILL.md`: the comprehension questions after the design pull request is open; an eval for them.
- `core/skills/reviewing-increments/SKILL.md`: the section on comprehension questions becomes a section on deviations, with the sketch on deviation. The eval `asks-questions-outside-the-repo` becomes an eval where the code decides something that the note left open.
- `adapters/claude-code/agents/increment-reviewer.md`: the deviations and the sketch go into Claude's answer.
- `AGENTS.md`: step 1 gets the questions; step 4 loses them and gets the deviations.
- `docs/principles.md`, rule 19: the new timing; enforced by the two skills.
- `README.md`: the status.

## Steps

1. Write `design.md` with the skill `designing-increments`: Goal; Out of scope (docs for people, see the other plan); Existing tools; Interface (what the engineer sees in the conversation, and when); Scenarios, demonstrated through the evals of both skills, and manual for the engineer's experience; Rollout. Run `factory check-design docs/increments/04-design-questions`, have a fresh context explain each scenario, open the design pull request, turn on auto-merge, and stop.
2. After the merge: create the branch. Load `authoring-skills` and read its `references/testing.md` completely. Change both skills and their evals. Run `factory validate-skills --strict` and `go test ./...`.
3. Run the evals, old version against new version, all cases started together. Show the engineer the review viewer before judging the outputs. A checkpoint comes before the run, because the run uses much of the usage limit.
4. Change `AGENTS.md`, rule 19, the agent prompt and the README. One topic per commit, and push after each.
5. When the stop hook blocks: the reviewer writes `evidence/04-design-questions/review.md`. Then the retro.

## Checkpoints
- 2026-11-03T00:54Z step 1 draft written, check-design ok; claude/04-design-questions-design@f086344; next: fresh context explains the scenarios, then the design pull request; open: none
- 2026-11-03T00:56Z step 1 done: design pull request https://github.com/HannesT117/software-factory/pull/43 open, scenarios sharpened after the fresh-context check; claude/04-design-questions-design@1353a3e; next: step 2 after the engineer approves and #43 merges; open: none
- 2026-11-03T03:33Z step 2 done: both skills, the reviewer prompt and 5 eval cases changed or added; claude/04-design-questions@5deb95d; next: step 3, evals old against new with `claude plugin eval` (1 run per case and version), then the review viewer for the engineer; open: none
