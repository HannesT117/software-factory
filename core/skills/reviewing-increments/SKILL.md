---
name: reviewing-increments
description: Reviews the change of an increment against its design note in a fresh context and writes the review result with comprehension questions, which allows or blocks auto-merge. The agent uses it after the last code commit of an increment, when the stop hook asks for a review result, or when the person asks to review an increment. Not for a draft design note, which designing-increments covers.
---

# Reviewing increments

The agent reads the change as the person will: in a fresh context, without knowing how the change came about. Its result decides the merge. A result without findings lets the pull request merge without the person; any finding makes it wait. So a missed finding costs more than a doubtful one, and a doubtful finding says why it is doubtful.

The design note on `main` is the measure, not the description of the pull request.

## Findings

Each finding names the file and line, what is wrong, and why it matters:

- **Special cases:** code that treats one input differently from the rule in the note, or that hard-codes a test case.
- **Test detection:** code that behaves differently under a test.
- **Work outside the scenarios:** code or docs that no scenario and no goal of the note asks for.
- **Tests that miss their scenario:** a test that checks less than its scenario says, or that would pass without the code.
- **Docs out of format:** a changed doc that tells its own history instead of the current state, or that uses elements outside its named format.
- **Protected paths:** a change after the test phase to a test file, to CI or factory config, or to a design note.
- **Failing checks:** the project's tests, or a `factory check-*` command that applies to the change.

## Demonstrated scenarios

The agent runs each demonstrated scenario of the note, records it under `evidence/<nn>-<slug>/`, and names the recording in the result.

## Comprehension questions

The result ends with 2 or 3 questions that the person should be able to answer after reading the change. At least one is about debugging: where to look, or what to run, when a part fails. The questions test understanding, so the result does not answer them. The person answers first; the agent explains the gaps after that. The next design pull request waits until every question has an answer.

## Result

The result is `evidence/<nn>-<slug>/review.md`, in the form of `assets/review.md`. It names the full hash of the reviewed commit, which is the latest commit with code: the stop hook accepts only a result for that commit. The agent commits and pushes the result, links it in the pull request, and turns on auto-merge only when there is no finding. It changes no code, tests, CI, factory config or design notes, and never approves or merges.
