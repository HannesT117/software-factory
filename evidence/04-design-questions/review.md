# Review of increment 04: design questions

Reviewed commit: bc098374c4b0a6dc2b0edd02005256cec29cad38
Pull request: https://github.com/HannesT117/software-factory/pull/66
Result: no findings

## Findings

- **Special cases:** none; the grader names the fixture's file and wordings, as a grader of an eval case does.
- **Test detection:** none.
- **Work outside the scenarios:** none; since the review of 01d1869 the branch changes only line 7 of grader 01 and its record in `evidence/04-design-questions/evals.md`.
- **Tests that miss their scenario:** none. Line 7 of `core/skills/reviewing-increments/evals/explains-deviations/graders/01-the-result-reports-the-wording-of-known.md` is now the PASS sentence that the review of 01d1869 proposed, byte for byte, and the rest of the file is unchanged; the judge accepts the two correct findings in other words that the earlier sentence failed, and still rejects the 9 broken results of that review.
- **Docs out of format:** none. `evidence/04-design-questions/evals.md:115` now says that run 1 also fails grader 05 and that run 2 fails only grader 01, as the kept result shows; the new lines on the first wording, the 23 built results and the two written results match the results checked below.
- **Protected paths:** none; against `origin/main` (2461855) the branch changes no file of the tests group, of the CI and config group, and no design note.
- **Failing checks:** none from this change. The required job `checks` fails in `check-tests` with 7 errors, the missing tests of 05.1 to 05.6 and the test of 02.7; `origin/main` (2461855) alone gives the same 7 errors.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 3 checked`
- `factory check-tests --base origin/main`: `check-tests: 7 errors`, the same as in the review of 01d1869 on `origin/main` alone
- `factory hook stop` on this branch before this result: blocks and names bc098374c4b0a6dc2b0edd02005256cec29cad38
- `git merge-tree --write-tree HEAD origin/main`: no conflict
- `cmp` of line 7 of grader 01 with the PASS sentence that the review of 01d1869 checked: equal; the other lines equal those of 01d1869
- `claude plugin eval` 2.1.321 with that sentence and the default judge, on the 23 results built for the review of 01d1869, one run each: the 9 broken results fail, the 14 correct ones pass
- `claude plugin eval` 2.1.321 with grader 01 of this commit and the default judge, on the review results that runs 0 and 1 of `evals.md` wrote, equal byte for byte to the kept files, placed by a scaffold, with a prompt that uses no tool: both pass in 2 of 2 runs, with 3 of 3 votes
- the kept result of the three runs in `evals.md`: run 1 fails grader 05 on the Mermaid sketch of its final answer, and run 2 fails only grader 01
- CI on ffe45c5: `test` and `tests-start-red` pass; `checks` fails in `check-tests` with the 7 errors above; `cross-review` waits for a review result for bc09837

Each `factory` command ran with a binary built from the branch.

## Demonstrated scenarios

- **scenario 04.2**: `evidence/04-design-questions/evals.md`, section The wording as a finding of its own; with grader 01 of this commit, the two results that the eval's reviewer wrote pass, the built broken results fail, and a run without a written result fails. Not run again as a whole case: `reviewing-increments` is unchanged on this branch, and the change is line 7 of grader 01, which the runner judged above on the written and the built results. The earlier runs of the case are in `evidence/04-design-questions/review-runs.md`, Sixth review and Seventh review, and in `evals.md`.
