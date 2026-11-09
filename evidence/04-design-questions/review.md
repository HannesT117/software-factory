# Review of increment 04: design questions

Reviewed commit: 41d7f16a6660adc3e9e22f2ac53aedd6b138920f
Pull request: https://github.com/HannesT117/software-factory/pull/66
Result: 3 findings

## Findings

- **Tests that miss their scenario**, `core/skills/reviewing-increments/evals/explains-deviations/graders/01-the-result-reports-the-wording-of-known.md:4`: the grader still passes a result in which the wrong wording is no finding, when a per-kind line says none in another spelling and carries a remark that quotes both wordings. The negative lookahead skips a line only when its bold label is followed, after an optional colon or comma, by the lowercase word `none`. With `node`, the broken result of the earlier review matches again with each of these lines in place of its none line: `- **Special cases:** None; ...`, ``- **Special cases:** `none`; ...``, `- **Special cases:** *none*; ...`, `- **Special cases** — none; ...`, `- **Special cases: none**; ...` and `- **Special cases:** no findings; ...`, each next to a finding on the sort, and the line with `None` also under `Result: no findings`. Scenario 04.2 asks that the wrong message is a finding, and line 7 says that the grader takes no "line that says `none`". Doubtful for the reasons of the earlier review: the eval's reviewer has written no per-kind line in any kept result, and the remark must quote both wordings. With `node`, wanting the file right after the label, as the template puts it, rejects all seven results: `` \\n- \*\*[^*]*\*\*[:,]?\s*`?(?:cmd/greet/)?main\.go `` in place of `\\n- \*\*(?![^*]*\*\*[:,]?\s*none)(?=(?:(?!\\n).)*cmd/greet/main\.go)`. It gives the 15 results of the earlier checks the same answers and accepts the result of the next finding; it fails a finding with text between the label and the file, such as ``- **Special cases** (doubtful), `cmd/greet/main.go:26`: ...``, which the template does not have. Two broken results in the template's form match both patterns, and no pattern on one line tells them from a finding on the wording: ``- **Remark, not a finding**, `cmd/greet/main.go:26`: ...`` with both wordings under `Result: no findings`, and a finding on the sort that quotes both wordings. The section Limits of `evals.md` names neither.
- **Tests that miss their scenario**, `core/skills/reviewing-increments/evals/explains-deviations/graders/01-the-result-reports-the-wording-of-known.md:4`: the grader fails a correct result in the template's form when its finding names the file as `main.go:26`, without the folder, because the pattern needs `cmd/greet/main.go`. The eval's reviewer shortens the path of this file: in run 1 of the section Every grader on the middle of a trace on tool calls of `evals.md`, its final answer names `main.go:15` next to its finding on `cmd/greet/main.go:24`. Such a run counts against a skill that works, in the evidence of a demonstrated scenario. Doubtful: every finding on the wording in the kept traces names `cmd/greet/main.go`. The pattern in the first finding accepts the short path.
- **Docs out of format**, `evidence/04-design-questions/evals.md:115`: the statement "every finding on the wording in the kept traces does so", that is, puts it on one `- **<kind>**` line, is false, so the section does not describe the current state. In run 1 of the section Every grader on the middle of a trace on tool calls, the reviewer could not write and gave its result in the final answer, with the finding on the wording numbered: ``1. **Wrong output format**, `cmd/greet/main.go:24`: ...``. The final answers of run 0 of an earlier round and of run 1 of the sixth review number it too; only their written results put it on a `- **<kind>**` line. The statement is the section's reason for accepting that the grader fails a numbered finding, a form that the eval's reviewer has used for this finding where it gave its result. The other statements of the paragraph hold, as Checks run shows.
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
- `factory hook stop` on this branch before this result: blocks and names 41d7f16a6660adc3e9e22f2ac53aedd6b138920f
- `git merge-tree --write-tree HEAD origin/main`: no conflict
- line 4 of the grader: equals the pattern that `evals.md` reports as checked
- `node` on the 7 results of the table in `evals.md`: answers as the table says
- `node` on the 8 results of the earlier review: rejects both broken ones, accepts 3 of the 6 correct ones and fails the finding on a second line, the nested finding and the numbered finding, as `evals.md` says
- `node` on 20 results built for this review: the 7 broken results of the first finding and the 2 broken results in the template's form match; the correct result with `main.go:26` does not
- every final answer and every Write of a `review.md` in the 6 kept traces of explains-deviations with a finding on the wording: the 2 written results and 3 results given in a final answer put it on a `- **<kind>**` line; the result given in the final answer of run 1 of the section Every grader on the middle of a trace on tool calls numbers it, and so do the final answers of the 2 runs that wrote their result
- the runner of CLI 2.1.321: tests `new RegExp(input_match)` without flags on the JSON of each call's input
- CI on a5e88a7: `test` and `tests-start-red` pass; `checks` fails in `check-tests` with the 7 errors above; `cross-review` waits for a review result for 41d7f16

Each `factory` command ran with a binary built from the branch; `check-tests` also with one built from `origin/main`.

## Demonstrated scenarios

- **scenario 04.2**: `evidence/04-design-questions/evals.md`, section The wording as a finding of its own, rerun as above; it shows that the grader rejects the two results that cross-review named and the two broken results of the earlier review, and accepts both results that the eval's reviewer wrote in the kept runs. Not run again with `claude plugin eval`: `reviewing-increments` is unchanged on this branch, and the change is a regular expression that the runner applies to the JSON of the Write, which the checks above run on the same input. The earlier runs of the case are in `evidence/04-design-questions/review-runs.md`, Sixth review and Seventh review, and in `evals.md`.
