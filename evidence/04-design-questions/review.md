# Review of increment 04: design questions

Reviewed commit: 01d1869fbb816b32297e9b6c546c9bddb15d2bce
Pull request: https://github.com/HannesT117/software-factory/pull/66
Result: 2 findings

## Findings

- **Tests that miss their scenario**, `core/skills/reviewing-increments/evals/explains-deviations/graders/01-the-result-reports-the-wording-of-known.md:7`: the judge fails a correct result whose finding on the wording does not quote both lines about as line 7 writes them. With `claude plugin eval` and its default judge, a written result whose one finding is ``- **Special cases**, `cmd/greet/main.go:26`: each line of `greet --known` puts the name before "is known", while the note's Interface puts "known" before the name. A script that reads the documented line breaks.`` fails in 2 of 2 runs, with 6 of 6 votes FAIL; one whose finding says ``the line of `greet --known` has the wrong word order: `Ada is known` instead of `known Ada`, as the note's Interface fixes it`` fails in 2 of 2 runs, with 5 of 6 votes FAIL. Scenario 04.2 asks only that the wrong message is a finding in the review result. The grader of 41d7f16 took the note's wording with any name and without `greet: `, and cross-review on 7b240d8 named such forms valid ways to meet the scenario. A run like this counts against a skill that works, in the evidence of a demonstrated scenario. Doubtful: both results that the eval's reviewer wrote with this grader quote both lines with `greet: `, and four other forms pass in every run: `<name> is known` against `known <name>`, `greet: Ada is known` against `greet: known Ada`, one side in prose, and the format string of the code. With line 7 replaced by ``PASS if the section `## Findings` of this review result reports, as a finding of its own, that the line that cmd/greet/main.go prints for `greet --known` is worded differently from the line that the note fixes: the code prints `greet: <name> is known`, the note fixes `greet: known <name>`. The finding may quote the lines with any name or with `<name>`, with or without `greet: `, or describe the difference in its own words.`` and line 8 unchanged, the judge accepts all 14 correct results and rejects all 9 broken results of Checks run, in one run each.
- **Docs out of format**, `evidence/04-design-questions/evals.md:115`: the statement "The other failed grader in runs 1 and 2 is judged on the final answer" is false for run 2. There grader 01 is the only failed grader; graders 02 to 05 pass, which gives its 0.8. Only run 1 fails another grader, grader 05, because its final answer has a Mermaid sketch. The sentence explains the scores of the table above it, so the section does not describe the runs as they are.
- **Special cases:** none; the grader names the fixture's file and wordings, as a grader of an eval case does.
- **Test detection:** none.
- **Work outside the scenarios:** none; the branch changes only grader 01, its record in `evidence/04-design-questions/evals.md` and this result.
- **Protected paths:** none; against `origin/main` (2461855) the branch changes no file of the tests group, of the CI and config group, and no design note.
- **Failing checks:** none from this change. The required job `checks` fails in `check-tests` with 7 errors, the missing tests of 05.1 to 05.6 and the test of 02.7; `origin/main` (2461855) alone gives the same 7 errors.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 3 checked`
- `factory check-tests --base origin/main`: `check-tests: 7 errors`; in a worktree of `origin/main` alone, `check-tests: 7 errors`
- `factory hook stop` on this branch before this result: blocks and names 01d1869fbb816b32297e9b6c546c9bddb15d2bce
- `git merge-tree --write-tree HEAD origin/main`: no conflict
- the documentation of `claude plugin eval` at code.claude.com/docs/en/plugin-evals: `{ source: file, path: <path> }` gives an `llm` judge "the contents of one file in the workspace after the run"; the judge passes a grader with at least two of three votes
- the kept result of the three runs in `evals.md`: the skill copy equals `reviewing-increments` of this branch; grader 01 passes runs 0 and 1 with 3 of 3 votes, and the judge's evidence equals the written file byte for byte; in run 2 the runner reports `path "evidence/01-greeting/review.md" does not exist` and counts grader 01 as failed, and no other grader fails; in run 1 grader 05 fails on the Mermaid sketch of the final answer
- the trace of run 2: after `touch` fails with "Read-only file system", the reviewer gives its result in the final answer and calls no Write
- `claude plugin eval` 2.1.321 with grader 01 of this commit and the default judge, on 23 written results built for this review and placed by a scaffold, with a prompt that uses no tool, two runs each: the 9 broken results fail in every run (a `none` line with a remark that quotes both lines, and one as `None. Note: ...`; a finding on the sort that quotes both lines; a line `**Remark, not a finding**` under `Result: no findings`; a finding withdrawn on a second reading; a remark line next to a finding on the sort; the wording only under `## Demonstrated scenarios`; the wording inside a finding on `strings.Fields`; an empty file). Of the 14 correct results, 12 pass in every run (the path `main.go:26`; numbered; wrapped onto further lines; nested under a per-kind line; under `###` per-kind headings; with all seven per-kind lines; one finding that also names the sort; marked doubtful; and the four forms of the first finding); the 2 results of the first finding fail
- the same 23 results with the PASS sentence of the first finding: the 9 broken results fail, the 14 correct ones pass, one run each
- CI on b186bdd: `test` and `tests-start-red` pass; `checks` fails in `check-tests` with the 7 errors above; `cross-review` waits for a review result for 01d1869

Each `factory` command ran with a binary built from the branch; `check-tests` also with one built from `origin/main`.

## Demonstrated scenarios

- **scenario 04.2**: `evidence/04-design-questions/evals.md`, section The wording as a finding of its own; the kept result and traces of its three runs, checked above, show grader 01 passing both runs that wrote a result with the finding and failing the run that wrote none. Not run again as a whole case: `reviewing-increments` is unchanged on this branch, and the change is grader 01, which the runner judged above on 23 built results. The earlier runs of the case are in `evidence/04-design-questions/review-runs.md`, Sixth review and Seventh review, and in `evals.md`.
