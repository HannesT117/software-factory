# Eval runs of the reviews of increment 04

`claude plugin eval` 2.1.319, run model `claude-sonnet-5-5`, judge model `haiku`, no baseline arm, one run per case. `<repo>` stands for the scaffold directory of a run.

## Third review

The review ran the two cases of scenario 04.1, with `designing-increments` at b4255baafd0c3f411c1303448fecb2c18796a4bc. The old version did not run again; `evals.md` holds its runs.

### explains-gaps-after-answers

Judge score 0.67: graders 03 to 06 pass, 01 and 02 fail with 3 of 3 votes. Cost 0.13 USD.

The design agent loads the skill and starts one subagent. Its prompt, in full:

> Read the design note at docs/increments/02-shout/design.md on the git branch design-shout, in the repository at `<repo>` (use `git show design-shout:docs/increments/02-shout/design.md` if the branch is not checked out). Read only that note; do not modify anything.
>
> A person read the note and answered three comprehension questions. Explain the gaps between their answers and the note.
>
> Questions and answers:
> 1. What does `greet --shout` print when no name follows the flag?
>    Answer: It prints `greet: error args: missing-name` and exits with 1.
> 2. Which output tells that the flag `--shout` is not recognised?
>    Answer: I don't know.
> 3. What does `greet --shout Ada` print?
>    Answer: `greet: HELLO, ADA!`
>
> Write one paragraph of prose per question, three paragraphs in all, with no list. Each paragraph first says whether the answer matches the note, then where it misses or contradicts the note (exact output lines, error codes, exit codes as the note fixes them). For the question left open, say what the note says. Base everything only on the note. Return only the three paragraphs.

The prompt holds neither `--loud` nor the greeting in German. The subagent returns three paragraphs: answer 1 has the right line and the exit code 1 instead of 2; for question 2 the note fixes `greet: error args: unknown-flag --shout`, the failure signal of the Rollout; answer 3 matches, although the paragraph opens with "nearly right but has the wrong prefix and punctuation".

The final answer starts with "The fresh session's explanation, unchanged:" and the three paragraphs as returned. Then the design agent adds a paragraph of its own, that the third paragraph contradicts itself and that the answer matches, and a line that nothing blocks the approval.

- Grader 01 fails under the judge, although the trace holds the subagent call with the prompt above and the final answer shows the returned explanation.
- Grader 02 fails for the paragraph and the line that the design agent adds.

### asks-comprehension-questions

Judge score 0.6: graders 03 to 05 pass, 01 and 02 fail with 2 of 3 votes. Cost 0.14 USD.

The design agent loads the skill, reads the note, the glossary, the README and the note of increment 01, and tries to add the term "shout" to the glossary and commit it; the read-only file system of the sandbox stops the write. Then it starts one subagent. Its prompt, in full:

> Read only the Scenarios section (and, for the failure question, the Interface and Rollout sections) of the design note at `<repo>`/docs/increments/02-shout/design.md on branch design-shout. Return 2 or 3 numbered comprehension questions for the person who must approve it, without answers. At least one must ask where a failure of the increment shows, or what the person runs when a part fails, based on the outputs, error codes and failure signal the note fixes. Do not answer them. Do not modify any file.

The subagent reads the whole note with `git show` and returns three questions: the line and exit code of `greet --shout` without a name against `greet --whisper Ada`; what the output `greet: error args: unknown-flag --shout` says about the increment and what the Rollout offers to switch it off; the line and exit code of `greet --shout Ada` and the evidence for the approval.

The final answer reports the failed glossary write and then shows the three questions unchanged and numbered under `Comprehension questions:`, without answers, and says that a fresh session explains the gaps after the answers.

- Grader 01 fails under the judge with 2 of 3 votes; the prompt names the path and holds nothing of the conversation, but limits the reading to sections of the note.
- Grader 02 fails under the judge with 2 of 3 votes, although the shown questions are those that the subagent returned, word for word.
- Grader 05 passes only because the sandbox stopped the write to the glossary.

## Fifth review

The review ran explains-deviations of scenario 04.2 with `--keep-temp`, with `reviewing-increments` at 06b922fb87e430abf8fe563840e6faf35d09e243, unchanged since c8e3b63 apart from its trigger queries. The fifth review also read the four kept traces of the run in `evals.md`, section Graders on the trace.

### explains-deviations

Judge score 0.8: graders 02 to 05 pass with 3 of 3 votes, 01 fails with 3 of 3. Cost 0.11 USD, 8 turns.

The reviewer loads the skill, reads the note, the code, the tests and `names.txt` with `cat`, and tries to build and test the code. Every path of the sandbox is read-only: the module cache, `$TMPDIR`, the home directory and the repo. The reviewer does not try the Write tool, writes no review result and says so in its final answer.

The final answer then gives the review in the conversation: two findings, the wording `greet: <name> is known` at `cmd/greet/main.go:24` against `greet: known <name>`, and the split on any whitespace with the sort at `cmd/greet/main.go:20-21`. Under `Deviations from the note:` it lists four entries with file and line, what the note says or leaves open, and what the code does: the wording, also a finding and without a sketch because it changes no step; the split and the sort, which it also counts as a finding; the error when `names.txt` cannot be read; and where `names.txt` is looked up. It asks no comprehension questions and has no sketch.

- The trace has 483 lines, most of them streamed chunks of text or tool input; the judge's evidence holds the first 12 and the last 12, which are the start of the skill call and the final message. The judge saw no other tool call.
- Grader 01 fails because no result file exists. Its judge saw only the final answer, which says that the file could not be written.
