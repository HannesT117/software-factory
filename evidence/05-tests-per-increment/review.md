# Review of increment 05: tests per increment

Reviewed commit: 16e4808119f0ab47fd9adabffe583e1f021d587b
Pull request: https://github.com/HannesT117/software-factory/pull/69
Result: 0 blocking, 0 notes

## Findings

None. The two notes of the round on e30dda0 are answered: the workflow patch lists the added tests with `--no-renames`, so `scenario-06.15.txtar` counts as added and its listing as unchanged behaviour is read (`.github/workflows/test.yml:73` after the patch); rule 27 in `docs/principles.md:44` says that a missing test fails only the pull requests of its own increment.

## Checks run

Each `factory` command ran with a binary built from 16e4808, unless the line says otherwise. `simbase` is `origin/main` with the design notes of #67 (5af365e) and #68 (12c6316) merged; it is not an ancestor of the branch, so its merge base with the branch is 5af365e, and the note of 06 counts as changed against it. The note of 05 is the one of #67; the code of 05 did not change since e30dda0.

- `go vet ./...`: no output
- `go test -count=1 -v ./test/acceptance/`: `ok  	github.com/HannesT117/software-factory/test/acceptance`, 47 scripts pass
- `factory check-tests --base simbase`: `check-tests: ok 47 checked`
- `factory check-tests --base origin/main`: `check-tests: 16 errors`: `unknown-id` for 05.7 and the 13 tests of 06, `missing-test` for 02.26 and 05.6, because the notes of #67 and #68 are not on `origin/main`
- `factory check-design docs/increments/05-tests-per-increment`: `check-design: ok 8 checked`
- `factory check-design --base simbase`: `check-design: ok 15 checked`; `--base origin/main`: `check-design: ok 23 checked`
- `factory check-pr --base simbase`: `check-pr: 46 errors`, `mixed-pr` on every changed path, with `.github/CODEOWNERS`, as the engineer allowed, and with the note of 06 and `docs/research/2026-11-factory-review.md` from the merge base above; `--base origin/main`: `check-pr: 47 errors`. Neither prints `tests end at`.
- `factory check-review --base simbase` with these results committed: `check-review: reviewed 16e4808119f0ab47fd9adabffe583e1f021d587b in evidence/05-tests-per-increment/review.md, round 2`, the same line for `evidence/06-gate-and-status/review.md`, and `check-review: 1 errors` for the blocking finding of increment 06
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory hook stop` before these results: blocks and names 16e4808119f0ab47fd9adabffe583e1f021d587b
- the workflow patch of the first comment of the pull request (SHA-256 `a68f361b68ad6da9a3ad68a154278b713b4eb1df2048ba6d3a7abb9332cc17b3`) applies on 16e4808
- `git merge-tree --write-tree HEAD origin/main`: no conflict

## Demonstrated scenarios

- **scenario 05.8**: `evidence/05-tests-per-increment/tests-start-red.md`. The Python of the job `tests-start-red` from the patched `.github/workflows/test.yml`, run with `BASE=simbase` on the results of dbe5e8db3403affb12c89cee5d6388c6b805321f, where the tests end, with the 19 added tests that `--no-renames` lists: `tests-start-red: ok 47 checked`, with the tests of 05.4, 05.5 and 06.15 passing before the code and the 16 other added tests failing. On the same results with the test of 05.4 set to failed: `tests-start-red: error TestAcceptance/05-tests-per-increment/scenario-05.4: fail, want pass` and `tests-start-red: 1 errors`, exit 1. Against `origin/main`, without the lists, 05.4, 05.5 and 06.15 get `pass, want fail`. On this pull request the job finds no end commit, skips the check and passes, so CI does not run this scenario here.
