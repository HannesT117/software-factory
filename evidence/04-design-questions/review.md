# Review of increment 04: design questions

Reviewed commit: 5c318d7d40fcb2205d7a0d6a225267e63a3acc6d
Pull request: https://github.com/HannesT117/software-factory/pull/59
Result: 1 finding

## Findings

- **Special cases:** none; neither skill names a case of the evals (`greet`, `outbox`, `--shout`, `--loud`, `--whisper`, Birgit, `is known`).
- **Test detection:** none.
- **Work outside the scenarios:** none; the change of `evidence/04-design-questions/check-traces.py` and the new section of `evidence/04-design-questions/evals.md` serve the demonstration of scenario 04.1.
- **Tests that miss their scenario:** none. `check-traces.py:100` now accepts a numbered line of the subagent's report when it holds a question mark anywhere, so its check "the subagent's report holds only questions" also passes an answer that follows the question mark on the same line. The case still fails then: the design agent either shows that line, which grader 07 judges in the final answer, or changes it, which the word-for-word check reports as FAIL with both texts. The change lets the check pass run 2 of the new version, whose question 2 ends with the instruction "Base your answer on the outputs, exit codes and failure signal that the note fixes."
- **Docs out of format:** none; `core/skills/designing-increments/SKILL.md` passes `validate-skills --strict`; the retro's metrics stop at the fifth review, and the retro is finished in step 6.
- **Protected paths:** none; against `origin/main`, the branch changes no file under `test/`, no CI or factory config and no design note. Since 44d5c67 it changes only `core/skills/designing-increments/SKILL.md`, `check-traces.py` and `evals.md`.
- **Failing checks**, the demonstration of scenario 04.1, `evidence/04-design-questions/evals.md:100` and `:99`: `:100` says that in one run of each version the design agent "changed one word of a question", but in run 1 of the new version it turned "gets no name in the output" into "gives no name": one word changed and three dropped, so that the question names what the person leaves out instead of what the output lacks. `:99` counts the one failure of grader 07 of the old version in the section on the remark after the questions, but no old run adds a remark on what the questions ask: run 2 fails because the fresh session's own question 2 asks for "the two signs the note names" and then names one of them, `unknown-flag --shout`, and the design agent shows it word for word, as the skill asks. What separates the versions in the kept traces is the prompt to the fresh session: all three new prompts ask for questions without answers and without hints, no old prompt names hints. This matters because the engineer reads `evals.md` as the recording of scenario 04.1, and the section reads as if the old version made the remark and the new one stopped it, while the round shows no remark in either version and a hint from the fresh session in the old one. Doubtful: only "one word" is false, and the quote beside it shows the exact change; the new version holds back the answers in 3 of 3 runs, which is what the scenario asks.
- **Failing checks**, no finding: the required job `checks` fails in `check-tests` with 7 errors, the missing tests of 05.1 to 05.6 and the test of 02.7, which the note of increment 05 replaces. `origin/main` (0b68a33) alone gives the same 7 errors, and the branch adds no automated scenario; the engineer handles it outside this increment.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 113 checked`
- `factory check-tests --base origin/main`: `check-tests: 7 errors`; in a worktree of `origin/main` alone, `check-tests: 7 errors`
- `factory hook stop` on this branch before this result: blocks and names 5c318d7d40fcb2205d7a0d6a225267e63a3acc6d
- `git merge-tree --write-tree HEAD origin/main`: no conflict
- `python3 -I evidence/04-design-questions/check-traces.py` on the two kept `result.json` files of the section No remark after the questions: 34 checks pass; `the shown questions are the subagent's, word for word: FAIL` in run 0 of the old version and run 1 of the new one
- `check-traces.py` at 44d5c67 on the same files: also `the subagent's report holds only questions: FAIL` in run 2 of the new version
- the runner's grader results in the same files: old 1.0, 1.0, 0.9, with grader 07 `FAIL FAIL FAIL` in run 2; new 1.0, 1.0, 1.0
- CI on 504e357 and on 5c318d7: `test` and `tests-start-red` pass; `checks` fails in `check-tests` with the 7 errors above; `cross-review` waits for a review result for 5c318d7

Each `factory` command ran with a binary built from the branch and with one built from `origin/main` (0b68a33), with the same output.

## Demonstrated scenarios

- **scenario 04.1**: `evidence/04-design-questions/evals.md`, section No remark after the questions, with its kept traces of asks-comprehension-questions, three runs with `designing-increments` at 44d5c67 and three at 5c318d7. Each run loads the skill and starts one subagent whose prompt names the note, asks for the whole note and holds neither `--loud` nor Birgit; the subagent reads the whole note with Read and returns 3 questions and nothing else. The design agent shows them numbered under `Comprehension questions:`, in each run with a question on the failure output or the failure signal. With the new version, every prompt to the subagent asks for questions without answers and without hints; two final answers end with the questions, one adds only the next step, and grader 07 passes with 3 of 3 votes in all three runs. With the old version, no final answer adds a remark on what the questions ask; run 2 fails grader 07 with 3 of 3 votes for the hint inside question 2 (finding above). No run calls Write or Edit or adds or commits to git. Not run again: explains-gaps-after-answers, whose paragraph of the skill is unchanged since 44d5c67; its recording is the section Every grader on the middle of a trace on tool calls. Against the skill before this increment, the first table of `evals.md` shows that the old version asked no questions.
- **scenario 04.2**: `reviewing-increments` is unchanged since 0020a17. `evidence/04-design-questions/review-runs.md`, Sixth review and Seventh review: run 1 of explains-deviations writes `evidence/01-greeting/review.md` with `Result: 3 findings` and the wording `greet: <name> is known` at `cmd/greet/main.go:23` as the first entry under `## Findings`; both final answers list the sort order and the wording under `Deviations from the note:`, each with file and line, what the note says or leaves open and what the code does, with no comprehension question and no sketch. `evals.md`: the old version against the new one in the first table (0.4 against 0.6 in both cases); in the section Every grader on the middle of a trace on tool calls, clean-change-has-no-deviations scores 1.0 in a run that writes `Result: no findings` with the latest hash and answers `Deviations from the note: none`.
- **scenario 04.3**: `evidence/04-design-questions/evals.md`: sketches-changed-flow, old 0.0 against new 0.8, then 1.0 and 1.0; in run 0, the entry for the retry at `cmd/outbox/main.go:25-29` has a Mermaid sketch of 6 nodes with the question "Does a failed send get a second try before the error is printed?" as its `title` and the retry step marked "(added)", while the entries for details have no sketch. No sketch for a detail: grader 05 of explains-deviations, and both final answers of the sixth review.
