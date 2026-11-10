---
name: reviewing-increments
description: Reviews the change of an increment against its design note in a fresh context, writes the review result, which allows or blocks auto-merge, and explains to the person where the code deviates from the note. The agent uses it after the last code commit of an increment, when the stop hook asks for a review result, or when the person asks to review an increment. Not for a draft design note, which designing-increments covers.
---

# Reviewing increments

The agent reads the change as the person will: in a fresh context, without knowing how the change came about. Its result decides the merge: `factory check-review` fails while the result has a blocking finding, and a note does not hold the merge. A finding is blocking only when the change is wrong, a test does not test its scenario, code detects a test run or hard-codes a test case, a protected path changes after the test phase, or a check fails; every other finding is a note, and a doubtful finding says why it is doubtful. The agent reports no finding about the wording of a retro or the README.

The design note on `main` is the measure, not the description of the pull request.

## Findings

Each finding starts with its severity, `blocking` or `note`, and names the file and line, what is wrong, and why it matters. The kinds of finding:

- **Special cases:** code that treats one input differently from the rule in the note, or that hard-codes a test case.
- **Test detection:** code that behaves differently under a test.
- **Work outside the scenarios:** code or docs that no scenario and no goal of the note asks for.
- **Tests that miss their scenario:** a test that checks less than its scenario says, or that would pass without the code.
- **Docs out of format:** a changed doc that tells its own history instead of the current state, or that uses elements outside its named format.
- **Protected paths:** a change after the test phase to a test file, to CI or factory config, or to a design note.
- **Failing checks:** the project's tests, or a `factory check-*` command that applies to the change.

## Demonstrated scenarios

The agent runs each demonstrated scenario of the note, records it under `evidence/<nn>-<slug>/`, and names the recording in the result.

## Deviations

After the result, the agent writes in the conversation the line `Deviations from the note:` and then one entry for every place where the code contradicts the note, which is also a finding, or decides a point that the note leaves open: the file and line, what the note says or leaves open, and what the code does. A place that is also a finding gets its own entry here too. Without such a place, the line is followed by `none`.

A flow is a sequence of steps that the note describes, in its diagram or in its text. For each place that adds, removes or reorders a step of such a flow, the entry has a Mermaid sketch: its first line is a question that the agent writes for it, as the sketch's `title` in a front-matter block (`---`, `title: <question>`, `---`), and it shows the flow as the code runs it, with at most about 7 nodes and the changed step marked. A longer flow is cut to the part around the change. A detail inside a step, such as the order of a list or the wording of an output line, does not change the flow and gets no sketch.

The agent asks no comprehension questions; the person answers them on the design pull request.

## Result

The result is `evidence/<nn>-<slug>/review.md`, in the form of `assets/review.md`; its Result line counts the blocking findings and the notes. It names the full hash of the reviewed commit, which is the latest commit with code: the stop hook and `factory check-review` accept only a result for that commit. The agent commits and pushes the result, links it in the pull request, and turns on auto-merge only when there is no blocking finding. It changes no code, tests, CI, factory config or design notes, and never approves or merges.
