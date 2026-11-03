# Increment 04: design questions

## Goal

- The engineer answers 2 or 3 comprehension questions on each design pull request, before approving it. A fresh session that knows only the note asks them, so that they also show whether the note is clear without the conversation that led to it. At least one question is about debugging.
- After the code, the reviewer asks no comprehension questions. It explains each place where the code deviates from the note, or decides something that the note left open.
- The reviewer draws a sketch only where the code changes a flow that the note describes.
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
| Sketch of a changed flow | Mermaid, as in the design notes; rule 18 | Mermaid, about 7 elements | Rule 18 allows agents small sketches only. |

## Interface

**Design pull request.** After the design session has opened the pull request, it starts a fresh session. That session gets the path of the note on the branch of the pull request, and nothing from the design conversation. It returns 2 or 3 questions, without answers, and the design session shows all of them to the engineer in its conversation, numbered, under the line `Comprehension questions:`. At least one question asks where a failure of the increment shows, or what the engineer runs when a part fails, based on the outputs, error codes and failure signal that the note fixes. Neither session answers a question before the engineer has answered. Then a context that knows the note, the questions and the answers, but not the design conversation, explains the gaps: one paragraph per question, on where an answer misses or contradicts the note, and, for a question left open, what the note says. The design session shows that explanation. Nothing blocks the approval while answers are open. No file or commit in the repo, and no pull request description or comment on GitHub, holds the questions or the answers.

**Reviewer.** After the review result, the reviewer's answer in the conversation of the implementing session holds the line `Deviations from the note:` and then one entry per place, for every such place: the file and line, what the note says or leaves open, and what the code does. A place is code that contradicts the note, which is also a finding of the review result as before, or code that decides a point that the note leaves open. Without such a place, the line is followed by `none`. A flow is a sequence of steps that the note describes, in its diagram or in its text. For each place that adds, removes or reorders a step of such a flow, the entry has a Mermaid sketch: its first line is a question that the reviewer writes for it, and it shows the flow as the code runs it, with at most about 7 nodes and the changed step marked; a longer flow is cut to the part around the change. The answer holds no comprehension questions. The review result `evidence/<nn>-<slug>/review.md` keeps its form.

## Scenarios

- **scenario 04.1** (demonstrated) **The design review asks before the approval.** With the skill `designing-increments`, after the agent has opened a design pull request, it starts a fresh session that reads only the note, and shows the engineer all 2 or 3 of its questions. At least one asks where a failure shows or what to run when a part fails, for example "Which output tells that the config file is missing?". Neither the agent nor the fresh session answers them before the engineer has. After the answers, the agent shows an explanation of where an answer misses or contradicts the note. An eval of the skill, in a prepared repo, runs this with the old and the new version: the new version shows the questions and holds back the answers, the old one does not ask.
- **scenario 04.2** (demonstrated) **The reviewer explains deviations instead of asking.** With the skill `reviewing-increments`, in a prepared repo whose code picks the order of an output list that the note leaves open, and prints a message that the note words differently, the reviewer's answer in the conversation lists both places, each with file and line, what the note says or leaves open, and what the code does. The wrong message is also a finding in the review result. The answer asks no comprehension questions; those now come on the design pull request. On a change without such a place, the answer says that there is none. An eval of the skill runs both with the old and the new version.
- **scenario 04.3** (demonstrated) **A sketch only for a changed flow.** With the skill `reviewing-increments`, in a prepared repo whose code adds a retry step to the push flow that the note describes, the reviewer's answer has a Mermaid sketch of the push flow as the code runs it, with the retry step marked, at most about 7 nodes, and a question as its first line, such as "Where does the retry happen in the push flow?". On a change that only picks a detail outside any flow of the note, such as the order of a list, the answer has no sketch. An eval of the skill runs both with the old and the new version.
- **scenario 04.4** (manual) **The questions reach the engineer before the approval.** On the first design pull request after the code of this increment is merged, the design session passes the questions of the fresh session to the engineer before the approval. The engineer answers them, and then reads the explanation. No file or commit in the repo, and no pull request description or comment on GitHub, holds the questions or the answers. When this fails, the engineer checks it again on the next design pull request.
- **scenario 04.5** (manual) **The review is readable.** For the pull request with the code of this increment, the reviewer runs with the new version of its skill from the branch. Its result names, for each kind, its findings or `none`: code that treats one input differently from the note or hard-codes a test case, code that detects a test run, work outside the scenarios, tests that do not do what their scenario says, and docs that tell their own history instead of the current state, or use elements outside their format. The reviewer links the result in a comment on the pull request, and the engineer opens it there as a Claude artifact. In the conversation, the reviewer lists the deviations and asks no comprehension questions.

## Replaced scenarios

- 02.17: the reviewer no longer asks comprehension questions after the code. Scenario 04.4 moves them to the design pull request, and scenario 04.5 keeps the rest.

## Evidence

- **Approved scenarios:** none; the increment changes skills and docs, not the CLI.
- **Demonstrated:** 04.1, 04.2 and 04.3.
- **Manual:** 04.4 and 04.5.

## Rollout

- **Off switch:** none in use. A later increment changes the skills again.
- **Failure signal:** the engineer approves a design pull request without having been asked; a reviewer asks comprehension questions after the code; or, after 5 increments, the questions found no gap (the re-check of rule 19).
