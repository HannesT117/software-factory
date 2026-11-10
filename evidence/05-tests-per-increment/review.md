# Review of increment 05: tests per increment

Reviewed commit: e30dda011b555807d757e4a6753d5402b888d486
Pull request: https://github.com/HannesT117/software-factory/pull/69
Result: 0 blocking, 2 notes

## Findings

- **note** Special cases, `.github/workflows/test.yml:73` after the patch in the first comment of the pull request: `git diff --name-only --diff-filter=A` detects renames, so an added test file that git pairs with a deleted one does not count as added. On this branch `scenario-06.15.txtar` is a rename (`R051`) of the deleted `scenario-02.26.txtar`; the job then wants it to pass before the code, and the list under `## Unchanged behaviour` that names it is never read for it. The outcome is the same here, because the note of 06 lists 06.15. A later pull request that moves a test file and rewrites it gets no check that the new test fails before the code, which scenario 05.8 asks for each added test that no note lists. Doubtful: the line comes from increment 02 and the patch keeps it; `--no-renames`, which every git call of the Go code passes, would close it.
- **note** Docs out of format, `docs/principles.md:44`: rule 27 changes in this pull request, and its "Enforced by" still says that `factory check-tests --base main` fails while a design note on `main` has no tests. Since this increment it fails only in the pull requests of that increment and names the others as `waiting`. Rule 36 asks docs to describe the current state, and rule 27 is where a reader looks up what holds the order of the two pull requests.

## Checks run

Each `factory` command ran with a binary built from the branch, whose code equals e30dda0, unless the line says otherwise. `simbase` is `origin/main` with the design notes of #67 and #68 merged; it is not an ancestor of the branch, so its merge base with the branch is 5af365e, the head of #67, and the note of 06 counts as changed against it.

- `go vet ./...`: no output
- `go test $(go list ./... | grep -v /test/acceptance)`: no package outside the acceptance tests has test files
- `go test -count=1 -v ./test/acceptance/`: `ok  	github.com/HannesT117/software-factory/test/acceptance`, 47 scripts pass
- `factory check-tests --base simbase`: `check-tests: ok 47 checked`
- `factory check-tests --base origin/main`: `check-tests: 16 errors`: `unknown-id` for 05.7 and the 13 tests of 06, `missing-test` for 02.26 and 05.6, because the notes of #67 and #68 are not on `origin/main`
- `factory check-design docs/increments/05-tests-per-increment`: `check-design: ok 8 checked`
- `factory check-design --base simbase`: `check-design: ok 15 checked`, the note of 06; `--base origin/main`: `check-design: ok 23 checked`, both notes
- `factory check-pr --base simbase`: `check-pr: 46 errors`, `mixed-pr` on every changed path, with `.github/CODEOWNERS`, as the engineer allowed, and with the note of 06 and `docs/research/2026-11-factory-review.md` from the merge base above; `--base origin/main`: `check-pr: 47 errors`, with both notes. Neither prints `tests end at`.
- `factory check-review --base simbase` with this result committed: `check-review: reviewed e30dda011b555807d757e4a6753d5402b888d486 in evidence/05-tests-per-increment/review.md, round 1` and `check-review: ok 1 checked`; it does not read `evidence/06-gate-and-status/review.md` (finding of increment 06)
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory hook stop` before this result: blocks and names e30dda011b555807d757e4a6753d5402b888d486
- with `factory` built from `simbase`, as CI builds it from the base: `check-tests --base simbase`: `check-tests: ok 47 checked`; `check-adr`: `check-adr: ok 4 checked`; `check-design --base simbase` ends with a usage error, and `check-review` is no command yet, as the Rollout of the note of 06 says
- `git merge-tree --write-tree HEAD origin/main` and `HEAD simbase`: no conflict

## Demonstrated scenarios

- **scenario 05.8**: `evidence/05-tests-per-increment/tests-start-red.md`. The Python of the job `tests-start-red` from the patched `.github/workflows/test.yml`, run with `BASE=simbase` on the results of dbe5e8db3403affb12c89cee5d6388c6b805321f, where the tests end: `tests-start-red: ok 47 checked`, with the tests of 05.4 and 05.5 passing before the code and the 16 other added tests failing. On the same results with the test of 05.4 set to failed: `tests-start-red: error TestAcceptance/05-tests-per-increment/scenario-05.4: fail, want pass` and `tests-start-red: 1 errors`, exit 1. Against `origin/main`, without the lists, 05.4 and 05.5 get `pass, want fail`. `check-pr` names the end commit only in a scratch clone with the notes on the base and without `.github/CODEOWNERS` in the range; on this pull request the job finds no end commit, skips the check and passes, so CI does not run this scenario here.
