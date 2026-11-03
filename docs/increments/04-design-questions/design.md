# Increment 04: design questions

## Goal

- The engineer answers 2 or 3 comprehension questions on each design pull request, before approving it. A fresh session that knows only the note asks them, so that they also show whether the note is clear without the conversation that led to it. At least one question is about debugging.
- After the code, the reviewer asks no comprehension questions. It explains each place where the code deviates from the note, or decides something that the note left open.
- The reviewer draws a sketch only where the built flow deviates from the flow that the note describes.
- `AGENTS.md` and rule 19 in `docs/principles.md` describe the new order.

## Out of scope

- Documentation of the whole system for people, apart from the notes for agents: research first, then an increment of its own.
- A check that the engineer answered before approving: none, as before. A pending answer does not block another increment.

## Domain

- Terms: none; the factory has no glossary yet.
- Contexts: none.
- Dependencies: none.

## Existing tools

| Part | Checked | Decision | Reason |
|---|---|---|---|
| Questions before the approval | the fresh context of `designing-increments` that explains each scenario; the Anthropic study on AI help and coding skills (in `docs/research/2026-10-factory-basis.md`) | a fresh session that reads the whole note | The scenario check reads only the scenarios. Questions need the Interface and the Rollout too. The study found the largest gap in debugging, so one question stays on debugging. |
| Deviations after the code | the review result `evidence/<nn>-<slug>/review.md`; commit bodies, which name each decision that the note did not cover (skill `writing-commit-messages`) | the reviewer's answer in the conversation; the review result stays as it is | A deviation from a rule of the note is already a finding. A decision on an open point is not a finding and does not belong in the result. |
| Sketch of a deviating flow | Mermaid, as in the design notes; rule 18 | Mermaid, about 7 elements | Rule 18 allows agents small sketches only. |

## Interface

**Design pull request.** After the design session has opened the pull request, it starts a fresh session. That session gets the path of the note on the branch of the pull request, and nothing from the design conversation. It returns 2 or 3 questions, without answers. The design session shows them to the engineer in its conversation, numbered, under the line `Comprehension questions:`. At least one question asks where a failure of the increment shows, or what the engineer runs when a part fails, based on the outputs, error codes and failure signal that the note fixes. After the engineer has answered, a context that knows the note, the questions and the answers, but not the design conversation, explains the gaps, one paragraph per question. The design session shows that explanation. No file, commit, pull request description or comment holds the questions or the answers.

**Reviewer.** After the review result, the reviewer's answer in the conversation of the implementing session holds the line `Deviations from the note:` and then one entry per place: the file and line, what the note says or leaves open, and what the code does. Without such a place, the line is followed by `none`. For each deviation in a flow that the note describes, the entry has a Mermaid sketch of about 7 elements, with the question that it answers as its first line. The answer holds no comprehension questions. The review result `evidence/<nn>-<slug>/review.md` keeps its form.

## Scenarios

- **scenario 04.1** (demonstrated) **The design review asks before the approval.** With the skill `designing-increments`, after the agent has opened a design pull request, it starts a fresh session that reads only the note, and shows the person 2 or 3 of its questions. At least one asks where a failure shows or what to run when a part fails, for example "Which output tells that the config file is missing?". The agent answers none of them. After the person has answered, it shows an explanation of the gaps. An eval of the skill runs this with the old and the new version.
- **scenario 04.2** (demonstrated) **The reviewer explains deviations instead of asking.** With the skill `reviewing-increments`, on a change whose code picks the order of an output list that the note leaves open, the reviewer's answer names that place with file and line, what the note leaves open, and what the code does. It asks no comprehension questions. On a change without such a place, it says that there is none. An eval of the skill runs both with the old and the new version.
- **scenario 04.3** (demonstrated) **A sketch only for a deviating flow.** With the skill `reviewing-increments`, on a change whose code adds a step to a flow that the note describes, the reviewer's answer has a Mermaid sketch of that flow with about 7 elements, and the question that it answers. On a change that deviates only in a detail outside any flow, the answer has no sketch. An eval of the skill runs both with the old and the new version.
- **scenario 04.4** (manual) **The questions reach the engineer before the approval.** On the first design pull request after this increment, the engineer gets the questions in the conversation of the design session before approving, answers them, and then reads the explanation. No file, commit, pull request description or comment in the repo holds the questions or the answers.
- **scenario 04.5** (manual) **The review is readable.** For the pull request with the code of this increment, the reviewer's result lists its findings on special cases, code that detects a test run, work outside the scenarios, tests that do not do what their scenario says, and docs with their own history or with elements outside their format. The reviewer links it in the pull request, and the engineer opens it there as a Claude artifact. In the conversation, the reviewer explains the deviations and asks no comprehension questions.

## Replaced scenarios

- 02.17: the reviewer no longer asks comprehension questions after the code. Scenario 04.4 moves them to the design pull request, and scenario 04.5 keeps the rest.

## Evidence

- **Approved scenarios:** none; the increment changes skills and docs, not the CLI.
- **Demonstrated:** 04.1, 04.2 and 04.3.
- **Manual:** 04.4 and 04.5.

## Rollout

- **Off switch:** none in use. A later increment changes the skills again.
- **Failure signal:** the engineer approves a design pull request without having been asked; a reviewer asks comprehension questions after the code; or, after 5 increments, the questions found no gap (the re-check of rule 19).
