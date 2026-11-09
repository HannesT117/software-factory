# Review of increment 04: design questions

Reviewed commit: 86c07b7b83f0eb0f29c6a8f95a132be3a04fc698
Pull request: https://github.com/HannesT117/software-factory/pull/66
Result: 2 findings

## Findings

- **Tests that miss their scenario**, `core/skills/reviewing-increments/evals/explains-deviations/graders/01-the-result-reports-the-wording-of-known.md:4`: the grader passes a result in which the wrong wording is no finding, when a `- **` line under `## Findings` quotes it as a remark. Scenario 04.2 asks that "the wrong message is also a finding in the review result". The regular expression checks only that one line under `## Findings` starts with `- **` and holds `cmd/greet/main.go`, `is known` and `greet: known <` or `greet: known A`; it does not check that the line reports a finding. A result that names each kind with its findings or `none`, the form that scenario 04.5 asks for and that the review results of this repo use, matches with `Result: no findings` and the line ``- **Special cases:** none; `cmd/greet/main.go:26` prints `greet: <name> is known`, the note fixes `greet: known <name>`; only an open detail, see the deviations.`` So does a result whose one finding is the sort and whose `- **Docs out of format:** none` line carries the same remark. This is the case of the cross-review, "none" with a remark that quotes the output, in the per-kind form; the grader's text on line 7 calls such a line "a finding line". Doubtful: the remark must quote both wordings, and the results that the eval's reviewer has written so far use `- **<kind>**, ...` lines for findings only. With `node`, a negative lookahead after `\\n- \*\*` that rejects a bold label followed by `none` rejects both results and keeps the seven results of the recorded check as they are.
- **Tests that miss their scenario**, `core/skills/reviewing-increments/evals/explains-deviations/graders/01-the-result-reports-the-wording-of-known.md:4`: the grader fails results in which the wrong wording is a finding, when the finding does not quote the note's wording in one of two forms on one top-level line. It needs `greet: known ` followed by `<` or `A` on the same `- **` line as `cmd/greet/main.go` and `is known`. With `node`, these results with a finding on the wording do not match: the note quoted as `known <name>`, without the prefix; the note's line for the first name of `names.txt`, `greet: known Grace`, because the class `[<A]` admits `<name>` and `Ada` but not the fixture's other names Grace and Linus; a finding whose sentence goes on in a second, indented line; the finding as a nested bullet under a per-kind line `- **Special cases:**`; and a numbered finding `1. **...**`. Scenario 04.2 asks only that the wrong message is a finding. Such a run counts against a skill that works, in the evidence of a demonstrated scenario. Doubtful: the template `assets/review.md` puts a finding on one `- **<kind>**` line, and all four findings on the wording in the kept traces, two written with Write and two given only in a final answer, quote `greet: known <name>` on such a line. Run 1 of the sixth review numbers its findings in the final answer, not in the written result.
- **Special cases:** none; the grader names the fixture's file and wording, as a grader of an eval case does; the class `[<A]` is part of the second finding.
- **Test detection:** none.
- **Work outside the scenarios:** none; the branch changes only grader 01 and its record in `evidence/04-design-questions/evals.md`.
- **Docs out of format:** none. The new section of `evals.md` states what the check showed: rerun with `node` on its `pattern.txt` and `cases.json`, all seven results match as its table says, and the pattern equals line 4 of the grader. The two kept inputs equal the JSON of the Write in their traces, and the CLI 2.1.321 tests `new RegExp(input_match)` without flags against the JSON of the input, as the section says.
- **Protected paths:** none; against `origin/main` (2461855) the branch changes no file under `test/`, `.github/` or `.factory/` and no design note.
- **Failing checks:** none from this change. The required job `checks` fails in `check-tests` with 7 errors, the missing tests of 05.1 to 05.6 and the test of 02.7; `origin/main` (2461855) alone gives the same 7 errors.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 2 checked`
- `factory check-tests --base origin/main`: `check-tests: 7 errors`; in a worktree of `origin/main` alone, `check-tests: 7 errors`
- `factory hook stop` on this branch before this result: blocks and names 86c07b7b83f0eb0f29c6a8f95a132be3a04fc698
- `git merge-tree --write-tree HEAD origin/main`: no conflict
- `node` with the check of `evals.md` (`pattern.txt`, `cases.json`): 7 of 7 as its table says; the pattern at 2461855 matches the two built results of the cross-review, `none` with a remark and the sort only
- `node` on every Write, Edit or Bash call that names a `review.md` in the kept traces of the eval runs: two Writes of `evidence/01-greeting/review.md` for explains-deviations, run 0 of the earlier round and run 1 of the sixth review; both match the old and the new pattern
- `node` on 8 results built for this review: 2 results without the finding match (first finding), 5 results with the finding do not match (second finding), 1 result with the finding on a per-kind line matches
- CI on e7c874e: `test` and `tests-start-red` pass; `checks` fails in `check-tests` with the 7 errors above; `cross-review` waits for a review result for 86c07b7

Each `factory` command ran with a binary built from the branch; `check-tests` also with one built from `origin/main`.

## Demonstrated scenarios

- **scenario 04.2**: `evidence/04-design-questions/evals.md`, section The wording as a finding of its own, rerun as above; it shows that the new grader rejects the two results that cross-review named and accepts both results that the reviewer wrote in the kept runs. Not run again with `claude plugin eval`: `reviewing-increments` is unchanged on this branch, and the change is a regular expression that the runner applies to the JSON of the Write, which the checks above run on the same input. The earlier runs of the case are in `evidence/04-design-questions/review-runs.md`, Sixth review and Seventh review, and in `evals.md`.
